<script lang="ts">
	import AgGridSvelte5Component from 'ag-grid-svelte5';
	import { ClientSideRowModelModule } from '@ag-grid-community/client-side-row-model';
	import { themeQuartz } from '@ag-grid-community/theming';
	import type { GridOptions } from '@ag-grid-community/core';
	import { Button } from '$lib/components/shadcn';

	interface VideoRow {
		id: number;
		title: string;
		duration: string;
		rating: number;
		published: string;
	}

	let rowData = $state<VideoRow[]>([
		{ id: 1, title: 'Introduction to Svelte 5', duration: '15:42', rating: 92, published: '2026-01-15' },
		{ id: 2, title: 'Building with SvelteKit', duration: '23:18', rating: 87, published: '2026-01-20' },
		{ id: 3, title: 'Tailwind CSS Deep Dive', duration: '31:05', rating: 95, published: '2026-01-25' },
		{ id: 4, title: 'GraphQL Fundamentals', duration: '18:33', rating: 78, published: '2026-02-01' },
		{ id: 5, title: 'TypeScript Best Practices', duration: '27:41', rating: 91, published: '2026-02-03' },
		{ id: 6, title: 'State Management Patterns', duration: '22:15', rating: 84, published: '2026-02-05' },
		{ id: 7, title: 'Testing Svelte Components', duration: '19:28', rating: 89, published: '2026-02-07' },
		{ id: 8, title: 'Performance Optimization', duration: '25:52', rating: 93, published: '2026-02-10' },
		{ id: 9, title: 'Deployment Strategies', duration: '16:44', rating: 81, published: '2026-02-12' },
		{ id: 10, title: 'Accessibility in Web Apps', duration: '20:36', rating: 96, published: '2026-02-15' },
		{ id: 11, title: 'Advanced Svelte Patterns', duration: '28:19', rating: 90, published: '2026-02-18' },
		{ id: 12, title: 'API Design Principles', duration: '24:07', rating: 86, published: '2026-02-20' },
	]);

	const modules = [ClientSideRowModelModule];

	const theme = themeQuartz.withParams({
		fontFamily: 'Inter, sans-serif',
		fontSize: 14,
		headerFontSize: 14,
	});

	const gridOptions: GridOptions<VideoRow> = {
		columnDefs: [
			{ field: 'id', headerName: 'ID', width: 80, sortable: true },
			{ field: 'title', headerName: 'Title', flex: 2, filter: true, sortable: true },
			{ field: 'duration', headerName: 'Duration', width: 120, sortable: true },
			{ field: 'rating', headerName: 'Rating', width: 100, filter: 'agNumberColumnFilter', sortable: true },
			{ field: 'published', headerName: 'Published', width: 130, sortable: true },
		],
		pagination: true,
		paginationPageSize: 5,
		paginationPageSizeSelector: [5, 10, 25],
		rowSelection: { mode: 'multiRow' },
		defaultColDef: {
			resizable: true,
		},
		getRowId: (params) => String(params.data?.id ?? ''),
		domLayout: 'autoHeight',
	};

	let nextId = rowData.length + 1;

	function addRow() {
		rowData = [...rowData, {
			id: nextId++,
			title: `New Video ${nextId - 1}`,
			duration: '10:00',
			rating: 80,
			published: new Date().toISOString().split('T')[0],
		}];
	}
</script>

<div class="space-y-4">
	<h2 class="text-xl font-semibold">AG Grid Svelte 5 Validation</h2>

	<div class="p-4 bg-muted rounded-lg">
		<h3 class="font-medium mb-2">Validation Checklist:</h3>
		<ul class="space-y-1 text-sm">
			<li>[ ] <strong>Sorting:</strong> Click column headers to sort asc/desc</li>
			<li>[ ] <strong>Filtering:</strong> Click filter icon, enter filter values</li>
			<li>[ ] <strong>Pagination:</strong> Use page controls at bottom</li>
			<li>[ ] <strong>Column Resize:</strong> Drag column borders to resize</li>
			<li>[ ] <strong>Row Selection:</strong> Click rows to select (checkboxes)</li>
			<li>[ ] <strong>Reactivity:</strong> Click "Add Row" button, new row appears</li>
		</ul>
	</div>

	<Button variant="secondary" onclick={addRow}>Add Row (Test Reactivity)</Button>

	<div class="w-full">
		<AgGridSvelte5Component {gridOptions} {rowData} {theme} {modules} />
	</div>
</div>
