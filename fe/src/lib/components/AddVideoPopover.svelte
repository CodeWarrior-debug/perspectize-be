<script lang="ts">
	import { createMutation, useQueryClient } from '@tanstack/svelte-query';
	import { toast } from 'svelte-sonner';
	import { Popover, PopoverContent, PopoverTrigger, buttonVariants, Button, Input, Label } from '$lib/components/shadcn';
	import { graphqlClient } from '$lib/queries/client';
	import { CREATE_CONTENT_FROM_YOUTUBE } from '$lib/queries/content';
	import { validateYouTubeUrl } from '$lib/utils/youtube';
	import PlusIcon from '@lucide/svelte/icons/plus';

	// Reactive state
	let open = $state(false);
	let url = $state('');
	let error = $state('');

	// Reset form when popover opens
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
			// Dispatch custom event for ActivityTable to refetch
			window.dispatchEvent(new CustomEvent('content-added'));
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

<Popover bind:open>
	<PopoverTrigger class={buttonVariants({ size: "default" })}>
		<PlusIcon class="size-4" />
		Add Video
	</PopoverTrigger>
	<PopoverContent align="end" sideOffset={8}>
		<form onsubmit={handleSubmit}>
			<div class="space-y-4">
				<div>
					<h3 class="font-semibold text-base">Add Video</h3>
					<p class="text-muted-foreground text-sm mt-1">
						Paste a YouTube URL to add it to your library.
					</p>
				</div>

				<div class="space-y-2">
					<Label for="url">YouTube URL</Label>
					<Input
						id="url"
						type="text"
						placeholder="https://www.youtube.com/watch?v=..."
						bind:value={url}
						disabled={mutation.isPending}
						autocomplete="off"
					/>
					{#if error}
						<p class="text-sm text-red-600">{error}</p>
					{/if}
				</div>

				<div class="flex gap-2 justify-end">
					<Button
						type="button"
						variant="outline"
						size="sm"
						onclick={() => (open = false)}
						disabled={mutation.isPending}
					>
						Cancel
					</Button>
					<Button type="submit" size="sm" disabled={mutation.isPending || !url.trim()}>
						{mutation.isPending ? 'Adding...' : 'Add Video'}
					</Button>
				</div>
			</div>
		</form>
	</PopoverContent>
</Popover>
