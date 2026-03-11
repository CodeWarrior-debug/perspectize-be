<script lang="ts">
	import AgGridSvelte5Component from 'ag-grid-svelte5';
	import { ClientSideRowModelModule } from '@ag-grid-community/client-side-row-model';
	import { themeQuartz } from '@ag-grid-community/theming';
	import type { GridApi, GridOptions, SortChangedEvent, FilterChangedEvent, ColDef, CellClickedEvent } from '@ag-grid-community/core';
	import { createQuery, keepPreviousData } from '@tanstack/svelte-query';
	import { graphqlRequest } from '$lib/queries/client';
	import { LIST_CONTENT, type ContentItem, type ContentResponse } from '$lib/queries/content';
	import {
		LIST_PERSPECTIVES_BY_USER,
		type ListPerspectivesByUserResponse,
		type PerspectiveItem,
	} from '$lib/queries/perspectives';
	import { queryKeys } from '$lib/queries/keys';
	import { getSelectedUserId } from '$lib/stores/userSelection.svelte';
	import {
		itemCellRenderer,
		typeCellRenderer,
		perspectiveCellRenderer,
		PerspectiveHeaderRenderer,
		durationValueGetter,
		durationFilterValueGetter,
		parseDurationInput,
		formatDurationSeconds,
		dateValueFormatter,
		formatCount,
		formatCountExact,
		formatPublishDate,
		formatTags,
		truncateDescription,
		contentRowId,
		headerMinWidth,
	} from '$lib/utils/formatting';
	import {
		SORT_FIELD_MAP,
		resolveSortField,
		resolveSortOrder,
		capitalizeContentType,
		durationComparator,
		computeNextPage,
		computePrevPage,
	} from '$lib/utils/grid-config';
	import { TagsTooltip } from '$lib/components/TagsTooltip';
	import { DescriptionTooltip } from '$lib/components/DescriptionTooltip';
	import FilterChips from '$lib/components/FilterChips.svelte';
	import PerspectivePopover from '$lib/components/PerspectivePopover.svelte';

	// Popover state for Perspectize column
	let popoverOpen = $state(false);
	let popoverContentId = $state<number | null>(null);
	let popoverContentName = $state('');
	let popoverExistingPerspective = $state<PerspectiveItem | null>(null);

	// State management
	let gridApi = $state<GridApi | null>(null);
	let gridReady = $state(false);
	let pageSize = $state(10);
	let currentPage = $state(0);
	let cursors = $state<(string | null)[]>([null]); // Stack of cursors for pagination
	let sortBy = $state<string>('UPDATED_AT');
	let sortOrder = $state<'ASC' | 'DESC'>('DESC');
	let filterText = $state<string>('');
	let debounceTimer: ReturnType<typeof setTimeout>;
	let activeFilterModel = $state<Record<string, any>>({});
	// Responsive tier: 'xs' (<445px), 'sm' (445-639px), 'md' (640-899px), 'lg' (900px+)
	let responsiveTier = $state<'xs' | 'sm' | 'md' | 'lg'>('lg');
	const isMobile = $derived(responsiveTier === 'xs' || responsiveTier === 'sm');

	// Selected user for perspectives query
	const selectedUserId = $derived(getSelectedUserId());

	// TanStack Query for user's perspectives — used to determine +/glasses icon per row
	const perspectivesQuery = createQuery(() => ({
		queryKey: queryKeys.perspectives.listByUser(selectedUserId ?? 0),
		queryFn: () =>
			graphqlRequest<ListPerspectivesByUserResponse>(LIST_PERSPECTIVES_BY_USER, {
				userID: selectedUserId,
			}),
		enabled: selectedUserId !== null,
		staleTime: 60 * 1000,
	}));

	// O(1) lookup map: contentID → PerspectiveItem
	const perspectivesByContentId = $derived(
		(() => {
			const map = new Map<string, PerspectiveItem>();
			const items = perspectivesQuery.data?.perspectives?.items ?? [];
			for (const p of items) {
				if (p.contentID) map.set(p.contentID, p);
			}
			return map;
		})(),
	);

	// TanStack Query for data fetching
	let currentCursor = $derived(cursors[currentPage]);

	const contentQuery = createQuery(() => ({
		queryKey: queryKeys.content.list({
			sortBy,
			sortOrder,
			search: filterText,
			first: pageSize,
			after: currentCursor,
		}),
		queryFn: async () => {
			const response = await graphqlRequest<ContentResponse>(LIST_CONTENT, {
				first: pageSize,
				after: currentCursor,
				sortBy,
				sortOrder,
				filter: filterText ? { search: filterText } : undefined,
				includeTotalCount: true,
			});

			// Update cursors stack for next page
			if (response.content.pageInfo.hasNextPage && response.content.pageInfo.endCursor) {
				if (cursors.length === currentPage + 1) {
					cursors = [...cursors, response.content.pageInfo.endCursor];
				}
			}

			return response;
		},
		placeholderData: keepPreviousData,
		staleTime: 60 * 1000,
	}));

	// Derived values from query
	const rowData = $derived(contentQuery.data?.content.items ?? []);
	const totalCount = $derived(contentQuery.data?.content.totalCount ?? 0);
	const loading = $derived(contentQuery.isLoading || contentQuery.isPlaceholderData);
	const error = $derived(contentQuery.error);

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

	// flex = clamp-like: proportional sizing with min/max constraints
	// minWidth is auto-derived from headerName unless explicitly set (e.g. Item = 200)
	const columnDefs: ColDef<ContentItem>[] = ([
		{
			colId: 'perspectize',
			headerName: '',
			headerComponent: PerspectiveHeaderRenderer,
			headerTooltip: 'Perspectize — add or edit your perspective',
			flex: 0,
			width: 50,
			minWidth: 50,
			maxWidth: 50,
			sortable: false,
			filter: false,
			resizable: false,
			cellRenderer: perspectiveCellRenderer,
			cellStyle: { display: 'flex', alignItems: 'center', justifyContent: 'center', cursor: 'pointer', padding: 0 },
		},
		{
			colId: 'item',
			headerName: 'Item',
			flex: 2,
			minWidth: 200,

			filter: false, // Search handled by page-level search input
			cellRenderer: itemCellRenderer,
			tooltipValueGetter: (params) => params.data?.name ?? '',
			headerTooltip: 'Video title and thumbnail from YouTube API',
		},
		{
			colId: 'type',
			headerName: 'Type',
			flex: 0.5,
			maxWidth: 100,

			filter: 'agTextColumnFilter',
			valueGetter: (params) => capitalizeContentType(params.data?.contentType),
			filterValueGetter: (params) => {
				return params.data?.contentType?.toLowerCase() ?? '';
			},
			cellRenderer: typeCellRenderer,
			headerTooltip: 'Content type',
		},
		{
			colId: 'duration',
			headerName: 'Length',
			flex: 0.7,
			maxWidth: 120,

			filter: 'agNumberColumnFilter',
			filterParams: {
				allowedCharPattern: '\\d\\:',
				numberParser: parseDurationInput,
				numberFormatter: (value: number | null) =>
					value == null ? null : formatDurationSeconds(value),
			},
			valueGetter: durationValueGetter,
			filterValueGetter: durationFilterValueGetter,
			comparator: durationComparator,
			headerTooltip: 'Video duration from YouTube API',
		},
		{
			colId: 'views',
			field: 'viewCount',
			headerName: 'Views',
			flex: 0.8,
			maxWidth: 130,

			filter: 'agNumberColumnFilter',
			valueFormatter: (params) => formatCount(params.value),
			tooltipValueGetter: (params) => formatCountExact(params.data?.viewCount ?? null),
			headerTooltip: 'View count from YouTube API',
		},
		{
			colId: 'likes',
			field: 'likeCount',
			headerName: 'Likes',
			flex: 0.8,
			maxWidth: 130,

			filter: 'agNumberColumnFilter',
			valueFormatter: (params) => formatCount(params.value),
			tooltipValueGetter: (params) => formatCountExact(params.data?.likeCount ?? null),
			headerTooltip: 'Like count from YouTube API',
		},
		{
			colId: 'publishDate',
			field: 'publishedAt',
			headerName: 'Date',
			flex: 1,
			maxWidth: 150,

			filter: 'agDateColumnFilter',
			filterValueGetter: (params) => {
				const val = params.data?.publishedAt;
				return val ? new Date(val) : null;
			},
			valueFormatter: (params) => formatPublishDate(params.value),
			headerTooltip: 'Publish date from YouTube API',
		},
		{
			colId: 'channel',
			field: 'channelTitle',
			headerName: 'Channel',
			flex: 1.2,
			maxWidth: 200,

			filter: 'agTextColumnFilter',
			headerTooltip: 'Channel name from YouTube API',
		},
		{
			colId: 'tags',
			field: 'tags',
			headerName: 'Tags',
			flex: 1.5,
			maxWidth: 250,
			sortable: false,
			filter: 'agTextColumnFilter',
			filterValueGetter: (params) => formatTags(params.data?.tags ?? null),
			valueFormatter: (params) => formatTags(params.value),
			tooltipComponent: TagsTooltip,
			tooltipField: 'tags',
			headerTooltip: 'Tags from YouTube API',
		},
		{
			colId: 'description',
			field: 'description',
			headerName: 'Description',
			flex: 2,
			sortable: false,
			filter: 'agTextColumnFilter',
			valueFormatter: (params) => truncateDescription(params.value, 80),
			tooltipComponent: DescriptionTooltip,
			tooltipField: 'description',
			headerTooltip: 'Video description from YouTube API',
			hide: true,
		},
		{
			colId: 'updatedAt',
			field: 'updatedAt',
			headerName: 'Updated',
			flex: 1,
			maxWidth: 150,

			filter: 'agDateColumnFilter',
			filterValueGetter: (params) => {
				const val = params.data?.updatedAt;
				return val ? new Date(val) : null;
			},
			valueFormatter: dateValueFormatter,
			headerTooltip: 'Last updated in Perspectize',
			hide: true,
		},

		{
			colId: 'createdAt',
			field: 'createdAt',
			headerName: 'Date Added',
			flex: 1,
			maxWidth: 150,

			filter: 'agDateColumnFilter',
			filterValueGetter: (params) => {
				const val = params.data?.createdAt;
				return val ? new Date(val) : null;
			},
			valueFormatter: dateValueFormatter,
			headerTooltip: 'Date added to Perspectize',
			hide: true,
		},
	] as ColDef<ContentItem>[]).map((col) => ({
		...col,
		minWidth: col.minWidth ?? headerMinWidth(col.headerName ?? '', col.filter !== false),
	}));

	const gridOptions: GridOptions<ContentItem> = {
		columnDefs,
		pagination: false, // Manual pagination
		defaultColDef: {
			resizable: true,

			tooltipValueGetter: (params) => {
				return params.valueFormatted ?? params.value ?? '';
			},
		},
		tooltipShowDelay: 1000,
		tooltipInteraction: true,
		getRowId: contentRowId,
		domLayout: 'normal',
		suppressCellFocus: true,
		context: { perspectivesByContentId: new Map() },
		onCellClicked: (event: CellClickedEvent<ContentItem>) => {
			if (event.colDef.colId !== 'perspectize') return;
			if (!event.data) return;

			const contentId = parseInt(String(event.data.id), 10);
			const existing = perspectivesByContentId.get(String(event.data.id)) ?? null;

			popoverContentId = contentId;
			popoverContentName = event.data.name;
			popoverExistingPerspective = existing;
			popoverOpen = true;
		},
		onGridReady: (params) => {
			gridApi = params.api;
			gridReady = true;
		},
		onSortChanged: (event: SortChangedEvent) => {
			const sortModel = event.api
				.getColumnState()
				.filter((col) => col.sort)
				.sort((a, b) => (a.sortIndex ?? 0) - (b.sortIndex ?? 0));

			if (sortModel.length > 0) {
				const col = sortModel[0];
				sortBy = resolveSortField(col.colId);
				sortOrder = resolveSortOrder(col.sort);
			} else {
				sortBy = 'UPDATED_AT';
				sortOrder = 'DESC';
			}

			// Reset to first page (query auto-refetches via key change)
			currentPage = 0;
			cursors = [null];
		},
		onFilterChanged: (event: FilterChangedEvent) => {
			// Immediate: update chip display
			activeFilterModel = event.api.getFilterModel();

			// Debounce filter changes for server-side search
			clearTimeout(debounceTimer);
			debounceTimer = setTimeout(() => {
				const filterModel = event.api.getFilterModel();
				// Only send Item column filter to server search
				const itemFilter = (filterModel as Record<string, any>)['item']?.filter ?? '';
				if (itemFilter !== filterText) {
					filterText = itemFilter;
					// Reset to first page (query auto-refetches via key change)
					currentPage = 0;
					cursors = [null];
				}
				// Other column filters (Type, Length, etc.) work client-side via AG Grid
			}, 500);
		},
		overlayNoRowsTemplate: '<div class="py-12 text-center text-muted-foreground">No items</div>',
	};

	function handleNextPage() {
		currentPage = computeNextPage(currentPage, totalCount, pageSize);
	}

	function handlePrevPage() {
		currentPage = computePrevPage(currentPage);
	}

	function handlePageSizeChange(newSize: number) {
		pageSize = newSize;
		currentPage = 0;
		cursors = [null];
	}

	// Responsive breakpoint detection — 4 tiers for progressive column reveal
	$effect(() => {
		if (typeof window === 'undefined') return;
		const mqSm = window.matchMedia('(min-width: 445px)');
		const mqMd = window.matchMedia('(min-width: 640px)');
		const mqLg = window.matchMedia('(min-width: 900px)');
		const update = () => {
			if (mqLg.matches) responsiveTier = 'lg';
			else if (mqMd.matches) responsiveTier = 'md';
			else if (mqSm.matches) responsiveTier = 'sm';
			else responsiveTier = 'xs';
		};
		update();
		mqSm.addEventListener('change', update);
		mqMd.addEventListener('change', update);
		mqLg.addEventListener('change', update);
		return () => {
			mqSm.removeEventListener('change', update);
			mqMd.removeEventListener('change', update);
			mqLg.removeEventListener('change', update);
		};
	});

	// Update AG Grid context reactively so perspectiveCellRenderer can access the map
	$effect(() => {
		if (gridApi) {
			gridApi.setGridOption('context', { perspectivesByContentId });
			gridApi.refreshCells({ columns: ['perspectize'], force: true });
		}
	});

	// Responsive column visibility — progressive reveal by tier
	// xs (<445px):  Perspectize, Item, Type
	// sm (445-639): Perspectize, Item, Type, Channel
	// md (640-899): Perspectize, Item, Type, Channel, Duration, Date
	// lg (900+):    Perspectize, Item, Type, Channel, Duration, Date, Views, Likes, Tags
	$effect(() => {
		if (!gridApi || !gridReady) return;
		const api = gridApi;
		const tier = responsiveTier;
		requestAnimationFrame(() => {
			const alwaysVisible = ['item', 'type', 'perspectize'];
			const smCols = ['channel'];
			const mdCols = ['duration', 'publishDate'];
			const lgCols = ['views', 'likes', 'tags'];
			const alwaysHidden = ['description', 'updatedAt', 'createdAt'];

			api.setColumnsVisible(alwaysVisible, true);
			api.setColumnsVisible(alwaysHidden, false);
			api.setColumnsVisible(smCols, tier !== 'xs');
			api.setColumnsVisible(mdCols, tier === 'md' || tier === 'lg');
			api.setColumnsVisible(lgCols, tier === 'lg');
		});
	});

	// Switch to autoHeight on mobile — eliminates empty gap below last row
	$effect(() => {
		if (!gridApi || !gridReady) return;
		gridApi.setGridOption('domLayout', isMobile ? 'autoHeight' : 'normal');
	});

	// Update loading state reactively
	$effect(() => {
		if (gridApi) {
			gridApi.setGridOption('loading', loading);
		}
	});

	// Re-evaluate flex column widths when the grid container resizes
	// (e.g. DevTools panel open/close, sidebar toggle)
	let gridContainer = $state<HTMLDivElement | null>(null);
	$effect(() => {
		if (!gridContainer || !gridApi || !gridReady) return;
		const api = gridApi;
		const observer = new ResizeObserver(() => {
			api.sizeColumnsToFit();
		});
		observer.observe(gridContainer);
		return () => observer.disconnect();
	});
