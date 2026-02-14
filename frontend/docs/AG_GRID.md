# AG Grid Svelte 5 Setup

The `ag-grid-svelte5` wrapper bundles AG Grid v32.x internally. **Do NOT install `ag-grid-community` separately** — it causes version conflicts.

## Installation

```bash
# Pinned to 32.2.x — latest 32.x is 32.3.9 (check before upgrading)
pnpm add ag-grid-svelte5 @ag-grid-community/core@32.2.1 @ag-grid-community/client-side-row-model@32.2.1 @ag-grid-community/theming@32.2.0
```

## Usage

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

## Do NOT

- Import from `ag-grid-community` (use `@ag-grid-community/*`)
- Import AG Grid CSS (use `themeQuartz.withParams()`)
- Use `AgGridSvelte` (use `AgGridSvelte5Component`)
