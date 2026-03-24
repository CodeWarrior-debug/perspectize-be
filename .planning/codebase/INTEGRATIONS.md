# External Integrations

**Analysis Date:** 2026-02-16

## APIs & External Services

**YouTube Data API v3:**
- **Purpose:** Fetch video metadata (title, description, duration, view counts) for content enrichment
- **SDK/Client:** Custom HTTP client implementation in `internal/adapters/youtube/client.go`
- **Auth:** Environment variable `YOUTUBE_API_KEY` (set via `.env` or config.json)
- **Endpoint:** `https://www.googleapis.com/youtube/v3/videos`
- **Status:** Active integration with functional client implementation
  - Fetches: `snippet`, `statistics`, `contentDetails` parts
  - Extracts: title, description, channelTitle, publishedAt, tags, duration, viewCount, likeCount, commentCount
  - Returns: Structured `VideoMetadata` to domain services via port interface `services.VideoMetadata`
  - Parser: `internal/adapters/youtube/parser.go` handles URL parsing for YouTube links in multiple formats (youtube.com, youtu.be, shorts, live, etc.)
- **Error Handling:** Returns custom domain error `domain.ErrYouTubeAPI` on API failures, `domain.ErrNotFound` if video not found
- **Response Storage:** Full API response stored as JSONB in `content.response` database column for audit trail
- **Config:** YouTube API key optional but warnings logged if missing (graceful degradation)

## Data Storage

**Databases:**

**PostgreSQL 18:**
- **Connection:** Via environment variable `DATABASE_URL` or config file + `DATABASE_PASSWORD` env var
- **Client:** GORM 1.31.1 ORM with pgx/v5 driver
- **Connection Details:**
  - Default local: `postgres://testuser:testpass@localhost:5432/testdb?sslmode=disable`
  - Docker container: `perspectize-postgres-go` (via docker-compose.yml)
  - Port: 5432
  - Sevalla deployment: May require `?sslmode=disable`
- **Advanced Features Used:**
  - JSONB columns: `content.response` stores YouTube API responses
  - Array types: PostgreSQL arrays for tags, labels, parts
  - Custom domains: Type validation at database level
  - Triggers: Automatic timestamp updates
- **Connection Pool:**
  - Max open: 25
  - Max idle: 5
  - Max lifetime: 5 minutes
  - Configurable via environment (pkg/database/pool_config.go)
- **Slow Query Logging:** Queries >100ms logged via `database.RegisterSlowQueryLogger()`
- **Migrations:** Located in `backend/migrations/`
  - Tool: golang-migrate
  - Format: `{sequence}_{description}.{up|down}.sql`
  - Current migrations (6 total):
    - 000001_create_content.sql - Initial content table
    - 000002_update_response_jsonb.sql - Change response to JSONB
    - 000003_update_length_numeric.sql - Numeric length columns
    - 000004_add_perspectives_users.sql - Perspectives and users tables
    - 000005_add_user_timestamps.sql - User created/updated timestamps

**File Storage:**
- Not applicable - Local filesystem only for migrations and config files
- All user data stored in PostgreSQL
- Static assets (images, icons) served from `frontend/static/`

**Caching:**
- Application-level: TanStack Query on frontend with `@tanstack/svelte-query` v6.0.18
- Server-side: None configured (direct DB queries via GORM)
- No Redis, Memcached, or in-memory cache layer

## Authentication & Identity

**Auth Provider:**
- Custom implementation (no external identity provider configured)
- Current approach: User system in domain with ID-based references
- `internal/core/domain/` defines user entity with roles (ADMIN, SENTINEL, DEFAULT)
- `internal/adapters/repositories/postgres/gorm_user_repository.go` manages persistence
- GraphQL mutations: `createUser`, `updateUser`, `deleteUser`
- GraphQL queries: `userByID`, `userByUsername`, `users`

**Status:** Basic user CRUD operations, no OAuth/JWT/session management currently implemented

## Monitoring & Observability

**Error Tracking:**
- Not configured - No Sentry, Rollbar, or error aggregation service
- Errors logged via structured JSON logging

**Logs:**
- **Framework:** Go standard library `log/slog` (structured JSON logging)
- **Configured in:** `pkg/logger/` setup for Sevalla log viewer compatibility
- **Frontend logging:** Browser console for GraphQL endpoint configuration and debug messages
- **Backend startup logs:** Database version, connection health, config validation, slow query detection
- **Environment-aware:** Suppresses `.env` file warnings in production (APP_ENV=production)

## CI/CD & Deployment

**Hosting:**
- Target platform: Sevalla (custom hosting service)
- Supports environment variables at build time (for VITE_GRAPHQL_URL)
- Database connectivity via DATABASE_URL to external PostgreSQL endpoint
- Frontend: Static SvelteKit build served via adapter-static

**CI Pipeline:**
- Not detected - No GitHub Actions, GitLab CI, or other CI pipeline in repo
- Manual testing and deployment expected
- Backend: `go build ./...` for compilation verification
- Frontend: `pnpm run build` for static build verification
- Tests: `go test ./...` and `pnpm run test:run` for verification

