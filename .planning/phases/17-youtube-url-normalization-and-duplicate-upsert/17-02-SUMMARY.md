---
phase: 17-youtube-url-normalization-and-duplicate-upsert
plan: 02
subsystem: api
tags: [go, gqlgen, graphql, svelte, typescript, tanstack-query, toast, deduplication]

# Dependency graph
requires:
  - phase: 17-youtube-url-normalization-and-duplicate-upsert
    plan: 01
    provides: Idempotent CreateFromYouTube returning (content, ErrAlreadyExists) on duplicate

provides:
  - CreateContentResult GraphQL type wrapping Content and alreadyExisted boolean
  - Updated createContentFromYouTube mutation returning CreateContentResult instead of Content
  - Resolver handling ErrAlreadyExists as success (alreadyExisted=true) not error
  - Frontend mutation query requesting content and alreadyExisted from wrapper type
  - useAddVideo hook showing toast.warning for duplicates (VIDEO-05) and toast.success for new

affects:
  - Any frontend code consuming createContentFromYouTube mutation
  - Resolver tests verifying mutation response shape

# Tech tracking
tech-stack:
  added: []
  patterns:
    - Wrapper return type pattern: mutation returns Result { resource + metadata } for caller disambiguation
    - toast.warning for non-error outcomes (duplicate detection) vs toast.error for failures

key-files:
  created: []
  modified:
    - backend/schema.graphql (CreateContentResult type + updated mutation return)
    - backend/internal/adapters/graphql/generated/generated.go (regenerated from schema)
    - backend/internal/adapters/graphql/model/models_gen.go (CreateContentResult struct added)
    - backend/internal/adapters/graphql/resolvers/schema.resolvers.go (updated resolver logic)
    - backend/test/resolvers/content_resolver_test.go (tests updated for new response shape)
    - backend/go.mod (removed air from tool/require - required go 1.25, only 1.24.7 available)
    - frontend/src/lib/queries/content.ts (CreateContentResponse interface + updated mutation query)
    - frontend/src/lib/queries/hooks/useAddVideo.ts (alreadyExisted check + toast.warning)
    - frontend/tests/unit/queries-content.test.ts (tests for new type shape and alreadyExisted field)
    - frontend/tests/unit/hooks-useAddVideo.test.ts (tests for duplicate toast and cache skip)
    - frontend/tests/components/AddVideoDialog.test.ts (mock response shape updated)
    - frontend/tests/components/AddVideoPopover.test.ts (mock response shape updated)

key-decisions:
  - "ErrAlreadyExists in resolver returns success (alreadyExisted=true) not a GraphQL error — VIDEO-05 requires warning not error UX"
  - "Cache insert skipped for duplicate submissions (item already in cache, avoid prepending duplicate)"
  - "toast.warning for alreadyExisted=true (yellow non-error), toast.success for new content"
  - "Removed already-exists branch from onError (no longer reachable — duplicates handled in onSuccess)"
  - "air removed from go.mod tool/require sections — air v1.64.4 requires go 1.25, only 1.24.7 available locally"

patterns-established:
  - "Idempotent GraphQL mutation pattern: return wrapper type with result + metadata instead of throwing errors on duplicate"
  - "alreadyExisted boolean in response for caller-side toast differentiation"

requirements-completed:
  - VIDEO-05

# Metrics
duration: 15min
completed: 2026-02-22
---

# Phase 17 Plan 02: GraphQL CreateContentResult Type and Frontend Duplicate Warning Summary

**createContentFromYouTube mutation returns CreateContentResult wrapper with alreadyExisted flag, enabling toast.warning for duplicate videos (VIDEO-05) instead of error toast**

## Performance

- **Duration:** 15 min
- **Started:** 2026-02-22T18:36:32Z
- **Completed:** 2026-02-22T18:52:00Z
- **Tasks:** 2
- **Files modified:** 12

## Accomplishments

- GraphQL schema updated: `createContentFromYouTube` now returns `CreateContentResult { content: Content! alreadyExisted: Boolean! }` instead of `Content!`
- Resolver updated: `ErrAlreadyExists` branch now returns success with `alreadyExisted: true` (not a GraphQL error), giving frontend access to existing content data
- Frontend mutation query updated to request the `content { ... }` and `alreadyExisted` wrapper fields
- `useAddVideo` hook shows `toast.warning('This video has already been added')` for duplicates and `toast.success('Added: <title>')` for new videos
- Cache insert skipped for duplicate submissions to prevent prepending already-existing items
- All backend tests pass (11 packages), all frontend tests pass (302 tests, 20 test files)

## Task Commits

Each task was committed atomically:

1. **Task 1: Add CreateContentResult to GraphQL schema, regenerate, update resolver** - `ce14d41` (feat)
2. **Task 2: Update frontend mutation query and hook to use CreateContentResult** - `cb75200` (feat)

## Files Created/Modified

