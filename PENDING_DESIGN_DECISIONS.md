# Pending Decisions

Patches and deferred choices that work today but need a proper long-term solution.

## Design

| ID | Problem | Patch | SHA | Long-term Solution |
|----|---------|-------|-----|-------------------|
| PDD-001 | Sticky header used `bg-background` but `--color-background` was never defined in `app.css`, causing transparent header and scroll bleed-through | Changed `bg-background` to `bg-white` in `Header.svelte` | b42c457 | Define complete color theme in `app.css` (`--color-background`, `--color-foreground`, `--color-border`, etc.) so semantic utilities resolve correctly and support dark mode. Revert header to `bg-background`. |

## Data & Pagination

| ID | Problem | Patch | Issue | Long-term Solution |
|----|---------|-------|-------|-------------------|
| PDD-002 | `ListContent` query fetches `first: 100` but AG Grid only shows 10 per page. Client-side pagination/filtering works over the pre-fetched set, but total count isn't exposed and fetch size isn't tied to page size. | `first: 100` hardcoded — covers all client-side page sizes (10/25/50) with room for search filtering. Acceptable for MVP. | [#56](https://github.com/CodeWarrior-debug/perspectize-be/issues/56) | **Adaptive prefetch with exposed total count:** 1) Expose `totalCount` to the UI — if not provided by query, don't show. Total count = total available server-side, not total loaded. 2) Allow items-per-page to be user-specified (max 100). 3) Prefetch a multiple of page size: 1–10 items/page → fetch 5x (up to 50); 11–33 items/page → fetch 3x (up to 100); 34–100 items/page → fetch 1.5x (up to 100). This keeps client-side filtering snappy while avoiding unbounded overfetch. |
