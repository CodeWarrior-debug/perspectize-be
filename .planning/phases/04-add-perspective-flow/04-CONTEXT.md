# Phase 4: Add Perspective Flow - Context

**Gathered:** 2026-02-18 (updated from 2026-02-16)
**Status:** Ready for planning

<domain>
## Phase Boundary

Users can create perspectives on videos with ratings, Like text, Review text, and validation. Phase 4 also includes the perspective schema evolution (simplified reference columns, optional fields, JSONB custom fields) to support the full Perspectize vision.

</domain>

<decisions>
## Implementation Decisions

### Domain Model — Content vs Perspective

- **Content** = things that exist independently (YouTube video, article, claim, team, player, etc.). External source material. Content hierarchy (team → players) is modeled as separate content items.
- **Perspective** = always a user's take. Can reference one content item, one primary perspective, and optionally other related perspectives.
- A "claim" is Content (something standalone and ratable), NOT a Perspective.
- A perspective-on-perspective is the correct term (not "reply") — e.g., critiquing someone's rating style, consistency, or thoroughness.

### Reference Architecture — Simplified FK + Array (LOCKED)

**Replaces the former many-to-many join table decision.**

Columns on `perspectives` table:
- **`content_id`** (int FK, optional) — references one content item. Already exists in current schema.
- **`primary_perspective_id`** (int FK, optional, NEW) — references one other perspective as the primary relationship.
- **`related_perspective_ids`** (int array, optional, NEW) — additional related perspective IDs. Max 50 entries (application-level cap, revisit if needed).

**Why this over a join table:**
- Single atomic INSERT (no transactional coupling between perspective + reference rows)
- Zero JOINs for the common case (all reference data on the row itself)
- Real FK constraints on `content_id` and `primary_perspective_id` (PostgreSQL enforces referential integrity)
- No polymorphic `ref_id` column (which can't have FK constraints)
- Native PostgreSQL array with GIN index for reverse lookups

**No bridge tables needed.** All reference data lives on the perspective row.

**Supports all use cases:**
- Rate a video: set `content_id` (existing behavior, unchanged)
- Perspective-on-perspective: set `primary_perspective_id`
- Both: set `content_id` + `primary_perspective_id`
- Synthesis ("these takes together..."): set `primary_perspective_id` + populate `related_perspective_ids`

**Stale reference handling:** If a perspective in `related_perspective_ids` is deleted, the ID remains in the array. Resolver returns null for deleted IDs; frontend skips them. No cleanup job needed. The FK on `primary_perspective_id` gets standard PostgreSQL cascade behavior.

**GIN index required** on `related_perspective_ids` for reverse lookups (`WHERE X = ANY(related_perspective_ids)`).

### What We Are NOT Supporting

1. **No content-to-content relationships** in the perspective model — content hierarchy is a separate domain concern for a future phase
2. **No perspective threading/chains UI** — perspectives can reference each other, but no conversation view or tree rendering
3. **No cross-referencing picker UI** — no "select 3 perspectives to synthesize" interface
4. **No metadata on references** — no "type of relationship" between perspectives (rebuttal vs endorsement). Just "related."
5. **No recursive queries** — we won't walk chains of perspective→perspective→perspective

### Key Type — int (unchanged)

- Keep integer primary keys for perspectives (consistent with users and content)
- No need for UUID — not distributed, auth gates access (Phase 12), well within int4 capacity
- Migrate to bigint later if needed

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

### Schema Migration Strategy

- **Additive migration** — add `primary_perspective_id` and `related_perspective_ids` columns to existing `perspectives` table
- Existing `content_id` column unchanged — no migration risk for current data
- Add GIN index on `related_perspective_ids`
- Zero impact on existing frontend queries and components

### Claude's Discretion

- Stepper click granularity (0.10, 0.25, or 0.50)
- Exact popover sizing and layout
- Loading skeleton design
- Error state handling
- `ON DELETE` behavior for `primary_perspective_id` FK (SET NULL vs CASCADE)

</decisions>

<specifics>
## Specific Ideas

- The "Perspectize glasses" icon is the column header; "+" icon per row is the action trigger
- Term "perspective-on-perspective" preferred over "reply" or "response"
- Sarah's scenario: critiquing James's inconsistent ratings across two videos — set `primary_perspective_id` to James's first perspective, add his second to `related_perspective_ids`
- Content hierarchy (team → players) modeled as separate content items, not perspective references

</specifics>

<deferred>
## Deferred Ideas

- Standalone claims as content type (content without external source) — future content type phase
- Perspective threading UI (viewing chains of perspective-on-perspectives) — future phase
- Cross-referencing UI (selecting multiple items to compare) — future phase
- Custom field definitions/schema per user or community — future phase
- Content-to-content relationships (team → players, playlist → videos) — future content hierarchy phase
- Metadata on perspective references (rebuttal vs endorsement classification) — future phase if needed
- Array cap revisit (currently 50, may adjust based on usage patterns)

</deferred>

---

*Phase: 04-add-perspective-flow*
*Context gathered: 2026-02-18*
