# Coding Conventions

**Analysis Date:** 2026-02-04

## Naming Patterns

**Files:**
- Lowercase with underscores: `user_service.go`, `user_repository.go`, `content_repository.go`
- Domain models: `user.go`, `content.go`, `perspective.go`, `pagination.go`, `errors.go`
- Test files: `{subject}_test.go` placed in `test/{category}/` (e.g., `test/services/user_service_test.go`)
- Generated files: prefixed with context like `models_gen.go`, `generated.go`, `schema.resolvers.go`

**Functions:**
- Public functions: PascalCase starting with verb or noun (e.g., `Create()`, `GetByID()`, `NewUserService()`)
- Private functions: camelCase (e.g., `userRowToDomain()`, `contentTypeToDBValue()`)
- Receivers: single lowercase letter (e.g., `s *UserService`, `r *UserRepository`, `m *mockUserRepository`)
- Constructor functions: `New{Type}` pattern (e.g., `NewUserService()`, `NewUserRepository()`)

**Variables:**
- Local variables: camelCase (e.g., `userID`, `existing`, `username`, `createdAt`)
- Constants: SCREAMING_SNAKE_CASE or PascalCase based on context
  - Domain enums: UPPERCASE (e.g., `SortOrderAsc`, `ContentTypeYouTube`, `ContentSortByName`)
  - Error variables: PascalCase with `Err` prefix (e.g., `ErrNotFound`, `ErrAlreadyExists`)
- Pointers for optional fields: `*int`, `*string` (e.g., `URL *string`, `Length *int`)
- Slice for collections: `[]int`, `[]string` (e.g., `Parts []int`, `Labels []string`)

**Types:**
- Structs: PascalCase (e.g., `UserRepository`, `ContentService`, `PaginatedContent`)
- Enums: PascalCase type name with UPPERCASE constants (e.g., `type ContentType string; const ContentTypeYouTube ContentType = "YOUTUBE"`)
- Input types: `{Name}Input` pattern (e.g., `CreatePerspectiveInput`)
- Database mapping types: lowercase with suffix (e.g., `userRow`, `contentRow`)
- Mock implementations: `mock{Name}` pattern (e.g., `mockUserRepository`, `mockContentRepository`)

## Code Style

**Formatting:**
- Uses `gofmt` (enforced by `make fmt`)
- All Go files must be `gofmt`-compliant
- Lines can exceed 80 characters when necessary (e.g., SQL queries, long error messages)

**Linting:**
- Uses `golangci-lint` (run with `make lint`)
- Config not explicitly checked in (uses default rules)
- Errors must be wrapped with `fmt.Errorf("message: %w", err)` format

**Whitespace:**
- Tab indentation (Go standard)
- Blank lines separate logical sections within functions
- Imports grouped with blank lines: standard library, third-party, local

## Import Organization

**Order:**
1. Standard library imports (e.g., `context`, `fmt`, `errors`)
2. Third-party imports (e.g., `github.com/stretchr/testify`, `github.com/jmoiron/sqlx`)
3. Local imports from `github.com/yourorg/perspectize-go/` (domain, ports, adapters, config)

**Path Aliases:**
- No custom aliases used in codebase
- Long import paths not aliased (full paths used)
- `portservices` import alias used for `internal/core/ports/services` to avoid confusion with `internal/core/services`

## Error Handling

**Patterns:**
- Domain layer defines error variables: `var ErrNotFound = errors.New("resource not found")` (file: `internal/core/domain/errors.go`)
- Service/adapter layers wrap errors: `fmt.Errorf("context message: %w", err)` to preserve error chain
- Use `errors.Is(err, domain.ErrNotFound)` to check domain errors
- SQL `sql.ErrNoRows` is converted to domain `ErrNotFound` at repository level
- Do NOT use custom error types; use sentinel errors (defined in `domain/errors.go`)

**Error Translation:**
- Domain errors have "business meaning" (e.g., `ErrInvalidInput`, `ErrAlreadyExists`)
- Services wrap domain errors with contextual messages: `fmt.Errorf("%w: username already taken", domain.ErrAlreadyExists)`
- Repositories wrap database errors with operation context: `fmt.Errorf("failed to get user by id: %w", err)`
- Resolvers translate domain errors to GraphQL responses (error translation may differ per resolver)

**Validation Pattern:**
- Input validation in service methods, not repositories
- Validate constraints at service layer: `if len(username) > 24 { return nil, fmt.Errorf("%w: ...", domain.ErrInvalidInput) }`
- Return early on validation failure

## Logging

**Framework:** Standard library `log` package (not `slog`)
- File: `cmd/server/main.go` uses `log.Println()` and `log.Printf()` for startup/connection messages
- No structured logging yet in code (CLAUDE.md documents intent to use `slog`, but not yet applied)
- Logging is sparse: only critical operations (DB connection, version checks)

