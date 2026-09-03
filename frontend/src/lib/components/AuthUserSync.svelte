<script lang="ts">
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import { useClerkContext } from 'svelte-clerk';
	import { graphqlRequest } from '$lib/queries/client';
	import { ME, type MeResponse } from '$lib/queries/users';
	import { setSelectedUserId, clearUserSelection } from '$lib/stores/userSelection.svelte';

	const queryClient = useQueryClient();
	const clerk = useClerkContext();

	const clerkUserId = $derived(clerk.auth.userId);

	const meQuery = createQuery(() => ({
		queryKey: ['me', clerkUserId],
		queryFn: () => graphqlRequest<MeResponse>(ME),
		enabled: clerk.isLoaded && !!clerkUserId,
		staleTime: 5 * 60 * 1000,
	}));

	// Tracks the last clerkUserId we reacted to, so we only clear/resync on an
	// actual transition (sign-in, sign-out, or account switch on a shared
	// device) — not on every unrelated re-render.
	let lastSyncedClerkUserId: string | null | undefined = undefined;

	$effect(() => {
		if (!clerk.isLoaded) return;
		const currentId = clerkUserId ?? null;
		if (currentId === lastSyncedClerkUserId) return;

		lastSyncedClerkUserId = currentId;
		// content.lists(), perspectives.listByUser(), etc. are all user-scoped —
		// clear everything, not just the me query, so nothing leaks across
		// accounts on a shared device.
		queryClient.clear();

		if (currentId === null) {
			clearUserSelection();
		}
	});

	$effect(() => {
		if (clerk.isLoaded && clerkUserId && meQuery.data?.me) {
			setSelectedUserId(parseInt(meQuery.data.me.id, 10));
		}
	});
</script>
