# Architecture

**Analysis Date:** 2026-09-04

## Pattern Overview

**Overall:** Monorepo with two independently deployed stacks: a Go GraphQL API backend (`backend/`) built with **Hexagonal Architecture** (Ports and Adapters), and a SvelteKit SPA frontend (`frontend/`) consuming that API over GraphQL.

**Key Characteristics:**
- Backend domain logic has zero framework/infra dependencies; all I/O goes through port interfaces
- Schema-first GraphQL (gqlgen) — `backend/schema.graphql` is the source of truth, generated code lives in `internal/adapters/graphql/generated/`
- Frontend is a client-rendered SPA (`ssr = false`, `csr = true`) — no server-side rendering, talks to the backend purely via GraphQL over HTTP
- Auth is delegated to Clerk (JWT bearer tokens on the frontend, Clerk SDK verification + webhook sync on the backend)
- Multi-dimensional rating domain: "Perspectives" hold `quality` and `agreement` as 0–10000 integers (0.01% precision, avoids float issues) rather than binary like/dislike

## Layers (Backend — `backend/`)

**Domain Layer (Core):**
- Purpose: Pure business entities and rules, no external dependencies
- Location: `backend/internal/core/domain/`
- Contains: `content.go`, `perspective.go`, `user.go`, `auth.go`, `claims.go`, `errors.go`, `pagination.go`
- Depends on: nothing (standard library only)
- Used by: services, ports

**Ports Layer:**
- Purpose: Interfaces defining contracts that adapters must satisfy
- Location: `backend/internal/core/ports/repositories/` (e.g. `content_repository.go`, `user_repository.go`, `perspective_repository.go`) and `backend/internal/core/ports/services/` (e.g. `content_service.go`, `auth_service.go`)
- Depends on: domain
- Used by: services (define contracts), adapters (implement contracts)

**Services Layer (Business Logic):**
- Purpose: Orchestrates domain rules, validation, cross-entity logic
- Location: `backend/internal/core/services/` — `content_service.go`, `user_service.go`, `perspective_service.go`, `auth_service.go`
- Depends on: domain, ports (interfaces only — never concrete adapters)
- Used by: GraphQL resolvers, directives

**Primary Adapter Layer (Driving — inbound):**
- Purpose: Exposes the application to the outside world
- Location: `backend/internal/adapters/graphql/` — `resolvers/` (`resolver.go`, `schema.resolvers.go`, `helpers.go`), `directives/` (`auth.go` — `@auth`/`@owner` directive implementations), `generated/` (gqlgen output, do not hand-edit), `model/` (GraphQL input/output structs)
- Depends on: services (via constructor injection)
- Used by: `cmd/server/main.go` (wires resolver into gqlgen handler)

**Secondary Adapter Layer (Driven — outbound):**
- Purpose: Implements port interfaces against real infrastructure
- Location: `backend/internal/adapters/repositories/postgres/` (GORM-backed repos: `gorm_content_repository.go`, `gorm_user_repository.go`, `gorm_perspective_repository.go`, `gorm_mappers.go`, `gorm_models.go`, `helpers.go`), `backend/internal/adapters/youtube/` (`client.go`, `cache.go`, `parser.go` — YouTube Data API client with in-memory TTL cache), `backend/internal/adapters/auth/` (`auth.go` middleware, `webhook_handler.go` for Clerk webhooks, `context.go`, `claims.go`)
- Depends on: domain, ports (implements them)
- Used by: `cmd/server/main.go` (constructed and injected into services)

**Infrastructure/Wiring Layer:**
- Purpose: Composition root — constructs adapters, services, and the HTTP server, wires everything together
- Location: `backend/cmd/server/main.go`
- Depends on: everything (only place allowed to import both adapters and services concretely)

**Cross-cutting Support (`pkg/`):**
- `backend/pkg/database/` — GORM connection setup (`postgres.go`), pool config, slow-query logging, DB stats endpoint (`stats.go`)
- `backend/pkg/graphql/` — `intid.go` (custom `IntID` scalar), `timing.go` (operation timing middleware for gqlgen)
- `backend/pkg/logger/` — structured `slog` JSON setup (`logger.go`)
- `backend/pkg/middleware/` — generic HTTP middleware not specific to auth: `recovery.go` (panic recovery → JSON via slog), `timing.go` (request timing)
- `backend/internal/config/` — `config.go` (env/JSON config loading), `security.go` (CORS origins, rate limits, Clerk secrets), `validation.go` (DATABASE_URL validation)
- `backend/internal/adapters/web/middleware/` — HTTP-layer middleware wired in `main.go`: `auth.go`, `contenttype.go` (CSRF via Content-Type), `ratelimit.go`, `secureheaders.go` (HSTS, X-Frame-Options, etc.)

