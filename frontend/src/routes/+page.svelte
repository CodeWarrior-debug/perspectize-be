<script lang="ts">
	import ActivityTable from '$lib/components/ActivityTable.svelte';
	import { Input } from '$lib/components/shadcn';
	import SearchIcon from '@lucide/svelte/icons/search';

	let searchText = $state('');
</script>

<div class="flex flex-col h-[calc(100vh-4rem)]">
	<!-- Page Header: Title + Search -->
	<div class="px-4 md:px-6 lg:px-8 py-4 md:py-6">
		<div class="flex flex-col sm:flex-row sm:items-end sm:justify-between gap-3">
			<div>
				<h1 class="text-2xl md:text-3xl font-semibold text-foreground">Activity</h1>
				<p class="text-sm text-muted-foreground mt-1">Recently updated content</p>
			</div>
			<div class="relative w-full sm:w-64 md:w-80">
				<SearchIcon class="absolute left-3 top-1/2 -translate-y-1/2 size-4 text-muted-foreground pointer-events-none" />
				<Input type="text" placeholder="Search content..." bind:value={searchText} class="pl-9" />
			</div>
		</div>
	</div>

	<!-- Table Card -->
	<div class="flex-1 min-h-0 px-4 md:px-6 lg:px-8 pb-4">
		<div class="border rounded-lg shadow-sm overflow-hidden h-full flex flex-col">
			<ActivityTable {searchText} />
		</div>
	</div>
</div>
