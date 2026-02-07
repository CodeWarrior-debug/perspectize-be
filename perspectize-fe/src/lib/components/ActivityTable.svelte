<script lang="ts">
	import AgGridSvelte5Component from 'ag-grid-svelte5';
	import { ClientSideRowModelModule } from '@ag-grid-community/client-side-row-model';
	import { themeQuartz } from '@ag-grid-community/theming';
	import type { GridApi, GridOptions } from '@ag-grid-community/core';

	interface ContentRow {
		id: string;
		name: string;
		url: string | null;
		contentType: string;
		length: number | null;
		lengthUnits: string | null;
		createdAt: string;
		updatedAt: string;
	}

	let { rowData = [], loading = false, searchText = '' } = $props<{
		rowData: ContentRow[];
		loading?: boolean;
		searchText?: string;
	}>();

	let gridApi = $state<GridApi | null>(null);

	const modules = [ClientSideRowModelModule];

	const theme = themeQuartz.withParams({
		fontFamily: 'Inter, sans-serif',
		fontSize: 14,
		headerFontSize: 14,
	});

	// Convert length + lengthUnits to display format
	function formatDuration(length: number | null, lengthUnits: string | null): string {
		if (length === null) return '—';

		if (lengthUnits === 'seconds') {
			const minutes = Math.floor(length / 60);
			const seconds = length % 60;
			return `${minutes}:${seconds.toString().padStart(2, '0')}`;
		}

		return `${length} ${lengthUnits}`;
	}

	// Format dates to locale string
	function formatDate(isoString: string): string {
		return new Date(isoString).toLocaleDateString('en-US', {
			year: 'numeric',
			month: 'short',
			day: 'numeric'
		});
	}

	const gridOptions: GridOptions<ContentRow> = {
		columnDefs: [
			{
				field: 'name',
				headerName: 'Title',
				flex: 2,
				sortable: true,
				filter: true,
				cellRenderer: (params: { data?: ContentRow }) => {
					if (!params.data) return '';
					if (params.data.url) {
						return `<a href="${params.data.url}" target="_blank" rel="noopener noreferrer" class="text-primary hover:underline">${params.data.name}</a>`;
					}
					return params.data.name;
				}
			},
			{
				field: 'contentType',
				headerName: 'Type',
				width: 100,
				sortable: true,
			},
			{
				headerName: 'Duration',
				width: 120,
				sortable: false,
				valueGetter: (params) => {
					if (!params.data) return '—';
					return formatDuration(params.data.length, params.data.lengthUnits);
				}
			},
			{
				field: 'createdAt',
				headerName: 'Date Added',
				width: 140,
				sortable: true,
				valueFormatter: (params) => {
					return params.value ? formatDate(params.value) : '—';
				}
			},
			{
				field: 'updatedAt',
				headerName: 'Last Updated',
				width: 140,
				sortable: true,
				sort: 'desc',
				valueFormatter: (params) => {
					return params.value ? formatDate(params.value) : '—';
				}
			}
		],
		pagination: true,
		paginationPageSize: 10,
		paginationPageSizeSelector: [10, 25, 50],
		defaultColDef: {
			resizable: true,
			sortable: true,
		},
		getRowId: (params) => String(params.data?.id ?? ''),
		domLayout: 'autoHeight',
		suppressCellFocus: true,
		onGridReady: (params) => {
			gridApi = params.api;
		}
	};

	// Update loading state reactively
	$effect(() => {
		if (gridApi) {
			gridApi.setGridOption('loading', loading);
		}
	});

	// Update Quick Filter when searchText changes
	$effect(() => {
		if (gridApi) {
			gridApi.setGridOption('quickFilterText', searchText);
		}
	});
</script>

<div class="w-full overflow-x-auto">
	<AgGridSvelte5Component {gridOptions} {rowData} {theme} {modules} />
</div>
