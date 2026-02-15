<script lang="ts">
	import { browser } from '$app/environment';
	import { onMount } from 'svelte';
	import { QueryClient, QueryClientProvider } from '@tanstack/svelte-query';
	import { Toaster } from 'svelte-sonner';
	import favicon from '$lib/assets/favicon.svg';
	import Header from '$lib/components/Header.svelte';
	import { reportWebVitals } from '$lib/vitals';
	import '../app.css';

	// CRITICAL: Disable queries on server to prevent post-SSR execution
	const queryClient = new QueryClient({
		defaultOptions: {
			queries: {
				enabled: browser, // Only run queries in browser
				staleTime: 60 * 1000, // 1 minute
				retry: 1
			}
		}
	});

	let { children } = $props();

	onMount(() => {
		reportWebVitals();
	});
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<QueryClientProvider client={queryClient}>
	<!-- Toast notifications: top-right, 2s auto-dismiss per requirements -->
	<Toaster position="top-right" duration={2000} richColors />

	<div class="min-h-screen bg-background text-foreground">
		<Header />
		{@render children()}
	</div>
</QueryClientProvider>
