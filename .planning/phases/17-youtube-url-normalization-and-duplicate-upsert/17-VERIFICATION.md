---
phase: 17-youtube-url-normalization-and-duplicate-upsert
verified: 2026-02-22T19:05:00Z
status: passed
score: 13/13 must-haves verified
re_verification: false
---

# Phase 17: YouTube URL Normalization and Duplicate Upsert — Verification Report

**Phase Goal:** Normalize all YouTube URL variants to canonical form (`https://www.youtube.com/watch?v=<videoID>`), implement idempotent upsert (return existing content instead of error on duplicate), and expose `alreadyExisted` signal to frontend for VIDEO-05 duplicate warning toast
**Verified:** 2026-02-22T19:05:00Z
**Status:** PASSED
**Re-verification:** No — initial verification

---

## Goal Achievement

### Observable Truths

Success criteria from ROADMAP.md Phase 17:

| # | Truth | Status | Evidence |
|---|-------|--------|----------|
| 1 | All YouTube URL variants normalize to `https://www.youtube.com/watch?v=<videoID>` before storage | VERIFIED | `youtube.NormalizeYouTubeURL(videoID)` called at line 39 in `content_service.go`; `ExtractVideoID` runs first to pull the ID from any URL variant, then canonical form is constructed |
| 2 | `CreateFromYouTube` returns existing content on duplicate URL (not an error) — M-03 resolved | VERIFIED | Service returns `(existing, domain.ErrAlreadyExists)` — content is non-nil on duplicate (lines 43-44 and 75-76); resolver converts this to GraphQL success with `alreadyExisted: true` |
| 3 | Repository uses `INSERT ON CONFLICT DO NOTHING` for atomic deduplication (no TOCTOU race) | VERIFIED | `GetOrCreateByURL` at line 75 uses `clause.OnConflict{Columns: []clause.Column{{Name: "url"}}, DoNothing: true}` when `refreshOnConflict=false`; called with `refreshOnConflict=true` from service for metadata refresh path, which uses `DoUpdates` |
| 4 | `createContentFromYouTube` GraphQL mutation returns `CreateContentResult { content, alreadyExisted }` | VERIFIED | `schema.graphql` line 73-76 defines `type CreateContentResult { content: Content! alreadyExisted: Boolean! }`; mutation at line 204 returns `CreateContentResult!`; `models_gen.go` line 52 has generated `CreateContentResult` struct |
| 5 | Frontend shows warning toast when video already exists (VIDEO-05) | VERIFIED | `useAddVideo.ts` line 27: `toast.warning('This video has already been added')` triggered when `result?.alreadyExisted` is true |
| 6 | Frontend shows success toast when new video is added | VERIFIED | `useAddVideo.ts` line 29: `toast.success(\`Added: ${newItem?.name ?? 'video'}\`)` triggered in the else branch |
| 7 | Data migration backfills existing raw URLs to canonical form | VERIFIED | `backend/migrations/000012_backfill_canonical_urls.up.sql` exists and contains `UPDATE content SET url = 'https://www.youtube.com/watch?v=' || ...` with regex extraction for all URL variants |
| 8 | All backend and frontend tests pass | VERIFIED | Backend: 11 packages pass (`go test ./...`); Frontend: 302 tests across 20 test files pass (`pnpm run test:run`) |

**Score:** 8/8 success criteria verified (13/13 artifact+link checks below also pass)

---

### Required Artifacts

#### Plan 17-01 Artifacts

