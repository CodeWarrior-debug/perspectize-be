# Project State

## Project Reference

See: .planning/PROJECT.md (updated 2026-02-04)

**Core value:** Users can easily submit their perspective on a YouTube video and browse others' perspectives in a way that keeps them in control.
**Current focus:** Phase 3 complete — Phase 3.1 (Dialog UX Polish) next

## Current Position

Phase: 3.1 of 5 (Design Token System) — IN PROGRESS
Plan: 1/3 complete
Status: Design token foundation established (27 color tokens, Charter font, dual typography)
Last activity: 2026-02-10 — Completed 03.1-01-PLAN.md

Progress: [███████░░░] 71%

## Performance Metrics

**Velocity:**
- Total plans completed: 12
- Average duration: 5.1 min
- Total execution time: 1.0 hours

**By Phase:**

| Phase | Plans | Total | Avg/Plan |
|-------|-------|-------|----------|
| 01-foundation | 5 | 36 min | 7 min |
| 02-data-layer-activity | 2 | 9 min | 4.5 min |
| 03-add-video-flow | 2 | 8 min | 4 min |
| 03.1-design-token-system | 1 | 3 min | 3 min |

**Recent Trend:**
- Last 5 plans: 6 min, 3 min, 3 min, 5 min, 3 min (avg: 4.0 min)
- Trend: Improving — Phase 3.1 started efficiently with focused CSS work

*Updated after each plan completion*

## Accumulated Context

### Decisions

Decisions are logged in PROJECT.md Key Decisions table.
Recent decisions affecting current work:

- [Roadmap]: 5 phases derived from 37 v1 requirements
- [Roadmap]: AG Grid validation in Phase 1 to derisk community wrapper early
- [Stack]: SvelteKit + TanStack Query + TanStack Form + AG Grid + shadcn-svelte + Tailwind CSS
- [01-01]: Downgraded Vite from 7.3.1 to 6.4.1 for Tailwind CSS v4 compatibility
- [01-01]: Upgraded Tailwind CSS from 4.0.0 to 4.1.18 to fix build errors
- [01-01]: Custom navy primary color: oklch(0.216 0.006 56.043) = #1a365d (CORRECTED in 03.1-01: hex #1a365d, not oklch)
- [03.1-01]: Hex color tokens (not oklch) — AI-generated oklch values unreliable, DESIGN_SPEC.md hex is source of truth
- [03.1-01]: Dual-font system: Geist (sans) for UI/headings, Charter (serif) for body/content text
- [03.1-01]: Charter font from charter-webfont npm package (official Matthew Carter distribution)
- [03.1-01]: disabled opacity via utility class (disabled:opacity-50) not @theme token
- [01-01]: Type-based organization over feature-based for flat structure
- [01-02]: Header height fixed at h-16 for consistency
- [01-02]: PageWrapper is opt-in for pages, not forced in layout
- [01-02]: Responsive padding scale: px-4 (mobile) → px-6 (md) → px-8 (lg)
- [01-03]: TanStack Query enabled: browser for SSG compatibility
- [01-03]: Toast notifications at top-right with 2s auto-dismiss
- [01-03]: Vitest with jsdom environment and SvelteKit mocks
- [01-05]: Browser resolve condition in vite.config.ts for Svelte 5 component testing
- [01-05]: Branch coverage threshold at 75% (vs 80% for others) due to Svelte compiler default parameter branches
- [01-05]: createRawSnippet pattern for testing Svelte 5 components with children props
- [02-01]: Setter-based session sync for Svelte 5 stores (not $effect) for testability
- [02-01]: Simple ListAll repository method (no pagination) for small datasets
- [02-01]: Namespaced session storage key: perspectize:selectedUserId
- [02-02]: TanStack Query function wrapper pattern for Svelte 5: createQuery(() => ({ ... }))
- [02-02]: Client-side AG Grid pagination (first: 100) for Phase 2 simplicity
- [02-02]: AG Grid Quick Filter for text search (no backend search yet)
- [02.1-02]: AG Grid onFirstDataRendered for deferred column visibility (avoids postConstruct bean init race)
- [02.1-02]: isGridInitialized flag gates onGridSizeChanged calls during AG Grid initialization
- [03-01]: shadcn components in shadcn/ directory (moved from default ui/ for consistency)
- [03-01]: URL constructor for YouTube validation (no regex, avoids catastrophic backtracking)
- [03-01]: Support 4 YouTube hosts: youtube.com, www.youtube.com, youtu.be, m.youtube.com
- [03-02]: TanStack Query createMutation with onSuccess/onError callbacks for dialog flow
- [03-02]: Error mapping via string matching (duplicate/already exists, invalid/not found, generic)
- [03-02]: Dialog stays open on error, closes only on success
- [03-02]: $effect for form reset when dialog opens (not on close)
- [Infra]: CLAUDE.md split into root + perspectize-go/CLAUDE.md + perspectize-fe/CLAUDE.md for package-level context loading
- [Infra]: Go module renamed from `github.com/yourorg/perspectize-go` to `github.com/CodeWarrior-debug/perspectize-be/perspectize-go` (30 files, all 78 tests pass)
- [Infra]: Docs delegated to docs/ directory: VERIFICATION.md, DOMAIN_GUIDE.md, GO_PATTERNS.md, GITHUB_PROJECTS.md, GSD_BRANCHING.md
- [Infra]: qmd .planning/ collection added with stable-vs-live convention
- [Infra]: All three CLAUDE.md files scored 95/100 (A) quality after optimization
- [Frontend]: Svelte 5 runes, SvelteKit routing, TanStack Query patterns documented in perspectize-fe/CLAUDE.md

