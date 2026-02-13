# Roadmap: Perspectize v1.0 Frontend MVP

## Overview

This roadmap delivers a functional SvelteKit frontend for Perspectize, enabling users to discover YouTube videos, add new videos via URL paste, and submit perspectives with ratings. The journey starts with project scaffolding and design system setup, progresses through the Activity page with AG Grid, then adds video and perspective creation flows, culminating in testing and deployment.

## Phases

**Phase Numbering:**
- Integer phases (1, 2, 3): Planned milestone work
- Decimal phases (2.1, 2.2): Urgent insertions (marked with INSERTED)

Decimal phases appear between their surrounding integers in numeric order.

- [x] **Phase 1: Foundation** - SvelteKit project setup with design system, tooling, and navigation skeleton
- [x] **Phase 2: Data Layer + Activity** - TanStack Query integration, AG Grid table, user selector
- [x] **Phase 2.1: Mobile Responsive Fixes** - Fix header overflow, pagination bar, and table layout at 375px (INSERTED)
- [x] **Phase 3: Add Video Flow** - YouTube URL paste, auto-fetch metadata, toast notifications
- [x] **Phase 3.1: Design Token System** - Implement all 27 Figma color variables, Geist + Charter typography, rating colors in code (INSERTED — rescoped)
- [ ] **Phase 3.2: Activity Page Beta Quality** - Rebuild Activity page to beta quality with server-side ops, new columns, popover dialog, data provenance (INSERTED)
- [ ] **Phase 3.3: Repository Rename & Folder Restructure** - Rename repo to perspectize, perspectize-go → backend, perspectize-fe → fe, update imports and Sevalla (INSERTED)
- [ ] **Phase 4: Add Perspective Flow** - TanStack Form with ratings, Like, Review, validation
- [ ] **Phase 5: Testing + Deployment** - Test coverage, CI/CD, hosting, CORS configuration

## Phase Details

