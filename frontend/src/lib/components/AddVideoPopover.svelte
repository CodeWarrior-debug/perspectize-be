<script lang="ts">
	import { Input, Label } from '$lib/components/shadcn';
	import FormPopover from '$lib/components/FormPopover.svelte';
	import { useAddVideo } from '$lib/queries/hooks/useAddVideo';
	import { validateYouTubeUrl } from '$lib/utils/youtube';
	import { createQuery } from '@tanstack/svelte-query';
	import { graphqlClient } from '$lib/queries/client';
	import { GET_USER_BY_USERNAME, type UserByUsernameResponse } from '$lib/queries/users';
	import { queryKeys } from '$lib/queries/keys';
	import PlusIcon from '@lucide/svelte/icons/plus';

	// Reactive state
	let open = $state(false);
	let url = $state('');
	let error = $state('');
	let anonymous = $state(false);

	// Fetch the [anonymous] sentinel user
	const anonymousUserQuery = createQuery(() => ({
		queryKey: [...queryKeys.users.all(), 'anonymous'],
		queryFn: () =>
			graphqlClient.request<UserByUsernameResponse>(GET_USER_BY_USERNAME, {
				username: '[anonymous]'
			}),
		staleTime: Infinity,
	}));

	const anonymousUserId = $derived(
		anonymousUserQuery.data?.userByUsername
			? parseInt(anonymousUserQuery.data.userByUsername.id, 10)
			: null
	);

	// Reset form when popover opens
	$effect(() => {
		if (open) {
			url = '';
			error = '';
			anonymous = false;
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
	function handleSubmit() {
		if (!validateYouTubeUrl(url)) {
			error = 'Please enter a valid YouTube URL';
			return;
		}

		if (anonymous && anonymousUserId === null) {
			error = 'Anonymous user not available';
			return;
		}

		error = '';
		mutation.mutate({
			url,
			userIdOverride: anonymous && anonymousUserId !== null ? anonymousUserId : undefined
		});
	}
</script>

<FormPopover
	bind:open
	triggerLabel="Add Video"
	title="Add Video"
	description="Paste a YouTube URL to add it to your library."
	submitLabel="Add Video"
	pendingLabel="Adding..."
	isPending={mutation.isPending}
	isSubmitDisabled={!url.trim()}
	onSubmit={handleSubmit}
>
	{#snippet triggerIcon()}
		<PlusIcon class="size-4" />
	{/snippet}
	{#snippet formFields()}
		<div class="space-y-3">
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
			<label class="flex items-center gap-2 text-sm cursor-pointer">
				<input
					type="checkbox"
					bind:checked={anonymous}
					disabled={mutation.isPending}
					class="rounded border-input"
				/>
				<span class="text-muted-foreground">Add anonymously</span>
			</label>
		</div>
	{/snippet}
</FormPopover>
