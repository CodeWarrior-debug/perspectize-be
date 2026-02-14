<script lang="ts">
	import { Popover, PopoverContent, PopoverTrigger, buttonVariants, Button, Input, Label } from '$lib/components/shadcn';
	import { useAddVideo } from '$lib/queries/hooks/useAddVideo';
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

	// Shared mutation hook
	const mutation = useAddVideo();

	// Close popover on success
	$effect(() => {
		if (mutation.isSuccess) {
			open = false;
		}
	});

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
