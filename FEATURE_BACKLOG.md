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

## Authentication Architecture Design

Discovered during frontend caching review (2026-02-14). The GraphQL client (`frontend/src/lib/queries/client.ts`) has empty `headers: {}` — no auth tokens, no CSRF protection, no per-user cache scoping. Designing the auth architecture involves:

- **Token strategy:** JWT vs session cookies vs OAuth2
- **GraphQL client auth hook:** Dynamic header injection via `requestMiddleware` or client factory
- **Cache scoping:** TanStack Query keys need user identity dimension (e.g., `['content', userId]`)
- **Cache invalidation on logout:** `queryClient.clear()` to prevent data leakage between users
- **CSRF protection:** Backend middleware + frontend token handling
- **Secure token storage:** httpOnly cookies preferred over localStorage/sessionStorage

**Dependencies:** Should be planned alongside Phase 9 (Security Hardening) which covers backend auth middleware.

**Source:** Frontend caching review Finding #4 (CVSS 8.1), Finding #2 (no auth).

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

---

## Compress/Trim YouTube Raw JSONB Response

The `content.response` JSONB column stores the full YouTube Data API response and accounts for **93.7% of all content table data**. At 49 rows this is manageable but will scale poorly.

**Per-column byte analysis (49 rows):**

| Column | Total Bytes | % of Row Data |
|--------|------------|---------------|
| response (jsonb) | 118 KB | 93.7% |
| name | 2.4 KB | 1.9% |
| url | 2.2 KB | 1.7% |
| row overhead | 1.5 KB | 1.2% |
| all other columns | ~1.6 KB | 1.3% |

Average response: **2,469 bytes/row**. All other columns combined: **136 bytes/row**.

**Options:**
1. **Trim on ingest** — Store only the JSONB paths the app actually reads (`snippet.title`, `snippet.channelTitle`, `snippet.publishedAt`, `snippet.description`, `snippet.tags`, `statistics.*`) and drop unused nested objects (`contentDetails`, `status`, `topicDetails`, `recordingDetails`, etc.)
2. **Extract to columns** — Promote frequently queried JSONB paths into proper columns (the GraphQL schema already exposes `viewCount`, `likeCount`, `commentCount`, `channelTitle`, `publishedAt`, `tags`, `description` as resolved fields). Keep a trimmed `response` as fallback.
3. **Compress** — Use `pg_lz_compress` or application-level compression for the raw response if full fidelity is needed for audit/replay.

**Priority:** Low — not a problem at current scale (49 rows, 8 MB DB). Revisit when content table approaches 1,000+ rows.
