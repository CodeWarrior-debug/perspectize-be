<script lang="ts">
	import type { VideoItem } from '$lib/services/youtubeApi';
	import { toWatchUrl } from '$lib/services/youtubeApi';
	import { Button } from '$lib/components/shadcn';
	import VideoCard from './VideoCard.svelte';

	let {
		items,
		nextPageToken,
		onLoadMore,
		libraryUrls,
		label,
		onAdd,
		pendingId = null,
		isLoading = false,
		isLoadingMore = false,
	}: {
		items: VideoItem[];
		nextPageToken?: string;
		onLoadMore: () => void;
		libraryUrls: Set<string>;
		label?: string;
		onAdd: (videoId: string) => void;
		pendingId?: string | null;
		isLoading?: boolean;
		isLoadingMore?: boolean;
	} = $props();
</script>

<!-- Named generically: renders both search results and trending results (and future Browse results). -->
<div class="flex flex-col gap-4">
	{#if label}
		<p class="text-sm font-medium text-muted-foreground">{label}</p>
	{/if}

	{#if isLoading}
		<div class="flex flex-col gap-4" aria-busy="true" aria-label="Loading videos">
			{#each Array(4) as _, i (i)}
				<div class="flex gap-4 p-3 border border-border rounded-lg bg-card animate-pulse">
					<div class="w-full sm:w-80 h-45 sm:h-[180px] bg-muted rounded shrink-0"></div>
					<div class="flex-1 space-y-2 py-1">
						<div class="h-4 bg-muted rounded w-3/4"></div>
						<div class="h-3 bg-muted rounded w-1/3"></div>
						<div class="h-3 bg-muted rounded w-1/4"></div>
						<div class="h-3 bg-muted rounded w-full"></div>
					</div>
				</div>
			{/each}
		</div>
	{:else if items.length === 0}
		<p class="text-sm text-muted-foreground text-center py-12">No results found</p>
	{:else}
		<div class="flex flex-col gap-4">
			{#each items as video (video.id)}
				<VideoCard
					{video}
					isInLibrary={libraryUrls.has(toWatchUrl(video.id))}
					isPending={pendingId === video.id}
					{onAdd}
				/>
			{/each}
		</div>

		{#if nextPageToken}
			<div class="flex justify-center py-4">
				<Button variant="outline" onclick={onLoadMore} disabled={isLoadingMore}>
					{isLoadingMore ? 'Loading...' : 'Load More'}
				</Button>
			</div>
		{/if}
	{/if}
</div>
