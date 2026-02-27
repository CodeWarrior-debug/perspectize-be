---
phase: 04-add-perspective-flow
plan: 02
subsystem: frontend-ui-agrid
tags: [perspective, ag-grid, svelte5, rating-input, modal, cell-renderer]
dependency_graph:
  requires:
    - frontend/src/lib/queries/hooks/useCreatePerspective.ts
    - frontend/src/lib/queries/hooks/useUpdatePerspective.ts
    - frontend/src/lib/queries/perspectives.ts
  provides:
    - frontend/src/lib/utils/ratings.ts
    - frontend/src/lib/components/RatingInput.svelte
    - frontend/src/lib/components/PerspectivePopover.svelte
  affects:
    - frontend/src/lib/components/ActivityTable.svelte
    - frontend/src/lib/utils/formatting.ts
tech_stack:
  added: []
  patterns:
    - "AG Grid context reactivity: setGridOption('context', ...) + refreshCells() in $effect for cell renderer data access"
    - "Centered modal Dialog (not anchored popover) for perspective form — per Figma Make decision"
    - "Hold-to-repeat stepper: 300ms initial delay then 75ms interval for rating increment/decrement"
    - "hasInteracted pattern: rating shows gray 5.000 default until touched, submits null if never interacted"
    - "perspectivesByContentId Map pattern: O(1) lookup from contentID to PerspectiveItem"
    - "Thumbs toggle: THUMBS_UP / THUMBS_DOWN enum, neither selected by default, click to toggle/deselect"
key_files:
  created:
    - frontend/src/lib/utils/ratings.ts
    - frontend/src/lib/components/RatingInput.svelte
    - frontend/src/lib/components/PerspectivePopover.svelte
    - frontend/tests/unit/utils-ratings.test.ts
  modified:
    - frontend/src/lib/components/ActivityTable.svelte
    - frontend/src/lib/utils/formatting.ts
decisions:
  - "Rating storage uses display * 1000 mapping (0.000-10.000 display, 0-10000 storage, 3 decimal places) — supersedes earlier 2-decimal spec per CONTEXT.md resolved blocker A"
  - "Perspectize column trigger is click-based for both create and edit — not hover — per CONTEXT.md resolved blocker B"
  - "PerspectivePopover uses shadcn Dialog (centered modal) for all viewport sizes, not Popover anchored to cell — per Figma Make spec"
  - "Phase 4 form contains no text inputs (no Review textarea, no Title input) — only 4 ratings + thumbs toggle + Add More expansion"
  - "parseLike helper function used instead of $state initialization with arrow function — resolves Svelte 5 type error"
  - "Rating state initialized to null (not from existingPerspective) to avoid Svelte 5 state_referenced_locally warnings; $effect syncs values on mount"
metrics:
  duration: 6 min
  completed_date: "2026-02-27"
  tasks: 2
  files: 6
---

# Phase 4 Plan 02: Perspective Form UI + AG Grid Perspectize Column Summary

Core user-facing perspective flow: RatingInput stepper component with hold-to-repeat and progress bar, PerspectivePopover Dialog with 2x2 rating grid and thumbs toggle, Perspectize column in AG Grid with "+"/glasses icons, click-based create/edit modal, and full perspectives-by-user query wiring.

## Tasks Completed

| Task | Name | Commit | Key Files |
|------|------|--------|-----------|
| 1 | RatingInput component, rating utilities, PerspectivePopover | `b4ea876` | ratings.ts, RatingInput.svelte, PerspectivePopover.svelte, utils-ratings.test.ts |
| 2 | AG Grid Perspectize column with cell renderers and data wiring | `d844544` | ActivityTable.svelte, formatting.ts |

## Verification Results

1. `cd frontend && pnpm run test:run` — 385 tests pass (25 test files, +24 new rating utility tests)
2. `cd frontend && pnpm run check` — 3 pre-existing errors only (same as 04-01-SUMMARY: +page.svelte:27, FormPopover.test.ts:37-38); 0 new errors introduced
3. `perspectiveCellRenderer` and `perspectiveHeaderRenderer` added to formatting.ts
4. Perspectize column is first column (50px, always visible) with glasses/plus icons
5. `LIST_PERSPECTIVES_BY_USER` query in ActivityTable with `perspectivesByContentId` Map for O(1) lookup
6. `onCellClicked` opens PerspectivePopover (create or edit mode based on existing perspective)
7. AG Grid context updated reactively with `$effect` → refreshCells on perspectivesByContentId change

## Visual Verification

