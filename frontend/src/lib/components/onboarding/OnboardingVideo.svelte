<script lang="ts">
	let {
		src,
		label = 'Watch',
		class: className = '',
	}: {
		src?: string | null;
		label?: string;
		class?: string;
	} = $props();

	let videoEl: HTMLVideoElement | undefined = $state();
	let playing = $state(false);

	const hasSrc = $derived(!!src?.trim());

	async function handlePlay() {
		if (!videoEl) return;
		try {
			await videoEl.play();
			playing = true;
		} catch {
			// Autoplay policies / missing codecs — leave controls visible
			playing = false;
		}
	}

	function handlePause() {
		playing = false;
	}
</script>

{#if hasSrc}
	<div class="relative w-full overflow-hidden rounded-lg bg-muted {className}">
		<!-- Click-to-play only; no autoplay sound. playsinline for mobile. -->
		<!-- svelte-ignore a11y_media_has_caption -->
		<video
			bind:this={videoEl}
			src={src!}
			class="block w-full max-h-[min(50vh,28rem)] object-contain bg-black"
			playsinline
			preload="metadata"
			controls={playing}
			onpause={handlePause}
			onended={handlePause}
		></video>
		{#if !playing}
			<button
				type="button"
				class="absolute inset-0 flex flex-col items-center justify-center gap-2 bg-black/40 text-white transition-colors hover:bg-black/50 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
				onclick={handlePlay}
			>
				<span
					class="flex h-14 w-14 items-center justify-center rounded-full bg-white/95 text-primary shadow-md"
					aria-hidden="true"
				>
					<svg class="ml-1 h-6 w-6" viewBox="0 0 24 24" fill="currentColor">
						<path d="M8 5v14l11-7z" />
					</svg>
				</span>
				<span class="text-sm font-medium tracking-wide">{label}</span>
			</button>
		{/if}
	</div>
{/if}
