# Architecture

**Analysis Date:** 2025-02-13

## Monorepo Overview

Perspectize is a **two-stack monorepo** with active development on both frontend and backend:

```
perspectize/ (repository root)
├── backend/     # Go GraphQL API backend (ACTIVE)
├── frontend/     # SvelteKit web app frontend (ACTIVE)
└── perspectize/     # Legacy C# (DEPRECATED - ignore)
```

**Key Principle:** Frontend and backend are loosely coupled via GraphQL HTTP API. No shared code between stacks. Stacks are independently deployable.

---

## Backend Architecture (backend/)

### Pattern: Hexagonal Architecture (Ports & Adapters)

The Go backend implements **Hexagonal Architecture** with strict dependency inversion:

```
┌─────────────────────────────────────────────────────┐
│         PRIMARY ADAPTERS (Driving/Input)            │
│    GraphQL Resolvers (gqlgen schema-first)          │
│    HTTP Handlers (chi router + CORS middleware)     │
└──────────────┬──────────────────────────────────────┘
               │ depends on
┌──────────────▼──────────────────────────────────────┐
│         CORE DOMAIN LAYER (Domain Logic)            │
│  ├── domain/       (Entities, enums, errors)        │
│  ├── services/     (Business logic, orchestration)  │
│  └── ports/        (Interface contracts)            │
│                                                      │
│  Critical Rule: Dependencies point INWARD ONLY      │
│  Domain never depends on adapters or frameworks     │
└──────────────┬──────────────────────────────────────┘
               │ implements
┌──────────────▼──────────────────────────────────────┐
│      SECONDARY ADAPTERS (Driven/Output)             │
│  ├── repositories/postgres/  (Database access)      │
│  └── youtube/                (External API client)  │
└─────────────────────────────────────────────────────┘
```

**Critical:** Services depend on port interfaces, not concrete adapters. This enables testability and clean separation.

### Layers

**Domain Layer** (`backend/internal/core/domain/`):
- **Purpose:** Pure business models and rules, zero external dependencies
- **Contains:**
  - `user.go` — User entity (ID, Username, Email, CreatedAt, UpdatedAt)
  - `content.go` — Content entity (ID, Name, URL, ContentType, Length, LengthUnits, Response JSON, timestamps)
  - `perspective.go` — Perspective entity with complex structure:
    - Claim (required, max 255 chars), UserID, ContentID (optional)
    - Ratings: Quality, Agreement, Importance, Confidence (0-10000, optional)
    - Metadata: Like, Privacy (PUBLIC/PRIVATE), Description, Category, ReviewStatus
    - Arrays: Parts (integer array), Labels (text array)
    - JSONB: CategorizedRatings (array of {category, rating})
    - Enums: Privacy, ReviewStatus, PerspectiveSortBy, SortOrder
  - `pagination.go` — Pagination types (ContentListParams, PaginatedContent, PerspectiveListParams, PaginatedPerspectives, cursor support)
  - `errors.go` — Domain error constants (ErrNotFound, ErrAlreadyExists, ErrInvalidInput, ErrInvalidURL, ErrYouTubeAPI, ErrInvalidRating, ErrDuplicateClaim)
- **Depends on:** Standard library only (time, encoding/json, errors)
- **Used by:** Services, repositories (for model conversion)
- **Key Pattern:** No external dependencies, no GORM tags, no database concerns

**Ports Layer** (`backend/internal/core/ports/`):
- **Purpose:** Define contracts (interfaces) between core and adapters
- **Repositories** (`backend/internal/core/ports/repositories/`):
  - `user_repository.go` — UserRepository interface (Create, GetByID, GetByUsername, GetByEmail, ListAll)
  - `content_repository.go` — ContentRepository interface (Create, GetByID, GetByURL, List with pagination)
  - `perspective_repository.go` — PerspectiveRepository interface (Create, GetByID, GetByUserAndClaim, Update, Delete, List with complex filtering)
- **Services** (`backend/internal/core/ports/services/`):
  - `youtube_client.go` — YouTubeClient interface (GetVideoMetadata, ExtractVideoID)
  - `user_service.go` — UserService interface (Create, GetByID, GetByUsername, ListAll)
  - `content_service.go` — ContentService interface (CreateFromYouTube, GetByID, ListContent)
  - `perspective_service.go` — PerspectiveService interface (Create, GetByID, Update, Delete, List)
