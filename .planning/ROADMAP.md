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
- [ ] **Phase 3: Add Video Flow** - YouTube URL paste, auto-fetch metadata, toast notifications
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
- [ ] 03-01-PLAN.md — YouTube URL validation utility, mutation definition, shadcn Dialog/Input/Label setup
- [ ] 03-02-PLAN.md — AddVideoDialog component with mutation, error handling, Header wiring, visual checkpoint

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
**Plans**: TBD

Plans:
- [ ] 05-01: Coverage verification and gap filling
- [ ] 05-02: Deployment host selection and CI/CD setup
- [ ] 05-03: CORS configuration and production verification

## Progress

**Execution Order:**
Phases execute in numeric order: 1 -> 2 -> 3 -> 4 -> 5

| Phase | Plans Complete | Status | Completed |
|-------|----------------|--------|-----------|
| 1. Foundation | 5/5 | Complete | 2026-02-07 |
| 2. Data Layer + Activity | 2/2 | Complete | 2026-02-07 |
| 3. Add Video Flow | 0/2 | Not started | - |
| 4. Add Perspective Flow | 0/2 | Not started | - |
| 5. Testing + Deployment | 0/3 | Not started | - |