**Dependency Rule:** Dependencies point inward only. `core/domain` never imports anything from `adapters/` or `pkg/`. `core/services` depends only on `core/ports` interfaces, never on concrete adapter types. `cmd/server/main.go` is the sole place that imports both adapters and services concretely to wire them together.

## Layers (Frontend — `frontend/src/`)

**Routing Layer:**
- Purpose: SvelteKit file-based routing, page composition
- Location: `frontend/src/routes/` — `+layout.svelte` (root layout: QueryClientProvider, Header, Toaster), `+layout.ts` (`ssr = false`, `csr = true`, `prerender = false`), `+page.svelte` (home page), `discover/+page.svelte` + `discover/+page.ts` (discover feature route)
- Depends on: components, queries, stores
- Used by: browser navigation only (SPA — no server rendering)

**Component Layer:**
- Purpose: Presentational and feature Svelte 5 components
- Location: `frontend/src/lib/components/` — feature components at top level (`ActivityTable.svelte`, `ActivityCardList.svelte`, `ActivityDetailsModal.svelte`, `AddVideoDialog.svelte`, `PerspectivePopover.svelte`, `RatingInput.svelte`, `Header.svelte`, `PageWrapper.svelte`, etc.), `shadcn/` (shadcn-svelte UI primitives: `button/`, `dialog/`, `drawer/`, `input/`, `label/`, `popover/`, `select/`), `discover/` (discover-page-specific components)
- Depends on: queries/hooks, stores, utils
- Used by: routes

**Query/Data Fetching Layer:**
- Purpose: GraphQL query/mutation definitions and TanStack Query integration
- Location: `frontend/src/lib/queries/` — `client.ts` (GraphQLClient instance, `VITE_GRAPHQL_URL`, `getAuthToken()`/`graphqlRequest()` helpers that attach Clerk bearer tokens), `keys.ts` (centralized query-key factory for cache invalidation), `content.ts`, `claims.ts`, `perspectives.ts`, `users.ts` (gql tagged-template query/mutation definitions), `hooks/` (TanStack Query wrapper hooks: `useAddVideo.ts`, `useCreateClaim.ts`, `useCreatePerspective.ts`, `useCreateUser.ts`, `useMe.svelte.ts`, `useUpdatePerspective.ts`, `useUpdateSourceData.ts`)
- Depends on: `graphql-request`, `@tanstack/svelte-query`
- Used by: components

**State Management Layer:**
- Purpose: Cross-component reactive state outside TanStack Query cache
- Location: `frontend/src/lib/stores/` — `userSelection.svelte.ts` (Svelte 5 rune-based store, `.svelte.ts` suffix required for rune usage outside `.svelte` files)
- Depends on: nothing internal
- Used by: components needing shared UI state (e.g. selected user for perspective filtering)

**Services Layer:**
- Purpose: Non-GraphQL external integrations
- Location: `frontend/src/lib/services/` — `youtubeApi.ts` (client-side YouTube helper)

**Utilities & Assets:**
- Location: `frontend/src/lib/utils/` — `formatting.ts`, `grid-config.ts` (AG Grid column/sort/pagination logic extracted as pure functions for testability), `gridUrlState.ts`, `ratings.ts`, `sanitize.ts`, `youtube.ts`, `native.ts`, `references.ts`, `activityItemCellRenderer.ts`; top-level `frontend/src/lib/utils.ts` (shared cn/class helpers), `frontend/src/lib/index.ts` (barrel export), `frontend/src/lib/vitals.ts`
- Location: `frontend/src/lib/assets/`, `frontend/src/assets/glasses_svgs/` — static image/SVG assets

**Global Styles:**
- Location: `frontend/src/app.css` (Tailwind v4 `@theme` design tokens), `frontend/src/app.html` (HTML shell, CSP config)

## Data Flow

