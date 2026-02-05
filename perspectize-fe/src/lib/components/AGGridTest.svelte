<script lang="ts">
	import AgGridSvelte from 'ag-grid-svelte5';

	// Test data simulating video content
	let rowData = $state([
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

	let gridOptions = $state({
		columnDefs: [
			{ field: 'id', headerName: 'ID', width: 80, sortable: true },
			{ field: 'title', headerName: 'Title', flex: 2, filter: 'agTextColumnFilter', sortable: true },
			{ field: 'duration', headerName: 'Duration', width: 120, sortable: true },
			{ field: 'rating', headerName: 'Rating', width: 100, filter: 'agNumberColumnFilter', sortable: true },
			{ field: 'published', headerName: 'Published', width: 130, filter: 'agDateColumnFilter', sortable: true },
		],
		rowData,
		pagination: true,
		paginationPageSize: 5,
		paginationPageSizeSelector: [5, 10, 25],
		rowSelection: 'multiple',
		defaultColDef: {
			resizable: true,
		},
		getRowId: (params: { data: { id: number } }) => String(params.data.id),
	});

	// Track validation results
	let validationResults = $state<Record<string, boolean | null>>({
		sorting: null,
		filtering: null,
		pagination: null,
		columnResize: null,
		rowSelection: null,
		reactivity: null,
	});

	function addRow() {
		const newId = rowData.length + 1;
		rowData = [...rowData, {
			id: newId,
			title: `New Video ${newId}`,
			duration: '10:00',
			rating: 80,
			published: new Date().toISOString().split('T')[0],
		}];
		gridOptions.rowData = rowData;
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

	<button
		class="px-4 py-2 bg-secondary text-secondary-foreground rounded hover:bg-secondary/80"
		onclick={addRow}
	>
		Add Row (Test Reactivity)
	</button>

	<div class="ag-theme-quartz" style="height: 400px; width: 100%;">
		<AgGridSvelte {gridOptions} />
	</div>
</div>
