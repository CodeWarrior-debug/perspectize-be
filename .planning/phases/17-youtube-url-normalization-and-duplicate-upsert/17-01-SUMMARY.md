---
phase: 17-youtube-url-normalization-and-duplicate-upsert
plan: 01
subsystem: api
tags: [go, gorm, postgresql, youtube, deduplication, upsert, url-normalization]

# Dependency graph
requires:
  - phase: 07.1-orm-migration-sqlx-to-gorm
    provides: GORM hex-clean pattern with separate domain/model types and mappers
  - phase: 03-add-video-flow
    provides: CreateFromYouTube service and YouTubeClient port

provides:
  - NormalizeYouTubeURL pure function in youtube adapter package
  - GetOrCreateByURL atomic upsert method on ContentRepository interface and GORM implementation
  - Idempotent CreateFromYouTube returning (existingContent, ErrAlreadyExists) on duplicate
  - Data migration 000012 backfilling existing raw URLs to canonical form

affects:
  - Any phase using ContentRepository interface (must implement GetOrCreateByURL in mocks)
  - Resolver tests for createContentFromYouTube mutation
  - Frontend duplicate handling (resolver still returns error on duplicate to GraphQL clients)

# Tech tracking
tech-stack:
  added:
    - gorm.io/gorm/clause (already in go.mod, first use of clause.OnConflict)
  patterns:
    - Atomic upsert via INSERT ON CONFLICT (url) DO NOTHING + RowsAffected check
    - Service-level canonical URL normalization before all DB operations
    - Returning (content, ErrAlreadyExists) for caller disambiguation (content available even on duplicate)

key-files:
  created:
    - backend/internal/adapters/youtube/parser.go (NormalizeYouTubeURL added)
    - backend/migrations/000012_backfill_canonical_urls.up.sql
    - backend/migrations/000012_backfill_canonical_urls.down.sql
  modified:
    - backend/internal/core/ports/repositories/content_repository.go (GetOrCreateByURL added)
    - backend/internal/adapters/repositories/postgres/gorm_content_repository.go (GetOrCreateByURL impl)
    - backend/internal/core/services/content_service.go (CreateFromYouTube refactored)
    - backend/test/youtube/parser_test.go (TestNormalizeYouTubeURL tests)
    - backend/test/services/content_service_test.go (tests rewritten for new behavior)
    - backend/test/services/user_service_test.go (mock updated for new interface)
    - backend/test/resolvers/content_resolver_test.go (mock updated + Success test fixed)

key-decisions:
  - "NormalizeYouTubeURL is a pure function (not a YouTubeClient method) to avoid polluting the port interface"
  - "CreateFromYouTube validates URL (ExtractVideoID) BEFORE GetByURL lookup — normalizes first, checks DB second"
  - "Service returns (content, ErrAlreadyExists) on duplicate — callers get both content AND can distinguish new vs existing"
  - "Resolver still returns nil + error on ErrAlreadyExists — API clients get error, no behavior change externally"
  - "GetOrCreateByURL uses clause.OnConflict{DoNothing: true} on url column for TOCTOU-safe atomic upsert"
  - "go.mod downgraded from go 1.25.0 to go 1.24.0 to allow tests with locally available go 1.24.7 toolchain"

patterns-established:
  - "Canonical URL normalization: all YouTube URL variants resolve to https://www.youtube.com/watch?v=<ID>"
  - "Atomic deduplication: INSERT ON CONFLICT DO NOTHING + RowsAffected==0 fallback fetch"
  - "Idempotent service methods: return existing resource + sentinel error for caller disambiguation"

requirements-completed:
  - M-03

# Metrics
duration: 8min
completed: 2026-02-22
---

# Phase 17 Plan 01: YouTube URL Normalization and Duplicate Upsert Summary

**YouTube URL normalization with NormalizeYouTubeURL pure function and atomic GetOrCreateByURL upsert eliminating TOCTOU race condition in duplicate detection**

## Performance

- **Duration:** 8 min
- **Started:** 2026-02-22T18:19:52Z
- **Completed:** 2026-02-22T18:27:43Z
- **Tasks:** 2
- **Files modified:** 11

## Accomplishments

- NormalizeYouTubeURL pure function maps any video ID to canonical `https://www.youtube.com/watch?v=<ID>` form
- CreateFromYouTube now normalizes input URL before all DB operations, eliminating variant-URL duplicates
- GetOrCreateByURL uses `clause.OnConflict{DoNothing: true}` on the `url` column for atomic race-condition-free upsert
- Service returns `(existingContent, ErrAlreadyExists)` on duplicate instead of `(nil, error)` — callers get the content
- Data migration 000012 backfills existing non-canonical YouTube URLs to canonical form
- All 7 CreateFromYouTube service tests rewritten and passing; all backend tests green

## Task Commits

Each task was committed atomically:

1. **Task 1 RED: Add NormalizeYouTubeURL tests** - `f908bb6` (test)
2. **Task 1 GREEN: Implement NormalizeYouTubeURL** - `18d6b6c` (feat)
3. **Task 2: Repository, service, migration, updated tests** - `2505eae` (feat)

_Note: TDD tasks have separate RED (test) and GREEN (implementation) commits_

## Files Created/Modified

