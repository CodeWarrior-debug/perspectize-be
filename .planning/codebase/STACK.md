# Technology Stack

**Analysis Date:** 2026-09-04

## Languages

**Primary:**
- Go 1.25 (min), `toolchain go1.26.0` — backend, `backend/`
- TypeScript ~5.9 — frontend, `frontend/src/`

**Secondary:**
- Svelte 5 (`.svelte` files, runes-only) — UI components, `frontend/src/lib/components/`
- SQL — migrations, `backend/migrations/*.sql`
- GraphQL SDL — API schema, `backend/schema.graphql`

## Runtime

**Environment:**
- Go 1.25/1.26 (`backend/go.mod`, pinned base image `golang:1.26-alpine` in `backend/Dockerfile`)
- Node.js (version not pinned via `.nvmrc`; managed via `frontend/package.json` `type: module`)

**Package Manager:**
- Go modules — `backend/go.mod` / `backend/go.sum` (lockfile present)
- pnpm — `frontend/package.json` / `frontend/pnpm-lock.yaml` (lockfile present, `frontend/pnpm-workspace.yaml`)

## Frameworks

**Core:**
- gqlgen v0.17.91 (schema-first GraphQL server) — `backend/gqlgen.yml`, `backend/schema.graphql`, generated code in `backend/internal/adapters/graphql/generated/`
- go-chi/chi v5 (HTTP router) + go-chi/cors + go-chi/httprate (rate limiting) — `backend/cmd/server/main.go`
- GORM v1.31 with `gorm.io/driver/postgres` over `jackc/pgx/v5` — `backend/internal/adapters/repositories/postgres/`
- SvelteKit `~2.70.2` + Svelte `^5.56.3` (runes-only, per `frontend/CLAUDE.md`) — `frontend/src/routes/`
- Vite `^7.3.6` — frontend build/dev server, `frontend/vite.config.ts`

**Testing:**
- testify v1.12.1 (Go unit/integration assertions) — `backend/test/`
- Vitest `^4.1.11` (unit project `unit`, browser project `browser` via `@vitest/browser` + Playwright provider) — `frontend/vitest.config.browser.ts`, `frontend/tests/`
- `@testing-library/svelte`, `@testing-library/jest-dom` — component testing
- jscpd — code duplication detection (`pnpm run test:duplication`)

**Build/Dev:**
- air v1.61.7 (Go hot-reload, `make dev`) — `backend/Makefile`
- golang-migrate (DB migrations, invoked via Makefile targets `migrate-up`/`migrate-down`)
- svelte-check + TypeScript — type checking, `frontend/tsconfig.json`
- Tailwind CSS v4 (`@tailwindcss/vite`) — `frontend/tailwind.config.ts`
- Capacitor 7/8 (`@capacitor/core`, `@capacitor/ios`, `@capacitor/android`) — mobile shell, `frontend/capacitor.config.ts`, `frontend/ios/`, `frontend/android/`

## Key Dependencies

**Critical:**
- `github.com/clerk/clerk-sdk-go/v2` v2.7.0 — auth verification, `backend/internal/adapters/auth/clerk_middleware.go`
- `github.com/svix/svix-webhooks` v1.99.1 — Clerk webhook signature verification, `backend/internal/adapters/auth/webhook_handler.go`
- `github.com/golang-jwt/jwt/v5` v5.3.1 — JWT handling
- `github.com/pilagod/gorm-cursor-paginator/v2` — cursor pagination (installed but hand-rolled cursor helpers in `helpers.go` currently used instead; full integration planned per `backend/CLAUDE.md`)
- `github.com/unrolled/secure` v1.17.0 — HTTP security headers middleware
- `graphql-request` `^7.4.0` + `graphql` `^16.13.1` — frontend GraphQL client, `frontend/src/lib/queries/client.ts`
- `@tanstack/svelte-query` `^6.1.8` — server-state/data-fetching, `frontend/src/lib/queries/`
- `svelte-clerk` `^1.1.9` — Clerk auth integration in SvelteKit
- `@ag-grid-community/*` v32.3.9 + `ag-grid-svelte5` — data grid (`ActivityTable`), `frontend/src/lib/components/`
- `@tiptap/*` v3.22.5 — rich text editor
- `dompurify` `^3.4.1` — HTML sanitization (used with Tiptap content)

**Infrastructure:**
- `go.opentelemetry.io/otel` + `otlptracehttp` + `otel/sdk` v1.46.0 — distributed tracing, `backend/cmd/server/main.go` (`initTracer`)
- `github.com/joho/godotenv` v1.5.1 — `.env` loading in local dev

## Configuration

**Environment:**
- Backend: `.env` (gitignored) loaded via godotenv; template at `backend/.env.example`. Config precedence: env vars > `backend/config/config.example.json` (see `backend/internal/config/config.go` `Load()`).
- Frontend: `.env` (gitignored); template at `frontend/.env.example`. Vite-exposed vars must be prefixed `VITE_`.
- Config loader: `backend/internal/config/config.go`, `backend/internal/config/security.go`, `backend/internal/config/validation.go`.

**Build:**
- Backend: `backend/Dockerfile`, `backend/Makefile`, `backend/gqlgen.yml`
- Frontend: `frontend/vite.config.ts`, `frontend/svelte.config.js`, `frontend/tailwind.config.ts`, `frontend/tsconfig.json`, `frontend/components.json` (shadcn-svelte)

## Platform Requirements

**Development:**
- Go toolchain 1.25+ (see `backend/CLAUDE.md` Go Version Management section)
- pnpm, run from `frontend/` directory (not repo root)
- No local Docker Postgres required for day-to-day dev — database is remote (Sevalla); `backend/docker-compose.yml` exists for local Postgres 18 fallback, CI uses Postgres 17 service container.

**Production:**
- Sevalla hosting (Dockerfile-based backend build, TLS termination + SSL via Cloudflare integration; see `backend/.env.example` and `SEVALLA_BACKEND_URL`)
- PostgreSQL 17 (`.docs/ARCHITECTURE.md`), connected via `DATABASE_URL` (Sevalla external endpoint, `?sslmode=disable`)

---

*Stack analysis: 2026-09-04*