**Database Migrations in Deployment:**
- Tool: golang-migrate
- Command: `migrate -path migrations -database "$DATABASE_URL" up`
- Must run before application startup to ensure schema is current

## Environment Configuration

**Required env vars:**

*Backend:*
- `DATABASE_URL` - PostgreSQL connection string (takes precedence over config file)
  - Example: `postgres://user:pass@host:5432/dbname?sslmode=disable`
  - Format: `postgres://username:password@hostname:port/dbname?sslmode=disable`

*Frontend (build-time):*
- `VITE_GRAPHQL_URL` - GraphQL endpoint URL (must be set at build time for production)
  - Example: `https://api.perspectize.com/graphql`
  - Fallback: `http://localhost:8080/graphql` in development

**Optional env vars:**

*Backend:*
- `YOUTUBE_API_KEY` - YouTube Data API v3 key (blank/empty if YouTube features not needed)
- `DATABASE_PASSWORD` - Alternative to embedding password in DATABASE_URL
- `APP_ENV` - Set to "production" to suppress .env file warnings
- `CONFIG_PATH` - Path to config JSON (defaults to `config/config.example.json`)

*Frontend:*
- None specified as required beyond VITE_GRAPHQL_URL

**Secrets location:**
- `.env` file (git-ignored) for local development - see `backend/.env.example`
- Environment variables passed via Sevalla deployment platform
- Config file `backend/config/config.example.json` for non-secret defaults (checked into repo)
- No external secrets vault (HashiCorp Vault, AWS Secrets Manager, etc.)

**Config file structure (`backend/config/config.example.json`):**
```json
{
  "server": {
    "port": 8080,
    "host": "0.0.0.0"
  },
  "database": {
    "host": "localhost",
    "port": 5432,
    "name": "perspectize",
    "user": "postgres",
    "sslmode": "require"
  },
  "youtube": {
    "api_key": ""
  },
  "logging": {
    "level": "info",
    "format": "json"
  }
}
```

**Env var Loading Order:**
1. `.env` file loaded by `godotenv.Load()` in `backend/cmd/server/main.go` (if exists)
2. Config file read from `CONFIG_PATH` or `config/config.example.json`
3. Env vars override config file values:
   - `DATABASE_PASSWORD` overrides config `database.password`
   - `DATABASE_URL` bypasses config entirely (full connection string)
   - `YOUTUBE_API_KEY` overrides config `youtube.api_key`

## Webhooks & Callbacks

**Incoming Webhooks:**
- Not implemented - No webhook endpoints for external events

**Outgoing Webhooks:**
- Not implemented - No outgoing event notifications

**GraphQL Subscriptions:**
- Not detected - No real-time subscription support configured (gqlgen supports it, not enabled)

## GraphQL Integration

**Frontend → Backend Communication:**
- Client: `graphql-request` v7.4.0 in `src/lib/queries/client.ts`
- Endpoint: Configurable via `VITE_GRAPHQL_URL` at build time
- Query definitions: Tagged template literals in `src/lib/queries/` (content.ts, users.ts, etc.)
- Integration with TanStack Query: `createQuery()` with function callback wrapper pattern
- Error handling: GraphQL errors surfaced to UI via TanStack Query state
- No custom middleware, interceptors, or authentication headers

**Schema Details:**
- Schema file: `backend/schema.graphql`
- Scalar types: `JSON` (mapped to graphql.Map), `IntID` (custom int scalar)
- Enums: UserRole (ADMIN, SENTINEL, DEFAULT), Privacy (PUBLIC, PRIVATE), ReviewStatus (PENDING, APPROVED, REJECTED)
- Pagination: Cursor-based with opaque base64 encoding (format: `cursor:<id>`)
- Query types defined: content, user, perspectives with filtering and sorting
- Mutation types: createContentFromYouTube, createUser, updateUser, deleteUser, createPerspective, updatePerspective, deletePerspective

---

## Integration Status Summary

| Integration | Status | Priority | Notes |
|-------------|--------|----------|-------|
| PostgreSQL 18 | Active | Critical | Primary data store, fully integrated with GORM |
| YouTube API v3 | Active | High | Content enrichment, custom HTTP client |
| godotenv | Active | Medium | Config loading, local dev only |
| Structured Logging (slog) | Active | Medium | JSON logs for Sevalla compatibility |
| GraphQL (gqlgen) | Active | Critical | Schema-first API definition and code generation |
| TanStack Query | Active | High | Frontend data fetching and caching |
| External Identity | Not Implemented | Medium | Basic user system exists, no OAuth/JWT |
| Error Tracking | Not Implemented | Low | Standard slog logging only |
| CI/CD | Not Implemented | High | Manual deployment expected |
| Redis/Cache | Not Implemented | Low | Direct DB queries, no server-side cache |
| WebSockets/Subscriptions | Not Implemented | Low | gqlgen supports, not yet enabled |

*Integration audit: 2026-02-16*
