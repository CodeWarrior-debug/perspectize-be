# Codebase Concerns

**Analysis Date:** 2026-02-16

## Tech Debt

**YouTube Client Lacks Dependency Injection for Testing:**
- Issue: `youtube.Client` constructor (`NewClient`) creates a hardcoded `http.Client{}` with no way to inject a test client. Tests cannot mock YouTube API responses without hitting the real API.
- Files: `backend/internal/adapters/youtube/client.go` (lines 24-29)
- Impact: Cannot fully unit-test YouTube metadata fetching. Current tests are either integration tests or skip comprehensive coverage.
- Fix approach: Add `NewClientWithHTTPClient(apiKey string, httpClient *http.Client)` constructor or `NewClientWithBaseURL(apiKey string, baseURL string)` to allow `httptest.Server` injection. See TODO comment in `backend/test/youtube/client_test.go` (line 353).

**Hand-Rolled Cursor Pagination Instead of Library:**
- Issue: Cursor encoding/decoding implemented manually via helpers in `backend/internal/adapters/repositories/postgres/helpers.go`. Works but reinvents the wheel.
- Files: `backend/internal/adapters/repositories/postgres/helpers.go` (lines 90-155), `gorm_content_repository.go`, `gorm_perspective_repository.go`
- Impact: More code to maintain, higher risk of pagination bugs. Library would provide type-safe field specification and automatic keyset query building.
- Fix approach: Integrate `gorm-cursor-paginator` library (already imported, partially used). Replace manual cursor logic with paginator's cursor handling. Simplify `List` method queries. See FEATURE_BACKLOG.md (HIGH PRIORITY).

**Incomplete GraphQL Filter Schema:**
- Issue: `ContentFilter` input type in GraphQL schema has TODO comment indicating missing filters. Currently only supports `contentType`, but search and date range filtering are planned but not implemented.
- Files: `backend/internal/adapters/graphql/generated/generated.go` (line 912 comment in schema)
- Impact: Limited filtering capability in Activity Table. Workaround: client-side filtering over fetched data (see client-side pagination below).
- Fix approach: Extend `ContentFilter` input with `search`, `dateRange`, and other fields. Add resolver logic to apply filters in GORM repository queries.

**Sticky Header Color Token Workaround:**
- Issue: Sticky header uses `bg-white` hardcoded instead of semantic `bg-background` token because `--color-background` wasn't defined in design system.
- Files: `frontend/src/lib/components/Header.svelte` (changed from `bg-background` to `bg-white` in commit b42c457)
- Impact: Cosmetic; header works but doesn't respect theme color scheme. Will break if dark mode is added.
- Fix approach: Define complete color theme in `frontend/src/app.css` with all semantic tokens (`--color-background`, `--color-foreground`, `--color-border`, etc.). Then revert header to `bg-background`.

**Client-Side Pagination Prefetch Strategy:**
- Issue: GraphQL `ListContent` query uses `first: pageSize` (configurable) to prefetch data, but AG Grid only displays 10 per page. Total fetch size is not adaptive.
- Files: `frontend/src/lib/components/ActivityTable.svelte` (line ~62: `first: pageSize`)
- Impact: Scalability concern at 100+ items. When content exceeds prefetched amount, client-side filtering breaks (only operates on loaded data). No total count exposed to UI.
- Fix approach: (1) Expose `totalCount` from GraphQL query so UI knows total server-side content. (2) Make prefetch adaptive or implement true server-side pagination. (3) Switch to server-side row model in AG Grid for seamless remote pagination. See FEATURE_BACKLOG.md "Server-Side Sorting and Filtering".

## Known Bugs

**Slow COUNT(*) and JSONB ORDER BY Queries:**
- Symptoms: `SELECT count(*) FROM "content"` takes 200-400ms; JSONB path ORDER BY queries take 150-200ms. Full GraphQL request latency 395-541ms at ~50 rows.
- Files: `backend/internal/adapters/repositories/postgres/gorm_content_repository.go`, slow query logger in `backend/pkg/database/postgres.go`
- Trigger: Any pagination request that calls `List()` repository method with sorting by JSONB fields (viewCount, likeCount, publishedAt).
- Workaround: Queries are slow but correct; will degrade further at 1000+ rows without indexing.
- Root cause: GORM generates `SELECT count(*)` without filtering applied, and JSONB path extraction (`response->'items'->0->'statistics'->>'viewCount'`) has no index.