Auto-approved per user request for fully autonomous execution.

Key expected behaviors (verifiable via browser):
- Perspectize column (leftmost, 50px) shows "+" on all rows when no perspectives exist
- Click "+" → Dialog opens with "Add Perspective" title, 4 rating inputs, thumbs up/down
- Rating inputs show gray 5.000 default; interaction turns to primary blue color
- Submit with all empty → toast error "Please fill in at least one field"
- Submit with rating → toast success "Perspective added", icon changes to glasses
- Click glasses → Dialog opens with "Edit Perspective" title, pre-populated fields
- Submit edit → toast success "Perspective updated"

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 1 - Bug] Svelte 5 type error: $state arrow function initialization for LikeValue**
- **Found during:** Task 1 TypeScript check
- **Issue:** `$state<LikeValue>(() => {...})` — TypeScript rejects this because the arrow function type doesn't match `LikeValue`
- **Fix:** Extracted `parseLike()` helper function, called before `$state()`: `$state<LikeValue>(parseLike(existingPerspective?.like))`; then simplified to `$state<LikeValue>(null)` with `$effect` handling initialization
- **Files modified:** `frontend/src/lib/components/PerspectivePopover.svelte`
- **Commit:** `b4ea876`

**2. [Rule 1 - Bug] Svelte 5 state_referenced_locally warnings for rating initialization**
- **Found during:** Task 1 TypeScript check
- **Issue:** Initializing `$state` from `existingPerspective?.quality` captured the initial prop value — Svelte warned this doesn't react to prop changes
- **Fix:** Initialize all rating states to `null`, rely on `$effect` (already present) to sync from `existingPerspective` on mount and change
- **Files modified:** `frontend/src/lib/components/PerspectivePopover.svelte`
- **Commit:** `b4ea876`

### Design Deviations from Plan (Intentional, Per CONTEXT.md)

**Plan vs CONTEXT.md conflicts resolved in favor of CONTEXT.md (locked decisions):**

1. **Modal behavior:** Plan specified "Popover on desktop / Dialog on mobile" — CONTEXT.md says "centered modal Dialog always" (Figma Make decision). Implemented as Dialog always.
2. **Edit trigger:** Plan specified `onCellMouseOver` for silhouette hover — CONTEXT.md resolved blocker B says click-based only. Implemented as `onCellClicked` for both create and edit.
3. **Like field:** Plan specified text input — CONTEXT.md decision D says thumbs up/down toggle. Implemented as THUMBS_UP/THUMBS_DOWN toggle buttons.
4. **No text inputs:** Plan mentioned Review textarea — CONTEXT.md decision F says no text inputs in Phase 4. Omitted Review textarea and Title input.
5. **Rating display:** Plan said 0.00 (2 decimals) — CONTEXT.md resolved blocker A says 0.000 (3 decimals, display * 1000 mapping). Implemented correctly.

## Must-Haves Verification

- [x] User can click '+' icon in the Perspectize column on a content row to open the perspective creation form
- [x] User can set Quality, Agreement, Importance, and Confidence ratings via number inputs with progress bar visualization (0.000-10.000 display, 0-10000 storage)
- [~] User can enter Like text — RESOLVED: Like is thumbs toggle (not text input), per CONTEXT.md decision D
- [x] User sees validation error toast if form is submitted with all fields empty
- [x] User sees success toast after perspective is created, attributed to selected user
- [x] User sees silhouette-with-glasses icon on rows where they already have a perspective
- [x] User can click silhouette-with-glasses icon to open pre-populated form for editing (click-based per CONTEXT.md, not hover)
- [x] Form renders as a centered Dialog (Figma Make decision — all viewports)

## Self-Check: PASSED

Files created verified:
- `frontend/src/lib/utils/ratings.ts` — EXISTS
- `frontend/src/lib/components/RatingInput.svelte` — EXISTS
- `frontend/src/lib/components/PerspectivePopover.svelte` — EXISTS
- `frontend/tests/unit/utils-ratings.test.ts` — EXISTS (24 tests)

Files modified verified:
- `frontend/src/lib/components/ActivityTable.svelte` — EXISTS (perspectize column + query added)
- `frontend/src/lib/utils/formatting.ts` — EXISTS (perspectiveCellRenderer + perspectiveHeaderRenderer added)

Commits verified:
- `b4ea876` feat(04-02): add RatingInput, PerspectivePopover, and rating utility tests — EXISTS
- `d844544` feat(04-02): add Perspectize column to AG Grid with perspective modal integration — EXISTS
