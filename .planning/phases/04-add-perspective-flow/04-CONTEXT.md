# Phase 4: Add Perspective Flow - Context

**Gathered:** 2026-02-18 (updated)
**Status:** Ready for planning

<domain>
## Phase Boundary

Users can create perspectives on content (ratings, Like, Review) AND create claims as new content entries — all from the Activity table row action. Includes perspective schema evolution (simplified reference columns, optional fields, JSONB custom fields), claim content type, and @reference system for inline content references.

</domain>

<decisions>
## Implementation Decisions

### Domain Model — Content vs Perspective (LOCKED)

- **Content** = things that exist independently (YouTube video, article, claim, team, player, etc.). External source material. Content hierarchy (team → players) is modeled as separate content items.
- **Perspective** = always a user's subjective take. Ratings, thumbs up/down, review text. Can reference one content item, one primary perspective, and optionally other related perspectives.
- **Claim** = a **Content entry** (type `claim`), NOT a perspective. "Bo Jackson ran 22.3 mph in game" is factual/assertional content that others can then create perspectives on (agree/disagree, rate confidence).
- A perspective-on-perspective is the correct term (not "reply") — e.g., critiquing someone's rating style, consistency, or thoroughness.

### Action Flow from Activity Table (LOCKED)

1. **No existing perspective:** Row shows "+" icon in the perspective column (header = Perspectize glasses icon)
2. **Click "+":** Opens perspective popover/dialog for that content item → user sets ratings, thumbs, text → **INSERT into `perspectives`** with `content_id`
3. **Has existing perspective:** Row shows silhouette-with-glasses icon instead of "+"
4. **Hover on silhouette:** Opens the form pre-populated with existing perspective data → edits are **UPDATE to existing `perspectives` row**
5. **Adding a claim** from within the form: Creates a **new `content` row** (type `claim`), not a perspective. The claim is associated with the parent content via @reference.

### Claim Creation (LOCKED)

- Claims are content entries with `content_type = 'claim'`
- Created from within the perspective form UI (exact trigger: Claude's discretion)
- A claim references its parent content (the video/item the user was looking at)
- Claims appear in the Activity table as their own rows (type: claim) — they're content like anything else
- Anyone can create perspectives (rate, agree/disagree) on a claim

### @Reference System (LOCKED)

Inline content references in claim text using `@this` / `@here` tokens:

- **Input:** User types `@this ran 22.3 mph in the 1987 game`
- **Storage:** Raw token preserved in text: `@this ran 22.3 mph in the 1987 game`
- **Display:** Token resolved to parent content name: `Bo Jackson ran 22.3 mph in the 1987 game`
- **Hover:** Hovering over resolved text shows it's a dynamic reference to the parent content

**Purpose:** Avoids redundant typing and keeps claims concise while maintaining the link to context. The parent content's name can change without updating claim text.

**Scope:** `@this` / `@here` references the parent content item only. More complex @mention patterns (referencing arbitrary content) are deferred.

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
6. **No arbitrary @mentions** — `@this`/`@here` only references parent content, not arbitrary content items

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
- **New row (no perspective):** "+" icon → click opens empty form → INSERT
- **Existing perspective row:** Silhouette-with-glasses icon → hover opens pre-populated form → UPDATE
- Column header: Perspectize glasses icon

### Schema Migration Strategy

- **Additive migration** — add `primary_perspective_id` and `related_perspective_ids` columns to existing `perspectives` table
- Add `claim` to `content_type` enum/values
- Existing `content_id` column unchanged — no migration risk for current data
- Add GIN index on `related_perspective_ids`
- Zero impact on existing frontend queries and components

### Claude's Discretion

- Stepper click granularity (0.10, 0.25, or 0.50)
- Exact popover sizing and layout
- Loading skeleton design
- Error state handling
- `ON DELETE` behavior for `primary_perspective_id` FK (SET NULL vs CASCADE)
- Claim creation trigger within the perspective form (button, tab, or section)
- @reference token parsing implementation (regex, custom parser, etc.)

</decisions>

<specifics>
## Specific Ideas

- The "Perspectize glasses" icon is the column header; "+" icon per row is the add trigger
- Silhouette-with-glasses icon replaces "+" when user already has a perspective on that row
- Hovering on silhouette-with-glasses opens the form pre-populated (edit mode)
- Term "perspective-on-perspective" preferred over "reply" or "response"
- Sarah's scenario: critiquing James's inconsistent ratings across two videos — set `primary_perspective_id` to James's first perspective, add his second to `related_perspective_ids`
- Content hierarchy (team → players) modeled as separate content items, not perspective references
- Bo Jackson example: video = content, "ran 22.3 mph" = claim (content type), user's agreement with that claim = perspective
- `@this ran 22.3 mph in the 1987 game` → displays as `Bo Jackson ran 22.3 mph in the 1987 game` with hover showing the reference link

</specifics>

<deferred>
## Deferred Ideas

- Perspective threading UI (viewing chains of perspective-on-perspectives) — future phase
- Cross-referencing UI (selecting multiple items to compare) — future phase
- Custom field definitions/schema per user or community — future phase
- Content-to-content relationships (team → players, playlist → videos) — future content hierarchy phase
- Metadata on perspective references (rebuttal vs endorsement classification) — future phase if needed
- Array cap revisit (currently 50, may adjust based on usage patterns)
- Arbitrary @mentions referencing any content item (not just parent) — future enhancement
- Claim-specific fields beyond text (e.g., source URL, evidence type) — future phase

</deferred>

---

*Phase: 04-add-perspective-flow*
*Context gathered: 2026-02-18*
