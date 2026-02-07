# Frontend: Perspectize SvelteKit

SvelteKit web app with TanStack Query, AG Grid, and shadcn-svelte.

## Architecture

```
perspectize-fe/src/
├── routes/           # SvelteKit file-based routing
│   ├── +layout.svelte   # Root layout (Header, QueryProvider)
│   ├── +layout.ts       # Layout load function
│   └── +page.svelte     # Home page
├── lib/
│   ├── components/      # Svelte components
│   │   ├── ui/          # shadcn-svelte primitives (button/)
│   │   ├── Header.svelte
│   │   ├── PageWrapper.svelte
│   │   └── AGGridTest.svelte
│   ├── queries/         # TanStack Query + graphql-request
│   │   ├── client.ts    # GraphQLClient setup (VITE_GRAPHQL_URL)
│   │   └── content.ts   # Content query definitions (gql)
│   ├── assets/          # Static assets
│   └── utils/           # Utility functions
├── app.css              # Global styles (Tailwind v4)
└── app.html             # HTML shell
```

## Development Commands

```bash
pnpm install          # Install dependencies
pnpm run dev          # Dev server (typically http://localhost:5173)
pnpm run build        # Production build
pnpm run preview      # Preview production build
pnpm run check        # Type-check (svelte-check + TypeScript)
pnpm run check:watch  # Type-check in watch mode
pnpm run test         # Run tests (Vitest, watch mode)
pnpm run test:run     # Run tests once
pnpm run test:coverage # Coverage report
```

## TanStack Query + GraphQL

Queries use `graphql-request` with TanStack Svelte Query. Pattern:

1. Define queries in `lib/queries/` using `gql` tagged templates
2. Client configured in `lib/queries/client.ts` — uses `VITE_GRAPHQL_URL` env var (defaults to `http://localhost:8080/graphql`)
3. Use `createQuery` in components with the shared `graphqlClient`

```svelte
<script lang="ts">
  import { createQuery } from '@tanstack/svelte-query';
  import { graphqlClient } from '$lib/queries/client';
  import { LIST_CONTENT } from '$lib/queries/content';

  const query = createQuery({ queryKey: ['content'], queryFn: () => graphqlClient.request(LIST_CONTENT) });
</script>
```

## AG Grid Svelte 5 Setup (CRITICAL)

The `ag-grid-svelte5` wrapper bundles AG Grid v32.x internally. **Do NOT install `ag-grid-community` separately** — it causes version conflicts.

**Correct setup:**
```bash
# Pinned to 32.2.x — latest 32.x is 32.3.9 (check before upgrading)
pnpm add ag-grid-svelte5 @ag-grid-community/core@32.2.1 @ag-grid-community/client-side-row-model@32.2.1 @ag-grid-community/theming@32.2.0
```

**Correct imports:**
```svelte
<script lang="ts">
  import AgGridSvelte5Component from 'ag-grid-svelte5';
  import { ClientSideRowModelModule } from '@ag-grid-community/client-side-row-model';
  import { themeQuartz } from '@ag-grid-community/theming';
  import type { GridOptions } from '@ag-grid-community/core';

  const modules = [ClientSideRowModelModule];
  const theme = themeQuartz.withParams({ fontFamily: 'Inter, sans-serif' });
</script>

<AgGridSvelte5Component {gridOptions} {rowData} {theme} {modules} />
```

**Do NOT:**
- Import from `ag-grid-community` (use `@ag-grid-community/*` packages)
- Import AG Grid CSS in app.css (use `themeQuartz.withParams()`)
- Use `AgGridSvelte` (use `AgGridSvelte5Component`)

## Self-Verification (Chrome DevTools MCP)

| Step | MCP Tool | Purpose |
|------|----------|---------|
| Navigate | `mcp__chrome-devtools__navigate_page` | Load frontend URL |
| Screenshot | `mcp__chrome-devtools__take_screenshot` | Visual verification |
| Snapshot | `mcp__chrome-devtools__take_snapshot` | DOM/component structure |
| Resize | `mcp__chrome-devtools__resize_page` | Responsive check (375px, 768px, 1024px) |
| Console | `mcp__chrome-devtools__list_console_messages` | Check for JS errors |
| Interact | `mcp__chrome-devtools__click` | Test buttons, toasts, navigation |