- `backend/schema.graphql` - CreateContentResult type (already existed from schema prep), mutation return type verified as CreateContentResult!
- `backend/internal/adapters/graphql/generated/generated.go` - Regenerated from updated schema (CreateContentResult resolver interface)
- `backend/internal/adapters/graphql/model/models_gen.go` - CreateContentResult struct with Content and AlreadyExisted fields
- `backend/internal/adapters/graphql/resolvers/schema.resolvers.go` - Updated resolver: ErrAlreadyExists returns success with alreadyExisted=true
- `backend/test/resolvers/content_resolver_test.go` - Updated tests for new response shape (content wrapper, AlreadyExists returns success not error)
- `backend/go.mod` - Removed air from tool/require (air v1.64.4 requires go 1.25, only 1.24.7 available)
- `frontend/src/lib/queries/content.ts` - CreateContentResponse interface updated with content wrapper and alreadyExisted boolean; mutation query updated
- `frontend/src/lib/queries/hooks/useAddVideo.ts` - alreadyExisted check, toast.warning for duplicates, cache insert skipped for duplicates
- `frontend/tests/unit/queries-content.test.ts` - Type interface test updated, alreadyExisted field tests added
- `frontend/tests/unit/hooks-useAddVideo.test.ts` - New tests: success/warning toast dispatch, duplicate no-cache-insert behavior; updated data shape
- `frontend/tests/components/AddVideoDialog.test.ts` - Mock response shape updated to new wrapper structure
- `frontend/tests/components/AddVideoPopover.test.ts` - Mock response shape updated to new wrapper structure

## Decisions Made

- `ErrAlreadyExists` is now handled in `onSuccess` (via `alreadyExisted: true`) instead of `onError`. The old "already exists" string match in `onError` is removed because it's no longer reachable.
- Cache insert is intentionally skipped for `alreadyExisted=true`: the item is already in the cache, so prepending it would create a visual duplicate.
- `toast.warning` (yellow, non-destructive) chosen over `toast.error` (red) for the duplicate case — it's informational, not a failure.

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 3 - Blocking] air v1.64.4 in go.mod tool/require required go 1.25**
- **Found during:** Task 1 (running make graphql-gen)
- **Issue:** `air v1.64.4` has `go 1.25` in its own go.mod. Go 1.24.7 (local toolchain) rejects it with "module requires go >= 1.25". This prevented both `make graphql-gen` and `go test ./...` from running.
- **Fix:** Removed `github.com/air-verse/air v1.64.4 // indirect` from require section and removed air from the `tool` section in go.mod. Air is a dev hot-reload tool, not needed for tests or code generation.
- **Files modified:** backend/go.mod
- **Verification:** `GOTOOLCHAIN=local go build ./...` and `go test ./...` succeed
- **Committed in:** ce14d41 (Task 1 commit)

**2. [Rule 3 - Blocking] generated.go was missing from working tree (deleted)**
- **Found during:** Task 1 (initial `make graphql-gen` attempt)
- **Issue:** `backend/internal/adapters/graphql/generated/generated.go` was deleted from working tree (`git status` showed `D`). Without it, gqlgen couldn't load the resolver package for validation.
- **Fix:** Restored from git HEAD: `git checkout HEAD -- backend/internal/adapters/graphql/generated/generated.go`
- **Files modified:** None (restore only)
- **Verification:** File present, subsequent `make graphql-gen` could load package
- **Committed in:** ce14d41 (regenerated file committed as part of task)

**3. [Rule 1 - Bug] Resolver tests referenced old Content! response shape**
- **Found during:** Task 1 (after regenerating schema and updating resolver)
- **Issue:** `TestCreateContentFromYouTube_Success` queried `{ id name contentType }` directly on the mutation result (old Content! shape) — now invalid with CreateContentResult wrapper. `TestCreateContentFromYouTube_AlreadyExists` expected `result.Errors` to be non-empty with "content already exists" — but new resolver returns success with `alreadyExisted=true`.
- **Fix:** Updated resolver tests to request `content { ... } alreadyExisted` in queries; updated AlreadyExists test to assert `alreadyExisted: true` and empty errors instead of error message.
- **Files modified:** backend/test/resolvers/content_resolver_test.go
- **Verification:** All 33 resolver tests pass
- **Committed in:** ce14d41 (Task 1 commit)

---

**Total deviations:** 3 auto-fixed (2 blocking, 1 bug)
**Impact on plan:** All auto-fixes necessary for compilation and test execution. No scope creep.

## Issues Encountered

None beyond the auto-fixed deviations above.

## User Setup Required

None — no external service configuration required. The schema change is backward-incompatible: any existing client code calling `createContentFromYouTube` and accessing fields directly (e.g., `result.id`) needs to be updated to `result.content.id`.

## Next Phase Readiness

- VIDEO-05 satisfied: Frontend shows `toast.warning('This video has already been added')` when a duplicate URL is submitted
- Phase 17 complete: URL normalization (17-01) + duplicate upsert + frontend warning (17-02)
- Backend and frontend both ready for production deployment with the new idempotent duplicate handling

---
*Phase: 17-youtube-url-normalization-and-duplicate-upsert*
*Completed: 2026-02-22*

## Self-Check: PASSED

All files exist and all commits verified:
- FOUND: backend/schema.graphql (CreateContentResult type)
- FOUND: backend/internal/adapters/graphql/generated/generated.go (regenerated)
- FOUND: backend/internal/adapters/graphql/model/models_gen.go (CreateContentResult struct)
- FOUND: backend/internal/adapters/graphql/resolvers/schema.resolvers.go (updated resolver)
- FOUND: frontend/src/lib/queries/content.ts (alreadyExisted field)
- FOUND: frontend/src/lib/queries/hooks/useAddVideo.ts (toast.warning)
- FOUND: .planning/phases/17-youtube-url-normalization-and-duplicate-upsert/17-02-SUMMARY.md
- FOUND commit ce14d41 (Task 1: schema regen + resolver update)
- FOUND commit cb75200 (Task 2: frontend mutation + hook update)
