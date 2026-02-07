<script lang="ts">
	import { createQuery } from '@tanstack/svelte-query';
	import { graphqlClient } from '$lib/queries/client';
	import { LIST_USERS } from '$lib/queries/users';
	import { setSelectedUserId, getSelectedUserId } from '$lib/stores/userSelection.svelte';

	interface User {
		id: string;
		username: string;
		email: string;
	}

	interface UsersResponse {
		users: User[];
	}

	const usersQuery = createQuery<UsersResponse>({
		queryKey: ['users'],
		queryFn: () => graphqlClient.request(LIST_USERS),
		staleTime: 5 * 60 * 1000, // 5 minutes
	});

	function handleChange(event: Event) {
		const target = event.target as HTMLSelectElement;
		const value = target.value;
		setSelectedUserId(value ? parseInt(value, 10) : null);
	}

	const currentUserId = $derived(getSelectedUserId());
</script>

<div>
	{#if $usersQuery.isLoading}
		<select class="h-9 rounded-md border border-input bg-background px-3 text-sm" disabled>
			<option>Loading users...</option>
		</select>
	{:else if $usersQuery.error}
		<select class="h-9 rounded-md border border-input bg-background px-3 text-sm text-destructive" disabled>
			<option>Error loading users</option>
		</select>
	{:else if $usersQuery.data}
		<select
			class="h-9 rounded-md border border-input bg-background px-3 text-sm"
			value={currentUserId ? String(currentUserId) : ''}
			onchange={handleChange}
		>
			<option value="">Select user...</option>
			{#each $usersQuery.data.users as user}
				<option value={user.id}>{user.username}</option>
			{/each}
		</select>
	{/if}
</div>
