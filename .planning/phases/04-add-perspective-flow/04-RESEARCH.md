# Phase 4: Add Perspective Flow - Research

**Researched:** 2026-02-18
**Domain:** Full-stack perspective CRUD — PostgreSQL schema migration, Go hexagonal architecture, SvelteKit + AG Grid popover UI
**Confidence:** HIGH (all findings verified from codebase inspection + existing patterns)

---

## Summary

This phase adds the "Perspectize" action to the Activity table: users can create and edit perspectives (ratings + text) on content items from a per-row popover. It also adds the claim content type and a simple @reference token system. The backend already has a complete perspective CRUD stack (domain, service, repository, resolvers, GraphQL schema). The frontend has zero perspective-facing UI beyond what the backend exposes — everything is greenfield on that side.

The primary backend work is a schema migration (additive: two new FK columns on `perspectives`, a GIN index, JSONB custom fields column) plus new `claim` content type support. The primary frontend work is a new AG Grid column with custom cell renderers for +/silhouette-with-glasses icons, a new `PerspectivePopover` component (larger than existing FormPopover), two mutation hooks (`useCreatePerspective` / `useUpdatePerspective`), and a `perspectives.ts` query file. A new `claims.ts` query file handles claim creation as a content mutation.

**Primary recommendation:** Follow the established hex-clean pattern exactly. The perspective service already validates ratings and handles all optional fields. Add the new schema columns, update domain/GORM models + mappers, extend the GraphQL schema, and build the frontend popover following the AddVideoPopover pattern.

---

## Standard Stack

### Core (already in project — no new installs needed)

| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| gqlgen | current (schema-first) | Go GraphQL | Already wired; add new fields, run `make graphql-gen` |
| GORM + pgx/v5 | current | PostgreSQL ORM | Hex-clean pattern established in codebase |
| golang-migrate | current | SQL migrations | All 11 existing migrations use this |
| TanStack Svelte Query v5 | current | Data fetching + cache | Existing pattern: `createMutation` with function wrapper |
| graphql-request | current | GraphQL client | `graphqlClient.request()` used in all hooks |
| svelte-sonner | current | Toast notifications | Used in `useAddVideo`, `useCreateUser` |
| AG Grid Svelte5 | v32.x | Table column | `ag-grid-svelte5` wrapper, `cellRenderer` pattern |
| shadcn-svelte | current | UI primitives | `Popover`, `Dialog`, `Input`, `Label`, `Button` |
| Lucide Svelte | current | Icons | `@lucide/svelte/icons/{name}` per-icon imports |

### Supporting

| Library | Version | Purpose | When to Use |
|---------|---------|---------|-------------|
| Vitest + Testing Library | current | Unit tests | TEST-05 requires component + hook unit tests |
| pilagod/gorm-cursor-paginator/v2 | current | Cursor pagination | Already used in perspective repository List() |

### No New Dependencies Required

All libraries are already installed. No `npm install` or `go get` needed.

---

## Architecture Patterns

### Recommended Project Structure (new files only)

```
backend/
├── migrations/
│   └── 000012_add_perspective_refs_claims.up.sql   # New columns, GIN index, custom_fields JSONB
│   └── 000012_add_perspective_refs_claims.down.sql
├── internal/core/domain/
│   └── perspective.go           # Add PrimaryPerspectiveID, RelatedPerspectiveIDs, CustomFields
│   └── content.go               # Add ContentTypeClaim constant
├── internal/adapters/repositories/postgres/
│   └── gorm_models.go           # Add PerspectiveModel fields + new IntArray type
│   └── gorm_mappers.go          # Update bidirectional mappers
│   └── gorm_perspective_repository.go  # No changes needed (Save() covers new columns)
├── schema.graphql               # Add fields to Perspective type + inputs; add claim to ContentType enum
└── internal/adapters/graphql/resolvers/
    └── schema.resolvers.go      # Update CreatePerspective + UpdatePerspective resolvers

frontend/src/lib/
├── queries/
│   ├── perspectives.ts          # New: GQL definitions + TS interfaces
│   └── keys.ts                  # Add perspectives key factory
├── queries/hooks/
│   ├── useCreatePerspective.ts  # New: createMutation following useAddVideo pattern
│   └── useUpdatePerspective.ts  # New: createMutation for edit mode
├── components/
│   └── PerspectivePopover.svelte # New: larger FormPopover-like, CREATE + UPDATE modes

frontend/tests/
├── unit/
│   └── hooks-useCreatePerspective.test.ts
│   └── hooks-useUpdatePerspective.test.ts
└── components/
    └── PerspectivePopover.test.ts
```