| Artifact | Expected | Status | Details |
|----------|----------|--------|---------|
| `backend/internal/adapters/youtube/parser.go` | `NormalizeYouTubeURL` pure function | VERIFIED | Function exists at line 13; returns `"https://www.youtube.com/watch?v=" + videoID`; substantive (not a stub) |
| `backend/internal/core/ports/repositories/content_repository.go` | `GetOrCreateByURL` on `ContentRepository` interface | VERIFIED | Method at line 18: `GetOrCreateByURL(ctx context.Context, content *domain.Content, refreshOnConflict bool) (*domain.Content, bool, error)` |
| `backend/internal/adapters/repositories/postgres/gorm_content_repository.go` | GORM implementation with `clause.OnConflict` | VERIFIED | `GetOrCreateByURL` at line 72; `clause.OnConflict` at line 75; handles both `DoNothing` and `DoUpdates` branches; re-fetches on conflict |
| `backend/internal/core/services/content_service.go` | Idempotent `CreateFromYouTube` returning `ErrAlreadyExists` on duplicate | VERIFIED | Calls `youtube.NormalizeYouTubeURL` at line 39; calls `GetOrCreateByURL` at line 70; returns `(created, domain.ErrAlreadyExists)` when `alreadyExisted=true` |
| `backend/migrations/000012_backfill_canonical_urls.up.sql` | Data migration with `UPDATE content` | VERIFIED | Contains `UPDATE content SET url = 'https://www.youtube.com/watch?v=' || ...` with CASE for youtu.be, watch?v=, embed/shorts/live paths |

#### Plan 17-02 Artifacts

| Artifact | Expected | Status | Details |
|----------|----------|--------|---------|
| `backend/schema.graphql` | `CreateContentResult` type with `content` and `alreadyExisted` | VERIFIED | Lines 73-76: `type CreateContentResult { content: Content! alreadyExisted: Boolean! }`; mutation at line 204 returns `CreateContentResult!` |
| `backend/internal/adapters/graphql/resolvers/schema.resolvers.go` | Resolver returning `CreateContentResult` | VERIFIED | `CreateContentFromYouTube` at line 22 returns `*model.CreateContentResult`; `ErrAlreadyExists` handled at lines 26-31 returning `AlreadyExisted: true` with content |
| `frontend/src/lib/queries/content.ts` | `CreateContentResponse` with `alreadyExisted`, mutation requests wrapper fields | VERIFIED | `CreateContentResponse` interface at lines 33-38 has `{ content: ContentItem; alreadyExisted: boolean }`; mutation query at lines 103-125 requests `content { ... }` and `alreadyExisted` |
| `frontend/src/lib/queries/hooks/useAddVideo.ts` | Hook checks `alreadyExisted`, shows `toast.warning` for duplicates | VERIFIED | `result?.alreadyExisted` checked at line 25; `toast.warning('This video has already been added')` at line 27; cache insert skipped for duplicates at line 34 |

---

### Key Link Verification

#### Plan 17-01 Key Links

| From | To | Via | Status | Details |
|------|----|-----|--------|---------|
| `content_service.go` | `parser.go` | `youtube.NormalizeYouTubeURL(videoID)` call | WIRED | Line 39: `canonicalURL := youtube.NormalizeYouTubeURL(videoID)` — import present at line 8, call present, result stored and used |
| `content_service.go` | `content_repository.go` | `s.repo.GetOrCreateByURL(ctx, content, true)` call | WIRED | Line 70: `created, alreadyExisted, err := s.repo.GetOrCreateByURL(ctx, content, true)` — return values all consumed |
| `gorm_content_repository.go` | `gorm.io/gorm/clause` | `clause.OnConflict{DoNothing: true}` for atomic upsert | WIRED | Import at line 13; `clause.OnConflict` struct constructed at line 75; `DoNothing: true` branch at line 81 |

#### Plan 17-02 Key Links

| From | To | Via | Status | Details |
|------|----|-----|--------|---------|
| `schema.resolvers.go` | `models_gen.go` | `model.CreateContentResult` struct usage | WIRED | `*model.CreateContentResult` as return type at line 22; `&model.CreateContentResult{...}` constructed at lines 27-31 and 45-48 |
| `useAddVideo.ts` | `content.ts` | `CreateContentResponse` type import | WIRED | Line 4: `import { CREATE_CONTENT_FROM_YOUTUBE, type CreateContentResponse, type ContentResponse } from '../content'` |
| `useAddVideo.ts` | `content.ts` | `alreadyExisted` field check in `onSuccess` | WIRED | Line 25: `if (result?.alreadyExisted)` — field from `CreateContentResponse.createContentFromYouTube.alreadyExisted` |