**Client-Side Filtering Doesn't Filter Full Dataset:**
- Symptoms: AG Grid's client-side sort/filter in ActivityTable only reorders/filters visible page. Filtering "Views" shows only the highest views on current page, not globally highest.
- Files: `frontend/src/lib/components/ActivityTable.svelte` (lines 52-79: query with `first: pageSize`, filter applied client-side)
- Trigger: Any sort/filter when total content exceeds current page size.
- Workaround: Currently acceptable because dataset is small (~50 items). Becomes obvious UX bug at 200+ items.
- Root cause: Server-side pagination + client-side sort/filter = filtering only what's loaded. Need server-side sort/filter parameters in GraphQL query.

**GraphQL Client Lacks Authentication Headers:**
- Symptoms: `graphqlClient` has empty `headers: {}` — no auth tokens, no CSRF protection. All users share the same unauthenticated client.
- Files: `frontend/src/lib/queries/client.ts` (line 15)
- Trigger: Currently not an issue because backend has no authentication enforcement, but will be critical when auth is added.
- Workaround: None — backend currently allows any request.
- Root cause: Auth architecture not yet designed. See FEATURE_BACKLOG.md "Authentication Architecture Design".

**CORS Allows All Origins:**
- Symptoms: Backend CORS middleware sets `Access-Control-Allow-Origin: *`, allowing any domain to make requests.
- Files: `backend/cmd/server/main.go` (lines 125-136)
- Trigger: Any cross-origin request is accepted.
- Workaround: For development only; must be restricted before production.
- Root cause: Development convenience setting left in place. CLAUDE.md notes: "Restrict to frontend's production origin before deploying."

## Security Considerations

**Unrestricted GraphQL Introspection:**
- Risk: GraphQL introspection is publicly available, exposing the full schema to anyone. No authentication required to discover field names, types, and resolver structure.
- Files: `backend/internal/adapters/graphql/generated/generated.go` (gqlgen default config), `backend/cmd/server/main.go` (GraphQL handler setup)
- Current mitigation: None. Introspection is enabled by default in gqlgen.
- Recommendations: (1) Disable introspection in production via gqlgen config (`introspection = false` or via handler options). (2) Implement authentication and only allow introspection for authenticated users in development environments. (3) Consider rate limiting to GraphQL endpoint.

**No Input Validation on User-Supplied Text:**
- Risk: Perspective claims, user descriptions, and other text fields are stored and displayed without sanitization. If HTML rendering is added later, XSS is possible.
- Files: All mutation resolvers in `backend/internal/adapters/graphql/resolvers/schema.resolvers.go` accept unvalidated text input.
- Current mitigation: Go's template system (if used) escapes by default. Currently no HTML rendering frontend-side.
- Recommendations: (1) Add length limits to text fields (e.g., max 5000 chars for claims). (2) Validate and sanitize on backend before storage. (3) Use content security policies in frontend HTML. (4) Never render user HTML without sanitization library (e.g., DOMPurify).

**Rate Limiting Not Implemented:**
- Risk: No rate limiting on GraphQL endpoint. A malicious actor can spam requests, exhausting database connections or causing DoS.
- Files: `backend/cmd/server/main.go` (no rate limiting middleware present)
- Current mitigation: None. Database connection pooling provides some soft limit (25 max open connections by default).
- Recommendations: (1) Add rate limiting middleware per IP/user (e.g., `golang.org/x/time/rate`). (2) Implement tiered limits: strict for anonymous users, relaxed for authenticated users. (3) Log and monitor for abuse patterns.

