# Backend: Perspectize Go

Go GraphQL API built with gqlgen, PostgreSQL 18, and Hexagonal Architecture.

## Architecture

```
perspectize-go/
├── cmd/server/       # Entry point
├── internal/
│   ├── core/         # domain/ (models), ports/ (interfaces), services/ (logic)
│   ├── adapters/     # graphql/ (primary), repositories/ (DB), youtube/ (API)
│   ├── config/       # Configuration loading
│   └── middleware/    # HTTP middleware
├── pkg/              # database/ (connection), graphql/ (IntID scalar)
└── migrations/       # SQL migration files
```

Full structure: [docs/ARCHITECTURE.md](../docs/ARCHITECTURE.md)

**Dependency Rule:** Dependencies point inward. Domain never depends on adapters. Adapters depend on domain ports.

### Hexagonal Architecture Guidelines

When implementing features:

1. **Start with domain** - Define models in `core/domain/`
2. **Define ports** - Create interfaces in `core/ports/`
3. **Implement business logic** - Write services in `core/services/`
4. **Add adapters** - Implement infrastructure in `adapters/`
5. **Wire dependencies** - Connect adapters to core in `cmd/server/main.go`

**Domain layer (`core/domain/`)** contains pure Go structs with no external dependencies. See [docs/DOMAIN_GUIDE.md](../docs/DOMAIN_GUIDE.md) for domain rules, entity details, and optional fields pattern.

## Technology Stack

- **Language:** Go 1.25+
- **GraphQL:** gqlgen (code generation, schema-first)
- **Database:** PostgreSQL 18 with sqlx + pgx/v5 driver
- **Migrations:** golang-migrate
- **Validation:** go-playground/validator
- **Testing:** testing + testify + sqlmock
- **Logging:** log/slog (structured logging)
- **Environment:** godotenv (.env file loading)

## Development Commands

### Initial Setup

```bash
go mod download
make docker-up
make migrate-up
cp .env.example .env
```

### Daily Development

```bash
make run          # Server on http://localhost:8080
make dev          # Hot-reload with air
make test         # Run all tests
make test-coverage # Coverage report → coverage.html
make fmt          # Format code
make lint         # Lint code
make graphql-gen  # Regenerate GraphQL after schema changes
```

### Database Migrations

```bash
make migrate-up       # Apply pending migrations
make migrate-down     # Rollback last migration
make migrate-create   # Create new migration (prompts for name)
make migrate-version  # Check current version
make migrate-force    # Force to specific version (recovery)
```

### Docker

```bash
make docker-up    # Start PostgreSQL
make docker-down  # Stop and remove
make docker-logs  # View PostgreSQL logs
```

## Configuration

Loaded from two sources (in order of precedence):
1. **Environment variables** (highest priority) - `DATABASE_URL`, `YOUTUBE_API_KEY`
2. **config/config.json** - non-secret configuration

Required: `DATABASE_URL` (e.g., `postgres://testuser:testpass@localhost:5432/testdb?sslmode=disable`)
Optional: `YOUTUBE_API_KEY`, `DATABASE_PASSWORD`

See `.env.example` for details.

**Production (Sevalla/Fly.io):** Sevalla connections may require `?sslmode=disable` and may succeed on second attempt.

## Database

PostgreSQL 18 features used: JSONB columns, array types, custom domains, triggers (auto timestamps).

Migration files in `migrations/`: `{sequence}_{description}.{up|down}.sql`

## GraphQL Schema

Schema-first approach in `schema.graphql`. After modifying:
1. Run `make graphql-gen`
2. Implement resolver logic in `internal/adapters/graphql/resolvers/`
3. Wire resolvers to domain services

## Testing

- **Unit tests:** Mock external dependencies, no DB required. Run with `make test`.
- **Integration tests:** Guarded with `t.Skip()` when DB unavailable. Auto-skip.
- **Environment isolation:** Tests loading config must clear env vars (`DATABASE_URL`, `DATABASE_PASSWORD`, `YOUTUBE_API_KEY`) via `t.Setenv("KEY", "")`. See `clearConfigEnvVars` in `test/config/config_test.go`.

## Code Style

- Follow [Effective Go](https://go.dev/doc/effective_go)
- `gofmt` formatting (enforced by `make fmt`)
- Explicit error handling, structured logging with `slog`
- Dependency injection via interfaces (ports)

## Adding a New Feature

1. Define domain model: `internal/core/domain/feature.go`
2. Define repository port: `internal/core/ports/repositories/feature_repository.go`
3. Implement service: `internal/core/services/feature_service.go`
4. Implement repository: `internal/adapters/repositories/postgres/feature_repository.go`
5. Update schema: `schema.graphql` → `make graphql-gen`
6. Implement resolver: `internal/adapters/graphql/resolvers/feature_resolver.go`
7. Wire in `cmd/server/main.go`
8. Write tests: `test/services/`, `test/repositories/`

## Error Handling & Database Queries

See [docs/GO_PATTERNS.md](../docs/GO_PATTERNS.md) for error wrapping pattern (domain → service → resolver) and sqlx query/transaction examples.

## Gotchas

### GraphQL Schema Defaults
gqlgen passes default values (e.g., `first: Int = 10`) as non-nil pointers, not `nil`. Tests should expect the default value.

### JSON Scalar Type
Use gqlgen's built-in `graphql.Map` scalar (configured in `gqlgen.yml` as `JSON`) for JSONB data.

### Cursor-Based Pagination
- Cursors are opaque base64 strings (format: `cursor:<id>`)
- Use keyset pagination (not OFFSET)
- Fetch `limit+1` rows to determine `hasNextPage`
- Whitelist sort columns to prevent SQL injection

### GraphQL Enum & ID Handling (REQUIRED)

**Always use gqlgen model binding** — never write switch statements for enum conversion.

1. Define domain enums with UPPERCASE values:
   ```go
   type SortOrder string
   const (
       SortOrderASC  SortOrder = "ASC"
       SortOrderDESC SortOrder = "DESC"
   )
   ```

2. Bind in `gqlgen.yml`:
   ```yaml
   models:
     SortOrder:
       model:
         - github.com/CodeWarrior-debug/perspectize-be/perspectize-go/internal/core/domain.SortOrder
   ```

3. For DB-stored enums, add repository converters (lowercase ↔ UPPERCASE).

4. Use the `IntID` custom scalar (`pkg/graphql/intid.go`) instead of `ID` with manual `strconv.Atoi`.

**When adding new enums:** Add UPPERCASE constants → bind in `gqlgen.yml` → add DB converter if stored → `make graphql-gen`

## Self-Verification

```bash
curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -d '{"query": "{ __typename }"}'
# Expect: {"data":{"__typename":"Query"}}
```
