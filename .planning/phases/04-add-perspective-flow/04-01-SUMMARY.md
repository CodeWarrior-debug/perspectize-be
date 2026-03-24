---
phase: 04-add-perspective-flow
plan: 01
subsystem: backend-domain-graphql-frontend-queries
tags: [perspective, migration, graphql, tanstack-query, mutation-hooks]
dependency_graph:
  requires: []
  provides:
    - backend/migrations/000013_add_perspective_refs_claims.up.sql
    - frontend/src/lib/queries/perspectives.ts
    - frontend/src/lib/queries/hooks/useCreatePerspective.ts
    - frontend/src/lib/queries/hooks/useUpdatePerspective.ts
  affects:
    - backend/schema.graphql
    - backend/internal/core/domain/perspective.go
    - backend/internal/adapters/repositories/postgres/gorm_models.go
    - frontend/src/lib/queries/keys.ts
tech_stack:
  added: []
  patterns:
    - "Additive DB migration pattern: new columns with defaults, GIN index for array reverse lookups"
    - "Fragment reuse pattern in GraphQL operations (PERSPECTIVE_FIELDS fragment)"
    - "At-least-one-field validation in service Create method"
key_files:
  created:
    - backend/migrations/000013_add_perspective_refs_claims.up.sql
    - backend/migrations/000013_add_perspective_refs_claims.down.sql
    - frontend/src/lib/queries/perspectives.ts
    - frontend/src/lib/queries/hooks/useCreatePerspective.ts
    - frontend/src/lib/queries/hooks/useUpdatePerspective.ts
    - frontend/tests/unit/queries-perspectives.test.ts
    - frontend/tests/unit/hooks-useCreatePerspective.test.ts
    - frontend/tests/unit/hooks-useUpdatePerspective.test.ts
  modified:
    - backend/internal/core/domain/perspective.go
    - backend/internal/core/domain/content.go
    - backend/internal/adapters/repositories/postgres/gorm_models.go
    - backend/internal/adapters/repositories/postgres/gorm_mappers.go
    - backend/schema.graphql
    - backend/internal/adapters/graphql/model/models_gen.go
    - backend/internal/adapters/graphql/generated/generated.go
    - backend/internal/adapters/graphql/resolvers/schema.resolvers.go
    - backend/internal/adapters/graphql/resolvers/helpers.go
    - backend/internal/core/ports/services/perspective_service.go
    - backend/internal/core/services/perspective_service.go
    - backend/test/services/perspective_service_test.go
    - frontend/src/lib/queries/keys.ts
decisions:
  - "Migration numbered 000013 (not 000012 as planned) — 000012 was already taken by backfill_canonical_urls migration from Phase 17"
  - "GraphQL fragment PERSPECTIVE_FIELDS defined in perspectives.ts to avoid repetition across 3 operations"
  - "useCreatePerspective also invalidates content.lists() on success (perspective status affects content rows in ActivityTable)"
  - "At-least-one-field validation added to Create service path only — Update allows field-only updates (no additional empty check)"
metrics:
  duration: 8 min
  completed_date: "2026-02-27"
  tasks: 2
  files: 21
---

# Phase 4 Plan 01: Backend Schema + Domain + Frontend Queries Summary

Established the full data pipeline for the Add Perspective flow: DB migration adds 4 new columns (primary_perspective_id FK, related_perspective_ids int[] with GIN index, custom_fields JSONB, review TEXT) to the perspectives table; domain/GORM/GraphQL models extended with new fields; CLAIM added to ContentType enum; two TanStack mutation hooks with toast feedback created; 45 new frontend tests pass.

## Tasks Completed

| Task | Name | Commit | Key Files |
|------|------|--------|-----------|
| 1 | Backend schema migration + domain/GORM/GraphQL extensions | `224592e` | 000013 migration, perspective.go, gorm_models.go, schema.graphql, models_gen.go |
| 2 | Frontend query definitions, mutation hooks, and tests | `94e07ab` | perspectives.ts, useCreatePerspective.ts, useUpdatePerspective.ts, 3 test files |

