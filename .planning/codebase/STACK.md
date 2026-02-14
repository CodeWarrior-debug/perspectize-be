# Technology Stack

**Analysis Date:** 2026-02-13

## Languages

**Primary:**
- Go 1.25.0 - Backend GraphQL API, server, database adapters, business logic
- TypeScript 5.9.3 - Frontend SvelteKit application, type-safe components and utilities
- SQL - PostgreSQL database migrations and queries
- GraphQL - API schema definition and code generation

**Supporting:**
- CSS - Tailwind v4 styling framework with global theme tokens
- HTML - SvelteKit template structure

## Runtime

**Environment:**
- Go 1.25.0 - Backend compiled binary execution
- Node.js (unspecified minor version) - Frontend development and build tooling
- PostgreSQL 18 - Primary data store (via Docker or external managed service)

**Package Manager:**
- `go mod` for backend dependencies
- `pnpm` for frontend dependencies (lockfile: `pnpm-lock.yaml`)

## Frameworks

**Backend:**
- gqlgen 0.17.86 - Schema-first GraphQL code generation and server
- chi/v5 2.5.5 - HTTP router and middleware
- pgx/v5 5.7.6 - PostgreSQL driver with connection pooling
- GORM 1.31.1 - Object-relational mapping (ORM migration in progress from sqlx)
  - gorm.io/gorm 1.31.1 - ORM core
  - gorm.io/driver/postgres 1.6.0 - PostgreSQL GORM driver

**Frontend:**
- SvelteKit 2.50.1 - Full-stack web framework with file-based routing
- Svelte 5.48.2 - Reactive UI components (Svelte 5 runes only, no Svelte 4 patterns)
- TanStack Query (SvelteKit) 6.0.18 - Server state management and data fetching
- Vite 6.4.1 - Build tool and dev server
- Tailwind CSS 4.1.18 - Utility-first CSS framework with `--color-*` theme variables

**UI Components:**
- shadcn-svelte - Headless component primitives (via bits-ui 2.15.5)
- AG Grid 32.2.x (pinned to v32, via ag-grid-svelte5 1.0.3 wrapper)
  - @ag-grid-community/core 32.2.1
  - @ag-grid-community/client-side-row-model 32.2.1
  - @ag-grid-community/theming 32.2.0

**Testing:**
- Vitest 4.0.18 - Frontend test runner with jsdom environment
- @testing-library/svelte 5.3.1 - Component testing utilities
- @testing-library/jest-dom 6.9.1 - DOM assertions

**Icons:**
- @lucide/svelte 0.561.0 - Icon library (per-icon imports for tree-shaking)

**Forms:**
- @tanstack/svelte-form 1.28.0 - Form state management

**Utilities:**
- graphql 16.12.0 - GraphQL client utilities
- graphql-request 7.4.0 - Minimal GraphQL client for TanStack Query integration
- clsx 2.1.1 - Conditional className utility
- tailwind-merge 3.4.0 - Tailwind CSS class merging
- svelte-sonner 1.0.7 - Toast notifications
- @internationalized/date 3.11.0 - Date handling utilities

## Key Dependencies

**Critical Backend:**
- gqlgen 0.17.86 - Schema-first GraphQL code generation with type binding and resolver generation
- pgx/v5 5.7.6 - PostgreSQL wire protocol driver with connection pooling and prepared statements
- GORM 1.31.1 - ORM for query building, dynamic filtering, and migration support (current migration from sqlx)
- chi/v5 2.5.5 - HTTP router and middleware stack
- testify 1.11.1 - Testing assertions and mocking utilities

**Critical Frontend:**
- graphql 16.12.0 - GraphQL client utilities and language support
- graphql-request 7.4.0 - Lightweight GraphQL client for TanStack Query
- @tanstack/svelte-query 6.0.18 - Server state management with caching and synchronization
- @tanstack/svelte-form 1.28.0 - Form state management
- @ag-grid-community/core 32.2.1 - AG Grid core functionality (pinned v32 for stability)

**Development/Build:**
- air 1.64.4 - Hot reload for backend development (watches `.go` and `.graphql` files)
- golangci-lint - Backend linting (invoked via `make lint`)
- golang-migrate - Schema migration tool (invoked via `make migrate-*` commands)
- svelte-check 4.3.5 - Svelte static type checking
- @sveltejs/adapter-auto 7.0.0 - SvelteKit deployment adapter
- @tailwindcss/vite 4.1.18 - Tailwind CSS Vite plugin integration

## Configuration

**Environment Variables:**
- `DATABASE_URL` - PostgreSQL connection string (required, format: `postgres://user:pass@host:5432/db`)
- `YOUTUBE_API_KEY` - Optional YouTube API v3 credentials
- `SEVALLA_BACKEND_URL` - Reference URL for Sevalla deployments
- `VITE_GRAPHQL_URL` - Frontend GraphQL endpoint (defaults to `http://localhost:8080/graphql`)

**Config Files:**
- `backend/gqlgen.yml` - GraphQL code generation mappings, scalar definitions, enum bindings
- `backend/.air.toml` - Air hot-reload watcher configuration (watches `.go`, `.graphql` files)
- `backend/.env` and `backend/.env.example` - Environment variable templates
- `frontend/vite.config.ts` - Vite build configuration with SvelteKit plugin, Tailwind integration, Vitest setup
- `frontend/tsconfig.json` - TypeScript strict mode configuration with bundler module resolution
- `backend/Makefile` - Build, test, migration, and development commands