### Pattern 1: Hex-Clean Domain Extension

**What:** Add new fields to domain model, GORM model, and mappers in sync — the established pattern from migration 000004 and existing perspective code.

**When to use:** Any new perspective column needs changes in all 3 layers.

**Example — domain.Perspective additions:**
```go
// Source: backend/internal/core/domain/perspective.go (extend existing struct)
type Perspective struct {
    // ... existing fields unchanged ...

    // Phase 4 additions
    PrimaryPerspectiveID   *int    // FK to perspectives.id (optional)
    RelatedPerspectiveIDs  []int   // int[] stored as PostgreSQL integer[] with GIN index
    CustomFields           []byte  // JSONB column for user-defined fields (json.RawMessage)
}
```

**Example — GORM model additions:**
```go
// Source: backend/internal/adapters/repositories/postgres/gorm_models.go (extend PerspectiveModel)
type PerspectiveModel struct {
    // ... existing fields unchanged ...

    // Phase 4 additions
    PrimaryPerspectiveID  *int       `gorm:"column:primary_perspective_id"`
    RelatedPerspectiveIDs Int64Array `gorm:"type:integer[];column:related_perspective_ids"`
    CustomFields          []byte     `gorm:"type:jsonb;column:custom_fields"`
}
```

**Example — migration:**
```sql
-- 000012_add_perspective_refs_claims.up.sql
ALTER TABLE public.perspectives
    ADD COLUMN IF NOT EXISTS primary_perspective_id integer NULL,
    ADD COLUMN IF NOT EXISTS related_perspective_ids integer[] NULL DEFAULT '{}',
    ADD COLUMN IF NOT EXISTS custom_fields jsonb NULL;

ALTER TABLE public.perspectives
    ADD CONSTRAINT perspectives_primary_perspective_fk
        FOREIGN KEY (primary_perspective_id)
        REFERENCES public.perspectives(id)
        ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS perspectives_related_perspective_ids_gin
    ON public.perspectives USING GIN (related_perspective_ids);
```

**Why SET NULL for ON DELETE on primary_perspective_id:** If a referenced perspective is deleted, the FK pointer becomes null automatically — the row continues to exist with just that reference cleared. CASCADE would delete the referencing perspective, which is destructive and undesirable.

### Pattern 2: GraphQL Schema Extension (additive)

**What:** Add new optional fields to existing `Perspective` type and `CreatePerspectiveInput`/`UpdatePerspectiveInput`. Do not break existing fields.

**Example — schema.graphql additions:**
```graphql
# Add to existing Perspective type
type Perspective {
  # ... existing fields ...
  primaryPerspectiveID: ID          # NEW optional
  relatedPerspectiveIDs: [ID!]      # NEW optional array
  customFields: JSON                # NEW optional JSONB
}

# Add to existing CreatePerspectiveInput
input CreatePerspectiveInput {
  # ... existing fields ...
  primaryPerspectiveID: IntID       # NEW optional
  relatedPerspectiveIDs: [IntID!]   # NEW optional array
  customFields: JSON                # NEW optional
  reviewText: String                # NEW — maps to existing `description` field
}

# Extend ContentType enum
enum ContentType {
  YOUTUBE
  CLAIM                             # NEW
}
```

After changes: `cd backend && make graphql-gen`

### Pattern 3: Shared Mutation Hook (useCreatePerspective)

**What:** TanStack `createMutation` with function wrapper, toast notifications, query cache invalidation. Mirrors `useAddVideo` exactly.

**When to use:** All perspective mutations follow this pattern.