- `backend/internal/adapters/youtube/parser.go` - Added NormalizeYouTubeURL pure function
- `backend/internal/core/ports/repositories/content_repository.go` - Added GetOrCreateByURL to interface
- `backend/internal/adapters/repositories/postgres/gorm_content_repository.go` - Implemented GetOrCreateByURL with clause.OnConflict
- `backend/internal/core/services/content_service.go` - Refactored CreateFromYouTube: normalize first, atomic upsert, return content on duplicate
- `backend/test/youtube/parser_test.go` - TestNormalizeYouTubeURL test table (3 cases)
- `backend/test/services/content_service_test.go` - Rewrote 6 tests, added TestCreateFromYouTube_ReturnExistingOnDuplicate and TestCreateFromYouTube_NormalizesURLVariants
- `backend/test/services/user_service_test.go` - Added GetOrCreateByURL to mockContentRepoForUser
- `backend/test/resolvers/content_resolver_test.go` - Added GetOrCreateByURL to mock, updated Success test
- `backend/migrations/000012_backfill_canonical_urls.up.sql` - Data migration for existing raw URLs
- `backend/migrations/000012_backfill_canonical_urls.down.sql` - Irreversible stub (data normalization)
- `backend/go.mod` - Downgraded go directive from 1.25.0 to 1.24.0 (local toolchain compatibility)

## Decisions Made

- NormalizeYouTubeURL is a pure function in the youtube adapter package (not a YouTubeClient port method) to keep the interface clean and avoid updating all mock implementations
- ExtractVideoID validation happens BEFORE GetByURL lookup — fail fast on invalid URLs, normalize before any DB interaction
- Service returns `(content, ErrAlreadyExists)` not `(nil, error)` — the caller (resolver) can access the existing content even on duplicate, while still being able to detect it was pre-existing
- Resolver behavior unchanged externally: still returns GraphQL error on duplicate (clients need no changes)
- go.mod downgraded from go 1.25.0 to go 1.24.0 to run tests with the locally available go 1.24.7 toolchain (no network access for toolchain download)

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 3 - Blocking] go.mod required go 1.25.0 but only 1.24.7 available locally**
- **Found during:** Task 1 (initial test run)
- **Issue:** `go test` attempted to download go 1.25.0 toolchain from Google storage (no network access)
- **Fix:** Downgraded `go 1.25.0` to `go 1.24.0` in backend/go.mod
- **Files modified:** backend/go.mod
- **Verification:** `GOTOOLCHAIN=local go test ./...` succeeds
- **Committed in:** 18d6b6c (Task 1 feat commit)

**2. [Rule 3 - Blocking] Other test mocks (resolver, user_service) didn't implement new GetOrCreateByURL interface method**
- **Found during:** Task 2 (initial `go test ./...` run)
- **Issue:** Adding GetOrCreateByURL to the ContentRepository interface broke compilation of resolver and user_service test files — their mockContentRepository structs didn't implement the new method
- **Fix:** Added GetOrCreateByURL method to `mockContentRepoForUser` in user_service_test.go and `mockContentRepository` in content_resolver_test.go; also updated TestCreateContentFromYouTube_Success to use getOrCreateByURLFn instead of createFn
- **Files modified:** backend/test/services/user_service_test.go, backend/test/resolvers/content_resolver_test.go
- **Verification:** All tests compile and pass
- **Committed in:** 2505eae (Task 2 commit)

---

**Total deviations:** 2 auto-fixed (2 blocking)
**Impact on plan:** Both auto-fixes necessary for compilation and test execution. No scope creep.

## Issues Encountered

None beyond the auto-fixed deviations above.

## User Setup Required

None - no external service configuration required. The data migration (000012) will need to be run against the database via `make migrate-up` when deploying.

## Next Phase Readiness

- M-03 resolved: CreateFromYouTube is now idempotent — returns existing content instead of error on duplicate
- URL normalization ensures all YouTube URL variants (youtu.be, shorts, live, mobile, embed) deduplicate correctly
- TOCTOU race condition eliminated via atomic ON CONFLICT DO NOTHING
- Data migration ready to normalize any pre-existing raw URLs in the database
- Resolver behavior unchanged (still returns GraphQL error on duplicate) — no frontend changes required

---
*Phase: 17-youtube-url-normalization-and-duplicate-upsert*
*Completed: 2026-02-22*

## Self-Check: PASSED

All files exist and all commits verified:
- FOUND: backend/internal/adapters/youtube/parser.go (NormalizeYouTubeURL)
- FOUND: backend/internal/core/ports/repositories/content_repository.go (GetOrCreateByURL)
- FOUND: backend/internal/adapters/repositories/postgres/gorm_content_repository.go (clause.OnConflict)
- FOUND: backend/internal/core/services/content_service.go (NormalizeYouTubeURL usage)
- FOUND: backend/migrations/000012_backfill_canonical_urls.up.sql (UPDATE content)
- FOUND: backend/migrations/000012_backfill_canonical_urls.down.sql
- FOUND: .planning/phases/17-youtube-url-normalization-and-duplicate-upsert/17-01-SUMMARY.md
- FOUND commit f908bb6 (test: NormalizeYouTubeURL tests)
- FOUND commit 18d6b6c (feat: NormalizeYouTubeURL implementation)
- FOUND commit 2505eae (feat: GetOrCreateByURL, refactored service, migration)
