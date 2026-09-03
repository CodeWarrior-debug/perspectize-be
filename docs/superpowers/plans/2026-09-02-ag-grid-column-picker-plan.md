# AG Grid Column Picker Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: superpowers:executing-plans (or subagent-driven-development). Steps use checkbox (`- [ ]`) syntax.

**Goal:** Add a gear-icon → modal column picker to `ActivityTable` so users can show/hide grid columns (session-only). Admins additionally get an "Internal" group (Content ID, Submitter, Source URL, Created/Updated at). The existing responsive `$effect` stays as the default and goes inert once the user makes a manual choice; a page refresh restores automatic layout.

**Architecture:** See `docs/superpowers/specs/2026-09-02-ag-grid-column-picker-design.md`. Frontend only — no backend/schema/migration changes (`me { role }`, `Content.addedByUserID`, `url`, `id`, `createdAt`, `updatedAt` all already exist server-side).

**Tech Stack:** Svelte 5 runes, `@tanstack/svelte-query` v6 (function-wrapper), `svelte-clerk` v1, `ag-grid-svelte5` (AG Grid v32 Community), shadcn-svelte `dialog`, Tailwind v4, Vitest + Testing Library.

## Global Constraints

- Branch: `feature/INI-ag-grid-column-picker` (already created).
- Svelte 5 runes only; no `$:` / `export let`.
- TanStack function-wrapper pattern for `createQuery`.
- No backend edits, no `make graphql-gen`.
- Column visibility stays controlled in the two places `frontend/CLAUDE.md` documents (colDef `hide` + responsive `$effect`) — this plan adds a third, lower-priority layer (`userColumnOverride`) that short-circuits the `$effect`, not the colDefs.
- Run `pnpm --dir frontend exec prettier --write` on touched files before committing; `pnpm --dir frontend run check` and `pnpm --dir frontend run test:run` must pass.
- One logical change per commit, conventional-commit messages.

---

## File Structure

**Create:**
- `frontend/src/lib/queries/hooks/useMe.svelte.ts` — `{ me, isAdmin }` from the cached `['me', clerkUserId]` query
- `frontend/src/lib/components/ColumnPickerDialog.svelte` — modal picker
- `frontend/tests/components/ColumnPickerDialog.test.ts`
- `frontend/tests/unit/grid-config.test.ts` (append if it exists)

**Modify:**
- `frontend/src/lib/queries/content.ts` — `addedByUserID` in `LIST_CONTENT` + `ContentItem`
- `frontend/src/lib/queries/users.ts` — `role` in `ME` + `Me` (+ `UserRole` type)
- `frontend/src/lib/utils/grid-config.ts` — `DATA_COLUMNS`, `INTERNAL_COLUMNS`, `togglableColIds`
- `frontend/src/lib/components/ActivityTable.svelte` — hidden colDefs (`id`, `addedByUserID`, `url`), move `createdAt`/`updatedAt` out of `alwaysHidden`, `useMe`, `userColumnOverride` state, `$effect` guard, gear button + `<ColumnPickerDialog>`
- `frontend/tests/unit/queries-content.test.ts` — assert `addedByUserID`
- `frontend/tests/unit/queries-users.test.ts` — assert `role` (create if absent)
- `frontend/tests/components/ActivityTable.test.ts` — gear button present when not `cardMode`
- `FEATURE_BACKLOG.md` — "Saved views" follow-up

---

## Tasks

### 1. Query + type changes
- [ ] `content.ts`: add `addedByUserID` to the `LIST_CONTENT` `items { ... }` selection and to `interface ContentItem` (`addedByUserID: string`).
- [ ] `users.ts`: add `export type UserRole = 'ADMIN' | 'SENTINEL' | 'DEFAULT';`, add `role: UserRole` to `interface Me`, add `role` to the `ME` query selection.
- [ ] Update `frontend/tests/unit/queries-content.test.ts` to assert the query string contains `addedByUserID`.
- [ ] Add/extend `frontend/tests/unit/queries-users.test.ts` to assert `ME` contains `role`.
- [ ] `pnpm --dir frontend run test:run` (unit) green for these files. Commit: `feat(content,users): expose addedByUserID and me.role to the client`.

### 2. Column registry in grid-config
- [ ] Add `DATA_COLUMNS`, `INTERNAL_COLUMNS` (readonly arrays of `{ colId, label }`) and `togglableColIds(isAdmin: boolean): string[]` to `frontend/src/lib/utils/grid-config.ts` per the spec.
- [ ] `frontend/tests/unit/grid-config.test.ts`: `togglableColIds(false)` has the 9 data colIds; `togglableColIds(true)` has 14; every `DATA_COLUMNS`/`INTERNAL_COLUMNS` colId is unique.
- [ ] Test green. Commit: `feat(grid-config): add togglable column registry`.

### 3. useMe hook
- [ ] Create `frontend/src/lib/queries/hooks/useMe.svelte.ts`: `useClerkContext()` for `clerkUserId`, `createQuery(() => ({ queryKey: ['me', clerkUserId], queryFn: () => graphqlRequest<MeResponse>(ME), enabled: clerk.isLoaded && !!clerkUserId, staleTime: 5*60*1000 }))`. Return an object with getters `me` (`meQuery.data?.me ?? null`) and `isAdmin` (`me?.role === 'ADMIN'`).
- [ ] No dedicated test (thin wrapper; exercised via ActivityTable/ColumnPickerDialog). Commit folded into task 5.