---

### Requirements Coverage

| Requirement | Source Plan | Description | Status | Evidence |
|-------------|------------|-------------|--------|----------|
| VIDEO-05 | 17-02-PLAN.md | User is warned via toast if video already exists in the system (duplicate detection) | SATISFIED | `useAddVideo.ts` shows `toast.warning('This video has already been added')` when `alreadyExisted=true`; resolver returns `alreadyExisted: true` instead of error; frontend test `hooks-useAddVideo.test.ts` covers this path |
| M-03 | 17-01-PLAN.md | `CreateFromYouTube` returns error instead of idempotent result (bug backlog item from ROADMAP.md) | SATISFIED | `CreateFromYouTube` now returns `(existingContent, ErrAlreadyExists)` — content is non-nil on duplicate; resolver converts to GraphQL success; 7 `TestCreateFromYouTube_*` service tests all pass including `TestCreateFromYouTube_ReturnExistingOnDuplicate` |

Note: M-03 is tracked in ROADMAP.md Bug Backlog (not in REQUIREMENTS.md), defined as "`CreateFromYouTube` returns error instead of idempotent result". ROADMAP.md Phase 17 line 597 explicitly states M-03 is resolved. The bug backlog entry at line 306 remains marked `[ ]` (open) — this is expected; ROADMAP backlog items are marked closed in phase Success Criteria, not in the backlog list.

---

### Anti-Patterns Found

No anti-patterns detected in key files:
- No `TODO/FIXME/HACK/PLACEHOLDER` comments in phase files
- No stub implementations (`return null`, empty bodies, console.log-only handlers)
- No unwired state variables
- No empty API handlers

One notable deviation from the 17-01 PLAN spec: the PLAN specified `GetOrCreateByURL(ctx, content)` (2 params), but the actual interface has `GetOrCreateByURL(ctx, content, refreshOnConflict bool)` (3 params). This is an intentional enhancement — the implementation supports both `DoNothing` (false) and `DoUpdates` (true) conflict strategies. The service calls with `refreshOnConflict=true` to update metadata on re-submit. This is more capable than the spec, not a regression.

---

### Human Verification Required

The following items require manual testing in a running environment:

#### 1. Duplicate URL Warning Toast (VIDEO-05 End-to-End)

**Test:** Add a YouTube video. Submit the same URL a second time.
**Expected:** Second submission shows a yellow warning toast "This video has already been added" (not an error toast, not a success toast). The video list does not show a duplicate item.
**Why human:** Toast appearance (color, position, auto-dismiss timing) and duplicate-free list state cannot be verified programmatically.

#### 2. URL Variant Deduplication

**Test:** Add a video via `https://www.youtube.com/watch?v=dQw4w9WgXcQ`. Then add the same video via `https://youtu.be/dQw4w9WgXcQ`.
**Expected:** Second submission shows "This video has already been added" warning. The video appears only once in the list.
**Why human:** Requires a running backend with a real database to verify the canonical URL lookup path works end-to-end.

#### 3. Data Migration Correctness

**Test:** Run `make migrate-up` against a database containing legacy YouTube URLs (e.g., `https://youtu.be/dQw4w9WgXcQ`). Verify the URL is updated to `https://www.youtube.com/watch?v=dQw4w9WgXcQ`.
**Expected:** All non-canonical YouTube URLs in the `content` table are normalized.
**Why human:** Requires a live PostgreSQL instance with pre-existing data.

---

### Gaps Summary

No gaps found. All 8 success criteria from ROADMAP.md Phase 17 are verified. Both requirements (VIDEO-05 and M-03) have clear implementation evidence and passing tests.

---

_Verified: 2026-02-22T19:05:00Z_
_Verifier: Claude (gsd-verifier)_
