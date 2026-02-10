<script lang="ts">
	import { createMutation, useQueryClient } from '@tanstack/svelte-query';
	import { toast } from 'svelte-sonner';
	import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription, DialogFooter, Button, Input, Label } from '$lib/components/shadcn';
	import { graphqlClient } from '$lib/queries/client';
	import { CREATE_CONTENT_FROM_YOUTUBE } from '$lib/queries/content';
	import { validateYouTubeUrl } from '$lib/utils/youtube';

	// Props
	let { open = $bindable(false) } = $props();

	// Reactive state
	let url = $state('');
	let error = $state('');

	// Reset form when dialog opens
	$effect(() => {
		if (open) {
			url = '';
			error = '';
		}
	});

	// Query client for cache invalidation
	const queryClient = useQueryClient();

	// TanStack Query mutation
	const mutation = createMutation(() => ({
		mutationFn: async (videoUrl: string) => {
			return graphqlClient.request(CREATE_CONTENT_FROM_YOUTUBE, {
				input: { url: videoUrl }
			});
		},
		onSuccess: (data: any) => {
			const name = data?.createContentFromYouTube?.name ?? 'video';
			toast.success(`Added: ${name}`);
			queryClient.invalidateQueries({ queryKey: ['content'] });
			open = false;
		},
		onError: (err: Error) => {
			const message = err.message.toLowerCase();
			if (message.includes('already exists')) {
				toast.error('This video has already been added');
			} else if (message.includes('invalid youtube url') || message.includes('video not found')) {
				toast.error('Invalid YouTube URL or video not found');
			} else {
				toast.error('Failed to add video. Please try again.');
			}
		}
	}));

	// Form submission handler
	function handleSubmit(e: Event) {
		e.preventDefault();

		if (!validateYouTubeUrl(url)) {
			error = 'Please enter a valid YouTube URL';
			return;
		}

		error = '';
		mutation.mutate(url);
	}
</script>

<Dialog bind:open>
	<DialogContent>
		<DialogHeader>
			<DialogTitle>Add Video</DialogTitle>
			<DialogDescription>Paste a YouTube URL to add it to your library.</DialogDescription>
		</DialogHeader>

		<form onsubmit={handleSubmit}>
			<div class="space-y-4 py-4">
				<div class="space-y-2">
					<Label for="url">YouTube URL</Label>
					<Input
						id="url"
						type="text"
						placeholder="https://www.youtube.com/watch?v=..."
						bind:value={url}
						disabled={mutation.isPending}
					/>
					{#if error}
						<p class="text-sm text-red-600">{error}</p>
					{/if}
				</div>
			</div>

			<DialogFooter>
				<Button
					type="button"
					variant="outline"
					onclick={() => (open = false)}
					disabled={mutation.isPending}
				>
					Cancel
				</Button>
				<Button type="submit" disabled={mutation.isPending || !url.trim()}>
					{mutation.isPending ? 'Adding...' : 'Add Video'}
				</Button>
			</DialogFooter>
		</form>
	</DialogContent>
</Dialog>
