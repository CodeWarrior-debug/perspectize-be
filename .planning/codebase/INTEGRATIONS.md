# External Integrations

**Analysis Date:** 2026-09-04

## APIs & External Services

**Video Metadata:**
- YouTube Data API v3 — fetches video snippet/statistics/contentDetails for ingested content
  - Client: `backend/internal/adapters/youtube/client.go` (`Client.GetVideoMetadata`)
  - Response parsing/trimming: `backend/internal/adapters/youtube/parser.go`
  - In-memory response cache (TTL-based, resets on restart, not shared across replicas): `backend/internal/adapters/youtube/cache.go`
  - Auth: `YOUTUBE_API_KEY` env var (backend), `VITE_YOUTUBE_API_KEY` (frontend, used on Discover page for search/trending)
  - Error sanitization strips API key from logged errors (`sanitizeYouTubeError` in `client.go`) — key is passed as a URL query param, so this matters
  - Optional integration: metadata fetching fails gracefully if key is unset

**Frontend consumes:**
- Backend GraphQL API — `frontend/src/lib/queries/client.ts` (`graphql-request` `GraphQLClient`), URL via `VITE_GRAPHQL_URL` (defaults `http://localhost:8080/graphql`)

## Data Storage

**Databases:**
- PostgreSQL 17 (production, hosted on Sevalla), PostgreSQL 18 image in local `backend/docker-compose.yml`, PostgreSQL 17 in CI (`.github/workflows/ci.yml`)
  - Connection: `DATABASE_URL` env var (preferred; Sevalla external proxy endpoint, e.g. `us-east1-001.proxy.sevalla.app`, requires `?sslmode=disable`), falls back to discrete `DatabaseConfig` fields in `backend/internal/config/config.go`
  - Driver/client: GORM (`gorm.io/driver/postgres`) over `jackc/pgx/v5`; connection setup in `backend/pkg/database/`
  - Migrations: `golang-migrate`, SQL files in `backend/migrations/`, run via `make migrate-up`/`migrate-down`
  - Database is remote — no local Docker Postgres needed for standard dev workflow

**File Storage:**
- None detected — no object storage (S3/GCS) integration found in backend or frontend dependencies

**Caching:**
- In-memory only: YouTube API response cache (`backend/internal/adapters/youtube/cache.go`), TTL controlled by `YOUTUBE_API_CACHE_TTL_SECONDS` (default 21600s / 6h). No Redis/Memcached.

## Authentication & Identity

**Auth Provider:**
- Clerk (`github.com/clerk/clerk-sdk-go/v2` backend, `svelte-clerk` frontend)
  - Backend middleware: `backend/internal/adapters/auth/clerk_middleware.go` — verifies session tokens, sets `clerk.SetKey(secCfg.ClerkSecretKey)` in `backend/cmd/server/main.go`
  - Secret: `CLERK_SECRET_KEY` (backend env var)
  - Publishable key: `VITE_CLERK_PUBLISHABLE_KEY` (frontend env var)
  - User sync via webhook (see Webhooks below)
  - Additional app-level auth: `JWT_SECRET` + `ACCESS_TOKEN_MINUTES` (`golang-jwt/jwt/v5`) — app-issued tokens, config in `backend/.env.example`
  - GraphQL-level authorization: custom `@auth`/`@owner` directives — `backend/internal/adapters/graphql/directives/auth.go` (resolves resource ownership by reading typed input structs via JSON-tag reflection, not raw maps)
  - Auth context helpers: `backend/internal/adapters/auth/context.go` (`ForContext`, `RequireAuth`)
  - Historical provider comparison research (Clerk chosen): `.planning/phases/12-authentication/12-RESEARCH-provider-comparison.md`

## Monitoring & Observability

**Error Tracking:**
- None detected (no Sentry/Bugsnag/Rollbar dependency found)

**Tracing:**
- OpenTelemetry — `go.opentelemetry.io/otel`, `otlptracehttp` exporter, `otel/sdk`
  - Setup: `initTracer()` in `backend/cmd/server/main.go`, activates only when `OTEL_EXPORTER_OTLP_ENDPOINT` env var is set
  - Exporter reads standard OTel env vars automatically: `OTEL_EXPORTER_OTLP_ENDPOINT`, `OTEL_EXPORTER_OTLP_HEADERS`, `OTEL_SERVICE_NAME`
  - Frontend has `web-vitals` `^5.3.0` dependency for Core Web Vitals collection (no external reporting endpoint wired beyond app code)

**Logs:**
- Structured logging via Go `log/slog` (backend); no external log aggregation service detected

## CI/CD & Deployment

**Hosting:**
- Sevalla (backend Dockerfile-based deploy; Dockerfile path `backend/Dockerfile`, build context `backend`)
  - `SEVALLA_BACKEND_URL` env var references deployed backend URL
  - TLS/SSL handled by Sevalla via Cloudflare — no `ListenAndServeTLS` in Go code

**CI Pipeline:**
- GitHub Actions — `.github/workflows/ci.yml` (backend test/lint/build against a Postgres 17 service container, `go test -race`, coverage upload to Codecov), `.github/workflows/frontend-test.yml` (frontend tests), `.github/workflows/codeql.yml` (CodeQL static analysis), `.github/workflows/trivy.yml` (dependency/container vulnerability scanning)
- Coverage reporting: Codecov (`codecov/codecov-action@v7`)

## Environment Configuration

**Required env vars (backend, `backend/.env.example`):**
- `DATABASE_URL` — Postgres connection string (Sevalla)
- `JWT_SECRET`, `ACCESS_TOKEN_MINUTES` — app JWT config
- `CLERK_SECRET_KEY` — Clerk backend API key
- `CLERK_WEBHOOK_SIGNING_SECRET` — Svix/Clerk webhook signature verification
- `RATE_LIMIT_PER_MIN` — go-chi/httprate limit
- `CORS_ORIGINS` — comma-separated allowed origins
- `YOUTUBE_API_KEY`, `YOUTUBE_API_CACHE_TTL_SECONDS` — optional YouTube integration
- `APP_ENV` — `development`/`production` (toggles HSTS, GraphQL introspection, playground)
- `SEVALLA_BACKEND_URL` — deployed backend URL reference

**Required env vars (frontend, `frontend/.env.example`):**
- `VITE_CLERK_PUBLISHABLE_KEY`
- `VITE_GRAPHQL_URL`
- `VITE_YOUTUBE_API_KEY`

**Secrets location:**
- Local: `.env` files (gitignored) in `backend/` and `frontend/`, never committed
- Production: environment variables configured in Sevalla platform (no secrets manager/vault integration detected)
- Rotation/incident-response guidance: `.docs/SECURITY.md`

## Webhooks & Callbacks

**Incoming:**
- `POST /webhooks/clerk` — Clerk user-sync webhook, verified via Svix signing secret (`CLERK_WEBHOOK_SIGNING_SECRET`)
  - Handler: `backend/internal/adapters/auth/webhook_handler.go`
  - Registered route: `backend/cmd/server/main.go` (`r.Post("/webhooks/clerk", webhookHandler.ServeHTTP)`)
  - Syncs Clerk user data into local `users` table: `backend/internal/adapters/repositories/postgres/gorm_user_repository.go`

**Outgoing:**
- None detected — no outgoing webhook dispatch found in codebase

---

*Integration audit: 2026-09-04*