- **Depends on:** Domain types only
- **Key Rule:** Core services depend on these interfaces; adapters implement these interfaces

**Services Layer** (`backend/internal/core/services/`):
- **Purpose:** Business logic orchestration, validation, error handling
- **Contains:**
  - `user_service.go` — UserService implementation
    - Create: username/email validation, uniqueness check, REGEX email validation, max length constraints
    - GetByID: input validation (positive ID), delegation to repo
    - GetByUsername: trimming, existence check
    - ListAll: delegation to repo
  - `content_service.go` — ContentService implementation
    - CreateFromYouTube: URL deduplication, YouTube ID extraction, metadata fetching, domain model creation
    - GetByID: input validation, repo delegation
    - ListContent: pagination bounds validation (1-100), repo delegation
  - `perspective_service.go` — PerspectiveService implementation
    - Create: comprehensive validation (claim required/max 255, user exists, ratings in range 0-10000, no duplicate claims per user)
    - GetByID: input validation
    - Update: similar validation as Create
    - Delete: cascade handling
    - List: pagination + complex filtering (UserID, ContentID, Privacy)
- **Depends on:** Domain models + port interfaces (never concrete adapters)
- **Pattern:** Constructor injection - `NewUserService(repo repositories.UserRepository) *UserService`
- **Error handling:** Explicit error returns with `fmt.Errorf("%w: context", domainError)` to preserve error type for `errors.Is()` checks
- **Validation Location:** Services validate ALL inputs before calling repo (input validation), then repo queries decide if resource exists

**GraphQL Adapter** (`backend/internal/adapters/graphql/`):
- **Purpose:** PRIMARY adapter — HTTP GraphQL API entry point for clients
- **Structure:**
  - `resolvers/resolver.go` — Dependency injection container (holds references to services)
    - Fields: ContentService, UserService, PerspectiveService (all injected at startup)
    - Created in main.go with wired services
  - `resolvers/schema.resolvers.go` — AUTO-GENERATED resolver implementations
    - Queries: users(), userByID(id), perspectives(...), content(...)
    - Mutations: createUser(input), createContentFromYouTube(url), createPerspective(input), updatePerspective(id, input), deletePerspective(id)
    - Transform GraphQL input → service input, call service, transform service response → GraphQL model
    - Error mapping: domain.ErrXxx → user-friendly GraphQL error strings
  - `resolvers/helpers.go` — Model conversion functions (domain → GraphQL models)
    - `domainToModel(content *domain.Content) *model.Content`
    - `userDomainToModel(user *domain.User) *model.User`
    - `perspectiveToModel(p *domain.Perspective) *model.Perspective`
  - `generated/generated.go` — gqlgen-generated executable schema (DO NOT EDIT)
  - `model/` — GraphQL type definitions (auto-generated from schema.graphql)
  - `schema.graphql` — GraphQL schema definition (schema-first approach)
- **Workflow:**
  1. Edit `schema.graphql` to add/modify types, queries, mutations
  2. Run `make graphql-gen` to auto-generate gqlgen code
  3. Implement resolver logic in `schema.resolvers.go`
  4. Wire in `cmd/server/main.go`
- **Key Patterns:**
  - Resolver methods return `(*model.Type, error)`
  - Error handling: map domain errors to strings, log unexpected errors with slog
  - Input validation happens in service layer, not resolver
  - Cursor pagination: decode base64 cursors, pass to repo

**Repository Adapters** (`backend/internal/adapters/repositories/postgres/`):
- **Purpose:** SECONDARY adapter — PostgreSQL database persistence via GORM
- **GORM Models** (`gorm_models.go`):
  - `UserModel` — mirrors `users` table (ID, Username, Email, CreatedAt, UpdatedAt)
  - `ContentModel` — mirrors `content` table (ID, Name, URL, ContentType, Length, LengthUnits, Response JSONB, timestamps)
  - `PerspectiveModel` — mirrors `perspectives` table with complex types:
    - Basic fields: ID, Claim, UserID, ContentID, Like, Quality/Agreement/Importance/Confidence, Privacy, Category, Description, ReviewStatus
    - Arrays: Parts (Int64Array), Labels (StringArray)
    - JSONB: CategorizedRatings (JSONBArray), column: categorized_ratings
    - Timestamps: CreatedAt, UpdatedAt
  - Each model has `TableName()` method for GORM
