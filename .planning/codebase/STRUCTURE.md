# Codebase Structure

**Analysis Date:** 2026-09-04

## Directory Layout

```
perspectize/                          # Monorepo root
├── backend/                          # Go GraphQL API (hexagonal architecture)
│   ├── cmd/server/                   # Entry point (main.go)
│   ├── config/                       # Runtime JSON config (config.example.json)
│   ├── internal/
│   │   ├── core/
│   │   │   ├── domain/               # Pure domain models + sentinel errors
│   │   │   ├── ports/
│   │   │   │   ├── repositories/     # Repository interfaces
│   │   │   │   └── services/         # Service interfaces
│   │   │   └── services/             # Business logic implementations
│   │   ├── adapters/
│   │   │   ├── graphql/
│   │   │   │   ├── directives/       # @auth/@owner directive implementations
│   │   │   │   ├── generated/        # gqlgen output — DO NOT hand-edit
│   │   │   │   ├── model/            # GraphQL input/output model structs
│   │   │   │   └── resolvers/        # Resolver implementations
│   │   │   ├── repositories/postgres/ # GORM repository implementations
│   │   │   ├── auth/                 # Clerk middleware, webhook handler
│   │   │   ├── web/middleware/       # HTTP middleware (CORS-adjacent, rate limit, headers)
│   │   │   └── youtube/              # YouTube Data API client + cache
│   │   └── config/                   # Config loading, security config, validation
│   ├── pkg/                          # Reusable, non-domain-specific packages
│   │   ├── database/                 # GORM connection, pool, slow-query log, stats
│   │   ├── graphql/                  # IntID scalar, gqlgen operation timing
│   │   ├── logger/                   # slog setup
│   │   └── middleware/               # Generic HTTP middleware (recovery, timing)
│   ├── migrations/                   # golang-migrate SQL migration files
│   ├── test/                         # Test suites, mirrors internal/ by concern
│   │   ├── config/ database/ domain/ graphql/ resolvers/ services/ youtube/
│   ├── coverage/                     # Generated coverage reports (not committed source)
│   ├── schema.graphql                # GraphQL schema (source of truth)
│   ├── gqlgen.yml                    # gqlgen codegen config (model bindings)
│   └── Dockerfile                    # Sevalla build (context=backend, path=backend/Dockerfile)
│
├── frontend/                         # SvelteKit SPA
│   ├── src/
│   │   ├── routes/                   # File-based routing
│   │   │   ├── +layout.svelte        # Root layout (QueryClientProvider, Header, Toaster)
│   │   │   ├── +layout.ts            # ssr=false, csr=true, prerender=false
│   │   │   ├── +page.svelte          # Home page
│   │   │   └── discover/             # Discover feature route
│   │   ├── lib/
│   │   │   ├── components/           # Svelte 5 components
│   │   │   │   ├── shadcn/           # shadcn-svelte primitives (button/, dialog/, drawer/, ...)
│   │   │   │   └── discover/         # Discover-page-specific components
│   │   │   ├── queries/              # GraphQL query/mutation defs + TanStack Query
│   │   │   │   └── hooks/            # TanStack Query hook wrappers (useXxx)
│   │   │   ├── services/             # Non-GraphQL client integrations (youtubeApi.ts)
│   │   │   ├── stores/               # Svelte 5 rune-based shared state (.svelte.ts)
│   │   │   ├── utils/                # Pure utility/helper functions
│   │   │   └── assets/               # Static assets bundled via lib
│   │   ├── assets/glasses_svgs/      # Additional static SVG assets
│   │   ├── app.css                   # Tailwind v4 design tokens (@theme)
│   │   └── app.html                  # HTML shell, CSP
│   ├── static/                       # Public static files (served as-is)
│   ├── tests/                        # Vitest tests (unit/, components/, browser/, fixtures/, helpers/)
│   ├── docs/                         # Frontend-specific docs (AG_GRID.md, DESIGN_SPEC.md, FIGMA.md)
│   ├── android/ ios/                 # Capacitor native shells (if mobile packaging is active)
│   └── build/ .svelte-kit/           # Generated build output — not source
│
├── .docs/                            # Monorepo-level reference docs (ARCHITECTURE, SECURITY, GO_PATTERNS, ...)
├── .github/                          # PR/issue templates, CI workflows
├── .planning/                        # GSD legacy planning artifacts (phases, roadmap, codebase docs)
├── docs/superpowers/                 # Active planning system: plans/ and specs/
├── .claude/                          # Claude Code config: agents, commands, skills, docs
└── graphify-out/                     # Knowledge graph output (graph.json, wiki, cache)
```