**Example:**
```typescript
// Source: pattern from frontend/src/lib/queries/hooks/useAddVideo.ts
import { createMutation, useQueryClient } from '@tanstack/svelte-query';
import { toast } from 'svelte-sonner';
import { graphqlClient } from '../client';
import { CREATE_PERSPECTIVE, type CreatePerspectiveResponse } from '../perspectives';
import { queryKeys } from '../keys';
import { getSelectedUserId } from '$lib/stores/userSelection.svelte';

export function useCreatePerspective() {
    const queryClient = useQueryClient();

    return createMutation(() => ({
        mutationFn: async (input: { contentId: string; /* ...other fields */ }) => {
            const userId = getSelectedUserId();
            if (userId === null) throw new Error('No user selected');
            return graphqlClient.request<CreatePerspectiveResponse>(CREATE_PERSPECTIVE, {
                input: { userID: userId, contentID: Number(input.contentId), /* ... */ }
            });
        },
        onSuccess: () => {
            toast.success('Perspective saved');
            queryClient.invalidateQueries({ queryKey: queryKeys.perspectives.lists() });
        },
        onError: (err: Error) => {
            toast.error('Failed to save perspective. Please try again.');
        }
    }));
}
```

### Pattern 4: AG Grid Custom Cell Renderer for Perspective Column

**What:** A `cellRenderer` returning an `HTMLElement` with either a "+" icon or a silhouette-with-glasses icon, depending on whether the row's content has a perspective from the selected user.

**Critical constraint:** AG Grid cell renderers in this project return `HTMLElement` (see `itemCellRenderer`, `typeCellRenderer` in `formatting.ts`). They are NOT Svelte components — they are plain DOM elements. Popover components CANNOT be embedded inside a cell renderer.

**Resolution pattern:** The AG Grid column shows icon-only state. Clicking/hovering the icon triggers the popover by setting a Svelte state variable (e.g., `activePerspectiveContentId`) that controls the popover `open` state. The popover lives in the parent `ActivityTable.svelte` template (or inline on the page), not inside the cell.

**Example column def:**
```typescript
// In ActivityTable.svelte columnDefs
{
    colId: 'perspective',
    headerName: '',                   // No header text — icon only via headerComponent
    headerTooltip: 'Your perspective',
    width: 48,
    minWidth: 48,
    maxWidth: 48,
    sortable: false,
    filter: false,
    cellRenderer: perspectiveCellRenderer,
    onCellClicked: (params) => {
        // Set active row for popover (open if no existing perspective)
        if (!params.data.userPerspective) {
            activePerspectiveRow = params.data;
            perspectivePopoverOpen = true;
        }
    },
    // Hover for edit mode handled differently — see Pitfall 2
}
```

**Header component:** The Perspectize glasses icon as header requires a `headerComponent`. The simplest implementation is an inline `headerComponent` function returning an `HTMLElement` with the SVG icon, consistent with how `typeCellRenderer` works.

### Pattern 5: "Perspectize" Per-User Perspective Data in the Content Query

**What:** The content query must include each row's perspective for the selected user. This requires either:
- A separate `perspectives` query filtered by `userId + contentId`, or
- Embedding perspective data in the content query response.

**Chosen approach (research recommendation):** Separate query. The `content` query returns content items. For the perspective column, maintain a parallel `perspectives` query filtered by `userID` (the selected user). The frontend merges the two by `contentId` via a `Map<contentId, perspective>`.

**Why not embedding:** The GraphQL `Content` type does not have a perspectives field, and adding one changes the content resolver significantly. Keeping queries separate maintains clean separation and aligns with the existing pattern where users and content are separate queries.

**How the merge works:**
```typescript
// In ActivityTable.svelte or a derived computation
const userPerspectivesMap = $derived(() => {
    const map = new Map<string, Perspective>();
    for (const p of perspectivesQuery.data?.perspectives.items ?? []) {
        if (p.contentID) map.set(p.contentID, p);
    }
    return map;
});

// When rendering the perspective column:
const perspective = userPerspectivesMap.get(params.data.id);
// Shows "+" if null, silhouette-glasses if defined
```

### Pattern 6: PerspectivePopover Component

**What:** A larger version of `FormPopover.svelte` with two modes — create (empty) and edit (pre-populated). It is NOT using `FormPopover.svelte` directly because the perspective form is significantly more complex (multiple rating inputs, two text areas, claim section, validation).

**Structure:**
```svelte
<!-- PerspectivePopover.svelte -->
<script lang="ts">
    let {
        contentId,         // The content item being perspectized
        contentName,       // For display
        existingPerspective = null,  // null = create mode; object = edit mode
        onClose,
        open = $bindable(false),
    } = $props();

    // Two mutation hooks, one active based on mode
    const createMutation = useCreatePerspective();
    const updateMutation = useUpdatePerspective();
    // ...rating state, validation, submit handler
</script>
```

