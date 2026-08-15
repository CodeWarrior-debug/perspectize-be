<script lang="ts">
	import type { VideoItem } from '$lib/services/youtubeApi';
	import { Button } from '$lib/components/shadcn';
	import { formatDate } from '$lib/utils/formatting';
	import CheckIcon from '@lucide/svelte/icons/check';

	let {
		video,
		isInLibrary = false,
		isPending = false,
		onAdd,
	}: {
		video: VideoItem;
		isInLibrary?: boolean;
		isPending?: boolean;
		onAdd: (videoId: string) => void;
	} = $props();

	const thumbnailUrl = $derived(
		video.thumbnails.medium?.url ?? video.thumbnails.high?.url ?? video.thumbnails.default?.url ?? '',
	);
</script>

<!-- Named generically: renders both search.list and videos.list (trending) items identically. -->
<div class="flex flex-col sm:flex-row gap-4 p-3 border border-border rounded-lg bg-card text-card-foreground shadow-sm">
	<img
		src={thumbnailUrl}
		alt={video.title}
		width="320"
		height="180"
		class="w-full sm:w-80 h-45 sm:h-[180px] object-cover rounded shrink-0"
	/>
	<div class="flex-1 min-w-0">
		<h3 class="font-medium line-clamp-2">{video.title}</h3>
		<p class="text-sm text-muted-foreground mt-1">{video.channelTitle}</p>
		<p class="text-xs text-muted-foreground mt-1">{formatDate(video.publishedAt)}</p>
		<p class="text-sm text-muted-foreground mt-2 line-clamp-2">{video.description}</p>
	</div>
	<div class="shrink-0 self-start">
		{#if isInLibrary}
			<Button variant="outline" size="sm" disabled>
				<CheckIcon class="size-4" />
				In Library
			</Button>
		{:else}
			<Button variant="default" size="sm" onclick={() => onAdd(video.id)} disabled={isPending}>
				{isPending ? 'Adding...' : 'Add to Library'}
			</Button>
		{/if}
	</div>
</div>
