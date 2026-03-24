# Frontend: Perspectize SvelteKit

SvelteKit web app with Svelte 5, TanStack Query, AG Grid, shadcn-svelte, and Tailwind v4.

## Architecture

```
frontend/src/
├── routes/              # SvelteKit file-based routing
│   ├── +layout.svelte   # Root layout (QueryClientProvider, Header, Toaster)
│   ├── +layout.ts       # Layout config (prerender = true)
│   └── +page.svelte     # Home page
├── lib/
│   ├── components/      # Svelte 5 components
│   │   ├── shadcn/      # shadcn-svelte primitives (button/)
│   │   ├── Header.svelte
│   │   ├── PageWrapper.svelte
│   │   └── AGGridTest.svelte
│   ├── queries/         # TanStack Query + graphql-request
│   │   ├── client.ts    # GraphQLClient (VITE_GRAPHQL_URL)
│   │   └── content.ts   # Content query definitions (gql)
│   ├── assets/          # Static assets (favicon)
│   └── utils/           # Utility functions
├── app.css              # Global styles (Tailwind v4)
└── app.html             # HTML shell
```

## Tailwind v4

Tailwind v4 uses `--color-*` prefix for theme variables (e.g., `--color-primary`), not bare `--primary` from v3/shadcn conventions.

## Commands

```bash
pnpm run dev          # Dev server (http://localhost:5173)
pnpm run check        # Type-check (svelte-check + TypeScript)
pnpm run test:run     # Tests once (CI/verification)
pnpm run test         # Tests in watch mode
```

**`pnpm exec` must run from `frontend/`** — running from repo root fails with `ERR_PNPM_RECURSIVE_EXEC_NO_PACKAGE`. Use `cd frontend && pnpm exec ...` or `pnpm --dir frontend exec ...`.

## Svelte 5 Patterns

This project uses **Svelte 5 runes** exclusively. Do not use Svelte 4 syntax.

| Svelte 5 (use this)                 | Svelte 4 (do NOT use)         |
| ----------------------------------- | ----------------------------- |
| `let count = $state(0)`             | `let count = 0` with `$:`     |
| `let doubled = $derived(count * 2)` | `$: doubled = count * 2`      |
| `let { data, children } = $props()` | `export let data`             |
| `$effect(() => { ... })`            | `onMount` / `$:` side effects |
| `{@render children()}`              | `<slot />`                    |
| `onclick={handler}`                 | `on:click={handler}`          |

**Additional rules:** Never use `$effect` for derivation (use `$derived`). Render children via `{@render children()}` with `let { children } = $props()`.

## TanStack Query + GraphQL

Queries use `graphql-request` with TanStack Svelte Query.

1. Define queries in `lib/queries/` using `gql` tagged templates
2. Client in `lib/queries/client.ts` — uses `VITE_GRAPHQL_URL` (defaults to `http://localhost:8080/graphql`)
3. QueryClientProvider wraps app in `+layout.svelte` with `enabled: browser` to prevent SSR queries

**Svelte 5 API (CRITICAL):** TanStack Query v5+ with Svelte 5 uses a **function wrapper** pattern. Query results are reactive objects, NOT stores (no `$` prefix).

```svelte
<script lang="ts">
	import { createQuery } from '@tanstack/svelte-query';
	import { graphqlClient } from '$lib/queries/client';
	import { LIST_CONTENT } from '$lib/queries/content';

	// Function wrapper pattern — pass a function returning options
	const query = createQuery(() => ({
		queryKey: ['content'],
		queryFn: () => graphqlClient.request(LIST_CONTENT),
	}));
</script>

<!-- Access as reactive object properties (NO $ prefix) -->
{#if query.isLoading}Loading...{/if}
{#if query.data}{JSON.stringify(query.data)}{/if}
```

**Do NOT:** Use `$query.data` (stores syntax) · Pass options object directly to `createQuery({...})` (must be function wrapper)

## Icons (Lucide)

Per-icon imports from `@lucide/svelte` for tree-shaking:

```svelte
<script lang="ts">
	import XIcon from '@lucide/svelte/icons/x';
	import PlusIcon from '@lucide/svelte/icons/plus';
</script>

<XIcon class="size-4" />
```

- Import path: `@lucide/svelte/icons/{kebab-case-name}`
- Default size in buttons: `size-4` (applied via button base). Override with explicit `size-*` class.

## Design Tokens

