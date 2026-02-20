# Phase 4: Add Perspective Flow - Context

**Gathered:** 2026-02-16
**Updated:** 2026-02-20
**Status:** Ready for planning (Phase 4A only — 4B depends on category/grouping architecture)

<domain>
## Phase Boundary

Users can create perspectives on videos with ratings, Like text, Review text, and validation. **Expanded scope:** Phase 4 also includes the perspective schema redesign (many-to-many references, self-joins, optional fields, JSONB custom fields) to support the full Perspectize vision.

**Phase 4 split (decided 2026-02-20):**
- **Phase 4A (independent):** Create perspective CRUD — form, ratings, Like/Review text, validation, backend mutation, PerspectivePopover UI, RatingInput component. Maps to plans 04-01 and popover portion of 04-02.
- **Phase 4B (depends on Phase 13/14):** How perspectives display in AG Grid with grouping — Perspectize column, expandable details, grouped views. Maps to AG Grid portion of 04-02.
- **Plan 04-03 (claims):** Independent — ties into reference architecture.

</domain>

<decisions>
## Implementation Decisions

### Domain Model — Content vs Perspective

- **Content** = things that exist independently (YouTube video, article, claim, etc.). External source material.
- **Perspective** = always a user's take. Can reference content, other perspectives, or both.
- A "claim" is Content (something standalone and ratable), NOT a Perspective.
- A perspective-on-perspective is the correct term (not "reply") — e.g., critiquing someone's rating style, consistency, or thoroughness.

### Reference Architecture — Many-to-Many (LOCKED)

- **PerspectiveReference join table:** `(perspective_id, ref_type ['content' | 'perspective'], ref_id)`
- A single perspective can reference 0+ content items AND 0+ other perspectives
- Replaces the current simple `content_id` FK on perspectives
- Supports:
  - Basic case: rate a video (1 content ref)
  - Perspective-on-perspective: critique someone's take (1 perspective ref)
  - Both: disagree with someone's take AND rate the video (1 content + 1 perspective ref)
  - Comparative: "these two videos contradict each other" (2+ content refs)
  - Synthesis: "these three takes together paint a picture" (2+ perspective refs)
- **This is core to Perspectize** — connecting perspectives across content is the product's value.

### Perspective Fields — All Optional + JSONB

- All standard fields (Quality, Agreement, Importance, Confidence, Like, Review) are OPTIONAL
- JSONB custom fields column for user-defined fields
- **Validation:** At least one field must be non-empty to submit (any rating OR any text)
- Ratings stored as integer 0–10000 in DB (hundredths precision), displayed as 0.00–10.00 to user

### Rating Input

- Number input with stepper/clicker control for quick adjustment without typing
- Display: 0.00 to 10.00 (two decimal places)
- Storage: integer 0 to 10000
- Step granularity: Claude's discretion (user deferred this)

### Form Location & Trigger

- Custom popover, bigger than Add Video's (needs room for ratings + text)
- Consistent with existing popover pattern but larger
- Becomes dialog at <768px (existing mobile pattern from FormPopover)
- **Trigger:** "+" icon in each Activity table row, under a column with Perspectize glasses icon as header
- Clicking "+" on a row opens the popover for that specific video

### Claude's Discretion

- Stepper click granularity (0.10, 0.25, or 0.50)
- Exact popover sizing and layout
- Loading skeleton design
- Error state handling
- Schema migration strategy (additive columns vs table replacement)

</decisions>

<cross_phase_architecture>
## Cross-Phase Architecture: Faceted Search & Grouping (decided 2026-02-20)

Phases 4B, 13, 14, and 15 share a unified data retrieval architecture. Decisions documented here apply across all four phases.

### Content Organization Model

**Primary category** (single, hierarchical):
- One category per content item via `category_id` FK
- Hierarchical path via PostgreSQL ltree extension
- Source: Google NL taxonomy (curated 20 of 27 categories) + user-created categories
- This is what AG Grid groups by (expandable tree rows)

**Labels** (multiple, flat):
- Many labels per content item
- Non-hierarchical — flat text for search/filtering, not grouping
- Case-insensitive comparison for searches
- Storage: separate labels table OR JSONB array with required `value` field + optional metadata (Claude's discretion)
- Supports faceted filtering alongside category grouping

### Faceted Search Pattern

The data retrieval flow:
1. **Search bar** — text query hits backend
2. **Category facet** — filter by taxonomy node (Sports > NFL)
3. **Label facets** — filter by labels (Mahomes, 2024 season)
4. **All combine** — backend returns intersection
5. **AG Grid displays** — groups by category hierarchy (expandable tree rows)
6. **Grid filters refine** — client-side further narrowing on returned results

### AG Grid Grouping Architecture

- **Expandable tree rows** — top level shows categories, expand to subcategories, expand to content items
- **API design:** Backend supports both flat queries (client-side grouping) and grouped/paginated queries (server-side grouping)
- **Start client-side** — fetch content with category paths, AG Grid groups client-side
- **Upgrade to server-side** — when data volume warrants, switch to Server-Side Row Model (API already supports it)
- **Toggle capability** — design API to support both patterns from day one

### Google NL Taxonomy — Research Spike Needed

**BLOCKER:** Before committing to the category architecture, a research spike is needed to understand:
- Full Google NL taxonomy depth (how many levels? what subcategories exist?)
- How to traverse layers (see all at one level, drill deeper, go higher)
- How to map YouTube content to taxonomy nodes
- Whether the taxonomy hierarchy maps cleanly to ltree paths
- Cost/API limits for automated classification

**Action:** Insert a spike phase (Phase 13.1 or similar) for taxonomy research before Phase 13 execution.

</cross_phase_architecture>

<specifics>
## Specific Ideas

- The "Perspectize glasses" icon is the column header; "+" icon per row is the action trigger
- Term "perspective-on-perspective" preferred over "reply" or "response"
- Sarah's scenario: critiquing James's inconsistent ratings across two videos — this is the canonical example of a multi-reference perspective (references 2 perspectives + 0 content, or 2 perspectives + 2 content items)
- NFL athlete example as canonical grouping use case: Patrick Mahomes video categorized under Sports > NFL > Chiefs, with labels ["Mahomes", "2024 season", "game film"]
- Faceted search like Amazon product filtering — combine search + category + labels to query backend

</specifics>

<deferred>
## Deferred Ideas

- Standalone claims as content type (content without external source) — future content type phase
- Perspective threading UI (viewing chains of perspective-on-perspectives) — future phase
- Cross-referencing UI (selecting multiple items to compare) — future phase
- Custom field definitions/schema per user or community — future phase

</deferred>

---

*Phase: 04-add-perspective-flow*
*Context gathered: 2026-02-16*
*Updated: 2026-02-20*
