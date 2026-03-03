---
phase: 09-security-hardening
plan: 02
subsystem: auth
tags: [authorization, ownership, graphql-directives, email-visibility]

# Dependency graph
requires:
  - phase: 09-01
    provides: JWT auth middleware, ForContext helper, @auth/@owner directive stubs
provides:
  - Full @owner directive with resource ownership checks via service.GetByID
  - User.email field resolver gated by authentication (H-10)
  - WithUserContext test helper for auth context setup
affects: [09-03, 09-04, 09-05, 09-06]

# Tech tracking
tech-stack:
  added: []
  patterns:
    - "gqlgen field resolver for auth-gated fields (User.email)"
    - "@owner directive with service injection for ownership validation"
    - "extractResourceID supports top-level args and nested input objects"
    - "WithUserContext exported helper for test auth context"

key-files:
  created:
    - backend/internal/adapters/graphql/directives/auth_test.go
  modified:
    - backend/internal/adapters/graphql/directives/auth.go
    - backend/internal/adapters/web/middleware/auth.go
    - backend/internal/adapters/graphql/resolvers/schema.resolvers.go
    - backend/internal/adapters/graphql/resolvers/helpers.go
    - backend/schema.graphql
    - backend/gqlgen.yml
    - backend/internal/adapters/graphql/generated/generated.go
    - backend/internal/adapters/graphql/model/models_gen.go

key-decisions:
  - "gqlgen field resolver for User.email instead of context-threaded helper — cleaner separation"
  - "email field made nullable (String! -> String) to return null when not authorized"
  - "WithUserContext exported from middleware for test reuse across packages"
  - "@owner directive extracts ID from both top-level args and nested input objects"

patterns-established:
  - "Field resolver pattern: gqlgen field resolver for auth-gated fields"
  - "Ownership check pattern: @owner directive + service.GetByID + field name routing"

requirements-completed: []

# Metrics
duration: 7min
completed: 2026-03-03
---

# Phase 09 Plan 02: Authorization & User Data Protection Summary

**@owner directive with service-backed ownership checks and email visibility gating via gqlgen field resolver**

## Performance

- **Duration:** 7 min
- **Started:** 2026-03-03T03:30:02Z
- **Completed:** 2026-03-03T03:37:19Z
- **Tasks:** 5 (3 committed, 2 already done from parallel Plan 03 work)
- **Files modified:** 8

## Accomplishments
- User.email field returns null unless authenticated user requests their own account (H-10)
- @owner directive performs real ownership checks via contentService.GetByID and perspectiveService.GetByID
- updatePerspective and deletePerspective mutations protected by @owner directive
- 10 unit tests covering auth and ownership directive behavior

## Task Commits

Each task was committed atomically:

1. **Task 1: Add email visibility check to User resolver** - `adeb75d` (feat)
2. **Task 2: Implement full @owner directive with resource ownership checks** - `f2d267f` (feat)
3. **Task 3: Add @owner directive to update/delete mutations** - included in `adeb75d` (schema changes bundled with Task 1)
4. **Task 4: Wire DirectiveRoot with service dependencies** - already wired by parallel Plan 03 execution
5. **Task 5: Verify build and test authorization** - `392fba4` (test)

## Files Created/Modified
- `backend/schema.graphql` - email nullable, @owner on perspective mutations
- `backend/gqlgen.yml` - User.email field resolver config
- `backend/internal/adapters/graphql/directives/auth.go` - Full @owner with service injection
- `backend/internal/adapters/graphql/directives/auth_test.go` - 10 unit tests for auth directives
- `backend/internal/adapters/web/middleware/auth.go` - Added WithUserContext test helper
- `backend/internal/adapters/graphql/resolvers/schema.resolvers.go` - Email field resolver with auth check
- `backend/internal/adapters/graphql/resolvers/helpers.go` - Email stored for field resolver access
- `backend/internal/adapters/graphql/generated/generated.go` - Regenerated with User field resolver
- `backend/internal/adapters/graphql/model/models_gen.go` - Email now *string (nullable)

## Decisions Made
- Used gqlgen field resolver for User.email (configured in gqlgen.yml) instead of threading context through userDomainToModel — cleaner gqlgen-native pattern
- Changed email from String! to String in schema to support returning null when unauthorized
- Added WithUserContext to middleware package as exported function for test reuse
- @owner directive uses extractResourceID helper that supports both top-level args (deletePerspective) and nested input objects (updatePerspective)

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 3 - Blocking] Schema email field changed to nullable**
- **Found during:** Task 1 (Email visibility check)
- **Issue:** email: String! cannot return null; gqlgen field resolver needs nullable field
- **Fix:** Changed to email: String in schema, regenerated code
- **Files modified:** backend/schema.graphql, generated code
- **Verification:** Build passes, field resolver returns nil for unauthorized requests

**2. [Rule 3 - Blocking] No deleteContent mutation in schema**
- **Found during:** Task 3 (Add @owner to mutations)
- **Issue:** Plan specified adding @owner to deleteContent but no such mutation exists
- **Fix:** Skipped deleteContent, applied @owner only to existing perspective mutations
- **Files modified:** backend/schema.graphql
- **Verification:** grep confirms @owner on updatePerspective and deletePerspective

**3. [Rule 3 - Blocking] Task 4 already completed by parallel execution**
- **Found during:** Task 4 (Wire DirectiveRoot)
- **Issue:** cmd/server/main.go already had NewDirectiveRoot(contentService, perspectiveService) from Plan 03 parallel execution
- **Fix:** No changes needed — verified wiring was correct
- **Files modified:** None
- **Verification:** grep confirms NewDirectiveRoot call with service params

---

**Total deviations:** 3 auto-fixed (3 blocking)
**Impact on plan:** All auto-fixes necessary for correctness. No scope creep.

## Issues Encountered
None beyond the documented deviations.

## User Setup Required
None - no external service configuration required.

## Next Phase Readiness
- Authorization layer complete with ownership checks
- Email visibility properly gated by authentication
- Test infrastructure established for directive testing
- Ready for Plan 03+ security hardening work

---
*Phase: 09-security-hardening*
*Completed: 2026-03-03*
