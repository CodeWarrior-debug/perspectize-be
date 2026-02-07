<script lang="ts">
	import AgGridSvelte5Component from 'ag-grid-svelte5';
	import { ClientSideRowModelModule } from '@ag-grid-community/client-side-row-model';
	import { themeQuartz } from '@ag-grid-community/theming';
	import type { GridApi, GridOptions } from '@ag-grid-community/core';
	import { nameCellRenderer, durationValueGetter, dateValueFormatter, contentRowId } from '$lib/utils/formatting';

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

	const gridOptions: GridOptions<ContentRow> = {
		columnDefs: [
			{
				field: 'name',
				headerName: 'Title',
				flex: 2,
				sortable: true,
				filter: true,
				cellRenderer: nameCellRenderer,
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
				valueGetter: durationValueGetter,
			},
			{
				field: 'createdAt',
				headerName: 'Date Added',
				width: 140,
				sortable: true,
				valueFormatter: dateValueFormatter,
			},
			{
				field: 'updatedAt',
				headerName: 'Last Updated',
				width: 140,
				sortable: true,
				sort: 'desc',
				valueFormatter: dateValueFormatter,
			}
		],
		pagination: true,
		paginationPageSize: 10,
		paginationPageSizeSelector: [10, 25, 50],
		defaultColDef: {
			resizable: true,
			sortable: true,
		},
		getRowId: contentRowId,
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
