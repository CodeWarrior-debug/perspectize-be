<script lang="ts">
	import { Select, SelectTrigger, SelectContent, SelectItem } from '$lib/components/shadcn';
	import type { SearchFilters } from '$lib/services/youtubeApi';

	let { filters = $bindable() }: { filters: SearchFilters } = $props();

	const DURATION_OPTIONS = [
		{ value: 'any', label: 'Any duration' },
		{ value: 'short', label: 'Under 4 minutes' },
		{ value: 'medium', label: '4-20 minutes' },
		{ value: 'long', label: 'Over 20 minutes' },
	] as const;

	const DATE_OPTIONS = [
		{ value: 'any', label: 'Any time' },
		{ value: 'hour', label: 'Last hour' },
		{ value: 'today', label: 'Today' },
		{ value: 'week', label: 'This week' },
		{ value: 'month', label: 'This month' },
		{ value: 'year', label: 'This year' },
	] as const;

	const ORDER_OPTIONS = [
		{ value: 'relevance', label: 'Relevance' },
		{ value: 'date', label: 'Upload date' },
		{ value: 'viewCount', label: 'View count' },
		{ value: 'rating', label: 'Rating' },
	] as const;

	type DateOptionValue = (typeof DATE_OPTIONS)[number]['value'];

	/**
	 * The Upload Date select needs a stable relative option ('week', 'month', ...)
	 * to display, but `filters.publishedAfter` is an absolute ISO instant computed
	 * at selection time. Track the relative option separately and derive
	 * publishedAfter from it whenever it changes.
	 */
	let dateOption = $state<DateOptionValue>('any');

	const durationValue = $derived(filters.videoDuration ?? 'any');
	const orderValue = $derived(filters.order ?? 'relevance');

	const durationLabel = $derived(DURATION_OPTIONS.find((o) => o.value === durationValue)?.label);
	const dateLabel = $derived(DATE_OPTIONS.find((o) => o.value === dateOption)?.label);
	const orderLabel = $derived(ORDER_OPTIONS.find((o) => o.value === orderValue)?.label);

	const hasActiveFilters = $derived(
		filters.videoDuration !== undefined || filters.publishedAfter !== undefined || orderValue !== 'relevance',
	);

	function handleDurationChange(value: string | undefined) {
		filters = {
			...filters,
			videoDuration: value && value !== 'any' ? (value as SearchFilters['videoDuration']) : undefined,
		};
	}

	function calculatePublishedAfter(option: DateOptionValue): string | undefined {
		if (option === 'any') return undefined;
		const now = new Date();
		switch (option) {
			case 'hour':
				now.setHours(now.getHours() - 1);
				break;
			case 'today':
				now.setDate(now.getDate() - 1);
				break;
			case 'week':
				now.setDate(now.getDate() - 7);
				break;
			case 'month':
				now.setMonth(now.getMonth() - 1);
				break;
			case 'year':
				now.setFullYear(now.getFullYear() - 1);
				break;
		}
		return now.toISOString();
	}

	function handleDateChange(value: string | undefined) {
		dateOption = (value as DateOptionValue | undefined) ?? 'any';
		filters = { ...filters, publishedAfter: calculatePublishedAfter(dateOption) };
	}

	function handleOrderChange(value: string | undefined) {
		filters = { ...filters, order: (value as SearchFilters['order']) ?? 'relevance' };
	}

	function handleClear() {
		dateOption = 'any';
		filters = { videoDuration: undefined, publishedAfter: undefined, order: 'relevance' };
	}
</script>

<div class="flex flex-wrap items-center gap-2">
	<Select type="single" value={durationValue} onValueChange={handleDurationChange}>
		<SelectTrigger class="w-44" aria-label="Duration">
			{durationLabel}
		</SelectTrigger>
		<SelectContent>
			{#each DURATION_OPTIONS as option (option.value)}
				<SelectItem value={option.value}>{option.label}</SelectItem>
			{/each}
		</SelectContent>
	</Select>

	<Select type="single" value={dateOption} onValueChange={handleDateChange}>
		<SelectTrigger class="w-36" aria-label="Upload date">
			{dateLabel}
		</SelectTrigger>
		<SelectContent>
			{#each DATE_OPTIONS as option (option.value)}
				<SelectItem value={option.value}>{option.label}</SelectItem>
			{/each}
		</SelectContent>
	</Select>

	<Select type="single" value={orderValue} onValueChange={handleOrderChange}>
		<SelectTrigger class="w-36" aria-label="Sort order">
			{orderLabel}
		</SelectTrigger>
		<SelectContent>
			{#each ORDER_OPTIONS as option (option.value)}
				<SelectItem value={option.value}>{option.label}</SelectItem>
			{/each}
		</SelectContent>
	</Select>

	{#if hasActiveFilters}
		<button
			type="button"
			onclick={handleClear}
			class="text-sm text-muted-foreground hover:text-foreground underline underline-offset-2"
		>
			Clear Filters
		</button>
	{/if}
</div>
