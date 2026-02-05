# External Integrations

**Analysis Date:** 2026-02-04

## APIs & External Services

**YouTube Data API v3:**
- **Purpose:** Fetch video metadata (title, description, duration, view counts) for content enrichment
- **SDK/Client:** Custom HTTP client implementation in `internal/adapters/youtube/client.go`
- **Auth:** Environment variable `YOUTUBE_API_KEY` (set via `.env` or config.json)
- **Endpoint:** `https://www.googleapis.com/youtube/v3/videos`
- **Status:** Active integration with functional client implementation
  - Fetches: `snippet`, `statistics`, `contentDetails` parts
  - Extracts: title, description, channelTitle, duration, viewCount, likeCount, commentCount
  - Returns: Structured `VideoMetadata` to domain services via port interface `services.VideoMetadata`
- **Error Handling:** Returns custom domain error `domain.ErrYouTubeAPI` on API failures
- **Response Storage:** Full API response stored as JSONB in `content.response` database column for audit trail

## Data Storage

**Databases:**

**PostgreSQL 18:**
- **Connection:** Via environment variable `DATABASE_URL` or config file settings
- **Client:** sqlx v1.4.0 (query builder) + pgx/v5 (driver)
- **Connection Details:**
  - Default local: `postgres://testuser:testpass@localhost:5432/testdb?sslmode=disable`
  - Docker container: `perspectize-postgres-go` (via docker-compose.yml)
  - Port: 5432
- **Advanced Features Used:**
  - JSONB columns: `content.response` stores YouTube API responses
  - Array types: `categorized_ratings` array column for perspective ratings
  - Custom domains: Type validation at database level
  - Triggers: Automatic timestamp updates
- **Connection Pool:**
  - Max open: 25
  - Max idle: 5
  - Max lifetime: 5 minutes
- **Migrations:** Located in `perspectize-go/migrations/`
  - Tool: golang-migrate
  - Format: `{sequence}_{description}.{up|down}.sql`
  - Current: 5 migrations (create content, update response jsonb, update length numeric, add perspectives/users, add user timestamps)

**File Storage:**
- Not applicable - Local filesystem only for migrations and config files
- All user data stored in PostgreSQL

**Caching:**
- Not implemented - No Redis, Memcached, or in-memory cache configured
- Query results served directly from PostgreSQL

## Authentication & Identity

**Auth Provider:**
- Custom implementation (no external identity provider configured)
- Current approach: User system in domain with ID-based references
- `internal/core/domain/user.go` defines user entity
- `internal/adapters/repositories/postgres/user_repository.go` manages persistence

**Status:** Basic user CRUD operations, no OAuth/JWT/session management currently implemented

## Monitoring & Observability

**Error Tracking:**
- Not configured - No Sentry, Rollbar, or error aggregation service
- Errors logged to stdout via standard library `log` package

**Logs:**
- **Framework:** Go standard library `log` package (not yet upgraded to slog/structured logging)
- **Current approach:** Printf-style logging in main.go
- **Logs captured:** Database connection, server startup, migrations
- **Configuration:** Defined in config.json `logging` field (level, format) but not yet consumed by application
- **WIP Note:** CLAUDE.md indicates `slog` (log/slog) should be used but not yet implemented in codebase

**Log Types Generated:**
- Database connectivity logs
- Server startup/shutdown
- Request/response logging (via gqlgen default handlers)

## CI/CD & Deployment

**Hosting:**
- Target platforms: Fly.io or Sevalla (mentioned in CLAUDE.md)
- Build process: `go build -o bin/perspectize-server cmd/server/main.go`
- Container image: Not yet built, would need Dockerfile

**CI Pipeline:**
- Not detected - No GitHub Actions, GitLab CI, or other CI pipeline configured
- Manual testing and deployment expected

**Database Migrations in CI/CD:**
- Tool: golang-migrate
- Command: `migrate -path migrations -database "$DATABASE_URL" up`
- Must run before application startup to ensure schema is current

## Environment Configuration

**Required env vars:**
- `DATABASE_URL` - PostgreSQL connection string (takes precedence over config file)
  - Example: `postgres://user:pass@host:5432/dbname?sslmode=disable`
  - Fallback: Constructed from config.json if not set
  - Note: Sevalla may require `?sslmode=disable` and connections may succeed on retry

**Optional env vars:**
- `YOUTUBE_API_KEY` - YouTube Data API v3 key (blank/empty if YouTube features not needed)
- `DATABASE_PASSWORD` - Alternative to embedding password in DATABASE_URL

**Secrets location:**
- `.env` file (git-ignored) for local development
- Environment variables passed via deployment platform (Fly.io secrets, Sevalla env vars)
- Config file `config/config.example.json` for non-secret defaults

**Config file structure (`perspectize-go/config/config.example.json`):**
```json
{
  "server": {
    "port": 8080,
    "host": "localhost"
  },
  "database": {
    "host": "localhost",
    "port": 5432,
    "name": "testdb",
    "user": "testuser",
    "sslmode": "disable"
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

**Env var Loading:**
1. `.env` file loaded by `godotenv.Load()` in `cmd/server/main.go`
2. Config file read from `config/config.example.json`
3. Env vars override config file values:
   - `DATABASE_PASSWORD` overrides config `database.password`
   - `DATABASE_URL` bypasses config entirely (uses full connection string)
   - `YOUTUBE_API_KEY` overrides config `youtube.api_key`

## Webhooks & Callbacks

**Incoming Webhooks:**
- Not implemented - No webhook endpoints for external events

**Outgoing Webhooks:**
- Not implemented - No outgoing event notifications

**GraphQL Subscriptions:**
- Not detected - No real-time subscription support configured

---

## Integration Status Summary

| Integration | Status | Priority | Notes |
|-------------|--------|----------|-------|
| PostgreSQL 18 | Active | Critical | Primary data store, fully integrated |
| YouTube API v3 | Active | High | Content enrichment, custom client |
| godotenv | Active | Medium | Config loading, local dev only |
| External Identity | Not Implemented | Medium | Basic user system exists, no OAuth/JWT |
| Error Tracking | Not Implemented | Low | Standard logging only |
| Structured Logging | WIP | Low | Config defined, slog not yet used |
| CI/CD | Not Implemented | High | Manual deployment expected |
| Redis/Cache | Not Implemented | Low | Direct DB queries |
| WebSockets | Configured | Low | gqlgen supports, not yet used |

*Integration audit: 2026-02-04*
