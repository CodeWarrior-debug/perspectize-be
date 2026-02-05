# Architecture

**Analysis Date:** 2026-02-04

## Pattern Overview

**Overall:** Hexagonal Architecture (Ports and Adapters)

**Key Characteristics:**
- Clean separation of domain logic from infrastructure
- Dependency injection through port interfaces
- Domain models have zero external dependencies
- Primary adapter is GraphQL (gqlgen), secondary adapters are PostgreSQL repositories and YouTube API client
- Explicit error handling with domain-specific error types

## Layers

**Domain Layer (Core):**
- Purpose: Pure business logic and domain models, completely independent of frameworks and infrastructure
- Location: `internal/core/domain/` and `internal/core/services/`
- Contains: Domain entities (`Content`, `Perspective`, `User`), constants/enums (`ContentType`, `Privacy`, `ReviewStatus`, `SortOrder`), domain errors, validation functions, pagination/filter types
- Depends on: Standard library only
- Used by: Services, adapters

**Ports Layer:**
- Purpose: Define interfaces (contracts) that adapters must implement
- Location: `internal/core/ports/repositories/` and `internal/core/ports/services/`
- Contains: `ContentRepository`, `UserRepository`, `PerspectiveRepository` interfaces; `YouTubeClient` service interface
- Depends on: Domain types
- Used by: Services, adapters

**Services Layer (Business Logic):**
- Purpose: Implement core business operations and orchestration
- Location: `internal/core/services/`
- Contains: `ContentService`, `UserService`, `PerspectiveService` with validation, error handling, and business rules
- Depends on: Domain types and port interfaces (repositories, external services)
- Used by: GraphQL resolvers

**GraphQL Adapter (Primary):**
- Purpose: HTTP entry point for the application using GraphQL
- Location: `internal/adapters/graphql/`
- Contains:
  - `resolvers/resolver.go` - Dependency injection container
  - `resolvers/schema.resolvers.go` - Query and Mutation resolver implementations
  - `resolvers/helpers.go` - Model conversion helpers
  - `generated/generated.go` - gqlgen auto-generated GraphQL execution code
  - `model/` - GraphQL model types (auto-generated from schema)
- Depends on: Services, domain types, gqlgen
- Used by: HTTP handler in main.go

**Repository Adapters (Secondary - Data Persistence):**
- Purpose: PostgreSQL database access implementations
- Location: `internal/adapters/repositories/postgres/`
- Contains:
  - `content_repository.go` - Cursor-based pagination, content CRUD
  - `user_repository.go` - User CRUD, lookups by ID/username/email
  - `perspective_repository.go` - Complex perspective queries with JSONB array handling
- Depends on: Domain types, sqlx, pq (PostgreSQL driver)
- Pattern: Each implements a port interface; handles DB type conversions (sql.NullString, pq.StringArray, JSONBArray)
- Used by: Services

**YouTube Adapter (Secondary - External API):**
- Purpose: YouTube Data API v3 integration
- Location: `internal/adapters/youtube/`
- Contains:
  - `client.go` - HTTP client for YouTube API; implements YouTubeClient port; fetches video metadata
  - `parser.go` - Utility functions: `ExtractVideoID()` parses URL formats, `ParseISO8601Duration()` converts duration format
- Depends on: Domain types, standard HTTP library
- Used by: ContentService

**Configuration Layer:**
- Purpose: Load and manage application configuration
- Location: `internal/config/config.go`
- Contains: Config structures (Server, Database, YouTube, Logging); `Load()` function that reads config.json and overrides with env vars
- Key flow: Environment variables override JSON config for secrets (`DATABASE_URL`, `DATABASE_PASSWORD`, `YOUTUBE_API_KEY`)
- Used by: main.go

**Database Utilities:**
- Purpose: Shared database connection and utilities
- Location: `pkg/database/postgres.go`
- Contains: `Connect()` - creates sqlx.DB with connection pooling; `Ping()` - health check
- Used by: main.go

**Custom GraphQL Scalars:**
- Purpose: Handle type conversion between GraphQL and Go
- Location: `pkg/graphql/intid.go`
- Contains: `IntID` scalar; `MarshalIntID()` and `UnmarshalIntID()` for string<->int conversion
- Used by: GraphQL generated code for ID fields

## Data Flow

**Content Creation from YouTube (Mutation):**

1. GraphQL mutation `createContentFromYouTube(url)` → `schema.resolvers.go:CreateContentFromYouTube()`
2. Call `ContentService.CreateFromYouTube(url, extractVideoIDFunc)`
3. Service validates: check if URL already exists via `ContentRepository.GetByURL()`
4. Extract video ID using `youtube.ExtractVideoID()` utility
5. Fetch metadata via `YouTubeClient.GetVideoMetadata(videoID)` (YouTube API adapter)
6. Parse ISO8601 duration via `youtube.ParseISO8601Duration()` utility
7. Create domain `Content` object with metadata
8. Persist via `ContentRepository.Create()` → PostgreSQL INSERT
9. Return `Content` domain object, converted to GraphQL model by `domainToModel()`
10. Return to client

