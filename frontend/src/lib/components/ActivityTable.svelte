<script lang="ts">
	import AgGridSvelte5Component from 'ag-grid-svelte5';
	import { ClientSideRowModelModule } from '@ag-grid-community/client-side-row-model';
	import { themeQuartz } from '@ag-grid-community/theming';
	import type { GridApi, GridOptions, SortChangedEvent, ColDef } from '@ag-grid-community/core';
	import { createQuery, keepPreviousData } from '@tanstack/svelte-query';
	import { graphqlClient } from '$lib/queries/client';
	import { LIST_CONTENT, type ContentItem, type ContentResponse } from '$lib/queries/content';
	import { queryKeys } from '$lib/queries/keys';
	import {
		itemCellRenderer,
		typeCellRenderer,
		durationValueGetter,
		dateValueFormatter,
		formatCount,
		formatCountExact,
		formatPublishDate,
		formatDateCompact,
		formatTags,
		truncateDescription,
		contentRowId,
	} from '$lib/utils/formatting';
	import { TagsTooltip } from '$lib/components/TagsTooltip';
	import { DescriptionTooltip } from '$lib/components/DescriptionTooltip';
	import ChevronLeftIcon from '@lucide/svelte/icons/chevron-left';
	import ChevronRightIcon from '@lucide/svelte/icons/chevron-right';

	// Props — searchText is lifted to the page level
	let { searchText = '' }: { searchText?: string } = $props();

	// GraphQL ContentSortBy to AG Grid colId mapping
	const SORT_FIELD_MAP: Record<string, string> = {
		item: 'NAME',
		type: 'NAME', // type not sortable in backend, fallback to NAME
		duration: 'NAME', // duration not sortable, fallback to NAME
		views: 'VIEW_COUNT',
		likes: 'LIKE_COUNT',
		publishDate: 'PUBLISHED_AT',
		channel: 'NAME', // channel not sortable, fallback
		createdAt: 'CREATED_AT',
		updatedAt: 'UPDATED_AT',
	};

	// State management
	let gridApi = $state<GridApi | null>(null);
	let gridReady = $state(false);
	let pageSize = $state(10);
	let currentPage = $state(0);
	let cursors = $state<(string | null)[]>([null]); // Stack of cursors for pagination
	let sortBy = $state<string>('UPDATED_AT');
	let sortOrder = $state<'ASC' | 'DESC'>('DESC');
	let filterText = $state<string>('');
	let searchDebounceTimer: ReturnType<typeof setTimeout>;
	// Responsive tier: 'xs' (<445px), 'sm' (445-639px), 'md' (640-899px), 'lg' (900px+)
	let responsiveTier = $state<'xs' | 'sm' | 'md' | 'lg'>('lg');

	// Sync external search prop to internal filter (debounced)
	$effect(() => {
		const text = searchText;
		clearTimeout(searchDebounceTimer);
		searchDebounceTimer = setTimeout(() => {
			if (text !== filterText) {
				filterText = text;
				currentPage = 0;
				cursors = [null];
			}
		}, 300);
		return () => clearTimeout(searchDebounceTimer);
	});

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
			const response = await graphqlClient.request<ContentResponse>(LIST_CONTENT, {
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
	const totalPages = $derived(Math.ceil(totalCount / pageSize) || 1);
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
	const columnDefs: ColDef<ContentItem>[] = [
		{
			colId: 'item',
			headerName: 'Item',
			flex: 3.5,
			minWidth: 200,
			sortable: true,
			filter: false, // Search handled by page-level search input
			cellRenderer: itemCellRenderer,
			tooltipValueGetter: (params) => params.data?.name ?? '',
			headerTooltip: 'Video title and thumbnail from YouTube API',
		},
		{
			colId: 'type',
			headerName: 'Type',
			flex: 0.5,
			minWidth: 70,
			maxWidth: 90,
			sortable: true,
			filter: 'agTextColumnFilter',
			valueGetter: (params) => {
				const t = params.data?.contentType;
				if (!t) return '';
				return t.charAt(0).toUpperCase() + t.slice(1).toLowerCase();
			},
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
			minWidth: 70,
			maxWidth: 120,
			sortable: true,
			filter: 'agNumberColumnFilter',
			valueGetter: durationValueGetter,
			comparator: (_valueA, _valueB, nodeA, nodeB) => {
				const a = nodeA?.data?.length ?? 0;
				const b = nodeB?.data?.length ?? 0;
				return a - b;
			},
			headerTooltip: 'Video duration from YouTube API',
		},
		{
			colId: 'views',
			field: 'viewCount',
			headerName: 'Views',
			flex: 0.8,
			minWidth: 70,
			maxWidth: 130,
			sortable: true,
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
			minWidth: 70,
			maxWidth: 130,
			sortable: true,
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
			minWidth: 90,
			maxWidth: 150,
			sortable: true,
			filter: 'agDateColumnFilter',
			filterValueGetter: (params) => {
				const val = params.data?.publishedAt;
				return val ? new Date(val) : null;
			},
			valueFormatter: (params) => {
				// Use compact format at md tier to prevent truncation
				if (responsiveTier === 'md') return formatDateCompact(params.value);
				return formatPublishDate(params.value);
			},
			headerTooltip: 'Publish date from YouTube API',
		},
		{
			colId: 'channel',
			field: 'channelTitle',
			headerName: 'Channel',
			flex: 1.2,
			minWidth: 100,
			maxWidth: 200,
			sortable: true,
			filter: 'agTextColumnFilter',
			headerTooltip: 'Channel name from YouTube API',
		},
		{
			colId: 'tags',
			field: 'tags',
			headerName: 'Tags',
			flex: 1.5,
			minWidth: 120,
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
			minWidth: 150,
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
			minWidth: 100,
			maxWidth: 150,
			sortable: true,
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
			minWidth: 100,
			maxWidth: 150,
			sortable: true,
			filter: 'agDateColumnFilter',
			filterValueGetter: (params) => {
				const val = params.data?.createdAt;
				return val ? new Date(val) : null;
			},
			valueFormatter: dateValueFormatter,
			headerTooltip: 'Date added to Perspectize',
			hide: true,
		},
	];

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
				sortBy = SORT_FIELD_MAP[col.colId ?? 'updatedAt'] ?? 'UPDATED_AT';
				sortOrder = col.sort === 'asc' ? 'ASC' : 'DESC';
			} else {
				sortBy = 'UPDATED_AT';
				sortOrder = 'DESC';
			}

			// Reset to first page (query auto-refetches via key change)
			currentPage = 0;
			cursors = [null];
		},
		onFilterChanged: () => {
			// AG Grid column filters work client-side only
			// Server-side search is handled via the searchText prop
		},
		overlayNoRowsTemplate: '<div class="py-12 text-center text-muted-foreground">No items</div>',
	};

	function handleNextPage() {
		if (currentPage < totalPages - 1) {
			currentPage += 1;
		}
	}

	function handlePrevPage() {
		if (currentPage > 0) {
			currentPage -= 1;
		}
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

	// Responsive column visibility — progressive reveal by tier
	// xs (<445px):  Item only (maximize title space)
	// sm (445-639): Item, Type
	// md (640-899): Item, Type, Duration, Date, Channel
	// lg (900+):    Item, Type, Duration, Date, Views, Likes, Channel, Tags
	$effect(() => {
		if (!gridApi || !gridReady) return;
		const api = gridApi;
		const tier = responsiveTier;
		requestAnimationFrame(() => {
			const alwaysVisible = ['item'];
			const smCols = ['type'];
			const mdCols = ['duration', 'publishDate', 'channel'];
			const lgCols = ['views', 'likes', 'tags'];
			const alwaysHidden = ['description', 'updatedAt', 'createdAt'];

			api.setColumnsVisible(alwaysVisible, true);
			api.setColumnsVisible(alwaysHidden, false);
			api.setColumnsVisible(smCols, tier !== 'xs');
			api.setColumnsVisible(mdCols, tier === 'md' || tier === 'lg');
			api.setColumnsVisible(lgCols, tier === 'lg');
		});
	});

	// Update loading state reactively
	$effect(() => {
		if (gridApi) {
			gridApi.setGridOption('loading', loading);
		}
	});
</script>

<div class="flex flex-col h-full">
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
		<!-- AG Grid -->
		<div class="flex-1 min-h-0" style="--ag-row-height: 44px; --ag-header-height: 40px;">
			<AgGridSvelte5Component {gridOptions} {rowData} {theme} {modules} />
		</div>
	{/if}

	<!-- Manual Pagination Controls -->
	<div class="flex items-center justify-between px-3 md:px-4 py-2 border-t border-border text-xs md:text-sm">
		<div class="flex items-center gap-2 md:gap-4">
			<span class="text-muted-foreground whitespace-nowrap">
				{totalCount} total
			</span>
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

		<div class="flex items-center gap-1 md:gap-2">
			<button
				onclick={handlePrevPage}
				disabled={currentPage === 0}
				aria-label="Previous page"
				class="inline-flex items-center justify-center min-w-[36px] min-h-[36px] rounded-md border border-input bg-background hover:bg-accent disabled:opacity-50 disabled:cursor-not-allowed"
			>
				<ChevronLeftIcon class="size-4" />
			</button>
			<span class="text-muted-foreground whitespace-nowrap px-1">
				Page {currentPage + 1} of {totalPages}
			</span>
			<button
				onclick={handleNextPage}
				disabled={currentPage >= totalPages - 1}
				aria-label="Next page"
				class="inline-flex items-center justify-center min-w-[36px] min-h-[36px] rounded-md border border-input bg-background hover:bg-accent disabled:opacity-50 disabled:cursor-not-allowed"
			>
				<ChevronRightIcon class="size-4" />
			</button>
		</div>
	</div>
</div>
