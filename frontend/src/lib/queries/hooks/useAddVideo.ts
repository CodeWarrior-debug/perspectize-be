import { createMutation, useQueryClient } from '@tanstack/svelte-query';
import { toast } from 'svelte-sonner';
import { graphqlClient } from '../client';
import { CREATE_CONTENT_FROM_YOUTUBE, type CreateContentResponse } from '../content';
import { queryKeys } from '../keys';
import { getSelectedUserId } from '$lib/stores/userSelection.svelte';

export function useAddVideo() {
	const queryClient = useQueryClient();

	return createMutation(() => ({
		mutationFn: async (url: string) => {
			const userId = getSelectedUserId();
			if (userId === null) {
				throw new Error('No user selected');
			}
			return graphqlClient.request<CreateContentResponse>(CREATE_CONTENT_FROM_YOUTUBE, {
				input: { url, userId }
			});
		},
		onSuccess: (data: CreateContentResponse) => {
			const name = data?.createContentFromYouTube?.name ?? 'video';
			toast.success(`Added: ${name}`);
			queryClient.invalidateQueries({ queryKey: queryKeys.content.lists() });
		},
		onError: (err: Error) => {
			const message = err.message.toLowerCase();
			if (message.includes('no user selected')) {
				toast.error('Please select a user first');
			} else if (message.includes('already exists')) {
				toast.error('This video has already been added');
			} else if (message.includes('invalid youtube url') || message.includes('video not found')) {
				toast.error('Invalid YouTube URL or video not found');
			} else {
				toast.error('Failed to add video. Please try again.');
			}
		}
	}));
}
