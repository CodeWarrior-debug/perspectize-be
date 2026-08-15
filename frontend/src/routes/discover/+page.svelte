<script lang="ts">
	import { createQuery } from '@tanstack/svelte-query';
	import PageWrapper from '$lib/components/PageWrapper.svelte';
	import SearchBar from '$lib/components/discover/SearchBar.svelte';
	import VideoResultsGrid from '$lib/components/discover/VideoResultsGrid.svelte';
	import { graphqlRequest } from '$lib/queries/client';
	import { LIST_CONTENT, type ContentResponse } from '$lib/queries/content';
	import { queryKeys } from '$lib/queries/keys';
	import { useAddVideo } from '$lib/queries/hooks/useAddVideo';
	import {
		fetchYouTubeSearch,
		fetchYouTubeTrending,
		toVideoItem,
		toWatchUrl,
		youtubeKeys,
		type VideoItem,
	} from '$lib/services/youtubeApi';

	let searchQuery = $state('');
	let debouncedQuery = $state('');

	// Internal view seam (not a visible toggle yet): a future Browse mode can
	// extend this union ('search' | 'trending' | 'browse') and slot into the
	// same layout without reworking the page.
	const view = $derived<'search' | 'trending'>(debouncedQuery ? 'search' : 'trending');

	const searchResult = createQuery(() => ({
		queryKey: youtubeKeys.search(debouncedQuery),
		queryFn: () => fetchYouTubeSearch({ query: debouncedQuery }),
		enabled: view === 'search',
		staleTime: 5 * 60 * 1000,
	}));

	const trendingResult = createQuery(() => ({
		queryKey: youtubeKeys.trending(),
		queryFn: () => fetchYouTubeTrending(),
		enabled: view === 'trending',
		staleTime: 60 * 60 * 1000,
	}));

	// Accumulated results (first page from the active query + any Load More pages).
	let allResults = $state<VideoItem[]>([]);
	let nextPageToken = $state<string | undefined>(undefined);
	let isLoadingMore = $state(false);
	let loadMoreError = $state<string | null>(null);

	// Sync accumulated results whenever the active query's (first-page) data
	// changes. Because each view has its own dedicated query/queryKey, this
	// also naturally resets accumulation when the search query or view changes.
	$effect(() => {
		if (view === 'search') {
			const data = searchResult.data;
			allResults = data ? data.items.map(toVideoItem) : [];
			nextPageToken = data?.nextPageToken;
		} else {
			const data = trendingResult.data;
			allResults = data ? data.items.map(toVideoItem) : [];
			nextPageToken = data?.nextPageToken;
		}
		loadMoreError = null;
	});

	async function handleLoadMore() {
		if (!nextPageToken || isLoadingMore) return;
		// Capture the query/view this request is for. If the user changes the
		// search query (or the view flips between search/trending) while this
		// request is in flight, the $effect above will already have reset
		// allResults/nextPageToken for the new query — the guard below stops
		// this now-stale response from clobbering that newer state.
		const requestView = view;
		const requestQuery = debouncedQuery;
		const requestToken = nextPageToken;

		isLoadingMore = true;
		loadMoreError = null;
		try {
			if (requestView === 'search') {
				const response = await fetchYouTubeSearch({ query: requestQuery, pageToken: requestToken });
				if (view === requestView && debouncedQuery === requestQuery) {
					allResults = [...allResults, ...response.items.map(toVideoItem)];
					nextPageToken = response.nextPageToken;
				}
			} else {
				const response = await fetchYouTubeTrending('US', requestToken);
				if (view === requestView && debouncedQuery === requestQuery) {
					allResults = [...allResults, ...response.items.map(toVideoItem)];
					nextPageToken = response.nextPageToken;
				}
			}
		} catch (err) {
			if (view === requestView && debouncedQuery === requestQuery) {
				loadMoreError = err instanceof Error ? err.message : 'Failed to load more results';
			}
		} finally {
			// Always clear the in-flight flag (regardless of whether the query
			// changed), otherwise a stale request would leave Load More stuck
			// in its loading state for the new query, which owns its own
			// nextPageToken via the $effect above.
			isLoadingMore = false;
		}
	}

	// Library content, used to determine which results are already added.
	// Note: capped at the API's max page size (100) — an MVP tradeoff for
	// libraries larger than 100 items (see 15-CONTEXT.md quota/caching notes).
	const libraryQuery = createQuery(() => ({
		queryKey: queryKeys.content.lists(),
		queryFn: () => graphqlRequest<ContentResponse>(LIST_CONTENT, { first: 100, includeTotalCount: false }),
	}));

	const libraryUrls = $derived(
		new Set(
			(libraryQuery.data?.content.items ?? [])
				.map((item) => item.url)
				.filter((url): url is string => Boolean(url)),
		),
	);

	const addVideo = useAddVideo();
	let pendingId = $state<string | null>(null);

	function handleAdd(videoId: string) {
		pendingId = videoId;
		addVideo.mutate(toWatchUrl(videoId), {
			onSettled: () => {
				pendingId = null;
			},
		});
	}

	const activeQuery = $derived(view === 'search' ? searchResult : trendingResult);
</script>

<PageWrapper>
	<div class="flex flex-col gap-6">
		<div>
			<h1 class="text-2xl md:text-3xl font-semibold text-foreground">Discover</h1>
			<p class="text-sm text-muted-foreground mt-1">Search YouTube and add videos to your library</p>
		</div>

		<SearchBar bind:value={searchQuery} bind:debouncedQuery />

		{#if activeQuery.isError}
			<p class="text-sm text-destructive">
				{activeQuery.error instanceof Error ? activeQuery.error.message : 'Failed to load videos. Please try again.'}
			</p>
		{:else}
			<VideoResultsGrid
				items={allResults}
				{nextPageToken}
				onLoadMore={handleLoadMore}
				{libraryUrls}
				label={view === 'trending' ? 'Showing Trending Content' : undefined}
				onAdd={handleAdd}
				{pendingId}
				isLoading={activeQuery.isLoading}
				{isLoadingMore}
			/>
			{#if loadMoreError}
				<p class="text-sm text-destructive">{loadMoreError}</p>
			{/if}
		{/if}
	</div>
</PageWrapper>