All tokens defined in `src/app.css` under `@theme`. Full set: primary, secondary, muted, accent, destructive, card, popover, border, input, ring, rating colors, brand.

- **Fonts:** Geist (UI/headings) + Charter (body/reading text)
- **Colors:** Hex values in `--color-*` custom properties (Tailwind v4 convention)
- **No external token pipeline** — tokens defined directly in CSS, consumed via Tailwind utilities

## AG Grid Svelte 5 Setup (CRITICAL)

Uses `ag-grid-svelte5` wrapper (bundles AG Grid v32.x). **Do NOT install `ag-grid-community` separately.** Import from `@ag-grid-community/*`, use `AgGridSvelte5Component`, style via `themeQuartz.withParams()`.

Full setup and examples: [docs/AG_GRID.md](docs/AG_GRID.md)

**Column visibility gotcha:** `ActivityTable.svelte` controls visibility in TWO places that must stay in sync: `hide: true` in colDef (initial default) and `$effect` with `setColumnsVisible()` (responsive override — runs on gridReady, always wins). When adding/changing columns, update BOTH. The responsive system uses 4 tiers: xs (<445px), sm (445-639px), md (640-899px), lg (900px+) — decide which tier(s) should show the new column. See [ADDING_AG_GRID_COLUMN.md](../.claude/docs/ADDING_AG_GRID_COLUMN.md) Decision 7.

## Figma Design Workflow

- **[docs/FIGMA.md](docs/FIGMA.md)** — Figma file reference (file keys, pages, variables, code↔Figma mapping)
- **[docs/FIGMA_VERIFICATION.md](docs/FIGMA_VERIFICATION.md)** — Verification guide for Figma Make outputs
- **[Code to Figma Canvas](../.claude/docs/CODE_TO_FIGMA_CANVAS.md)** — Capture running app into Figma to keep designs in sync

**Code-to-Figma capture gotchas:**
- CSP in `app.html` blocks `mcp.figma.com` — temporarily add to `script-src` and `connect-src`, revert after capture
- AG Grid canvas-rendered cells (thumbnails) don't serialize into Figma captures
- SPA hash navigation may not re-trigger auto-capture — reload or click "Send to Figma" manually

## Self-Verification (Chrome DevTools MCP)

| Step       | Tool                                          | Purpose                        |
| ---------- | --------------------------------------------- | ------------------------------ |
| Navigate   | `mcp__chrome-devtools__navigate_page`         | Load frontend URL              |
| Screenshot | `mcp__chrome-devtools__take_screenshot`       | Visual verification            |
| Snapshot   | `mcp__chrome-devtools__take_snapshot`         | DOM structure                  |
| Resize     | `mcp__chrome-devtools__resize_page`           | Desktop breakpoints (768+)     |
| Emulate    | `mcp__chrome-devtools__emulate`               | Mobile/tablet device emulation |
| Console    | `mcp__chrome-devtools__list_console_messages` | JS errors                      |
| Interact   | `mcp__chrome-devtools__click`                 | Buttons, navigation            |

**Mobile screenshots:** Use `emulate` (NOT `resize_page`) for mobile/tablet. AG Grid and CSS media queries only respond to viewport changes when `isMobile: true` is set. `resize_page` alone doesn't trigger responsive column hiding. Example: `emulate({ viewport: { width: 375, height: 812, deviceScaleFactor: 3, isMobile: true, hasTouch: true } })`. Reset with `emulate({ viewport: null, userAgent: null })` after. Use `resize_page` for desktop breakpoint comparisons (768px+).

## Testing Gotchas

**Date formatting timezone:** `formatDate`/`formatDateCompact` use `toLocaleDateString` (local timezone). In tests, use midday UTC times (`T12:00:00Z`) not midnight (`T00:00:00Z`) to avoid dates shifting to previous day in US timezones.

**AG Grid testing strategy:** AG Grid doesn't render in jsdom — no lifecycle hooks, no Grid API, no cell rendering. Test AG Grid logic by extracting pure functions into `$lib/utils/grid-config.ts` (sort mapping, pagination bounds, responsive tiers, comparators, column metadata). Test renderers/formatters via `$lib/utils/formatting.ts`. For grid integration (filter UI, sort clicks, responsive `$effect` blocks), use Playwright E2E or Vitest Browser Mode (future). See [ADDING_AG_GRID_COLUMN.md](../.claude/docs/ADDING_AG_GRID_COLUMN.md) testing section.
