<script lang="ts">
	import { QueryClientProvider, type QueryClient } from '@tanstack/svelte-query';
	import type { Component } from 'svelte';

	let { queryClient, component, props = {} } = $props<{
		queryClient: QueryClient;
		component: Component;
		props?: any;
	}>();
</script>

<QueryClientProvider client={queryClient}>
	<component.default {...props} />
</QueryClientProvider>
