# Feature Backlog

Ideas and future enhancements captured during development. Not committed to any milestone — evaluated when planning new work.

---

## Discover Page (New Page)

The v1 home page is an **Activity** page — a data table of user activity on videos already in the system. This is the correct approach for v1.

A future **Discover** page would be a separate page for finding new content outside the system:
- **Browse** — Show topics/tags from YouTube API endpoint, letting users explore categories
- **Search** — Live YouTube API search to discover new videos directly from YouTube

This is distinct from the Activity page's local search/filter. Discover reaches out to YouTube; Activity shows what's already tracked.

---

## gorm-cursor-paginator Integration (HIGH PRIORITY)

Phase 7.1 research recommended [gorm-cursor-paginator](https://github.com/pilagod/gorm-cursor-paginator) (226 stars) to replace hand-rolled cursor encoding in GORM repositories. The executor kept the existing `encodeCursor`/`decodeCursor` helpers from the sqlx migration instead of adding the library.

**Current state:** Hand-rolled cursor logic in `backend/internal/adapters/repositories/postgres/helpers.go` — works but misses benefits of the library (type-safe cursor fields, automatic keyset query building, simpler pagination setup).

**What to do:**
- Add `gorm-cursor-paginator` dependency
- Replace `encodeCursor`/`decodeCursor` and manual keyset WHERE clauses in `gorm_content_repository.go` and `gorm_perspective_repository.go`
- Simplify List methods to use paginator's built-in cursor handling
- Update tests

**Priority:** High — should be addressed in the next backend phase. Was part of the original 7.1 plan but was missed during execution.

---

## AG Grid Power Features Toolbar

Add a toolbar above `ActivityTable` with power-user grid controls. All features below use **AG Grid Community APIs** — no Enterprise license needed.

- **Clear all filters** — `gridApi.setFilterModel(null)`
- **Clear single column filter** — `gridApi.setColumnFilterModel('colId', null)`
- **Multi-column sort** — `multiSortKey: 'ctrl'` in gridOptions (hold Ctrl+click headers)
- **Column show/hide picker** — `gridApi.setColumnsVisible(['col1', 'col2'], true/false)`
- **Save/restore filter state** — `gridApi.getFilterModel()` / `gridApi.setFilterModel(saved)`
- **Save/restore column state** — `gridApi.getColumnState()` / `gridApi.applyColumnState({state})`

**References:**
- [AG Grid Filter API](https://www.ag-grid.com/javascript-data-grid/filter-api/)
- [AG Grid Column State](https://www.ag-grid.com/javascript-data-grid/column-state/)
- [AG Grid Multi-Sort](https://www.ag-grid.com/javascript-data-grid/row-sorting/#multi-column-sorting)
