# AG Grid Column Picker — Design

**Date:** 2026-09-02
**Status:** Approved
**Scope:** Frontend only (no backend/schema/migration changes)

## Summary

Add a user-facing control to show/hide columns in the `ActivityTable` AG Grid.
A gear icon in the table's bottom toolbar opens a **modal dialog** with a
checkbox list of columns. Selections are **session-only** (cleared on page
refresh). Admin users get an additional "Internal" group exposing IDs,
timestamps, and the source URL as optional columns.

## Problem

`ActivityTable.svelte` has a fixed column set whose visibility is driven
entirely by a responsive `$effect` keyed to four breakpoint tiers
(xs/sm/md/lg). Users cannot choose which columns they see. Admins have no way
to surface internal fields (content ID, submitter user ID, timestamps, source
URL) for debugging/support.

## Solution

### Visibility model — override layered over the responsive effect

- The responsive `$effect` (progressive reveal by tier) stays **unchanged**
  and remains the default behaviour on page load.
- New session state: `userColumnOverride = $state<Record<string, boolean> | null>(null)`
  — a `colId → visible` map, initially `null`.
- The responsive `$effect` gains one guard at the top:
  `if (userColumnOverride) return;`
- The first time the user toggles a column in the modal, `userColumnOverride`
  becomes a non-null map of every togglable column's current visibility, with
  the toggled column flipped. From that point the responsive `$effect` is
  inert for the rest of the session.
- Each toggle calls `gridApi.setColumnsVisible([colId], next)` immediately for
  live preview, and updates `userColumnOverride`.
- **Session-only:** no persistence. A page refresh resets `userColumnOverride`
  to `null` and automatic responsive layout resumes.
- The modal shows a standing notice whenever `userColumnOverride` is non-null:
  > "Columns are set manually for this session. Refresh the page to restore
  > automatic responsive layout."

### The picker UI

- **Trigger:** a gear icon button (`@lucide/svelte/icons/settings-2` or
  `sliders-horizontal`) in the bottom pagination/toolbar row of
  `ActivityTable.svelte`, left cluster, near the page-size `<select>`.
- **Surface:** a modal `Dialog` built on the existing
  `$lib/components/shadcn/dialog` primitives (same pattern as
  `AddVideoDialog.svelte`). New component:
  `frontend/src/lib/components/ColumnPickerDialog.svelte`.
- **Contents:**
  - Group **"Columns"** — one checkbox per data column:
    item, type, duration, views, likes, publishDate, channel, tags,
    description. Label = the colDef `headerName` (fallback to a friendly name
    for the blank-header cases; `item` → "Item", etc.).
  - Group **"Internal"** — rendered only when `isAdmin`. Checkboxes for:
    Content ID (`id`), Submitter user ID (`addedByUserID`), Source URL
    (`url`), Created at (`createdAt`), Updated at (`updatedAt`). All
    **unchecked by default**; none forced.
  - The Perspectize action column (`colId: 'perspectize'`) is **not listed**
    and always visible.
  - The manual-mode notice (above) when active.
- **Behaviour:** checkbox change → toggle that column live + set override.
  Closing the dialog keeps state. No "apply"/"cancel" — changes are immediate.

### Admin signal

- Add `role` to the `ME` GraphQL query and `Me` interface in
  `frontend/src/lib/queries/users.ts` (`role: UserRole` — backend already
  exposes `me { role }`, enum `ADMIN | SENTINEL | DEFAULT`).
- New hook `frontend/src/lib/queries/hooks/useMe.svelte.ts`: wraps
  `useClerkContext()` + the `['me', clerkUserId]` query (same key/staleTime as
  `AuthUserSync.svelte`, so the cache is shared — no extra request) and
  exposes `{ me, isAdmin }` where `isAdmin = me?.role === 'ADMIN'`.
- `ActivityTable.svelte` consumes `useMe()` and passes `isAdmin` to
  `ColumnPickerDialog`.

### Column definitions

- Add `addedByUserID` to the `LIST_CONTENT` selection set and the
  `ContentItem` interface in `frontend/src/lib/queries/content.ts`
  (`url`, `id`, `createdAt`, `updatedAt` are already selected).
- In `columnDefs` (`ActivityTable.svelte`):
  - `description` already exists as a hidden col — keep, list it in the
    "Columns" group.
  - `createdAt` / `updatedAt` already exist as hidden cols — move them from
    "always hidden" to the admin "Internal" group (still `hide: true` by
    default; the responsive `$effect`'s `alwaysHidden` list drops
    `createdAt`/`updatedAt` and they simply stay hidden unless an admin
    enables them via override).
  - New hidden colDefs: `id` (headerName "Content ID"), `addedByUserID`
    ("Submitter"), `url` ("Source URL"). `hide: true`, `sortable: false`,
    `filter: false`, plain text; `url` uses a value formatter that returns the
    raw string.
- Admin colDefs are **always present** in `columnDefs`. Gating is at the
  modal (only listed when `isAdmin`) plus a guard so a non-admin can never end
  up with them visible (the override map only ever contains togglable colIds,
  and for a non-admin the "Internal" ids are excluded from that set).

### Togglable-column registry

Introduce a single source of truth so the modal, the override-map builder,
and tests agree:

```ts
// in $lib/utils/grid-config.ts
export const DATA_COLUMNS = [
  { colId: 'item', label: 'Item' },
  { colId: 'type', label: 'Type' },
  { colId: 'duration', label: 'Length' },
  { colId: 'views', label: 'Views' },
  { colId: 'likes', label: 'Likes' },
  { colId: 'publishDate', label: 'Published' },
  { colId: 'channel', label: 'Channel' },
  { colId: 'tags', label: 'Tags' },
  { colId: 'description', label: 'Description' },
] as const;

export const INTERNAL_COLUMNS = [
  { colId: 'id', label: 'Content ID' },
  { colId: 'addedByUserID', label: 'Submitter' },
  { colId: 'url', label: 'Source URL' },
  { colId: 'createdAt', label: 'Created at' },
  { colId: 'updatedAt', label: 'Updated at' },
] as const;

export function togglableColIds(isAdmin: boolean): string[] {
  return [
    ...DATA_COLUMNS.map((c) => c.colId),
    ...(isAdmin ? INTERNAL_COLUMNS.map((c) => c.colId) : []),
  ];
}
```

## Data flow

1. `useMe()` resolves `isAdmin` from the cached `me` query.
2. On grid ready, the responsive `$effect` runs as today (override is `null`).
3. User clicks the gear → `ColumnPickerDialog` opens, reads current visibility
   from `gridApi.getColumnState()` to seed checkbox checked-state.
4. User toggles a checkbox → `onToggle(colId, next)` in `ActivityTable`:
   - if `userColumnOverride` is `null`, build it from
     `togglableColIds(isAdmin)` mapped to their current `getColumnState`
     visibility;
   - set `userColumnOverride[colId] = next`;
   - `gridApi.setColumnsVisible([colId], next)`.
5. Responsive `$effect` now early-returns; breakpoint changes no longer alter
   columns.
6. Page refresh → `userColumnOverride` back to `null` → step 2 again.

## Error / edge handling

- **Grid not ready / destroyed:** `onToggle` guards on `gridApi && gridReady`
  (matches existing `$effect` guards). Dialog trigger is disabled until
  `gridReady`.
- **`me` query pending/failed:** `isAdmin` is `false` → no "Internal" group.
  No error surfaced (non-critical).
- **`cardMode` active (<860px):** the grid is unmounted; the gear button is
  hidden in `cardMode` (card list has no columns). Override state persists and
  re-applies if the viewport widens back into grid mode (guard in the
  responsive `$effect` still holds).
- **Column state race:** toggles are synchronous against a live `gridApi`; no
  async apply.
- **All columns hidden:** allowed; AG Grid shows an empty grid body. Not
  prevented (user can re-enable via the still-open dialog).

## Testing

- `frontend/tests/unit/queries-content.test.ts` — assert `LIST_CONTENT`
  includes `addedByUserID`.
- `frontend/tests/unit/queries-users.test.ts` (new or existing) — assert `ME`
  includes `role`.
- `frontend/tests/unit/grid-config.test.ts` — `togglableColIds(false)` = 9
  data colIds; `togglableColIds(true)` = 14; `DATA_COLUMNS` /
  `INTERNAL_COLUMNS` colIds match the real `columnDefs` colIds (guard against
  drift).
- `frontend/tests/components/ColumnPickerDialog.test.ts` (new) — renders the
  data group; hides "Internal" group when `isAdmin=false`, shows it when
  `true`; toggling a checkbox fires `onToggle` with `(colId, boolean)`; the
  manual-mode notice renders when `overrideActive` prop is true.
- `frontend/tests/components/ActivityTable.test.ts` — existing suite must
  still pass (JSDOM: grid does not render; assert the gear button is present
  when not in card mode).
- No backend tests (no backend changes).

## Out of scope / follow-up

- **Saved views** (named, persisted column+filter+sort presets) — file in
  `FEATURE_BACKLOG.md`.
- Column reordering / resizing persistence.
- Per-column drag-and-drop ordering in the picker.
- Enterprise AG Grid tool panel / sidebar.

## Files touched

| File | Change |
|------|--------|
| `frontend/src/lib/queries/content.ts` | `addedByUserID` in query + `ContentItem` |
| `frontend/src/lib/queries/users.ts` | `role` in `ME` + `Me` |
| `frontend/src/lib/queries/hooks/useMe.svelte.ts` | **new** — `{ me, isAdmin }` |
| `frontend/src/lib/utils/grid-config.ts` | `DATA_COLUMNS`, `INTERNAL_COLUMNS`, `togglableColIds` |
| `frontend/src/lib/components/ColumnPickerDialog.svelte` | **new** — modal picker |
| `frontend/src/lib/components/ActivityTable.svelte` | gear button, dialog wiring, `userColumnOverride`, responsive `$effect` guard, new hidden colDefs, `useMe` |
| `frontend/tests/unit/queries-content.test.ts` | assertion |
| `frontend/tests/unit/queries-users.test.ts` | assertion (new if absent) |
| `frontend/tests/unit/grid-config.test.ts` | registry tests |
| `frontend/tests/components/ColumnPickerDialog.test.ts` | **new** |
| `FEATURE_BACKLOG.md` | Saved views follow-up |