**GraphQL Query (e.g. list content):**
1. Component calls a TanStack Query hook (function-wrapper pattern: `createQuery(() => ({...}))`) — `frontend/src/lib/queries/hooks/` or inline in component
2. Hook calls `graphqlClient.request(...)` (or `graphqlRequest()` for authenticated calls) from `frontend/src/lib/queries/client.ts`, using a query defined in `frontend/src/lib/queries/content.ts`
3. Request POSTed to `VITE_GRAPHQL_URL` (defaults `http://localhost:8080/graphql`) with `Authorization: Bearer <clerk-token>` if signed in
4. Backend chi router (`backend/cmd/server/main.go`) applies middleware stack (rate limit → CORS → security headers → content-type validation → auth middleware → request timer → recoverer) then routes to `/graphql` handler (gqlgen)
5. gqlgen dispatches to resolver method in `backend/internal/adapters/graphql/resolvers/schema.resolvers.go`
6. Resolver calls into a `core/services` method (e.g. `ContentService.List`)
7. Service applies business rules, calls repository port method
8. `postgres.GormContentRepository` (implementing the port) executes the GORM query, maps GORM model → domain model via `gorm_mappers.go`
9. Response flows back up: repository → service → resolver → gqlgen → JSON response
10. TanStack Query caches the response under a hierarchical key (`queryKeys.content.list(filters)`)

**Mutation with Auth Directive (e.g. update perspective):**
1. GraphQL schema field annotated `@auth` or `@owner` in `backend/schema.graphql`
2. gqlgen invokes directive resolver in `backend/internal/adapters/graphql/directives/auth.go` before the field resolver
3. `@owner` directive extracts the resource ID from the typed input struct via reflection (`fieldByJSONTag`) — NOT from a `map[string]interface{}` type assertion, which silently fails for real typed gqlgen inputs
4. Directive validates the authenticated user (set into context by `auth.Middleware` in `backend/internal/adapters/auth/auth.go`, which verifies the Clerk JWT) owns the resource, else returns `ErrForbidden`
5. On success, field resolver executes normally

**Clerk User Sync (webhook):**
1. Clerk sends webhook events (user created/updated) to `POST /webhooks/clerk`
2. Route bypasses the standard auth middleware — Svix signature verification is the auth mechanism (`backend/internal/adapters/auth/webhook_handler.go`)
3. Handler upserts the user via `UserRepository`

**State Management (frontend):**
- Server state (GraphQL data) lives entirely in TanStack Query's cache, keyed via the factory in `frontend/src/lib/queries/keys.ts`
- Non-server UI state (e.g. currently selected user) lives in Svelte 5 rune-based stores in `frontend/src/lib/stores/` (`.svelte.ts` file suffix enables `$state`/`$derived` outside components)
- No global client-side store framework (no Redux/Zustand equivalent) — TanStack Query + Svelte 5 runes covers both concerns

## Key Abstractions

**Domain Models:**
- Purpose: Represent core business entities independent of storage/transport
- Examples: `backend/internal/core/domain/content.go`, `backend/internal/core/domain/perspective.go`, `backend/internal/core/domain/user.go`
- Pattern: Plain Go structs, zero GORM/gqlgen tags — kept strictly separate from GORM persistence models

**GORM Models (Hex-Clean Separate Model Pattern):**
- Purpose: Persistence-layer representation, decoupled from domain models
- Examples: `backend/internal/adapters/repositories/postgres/gorm_models.go`
- Pattern: `gorm:` tagged structs; bidirectional mapping to/from domain models happens explicitly in `gorm_mappers.go` — domain layer never sees a GORM tag

**Repository Ports:**
- Purpose: Define storage contracts the domain/services depend on, implemented by adapters
- Examples: `backend/internal/core/ports/repositories/content_repository.go`, `.../user_repository.go`, `.../perspective_repository.go`
- Pattern: Interface in `core/ports`, concrete GORM implementation in `adapters/repositories/postgres/gorm_*_repository.go`

**Cursor-based Pagination:**
- Purpose: Stable pagination over large lists without OFFSET
- Examples: `backend/internal/core/domain/pagination.go`, `backend/internal/adapters/repositories/postgres/helpers.go` (`encodeCursor`/`decodeCursor`)
- Pattern: Opaque base64 cursor (`cursor:<id>`), keyset pagination, fetch `limit+1` rows to compute `hasNextPage`, sort columns whitelisted to prevent SQL injection