- **Mappers** (`gorm_mappers.go`):
  - Bidirectional conversion functions:
    - `userDomainToModel(*domain.User) *UserModel` — prepare for INSERT
    - `userModelToDomain(*UserModel) *domain.User` — extract from query result
    - `contentDomainToModel(*domain.Content) *ContentModel`
    - `contentModelToDomain(*ContentModel) *domain.Content`
    - `perspectiveModelToDomain(*PerspectiveModel) *domain.Perspective`
    - `perspectiveDomainToModel(*domain.Perspective) *PerspectiveModel`
  - Enum conversion: lowercase DB values ↔ UPPERCASE domain enums
  - Null/pointer handling: nil checks, pointer wrapping for optional fields
- **Implementations:**
  - `gorm_user_repository.go` — UserRepository
    - Create: domain → model, INSERT, auto-fill ID/timestamps
    - GetByID: SELECT WHERE id=?, map to domain, handle gorm.ErrRecordNotFound → domain.ErrNotFound
    - GetByUsername/GetByEmail: WHERE clause + model conversion
    - ListAll: ORDER BY username ASC, collect into domain array
  - `gorm_content_repository.go` — ContentRepository
    - Create: similar to user
    - GetByID, GetByURL: similar pattern
    - List: Pagination with cursor decoding, filter params, sort order, LIMIT limit+1 for hasNextPage
  - `gorm_perspective_repository.go` — PerspectiveRepository (most complex)
    - Create: with duplicate claim check, JSONB array marshaling
    - GetByID, GetByUserAndClaim: basic queries
    - Update: UPDATE ... SET on specific fields, timestamp handling
    - Delete: DELETE WHERE id=?
    - List: Complex GORM chaining
      - Build dynamic WHERE clauses for filters (UserID, ContentID, Privacy)
      - Apply ORDER BY sort_by, sort_order
      - Pagination: cursor decoding, keyset query (id > cursor)
      - LIMIT limit+1 to determine hasNextPage
      - Convert to domain models in loop
  - `helpers.go` — Custom type handlers
    - Int64Array: custom Scanner/Valuer for pq integer arrays
    - StringArray: custom Scanner/Valuer for pq text arrays
    - JSONBArray: custom Scanner/Valuer for JSONB array columns
