# Perspectize

## What This Is

A platform where people can input, browse, and discover perspectives on content — starting with YouTube videos. Perspectives range from a simple agree/disagree with a comment to detailed 0-1000 quality ratings with long-form reviews. The experience is designed to feel effortless to input and calm to browse — users should never feel overwhelmed, lost, interrupted, or out of control (the opposite of social media).

## Core Value

Users can easily submit their perspective on a YouTube video and browse others' perspectives in a way that keeps them in control.

## Requirements

### Validated

*Inferred from existing Go backend codebase:*

- ✓ GraphQL API with gqlgen (schema-first, code generation) — existing
- ✓ Content model with YouTube video metadata (title, description, thumbnails, duration, tags) — existing
- ✓ Perspective model with quality ratings (0-1000), agreement ratings, claims, review text — existing
- ✓ User model with username, email, display name — existing
- ✓ YouTube URL parsing and Data API v3 integration (auto-fetch video metadata) — existing
- ✓ Cursor-based pagination with configurable sort/filter on all list queries — existing
- ✓ PostgreSQL database with migrations (golang-migrate) — existing
- ✓ Hexagonal architecture (domain → ports → services → adapters) — existing
- ✓ Content CRUD operations via GraphQL — existing
- ✓ Perspective CRUD operations via GraphQL — existing
- ✓ User CRUD operations via GraphQL — existing

### Active

**Frontend (SvelteKit Web App):**
- [ ] Discover page: paginated data table of videos with search and filters
- [ ] Video table columns: thumbnail, title, duration, perspective count, avg quality, avg agreement, date added
- [ ] Color-coded score badges for quality/agreement ratings
- [ ] Summary footer with aggregate stats (total videos, total perspectives, avg scores)
- [ ] Add Video flow: paste YouTube URL → backend auto-fetches metadata → content created
- [ ] Add Perspective flow: select a video, fill in perspective form (simple or detailed)
- [ ] Simple perspective: agree/null/undecided/disagree + optional comment
- [ ] Detailed perspective: 0-1000 quality ratings, long paragraph answers, review option
- [ ] User selector dropdown (switch between existing users from database)
- [ ] Responsive layout: mobile-first, scaling to desktop
- [ ] Pagination with page size control

**Design System:**
- [ ] shadcn-svelte component library initialized with Radix 3.0 design tokens
- [ ] Custom theme matching Figma design system (primary navy #1a365d, Inter font, spacing/radius tokens)
- [ ] Data table component with sortable columns, search, and filters

**Backend (fill gaps for frontend):**
- [ ] Any missing GraphQL queries/mutations needed by the frontend
- [ ] Aggregate stats queries (total counts, average scores)

### Out of Scope

- Authentication/login system — v1 uses a simple user dropdown selector, no auth
- AI-assisted perspective input — future feature, not planned for v1
- Side-by-side perspective comparison — v2 feature
- Aggregated views (score distributions) — v2 feature
- Divergence/debate highlighting — v2 feature
- Native mobile app — v1 is responsive web only; Svelte mobile app planned later
- Content types beyond YouTube videos — architecture supports it, but v1 is YouTube only
- Real-time updates / subscriptions — not needed for v1
- Comments/social features — not in v1 scope

## Context

**Existing Backend:** Go 1.25+ backend with PostgreSQL 17, GraphQL (gqlgen), hexagonal architecture. Deployed on Sevalla. Has Content, Perspective, and User models with full CRUD. YouTube Data API integration for auto-fetching video metadata. Cursor-based pagination with sort/filter support.

**Migration:** This repo contains a legacy C# ASP.NET Core implementation in `perspectize/` that is being replaced by the Go backend in `backend/`. The C# code is reference only — all development is in Go.

**Existing Data:** Production database on Sevalla already has users and content data.

**Figma Assets:**
- Design system: Radix 3.0 with shadcn components — [Figma](https://www.figma.com/design/SyvrP9yYbrmCorofJK4Co8/Perspectize---Radix-3.0-Implementation?node-id=2143-2251)
- App design (reference): [Figma](https://www.figma.com/design/K1HaZLeNwCckWvhoyAfRhj/Perspectize-Youtube---Design-1?node-id=3-408&m=dev)
- Prototype: [Figma Site](https://spout-jolly-14050265.figma.site)

**Design Tokens (from Figma):**
- Primary: `#1a365d` (navy), Hover: `#2d3748`, Destructive: `#dc2626`
- Secondary: `#f5f5f5`, Accent: `#f7fafc`
- Font: Inter (full weight scale thin–black)
- Typography: Tailwind-like scale (text-xs through text-9xl)
- Spacing/radius via CSS custom properties (`--spacing/*`, `--border-radius/*`)

**Tracking:** GitHub Issues for all project management.

## Constraints

- **Tech Stack (Backend)**: Go + PostgreSQL + gqlgen — already built, not changing
- **Tech Stack (Frontend)**: SvelteKit + shadcn-svelte + Tailwind CSS — chosen for alignment with Figma design system
- **Hosting**: Inexpensive hosting required — GitHub Pages (static adapter) or similar for frontend; Sevalla for backend
- **Design System**: Must use shadcn-svelte with custom theme matching Radix 3.0 Figma tokens — avoid rework when design system evolves
- **Mobile-First**: Responsive design starting from mobile layouts, scaling up to desktop — need clear breakpoint strategy
- **Deployment**: SvelteKit with static adapter (no SSR) to enable cheap/free static hosting

## Current Milestone: v1.0 Frontend MVP

**Goal:** Build a functional SvelteKit frontend that lets users discover videos, add new videos, and submit perspectives.

**Target features:**
- Discover page with AG Grid data table (search, filter, sort)
- Add Video flow (paste YouTube URL → auto-fetch metadata)
- Add Perspective flow (multi-step form with dynamic fields)
- User selector dropdown (no auth for v1)
- shadcn-svelte design system matching Figma tokens

## Key Decisions

| Decision | Rationale | Outcome |
|----------|-----------|---------|
| SvelteKit for frontend | Modern, lightweight, good DX, Svelte mobile path later | — Pending |
| TanStack Query for data fetching | Official Svelte support, caching, GraphQL integration | — Pending |
| TanStack Form for forms | Multi-step support, dynamic fields, data-type picker needs | — Pending |
| AG Grid for data table | Feature-rich grid, handles sorting/filtering/pagination | — Pending |
| shadcn-svelte for components | Direct mapping from Radix 3.0 Figma design system, avoids rework | — Pending |
| Tailwind CSS for styling | Utility-first, pairs with shadcn-svelte | — Pending |
| Static adapter (no SSR) | Enables GitHub Pages hosting, all data via GraphQL client-side | — Pending |
| User dropdown instead of auth | Simplifies v1, leverages existing users in DB | — Pending |
| Mobile-first responsive | Figma design system TBD on mobile, need breakpoint strategy | — Pending |
| Monorepo (frontend in same repo) | Frontend and backend co-located for easier development | — Pending |

---
*Last updated: 2026-02-04 after milestone v1.0 scope confirmed*
