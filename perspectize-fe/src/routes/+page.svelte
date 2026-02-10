<script lang="ts">
	import { createQuery } from '@tanstack/svelte-query';
	import { graphqlClient } from '$lib/queries/client';
	import { LIST_CONTENT } from '$lib/queries/content';
	import PageWrapper from '$lib/components/PageWrapper.svelte';
	import ActivityTable from '$lib/components/ActivityTable.svelte';

	interface ContentItem {
		id: string;
		name: string;
		url: string | null;
		contentType: string;
		length: number | null;
		lengthUnits: string | null;
		createdAt: string;
		updatedAt: string;
	}

	interface ContentResponse {
		content: {
			items: ContentItem[];
			pageInfo: {
				hasNextPage: boolean;
				endCursor: string | null;
			};
			totalCount: number;
		};
	}

	let searchText = $state('');
	let debouncedSearchText = $state('');
	let debounceTimer: ReturnType<typeof setTimeout>;

	$effect(() => {
		clearTimeout(debounceTimer);
		debounceTimer = setTimeout(() => {
			debouncedSearchText = searchText;
		}, 300);
		return () => clearTimeout(debounceTimer);
	});

	const contentQuery = createQuery(() => ({
		queryKey: ['content', { first: 100, sortBy: 'UPDATED_AT', sortOrder: 'DESC' }],
		queryFn: () => graphqlClient.request(LIST_CONTENT, {
			first: 100,
			sortBy: 'UPDATED_AT',
			sortOrder: 'DESC'
		}),
		staleTime: 60 * 1000, // 1 minute
		retry: 1
	}));

	const rowData = $derived(contentQuery.data?.content.items ?? []);
</script>

<PageWrapper>
	<div class="space-y-6">
		<div>
			<h1 class="text-3xl font-bold mb-2">Activity</h1>
			<p class="text-muted-foreground">
				Recently updated content
			</p>
		</div>

		<!-- Search input -->
		<div>
			<input
				type="text"
				bind:value={searchText}
				placeholder="Search content..."
				class="w-full max-w-md px-4 py-2 border border-input rounded-md bg-background text-sm"
			/>
		</div>

		<!-- Content Table -->
		{#if contentQuery.isLoading}
			<div class="py-12 text-center text-muted-foreground">
				<p>Loading content...</p>
			</div>
		{:else if contentQuery.error}
			<div class="py-12 text-center text-destructive">
				<p>Error loading content: {contentQuery.error.message}</p>
				<button
					onclick={() => contentQuery.refetch()}
					class="mt-4 px-4 py-2 text-sm rounded-md border border-input bg-background hover:bg-accent"
				>
					Retry
				</button>
			</div>
		{:else if rowData.length === 0}
			<div class="py-12 text-center text-muted-foreground">
				<p>No content found</p>
			</div>
		{:else}
			<ActivityTable {rowData} loading={contentQuery.isLoading} searchText={debouncedSearchText} />
		{/if}
	</div>
</PageWrapper>