### Phase 1: Foundation
**Goal**: Establish project scaffolding with all core libraries, mobile-first design system, and navigation working end-to-end
**Depends on**: Nothing (first phase)
**Requirements**: SETUP-01, SETUP-02, SETUP-03, SETUP-04, SETUP-05, SETUP-06, SETUP-07, SETUP-08, SETUP-09, NAV-01, NAV-02, NAV-03, NAV-04, NAV-05, API-01, TEST-01, TEST-02
**Success Criteria** (what must be TRUE):
  1. Activity page has "Add Video" button in header (modal placeholder — full modal in Phase 3)
  2. Application loads with custom navy theme (#1a365d) and Inter font applied
  3. Toast notifications appear in top-right and auto-dismiss after 2 seconds
  4. AG Grid renders a test table (validation that wrapper works with Svelte 5)
  5. Layout works on iPhone SE (375px) and scales up to desktop
  6. Folder structure documented with example files in each folder
  7. Test coverage >80% on all foundation source files with enforced thresholds
**Plans**: 5 plans in 3 waves

Plans:
- [x] 01-01-PLAN.md - SvelteKit + Tailwind + shadcn + Inter font + folder structure docs
- [x] 01-02-PLAN.md - Mobile-first layout system (Header, PageWrapper, breakpoints)
- [x] 01-03-PLAN.md - TanStack Query + GraphQL client + Toast + Vitest
- [x] 01-04-PLAN.md - Navigation + AG Grid validation (checkpoint)
- [x] 01-05-PLAN.md - Comprehensive test coverage (>80%) with testing conventions

### Phase 2: Data Layer + Activity
**Goal**: Users can view recently updated content in an AG Grid table (Activity page, default landing)
**Depends on**: Phase 1
**Requirements**: ACT-01, ACT-02, ACT-03, ACT-04, ACT-05, ACT-06, USER-01, USER-02, API-03, TEST-03
**Success Criteria** (what must be TRUE):
  1. User can view Activity page showing most recently updated content (default sort by updatedAt desc)
  2. User can sort the table by any column (ascending/descending)
  3. User can filter the table by text search (searches title, description)
  4. User can paginate with page size selector (10/25/50) using cursor-based navigation
  5. User can select from existing users via dropdown, selection persists across navigation
**Plans**: 2 plans in 2 waves

Plans:
- [x] 02-01-PLAN.md — Backend users query + frontend query definitions + user selection store + tests
- [x] 02-02-PLAN.md — Activity page with AG Grid table, UserSelector in header, visual checkpoint

### Phase 2.1: Mobile Responsive Fixes (INSERTED)
**Goal**: Fix P1/P2 mobile responsive issues at 375px so the Activity page is fully usable on iPhone SE
**Depends on**: Phase 2
**Success Criteria** (what must be TRUE):
  1. Header fits at 375px — logo, user selector, and Add Video button all visible without clipping
  2. AG Grid pagination bar is readable and usable at 375px
  3. Table content aligns within viewport bounds (no left-shift overflow)
**Plans**: 2 plans in 2 waves

Plans:
- [x] 02.1-01-PLAN.md — Header responsive fix (min-w-0, truncate, shrink-0) + AG Grid pagination CSS override
- [x] 02.1-02-PLAN.md — Table responsive column hiding + visual checkpoint at 375px

### Phase 3: Add Video Flow
**Goal**: Users can add YouTube videos by pasting a URL, with automatic metadata fetch and feedback
**Depends on**: Phase 2
**Requirements**: VIDEO-01, VIDEO-02, VIDEO-03, VIDEO-04, VIDEO-05, TEST-04
**Success Criteria** (what must be TRUE):
  1. User can paste a YouTube URL and submit to add a video
  2. Backend auto-fetches video metadata (title, description, thumbnail, duration) on creation
  3. User sees success toast (top-right, 2s auto-dismiss) after video creation
  4. User sees error toast if URL is invalid or fetch fails
  5. User is warned via toast if video already exists (duplicate detection)
**Plans**: 2 plans in 2 waves

Plans:
- [x] 03-01-PLAN.md — YouTube URL validation utility, mutation definition, shadcn Dialog/Input/Label setup
- [x] 03-02-PLAN.md — AddVideoDialog component with mutation, error handling, Header wiring, visual checkpoint

### Phase 3.1: Design Token System (INSERTED — rescoped from Dialog UX Polish)
**Goal**: Implement all Figma design tokens in code — 27 color variables, Geist + Charter dual-font typography, rating colors — so that all shadcn components render correctly and Phase 3.2+ can build on a complete token foundation
**Depends on**: Phase 3
**Success Criteria** (what must be TRUE):
  1. All 22 Theme/Light color tokens from DESIGN_SPEC.md are defined in app.css and usable via Tailwind classes (bg-background, text-foreground, etc.)
  2. All 4 rating color tokens are defined (rating-positive, rating-neutral, rating-negative, rating-undecided)
  3. Charter font is loaded and available via font-serif utility class
  4. Geist + Charter dual-font system works (Geist for UI, Charter for body/content)
  5. Existing components (Dialog, Header, shadcn primitives) render correctly with the new tokens — no visual regressions
**Plans**: 2 plans in 2 waves

Plans:
- [x] 03.1-01-PLAN.md — Complete color token system (27 tokens) + Charter font setup in app.css
- [x] 03.1-02-PLAN.md — AG Grid branded theming + visual verification checkpoint

### Phase 3.2: Activity Page Beta Quality (INSERTED)
**Goal**: Get every element on the Activity page — table, header, layout, dialog — to beta quality. Rebuild on top of Phase 3.1 design tokens. Server-side sorting/filtering/pagination, new YouTube columns, popover dialog, data provenance infrastructure. Source-data columns only.
**Depends on**: Phase 3.1
**Success Criteria** (what must be TRUE):
  1. Server-side sorting — column header clicks send sort params to GraphQL API, backend returns pre-sorted data
  2. Per-column filters — no global search bar; each column header has its own filter control
  3. Server-side pagination — one page at a time from backend, total count displayed, page size selector (10/25/50)
  4. 6 default visible columns: Item (title+thumbnail), Type (YouTube icon), Length, Views, Likes, Date (YouTube publish date)
  5. Column chooser with hidden columns: Channel, Date Added, Date Updated, Tags, Description
  6. New YouTube fields exposed through GraphQL (view count, like count, channel, publish date, tags)
  7. Compact rows (~40-48px height) matching Figma reference (node 3:409)
  8. Sticky page header + sticky table header, table body scrolls independently
  9. Add Video dialog: popover-near-button pattern (no overlay), page stays interactive while open
  10. Data provenance visual infrastructure: columns grouped by source, tooltip on header hover, visual tier indicators
  11. Empty state: "No items yet - add the first one!" in table body area
**Plans**: 4 plans in 3 waves

Plans:
- [ ] 03.2-01-PLAN.md — Backend: expose YouTube fields (channelTitle, publishedAt, tags, description) + extend sort/filter
- [ ] 03.2-02-PLAN.md — Frontend: popover dialog redesign (replace modal with non-modal popover)
- [ ] 03.2-03-PLAN.md — Frontend: ActivityTable rewrite (server-side pagination, new columns, compact rows, sticky headers, provenance)
- [ ] 03.2-04-PLAN.md — Integration polish, test coverage, visual verification checkpoint

### Phase 3.3: Repository Rename & Folder Restructure (INSERTED)
**Goal**: Rename repository from perspectize-be to perspectize, restructure folders (perspectize-go → backend, perspectize-fe → fe), update all Go imports, fix CI/CD and Sevalla deployment pointers
**Depends on**: Phase 3.2
**Success Criteria** (what must be TRUE):
  1. GitHub repository renamed to `perspectize`
  2. `perspectize-go/` renamed to `backend/` with all Go import paths updated and tests passing
  3. `perspectize-fe/` renamed to `fe/` with all config paths updated and build passing
  4. Sevalla deployment updated to point to new directory structure
  5. All CLAUDE.md files, docs, and planning references updated to reflect new paths
  6. CI/CD (if any) updated for new paths
**Plans**: 3 plans in 3 waves

Plans:
- [ ] 03.3-01-PLAN.md — GitHub repo rename (checkpoint) + Go module path refactor (imports, gqlgen.yml, tests)
- [ ] 03.3-02-PLAN.md — Folder rename (git mv) + bulk path updates across CI/CD, CLAUDE.md, docs, planning
- [ ] 03.3-03-PLAN.md — Deployment config update (checkpoint) + push to GitHub + CI/CD verification

### Phase 4: Add Perspective Flow
**Goal**: Users can create perspectives on videos with ratings, Like text, and Review text
**Depends on**: Phase 3
**Requirements**: PERSP-01, PERSP-02, PERSP-03, PERSP-04, PERSP-05, PERSP-06, PERSP-07, PERSP-08, PERSP-09, USER-03, TEST-05
**Success Criteria** (what must be TRUE):
  1. User can select a video to add a perspective on
  2. User can set Quality, Agreement, Importance, and Confidence ratings via number inputs with progress bar visualization
  3. User can enter Like text and Review text (freeform)
  4. User sees validation error toasts before submission if form is invalid
  5. User sees success toast after perspective is created, attributed to selected user
**Plans**: 2 plans in 2 waves

Plans:
- [ ] 04-01-PLAN.md — Mutation definition, shadcn Progress/Textarea, RatingInput and VideoSelector components with tests
- [ ] 04-02-PLAN.md — AddPerspectiveDialog with TanStack Form, Header wiring, validation, visual checkpoint

### Phase 5: Testing + Deployment
**Goal**: Application is tested, deployed, and accessible via public URL with proper CORS
**Depends on**: Phase 4
**Requirements**: TEST-06, DEPLOY-01, DEPLOY-02, DEPLOY-03, API-02
**Success Criteria** (what must be TRUE):
  1. Final coverage verification shows >=80% of lines covered
  2. Frontend is deployed and accessible via public URL
  3. CI/CD pipeline runs tests and deploys automatically on push
  4. CORS is configured on backend to allow frontend origin
**Plans**: 3 plans in 2 waves

Plans:
- [x] 05-01-PLAN.md — Coverage verification — SKIPPED (thresholds already met: 87.6% stmts, 90.1% lines)
- [ ] 05-02-PLAN.md — DigitalOcean App Platform static site deployment (cleanup GitHub Pages artifacts, deploy frontend)
- [ ] 05-03-PLAN.md — CORS configuration with rs/cors and DigitalOcean App Platform frontend origin

---

## Post-MVP: Concerns Remediation (Phases 6–10)

Phases 6–10 address the 77 issues cataloged in `.planning/codebase/CONCERNS.md`. Ordered by dependency: fix errors first, then architecture, then schema, then security (which depends on clean architecture), then frontend. Each phase is a living checklist — items can be picked off incrementally.

- [ ] **Phase 6: Error Handling & Data Integrity** - Fix silent failures, error leakage, and config validation
- [ ] **Phase 7: Backend Architecture** - Hexagonal cleanup, dependency injection, server infrastructure
- [ ] **Phase 8: API & Schema Quality** - Fix pagination, GraphQL types, race conditions, nested resolvers
- [ ] **Phase 9: Security Hardening** - Authentication, rate limiting, query complexity, headers, HTTPS
- [ ] **Phase 10: Frontend Quality & Test Coverage** - XSS fix, codegen, error boundaries, cleanup, test gaps

### Phase 6: Error Handling & Data Integrity
**Goal**: Eliminate all silent failures so errors are visible, logged, and surfaced correctly to clients
**Depends on**: Phase 5 (CI/CD catches regressions)
**Source**: CONCERNS.md C-06, C-07, C-08, H-13, H-16, H-19, H-20, H-21, M-03, M-27
**Success Criteria** (what must be TRUE):
  1. All `json.Unmarshal` calls check and handle errors (C-06, C-08)
  2. All `strconv`/`time.Parse` calls check and handle errors (C-07, C-08)
  3. GraphQL error responses never expose database schema or internal details (H-13)
  4. Not-found handling is consistent across all resolvers — standardized pattern (H-16, M-07)
  5. `.env` load warns if file missing in dev; YouTube API key validated at startup (H-19, H-20)
  6. `WriteString` return value checked in IntID marshal (H-21)
  7. `CreateFromYouTube` returns existing item on duplicate instead of error (M-03)
  8. `formatDate` handles invalid input gracefully (M-27)
**Plans**: TBD

**Concern checklist:**
- [ ] C-06: Silent JSON unmarshal in perspective repository
- [ ] C-07: Silent duration parse in YouTube client
- [ ] C-08: Five silent parse failures in `domainToModel` helpers
- [ ] H-13: Sensitive data leaked in GraphQL errors (use generic client errors, log full errors server-side)
- [ ] H-16: Inconsistent not-found error handling across resolvers
- [ ] H-19: `.env` load failure silently ignored
- [ ] H-20: Empty YouTube API key not validated at startup
- [ ] H-21: `WriteString` return value ignored in IntID
- [ ] M-03: `CreateFromYouTube` returns error instead of idempotent result
- [ ] M-07: Inconsistent not-found (duplicate of H-16)
- [ ] M-27: `formatDate` silently produces "Invalid Date"

### Phase 7: Backend Architecture
**Goal**: Clean up hexagonal architecture violations, add proper dependency injection, and harden server infrastructure
**Depends on**: Phase 6 (error handling patterns established first)
**Source**: CONCERNS.md H-01, H-02, H-09, M-01, M-02, M-05, M-06, M-09, M-10, M-12, M-17
**Success Criteria** (what must be TRUE):
  1. No adapter-to-adapter imports — resolvers use service ports only (H-01, H-02)
  2. Service port interfaces defined; resolver depends on interfaces, not concrete types (H-02)
  3. Config path loaded from env var with sensible default (H-09)
  4. Single PostgreSQL driver (`pgx`) — `lib/pq` removed (M-01)
  5. DB pool settings configurable via env vars (M-02)
  6. YouTube `extractVideoID` injected via constructor, not function param (M-05)
  7. Request logging middleware installed (chi router or similar) (M-06)
  8. Graceful shutdown with SIGTERM handler (M-09)
  9. `/health` and `/ready` endpoints exist (M-10)
  10. DB credentials sanitized before logging (M-12)
  11. `DATABASE_URL` format validated at startup (M-17)
**Plans**: TBD

**Concern checklist:**
- [ ] H-01: Adapter-to-adapter coupling (resolver imports YouTube adapter directly)
- [ ] H-02: Resolver depends on concrete service types (no port interfaces)
- [ ] H-09: Hardcoded config path
- [ ] M-01: Dual PostgreSQL driver dependencies (`lib/pq` + `pgx`)
- [ ] M-02: Hardcoded database connection pool settings
- [ ] M-05: Function parameter instead of dependency injection
- [ ] M-06: No request logging middleware
- [ ] M-09: No graceful shutdown handler
- [ ] M-10: No health check endpoint
- [ ] M-12: DB credentials in logs on failure
- [ ] M-17: No `DATABASE_URL` format validation

### Phase 8: API & Schema Quality
**Goal**: Fix GraphQL schema types, pagination bugs, race conditions, and missing resolvers
**Depends on**: Phase 7 (clean architecture enables proper resolver changes)
**Source**: CONCERNS.md C-02, H-03, H-04, H-05, H-06, H-07, H-08, M-04, M-08, M-11, M-13, M-16
**Success Criteria** (what must be TRUE):
  1. Cursor pagination works correctly for all sort columns, not just ID (C-02)
  2. `ListAll` users has pagination with configurable limit (H-03)
  3. Timestamps use `DateTime` scalar with proper serialization (H-04)
  4. `contentType` uses `ContentType` enum, not `String` (H-05)
  5. Uniqueness enforced via DB constraints, not app-level TOCTOU checks (H-06, H-07)
  6. YouTube API response stored as structured fields, not raw JSON blob (H-08, M-13)
  7. `deletePerspective` uses `IntID` consistently (M-04)
  8. Perspective `user` and `content` nested fields resolve correctly (M-08)
  9. Input length validation on description, labels, categorizedRatings (M-11)
  10. Update checks `RowsAffected` for optimistic concurrency (M-16)
**Plans**: TBD

**Concern checklist:**
- [ ] C-02: Cursor pagination broken for non-ID sorts (compound keyset required)
- [ ] H-03: `ListAll()` users has no pagination (unbounded query)
- [ ] H-04: Timestamps as `String!` instead of `DateTime` scalar
- [ ] H-05: `contentType` uses `String!` instead of `ContentType` enum
- [ ] H-06: Race condition on perspective claim uniqueness check (TOCTOU)
- [ ] H-07: Race condition on user uniqueness check (TOCTOU)
- [ ] H-08: YouTube API response stored verbatim (~5KB bloat per item)
- [ ] M-04: `deletePerspective` uses `ID` scalar instead of `IntID`
- [ ] M-08: Missing nested field resolvers (perspective->user, perspective->content)
- [ ] M-11: Missing input length validation
- [ ] M-13: Unbounded JSON field (duplicate of H-08)
- [ ] M-16: Update does not check `RowsAffected`

### Phase 9: Security Hardening
**Goal**: Add authentication, authorization, rate limiting, and security headers to make the app safe for multi-user deployment
**Depends on**: Phase 8 (clean schema + architecture required before layering auth)
**Source**: CONCERNS.md C-01, C-04, C-05, C-09, C-10, H-10, H-11, H-12, H-14, H-15, H-25, M-14, M-15, M-28
**Success Criteria** (what must be TRUE):
  1. Authentication middleware validates JWT/session on all mutations (C-01)
  2. Authorization checks on all mutations — users can only modify their own data (C-01)
  3. GraphQL query complexity limit enforced (C-04)
  4. CORS restricted to explicit frontend origin (C-05 -- may already be done in Phase 5)
  5. GraphQL playground disabled in production (C-09)
  6. Introspection disabled in production (C-10)
  7. User email addresses only visible to authenticated user for their own account (H-10)
  8. Rate limiting middleware installed (H-11)
  9. YouTube API key never appears in logs or error responses (H-12)
  10. HTTP server has read/write/idle timeouts configured (H-15)
  11. TLS/HTTPS via reverse proxy or `ListenAndServeTLS` (H-14)
  12. Security headers set: `X-Content-Type-Options`, `X-Frame-Options`, HSTS (M-14)
  13. CSRF protection middleware installed (M-15)
  14. Content Security Policy header on frontend (H-25)
  15. Secrets managed via vault/rotation mechanism (M-28)
**Plans**: TBD

**Concern checklist:**
- [ ] C-01: No authentication or authorization (CRITICAL)
- [ ] C-04: No GraphQL query complexity limiting (DoS vector)
- [ ] C-05: Wildcard CORS configuration (may be resolved by Phase 5, plan 05-03)
- [ ] C-09: GraphQL playground exposed unconditionally
- [ ] C-10: GraphQL introspection enabled without restriction
- [ ] H-10: User email addresses exposed in public query
- [ ] H-11: No rate limiting
- [ ] H-12: YouTube API key exposure risk in logs/errors
- [ ] H-14: No HTTPS/TLS
- [ ] H-15: No HTTP server timeouts (Slowloris DoS)
- [ ] H-25: No Content Security Policy
- [ ] M-14: Missing security headers
- [ ] M-15: No CSRF protection
- [ ] M-28: No secret rotation or vault integration

### Phase 10: Frontend Quality & Test Coverage
**Goal**: Fix frontend vulnerabilities, add codegen, error boundaries, and close all test coverage gaps
**Depends on**: Phase 8 (schema fixes enable codegen; nested resolvers enable frontend cleanup)
**Source**: CONCERNS.md C-03, H-17, H-18, H-22, H-23, H-24, M-18-M-26, T-01-T-06, L-*
**Success Criteria** (what must be TRUE):
  1. AG Grid cellRenderer uses safe DOM APIs, no raw innerHTML interpolation (C-03)
  2. `+error.svelte` error boundary exists with retry UI (H-17, M-23)
  3. `hooks.client.ts` and `hooks.server.ts` catch unhandled errors (H-18)
  4. `prerender` set to false or properly leveraged (H-22)
  5. `graphql-codegen` generates TypeScript types from schema (H-23, M-18)
  6. GraphQL client has error interceptor, timeout, and auth header support (H-24)
  7. Search input debounced (M-22)
  8. Dead code removed: `AGGridTest.svelte`, unused stores, unused type guards (M-20, M-21, M-24)
  9. Retry config only retries 5xx/network errors (M-26)
  10. `PerspectiveService.Update()` has test coverage (T-01)
  11. Resolver tests exist for User/Perspective queries (T-02)
  12. `helpers.go` conversion tested with malformed input (T-03)
  13. Repository-layer tests exist (T-04)
  14. YouTube API client tested (T-05)
  15. `IntID` scalar tested (T-06)
**Plans**: TBD

**Concern checklist:**
- [ ] C-03: XSS vulnerability in AG Grid cellRenderer
- [ ] H-17: Missing `+error.svelte` error boundary
- [ ] H-18: Missing `hooks.client.ts` / `hooks.server.ts`
- [ ] H-22: `prerender = true` without SSR (architectural mismatch)
- [ ] H-23: No TypeScript types generated from GraphQL schema
- [ ] H-24: GraphQL client missing error/timeout infrastructure
- [ ] M-18: Duplicated type definitions across components
- [ ] M-19: No server-side pagination integration (hard-coded 100)
- [ ] M-20: `selectedUserId` store not consumed
- [ ] M-21: Unused type guards
- [ ] M-22: Search input not debounced
- [ ] M-23: No error recovery UI (no retry button)
- [ ] M-24: Dead code (`AGGridTest.svelte`)
- [ ] M-25: HTTP fallback for GraphQL endpoint
- [ ] M-26: Retry configuration retries all errors (should only retry 5xx)
- [ ] T-01: `PerspectiveService.Update()` -- zero tests
- [ ] T-02: No resolver tests
- [ ] T-03: No `helpers.go` conversion tests
- [ ] T-04: No repository-layer tests
- [ ] T-05: No YouTube API client tests
- [ ] T-06: No `IntID` scalar tests
- [ ] L-*: Low priority cleanup (see CONCERNS.md L-01 through L-22)

## Progress

**Execution Order:**
Phases execute in numeric order: 1 -> 2 -> 2.1 -> 3 -> 3.1 -> 3.2 -> 3.3 -> 4 -> 5 -> 6 -> 7 -> 8 -> 9 -> 10

| Phase | Plans Complete | Status | Completed |
|-------|----------------|--------|-----------|
| 1. Foundation | 5/5 | Complete | 2026-02-07 |
| 2. Data Layer + Activity | 2/2 | Complete | 2026-02-07 |
| 2.1 Mobile Responsive Fixes | 2/2 | Complete | 2026-02-07 |
| 3. Add Video Flow | 2/2 | Complete | 2026-02-07 |
| 3.1 Design Token System | 2/2 | Complete | 2026-02-12 |
| 3.2 Activity Page Beta Quality | 0/4 | Planned | - |
| 3.3 Repository Rename & Restructure | 0/3 | Planned | - |
| 4. Add Perspective Flow | 0/2 | Not started | - |
| 5. Testing + Deployment | 1/3 | In progress | - |
| 6. Error Handling & Data Integrity | 0/0 | Not started | - |
| 7. Backend Architecture | 0/0 | Not started | - |
| 8. API & Schema Quality | 0/0 | Not started | - |
| 9. Security Hardening | 0/0 | Not started | - |
| 10. Frontend Quality & Test Coverage | 0/0 | Not started | - |
