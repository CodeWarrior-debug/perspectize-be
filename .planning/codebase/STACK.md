# Technology Stack

**Analysis Date:** 2026-02-04

## Languages

**Primary:**
- Go 1.25+ - Primary language for all backend development (`perspectize-go/`)
- SQL - Database schema and migrations in PostgreSQL

**Secondary (Legacy):**
- C# with ASP.NET Core - Legacy implementation in `perspectize-be/` directory (do not modify, migration in progress)

## Runtime

**Environment:**
- Go 1.25.0 binary compiled application
- Runs as standalone HTTP server on configurable port (default: 8080)
- Docker containerization for PostgreSQL in development

**Package Manager:**
- Go Modules (`go.mod`, `go.sum`)
- Lockfile: Present (`go.sum`)
- Build output: Binary at `bin/perspectize-server` after `make build`

## Frameworks

**Core:**
- **gqlgen** v0.17.86 - GraphQL server framework (code generation, schema-first approach)
  - Generates executable schema from `schema.graphql`
  - Output: `internal/adapters/graphql/generated/generated.go`
  - Also generates resolver stubs and models
- **github.com/gorilla/websocket** v1.5.3 - WebSocket support (indirect dependency via gqlgen)

**Web Server:**
- Standard library `net/http` - HTTP request handling and routing
  - GraphQL endpoint: `/graphql`
  - GraphQL Playground: `/` (interactive IDE)

**Database:**
- **sqlx** v1.4.0 - SQL query builder and result mapper (wrapper around database/sql)
- **pgx** v5.7.6 - PostgreSQL driver (high-performance pure Go driver)
  - Loaded via stdlib interface in `pkg/database/postgres.go`
- **lib/pq** v1.10.9 - Alternative PostgreSQL driver (likely legacy, pgx is primary)

**Configuration:**
- **godotenv** v1.5.1 - Load `.env` files for local development

**Code Generation:**
- **gqlparser** v2.5.31 - GraphQL query/schema parser (used by gqlgen)

## Key Dependencies

**Critical:**
- **gqlgen** v0.17.86 - Why it matters: Eliminates boilerplate, enforces schema-first GraphQL development, auto-generates type-safe Go code
- **pgx/v5** v5.7.6 - Why it matters: Modern, performant PostgreSQL driver with prepared statement support; replaces legacy lib/pq
- **sqlx** v1.4.0 - Why it matters: Convenient query methods (`.Get()`, `.Select()`) and struct scanning; reduces SQL boilerplate
- **godotenv** v1.5.1 - Why it matters: Enables environment-based configuration without manual .env parsing

**Infrastructure:**
- **stretchr/testify** v1.11.1 - Testing assertions and mocking library
- **air** v1.64.4 - Hot-reload development server (`make dev` uses this)

**Indirect/Supporting:**
- **goccy/go-yaml** v1.19.2 - YAML parsing (used by gqlgen config)
- **hashicorp/golang-lru/v2** - LRU caching (dependency of gqlgen)
- **google/uuid** v1.6.0 - UUID generation (likely used in query cache keys)
- **pelletier/go-toml/v2** v2.2.4 - TOML parsing (gqlgen config support)
- **vektah/gqlparser/v2** v2.5.31 - GraphQL parsing utilities

## Configuration

**Environment:**
Configuration is loaded from two sources (in order of precedence):

1. **Environment variables** (highest priority) - Overrides config file settings for secrets
   - `DATABASE_URL` - Full PostgreSQL connection string (e.g., `postgres://user:pass@host:5432/db?sslmode=disable`)
   - `DATABASE_PASSWORD` - Optional override for just the password field
   - `YOUTUBE_API_KEY` - YouTube Data API v3 key (optional)

2. **config/config.example.json** - Non-secret application configuration
   - Default location: `perspectize-go/config/config.example.json`
   - Loaded in `cmd/server/main.go` via `config.Load()`
   - JSON structure defines server port, database host/port, logging level

**Build Configuration:**
- `.env` file (git-ignored) - Loaded automatically by Makefile and `main.go` via godotenv
- `.env.example` - Template showing required variables
- `gqlgen.yml` - GraphQL code generation configuration
  - Defines schema file location
  - Output paths for generated code
  - Model binding for custom types (`IntID`, enums, scalars)

**Makefile Variables:**
- Database URL with fallback: `DB_URL ?= $(or $(DATABASE_URL),postgres://testuser:testpass@localhost:5432/testdb?sslmode=disable)`
- Build target: `go build -o bin/perspectize-server cmd/server/main.go`

## Platform Requirements

**Development:**
- Go 1.25.0+
- Docker & Docker Compose (for PostgreSQL container)
- PostgreSQL 17 (via docker-compose.yml)
- GNU Make (for Makefile commands)
- golangci-lint (optional, for `make lint`)
- golang-migrate CLI (required for database migrations)
- air (included in `go.mod` as tool, for hot-reload development)

**Production:**
- PostgreSQL 17 database (external or managed)
- Linux/Unix environment (Go binary runs on macOS/Linux/Windows)
- Deployment target: Fly.io or Sevalla (mentioned in CLAUDE.md)
  - Special note: Sevalla connections may require `?sslmode=disable` and may succeed on second attempt
  - Cold starts ~10-50ms (important for serverless)

**Deployment:**
- Go binary size: Typically 10-20MB when compiled
- Memory footprint: ~20-50MB (vs ~100-300MB for Node.js)
- No external runtime required beyond PostgreSQL database

## Connection Pool Configuration

**Database (`pkg/database/postgres.go`):**
- Max open connections: 25
- Max idle connections: 5
- Connection max lifetime: 5 minutes
- These settings are hardcoded in the `Connect()` function

**HTTP Server:**
- Built on Go's net/http which manages goroutines automatically
- No explicit pool configuration (relies on Go scheduler)

---

*Stack analysis: 2026-02-04*