**Patterns:**
- Log at startup: configuration loading, database connection, version queries
- Log connection success with operational context (e.g., "Connecting to database at host:port/db...")
- Mask credentials in logs: "Connecting to database using DATABASE_URL..." (no printed values)
- Use `log.Fatalf()` for fatal errors that prevent startup
- Use `log.Println()` and `log.Printf()` for informational messages

## Comments

**When to Comment:**
- Function comments: Exported functions have comment starting with function name (e.g., `// Create inserts a new user record...`)
- Type comments: Structs have comment describing purpose (e.g., `// User represents a user who can create perspectives`)
- Complex logic: Inline comments explain non-obvious decisions (e.g., "Fetch limit+1 to determine hasNextPage")
- Conversions: Comments explain why conversions are needed (e.g., `// userRowToDomain converts a database row to a domain User`)
- Database queries: Comments label query steps (e.g., "Check if username already exists", "Validate email")

**JSDoc/TSDoc:**
- Not used (Go uses `//` comments only)
- Function signatures are self-documenting with type information

## Function Design

**Size:**
- Small focused functions (typical 20-50 lines)
- Services contain validation logic and orchestration (50-100 lines common for complex operations like `Create`)
- Repositories contain query logic (30-80 lines typical)
- Helper functions extracted for reuse (e.g., `contentTypeToDBValue()`, `userRowToDomain()`)

**Parameters:**
- Accept `context.Context` as first parameter in all service/repository methods
- Use `*Domain` pointers for input entities
- Use `interface{}` for dependencies (e.g., `repositories.UserRepository`, `portservices.YouTubeClient`)
- Optional inputs: pointer types for nullable fields (e.g., `First *int`, `After *string`)

**Return Values:**
- Single or dual returns
- Dual pattern: `(*Type, error)` - return nil resource on error
- Never return both resource and error set
- Error always checked with `if err != nil { return ... }`

## Module Design

**Exports:**
- All public types and functions exported (PascalCase)
- Private types/functions lowercase (only used within package)
- Constructor functions (`New{Type}`) always exported

**Barrel Files:**
- Not used (no index.go or barrel exports)
- Each package imports what it needs from other packages

**Dependency Injection:**
- Via constructor functions: `NewUserService(repo repositories.UserRepository) *UserService`
- Repository injected via interface: repositories are interfaces in `internal/core/ports/repositories/`
- Services injected into resolvers via constructor: `NewResolver(contentService, userService, perspectiveService)`
- All dependencies resolved in `cmd/server/main.go`

**Package Structure:**
- Domain layer (`internal/core/domain/`): entities, enums, errors (no dependencies on adapters)
- Ports (`internal/core/ports/`): repository/service interfaces (contracts)
- Services (`internal/core/services/`): business logic (depends on ports, not adapters)
- Adapters (`internal/adapters/`): implementations (depends on domain, not domain dependent)
- Config: loaded via `internal/config/`, environment variables override JSON

## Domain Enum Pattern

**Requirement:** Always use gqlgen model binding - never write switch statements to convert between GraphQL and domain enums.

**For value enums (SortOrder, Privacy, ContentType, etc.):**
1. Define domain enums with UPPERCASE values in `internal/core/domain/`:
   ```go
   type SortOrder string
   const (
       SortOrderAsc  SortOrder = "ASC"
       SortOrderDesc SortOrder = "DESC"
   )
   ```

2. Bind in `gqlgen.yml` to use Go types directly (eliminates manual mapping)

3. For DB-stored enums, add repository converter functions:
   ```go
   func contentTypeToDBValue(ct domain.ContentType) string {
       return strings.ToLower(string(ct))
   }
   func contentTypeFromDBValue(s string) domain.ContentType {
       return domain.ContentType(strings.ToUpper(s))
   }
   ```

**For ID fields in filters/inputs:**
Use the `IntID` custom scalar (`pkg/graphql/intid.go`) instead of `ID` with manual `strconv.Atoi`. This auto-converts string IDs to integers in filters.

## Pagination & Cursors

**Cursor Pattern:**
- Opaque base64-encoded strings: `base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("cursor:%d", id)))`
- Encoded format: `cursor:<id>` before base64 encoding
- Decoded in repository: validate format, extract ID, convert to int
- Helper functions: `encodeCursor()`, `decodeCursor()` in repository

**SQL Pagination:**
- Keyset pagination: `id > cursor_id` for forward, `id < cursor_id` for backward
- Fetch `limit+1` rows to determine `hasNextPage` without extra query
- Whitelist sort columns via switch function: `sortColumnName()` prevents SQL injection

## Optional Fields Pattern

**Nullable fields:** Use pointers for optional values
```go
type Perspective struct {
    Claim   string  // Required - always has value
    Quality *int    // Optional - nil means "not provided"
}

// Check if optional field is set
if p.Quality != nil {
    fmt.Println(*p.Quality)  // Dereference to get value
}

// Set an optional field
quality := 85
p.Quality = &quality
```

**Database mapping:** `sql.NullString`, `sql.NullInt64`, `sql.NullTime` for nullable columns

---

*Convention analysis: 2026-02-04*