**YouTube API Key Exposed in Logs:**
- Risk: YouTube API key is passed as constructor argument to `youtube.Client` without masking. If logged via slog, it could leak in error messages.
- Files: `backend/internal/adapters/youtube/client.go`, `backend/cmd/server/main.go` (line 96: logs "YOUTUBE_API_KEY is empty" but doesn't log the actual key)
- Current mitigation: Weak. Key is only mentioned in warning when empty; actual key is not logged directly.
- Recommendations: (1) Add API key masking in error messages (show only first 3 and last 3 chars). (2) Audit all slog calls to ensure no accidental key logging. (3) Use structured logging with redaction helpers for sensitive fields.

**System User (Sentinel) Can Be Modified But Guarded:**
- Risk: A sentinel user with ID -1 exists for operations without a real user (e.g., admin actions). UpdateUser checks for sentinel, but Delete operation must also be guarded.
- Files: `backend/internal/adapters/graphql/resolvers/schema.resolvers.go` (lines 92-113: `DeleteUser` resolver), `backend/internal/core/services/user_service.go`
- Current mitigation: `UpdateUser` includes sentinel check (line 82: `ErrSentinelUser`). Must verify `DeleteUser` service has same protection.
- Recommendations: (1) Ensure `DeleteUser` service method checks for sentinel user and rejects. (2) Add unit tests for sentinel protection on both Update and Delete. (3) Consider making sentinel ID system constant to prevent accidental deletion.

## Performance Bottlenecks

**JSONB Response Column Dominates Storage:**
- Problem: The `content.response` JSONB column stores full YouTube API responses and accounts for 93.7% of all content table data (118 KB of 126 KB at 49 rows).
- Files: `backend/migrations/000002_update_response_jsonb.up.sql` (creates response column), database queries that retrieve response
- Cause: YouTube API response includes unused fields (status, topicDetails, recordingDetails, etc.) alongside needed fields (snippet, statistics, contentDetails).
- Improvement path: (1) **Trim on ingest** — Store only JSONB paths the app reads (snippet.title, statistics.*, etc.), drop unused nested objects. (2) **Extract to columns** — Promote frequently queried JSONB paths (viewCount, likeCount, publishedAt, channelTitle) to dedicated columns for indexing and sort performance. (3) **Compress** — Use `pg_lz_compress` if full response audit trail needed.
- Priority: Low at current scale (49 rows, 8 MB DB). Revisit at 1000+ rows. Related to slow JSONB query issue above.

**N+1 Query Risk in Perspective Resolvers:**
- Problem: Resolving perspectives on content may fetch user info per perspective without batching (DataLoader).
- Files: `backend/internal/adapters/graphql/resolvers/schema.resolvers.go` (perspective resolvers)
- Cause: GraphQL resolvers resolve nested fields independently. At 10 perspectives, this could be 10 separate user queries.
- Improvement path: Implement DataLoader pattern in GraphQL resolver to batch user queries. Batch load all users referenced by a page of perspectives in one DB query.
- Priority: Low at current scale (small datasets). Becomes noticeable at 1000+ perspectives.

**Slow Pagination Count Query:**
- Problem: `SELECT count(*) FROM "content"` scans entire table and takes 200-400ms even with 50 rows. Will scale linearly with data.
- Files: GORM `List()` method in `gorm_content_repository.go`
- Cause: No indexes on frequently filtered columns; count query runs before filters applied.
- Improvement path: (1) Add B-tree indexes on frequently filtered columns (content_type). (2) Use estimated count for pagination (PostgreSQL `pg_stat_user_tables`) to avoid expensive exact counts. (3) Cache total count with TTL for pagination use case. (4) Use keyset pagination WITHOUT total count — many modern UIs don't show "Page X of Y", just "Next" button.
- Priority: Low now, High at 10K+ rows.

## Fragile Areas

**GraphQL Generated Code Contains Panics:**
- Files: `backend/internal/adapters/graphql/generated/generated.go` (lines 6006, 6094, 6185, 6233, etc.) — multiple `panic("unknown field ...")` statements in unmarshaling code.
- Why fragile: Auto-generated code is never edited directly. If schema changes break the unmarshaling logic, panics are not caught by tests. Schema validation should catch this, but panic is worse than returning error.
- Safe modification: (1) Do not edit generated file directly. (2) After schema changes, run `make graphql-gen` and test thoroughly. (3) Panic recovery middleware (`backend/pkg/middleware/recovery.go`) will catch panics and log them as errors.
- Test coverage: Panic recovery should catch panics. Verify middleware is wired in `backend/cmd/server/main.go` (line 122: `r.Use(perfmw.Recoverer)`).

**Activity Table Complex State Management:**
- Files: `frontend/src/lib/components/ActivityTable.svelte` (424 lines with cursor stack, page state, sort state, filter state, grid state)
- Why fragile: Multiple state variables (`cursors`, `currentPage`, `sortBy`, `sortOrder`, `filterText`) are tightly coupled. Changing pagination logic could break sort/filter behavior. Cursor stack management is error-prone.
- Safe modification: (1) Before changing pagination, write tests for next/prev page transitions with sorted data. (2) Document cursor stack management clearly (why `cursors[currentPage]` is used, how hasNextPage affects stack). (3) Consider extracting pagination logic into custom hook (e.g., `useCursorPagination`). (4) Add debug logging for cursor state during development.
- Test coverage: Likely no tests for ActivityTable. Add integration tests for pagination + sort combinations.

**Repository Filter Logic with String Conversion:**
- Files: `backend/internal/adapters/repositories/postgres/gorm_content_repository.go` (lines 94: `strings.ToLower(string(*params.Filter.ContentType))`)
- Why fragile: Enum-to-string conversion happens in repository. If domain enum changes, SQL queries will silently break (database stores lowercase, domain is uppercase by convention). No validation that converted string is valid enum value.
- Safe modification: (1) Use enum conversion helpers (`contentTypeToDBValue`, `contentTypeFromDBValue` in `helpers.go`). (2) Add tests that verify round-trip conversion: `domainEnum -> dbString -> domainEnum` should be identity.
- Test coverage: Check if tests exist for enum conversions. If not, add tests in test suite.

**Hard-Coded Query Parameters:**
- Files: `frontend/src/lib/components/ActivityTable.svelte` (lines 52-79: query with dynamic `first: pageSize`)
- Why fragile: If a request fails silently with "limit exceeded" error from backend, client won't know. No error boundaries around fetch.
- Safe modification: (1) Add error handling in ActivityTable's `queryFn` to log and display fetch failures. (2) Validate `pageSize` against server limit before sending. (3) Add comments explaining why specific limits are chosen.
- Test coverage: Add tests for fetch failure scenarios.

## Scaling Limits

**Database Connection Pool Scaling:**
- Current capacity: 25 max open connections (default in `backend/pkg/database/postgres.go`, line 26)
- Limit: At ~100 concurrent users with keep-alive, connection pool may saturate. Neon (PostgreSQL provider) has connection limits (e.g., 100 connections on free tier).
- Scaling path: (1) Increase `DB_MAX_OPEN_CONNS` env var (configurable, read in `backend/pkg/database/postgres.go` lines 36-39). (2) Implement connection pooling middleware (PgBouncer or similar) to multiplex connections. (3) For Fly.io deployment with Neon, consider managed connection pooling. (4) Optimize query duration to reduce connection hold time.

**In-Memory GraphQL Query Caching:**
- Current capacity: No backend cache. TanStack Query uses browser memory for `staleTime: 60 * 1000` (60-second cache per query).
- Limit: Browser memory limit ~50-100 MB for cached responses. At 1000+ items with large JSONB responses, frontend memory usage grows linearly.
- Scaling path: (1) Implement Redis cache on backend for hot data (frequently accessed content). (2) Implement incremental/subscription queries to reduce payload per request. (3) Compress responses (gzip already done by HTTP layer). (4) Implement pagination + infinite scroll instead of loading all items at once.

**CORS Wildcard at Scale:**
- Current capacity: Unlimited. Any origin can request.
- Limit: If frontend domain is compromised or a competitor tries to abuse the API, requests are not rejected.
- Scaling path: (1) Restrict CORS to specific frontend origin before production. (2) Implement per-IP rate limiting. (3) Add API key or JWT authentication (planned in Phase 9).

## Dependencies at Risk

**GORM Version Pinning:**
- Risk: GORM is deeply integrated into repository layer. If a security issue is found, upgrading GORM could introduce breaking changes.
- Impact: Forced to stay on vulnerable version if not managed carefully.
- Migration plan: (1) Keep GORM up-to-date with minor/patch versions regularly (no cost). (2) For major versions, plan a refactor cycle (e.g., GORM v2 to v3). (3) Alternatively, migrate to sqlc (type-safe SQL generation) to reduce ORM dependency.
- Priority: Low (GORM is stable). Monitor security advisories.

**gqlgen Code Generation Coupling:**
- Risk: GraphQL resolvers are tightly coupled to gqlgen-generated types (model.Content, model.User, etc.). If gqlgen version changes, type definitions could shift.
- Impact: Potential breaking changes on `make graphql-gen`.
- Migration plan: (1) Keep gqlgen version pinned (currently v0.17.86). (2) Before upgrading, run full test suite (`go test ./...`) to catch generated code changes. (3) Consider wrapping generated types in domain types to decouple resolvers.
- Priority: Low (gqlgen is stable). Currently pinned.

**No Automated Dependency Security Scanning:**
- Risk: Go modules could have vulnerabilities. No Dependabot or similar checking for security updates.
- Impact: Could ship with known vulnerable dependencies.
- Mitigation: (1) Enable Dependabot on GitHub (auto-creates PRs for updates). (2) Run `go mod tidy && go mod verify` before commits. (3) Periodically audit with `go list -u -m all`.
- Priority: Medium (standard practice).

## Missing Critical Features

**Authentication Not Implemented:**
- Problem: Backend accepts any request without auth. Frontend has no user context. All data is public.
- Blocks: (1) User-scoped data (my perspectives vs. others' perspectives). (2) Secure content creation (prevent spoofing). (3) Admin operations. (4) Data privacy compliance.
- Mitigation: Currently acceptable for MVP. Phase 9 (Security Hardening) will address. See FEATURE_BACKLOG.md "Authentication Architecture Design".

**Authorization / Permissions:**
- Problem: No role-based access control (RBAC). No way to prevent users from deleting others' content or modifying system data.
- Blocks: Multi-user features (collaboration, moderation, content ownership).
- Mitigation: Sentinel user mechanism exists (ErrSentinelUser) but comprehensive RBAC is missing.

**Error Boundaries in Frontend:**
- Problem: GraphQL query errors in ActivityTable are not caught or displayed to user. If fetch fails, user sees no data with no error message.
- Blocks: Good error UX (showing "Error loading content" vs. silent failure).
- Mitigation: Add error UI to ActivityTable to display `query.isError` and `query.error`.

## Test Coverage Gaps

**YouTube Client Integration:**
- What's not tested: Comprehensive YouTube API mocking. Current tests likely skip due to lack of dependency injection.
- Files: `backend/test/youtube/client_test.go` (has TODO at line 353 for refactoring)
- Risk: YouTube integration could have bugs that only surface in production or with real API changes.
- Priority: High — YouTube integration is core feature.

**Pagination Edge Cases:**
- What's not tested: Page boundary conditions (last page with fewer items, cursor validity), sort + pagination interaction, filter + pagination interaction.
- Files: `backend/test/services/`, repository test files
- Risk: Pagination could have off-by-one errors or lose data at page boundaries.
- Priority: High — pagination is visible to all users.

**ActivityTable Interaction Tests:**
- What's not tested: Sort by clicking headers, filter by typing, page navigation, grid state persistence.
- Files: `frontend/src/lib/components/ActivityTable.svelte`
- Risk: Frontend bugs only caught by manual testing. Regressions on refactor.
- Priority: Medium — requires browser testing (Playwright or Cypress).

**Error Handling Paths:**
- What's not tested: Service layer error conditions (validation failures, missing resources). Resolvers catch errors but may not handle all cases.
- Files: `backend/internal/adapters/graphql/resolvers/`, `backend/internal/core/services/`
- Risk: Unhandled errors could leak implementation details or fail silently.
- Priority: Medium — covered by integration tests but not isolated unit tests.

---

*Concerns audit: 2026-02-16*
