# Known Bugs

Issues discovered during development, testing, and codebase review (2026-02-07).

## Critical

| ID | Stack | Summary | Location | Source |
|----|-------|---------|----------|--------|
| C-01 | Go | No authentication or authorization — any client can CRUD any user's data | `cmd/server/main.go` | Arch, Code, Security, Silent |
| C-02 | Go | Cursor pagination broken for non-ID sort columns — keyset cursor only encodes `id`, produces wrong pages when sorting by name/date | `content_repository.go:207-336`, `perspective_repository.go:233-362` | Arch |
| C-03 | Go | No GraphQL query complexity or depth limiting — DoS vector | `cmd/server/main.go:75` | Code, Security |
| C-04 | Go | Wildcard CORS `Access-Control-Allow-Origin: *` allows any origin | `cmd/server/main.go:80` | Code, Security |
| C-05 | FE | XSS via innerHTML in AG Grid cellRenderer — `params.data.name` and `params.data.url` interpolated without escaping | `ActivityTable.svelte:64-70` | Arch, Code, Security, Silent |
| C-06 | Go | Silent JSON unmarshal failure drops categorized ratings — corrupted data silently omitted from responses | `perspective_repository.go:419-426` | Code, Silent |
| C-07 | Go | Silent duration parse failure defaults to 0 — indistinguishable from real zero-length video | `youtube/client.go:90-93` | Code, Silent |
| C-08 | Go | 5 silent parse failures in `domainToModel` — response, viewCount, likeCount, commentCount all silently nil on error | `resolvers/helpers.go:36-63` | Code, Silent |
| C-09 | Go | GraphQL Playground exposed unconditionally (no env check) | `cmd/server/main.go:92` | Security |
| C-10 | Go | GraphQL introspection enabled without restriction | `cmd/server/main.go:75` | Security |

## High

| ID | Stack | Summary | Location | Source |
|----|-------|---------|----------|--------|
| H-01 | Go | Resolver imports YouTube adapter directly — violates hexagonal dependency rule (adapter-to-adapter coupling) | `schema.resolvers.go:16,23` | Arch |
| H-02 | Go | Resolver depends on concrete service types, not interfaces — missing service port interfaces | `resolvers/resolver.go:12-16` | Arch |
| H-03 | Go | `ListAll()` users has no pagination — unbounded query result set | `user_repository.go:98-114` | Arch, Code, Security |
| H-04 | Go | GraphQL timestamps as `String!` instead of `DateTime` scalar — weak API contract | `schema.graphql:9-10,57-58,77-78` | Arch |
| H-05 | Go | `contentType` output uses `String!` instead of defined `ContentType` enum | `schema.graphql:68` | Arch |
| H-06 | Go | Race condition on duplicate claim check (read-then-write without transaction) | `perspective_service.go:91-97` | Code |
| H-07 | Go | Race condition on user uniqueness check (same read-then-write pattern) | `user_service.go:49-65` | Code |
| H-08 | Go | Entire YouTube API response stored verbatim in DB — bloat risk | `youtube/client.go:100` | Code |
| H-09 | Go | Hardcoded config path `config/config.example.json` | `cmd/server/main.go:27` | Code |
| H-10 | Go | User email addresses exposed without access control via `users` query | `schema.resolvers.go:302-315` | Security |
| H-11 | Go | No rate limiting on any endpoint | `cmd/server/main.go` | Security |
| H-12 | Go | YouTube API key exposure risk in URLs, error messages, and stored responses | `youtube/client.go:53-57,76` | Security |
| H-13 | Go | Sensitive data (DB schema info) leaked in GraphQL error messages via `%w` wrapping | `schema.resolvers.go:31,47,99,152,173` | Security |
| H-14 | Go | No HTTPS/TLS — `http.ListenAndServe` only | `cmd/server/main.go:99` | Security |
| H-15 | Go | No HTTP server timeouts (ReadTimeout, WriteTimeout, IdleTimeout) — Slowloris vector | `cmd/server/main.go:99` | Security |
| H-16 | Go | `nil, nil` returns for not-found hide errors from GraphQL clients (inconsistent with ContentByID) | `schema.resolvers.go:274,289,326` | Silent |
| H-17 | FE | No `+error.svelte` error boundary — unhandled errors show ugly default page | `src/routes/` (missing) | Silent |
| H-18 | FE | No `hooks.client.ts` or `hooks.server.ts` — errors outside TanStack invisible | `src/` (missing) | Silent |
| H-19 | Go | `.env` load failure silently ignored (`_ = godotenv.Load()`) | `cmd/server/main.go:24` | Silent |
| H-20 | Go | Empty YouTube API key not validated at startup — fails with cryptic 403 at runtime | `config.go`, `youtube/client.go:23` | Silent |
| H-21 | Go | `io.WriteString` return value ignored in IntID marshaler — response corruption risk | `intid.go:17` | Silent |
| H-22 | FE | `prerender = true` with `adapter-static` — SvelteKit used as SPA wrapper, no SSR | `+layout.ts:1` | Arch |
| H-23 | FE | No TypeScript types generated from GraphQL schema — manual duplication, drift risk | `+page.svelte`, `ActivityTable.svelte`, `UserSelector.svelte` | Arch, Code |
| H-24 | FE | GraphQL client has no error interceptor, timeout, or auth header infrastructure | `queries/client.ts:1-7` | Arch, Silent |
| H-25 | FE | No Content Security Policy | `app.html` | Security |
| H-26 | Both | No CI/CD pipeline or automated security scanning | `.github/` | Security |