- **Pattern:**
  - Each repository implements its port interface (compile-time check: `var _ repositories.UserRepository = (*GormUserRepository)(nil)`)
  - GORM injected via constructor: `NewGormUserRepository(db *gorm.DB)`
  - All methods accept `ctx context.Context` for cancellation
  - Errors returned explicitly (no panic)
  - Query building lazy (GORM doesn't execute until `.Error` or `.Scan()`)

**YouTube Adapter** (`backend/internal/adapters/youtube/`):
- **Purpose:** SECONDARY adapter — External YouTube Data API v3 integration
- **Client** (`client.go`):
  - NewClient(apiKey string) — constructor with default baseURL, http.Client
  - GetVideoMetadata(ctx context.Context, videoID string) — Makes HTTP request to YouTube API
    - Endpoint: `https://www.googleapis.com/youtube/v3/videos?part=snippet,statistics,contentDetails&id=...&key=...`
    - Parses response into YouTubeAPIResponse struct
    - Returns VideoMetadata (Title, Description, Duration in seconds, ChannelName, raw Response JSON)
    - Error mapping: non-200 → domain.ErrYouTubeAPI, no items → domain.ErrNotFound
    - Duration parsing: ISO8601 format → seconds, fallback to 0 with warning
  - ExtractVideoID(url string) — delegates to parser
- **Parser** (`parser.go`):
  - ExtractVideoID(url string) — REGEX patterns for youtube.com, youtu.be variants
  - ParseISO8601Duration(duration string) — converts P0DT1H23M45S to seconds
- **Implements:** YouTubeClient port interface (testable via mock)

**Infrastructure Layer:**
- **Configuration** (`backend/internal/config/`):
  - Config struct: Server (Port, Host), Database (Host, Port, Name, User, Password, SSLMode), YouTube (APIKey), Logging (Level, Format)
  - Load(configPath string) — Load JSON file, override with env vars (DATABASE_URL, YOUTUBE_API_KEY, DATABASE_PASSWORD)
  - GetDSN() — Format PostgreSQL DSN string
  - ValidateDatabaseURL, SanitizeDSN — security/validation helpers
- **Database** (`backend/pkg/database/`):
  - PoolConfig struct: MaxOpenConns (default 25), MaxIdleConns (default 5), ConnMaxLifetime (default 5m)
  - PoolConfigFromEnv() — Read from DB_MAX_OPEN_CONNS, DB_MAX_IDLE_CONNS, DB_CONN_MAX_LIFETIME env vars
  - ConnectGORM(dsn string, pool PoolConfig) — Create sql.DB with pgx driver, configure pool, wrap with GORM
  - PingGORM(ctx context.Context, db *gorm.DB) — Test connection
- **GraphQL Scalar** (`backend/pkg/graphql/`):
  - intid.go — Custom IntID scalar for type-safe integer IDs in GraphQL inputs
- **Server Bootstrap** (`backend/cmd/server/main.go`):
  - Load .env via godotenv
  - Load config from CONFIG_PATH env var or default
  - Validate DATABASE_URL, YOUTUBE_API_KEY
  - Connect to PostgreSQL with pool config
  - Initialize adapters: YouTubeClient, GormUserRepository, GormContentRepository, GormPerspectiveRepository
  - Initialize services: UserService, ContentService, PerspectiveService (inject repos)
  - Create resolver: NewResolver(contentService, userService, perspectiveService)
  - Setup chi router with middleware:
    - middleware.RequestID — attach unique ID to each request
    - middleware.RealIP — extract real IP from X-Forwarded-For
    - middleware.Logger — log HTTP requests (uses slog)
    - middleware.Recoverer — panic recovery
    - CORS middleware (lines 115-126) — allows all origins
  - Register endpoints:
    - GET /health — liveness probe (always 200)
    - GET /ready — readiness probe (ping DB, return 503 if unavailable)
    - POST /graphql — GraphQL execution
    - GET / — GraphQL Playground (non-production only)
  - HTTP server config: ReadTimeout 15s, WriteTimeout 15s, IdleTimeout 60s
  - Graceful shutdown: SIGINT/SIGTERM → 30s context timeout for cleanup

## Data Flow

### GraphQL Mutation Flow: Create Perspective

```
Client: mutation { createPerspective(input: {...}) { id claim ... } }
  ↓
HTTP POST /graphql (chi router)
  ↓
gqlgen handler: Parse GraphQL query, lookup mutation resolver
  ↓
schema.resolvers.go: CreatePerspective(ctx, input model.CreatePerspectiveInput)
  ↓
Model conversion: model.CreatePerspectiveInput → portservices.CreatePerspectiveInput
  ↓
PerspectiveService.Create(ctx, serviceInput)
  ├─ Validate claim (not empty, ≤255 chars)
  ├─ Validate ratings (nil or 0-10000)
  ├─ UserRepository.GetByID(ctx, userID) — verify user exists
  ├─ PerspectiveRepository.GetByUserAndClaim(ctx, userID, claim) — check duplicate
  ├─ Construct domain.Perspective
  └─ PerspectiveRepository.Create(ctx, perspective)
      ↓
      gorm_perspective_repository.go: Create
      ├─ perspectiveDomainToModel(perspective) → PerspectiveModel
      ├─ db.WithContext(ctx).Create(model) — GORM INSERT
      ├─ GORM auto-fills ID, CreatedAt, UpdatedAt
      └─ perspectiveModelToDomain(model) → domain.Perspective
      ↓
      (back up call stack)
  ↓
Resolver: perspectiveToModel(domainPerspective) → model.Perspective
  ↓
gqlgen: Serialize to JSON
  ↓
HTTP 200 JSON response to client
```

### Query Flow: List Perspectives with Pagination

```
Client: query { perspectives(first: 10, after: "cursor123") { items { ... } pageInfo { ... } } }
  ↓
HTTP POST /graphql
  ↓
schema.resolvers.go: Perspectives(ctx, input model.PerspectiveListInput)
  ↓
Model conversion: model.PerspectiveListInput → domain.PerspectiveListParams
  ├─ Decode after/before cursors from base64
  ├─ Build PerspectiveFilter struct from input
  └─ Set SortBy, SortOrder
  ↓
PerspectiveService.List(ctx, params)
  ├─ Validate pagination bounds (1-100 for first/last)
  └─ PerspectiveRepository.List(ctx, params)
      ↓
      gorm_perspective_repository.go: List
      ├─ Start GORM query chain: db.WithContext(ctx)
      ├─ Add WHERE filters:
      │  ├─ WHERE user_id = ? (if params.Filter.UserID set)
      │  ├─ WHERE content_id = ? (if params.Filter.ContentID set)
      │  └─ WHERE privacy = ? (if params.Filter.Privacy set)
      ├─ Apply cursor filter: WHERE id > cursor_id (keyset pagination)
      ├─ ORDER BY sort_by, sort_order
      ├─ LIMIT first+1 (to determine hasNextPage)
      ├─ Execute query: .Find(&models)
      ├─ Loop: perspectiveModelToDomain(models[i])
      ├─ Build PaginatedPerspectives:
      │  ├─ Items: domain perspectives (truncate to limit if over)
      │  ├─ HasNext: len(models) > first
      │  ├─ StartCursor, EndCursor: base64-encode first/last IDs
      │  └─ TotalCount: optional COUNT(*) if requested
      └─ Return PaginatedPerspectives
  ↓
Resolver: Convert domain.PaginatedPerspectives → model.PerspectiveConnection
  ├─ Items → Edges (with cursors)
  ├─ PageInfo (hasNextPage, startCursor, endCursor)
  └─ TotalCount (if present)
  ↓
gqlgen: Serialize to JSON
  ↓
HTTP 200 JSON response to client
```

### YouTube Content Creation Flow

```
Client: mutation { createContentFromYouTube(url: "https://youtube.com/watch?v=abc123") { ... } }
  ↓
schema.resolvers.go: CreateContentFromYouTube(ctx, input)
  ↓
ContentService.CreateFromYouTube(ctx, url)
  ├─ ContentRepository.GetByURL(ctx, url) — check if already exists
  ├─ YouTubeClient.ExtractVideoID(url) — parse URL to get video ID
  ├─ YouTubeClient.GetVideoMetadata(ctx, videoID)
  │  ↓
  │  YouTube Adapter (youtube/client.go)
  │  ├─ Format API endpoint with video ID and API key
  │  ├─ HTTP GET request to googleapis.com
  │  ├─ Parse JSON response into YouTubeAPIResponse
  │  ├─ Extract: Title, Description, Duration (ISO8601 → seconds), ChannelName
  │  └─ Return VideoMetadata { Title, Description, Duration, ChannelName, Response (raw JSON) }
  │  ↓
  ├─ Construct domain.Content with YouTube metadata
  └─ ContentRepository.Create(ctx, content)
      ↓
      gorm_content_repository.go
      ├─ contentDomainToModel(content) → ContentModel
      ├─ db.Create(model) — INSERT, GORM auto-fills ID and timestamps
      └─ contentModelToDomain(model) → domain.Content
      ↓
  ↓
Resolver: domainToModel(content) → model.Content
  ↓
gqlgen: Serialize to JSON
  ↓
HTTP 200 JSON response to client
```

### State Management

- **Domain Models:** Immutable; created in services, passed to repos
- **Service Layer:** Stateless; validates, orchestrates, delegates to repos
- **Repository Layer:** Stateless; query builders with GORM lazy evaluation
- **HTTP Server:** Graceful shutdown with 30s timeout for in-flight requests
- **Database Connection:** Connection pool managed by GORM (pgx driver)
- **Middleware Context:** Request ID attached via chi middleware, flows through all layers

## Key Abstractions

**Repository Pattern:**
- **Purpose:** Isolate database access logic, enable testing via mocks
- **Examples:** `backend/internal/core/ports/repositories/user_repository.go`, `content_repository.go`, `perspective_repository.go`
- **Pattern:**
  - Port interface defines contract (CRUD operations + specialized queries)
  - Implementation in `adapters/repositories/postgres/` uses GORM
  - Testable: mock repository implementation with predictable returns
  - Swappable: replace GORM with another ORM or database

**Service Ports:**
- **Purpose:** Isolate external integrations and service implementations
- **Examples:** `backend/internal/core/ports/services/youtube_client.go`
- **Pattern:**
  - Port interface defines contract (methods for external integration)
  - Adapter (`backend/internal/adapters/youtube/client.go`) implements via HTTP
  - Testable: mock YouTube client for unit tests, real client for integration tests
  - Swappable: replace HTTP client with gRPC client or database-backed service

**GORM Model Mapping (ORM Separation):**
- **Purpose:** Decouple domain models from database representation
- **Files:**
  - `backend/internal/adapters/repositories/postgres/gorm_models.go` — GORM models with `gorm:` tags
  - `backend/internal/adapters/repositories/postgres/gorm_mappers.go` — conversion functions
- **Pattern:**
  - Domain models: Pure Go structs, no ORM tags (user.go, content.go, perspective.go)
  - GORM models: Database-specific structs with GORM column/type/constraint tags (UserModel, ContentModel, PerspectiveModel)
  - Mappers: Bidirectional conversion (userDomainToModel, userModelToDomain)
  - Benefit: Domain can evolve independently of DB schema; easy to change ORM

**GraphQL Model Adaptation:**
- **Purpose:** Decouple GraphQL API contracts from domain/database models
- **Files:**
  - `backend/internal/adapters/graphql/schema.graphql` — Schema definition (source of truth)
  - `backend/internal/adapters/graphql/model/` — Auto-generated GraphQL types
  - `backend/internal/adapters/graphql/resolvers/helpers.go` — Conversion functions
- **Pattern:**
  - Schema-first: GraphQL schema defined in schema.graphql
  - Code generation: `make graphql-gen` auto-generates model types
  - Resolvers: Implement by converting domain → GraphQL models
  - Benefit: Schema changes don't ripple through domain; easy to support multiple API versions

**Dependency Injection:**
- **Purpose:** Wire components at startup, enable testing without global state
- **Pattern in main.go (lines 91-102):**
  ```go
  youtubeClient := youtube.NewClient(cfg.YouTube.APIKey)
  contentRepo := postgres.NewGormContentRepository(db)
  userRepo := postgres.NewGormUserRepository(db)
  perspectiveRepo := postgres.NewGormPerspectiveRepository(db)

  contentService := services.NewContentService(contentRepo, youtubeClient)
  userService := services.NewUserService(userRepo)
  perspectiveService := services.NewPerspectiveService(perspectiveRepo, userRepo)

  resolver := resolvers.NewResolver(contentService, userService, perspectiveService)
  ```
- **Pattern:** Constructor injection (pass dependencies as args, no global service locator)
- **Benefit:** Testable (inject mocks), explicit dependencies, easy to trace call chains

**Cursor-Based Pagination:**
- **Purpose:** Stable pagination across concurrent database modifications (unlike OFFSET-based pagination which can skip/duplicate rows)
- **Files:**
  - `backend/internal/core/domain/pagination.go` — PaginatedContent, PerspectiveListParams, cursors
  - `backend/internal/adapters/repositories/postgres/gorm_content_repository.go` — List implementation
  - `backend/internal/adapters/repositories/postgres/gorm_perspective_repository.go` — Complex List with filters
- **Pattern:**
  - Cursor = base64-encoded ID (opaque, allows changing encoding)
  - Keyset query: `WHERE id > :cursor_id ORDER BY created_at DESC LIMIT limit+1`
  - Fetch limit+1 to determine hasNextPage
  - Return startCursor (first item ID) and endCursor (last item ID)
- **Benefit:** Stable across concurrent writes, scales to large datasets

## Entry Points

**HTTP Server Bootstrap** (`backend/cmd/server/main.go` lines 28-182):
- Loads .env file
- Loads config (JSON file + env var overrides)
- Validates DATABASE_URL, YOUTUBE_API_KEY
- Connects to PostgreSQL with pool configuration
- Creates adapters (YouTube client, repositories)
- Creates services (inject adapters)
- Creates GraphQL resolver (inject services)
- Registers HTTP routes with chi router
- Sets up middleware (request ID, logging, CORS, panic recovery)
- Listens on configured port (default 8080)
- Graceful shutdown on SIGINT/SIGTERM

**HTTP Endpoints:**
- `POST /graphql` — GraphQL mutations and queries
- `GET /` — GraphQL Playground (non-production only)
- `GET /health` — Liveness probe
- `GET /ready` — Readiness probe

**GraphQL Resolvers** (from `schema.graphql`):
- **Queries:**
  - `users() → [User!]!` — List all users
  - `userByID(id: ID!) → User` — Single user by ID
  - `content(...) → ContentConnection!` — List content with pagination
  - `perspectives(...) → PerspectiveConnection!` — List perspectives with complex filtering
- **Mutations:**
  - `createUser(input: CreateUserInput!) → User!` — Create user
  - `createContentFromYouTube(url: String!) → Content!` — Create content from YouTube URL
  - `createPerspective(input: CreatePerspectiveInput!) → Perspective!` — Create perspective
  - `updatePerspective(id: ID!, input: UpdatePerspectiveInput!) → Perspective!` — Update perspective
  - `deletePerspective(id: ID!) → Boolean!` — Delete perspective

## Error Handling

**Strategy:** Domain-defined error types, wrapped with context, mapped to GraphQL errors in resolvers

**Domain Errors** (`backend/internal/core/domain/errors.go`):
- `ErrNotFound` — Resource doesn't exist
- `ErrAlreadyExists` — Duplicate resource
- `ErrInvalidInput` — Validation failure
- `ErrInvalidURL` — Bad URL format
- `ErrYouTubeAPI` — YouTube API failure
- `ErrInvalidRating` — Rating out of range
- `ErrDuplicateClaim` — User already has this claim

**Service Layer Wrapping:**
- Services wrap domain errors with context using `fmt.Errorf("%w: context", domainError)`
- Wrapping preserves error type for `errors.Is()` checks
- Example from user_service.go:
  ```go
  if len(username) > 24 {
      return nil, fmt.Errorf("%w: username must be 24 characters or less", domain.ErrInvalidInput)
  }
  ```

**Resolver Error Handling:**
- Resolvers check error type and return user-friendly messages
- Expected errors (validation) return GraphQL error strings
- Unexpected errors logged with slog and return generic message
- Example from schema.resolvers.go:
  ```go
  if errors.Is(err, domain.ErrAlreadyExists) {
      return nil, fmt.Errorf("user already exists: %w", err)
  }
  slog.Error("creating user failed", "error", err)
  return nil, fmt.Errorf("failed to create user")
  ```

**Repository Error Handling:**
- Repositories map database errors to domain errors
- Example: `gorm.ErrRecordNotFound` → `domain.ErrNotFound`
- Other database errors returned wrapped with context

**Context Propagation:**
- Context passed through all layers (main → chi → resolver → service → repo)
- Enables cancellation and timeout support
- HTTP server timeouts: 15s read/write, 60s idle
- Graceful shutdown: 30s context timeout

## Cross-Cutting Concerns

**Logging:**
- **Framework:** `log/slog` (structured logging, standard library)
- **Usage:**
  - `slog.Info()` — Configuration loaded, server started
  - `slog.Warn()` — Missing optional config (YOUTUBE_API_KEY)
  - `slog.Error()` — Unexpected errors in resolvers, startup failures
- **Location:** `backend/cmd/server/main.go` (lines 57-83 for startup logs)
- **Middleware:** `chi/middleware.Logger` — Automatic HTTP request logging

**Validation:**
- **Location:** Service layer (before database operations)
- **Patterns:**
  - Input validation: UserService validates username/email format, length, uniqueness
  - Business rule validation: PerspectiveService validates ratings (0-10000), duplicate claims
  - Constraint validation: ContentService validates YouTube URL format
  - Database constraints: UNIQUE, NOT NULL, CHECK constraints enforce at DB level
- **Error mapping:** Return domain errors for validation failures

**Authentication:**
- **Current State:** Not implemented (open API for MVP)
- **Future:** Phase 5 will add OAuth 2.0 with YouTube, JWT tokens, refresh token rotation

**CORS:**
- **Current:** Middleware in main.go allows all origins (`*`)
- **Future:** Phase 5 will restrict to frontend production origin
- **Middleware:** Lines 115-126 in main.go

**Database Transactions:**
- **Current:** Single repository operations, no multi-repo transactions
- **Pattern:** Each Create/Update/Delete is atomic within GORM
- **Future:** May need transaction support for complex operations (e.g., create user + perspective simultaneously)

**Request Scoping:**
- **Context:** Attached to requests by chi middleware, flows through all layers
- **Request ID:** `chi/middleware.RequestID` — unique ID per request for tracing
- **Real IP:** `chi/middleware.RealIP` — extract X-Forwarded-For for load-balanced scenarios
- **Timeouts:** Server-level (15s read/write), graceful shutdown (30s)

---

*Architecture analysis: 2025-02-13*