**Perspective Creation (Mutation):**

1. GraphQL mutation `createPerspective(input)` → `schema.resolvers.go:CreatePerspective()`
2. Convert GraphQL input to service input type
3. Call `PerspectiveService.Create(input)`
4. Service validation:
   - Validate claim (required, max 255 chars)
   - Validate user exists via `UserRepository.GetByID()`
   - Validate all ratings in range [0, 10000] via `domain.ValidateRating()`
   - Check for duplicate claim via `PerspectiveRepository.GetByUserAndClaim()`
5. Set default privacy if not provided
6. Create domain `Perspective` object
7. Persist via `PerspectiveRepository.Create()` → PostgreSQL INSERT
8. Handle JSONB array conversion for categorized ratings (serialize to JSON strings)
9. Return persisted Perspective, converted to GraphQL model

**Perspective List Query with Pagination:**

1. GraphQL query `perspectives(first, after, sortBy, filter)` → `schema.resolvers.go:PerspectivesResolver()`
2. Call `PerspectiveService.ListPerspectives(params)`
3. Service validates pagination bounds (1-100 items)
4. Call `PerspectiveRepository.List(params)`
5. Repository builds SQL query:
   - Cursor-based pagination: decode `after` cursor (base64 format "cursor:{id}") to find starting point
   - Fetch `limit+1` rows to determine `hasNextPage` without extra query
   - Apply filters (userID, contentID, privacy)
   - Apply sort (CREATED_AT/UPDATED_AT/CLAIM, ASC/DESC)
   - Use keyset pagination (WHERE id > cursor_id) for performance
6. Convert perspectiveRow to domain Perspective objects
7. Unmarshal categorized_ratings from JSONB array
8. Return PaginatedPerspectives with cursors and optional total count

**State Management:**

- Domain objects are immutable data structures (no setters)
- Services receive input objects and return new domain objects
- Repositories receive domain objects and return database-stored versions with populated timestamps/IDs
- GraphQL resolvers convert between models and domain types at boundaries
- No shared state between requests; each request has its own service/repository instances

## Key Abstractions

**Port Interfaces (Contracts):**

**`ContentRepository` interface** (`internal/core/ports/repositories/content_repository.go`)
- Purpose: Define contract for content persistence layer
- Methods:
  - `Create(ctx, content)` - Insert new content, return with generated ID/timestamps
  - `GetByID(ctx, id)` - Retrieve single content
  - `GetByURL(ctx, url)` - Check for duplicate content by URL
  - `List(ctx, params)` - Cursor-based pagination with filtering and sorting
- Implemented by: `postgres.ContentRepository` in `internal/adapters/repositories/postgres/content_repository.go`

**`UserRepository` interface** (`internal/core/ports/repositories/user_repository.go`)
- Purpose: Define contract for user persistence
- Methods:
  - `Create(ctx, user)` - Insert new user
  - `GetByID(ctx, id)` - Retrieve user by ID
  - `GetByUsername(ctx, username)` - Retrieve user by username
  - `GetByEmail(ctx, email)` - Retrieve user by email
- Implemented by: `postgres.UserRepository` in `internal/adapters/repositories/postgres/user_repository.go`

**`PerspectiveRepository` interface** (`internal/core/ports/repositories/perspective_repository.go`)
- Purpose: Define contract for perspective persistence with complex queries
- Methods:
  - `Create(ctx, perspective)` - Insert new perspective with JSONB array handling
  - `GetByID(ctx, id)` - Retrieve single perspective
  - `GetByUserAndClaim(ctx, userID, claim)` - Check for duplicate claim
  - `Update(ctx, perspective)` - Update existing perspective
  - `Delete(ctx, id)` - Remove perspective
  - `List(ctx, params)` - Cursor-based pagination with user/content/privacy filters
- Implemented by: `postgres.PerspectiveRepository` in `internal/adapters/repositories/postgres/perspective_repository.go`

**`YouTubeClient` interface** (`internal/core/ports/services/youtube_client.go`)
- Purpose: Define contract for YouTube API interactions
- Methods:
  - `GetVideoMetadata(ctx, videoID) → VideoMetadata, error` - Fetch title, description, duration, channel name
- Implemented by: `youtube.Client` in `internal/adapters/youtube/client.go`

**Domain Error Types:**

Located in `internal/core/domain/errors.go`:
- `ErrNotFound` - Resource does not exist
- `ErrAlreadyExists` - Resource already exists (duplicate)
- `ErrInvalidInput` - Validation failed
- `ErrInvalidURL` - Malformed YouTube URL
- `ErrYouTubeAPI` - External API error
- `ErrInvalidRating` - Rating outside [0, 10000] range
- `ErrDuplicateClaim` - User already created this claim

Used throughout services for explicit error handling:
```go
if errors.Is(err, domain.ErrNotFound) {
    // Handle not found
}
```

**Pagination Abstractions:**