## Directory Purposes

**`backend/internal/core/domain/`:**
- Purpose: Pure business entities and domain-level errors — no external dependencies (no GORM, no gqlgen)
- Contains: `content.go`, `perspective.go`, `user.go`, `auth.go`, `claims.go`, `errors.go`, `pagination.go`
- Key files: `errors.go` (all sentinel errors used across the backend)

**`backend/internal/core/ports/`:**
- Purpose: Interfaces defining what adapters must implement
- Contains: `repositories/` (data-access contracts), `services/` (service contracts, mostly for auth)

**`backend/internal/core/services/`:**
- Purpose: Business logic — validation, orchestration, rules — depends only on ports (never concrete adapters)
- Key files: `content_service.go`, `user_service.go`, `perspective_service.go`, `auth_service.go`

**`backend/internal/adapters/graphql/resolvers/`:**
- Purpose: gqlgen resolver implementations, the primary (driving) adapter
- Key files: `resolver.go` (struct + constructor wiring services), `schema.resolvers.go` (generated stubs filled in by hand), `helpers.go` (cursor encode/decode, sort whitelisting)

**`backend/internal/adapters/repositories/postgres/`:**
- Purpose: GORM-backed implementations of repository ports (Hex-Clean Separate Model Pattern)
- Key files: `gorm_models.go` (persistence structs), `gorm_mappers.go` (domain ↔ GORM conversion), `gorm_content_repository.go`, `gorm_user_repository.go`, `gorm_perspective_repository.go`, `helpers.go`
- Note: `.sqlx.bak` files present (`content_repository.go.sqlx.bak`, etc.) are legacy sqlx implementations kept as reference, superseded by GORM — do not extend them

**`backend/migrations/`:**
- Purpose: Version-controlled SQL migrations (golang-migrate format: `NNNNNN_description.up.sql` / `.down.sql`)
- Note: always check existing files with `ls backend/migrations/ | tail -5` before assigning a new number — plan-specified numbers can be stale

**`frontend/src/lib/queries/`:**
- Purpose: All GraphQL communication — query/mutation string definitions, the GraphQL client, cache-key factory
- Key files: `client.ts` (GraphQLClient + auth header injection), `keys.ts` (query key factory), `content.ts`/`claims.ts`/`perspectives.ts`/`users.ts` (per-entity `gql` definitions), `hooks/` (TanStack Query wrapper hooks per mutation)

**`frontend/src/lib/components/shadcn/`:**
- Purpose: shadcn-svelte UI primitives, installed via CLI
- Note: CLI sometimes installs into a `ui/` directory instead of `shadcn/` despite `components.json` alias config — after installing, verify location and move if needed; always add new components to `shadcn/index.ts` barrel export

**`frontend/tests/`:**
- Purpose: Vitest test suite
- Contains: `unit/`, `components/`, `browser/` (Vitest Browser Mode), `fixtures/`, `helpers/` (e.g. `TestWrapper.svelte` for dynamic component testing)

## Key File Locations

**Entry Points:**
- `backend/cmd/server/main.go`: backend HTTP server composition root
- `frontend/src/routes/+layout.svelte` + `+layout.ts`: frontend app shell

**Configuration:**
- `backend/config/config.example.json`, `backend/internal/config/config.go`: backend config loading (env > JSON)
- `backend/internal/config/security.go`: CORS origins, rate limits, Clerk secret loading
- `backend/gqlgen.yml`: GraphQL codegen config, model bindings
- `frontend/vite.config.ts`, `frontend/svelte.config.js`: frontend build config
- `frontend/src/app.css`: design tokens

**Core Logic:**
- `backend/internal/core/services/`: all backend business logic
- `backend/schema.graphql`: API contract
- `frontend/src/lib/utils/grid-config.ts`: AG Grid pure-function logic (sort mapping, pagination bounds, responsive tiers) — extracted here specifically for testability since AG Grid itself doesn't render in jsdom

**Testing:**
- `backend/test/`: mirrors `internal/` by concern (config, database, domain, graphql, resolvers, services, youtube)
- `frontend/tests/`: unit, components, browser, fixtures, helpers

