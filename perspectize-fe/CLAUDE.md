# Frontend: Perspectize SvelteKit

SvelteKit web app with Svelte 5, TanStack Query, AG Grid, shadcn-svelte, and Tailwind v4.

## Architecture

```
perspectize-fe/src/
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
pnpm install          # Install dependencies
pnpm run dev          # Dev server (http://localhost:5173)
pnpm run build        # Production build
pnpm run preview      # Preview production build
pnpm run check        # Type-check (svelte-check + TypeScript)
pnpm run check:watch  # Type-check in watch mode
pnpm run test         # Tests (Vitest, watch mode)
pnpm run test:run     # Tests once
pnpm run test:coverage # Coverage report
```

## Svelte 5 Patterns

This project uses **Svelte 5 runes** exclusively. Do not use Svelte 4 syntax.

| Svelte 5 (use this) | Svelte 4 (do NOT use) |
|----------------------|-----------------------|
| `let count = $state(0)` | `let count = 0` with `$:` |
| `let doubled = $derived(count * 2)` | `$: doubled = count * 2` |
| `let { data, children } = $props()` | `export let data` |
| `$effect(() => { ... })` | `onMount` / `$:` side effects |
| `{@render children()}` | `<slot />` |
| `onclick={handler}` | `on:click={handler}` |

**Key conventions established in this codebase:**

```svelte
<script lang="ts">
  // Props via $props() destructuring
  let { optional = 'default', required } = $props();

  // Reactive state
  let items = $state<Item[]>([]);

  // Derived values (never use $effect for derivation)
  let total = $derived(items.length);

  // Render children via snippet
  let { children } = $props();
</script>

{@render children()}
```

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
    queryFn: () => graphqlClient.request(LIST_CONTENT)
  }));
</script>

<!-- Access as reactive object properties (NO $ prefix) -->
{#if query.isLoading}Loading...{/if}
{#if query.data}{JSON.stringify(query.data)}{/if}
```

**Do NOT:** Use `$query.data` (stores syntax) · Pass options object directly to `createQuery({...})` (must be function wrapper)

## AG Grid Svelte 5 Setup (CRITICAL)

The `ag-grid-svelte5` wrapper bundles AG Grid v32.x internally. **Do NOT install `ag-grid-community` separately** — it causes version conflicts.

```bash
# Pinned to 32.2.x — latest 32.x is 32.3.9 (check before upgrading)
pnpm add ag-grid-svelte5 @ag-grid-community/core@32.2.1 @ag-grid-community/client-side-row-model@32.2.1 @ag-grid-community/theming@32.2.0
```

```svelte
<script lang="ts">
  import AgGridSvelte5Component from 'ag-grid-svelte5';
  import { ClientSideRowModelModule } from '@ag-grid-community/client-side-row-model';
  import { themeQuartz } from '@ag-grid-community/theming';
  import type { GridOptions } from '@ag-grid-community/core';

  const modules = [ClientSideRowModelModule];
  const theme = themeQuartz.withParams({ fontFamily: 'Inter, sans-serif' });
  let rowData = $state<MyRow[]>([]);
  const gridOptions: GridOptions<MyRow> = { columnDefs: [...] };
</script>

<AgGridSvelte5Component {gridOptions} {rowData} {theme} {modules} />
```

**Do NOT:** Import from `ag-grid-community` (use `@ag-grid-community/*`) · Import AG Grid CSS (use `themeQuartz.withParams()`) · Use `AgGridSvelte` (use `AgGridSvelte5Component`)

## Self-Verification (Chrome DevTools MCP)

| Step | Tool | Purpose |
|------|------|---------|
| Navigate | `mcp__chrome-devtools__navigate_page` | Load frontend URL |
| Screenshot | `mcp__chrome-devtools__take_screenshot` | Visual verification |
| Snapshot | `mcp__chrome-devtools__take_snapshot` | DOM structure |
| Resize | `mcp__chrome-devtools__resize_page` | Responsive (375/768/1024px) |
| Console | `mcp__chrome-devtools__list_console_messages` | JS errors |
| Interact | `mcp__chrome-devtools__click` | Buttons, navigation |
