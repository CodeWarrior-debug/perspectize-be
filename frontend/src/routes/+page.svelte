<script lang="ts">
	import ActivityTable from '$lib/components/ActivityTable.svelte';
</script>

<div class="flex flex-col h-[calc(100vh-4rem)]">
	<div class="px-4 md:px-6 lg:px-8 py-6">
		<h1 class="text-3xl font-bold">Activity</h1>
	</div>
	<div class="flex-1 min-h-0 px-4 md:px-6 lg:px-8 pb-4">
		<ActivityTable />
	</div>
</div>