</script>

<div class="flex flex-col h-full gap-4">
	<!-- Error State -->
	{#if contentQuery.isError}
		<div class="flex-1 min-h-0 flex items-center justify-center">
			<div class="text-center py-12 px-4">
				<p class="text-muted-foreground mb-4">Failed to load content. Please try again.</p>
				<button
					onclick={() => contentQuery.refetch()}
					class="px-4 py-2 text-sm font-medium border border-input rounded-md bg-background hover:bg-accent"
				>
					Retry
				</button>
			</div>
		</div>
	{:else}
		<!-- Active Filter Chips -->
		<FilterChips {gridApi} filterModel={activeFilterModel} />

		<!-- AG Grid -->
		<div bind:this={gridContainer} class="{isMobile ? 'overflow-y-auto' : 'flex-1'} min-h-0" style="--ag-row-height: 44px; --ag-header-height: 40px;">
			<AgGridSvelte5Component {gridOptions} {rowData} {theme} {modules} />
		</div>
	{/if}

	<!-- Manual Pagination Controls -->
	<div
		class="shrink-0 flex flex-col md:flex-row items-start md:items-center justify-between gap-2 md:gap-0 px-2 md:px-4 py-2 border-t border-border text-xs md:text-sm"
	>
		<div class="flex items-center gap-2 md:gap-4">
			<div class="text-muted-foreground">
				{totalCount} total
			</div>
			<div class="hidden md:flex items-center gap-2">
				<label for="pageSize" class="text-muted-foreground">Page size:</label>
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
				<span class="hidden sm:inline">Previous</span><span class="sm:hidden">&lt;</span>
			</button>
			<span class="text-muted-foreground">
				Page {currentPage + 1} of {Math.ceil(totalCount / pageSize) || 1}
			</span>
			<button
				onclick={handleNextPage}
				disabled={currentPage >= Math.ceil(totalCount / pageSize) - 1}
				class="px-3 py-1 text-sm border border-input rounded-md bg-background hover:bg-accent disabled:opacity-50 disabled:cursor-not-allowed"
			>
				<span class="hidden sm:inline">Next</span><span class="sm:hidden">&gt;</span>
			</button>
		</div>
	</div>
</div>

<!-- Perspective create/edit modal — rendered outside the grid for correct portal behavior -->
{#if popoverOpen && popoverContentId !== null}
	<PerspectivePopover
		contentId={popoverContentId}
		contentName={popoverContentName}
		existingPerspective={popoverExistingPerspective}
		userId={selectedUserId ?? 0}
		bind:open={popoverOpen}
		onClose={() => {
			popoverOpen = false;
		}}
	/>
{/if}
