<script lang="ts">
	/**
	 * Thumbs — outlined thumbs up/down toggle.
	 * Down on the left, up on the right. Outlined SVGs that fill on selection.
	 * Green fill for up, red fill for down.
	 */
	type LikeValue = 'THUMBS_UP' | 'THUMBS_DOWN' | null;

	let {
		value,
		onChange,
		size = 'md',
	}: {
		value: LikeValue;
		onChange: (val: LikeValue) => void;
		size?: 'md' | 'lg';
	} = $props();

	const dim = $derived(size === 'lg' ? 28 : 22);
	const gap = $derived(size === 'lg' ? 'gap-3.5' : 'gap-2.5');
	const isUp = $derived(value === 'THUMBS_UP');
	const isDown = $derived(value === 'THUMBS_DOWN');
</script>

<div class="flex items-center justify-center {gap}">
	<!-- Thumbs down (left) -->
	<button
		type="button"
		onclick={() => onChange(isDown ? null : 'THUMBS_DOWN')}
		aria-pressed={isDown}
		aria-label="Thumbs down"
		class="flex items-center justify-center p-1.5 rounded-lg border-none bg-transparent cursor-pointer transition-transform duration-[120ms] ease-out"
		style="transform: {isDown ? 'scale(1.05)' : 'scale(1)'};"
	>
		<svg
			width={dim}
			height={dim}
			viewBox="0 0 24 24"
			fill={isDown ? '#b91c1c' : 'none'}
			stroke={isDown ? '#b91c1c' : 'var(--color-muted-foreground)'}
			stroke-width={isDown ? 1.6 : 2}
			stroke-linecap="round"
			stroke-linejoin="round"
		>
			<path d="M17 14V2" />
			<path d="M9 18.12 10 14H4.17a2 2 0 0 1-1.92-2.56l2.33-8A2 2 0 0 1 6.5 2H17v12l-4.34 9.66a1 1 0 0 1-1.66-.43z" />
		</svg>
	</button>

	<!-- Thumbs up (right) -->
	<button
		type="button"
		onclick={() => onChange(isUp ? null : 'THUMBS_UP')}
		aria-pressed={isUp}
		aria-label="Thumbs up"
		class="flex items-center justify-center p-1.5 rounded-lg border-none bg-transparent cursor-pointer transition-transform duration-[120ms] ease-out"
		style="transform: {isUp ? 'scale(1.05)' : 'scale(1)'};"
	>
		<svg
			width={dim}
			height={dim}
			viewBox="0 0 24 24"
			fill={isUp ? '#15803d' : 'none'}
			stroke={isUp ? '#15803d' : 'var(--color-muted-foreground)'}
			stroke-width={isUp ? 1.6 : 2}
			stroke-linecap="round"
			stroke-linejoin="round"
		>
			<path d="M7 10v12" />
			<path d="M15 5.88 14 10h5.83a2 2 0 0 1 1.92 2.56l-2.33 8A2 2 0 0 1 17.5 22H7V10l4.34-9.66a1 1 0 0 1 1.66.43z" />
		</svg>
	</button>
</div>