### 4. ColumnPickerDialog component
- [ ] Create `frontend/src/lib/components/ColumnPickerDialog.svelte`. Props: `open: boolean`, `onOpenChange: (v: boolean) => void`, `isAdmin: boolean`, `visibility: Record<string, boolean>` (colId → currently visible), `overrideActive: boolean`, `onToggle: (colId: string, next: boolean) => void`.
- [ ] Build on `$lib/components/shadcn/dialog` (follow `AddVideoDialog.svelte`). Title "Columns". Render a checkbox list for `DATA_COLUMNS`; when `isAdmin`, a second labelled group for `INTERNAL_COLUMNS`. Each checkbox `checked={visibility[colId] ?? false}`, `onchange` → `onToggle(colId, e.currentTarget.checked)`.
- [ ] When `overrideActive`, show the notice text: "Columns are set manually for this session. Refresh the page to restore automatic responsive layout."
- [ ] `frontend/tests/components/ColumnPickerDialog.test.ts`: data group renders; Internal group hidden when `isAdmin=false` and shown when `true`; toggling a checkbox calls `onToggle` with `(colId, boolean)`; notice shows only when `overrideActive`.
- [ ] Prettier + `check` + component test green. Commit: `feat(ColumnPickerDialog): modal column show/hide picker`.

### 5. Wire into ActivityTable
- [ ] Add hidden colDefs to `columnDefs`: `{ colId: 'id', field: 'id', headerName: 'Content ID', hide: true, sortable: false, filter: false, minWidth: 260 }`; `{ colId: 'addedByUserID', field: 'addedByUserID', headerName: 'Submitter', hide: true, sortable: false, filter: false, minWidth: 120 }`; `{ colId: 'url', field: 'url', headerName: 'Source URL', hide: true, sortable: false, filter: false, minWidth: 240 }`. Place after `description`.
- [ ] In the responsive `$effect`, remove `createdAt`/`updatedAt` from `alwaysHidden` (they remain `hide: true` in their colDefs, so they stay hidden until an admin enables them). Keep `description` in `alwaysHidden`.
- [ ] Add `import { useMe } from '$lib/queries/hooks/useMe.svelte';` and `const meCtx = useMe();` — use `meCtx.isAdmin`.
- [ ] Add `let userColumnOverride = $state<Record<string, boolean> | null>(null);` and `const overrideActive = $derived(userColumnOverride !== null);`
- [ ] Guard the responsive `$effect`: `if (userColumnOverride) return;` immediately after the `if (!gridApi || !gridReady) return;` line.
- [ ] `let columnPickerOpen = $state(false);`
- [ ] `function currentVisibility(): Record<string, boolean>` — from `gridApi.getColumnState()`, map `colId → !state.hide`, restricted to `togglableColIds(meCtx.isAdmin)`.
- [ ] `function handleColumnToggle(colId, next)`: guard `gridApi && gridReady`; if `userColumnOverride === null` seed it from `currentVisibility()`; `userColumnOverride = { ...userColumnOverride, [colId]: next }`; `gridApi.setColumnsVisible([colId], next)`.
- [ ] In the bottom toolbar left cluster (near the page-size `<select>`), add a gear icon button (`@lucide/svelte/icons/sliders-horizontal`), `aria-label="Choose columns"`, `disabled={!gridReady}`, hidden when `cardMode` (`{#if !cardMode}`), `onclick={() => (columnPickerOpen = true)}`.
- [ ] Render `<ColumnPickerDialog open={columnPickerOpen} onOpenChange={(v) => (columnPickerOpen = v)} isAdmin={meCtx.isAdmin} visibility={columnPickerOpen ? currentVisibility() : {}} overrideActive={overrideActive} onToggle={handleColumnToggle} />`.
- [ ] `frontend/tests/components/ActivityTable.test.ts`: assert a `[aria-label="Choose columns"]` button renders (grid path, not card mode). Keep the rest of the suite passing.
- [ ] Prettier + `check` + `test:run` (full) green. Commit: `feat(ActivityTable): session column picker with admin internal columns`.

### 6. Backlog note
- [ ] Append to `FEATURE_BACKLOG.md`: "**Saved views** — named, persisted presets of column visibility + filters + sort for the ActivityTable (localStorage or backend). Prompted by the session-only column picker (2026-09-02)."
- [ ] Commit: `docs(backlog): add Saved views follow-up`.

### 7. Full verification
- [ ] `pnpm --dir frontend run check` — zero errors.
- [ ] `pnpm --dir frontend run test:run` — all pass; report the summary counts.
- [ ] No backend verification needed (no backend changes).
- [ ] Grep the repo for any remaining `qmd` references introduced/missed; confirm none in tracked files outside `.planning/STATE.md` (historical) and `.claude/worktrees/*` (other checkouts).

---

## Manual QA (post-merge, not blocking)

- Desktop: gear opens modal; toggling "Tags" hides/shows it live; notice appears after first toggle; refresh restores responsive layout.
- Resize to `md` after a manual toggle → columns do NOT snap back (override holds).
- Admin account: "Internal" group visible; enabling "Content ID" shows the UUID column. Non-admin: no "Internal" group.
- `< 860px` (cardMode): gear button absent.
