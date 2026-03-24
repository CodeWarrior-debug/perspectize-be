---
phase: 09-security-hardening
plan: 01
subsystem: auth
tags: [jwt, golang-jwt, middleware, graphql-directives, httponly-cookie, hs256]

# Dependency graph
requires:
  - phase: 07-backend-architecture
    provides: chi router, middleware stack, hexagonal architecture
provides:
  - JWT generation/validation service (AuthService)
  - HTTP auth middleware with httpOnly cookie extraction
  - GraphQL @auth and @owner directives
  - Security config with JWT_SECRET from env
  - All mutations protected by @auth directive
affects: [09-02, 09-03, 09-04, 09-05, 09-06, 10-user-registration]

# Tech tracking
tech-stack:
  added: [golang-jwt/jwt v5, golang.org/x/crypto]
  patterns: [middleware-authenticates-directives-authorize, httponly-cookie-jwt, port-interface-for-auth-service]

key-files:
  created:
    - backend/internal/core/domain/auth.go
    - backend/internal/core/ports/services/auth_service.go
    - backend/internal/core/services/auth_service.go
    - backend/internal/adapters/web/middleware/auth.go
    - backend/internal/adapters/graphql/directives/auth.go
    - backend/internal/config/security.go
  modified:
    - backend/cmd/server/main.go
    - backend/schema.graphql
    - backend/go.mod
    - backend/go.sum
    - backend/.env.example
    - backend/test/resolvers/content_resolver_test.go
    - backend/internal/adapters/graphql/generated/generated.go

key-decisions:
  - "HS256 JWT signing with golang-jwt/jwt v5 (maintained successor to dgrijalva/jwt-go)"
  - "Middleware authenticates, directives authorize pattern - middleware never blocks unauthenticated requests"
  - "Dev fallback JWT secret for local development, fail-fast in production if missing or <32 bytes"
  - "ForContext returns (userID, bool) for simple resolver/directive extraction"

patterns-established:
  - "Auth middleware → directive pattern: middleware stores claims in context, directives enforce auth on fields"
  - "Port interface for AuthService: enables mock auth service in tests"
  - "Test auth pattern: mockAuthService + auth_token cookie for resolver tests with @auth directives"

requirements-completed: []

# Metrics
duration: 6min
completed: 2026-03-02
---

# Phase 09 Plan 01: JWT Authentication Infrastructure Summary

**JWT auth with HS256 signing via golang-jwt/jwt v5, httpOnly cookie middleware, and GraphQL @auth/@owner directives protecting all mutations**

## Performance

- **Duration:** 6 min
- **Started:** 2026-03-03T03:21:54Z
- **Completed:** 2026-03-03T03:27:41Z
- **Tasks:** 7
- **Files modified:** 13

## Accomplishments
- JWT generation and validation service with HS256 signing, issuer claim, and configurable TTL
- HTTP auth middleware extracts JWT from httpOnly cookie and stores user context
- GraphQL @auth and @owner directives enforce authentication on all 8 mutations
- Security config loads JWT_SECRET from environment with production validation (>=32 bytes)
- Auth infrastructure wired in server startup (middleware + directives + config)

## Task Commits

Each task was committed atomically:

1. **Task 1: Install JWT dependencies and create auth domain models** - `3fe64c0` (feat)
2. **Task 2: Create AuthService for JWT generation and validation** - `84347fd` (feat)
3. **Task 3: Create authentication middleware for HTTP layer** - `3ee0fe0` (feat)
4. **Task 4: Create GraphQL directives for authorization** - `0911de7` (feat)
5. **Task 5: Wire authentication in server startup** - `5bf3319` (feat)
6. **Task 6: Add @auth directive to protected mutations** - `fe79a7a` (feat)
7. **Task 7: Verify build and add .env.example entry** - `5cac409` (feat)

## Files Created/Modified
- `backend/internal/core/domain/auth.go` - Claims struct with UserID, Email, RegisteredClaims
- `backend/internal/core/ports/services/auth_service.go` - AuthService interface (GenerateAccessToken, ValidateToken)
- `backend/internal/core/services/auth_service.go` - HS256 JWT implementation with secret validation
- `backend/internal/adapters/web/middleware/auth.go` - Cookie-based auth middleware + ForContext helper
- `backend/internal/adapters/graphql/directives/auth.go` - @auth and @owner directive implementations
- `backend/internal/config/security.go` - Security config loading from env vars
- `backend/cmd/server/main.go` - Auth service creation, middleware wiring, directive configuration
- `backend/schema.graphql` - Added directive declarations and @auth on all mutations
- `backend/.env.example` - JWT_SECRET and ACCESS_TOKEN_MINUTES documentation
- `backend/test/resolvers/content_resolver_test.go` - Mock auth service and cookie injection for tests

## Decisions Made
- Used HS256 (HMAC) signing method for simplicity; RS256 can be considered later for distributed verification
- Middleware passes through unauthenticated requests (no blocking) - directives enforce auth per-field
- Dev fallback secret auto-generated for local development; production requires explicit JWT_SECRET
- ForContext returns simple (int, bool) tuple rather than full claims for resolver ergonomics
- Owner directive initially only checks authentication; full ownership validation deferred to Plan 02

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 1 - Bug] Fixed resolver tests failing after @auth directive addition**
- **Found during:** Task 7 (Build verification)
- **Issue:** Existing resolver tests used bare Config{Resolvers: resolver} without directives, causing "directive auth is not implemented" errors on mutation tests
- **Fix:** Wired Auth/Owner directives in test server setup, added mockAuthService, included auth_token cookie in executeGraphQL requests
- **Files modified:** backend/test/resolvers/content_resolver_test.go
- **Verification:** All tests pass (go test ./...)
- **Committed in:** 5cac409 (Task 7 commit)

---

**Total deviations:** 1 auto-fixed (1 bug fix)
**Impact on plan:** Essential fix - tests would fail without directive wiring. No scope creep.

## Issues Encountered
None beyond the test fix documented above.

## User Setup Required
None - JWT_SECRET defaults to an insecure value in development. For production, set:
- `JWT_SECRET` (minimum 32 bytes)
- `ACCESS_TOKEN_MINUTES` (optional, defaults to 15)

## Next Phase Readiness
- Auth infrastructure complete and ready for Plan 02 (ownership checks, rate limiting)
- No login/register endpoints yet (deferred to Phase 10)
- All mutations protected - frontend will need to handle "access denied" responses

---
*Phase: 09-security-hardening*
*Completed: 2026-03-02*
