---
phase: 09-security-hardening
plan: 03
subsystem: api
tags: [rate-limiting, cors, csrf, introspection, graphql-security, httprate]

requires:
  - phase: 09-01
    provides: "JWT auth infrastructure with middleware and directives"
provides:
  - "Rate limiting middleware (httprate sliding window)"
  - "CORS restriction with configurable origins"
  - "Content-Type validation for CSRF protection"
  - "Query complexity limit (FixedComplexityLimit 500)"
  - "Introspection disabled in production"
affects: [deployment, frontend-cors]

tech-stack:
  added: [go-chi/httprate v0.15.0, go-chi/cors v1.2.2]
  patterns: [handler.New with selective extensions, middleware ordering for security]

key-files:
  created:
    - backend/internal/adapters/web/middleware/ratelimit.go
    - backend/internal/adapters/web/middleware/contenttype.go
    - backend/internal/adapters/web/middleware/contenttype_test.go
  modified:
    - backend/cmd/server/main.go
    - backend/internal/config/security.go
    - backend/.env.example

key-decisions:
  - "handler.New instead of NewDefaultServer for selective extension control (introspection toggle)"
  - "Rate limiting before auth in middleware stack to prevent auth DoS"
  - "Content-Type validation as CSRF protection (simpler than gorilla/csrf for API-only GraphQL)"
  - "Wildcard CORS default for dev, explicit origins required in production"

patterns-established:
  - "Middleware ordering: RequestID -> RealIP -> RateLimit -> CORS -> ContentType -> Auth -> Timer -> Recoverer"
  - "getEnvStringSlice helper for comma-separated env var lists"

requirements-completed: []

duration: 4min
completed: 2026-03-02
---

# Phase 09 Plan 03: API Protection Layer Summary

**Rate limiting (httprate), CORS restriction, Content-Type CSRF protection, query complexity limit (500), and introspection control via handler.New with selective gqlgen extensions**

## Performance

- **Duration:** 4 min
- **Started:** 2026-03-03T03:30:12Z
- **Completed:** 2026-03-03T03:34:22Z
- **Tasks:** 7
- **Files modified:** 8

## Accomplishments
- Rate limiting middleware with sliding window counter (httprate) wired before auth
- CORS restricted to configurable origins via go-chi/cors (C-05)
- Content-Type validation rejects non-JSON POSTs for CSRF protection (M-15)
- Query complexity limited to 500 via FixedComplexityLimit (C-04)
- Introspection disabled in production via selective extension loading (C-10)
- 5 unit tests for Content-Type validation middleware

## Task Commits

Each task was committed atomically:

1. **Task 1: Install dependencies** - `5ad4c9d` (chore)
2. **Task 2: Rate limiting middleware** - `585972a` (feat)
3. **Task 3: Content-Type validation middleware** - `0e317e9` (feat)
4. **Task 4: Security config update** - `93684d5` (feat)
5. **Task 5: Wire API protection layers** - `fc18620` (feat)
6. **Task 6: Update .env.example** - `b0805c1` (docs)
7. **Task 7: Tests and verification** - `b7514c0` (test)

## Files Created/Modified
- `backend/internal/adapters/web/middleware/ratelimit.go` - GlobalRateLimit and GraphQLRateLimit functions
- `backend/internal/adapters/web/middleware/contenttype.go` - ContentTypeValidation middleware
- `backend/internal/adapters/web/middleware/contenttype_test.go` - 5 unit tests for Content-Type validation
- `backend/cmd/server/main.go` - Wired all protection layers, switched to handler.New
- `backend/internal/config/security.go` - Added RateLimitPerMin, CORSOrigins, helpers
- `backend/.env.example` - RATE_LIMIT_PER_MIN and CORS_ORIGINS documentation
- `backend/go.mod` - httprate v0.15.0, cors v1.2.2
- `backend/go.sum` - Dependency checksums

## Decisions Made
- Used `handler.New` instead of `handler.NewDefaultServer` to selectively include gqlgen extensions (enables introspection toggle)
- Rate limiting placed before auth middleware to prevent DoS on authentication endpoint
- Content-Type validation chosen over gorilla/csrf for CSRF protection (simpler for API-only GraphQL)
- Default CORS wildcard for development, explicit origins required in production via CORS_ORIGINS env

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 1 - Bug] Fixed NewDirectiveRoot missing service arguments**
- **Found during:** Task 7 (Build verification)
- **Issue:** `directives.NewDirectiveRoot()` was called without required `contentService` and `perspectiveService` arguments (pre-existing from 09-01)
- **Fix:** Added service arguments to both `cmd/server/main.go` and `test/resolvers/content_resolver_test.go`
- **Files modified:** `backend/cmd/server/main.go`, `backend/test/resolvers/content_resolver_test.go`
- **Verification:** `go build ./cmd/server/` and `go test ./...` pass
- **Committed in:** `b7514c0` (Task 7 commit)

---

**Total deviations:** 1 auto-fixed (1 bug fix)
**Impact on plan:** Bug fix necessary for build correctness. No scope creep.

## Issues Encountered
- go-chi/cors v1.2.2 needed re-installation after initial `go get` did not persist in go.mod (resolved by running `go get` again after import was added to source)

## User Setup Required
None - no external service configuration required. CORS_ORIGINS defaults to wildcard for development. Production deployment should set `CORS_ORIGINS` to explicit frontend URL.

## Next Phase Readiness
- API protection layer complete (C-04, C-05, C-10, H-11, M-15)
- Ready for remaining security hardening plans (02-06)
- Production deployment should configure CORS_ORIGINS and RATE_LIMIT_PER_MIN env vars

---
*Phase: 09-security-hardening*
*Completed: 2026-03-02*