## Medium

| ID | Stack | Summary | Location | Source |
|----|-------|---------|----------|--------|
| M-01 | Go | Dual PostgreSQL driver dependencies (`lib/pq` and `pgx/v5`) | `go.mod:10-11` | Arch |
| M-02 | Go | Hardcoded connection pool settings (25 max open, 5 idle) | `pkg/database/postgres.go:21-23` | Arch |
| M-03 | Go | `CreateFromYouTube` returns `ErrAlreadyExists` instead of idempotent return | `content_service.go:30-36` | Arch |
| M-04 | Go | `deletePerspective` uses `ID` scalar instead of `IntID` (inconsistent) | `schema.graphql:185` | Arch |
| M-05 | Go | `CreateFromYouTube` accepts extractVideoID as function param instead of constructor injection | `content_service.go:28` | Arch |
| M-06 | Go | No request logging or middleware chain — uses `net/http` default mux, not chi | `cmd/server/main.go:91-94` | Arch |
| M-07 | Go | Inconsistent not-found error handling across resolvers | `schema.resolvers.go:188-189 vs 274 vs 327` | Code |
| M-08 | Go | Nested GraphQL field resolvers missing — `user`/`content` on Perspective silently return null | `resolvers/helpers.go:70-107` | Code |
| M-09 | Go | No graceful shutdown handler | `cmd/server/main.go:99` | Code |
| M-10 | Go | No health check endpoint (`/health`, `/ready`) | `cmd/server/main.go:91-93` | Code, Security |
| M-11 | Go | Missing input length validation on description, labels[], parts[], categorizedRatings[] | `perspective_service.go`, `user_service.go` | Security |
| M-12 | Go | DB credentials may appear in logs on connection failure | `cmd/server/main.go:43-44`, `config.go:83` | Security |
| M-13 | Go | Unbounded `response: JSON` field exposes full YouTube API response per item | `schema.graphql:77` | Security |
| M-14 | Go | No security headers (X-Content-Type-Options, X-Frame-Options, HSTS, Cache-Control) | `cmd/server/main.go` | Security |
| M-15 | Go | No CSRF protection (moot with CORS wildcard, relevant once fixed) | `cmd/server/main.go:93` | Security |
| M-16 | Go | Perspective Update does not check RowsAffected (unlike Delete) — TOCTOU race | `perspective_repository.go:187-209` | Silent |
| M-17 | Go | No DATABASE_URL format validation | `config.go:79-84` | Silent |
| M-18 | FE | Duplicated type definitions across components (ContentItem, ContentRow, User) | `+page.svelte`, `ActivityTable.svelte`, `UserSelector.svelte` | Arch, Code |
| M-19 | FE | Content query always fetches 100 items — no server-side pagination integration with AG Grid | `+page.svelte:33-34` | Arch, Code |
| M-20 | FE | `selectedUserId` store wired but never consumed by any query | `userSelection.svelte.ts`, `+page.svelte` | Arch, Code |
| M-21 | FE | `ContentResponse`/`ContentItem` interfaces declared but never used as type guards | `+page.svelte:8-28` | Code |
| M-22 | FE | Search input not debounced — AG Grid filtering on every keystroke | `+page.svelte:30`, `ActivityTable.svelte:130-133` | Code |
| M-23 | FE | No error recovery UI (retry button) on error states | `+page.svelte:70-73`, `UserSelector.svelte:37-40` | Code, Silent |
| M-24 | FE | `AGGridTest.svelte` dead code in production component tree | `src/lib/components/AGGridTest.svelte` | Code |
| M-25 | FE | GraphQL endpoint fallback uses HTTP not HTTPS | `queries/client.ts:3` | Security |
| M-26 | FE | `retry: 1` retries all errors including 4xx (should only retry network/5xx) | `+layout.svelte:15`, `+page.svelte:40` | Silent |
| M-27 | FE | `formatDate` silently produces "Invalid Date" string for bad input | `ActivityTable.svelte:48-53` | Silent |
| M-28 | Both | No secret rotation or vault integration | - | Security |

## Low