**Cursor-based pagination** (keyset pagination):
- Location: `internal/core/domain/pagination.go` (types), `internal/adapters/repositories/postgres/` (implementation)
- Cursor encoding: Base64-encoded string in format "cursor:{id}"
- SQL pattern: `WHERE id > $cursor AND ... ORDER BY created_at DESC LIMIT limit+1`
- Avoids OFFSET performance issues; supports forward traversal
- Determines `hasNextPage` by fetching `limit+1` rows and checking length

**Type Conversion Abstractions:**

**Custom scalar handling:**
- `IntID` scalar in `pkg/graphql/intid.go` - Converts GraphQL string IDs to Go int
- Marshalers: `MarshalIntID()` for Go→GraphQL, `UnmarshalIntID()` for GraphQL→Go
- Eliminates manual strconv.Atoi calls in resolvers

**Database type converters:**
- Content enum: `contentTypeToDBValue()`, `contentTypeFromDBValue()` - Converts domain.ContentType (UPPERCASE) to lowercase for DB
- Perspective privacy: `privacyToDBValue()` (implicit pattern) - Stores as lowercase in DB
- Null handling: `toNullString()`, `toNullInt64()` - Convert pointers to sql.Null* types
- Array handling: `pq.StringArray` for labels, `pq.Int64Array` for parts, `JSONBArray` for categorized_ratings

**Model Conversion:**

- Location: `internal/adapters/graphql/resolvers/helpers.go`
- Pattern: `domainToModel()`, `userDomainToModel()`, `perspectiveDomainToModel()` functions
- Purpose: Convert domain objects to GraphQL model types before returning to client
- Domain types have zero GraphQL knowledge; conversion happens at adapter boundary

## Entry Points

**HTTP Server Entry Point:**
- Location: `cmd/server/main.go`
- Triggers:
  1. Load .env file (optional via godotenv)
  2. Load config from config/config.example.json, override with env vars
  3. Connect to PostgreSQL via `database.Connect(dsn)`
  4. Test connection with `database.Ping()` and version query
  5. Initialize adapters (YouTube client, repositories)
  6. Initialize services with adapter dependencies
  7. Initialize GraphQL resolver with service dependencies
  8. Create gqlgen handler: `handler.NewDefaultServer(generated.NewExecutableSchema())`
  9. Register HTTP routes:
     - `GET /` → GraphQL Playground
     - `POST /graphql` → GraphQL execution
  10. Listen on configured port (default 8080)

**GraphQL Request Flow:**
- Client sends mutation/query → HTTP POST /graphql
- gqlgen parses request, validates against schema
- Routes to appropriate resolver in `schema.resolvers.go`
- Resolver calls service methods
- Service calls repository/external service methods via port interfaces
- Results propagated back through layers, converted to GraphQL models at boundary

## Error Handling

**Strategy:** Explicit error returns with domain error types; no exceptions

**Patterns:**

**Service layer:**
```go
// Domain error
func (s *UserService) Create(ctx context.Context, username, email string) (*domain.User, error) {
    if username == "" {
        return nil, fmt.Errorf("%w: username is required", domain.ErrInvalidInput)
    }
    // ...
}

// Check if error is specific domain error
existing, err := s.repo.GetByUsername(ctx, username)
if err != nil && !errors.Is(err, domain.ErrNotFound) {
    return nil, fmt.Errorf("failed to check username: %w", err)
}
```

**Repository layer:**
```go
// Map database "not found" to domain error
if errors.Is(err, sql.ErrNoRows) {
    return nil, domain.ErrNotFound
}
return nil, fmt.Errorf("failed to get user: %w", err)
```

**Resolver layer (GraphQL boundary):**
```go
user, err := r.UserService.Create(ctx, input.Username, input.Email)
if err != nil {
    if errors.Is(err, domain.ErrAlreadyExists) {
        return nil, fmt.Errorf("user already exists: %w", err)
    }
    if errors.Is(err, domain.ErrInvalidInput) {
        return nil, fmt.Errorf("invalid input: %w", err)
    }
    return nil, fmt.Errorf("failed to create user: %w", err)
}
```

**Context propagation:**
- All repository and service methods accept `context.Context` as first parameter
- Enables request cancellation, deadlines, and request-scoped values
- Database operations respect context timeouts

## Cross-Cutting Concerns

**Logging:** Not yet implemented. Placeholder: use `log/slog` (per CLAUDE.md).

**Validation:**
- **Input validation:** Done in service layer before database operations
  - Username: required, max 24 chars
  - Email: required, valid format (regex), unique
  - Claim: required, max 255 chars, unique per user
  - Ratings: 0-10000 range (validated per field)
- **Type validation:** Handled by domain constants/enums (ContentType, Privacy, etc.)
- **Database constraints:** Enforce NOT NULL, UNIQUE, CHECK constraints at DB level

**Authentication:** Not yet implemented. No auth middleware or user session handling in current codebase.

**Rate Limiting:** Not implemented.

**Database Connection Pooling:**
- Location: `pkg/database/postgres.go`
- Settings:
  - MaxOpenConns: 25
  - MaxIdleConns: 5
  - ConnMaxLifetime: 5 minutes

---

*Architecture analysis: 2026-02-04*
