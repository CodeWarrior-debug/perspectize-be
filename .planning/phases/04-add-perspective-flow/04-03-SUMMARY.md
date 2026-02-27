---
phase: 04-add-perspective-flow
plan: 03
subsystem: backend-graphql-frontend-mutation-hooks
tags: [claims, graphql, tanstack-query, at-reference, perspective-popover, svelte5]
dependency_graph:
  requires:
    - backend/schema.graphql (createPerspective mutation from 04-01)
    - frontend/src/lib/components/PerspectivePopover.svelte (from 04-02)
    - frontend/src/lib/queries/hooks/useCreatePerspective.ts (from 04-01)
  provides:
    - backend/schema.graphql (createClaim mutation + CreateClaimInput)
    - backend/internal/core/services/content_service.go (CreateClaim method)
    - frontend/src/lib/queries/claims.ts
    - frontend/src/lib/queries/hooks/useCreateClaim.ts
    - frontend/src/lib/utils/references.ts
  affects:
    - frontend/src/lib/components/PerspectivePopover.svelte
    - backend/internal/adapters/graphql/resolvers/schema.resolvers.go
    - backend/internal/core/ports/services/content_service.go
tech_stack:
  added: []
  patterns:
    - "Claim as Content pattern: claims stored as content rows with content_type=claim and JSONB response containing parentContentId + raw text"
    - "@reference token pattern: tokens stored raw in DB, resolved at display time via resolveAtReference utility"
    - "Separate mutation pattern: claim creation fires independently from perspective submission (decision J)"
    - "Popover stays open after claim creation: claimText resets but Dialog remains visible"
key_files:
  created:
    - backend/internal/core/ports/services/content_service.go (CreateClaimInput struct + CreateClaim method added to interface)
    - frontend/src/lib/queries/claims.ts
    - frontend/src/lib/queries/hooks/useCreateClaim.ts
    - frontend/src/lib/utils/references.ts
    - frontend/tests/unit/utils-references.test.ts
  modified:
    - backend/schema.graphql (CreateClaimInput type + createClaim mutation)
    - backend/internal/adapters/graphql/generated/generated.go (regenerated)
    - backend/internal/adapters/graphql/model/models_gen.go (regenerated)
    - backend/internal/adapters/graphql/resolvers/schema.resolvers.go (CreateClaim resolver implemented)
    - backend/internal/core/services/content_service.go (CreateClaim method + encoding/json import)
    - frontend/src/lib/components/PerspectivePopover.svelte (claim section wired up)
decisions:
  - "CreateClaim validates parentContentID exists via repo.GetByID before creating (prevents orphan claims)"
  - "Claim name is set to the raw claim text (display = name field for claim rows in Activity table)"
  - "Response JSONB stores {parentContentId, text} for future display-time resolution without re-querying"
  - "isClaimPending tracked separately from isPending so perspective submit button is unaffected during claim creation"
  - "Create Claim button disabled when claimText.trim() is empty (client-side validation)"
metrics:
  duration: 5 min
  completed_date: "2026-02-27"
  tasks: 2
  files: 11
---

# Phase 4 Plan 03: Claim Creation Flow Summary

createClaim GraphQL mutation (backend end-to-end) + useCreateClaim TanStack hook + @reference token utilities + wired claim section in PerspectivePopover.

## Tasks Completed

| Task | Name | Commit | Key Files |
|------|------|--------|-----------|
| 1 | Backend createClaim mutation + @reference utilities | `ad8ea41` | schema.graphql, content_service.go, schema.resolvers.go, references.ts, utils-references.test.ts |
| 2 | Frontend claim creation hook + PerspectivePopover claim trigger | `7619c7e` | claims.ts, useCreateClaim.ts, PerspectivePopover.svelte |

## Verification Results

1. `cd backend && go build ./... && go test ./...` — zero errors, all tests pass
2. `cd frontend && pnpm run test:run` — 398 tests pass (26 test files, 13 new reference utility tests)
3. `cd frontend && pnpm run check` — 3 pre-existing type errors only (same as 04-01/04-02 — out of scope)
4. createClaim mutation visible in backend/schema.graphql
5. CreateClaim resolver implemented in schema.resolvers.go (delegates to ContentService.CreateClaim)
6. ContentService.CreateClaim validates text, userID, parentContentID; verifies parent exists; stores claim as content row with type=CLAIM and JSONB response {parentContentId, text}
7. resolveAtReference and hasAtReference exported from references.ts with 13 unit tests
8. PerspectivePopover "+ Add More..." section wired to useCreateClaim — textarea bound, helper text uses @this, Create Claim button triggers mutation

## Deviations from Plan

None — plan executed exactly as written.

## Must-Haves Verification

- [x] User can create a claim from within the PerspectivePopover via a dedicated trigger (button/section — Claude's discretion): "+ Add More..." button expands claim section with "Create Claim" button
- [x] createClaim mutation accepts text, userID, and parentContentID, creates a content row with content_type='claim'
- [x] Claim text containing @this or @here tokens is stored raw and resolved to parent content name at display time (via resolveAtReference utility)
- [x] Newly created claims appear in the Activity table as their own rows with type 'claim' (content list cache invalidated on success)
- [x] Backend stores parentContentID in the content response JSONB column as {parentContentId, text}

## Self-Check: PASSED

Files created verified:
- `frontend/src/lib/utils/references.ts` — EXISTS
- `frontend/tests/unit/utils-references.test.ts` — EXISTS (13 tests)
- `frontend/src/lib/queries/claims.ts` — EXISTS
- `frontend/src/lib/queries/hooks/useCreateClaim.ts` — EXISTS

Files modified verified:
- `backend/schema.graphql` — EXISTS (createClaim mutation + CreateClaimInput)
- `backend/internal/adapters/graphql/resolvers/schema.resolvers.go` — EXISTS (CreateClaim resolver implemented)
- `backend/internal/core/services/content_service.go` — EXISTS (CreateClaim method)
- `backend/internal/core/ports/services/content_service.go` — EXISTS (CreateClaimInput + CreateClaim interface method)
- `frontend/src/lib/components/PerspectivePopover.svelte` — EXISTS (claim section wired)

Commits verified:
- `ad8ea41` feat(04-03): add createClaim mutation, service, and @reference utility tests — EXISTS
- `7619c7e` feat(04-03): add useCreateClaim hook and wire claim creation in PerspectivePopover — EXISTS
