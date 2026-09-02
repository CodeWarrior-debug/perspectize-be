<script lang="ts">
	import PlayIcon from '@lucide/svelte/icons/play';
	import { extractVideoIdFromUrl, formatDuration } from '$lib/utils/formatting';

	interface CardRow {
		id: string | number;
		name: string;
		url: string | null;
		channelTitle: string | null;
		length: number | null;
		lengthUnits: string | null;
	}

	let {
		rowData,
		onOpenDetails,
	}: {
		rowData: CardRow[];
		onOpenDetails: (contentId: string) => void;
	} = $props();

	function thumbSrc(row: CardRow): string | null {
		const videoId = extractVideoIdFromUrl(row.url);
		return videoId ? `https://i.ytimg.com/vi/${videoId}/hqdefault.jpg` : null;
	}

	function handleThumbClick(row: CardRow, e: MouseEvent) {
		e.stopPropagation();
		if (row.url) window.open(row.url, '_blank', 'noopener,noreferrer');
	}
</script>

<div data-testid="activity-card-list" class="flex flex-col gap-2.5 px-2 py-2">
	{#each rowData as row (row.id)}
		<div
			class="flex items-center gap-3 rounded-lg border border-border bg-card p-2.5 hover:bg-primary/[0.06]"
		>
			<button
				type="button"
				data-testid={`card-thumb-${row.id}`}
				title="Open original content in new tab"
				class="relative h-16 w-24 flex-none overflow-hidden rounded-md bg-muted"
				onclick={(e) => handleThumbClick(row, e)}
			>
				{#if thumbSrc(row)}
					<img src={thumbSrc(row)} alt="" class="h-full w-full object-cover" />
				{/if}
				<span
					class="absolute right-1 bottom-1 flex items-center justify-center rounded bg-[rgba(23,23,23,0.65)] p-1"
				>
					<PlayIcon class="size-2.5 fill-white text-white" />
				</span>
			</button>

			<button
				type="button"
				title="View content data + details"
				class="min-w-0 flex-1 text-left"
				onclick={() => onOpenDetails(String(row.id))}
			>
				<div
					class="line-clamp-2 font-[family-name:var(--font-family-serif)] text-sm leading-tight font-semibold text-foreground"
				>
					{row.name}
				</div>
				<div class="mt-1.5 flex items-center gap-2 text-xs text-muted-foreground">
					{#if row.channelTitle}
						<span>{row.channelTitle}</span>
						<span>&middot;</span>
					{/if}
					<span>{formatDuration(row.length, row.lengthUnits)}</span>
				</div>
			</button>
		</div>
	{/each}
</div>