**Configuration Source Priority:**
1. Environment variables (highest priority, overrides all)
2. `.env` file (if present)
3. Hardcoded defaults in config packages

## Platform Requirements

**Development:**
- Go 1.25.0 (or compatible)
- PostgreSQL 18 (via Docker Compose or local installation)
- Docker and Docker Compose
- Node.js (any recent version) with pnpm package manager
- migrate CLI (golang-migrate) for schema management
- golangci-lint for code linting
- GNU Make for build commands

**Build:**
- `go build` for backend binary compilation
- `vite build` for frontend static asset generation
- `go run github.com/99designs/gqlgen generate` for GraphQL code generation

**Production Deployment:**
- PostgreSQL 18 (external managed service or self-hosted)
- Go compiled binary execution environment
- Node.js (optional, if using Node.js hosting for frontend)
- SvelteKit-compatible hosting (Vercel, Netlify, static file server, or Node.js server)
- Environment variables for `DATABASE_URL`, `YOUTUBE_API_KEY`, `VITE_GRAPHQL_URL`

**Special Deployment Notes:**
- Sevalla/Fly.io: PostgreSQL connections may require `?sslmode=disable` and may succeed on second attempt
- Cold starts: ~10-50ms for Go binary
- Binary size: Typically 10-20MB when compiled
- Memory footprint: ~20-50MB (lean compared to Node.js ~100-300MB)

## Database

**Engine:** PostgreSQL 18 (Docker image: `postgres:18`)

**Drivers:**
- pgx/v5 5.7.6 (native PostgreSQL wire protocol, high-performance)
- GORM ORM layer (active migration from sqlx)

**Connection Details:**
- Uses `DATABASE_URL` environment variable for connection string
- Docker service: `perspectize-postgres-go` on port 5432
- Local development defaults: user `testuser`, password `testpass`, database `testdb`
- Connection pooling: Max open 25, max idle 5, max lifetime 5 minutes (hardcoded in `pkg/database/postgres.go`)

**Migrations:**
- Tool: golang-migrate (via `make migrate-*` commands)
- Location: `backend/migrations/`
- Commands: `make migrate-up`, `make migrate-down`, `make migrate-create`, `make migrate-version`, `make migrate-force`

## Build & Development Tools

**Backend:**
- Makefile targets: `build`, `run`, `dev` (hot-reload), `test`, `test-coverage`, `graphql-gen`, `fmt`, `lint`, `docker-up`, `docker-down`, `clean`
- Docker Compose for local PostgreSQL
- Air v1.64.4 hot-reload (watches `go` and `graphql` files, excludes `_test.go`)
- golang-migrate for versioned SQL migrations

**Frontend:**
- Vite dev server (`pnpm run dev` on http://localhost:5173)
- Vitest test runner with jsdom environment, globals enabled
- TypeScript strict checking (`pnpm run check`)
- Code duplication detection (`pnpm run test:duplication`)
- Coverage: v8 provider with thresholds (lines 80%, functions 80%, branches 75%, statements 80%)

## Docker & Containerization

**Docker Compose (`backend/docker-compose.yml`):**
- Service: PostgreSQL 18
- Container name: `perspectize-postgres-go`
- Port: 5432
- Health check: `pg_isready` with 10s interval, 5s timeout, 5 retries
- Volume: Named volume `postgres_data` for persistence

**Commands:**
- `make docker-up` - Start PostgreSQL container (with 3s wait)
- `make docker-down` - Stop container
- `make docker-logs` - View PostgreSQL logs

## Testing Infrastructure

**Frontend (Vitest 4.0.18):**
- Environment: jsdom (browser DOM simulation)
- Globals: true (no need to import `describe`, `it`, `expect`)
- Test patterns: `src/**/*.{test,spec}.{js,ts}` and `tests/**/*.{test,spec}.{js,ts}`
- Setup files: `tests/setup.ts` (Vitest configuration)
- Coverage provider: v8
- Coverage thresholds: 80% lines/functions/statements, 75% branches
- Coverage exclusions: `node_modules`, `.svelte-kit`, `*.d.ts`, `*.config.*`, `src/lib/components/shadcn/**`, `src/routes/**`

**Backend:**
- Runner: `go test ./...`
- Coverage command: `make test-coverage` (outputs `coverage.html`)
- Integration tests: Auto-skip when database unavailable (`t.Skip()`)
- Mocking: testify with custom mock implementations in `test/` directory

## GraphQL

**Schema:** Schema-first approach in `backend/schema.graphql`

**Code Generation:**
- Tool: gqlgen 0.17.86
- Command: `make graphql-gen`
- Configuration: `backend/gqlgen.yml` defines:
  - Schema files
  - Output paths for generated code
  - Model bindings (custom scalars like `IntID`, domain enums)
  - Resolver layout and naming conventions
  - Type mappings (e.g., `JSON` → `graphql.Map`)

**Generated Code Locations:**
- Executable schema: `internal/adapters/graphql/generated/generated.go`
- Models: `internal/adapters/graphql/model/models_gen.go`
- Resolver stubs: `internal/adapters/graphql/resolvers/{name}.resolvers.go`

**GraphQL Endpoints:**
- Query/Mutation: `POST /graphql`
- Playground (dev only): `GET /` (interactive IDE)

---

*Stack analysis: 2026-02-13*