### Roadmap Evolution

- Phase 02.1 inserted after Phase 2: Mobile Responsive Fixes (URGENT) — P1 issues: header overflow/clipping at 375px, pagination bar broken, table left-shift overflow
- Phase 03.1 inserted after Phase 3: Dialog UX Polish — Gray overlay too aggressive, modal translucent/hard to read, needs redesign with shadcn best practices

### Pending Todos

None yet.

### Known Bugs

- **Add Video dialog UX (P2):** Full-screen gray overlay too aggressive, modal content translucent/hard to read, poor visual hierarchy. **CRITICAL FIX APPLIED:** Transparent background bug fixed in 03.1-01 (bg-background now resolves to #ffffff). Remaining work: reduce overlay opacity, improve layout spacing (Phase 3.1 plans 02-03).

### Blockers/Concerns

None — Phase 3 complete, Phase 3.1 (Dialog UX) inserted before Phase 4.

## Session Log

### 2026-02-06 — CLAUDE.md Audit & Optimization

**Branch:** `feature/INI-37-plan-01-04-navigation-ag-grid`

**Work completed:**
1. **CLAUDE.md audit skill created** — custom-claude-improver skill with instruction counting, context budget analysis, and session-based compliance checking
2. **CLAUDE.md split** — Monolithic CLAUDE.md (683 lines, 372 instructions) split into root + perspectize-go/CLAUDE.md + perspectize-fe/CLAUDE.md
3. **Content delegated to docs/** — Created docs/VERIFICATION.md, docs/DOMAIN_GUIDE.md, docs/GO_PATTERNS.md, docs/GITHUB_PROJECTS.md, docs/GSD_BRANCHING.md
4. **Go module renamed** — `github.com/yourorg/perspectize-go` to `github.com/CodeWarrior-debug/perspectize-be/perspectize-go` (30 files, all 78 tests pass)
5. **qmd .planning/ collection** — Added .planning/ as separate qmd collection with stable-vs-live convention
6. **Quality scores** — All three CLAUDE.md files at 95/100 (A). Session compliance: root+backend 122/200, root+frontend 145/200
7. **Frontend patterns documented** — Svelte 5 runes, SvelteKit routing, TanStack Query patterns

**Commits:**
- `8fa9784` docs: apply CLAUDE.md audit fixes and split into package-level files
- `88b39b6` refactor: rename Go module path to match GitHub repository
- `f533df5` chore: add coverage/ to frontend .gitignore
- `6619f98` docs: add .planning/ qmd collection with stable vs live convention
- `21abf6a` docs: optimize all CLAUDE.md files to 95/100 quality score

### 2026-02-05 — Phase 1 Foundation (Plans 01-01 through 01-03)

### 2026-02-07 — Phase 03: Add Video Flow

**Branch:** `feature/INI-phase-03-add-video-flow`

**Work completed:**
1. **Plan 01: Foundation Components** — shadcn Dialog/Input/Label, YouTube URL validation (URL constructor), CREATE_CONTENT_FROM_YOUTUBE mutation definition
2. **Plan 02: AddVideoDialog + Header Wiring** — Complete dialog with TanStack Query mutation, error mapping (duplicate/invalid/generic), Header button wired to dialog
3. **Self-verification** — Chrome DevTools MCP: add video flow, validation errors, duplicate detection, mobile 375px responsiveness
4. **95 tests passing** — 12 test files, 26 new tests for Phase 3

**Commits:**
- `9bea0d4` feat(03-01): install shadcn Dialog, Input, Label components
- `9255b47` feat(03-01): add YouTube validation utility and mutation definition
- `8e330d3` docs(03-01): complete foundation components plan
- `713f00c` feat(03-02): create AddVideoDialog component with mutation and error handling
- `21faae0` feat(03-02): wire Header to AddVideoDialog and add tests

## Session Continuity

Last session: 2026-02-10
Stopped at: Completed 03.1-01-PLAN.md (Design Token System foundation)
Resume file: None

### 2026-02-07 — Plan 01-05: Test Coverage

**Branch:** `feature/INI-45-plan-01-05-test-coverage`

**Work completed:**
1. **Comprehensive test suite** — 42 tests across 6 test files with zero failures
2. **100% coverage** — Statement/line/function coverage at 100%, branch coverage at 75%
3. **Test helpers created** — Reusable Svelte 5 component testing utilities in tests/helpers/render.ts
4. **Coverage thresholds enforced** — 80% lines/functions/statements, 75% branches in vite.config.ts
5. **Testing patterns established** — Unit tests (tests/unit/), component tests (tests/components/), createRawSnippet for children props

**Commits:**
- `7c5f90b` test(01-05): add unit tests and test helpers
- `828617a` feat(01-05): add component tests and coverage thresholds

**Tests created:**
- tests/unit/utils.test.ts (9 tests - cn utility)
- tests/unit/queries-client.test.ts (3 tests - GraphQL client)
- tests/unit/queries-content.test.ts (11 tests - query definitions)
- tests/unit/shadcn-barrel.test.ts (4 tests - barrel exports)
- tests/components/Header.test.ts (9 tests - Header component)
- tests/components/PageWrapper.test.ts (6 tests - PageWrapper component)

### 2026-02-07 — Plan 02-01: Data Layer Foundation

**Branch:** `feature/INI-47-phase-02-data-layer-activity`

**Work completed:**
1. **Backend users list query** — GraphQL `users: [User!]!` query implemented through full hexagonal architecture stack
2. **Frontend query definitions** — LIST_USERS query created, LIST_CONTENT updated with sort/pagination params
3. **User selection store** — Session-persistent store with Svelte 5 runes using setter-based sync pattern
4. **Unit tests** — 13 new tests for queries and store (total: 55 frontend tests, 78+ backend tests)
5. **Svelte 5 runes pattern** — Established testable pattern for module-level stores without $effect

**Commits:**
- `4575b3d` feat(02-01): add users list query to backend
- `9e5de4f` feat(02-01): add frontend query definitions and user selection store
- `c003722` test(02-01): add unit tests for frontend data layer

**Key patterns established:**
- Setter-based session sync for Svelte 5 stores (not $effect) for testability
- GraphQL query parameter expansion pattern (sort, pagination, includeTotalCount)
- Repository ListAll pattern for fetching all records without pagination

**Files created:**
- perspectize-fe/src/lib/queries/users.ts (LIST_USERS query)
- perspectize-fe/src/lib/stores/userSelection.svelte.ts (session-persistent store)
- perspectize-fe/tests/unit/queries-users.test.ts
- perspectize-fe/tests/unit/stores-userSelection.test.ts

**Duration:** 6 min

### 2026-02-07 — Plan 02-02: Activity Page UI

**Branch:** `feature/INI-47-phase-02-data-layer-activity`

**Work completed:**
1. **ActivityTable component** — AG Grid wrapper with 5 columns (Title, Type, Duration, Date Added, Last Updated), Quick Filter search, pagination (10/25/50)
2. **UserSelector component** — User dropdown with TanStack Query fetching users list, session storage persistence
3. **Activity page** — TanStack Query integration, search input, ActivityTable component with real content data
4. **Header update** — UserSelector added before Add Video button
5. **Component tests** — ActivityTable (4 tests), UserSelector (2 tests), Header tests updated (9 tests)

**Commits:**
- `d4ff7a2` feat(02-02): create ActivityTable and UserSelector components
- `45c3af2` feat(02-02): wire Activity page and update Header, add component tests

**Key patterns established:**
- TanStack Query function wrapper for Svelte 5: `createQuery(() => ({ ... }))`
- TanStack Query returns reactive object (not store) - access directly, not with $
- AG Grid reactive loading state via $effect
- Component mocking strategy for tests with QueryClient dependencies

**Files created:**
- perspectize-fe/src/lib/components/ActivityTable.svelte (AG Grid wrapper)
- perspectize-fe/src/lib/components/UserSelector.svelte (User dropdown)
- perspectize-fe/tests/components/ActivityTable.test.ts
- perspectize-fe/tests/components/UserSelector.test.ts
- perspectize-fe/tests/helpers/TestWrapper.svelte

**Duration:** 3 min
