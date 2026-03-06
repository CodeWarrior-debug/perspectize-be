<script lang="ts">
	import AgGridSvelte5Component from 'ag-grid-svelte5';
	import { ClientSideRowModelModule } from '@ag-grid-community/client-side-row-model';
	import { themeQuartz } from '@ag-grid-community/theming';
	import type {
		GridApi,
		GridOptions,
		SortChangedEvent,
		FilterChangedEvent,
		ColDef,
		CellClickedEvent,
	} from '@ag-grid-community/core';
	import { createQuery, keepPreviousData } from '@tanstack/svelte-query';
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { graphqlClient } from '$lib/queries/client';
	import { LIST_CONTENT, type ContentItem, type ContentResponse } from '$lib/queries/content';
	import {
		LIST_PERSPECTIVES_BY_USER,
		type ListPerspectivesByUserResponse,
		type PerspectiveItem,
	} from '$lib/queries/perspectives';
	import { queryKeys } from '$lib/queries/keys';
	import { getSelectedUserId } from '$lib/stores/userSelection.svelte';
	import {
		parseGridParams,
		serializeGridParams,
		COL_TO_SORT,
		urlParamsToGraphQLFilter,
		urlParamsToFilter,
		filterToUrlParams,
	} from '$lib/utils/gridUrlState';
	import type { DataMode, GridParams } from '$lib/utils/gridUrlState';
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
	import { TagsTooltip } from '$lib/components/TagsTooltip';
	import { DescriptionTooltip } from '$lib/components/DescriptionTooltip';
	import DataModeToggle from '$lib/components/DataModeToggle.svelte';
	import FilterChips from '$lib/components/FilterChips.svelte';
	import PerspectivePopover from '$lib/components/PerspectivePopover.svelte';

	// Popover state for Perspectize column
	let popoverOpen = $state(false);
	let popoverContentId = $state<number | null>(null);
	let popoverContentName = $state('');
	let popoverExistingPerspective = $state<PerspectiveItem | null>(null);

	// ---------------------------------------------------------------------------
	// URL-derived state
	// ---------------------------------------------------------------------------

	// Derive grid params from current URL (reactive to URL changes)
	const gridParams = $derived(parseGridParams(page.url.searchParams));

	// Individual derived fields from URL
	const mode = $derived(gridParams.mode);
	const sortBy = $derived(COL_TO_SORT[gridParams.sort] ?? 'UPDATED_AT');
	const sortOrder = $derived(gridParams.dir === 'asc' ? 'ASC' : 'DESC');
	const pageNum = $derived(gridParams.page); // 1-indexed
	const pageSize = $derived(gridParams.pageSize);
	const searchText = $derived(gridParams.q);
	const filters = $derived(gridParams.filters);

	// ---------------------------------------------------------------------------
	// URL update helper
	// ---------------------------------------------------------------------------

	function updateUrl(changes: Partial<GridParams>, replace = true) {
		const updated = { ...gridParams, ...changes };
		const search = serializeGridParams(updated);
		const url = search ? `?${search}` : page.url.pathname;
		goto(url, { replaceState: replace, keepFocus: true, noScroll: true });
	}

	// ---------------------------------------------------------------------------
	// Grid state
	// ---------------------------------------------------------------------------

	let gridApi = $state<GridApi | null>(null);
	let gridReady = $state(false);
	let debounceTimer: ReturnType<typeof setTimeout>;
	let skipNextSortEvent = $state(false);
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
			graphqlClient.request<ListPerspectivesByUserResponse>(LIST_PERSPECTIVES_BY_USER, {
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

	// Cursor stack for cursor-based pagination
	// Index = page number (1-indexed: cursors[0] = null for page 1, cursors[1] = cursor for page 2, etc.)
	let cursors = $state<(string | null)[]>([null]);
	const currentCursor = $derived(cursors[pageNum - 1] ?? null); // 1-indexed page → 0-indexed cursor array

	// ---------------------------------------------------------------------------
	// Data fetching (mode-conditional)
	// ---------------------------------------------------------------------------

	// In server-side mode, build GraphQL filter from URL params + search
	// In client-side mode, pass search as simple filter (no column filters)
	const graphqlFilter = $derived(
		mode === 'all' ? urlParamsToGraphQLFilter(filters, searchText) : searchText ? { search: searchText } : undefined,
	);

	const contentQuery = createQuery(() => ({
		queryKey: queryKeys.content.list({
			sortBy: mode === 'all' ? sortBy : 'UPDATED_AT',
			sortOrder: mode === 'all' ? sortOrder : 'DESC',
			search: mode === 'all' ? searchText : '',
			first: pageSize,
			after: currentCursor,
			filter: mode === 'all' ? (graphqlFilter as Record<string, unknown>) : undefined,
			mode,
		}),
		queryFn: async () => {
			const response = await graphqlClient.request<ContentResponse>(LIST_CONTENT, {
				first: mode === 'all' ? pageSize : 100, // Load more in client mode for client-side filtering
				after: mode === 'all' ? currentCursor : null,
				sortBy: mode === 'all' ? sortBy : 'UPDATED_AT',
				sortOrder: mode === 'all' ? sortOrder : 'DESC',
				filter: graphqlFilter,
				includeTotalCount: true,
			});

			// Update cursor stack for page navigation (server-side mode only)
			if (mode === 'all' && response.content.pageInfo.hasNextPage && response.content.pageInfo.endCursor) {
				if (cursors.length === pageNum) {
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
	const hasActiveFilters = $derived(Object.keys(filters).length > 0 || searchText !== '');

	// ---------------------------------------------------------------------------
	// Mode switch handler
	// ---------------------------------------------------------------------------

	function handleModeToggle(newMode: DataMode) {
		// Reset pagination when switching modes
		cursors = [null];

		// When switching Loaded → All: sync AG Grid filter state to URL params
		if (newMode === 'all' && gridApi) {
			const filterModel = gridApi.getFilterModel();
			const urlFilters = filterToUrlParams(filterModel as Record<string, unknown>);
			updateUrl({ mode: newMode, page: 1, filters: urlFilters });
		} else {
			updateUrl({ mode: newMode, page: 1 });
		}
	}

	// ---------------------------------------------------------------------------
	// Pagination handlers
	// ---------------------------------------------------------------------------

	function handleNextPage() {
		if (pageNum < Math.ceil(totalCount / pageSize)) {
			updateUrl({ page: pageNum + 1 });
		}
	}

	function handlePrevPage() {
		if (pageNum > 1) {
			updateUrl({ page: pageNum - 1 });
		}
	}

	function handlePageSizeChange(newSize: number) {
		cursors = [null];
		updateUrl({ pageSize: newSize, page: 1 });
	}

	// ---------------------------------------------------------------------------
	// AG Grid column definitions
	// ---------------------------------------------------------------------------

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
	const columnDefs: ColDef<ContentItem>[] = (
		[
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

				filter: 'agTextColumnFilter',
				filterValueGetter: (params) => params.data?.name ?? '',
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
				maxWidth: 120,

				filter: 'agNumberColumnFilter',
				filterParams: {
					allowedCharPattern: '\\d\\:',
					numberParser: parseDurationInput,
					numberFormatter: (value: number | null) => (value == null ? null : formatDurationSeconds(value)),
				},
				valueGetter: durationValueGetter,
				filterValueGetter: durationFilterValueGetter,
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
		] as ColDef<ContentItem>[]
	).map((col) => ({
		...col,
		minWidth: col.minWidth ?? headerMinWidth(col.headerName ?? '', col.filter !== false),
	}));

	// ---------------------------------------------------------------------------
	// AG Grid options (mode-conditional event handlers)
	// ---------------------------------------------------------------------------

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
			// In "Loaded" mode, AG Grid handles client-side sort — skip URL update
			if (mode === 'loaded') return;
			// Skip if we triggered this event programmatically (to avoid loop)
			if (skipNextSortEvent) {
				skipNextSortEvent = false;
				return;
			}

			// Server-side: update URL → triggers refetch
			const sortModel = event.api
				.getColumnState()
				.filter((col) => col.sort)
				.sort((a, b) => (a.sortIndex ?? 0) - (b.sortIndex ?? 0));

			if (sortModel.length > 0) {
				const col = sortModel[0];
				cursors = [null]; // Reset cursor stack
				updateUrl({ sort: col.colId ?? 'updatedAt', dir: col.sort === 'asc' ? 'asc' : 'desc', page: 1 });
			} else {
				cursors = [null];
				updateUrl({ sort: 'updatedAt', dir: 'desc', page: 1 });
			}
		},
		onFilterChanged: (event: FilterChangedEvent) => {
			// Immediate: update chip display
			activeFilterModel = event.api.getFilterModel();

			// In "Loaded" mode, AG Grid handles client-side filter — skip URL update
			if (mode === 'loaded') return;

			// Server-side: debounce → convert filter model → update URL
			clearTimeout(debounceTimer);
			debounceTimer = setTimeout(() => {
				const filterModel = event.api.getFilterModel();
				const urlFilters = filterToUrlParams(filterModel as Record<string, unknown>);
				cursors = [null];
				updateUrl({ filters: urlFilters, page: 1 });
			}, 500);
		},
		// Prevent AG Grid from client-side re-sorting data that arrives pre-sorted from server
		postSortRows: () => {
			// In "All Items" mode, data is server-sorted — suppress AG Grid re-sort
			// In "Loaded" mode, let AG Grid sort normally (default behavior)
			if (mode === 'all') return;
		},
		overlayNoRowsTemplate: '<div class="py-12 text-center text-muted-foreground">No content yet</div>',
	};

	// ---------------------------------------------------------------------------
	// Effects
	// ---------------------------------------------------------------------------

	// Sync AG Grid sort indicators to URL state when switching to "All Items" mode
	// Uses skipNextSortEvent flag to prevent onSortChanged from re-updating the URL
	$effect(() => {
		if (!gridApi || !gridReady || mode !== 'all') return;
		skipNextSortEvent = true;
		gridApi.applyColumnState({
			state: [{ colId: gridParams.sort, sort: gridParams.dir }],
			defaultState: { sort: null },
		});
	});

	// Restore AG Grid filter state from URL on mount and mode changes
	$effect(() => {
		if (!gridApi || !gridReady) return;
		const filterModel = urlParamsToFilter(gridParams.filters);
		gridApi.setFilterModel(Object.keys(filterModel).length > 0 ? filterModel : null);
	});

	// Update loading state reactively
	$effect(() => {
		if (gridApi) {
			gridApi.setGridOption('loading', loading);
		}
	});

	// Update empty state message based on active filters
	$effect(() => {
		if (!gridApi) return;
		const template = hasActiveFilters
			? '<div class="py-12 text-center"><p class="text-muted-foreground mb-1">No results match your filters</p><p class="text-xs text-muted-foreground/70">Try adjusting or clearing your filters</p></div>'
			: '<div class="py-12 text-center text-muted-foreground">No content yet</div>';
		gridApi.setGridOption('overlayNoRowsTemplate', template);
	});

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
	<!-- Active Filter Chips — always visible so users can clear filters even during errors -->
	<FilterChips {gridApi} filterModel={activeFilterModel} />

	<!-- Error State -->
	{#if contentQuery.isError}
		<div class="flex-1 min-h-0 flex items-center justify-center">
			<div class="text-center py-12 px-4">
				<p class="text-muted-foreground mb-2">Failed to load content. Please try again.</p>
				{#if hasActiveFilters}
					<p class="text-xs text-muted-foreground/70 mb-4">
						Your active filters may be causing this issue. Try clearing them.
					</p>
				{/if}
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
		<div
			bind:this={gridContainer}
			class="{isMobile ? 'overflow-y-auto' : 'flex-1'} min-h-0"
			style="--ag-row-height: 44px; --ag-header-height: 40px;"
		>
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
			<!-- Data Mode Toggle -->
			<DataModeToggle {mode} loadedCount={rowData.length} onToggle={handleModeToggle} />
			{#if mode === 'all'}
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
			{/if}
		</div>

		{#if mode === 'all'}
			<div class="flex items-center gap-2">
				<button
					onclick={handlePrevPage}
					disabled={pageNum <= 1}
					class="px-3 py-1 text-sm border border-input rounded-md bg-background hover:bg-accent disabled:opacity-50 disabled:cursor-not-allowed"
				>
					<span class="hidden sm:inline">Previous</span><span class="sm:hidden">&lt;</span>
				</button>
				<span class="text-muted-foreground">
					Page {pageNum} of {Math.ceil(totalCount / pageSize) || 1}
				</span>
				<button
					onclick={handleNextPage}
					disabled={pageNum >= Math.ceil(totalCount / pageSize)}
					class="px-3 py-1 text-sm border border-input rounded-md bg-background hover:bg-accent disabled:opacity-50 disabled:cursor-not-allowed"
				>
					<span class="hidden sm:inline">Next</span><span class="sm:hidden">&gt;</span>
				</button>
			</div>
		{/if}
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
