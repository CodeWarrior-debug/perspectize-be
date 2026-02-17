# Phase 4: Add Perspective Flow - Context

**Gathered:** 2026-02-16
**Status:** In progress — architecture locked, UI decisions partial, performance discussion pending

<domain>
## Phase Boundary

Users can create perspectives on videos with ratings, Like text, Review text, and validation. **Expanded scope:** Phase 4 also includes the perspective schema redesign (many-to-many references, self-joins, optional fields, JSONB custom fields) to support the full Perspectize vision.

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

<specifics>
## Specific Ideas

- The "Perspectize glasses" icon is the column header; "+" icon per row is the action trigger
- Term "perspective-on-perspective" preferred over "reply" or "response"
- Sarah's scenario: critiquing James's inconsistent ratings across two videos — this is the canonical example of a multi-reference perspective (references 2 perspectives + 0 content, or 2 perspectives + 2 content items)

</specifics>

<open_questions>
## Open Questions (Resume Here)

### Performance Concerns
User wants to discuss performance implications of:
- Many-to-many join table for references (query complexity, N+1, indexing)
- JSONB custom fields (indexing, querying, schema-less data)
- Recursive queries for perspective threads/chains
- Impact on AG Grid / Activity table rendering with richer data

### UI/UX Still to Discuss
- What the popover layout looks like (rating arrangement, text field placement)
- Whether to show existing perspectives on a video before adding yours
- How perspectives display in the Activity table (new column? expandable row?)
- Empty states (video with no perspectives yet)

### Schema Migration
- How to migrate from current `content_id` FK to PerspectiveReference join table
- Whether to keep backward compatibility or do a clean migration
- Impact on existing frontend queries and components

</open_questions>

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
*Status: PARTIAL — resume discussion before planning*