## Naming Conventions

**Backend files:**
- Go source: `snake_case.go` (e.g. `content_service.go`, `gorm_content_repository.go`)
- Tests: `*_test.go` co-located logically under `backend/test/<concern>/` (not co-located with source files)
- GORM repository implementations prefixed `gorm_` to distinguish from the port interface and legacy `.sqlx.bak` files

**Frontend files:**
- Svelte components: `PascalCase.svelte` (e.g. `ActivityTable.svelte`, `AddVideoDialog.svelte`)
- TypeScript modules: `camelCase.ts` (e.g. `formatting.ts`, `gridUrlState.ts`)
- Svelte 5 rune-based `.ts` files needing runes outside a `.svelte` file: `*.svelte.ts` suffix (e.g. `userSelection.svelte.ts`, `useMe.svelte.ts`)
- TanStack Query hooks: `useXxx.ts` / `useXxx.svelte.ts`
- SvelteKit route files: fixed names `+page.svelte`, `+page.ts`, `+layout.svelte`, `+layout.ts` per SvelteKit convention

**Directories:**
- Backend: lowercase, singular-plural per Go convention (`domain`, `ports`, `services`, `repositories`)
- Frontend: lowercase (`components`, `queries`, `stores`, `utils`)

## Where to Add New Code

**New Backend Feature (end-to-end):**
1. Domain model: `backend/internal/core/domain/<feature>.go`
2. Repository port: `backend/internal/core/ports/repositories/<feature>_repository.go`
3. Service: `backend/internal/core/services/<feature>_service.go`
4. Repository impl: `backend/internal/adapters/repositories/postgres/gorm_<feature>_repository.go`
5. Schema: edit `backend/schema.graphql`, run `make graphql-gen`
6. Resolver: `backend/internal/adapters/graphql/resolvers/schema.resolvers.go` (fill in generated stub)
7. Wire: construct + inject in `backend/cmd/server/main.go`
8. Tests: `backend/test/services/`, `backend/test/repositories/` (or matching concern directory)

**New Frontend Feature/Component:**
- Component: `frontend/src/lib/components/<Feature>.svelte` (or a feature subdirectory like `discover/` if it's page-scoped)
- GraphQL query/mutation: add to relevant file in `frontend/src/lib/queries/` (or new file per entity), add cache keys to `frontend/src/lib/queries/keys.ts`
- Data-fetching hook: `frontend/src/lib/queries/hooks/use<Feature>.ts`
- Route: `frontend/src/routes/<feature>/+page.svelte` + `+page.ts`

**Utilities:**
- Backend shared helpers: `backend/pkg/` (framework-agnostic) or `backend/internal/adapters/graphql/resolvers/helpers.go` (resolver-specific)
- Frontend shared helpers: `frontend/src/lib/utils/` (pure functions), `frontend/src/lib/utils.ts` (class/cn helpers)

## Special Directories

**`backend/internal/adapters/graphql/generated/`:**
- Purpose: gqlgen-generated resolver interfaces, models, executable schema
- Generated: Yes (via `make graphql-gen`)
- Committed: Yes (checked in per gqlgen convention, but never hand-edited)

**`backend/coverage/`, `frontend/coverage/`:**
- Purpose: Generated test coverage reports
- Generated: Yes
- Committed: No (build artifact)

**`frontend/.svelte-kit/`, `frontend/build/`:**
- Purpose: SvelteKit build output
- Generated: Yes
- Committed: No

**`frontend/android/`, `frontend/ios/`:**
- Purpose: Capacitor native shell projects for mobile packaging
- Generated: Partially (scaffolded by Capacitor CLI, then customized)
- Committed: Yes

**`.planning/`:**
- Purpose: Legacy GSD planning artifacts (phases, roadmap, this codebase-mapping output)
- Generated: Mixed (some hand-written, some generated by GSD commands)
- Committed: Yes, but new work should use `docs/superpowers/` instead — see root `CLAUDE.md`

**`graphify-out/`:**
- Purpose: Knowledge graph of the codebase (`graph.json`, wiki, per-date snapshots) used by the `graphify` tool for code navigation
- Generated: Yes (`graphify update .`)
- Committed: Appears versioned by date subdirectories — treat as generated cache, do not hand-edit

---

*Structure analysis: 2026-09-04*