| ID | Stack | Summary | Location | Source |
|----|-------|---------|----------|--------|
| L-01 | Go | `log` package used instead of `slog` (docs say structured logging) | `cmd/server/main.go` | Arch, Code |
| L-02 | Go | Tests in external `test/` directory instead of co-located `_test.go` | `perspectize-go/test/` | Arch |
| L-03 | Go | `config.example.json` used as runtime config name | `cmd/server/main.go:27` | Arch |
| L-04 | Go | No database indexes on `user_id`/`content_id` foreign keys in perspectives | `migrations/000004` | Arch |
| L-05 | Go | `content.name` UNIQUE constraint too restrictive — two videos can share a title | `migrations/000001:11` | Arch |
| L-06 | Go | Repeated `strconv.Atoi` pattern in resolvers — `IntID` scalar only used for inputs | `schema.resolvers.go:160,181,266,319` | Code |
| L-07 | Go | Duplicated pagination logic across content and perspective repositories | `content_repository.go`, `perspective_repository.go` | Code |
| L-08 | Go | Dead code: deprecated `toNullStringFromReviewStatus` | `perspective_repository.go:474-478` | Code |
| L-09 | Go | Duplicate `toNullInt64` / `toNullInt64FromIntPtr` functions | `content_repository.go:155`, `perspective_repository.go:440` | Code |
| L-10 | Go | Regex compiled on every `ExtractVideoID` call — should be package-level var | `youtube/parser.go:12-18` | Code |
| L-11 | Go | `*TEMP*` comment left in production code | `pkg/database/postgres.go:15` | Code |
| L-12 | Go | `docker-compose.yml` uses deprecated `version: '3.8'` | `docker-compose.yml:1` | Code |
| L-13 | Go | Sequential integer IDs enable enumeration (consider UUIDs for external identifiers) | `migrations/000004` | Security |
| L-14 | Go | Base64 cursor format trivially decodable — exposes internal IDs | `content_repository.go:163-182` | Security |
| L-15 | Go | No HTTP response size limit on YouTube client (`io.ReadAll`) | `youtube/client.go:70` | Security |
| L-16 | FE | `@tanstack/svelte-form` dependency unused | `package.json:38` | Arch |
| L-17 | FE | `utils/` directory contains only `.gitkeep` — utility functions inline in components | `src/lib/utils/.gitkeep` | Arch |
| L-18 | FE | `$lib/index.ts` empty — exports nothing | `src/lib/index.ts` | Code |
| L-19 | FE | `WithElementRef` type uses `any` for ref, class, children | `src/lib/utils.ts:8-12` | Code |
| L-20 | FE | Redundant `staleTime`/`retry` config — page duplicates layout defaults | `+page.svelte:39-40` vs `+layout.svelte:13-14` | Code |
| L-21 | FE | `parseInt` edge cases in UserSelector ID parsing | `UserSelector.svelte:26` | Code |
| L-22 | FE | Raw error messages displayed to users (may leak internal details) | `+page.svelte:71-73` | Security, Silent |

## UI Bugs (Phase 2.1)

| ID | Priority | Summary | Location | Found |
|----|----------|---------|----------|-------|
| BUG-001 | P1 | Header overflows at 375px — "Perspectize" truncated, "Add Video" button clipped | `Header.svelte` | 2026-02-07 |
| BUG-002 | P1 | Pagination bar broken at 375px — "Page Size:" truncated, count text wraps | `ActivityTable.svelte` (AG Grid) | 2026-02-07 |
| BUG-003 | P1 | Sticky header clipping persists on scroll at 375px | `Header.svelte` | 2026-02-07 |
| BUG-004 | P2 | No responsive header collapse — needs hamburger menu at mobile widths | `Header.svelte` | 2026-02-07 |
| BUG-005 | P2 | Table content left-shifted beyond viewport at 375px when scrolled down | `ActivityTable.svelte` | 2026-02-07 |
| BUG-006 | P3 | No visual affordance for hidden columns — no hint that horizontal scroll reveals more | `ActivityTable.svelte` | 2026-02-07 |

## Test Coverage Gaps

| ID | Priority | Summary | Location |
|----|----------|---------|----------|
| T-01 | P1 | `PerspectiveService.Update()` — 100-line complex mutation with zero tests | `perspective_service.go:163-265` |
| T-02 | P1 | No resolver tests for User/Perspective queries/mutations (only Content covered) | `schema.resolvers.go` |
| T-03 | P1 | No tests for `helpers.go` domain-to-model conversion (silent JSON parse) | `resolvers/helpers.go` |
| T-04 | P2 | No repository-layer tests (cursor encoding, SQL construction, sort whitelisting) | `content_repository.go`, `perspective_repository.go` |
| T-05 | P2 | No YouTube API client tests | `youtube/client.go` |
| T-06 | P2 | No tests for `IntID` scalar (handles 4 input types) | `pkg/graphql/intid.go` |
| T-07 | P2 | UserSelector component has no behavioral tests (only import check) | `UserSelector.test.ts` |
| T-08 | P3 | Mock duplication across 4 test files — needs shared `testutil/mocks.go` | `test/services/`, `test/resolvers/` |
| T-09 | P3 | Frontend query tests only check string content, not query validity | `queries-content.test.ts`, `queries-users.test.ts` |
| T-10 | P3 | Component tests overtest CSS classes (brittle) | various test files |
| T-11 | P3 | ActivityTable formatting functions untested (inline, should extract to utils) | `ActivityTable.svelte` |
| T-12 | P3 | Store `setSelectedUserId`/`clearUserSelection` not tested | `stores-userSelection.test.ts` |
