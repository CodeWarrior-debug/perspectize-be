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

---

## Multi-Content-Type Schema Design (JSONB Philosophy)

When adding content types beyond YouTube (books, articles, podcasts, etc.), the `content` table needs a clear column vs. JSONB strategy. The guiding principle: **shared fields get columns, type-specific fields stay in JSONB.**

### Column Promotion Candidates

Fields that are common across content types should be promoted to dedicated columns for indexing, sorting, and type safety:

| Promoted Column | YouTube Source | Book Source | Data Type (TBD) |
|----------------|---------------|-------------|------------------|
| `name` | snippet.title | volumeInfo.title | already exists |
| `image_url` | snippet.thumbnails | volumeInfo.imageLinks | already exists |
| `creator` | snippet.channelTitle | volumeInfo.authors (array) | text vs text[] vs jsonb — needs discussion |
| `published_at` | snippet.publishedAt | volumeInfo.publishedDate | timestamptz vs date — YouTube has full timestamps, books often only have year |
| `description` | snippet.description | volumeInfo.description | text (potentially long) |
| `length` | contentDetails.duration (ISO 8601) | volumeInfo.pageCount (integer) | text vs integer — needs discussion (heterogeneous units) |

**Open questions on data types:**
- `creator` — Single author vs. multiple authors (YouTube has one channel, books have co-authors). Store as `text` (comma-joined) or `text[]` or keep in JSONB?
- `published_at` — YouTube provides RFC 3339 timestamps; books may only have "2024" or "2024-03". Use `timestamptz` with day-level precision fallback, or `text`?
- `length` — Duration (PT15M33S) and page count (384) are fundamentally different units. Store as `text` with a `length_unit` column? Or keep type-specific in JSONB?

### JSONB-Only Fields (Type-Specific)

Fields that belong to one content type and don't generalize stay in the JSONB `metadata` column:

| Field | Content Type | Why JSONB |
|-------|-------------|-----------|
| `viewCount`, `likeCount`, `commentCount` | YouTube | Video engagement metrics — no book equivalent |
| `isbn`, `isbn13` | Book | Identifier unique to books |
| `publisher` | Book | Books have publishers; YouTube has channels (→ `creator`) |
| `pageCount` | Book | Could go either way — see `length` discussion above |
| `language` | Book | YouTube has `defaultAudioLanguage` but rarely useful |
| `categories` / `tags` | Both | Similar concept but different taxonomies per source |

### Book API Recommendation: Google Books API

For the first non-YouTube content type, **Google Books API** is the recommended data source:

- **Free tier** — 1,000 req/day without API key, higher with key
- **Shared infrastructure** — Same Google API ecosystem as YouTube (API key, client libs, error patterns, rate limits)
- **Rich metadata** — title, authors, description, thumbnails, page count, categories, publisher, ISBNs, ratings, preview links
- **Similar response structure** — `items[]` with `volumeInfo` parallels YouTube's `items[]` with `snippet`

**Other APIs evaluated:**
- **Open Library** (openlibrary.org) — Free, no key, 20M+ editions. Good fallback if Google quota is insufficient. Less structured metadata.
- **Library of Congress** — Free, authoritative for bibliographic data, but clunky API and inconsistent response formats. Better for archival than consumer use.
- **Goodreads** — API shut down by Amazon in 2020. Not available.
- **ISBNdb** — Paid ($9/mo+). Not worth it when free options exist.

**Priority:** Medium — design the schema strategy before adding the second content type. The column promotion migration (creator, published_at, description) should happen during the JSONB trim work above.

**Source:** Architecture discussion (2026-02-16)

---

## Server-Side Sorting and Filtering for Activity Table

Currently AG Grid sorting and filtering is client-side only — operates on the current page of rows. With server-side pagination (`limit`/`offset`), this means sort/filter only affects the visible page, not the full dataset.

**What to do:**
- Add `orderBy` and `filter` parameters to the `contents` GraphQL query
- Implement sort/filter in the backend resolver and GORM repository
- Wire AG Grid's `onSortChanged`/`onFilterChanged` events to re-fetch with server parameters
- Consider switching to AG Grid's server-side row model for seamless integration

**Priority:** Medium — noticeable UX gap once content exceeds one page (currently 52 items across 6 pages). Sorting "Views" only reorders the 10 visible rows, not all 52.

**Source:** Phase 3.2 UAT (2026-02-16)

---

## Slow COUNT(*) and JSONB Sort Queries

The `SELECT count(*) FROM "content"` query and JSONB-path ORDER BY queries are hitting 200-400ms, triggering GORM's slow query warning (>= 200ms). Observed at ~50 rows — will worsen at scale.

**Slow queries observed:**
- `SELECT count(*) FROM "content"` — 203ms to 382ms
- `SELECT * FROM "content" ORDER BY COALESCE(response->'items'->0->'snippet'->>'publishedAt', '') ASC` — 157ms
- Full GraphQL request latency: 395-541ms

**Potential fixes:**
1. **Add indexes** — GIN index on `response` JSONB, or expression indexes on frequently sorted JSONB paths (`(response->'items'->0->'statistics'->>'viewCount')::BIGINT`, `response->'items'->0->'snippet'->>'publishedAt'`)
2. **Extract to columns** — Promote JSONB-derived sort fields to dedicated columns (related to "Compress/Trim YouTube Raw JSONB Response" item above)
3. **Cache total count** — Avoid `SELECT count(*)` on every paginated request; use estimated count or cache with TTL
4. **Analyze** — Run `EXPLAIN ANALYZE` on slow queries to confirm whether it's sequential scan or JSONB extraction overhead

**Priority:** Low — manageable at 50 rows but a known scaling bottleneck. Consider alongside JSONB trimming/extraction work.

**Source:** Backend logs (2026-02-15)
