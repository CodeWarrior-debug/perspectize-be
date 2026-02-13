<script lang="ts">
	import AgGridSvelte5Component from 'ag-grid-svelte5';
	import { ClientSideRowModelModule } from '@ag-grid-community/client-side-row-model';
	import { themeQuartz } from '@ag-grid-community/theming';
	import type { GridApi, GridOptions, SortChangedEvent, FilterChangedEvent, ColDef } from '@ag-grid-community/core';
	import { graphqlClient } from '$lib/queries/client';
	import { LIST_CONTENT, type ContentItem, type ContentResponse } from '$lib/queries/content';
	import {
		itemCellRenderer,
		typeCellRenderer,
		durationValueGetter,
		dateValueFormatter,
		formatCount,
		formatPublishDate,
		formatTags,
		truncateDescription,
		contentRowId
	} from '$lib/utils/formatting';

	// GraphQL ContentSortBy to AG Grid colId mapping
	const SORT_FIELD_MAP: Record<string, string> = {
		'item': 'NAME',
		'type': 'NAME', // type not sortable in backend, fallback to NAME
		'duration': 'NAME', // duration not sortable, fallback to NAME
		'views': 'VIEW_COUNT',
		'likes': 'LIKE_COUNT',
		'publishDate': 'PUBLISHED_AT',
		'channel': 'NAME', // channel not sortable, fallback
		'createdAt': 'CREATED_AT',
		'updatedAt': 'UPDATED_AT'
	};

	// State management
	let gridApi = $state<GridApi | null>(null);
	let rowData = $state<ContentItem[]>([]);
	let loading = $state(true);
	let totalCount = $state(0);
	let pageSize = $state(10);
	let currentPage = $state(0);
	let cursors = $state<(string | null)[]>([null]); // Stack of cursors for pagination
	let sortBy = $state<string>('UPDATED_AT');
	let sortOrder = $state<'ASC' | 'DESC'>('DESC');
	let filterText = $state<string>('');
	let debounceTimer: ReturnType<typeof setTimeout>;

	const modules = [ClientSideRowModelModule];

	const theme = themeQuartz.withParams({
		fontFamily: "'Geist', system-ui, sans-serif",
		fontSize: 14,
		headerBackgroundColor: '#1a365d',
		headerTextColor: '#ffffff',
		headerFontWeight: 600,
		oddRowBackgroundColor: '#f7fafc',
		rowHoverColor: 'rgba(26, 54, 93, 0.06)',
		borderColor: '#d4d4d4',
		accentColor: '#1a365d',
		foregroundColor: '#171717',
		backgroundColor: '#ffffff',
		selectedRowBackgroundColor: 'rgba(26, 54, 93, 0.08)',
		columnHoverColor: 'rgba(26, 54, 93, 0.04)',
		headerColumnResizeHandleColor: 'rgba(255, 255, 255, 0.5)',
		rowHeight: 44,
		headerHeight: 40,
	});

	const columnDefs: ColDef<ContentItem>[] = [
		{
			colId: 'item',
			headerName: 'Item',
			flex: 2,
			sortable: true,
			filter: true,
			floatingFilter: true,
			cellRenderer: itemCellRenderer,
			headerTooltip: 'Video title and thumbnail from YouTube API'
		},
		{
			colId: 'type',
			headerName: 'Type',
			width: 80,
			sortable: false,
			cellRenderer: typeCellRenderer,
			headerTooltip: 'Content type (YouTube video)'
		},
		{
			colId: 'duration',
			headerName: 'Length',
			width: 100,
			sortable: false,
			valueGetter: durationValueGetter,
			headerTooltip: 'Video duration from YouTube API'
		},
		{
			colId: 'views',
			field: 'viewCount',
			headerName: 'Views',
			width: 100,
			sortable: true,
			floatingFilter: true,
			valueFormatter: (params) => formatCount(params.value),
			headerTooltip: 'View count from YouTube API'
		},
		{
			colId: 'likes',
			field: 'likeCount',
			headerName: 'Likes',
			width: 100,
			sortable: true,
			floatingFilter: true,
			valueFormatter: (params) => formatCount(params.value),
			headerTooltip: 'Like count from YouTube API'
		},
		{
			colId: 'publishDate',
			field: 'publishedAt',
			headerName: 'Date',
			width: 140,
			sortable: true,
			floatingFilter: true,
			valueFormatter: (params) => formatPublishDate(params.value),
			headerTooltip: 'Publish date from YouTube API'
		},
		// Hidden columns
		{
			colId: 'channel',
			field: 'channelTitle',
			headerName: 'Channel',
			width: 160,
			sortable: false,
			hide: true,
			headerTooltip: 'Channel name from YouTube API'
		},
		{
			colId: 'createdAt',
			field: 'createdAt',
			headerName: 'Date Added',
			width: 140,
			sortable: true,
			hide: true,
			valueFormatter: dateValueFormatter,
			headerTooltip: 'Date added to Perspectize'
		},
		{
			colId: 'updatedAt',
			field: 'updatedAt',
			headerName: 'Date Updated',
			width: 140,
			sortable: true,
			hide: true,
			valueFormatter: dateValueFormatter,
			headerTooltip: 'Last updated in Perspectize'
		},
		{
			colId: 'tags',
			field: 'tags',
			headerName: 'Tags',
			width: 200,
			sortable: false,
			hide: true,
			valueFormatter: (params) => formatTags(params.value),
			headerTooltip: 'Tags from YouTube API'
		},
		{
			colId: 'description',
			field: 'description',
			headerName: 'Description',
			flex: 1,
			sortable: false,
			hide: true,
			valueFormatter: (params) => truncateDescription(params.value, 100),
			headerTooltip: 'Video description from YouTube API'
		}
	];

	const gridOptions: GridOptions<ContentItem> = {
		columnDefs,
		pagination: false, // Manual pagination
		defaultColDef: {
			resizable: true,
		},
		getRowId: contentRowId,
		domLayout: 'normal',
		suppressCellFocus: true,
		onGridReady: (params) => {
			gridApi = params.api;
			fetchData();
		},
		onSortChanged: (event: SortChangedEvent) => {
			const sortModel = event.api.getColumnState()
				.filter(col => col.sort)
				.sort((a, b) => (a.sortIndex ?? 0) - (b.sortIndex ?? 0));

			if (sortModel.length > 0) {
				const col = sortModel[0];
				sortBy = SORT_FIELD_MAP[col.colId ?? 'updatedAt'] ?? 'UPDATED_AT';
				sortOrder = col.sort === 'asc' ? 'ASC' : 'DESC';
			} else {
				sortBy = 'UPDATED_AT';
				sortOrder = 'DESC';
			}

			// Reset to first page and fetch
			currentPage = 0;
			cursors = [null];
			fetchData();
		},
		onFilterChanged: (event: FilterChangedEvent) => {
			// Debounce filter changes
			clearTimeout(debounceTimer);
			debounceTimer = setTimeout(() => {
				const filterModel = event.api.getFilterModel();
				// Collect all filter values into a single search string
				const filters = Object.values(filterModel)
					.map((f: any) => f.filter)
					.filter(Boolean)
					.join(' ');
				filterText = filters;

				// Reset to first page and fetch
				currentPage = 0;
				cursors = [null];
				fetchData();
			}, 500);
		},
		overlayNoRowsTemplate: '<div class="py-12 text-center text-muted-foreground">No items yet - add the first one!</div>'
	};

	async function fetchData() {
		loading = true;
		try {
			const response = await graphqlClient.request<ContentResponse>(LIST_CONTENT, {
				first: pageSize,
				after: cursors[currentPage],
				sortBy,
				sortOrder,
				filter: filterText || undefined,
				includeTotalCount: true
			});

			rowData = response.content.items;
			totalCount = response.content.totalCount;

			// Update cursors stack for next page
			if (response.content.pageInfo.hasNextPage && response.content.pageInfo.endCursor) {
				if (cursors.length === currentPage + 1) {
					cursors = [...cursors, response.content.pageInfo.endCursor];
				}
			}
		} catch (error) {
			console.error('Error fetching content:', error);
			rowData = [];
			totalCount = 0;
		} finally {
			loading = false;
		}
	}

	function handleNextPage() {
		if (currentPage < Math.ceil(totalCount / pageSize) - 1) {
			currentPage += 1;
			fetchData();
		}
	}

	function handlePrevPage() {
		if (currentPage > 0) {
			currentPage -= 1;
			fetchData();
		}
	}

	function handlePageSizeChange(newSize: number) {
		pageSize = newSize;
		currentPage = 0;
		cursors = [null];
		fetchData();
	}

	// Update loading state reactively
	$effect(() => {
		if (gridApi) {
			gridApi.setGridOption('loading', loading);
		}
	}

	// Listen for content-added event from AddVideoPopover
	$effect(() => {
		const handleContentAdded = () => {
			// Reset to first page and refetch
			currentPage = 0;
			cursors = [null];
			fetchData();
		};

		window.addEventListener('content-added', handleContentAdded);
		return () => window.removeEventListener('content-added', handleContentAdded);
	});
</script>

<div class="flex flex-col h-full gap-4">
	<!-- AG Grid -->
	<div class="flex-1 min-h-0">
		<AgGridSvelte5Component {gridOptions} {rowData} {theme} {modules} />
	</div>

	<!-- Manual Pagination Controls -->
	<div class="flex items-center justify-between px-4 py-2 border-t border-border">
		<div class="flex items-center gap-4">
			<div class="text-sm text-muted-foreground">
				{totalCount} total items
			</div>
			<div class="flex items-center gap-2">
				<label for="pageSize" class="text-sm text-muted-foreground">Page size:</label>
				<select
					id="pageSize"
					value={pageSize}
					onchange={(e) => handlePageSizeChange(Number(e.currentTarget.value))}
					class="px-2 py-1 text-sm border border-input rounded-md bg-background"
				>
					<option value={10}>10</option>
					<option value={25}>25</option>
					<option value={50}>50</option>
				</select>
			</div>
		</div>

		<div class="flex items-center gap-2">
			<button
				onclick={handlePrevPage}
				disabled={currentPage === 0}
				class="px-3 py-1 text-sm border border-input rounded-md bg-background hover:bg-accent disabled:opacity-50 disabled:cursor-not-allowed"
			>
				Previous
			</button>
			<span class="text-sm text-muted-foreground">
				Page {currentPage + 1} of {Math.ceil(totalCount / pageSize) || 1}
			</span>
			<button
				onclick={handleNextPage}
				disabled={currentPage >= Math.ceil(totalCount / pageSize) - 1}
				class="px-3 py-1 text-sm border border-input rounded-md bg-background hover:bg-accent disabled:opacity-50 disabled:cursor-not-allowed"
			>
				Next
			</button>
		</div>
	</div>
</div>