## Verification Results

1. `cd backend && go build ./... && go test ./...` — zero errors, all tests pass (added 1 new test)
2. `cd frontend && pnpm run test:run` — 361 tests pass including 45 new ones
3. `cd frontend && pnpm run check` — 3 pre-existing type errors (AGGridTest.svelte, +page.svelte, FormPopover.test.ts) — out of scope, not introduced by this plan
4. Migration file exists at `backend/migrations/000013_add_perspective_refs_claims.up.sql`
5. New fields visible in `backend/internal/adapters/graphql/model/models_gen.go` after `make graphql-gen`

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 1 - Bug] Migration number conflict: used 000013 instead of 000012**
- **Found during:** Task 1
- **Issue:** Plan specified migration `000012_add_perspective_refs_claims` but `000012_backfill_canonical_urls` already existed from Phase 17
- **Fix:** Created `000013_add_perspective_refs_claims` instead
- **Files modified:** `backend/migrations/000013_*.sql` (created)
- **Commit:** `224592e`

**2. [Rule 1 - Bug] Fixed 2 existing service tests broken by at-least-one-field validation**
- **Found during:** Task 1 test run
- **Issue:** `TestPerspectiveCreate_Success` and `TestPerspectiveCreate_RepositoryError` created perspectives with zero content fields, failing the new validation
- **Fix:** Added `like: &like` to both test inputs; added `TestPerspectiveCreate_NoFieldsProvided` test to verify validation works
- **Files modified:** `backend/test/services/perspective_service_test.go`
- **Commit:** `224592e`

### Deferred Items (Pre-existing, Out of Scope)

The following type errors existed before this plan and are not introduced by these changes:
- `src/lib/components/AGGridTest.svelte:58` — rowData initial value capture warning
- `src/routes/+page.svelte:27` — Type 'string' not assignable to 'never'
- `tests/components/FormPopover.test.ts:37,38` — Snippet type mismatch in test

## Must-Haves Verification

- [x] Backend accepts createPerspective mutation with all optional fields (quality, agreement, importance, confidence, like, review text) and returns created perspective
- [x] Backend accepts updatePerspective mutation with partial fields and returns updated perspective
- [x] Backend validates that at least one field is non-empty on create (any rating OR any text)
- [x] Frontend useCreatePerspective hook calls GraphQL mutation, shows success/error toasts, and invalidates perspective cache
- [x] Frontend useUpdatePerspective hook calls GraphQL mutation with partial updates, shows success/error toasts
- [x] Database migration adds primary_perspective_id FK, related_perspective_ids int array with GIN index, and custom_fields JSONB column to perspectives table
- [x] Content type 'claim' is supported in domain constants and GraphQL enum

## Self-Check: PASSED

Files created/modified verified:
- `backend/migrations/000013_add_perspective_refs_claims.up.sql` — EXISTS
- `backend/migrations/000013_add_perspective_refs_claims.down.sql` — EXISTS
- `frontend/src/lib/queries/perspectives.ts` — EXISTS
- `frontend/src/lib/queries/hooks/useCreatePerspective.ts` — EXISTS
- `frontend/src/lib/queries/hooks/useUpdatePerspective.ts` — EXISTS
- `frontend/tests/unit/queries-perspectives.test.ts` — EXISTS (22 tests)
- `frontend/tests/unit/hooks-useCreatePerspective.test.ts` — EXISTS (11 tests)
- `frontend/tests/unit/hooks-useUpdatePerspective.test.ts` — EXISTS (12 tests)

Commits verified:
- `224592e` feat(04-01): extend perspective schema with refs, claim type, and custom fields — EXISTS
- `94e07ab` feat(04-01): add frontend perspective query definitions, mutation hooks, and tests — EXISTS