**Mobile behavior:** Becomes Dialog at <768px — same media query pattern as `FormPopover.svelte`. Copy the `$effect` media query pattern verbatim.

**Trigger mechanism:** Unlike FormPopover (which owns its trigger button), the PerspectivePopover is triggered externally by the AG Grid column click. The component receives `open` as a bindable prop and `contentId`/`existingPerspective` props set by the parent before opening.

### Pattern 7: @Reference Token Storage and Display

**What:** `@this` / `@here` tokens in claim text are stored raw. Display resolves them to the parent content's name. Implementation is simple string replacement at render time.

**Storage:** `INSERT` with raw text: `"@this ran 22.3 mph in the 1987 game"` — no parsing on write.

**Display resolver:**
```typescript
// Simple function — no library needed
export function resolveAtReference(
    text: string,
    parentContentName: string
): string {
    return text.replace(/@this|@here/g, parentContentName);
}
```

**Hover:** Wrap resolved text in a span with a tooltip showing "Reference to: {parentContentName}". Use a `DescriptionTooltip`-style pattern.

### Pattern 8: Claim Creation Flow

**What:** Within the `PerspectivePopover`, a "Add Claim" section (button or tab — Claude's discretion) opens a sub-form for entering claim text. On submit, calls `createContentFromYouTube`... wait — claims use a different mutation.

**Claim creation mutation:** Claims are content rows with `content_type = 'claim'`. The existing `createContentFromYouTube` mutation is YouTube-specific. Phase 4 needs either:
1. A new `createClaim` mutation in GraphQL (new resolver), or
2. Extending `createContentFromYouTube` to a generic `createContent` mutation.

**Recommendation:** Add a new `createClaim(input: CreateClaimInput!): Content!` mutation. Keep `createContentFromYouTube` unchanged to avoid breaking existing callers.

**`CreateClaimInput` fields:**
```graphql
input CreateClaimInput {
    text: String!          # The claim text (may contain @this/@here)
    userID: IntID!         # Who created it
    parentContentID: IntID! # The content item being claimed about
}
```

**Backend:** New resolver delegates to `ContentService.CreateClaim()` or a new claim-specific service method. The domain creates a `Content` row with `content_type = "claim"` and stores `parentContentID` as JSONB metadata (in the `response` column) since the `content` table has no FK column for parent content in this phase.

**Simpler alternative:** Since the `content` table's `response` column is JSONB, store `{ "parentContentId": 123, "text": "..." }` there. The claim's `name` field holds the display text (resolved from @reference at read time). This avoids any new columns.

### Anti-Patterns to Avoid

- **Svelte component as AG Grid cell renderer:** Cannot mount Svelte components in AG Grid cells in this project. Use `HTMLElement` returning functions and lift popover state to parent.
- **Embedding perspective data in content query:** Don't add `perspectives` to the `Content` GraphQL type — keeps the existing query shape clean and avoids N+1 issues.
- **FormPopover reuse for PerspectivePopover:** FormPopover's `onSubmit` signature and single-action flow don't fit the dual create/edit mode. Build `PerspectivePopover.svelte` as a standalone component that internally uses shadcn primitives directly.
- **Using Svelte 4 syntax:** All reactive state must use `$state`, `$derived`, `$effect`. No `export let`, no `$:`, no `on:event`.
- **Rating as float in GraphQL:** Keep ratings as `Int` in GraphQL (0–10000). Display conversion (divide by 100) is client-side only.

---

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| Integer array scanner for PostgreSQL | Custom string parser | Existing `Int64Array` type in `array_types.go` | Already handles `{1,2,3}` PostgreSQL format |
| Toast notifications | Custom notification component | `svelte-sonner` (`toast.success`, `toast.error`) | Already wired in layout |
| Popover/dialog mobile responsive | Custom modal logic | Copy media query `$effect` from `FormPopover.svelte` | Same 768px breakpoint pattern |
| Query cache invalidation | Manual state reset | `queryClient.invalidateQueries` with key factory | TanStack handles stale-while-revalidate |
| GraphQL type generation | Manual type writing | `make graphql-gen` after schema changes | gqlgen regenerates all model types |
| FK cascade behavior | Application-level deletion propagation | PostgreSQL `ON DELETE SET NULL` | DB enforces referential integrity |
| @reference resolution | Full mention/hashtag library | Simple `String.replace(/@this|@here/g, name)` | Scope is only two fixed tokens |

**Key insight:** The existing codebase has all the plumbing — new code is almost entirely additive. The repository's `Save()` (GORM) automatically handles new nullable columns once the GORM model has the fields.

---

## Common Pitfalls

### Pitfall 1: AG Grid Column Visibility Must Be Updated in TWO Places

**What goes wrong:** Adding the `perspective` column to `columnDefs` with `hide: false` shows it on desktop, but the `$effect` block that runs on `gridReady` overwrites visibility. If the new column is not added to the `$effect`'s `setColumnsVisible` calls, it may be hidden unexpectedly.

**Why it happens:** The `$effect` in `ActivityTable.svelte` calls `setColumnsVisible` with explicit column lists for both mobile and desktop — it always wins because it runs after grid ready.

**How to avoid:** Add `'perspective'` to the appropriate `setColumnsVisible(['perspective', ...others], true/false)` calls in BOTH the mobile and desktop branches of the `$effect`. Add to desktop-visible list.

**Warning signs:** Column appears in column menu but not in table after `gridReady`.

### Pitfall 2: Hover-to-Edit on an AG Grid Row Is Not Standard

**What goes wrong:** AG Grid has `onCellClicked` but no native "hover cell to open popover" API. Attempting to attach hover events directly to cell renderers leads to inconsistent behavior when the cell re-renders.

**Why it happens:** AG Grid recycles cell DOM nodes (virtual scrolling). Event listeners attached manually inside `cellRenderer` may be lost on re-render.

**How to avoid:** Use `onCellMouseOver` and `onCellMouseOut` grid events at the grid level (not per-cell) to detect hover, OR use a dedicated icon button inside the cell renderer that uses a `data-content-id` attribute, and attach a delegated event listener to the grid container. The recommended approach: in the `cellRenderer` function, return an `<button>` element. The cell click handler in `gridOptions.onCellClicked` dispatches to the right action based on whether `existingPerspective` is set.

**Alternative (simpler):** Use `onCellClicked` for both create and edit — single click opens the popover in create or edit mode. The hover UX can be approximated with a CSS `:hover` state on the cell icon (change opacity) without triggering the popover. This avoids the hover-open-popover complexity entirely and is more accessible.

**Warning signs:** Popover opens and closes unexpectedly, or doesn't open on hover.

### Pitfall 3: The `perspectives` Query Needs to Filter by Selected User

**What goes wrong:** If the perspectives query fetches all perspectives (no user filter), the perspective column would show all users' perspectives for a given content row, not just the current user's.

**Why it happens:** `PerspectiveFilter` accepts `userID` as optional. It must be passed.

**How to avoid:** The `useCreatePerspective` hook (and the perspectives query) must always filter by `getSelectedUserId()`. If no user is selected, show "+" icons in all rows (no perspective possible to display) and disable perspective creation with a toast.

**Warning signs:** Perspective icons show for wrong user's perspective.

### Pitfall 4: `make graphql-gen` Must Run After Schema Changes

**What goes wrong:** Editing `schema.graphql` does not automatically update `internal/adapters/graphql/model/models_gen.go` or the resolver method signatures. Compilation fails with "undefined field" or "missing method" errors.

**Why it happens:** gqlgen is schema-first; generated code is not auto-run.

**How to avoid:** After every `schema.graphql` change, run `cd backend && make graphql-gen`. Verify compilation with `go build ./...` immediately after.

**Warning signs:** `go build ./...` errors about missing resolver methods or unknown types.

### Pitfall 5: The Perspective Update Service Uses Partial Updates (Nil = No Change)

**What goes wrong:** When editing a perspective, sending a `nil` rating field means "don't change this field" — NOT "clear this field to null". If the UI allows removing a rating, the mutation input must explicitly set the field rather than omit it.

**Why it happens:** The existing `Update()` service in `perspective_service.go` skips nil fields. This is intentional for partial updates but means you cannot zero-out a rating by passing nil.

**How to avoid:** For the phase 4 use case (edit a perspective), pass all current field values in the update input, not just changed ones. The UI should send the current value of every field on every save — this is a "full update" pattern even though the service supports partial.

**Warning signs:** User removes a rating during edit but it persists after save.

### Pitfall 6: `content_type` Is VARCHAR, Not a PostgreSQL Enum

**What goes wrong:** Assuming `content_type = 'claim'` requires an `ALTER TYPE` enum migration.

**Why it happens:** Natural assumption when adding a new "enum value".

**How to avoid:** The `content` table's `content_type` is `varchar NOT NULL`. Adding `claim` as a valid value requires NO migration to the database column — only the Go domain constant and GraphQL `ContentType` enum need updating. The GraphQL enum `ContentType` is a gqlgen construct, not mapped from a DB enum.

**Warning signs:** Unnecessary migration that tries to `ALTER TYPE content_type ADD VALUE 'claim'` (the type doesn't exist as a PostgreSQL ENUM).

### Pitfall 7: `primary_perspective_id` FK Is a Self-Reference

**What goes wrong:** Writing the FK migration referencing `public.perspectives(id)` when the table is `perspectives` itself — this is a valid self-referential FK but requires the table to already exist (it does).

**How to avoid:** The FK must be added in a separate `ALTER TABLE` statement after the column is added (same migration file, but `ADD COLUMN` first, then `ADD CONSTRAINT`). PostgreSQL requires the column to exist before constraining it.

### Pitfall 8: Empty Perspective Validation Must Be Frontend AND Backend

**What goes wrong:** Only validating "at least one field non-empty" on the frontend. A direct API call could create an empty perspective row.

**How to avoid:** Add backend validation in `PerspectiveService.Create()`: if all ratings are nil AND `like` is nil AND `description` is nil AND `customFields` is nil, return `domain.ErrInvalidInput`. Frontend shows a toast (PERSP-08) before submission, backend enforces it as a safety net.

---

## Code Examples

Verified patterns from codebase inspection:

### Migration: Additive Column Addition (from migrations/000011)
```sql
-- Source: backend/migrations/000011_allow_null_user_email.up.sql (pattern)
ALTER TABLE public.perspectives
    ADD COLUMN IF NOT EXISTS primary_perspective_id integer NULL,
    ADD COLUMN IF NOT EXISTS related_perspective_ids integer[] NULL DEFAULT '{}',
    ADD COLUMN IF NOT EXISTS custom_fields jsonb NULL;

ALTER TABLE public.perspectives
    ADD CONSTRAINT perspectives_primary_perspective_fk
        FOREIGN KEY (primary_perspective_id)
        REFERENCES public.perspectives(id)
        ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS perspectives_related_ids_gin
    ON public.perspectives USING GIN (related_perspective_ids);
```

### GORM Model: Extending PerspectiveModel
```go
// Source: backend/internal/adapters/repositories/postgres/gorm_models.go (extend)
type PerspectiveModel struct {
    // ... existing fields ...
    PrimaryPerspectiveID  *int       `gorm:"column:primary_perspective_id"`
    RelatedPerspectiveIDs Int64Array `gorm:"type:integer[];column:related_perspective_ids"`
    CustomFields          []byte     `gorm:"type:jsonb;column:custom_fields"`
}
```

### Mapper: New Fields Pattern
```go
// Source: pattern from backend/internal/adapters/repositories/postgres/gorm_mappers.go
// In perspectiveModelToDomain():
p.PrimaryPerspectiveID = m.PrimaryPerspectiveID
if len(m.RelatedPerspectiveIDs) > 0 {
    p.RelatedPerspectiveIDs = make([]int, len(m.RelatedPerspectiveIDs))
    for i, v := range m.RelatedPerspectiveIDs {
        p.RelatedPerspectiveIDs[i] = int(v)
    }
}
p.CustomFields = m.CustomFields

// In perspectiveDomainToModel():
m.PrimaryPerspectiveID = p.PrimaryPerspectiveID
if len(p.RelatedPerspectiveIDs) > 0 {
    m.RelatedPerspectiveIDs = make(Int64Array, len(p.RelatedPerspectiveIDs))
    for i, v := range p.RelatedPerspectiveIDs {
        m.RelatedPerspectiveIDs[i] = int64(v)
    }
}
m.CustomFields = p.CustomFields
```

### Frontend: Perspectives Query File (new file)
```typescript
// Source: pattern from frontend/src/lib/queries/content.ts
import { gql } from 'graphql-request';

export interface PerspectiveItem {
    id: string;
    userID: string;
    contentID: string | null;
    quality: number | null;
    agreement: number | null;
    importance: number | null;
    confidence: number | null;
    like: string | null;
    description: string | null;
    primaryPerspectiveID: string | null;
    relatedPerspectiveIDs: string[];
    customFields: Record<string, unknown> | null;
    createdAt: string;
    updatedAt: string;
}

export interface PerspectivesResponse {
    perspectives: {
        items: PerspectiveItem[];
        pageInfo: { hasNextPage: boolean; hasPreviousPage: boolean; startCursor: string | null; endCursor: string | null; };
        totalCount: number;
    };
}

export interface CreatePerspectiveInput {
    userID: number;
    contentID?: number;
    quality?: number;
    agreement?: number;
    importance?: number;
    confidence?: number;
    like?: string;
    description?: string;
}

export const LIST_PERSPECTIVES = gql`
    query ListPerspectives($filter: PerspectiveFilter, $first: Int) {
        perspectives(filter: $filter, first: $first) {
            items {
                id
                userID
                contentID
                quality
                agreement
                importance
                confidence
                like
                description
                createdAt
                updatedAt
            }
            pageInfo { hasNextPage hasPreviousPage startCursor endCursor }
            totalCount
        }
    }
`;

export const CREATE_PERSPECTIVE = gql`
    mutation CreatePerspective($input: CreatePerspectiveInput!) {
        createPerspective(input: $input) {
            id
            userID
            contentID
            quality
            agreement
            importance
            confidence
            like
            description
            createdAt
            updatedAt
        }
    }
`;

export const UPDATE_PERSPECTIVE = gql`
    mutation UpdatePerspective($input: UpdatePerspectiveInput!) {
        updatePerspective(input: $input) {
            id
            quality
            agreement
            importance
            confidence
            like
            description
            updatedAt
        }
    }
`;
```

### Frontend: Query Keys Extension
```typescript
// Source: frontend/src/lib/queries/keys.ts (extend)
perspectives: {
    all: () => [...queryKeys.all, 'perspectives'] as const,
    lists: () => [...queryKeys.perspectives.all(), 'list'] as const,
    list: (filters: { userID?: number; contentID?: number }) =>
        [...queryKeys.perspectives.lists(), filters] as const,
},
```

### Frontend: Rating Input with Stepper (Svelte 5)
```svelte
<!-- Recommended stepper granularity: 0.25 (25 in integer units) -->
<!-- Source: Svelte 5 runes pattern from codebase -->
<script lang="ts">
    let quality = $state<number | null>(null);
    const STEP = 25; // 0.25 display steps = 25 integer units
    const MAX = 10000;

    function displayToInt(display: number): number {
        return Math.round(display * 100);
    }
    function intToDisplay(val: number | null): string {
        return val === null ? '' : (val / 100).toFixed(2);
    }
</script>

<input
    type="number"
    min="0"
    max="100"
    step="0.25"
    value={intToDisplay(quality)}
    oninput={(e) => {
        const v = parseFloat((e.target as HTMLInputElement).value);
        quality = isNaN(v) ? null : Math.min(MAX, Math.max(0, displayToInt(v)));
    }}
/>
```

### Backend: Empty Perspective Validation (add to Create service)
```go
// Source: pattern from backend/internal/core/services/perspective_service.go
// Add after existing validations in Create():
func isEmptyPerspective(input portservices.CreatePerspectiveInput) bool {
    return input.Quality == nil &&
        input.Agreement == nil &&
        input.Importance == nil &&
        input.Confidence == nil &&
        input.Like == nil &&
        input.Description == nil &&
        len(input.CategorizedRatings) == 0
}

// In Create():
if isEmptyPerspective(input) {
    return nil, fmt.Errorf("%w: at least one field must be provided", domain.ErrInvalidInput)
}
```

---

## State of the Art

| Old Approach | Current Approach | Impact for Phase 4 |
|--------------|------------------|-------------------|
| `claim` column in perspectives table | Removed in migration 000007 | Claim is now a content_type, not a perspective field |
| Join tables for references | FK + int[] array on row | No bridge table to create; simpler migration |
| Separate test mock for each repo method | All mocks implement full interface | Adding new interface methods (e.g., `FindByContentID`) requires updating ALL test mocks |
| Custom cursor encoding | `gorm-cursor-paginator/v2` | Perspective List() already uses paginator — no change needed |

**Deprecated/outdated:**
- `claim` field on `perspectives` table: removed in migration 000007. Do NOT reference it. Claims are now `content` rows.
- `perspectives_unique_user_claims` constraint: removed in migration 000007. No uniqueness constraint exists on perspectives currently.

---

## Open Questions

1. **How does the perspectives query attach to ActivityTable without causing two simultaneous paginated fetches?**
   - What we know: ActivityTable fetches content with cursor pagination. A second perspectives query (all user's perspectives) would be a separate network request.
   - What's unclear: Should perspectives fetch all at once (no pagination) or match the content page?
   - Recommendation: Fetch all user perspectives without pagination (remove `first` limit or set high, e.g., 1000) since each user likely has few perspectives. This avoids complex cursor coordination between the two queries. The `perspectives` query is much smaller than content.

2. **Claim content storage — `response` JSONB vs new column?**
   - What we know: The `content` table has a `response jsonb` column used for YouTube API response caching. It could store `{"parentContentId": 123, "claimText": "..."}`.
   - What's unclear: Whether using `response` JSONB for claim-specific data is semantically correct vs adding a `parent_content_id` FK column.
   - Recommendation: For Phase 4, use the existing `response` column as JSONB storage `{parentContentId, text}`. This is additive with zero migration risk. Adding a dedicated FK column is a future refinement when the claim model stabilizes.

3. **Stepper granularity choice (Claude's discretion):**
   - Recommendation: Use **0.25** step (25 integer units). Rationale: finer than 0.50 (gives 40 steps over 10.00) without the awkwardness of 0.10 (100 steps is too many for click-based adjustment). 0.25 is standard for star rating equivalents.

4. **ON DELETE behavior for `primary_perspective_id` FK (Claude's discretion):**
   - Recommendation: **SET NULL**. If the referenced perspective is deleted, the reference becomes null rather than cascade-deleting the entire referencing perspective. This preserves the perspective's own data.

---

## Sources

### Primary (HIGH confidence)

- Codebase: `backend/migrations/000001_create_content.up.sql` — confirmed `content_type` is varchar, not PostgreSQL ENUM
- Codebase: `backend/migrations/000004_add_perspectives_users.up.sql` — existing perspectives table schema
- Codebase: `backend/migrations/000007_remove_claim_add_system_user.up.sql` — confirmed `claim` column removed
- Codebase: `backend/internal/core/domain/perspective.go` — domain model fields
- Codebase: `backend/internal/core/services/perspective_service.go` — service Create/Update logic, validation patterns
- Codebase: `backend/internal/adapters/repositories/postgres/gorm_models.go` — GORM PerspectiveModel
- Codebase: `backend/internal/adapters/repositories/postgres/gorm_mappers.go` — mapper patterns for arrays
- Codebase: `backend/internal/adapters/repositories/postgres/array_types.go` — Int64Array, StringArray, JSONBArray
- Codebase: `backend/schema.graphql` — full current schema, confirmed perspectives mutations exist
- Codebase: `frontend/src/lib/components/FormPopover.svelte` — popover/dialog pattern, mobile breakpoint
- Codebase: `frontend/src/lib/components/ActivityTable.svelte` — AG Grid column patterns, responsive visibility
- Codebase: `frontend/src/lib/queries/hooks/useAddVideo.ts` — mutation hook pattern with cache update
- Codebase: `.claude/docs/ADDING_AG_GRID_COLUMN.md` — confirmed two-place visibility requirement

### Secondary (MEDIUM confidence)
- Backend CLAUDE.md + ARCHITECTURE.md — hexagonal architecture layer rules
- Frontend CLAUDE.md — Svelte 5 runes-only requirement, AG Grid svelte5 constraint

### Tertiary (LOW confidence)
- None required — all research sourced from codebase directly.

---

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH — all libraries already in use, no new dependencies
- Architecture: HIGH — all patterns verified from existing code
- Migration approach: HIGH — additive FK + array pattern confirmed viable by existing array usage
- AG Grid popover interaction: MEDIUM — hover-to-edit requires non-standard workaround; click-to-edit recommended as simpler alternative
- Claim storage via JSONB response column: MEDIUM — pragmatic for phase 4 but semantically imperfect

**Research date:** 2026-02-18
**Valid until:** 2026-03-20 (stable stack, 30-day window)