**IntID Scalar:**
- Purpose: Type-safe integer IDs in GraphQL filter/input fields (vs the built-in `ID` string scalar)
- Examples: `backend/pkg/graphql/intid.go`, bound in `backend/gqlgen.yml`
- Pattern: Custom gqlgen scalar; top-level query/mutation ID args still use plain `ID!` + `strconv.Atoi`

**Query Key Factory (frontend):**
- Purpose: Type-safe, hierarchical TanStack Query cache keys for predictable invalidation
- Examples: `frontend/src/lib/queries/keys.ts`
- Pattern: Nested object of key-builder functions per entity (`queryKeys.content.list(filters)`, `queryKeys.perspectives.detail(id)`)

## Entry Points

**Backend Server:**
- Location: `backend/cmd/server/main.go`
- Triggers: `go run ./cmd/server`, `make run`, `make dev` (air hot-reload), Docker/Sevalla deployment
- Responsibilities: load config/env, connect to Postgres (GORM), init Clerk SDK, construct repositories → services → resolver, build gqlgen handler with directives/complexity limits/persisted-query cache, build chi router with full middleware stack, register `/health`, `/ready`, `/graphql`, `/webhooks/clerk`, start HTTP server with graceful shutdown

**Frontend App Shell:**
- Location: `frontend/src/routes/+layout.svelte` + `frontend/src/routes/+layout.ts`
- Triggers: any page load (SPA, client-side only — `ssr = false`)
- Responsibilities: mount `QueryClientProvider`, render `Header`, `Toaster`, wrap page content

**GraphQL Schema (contract):**
- Location: `backend/schema.graphql`
- Triggers: `make graphql-gen` regenerates `backend/internal/adapters/graphql/generated/` after any schema edit
- Responsibilities: single source of truth for the API contract consumed by both resolver implementations and (implicitly) frontend query definitions

## Error Handling

**Strategy (backend):** Sentinel errors defined once in the domain layer, propagated up through services and translated to GraphQL errors at the resolver boundary.

**Patterns:**
- Domain sentinel errors: `backend/internal/core/domain/errors.go` — `ErrNotFound`, `ErrAlreadyExists`, `ErrInvalidInput`, `ErrInvalidURL`, `ErrYouTubeAPI`, `ErrInvalidRating`, `ErrSentinelUser`, `ErrDeleteSentinel`
- Services return these sentinel errors (via `errors.Is`-compatible wrapping); resolvers/directives map them to GraphQL error responses
- Auth-specific errors surface as `ErrUnauthorized`/`ErrForbidden` from the `@auth`/`@owner` directives (`backend/internal/adapters/graphql/directives/auth.go`)
- Panic recovery centralized in `backend/pkg/middleware/recovery.go` (structured JSON via `slog`, not chi's default logger)

**Strategy (frontend):** TanStack Query's built-in `isError`/`error` reactive state per query/mutation; no global error boundary framework beyond that.

## Cross-Cutting Concerns

**Logging:** `log/slog` structured JSON logging throughout the backend (`backend/pkg/logger/logger.go`, `RegisterSlowQueryLogger` in `pkg/database`, request timing in `pkg/middleware/timing.go` and `pkg/graphql/timing.go`). No backend `fmt.Println`/`log.Println` in request-handling code paths.

**Validation:** `go-playground/validator` struct-tag validation on backend inputs; `backend/internal/config/validation.go` validates `DATABASE_URL` format specifically.

**Authentication:** Clerk (hosted auth) — frontend obtains JWT via `window.Clerk.session.getToken()` (`frontend/src/lib/queries/client.ts`), backend verifies via Clerk SDK in `auth.Middleware` (`backend/internal/adapters/auth/auth.go`), field-level authorization enforced declaratively via `@auth`/`@owner` GraphQL directives rather than imperative checks scattered in resolvers.

**Rate limiting / security headers:** Applied as chi middleware in `backend/cmd/server/main.go` before auth (`apimw.GlobalRateLimit`, `cors.Handler`, `apimw.SecureHeaders`, `apimw.ContentTypeValidation`) — order is deliberate (rate limit before auth to prevent DoS from unauthenticated flood).

**Observability:** OpenTelemetry tracing optionally enabled via `OTEL_EXPORTER_OTLP_ENDPOINT` env var (`initTracer` in `main.go`); GraphQL operation timing middleware (`pkg/graphql/timing.go`); `/debug/db-stats` endpoint (non-production only) exposes connection pool stats.

---

*Architecture analysis: 2026-09-04*
