import { createMutation, useQueryClient } from '@tanstack/svelte-query';
import { toast } from 'svelte-sonner';
import { graphqlClient } from '../client';
import { CREATE_CONTENT_FROM_YOUTUBE, type CreateContentResponse, type ContentResponse } from '../content';
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
			const newItem = data?.createContentFromYouTube;
			toast.success(`Added: ${newItem?.name ?? 'video'}`);

			if (newItem) {
				// Insert new item at top of all active list caches
				queryClient.setQueriesData<ContentResponse>(
					{ queryKey: queryKeys.content.lists() },
					(oldData) => {
						if (!oldData) return oldData;
						return {
							content: {
								...oldData.content,
								items: [newItem, ...oldData.content.items],
								totalCount: (oldData.content.totalCount ?? 0) + 1,
							},
						};
					}
				);

				// Mark stale for eventual consistency (no immediate refetch)
				queryClient.invalidateQueries({
					queryKey: queryKeys.content.lists(),
					refetchType: 'none',
				});
			} else {
				// Fallback: full refetch if response shape is unexpected
				queryClient.invalidateQueries({ queryKey: queryKeys.content.lists() });
			}
		},
		onError: (err: Error) => {
			console.error('[AddVideo] mutation failed:', err);
			const message = err.message.toLowerCase();
			if (message.includes('no user selected')) {
				toast.error('Please select a user first');
			} else if (message.includes('already exists')) {
				toast.error('This video has already been added');
			} else if (message.includes('invalid youtube url') || message.includes('video not found')) {
				toast.error('Invalid YouTube URL or video not found');
			} else if (message.includes('load failed') || message.includes('failed to fetch')) {
				toast.error('Cannot reach the server. Check your connection and try again.');
			} else {
				toast.error('Failed to add video. Please try again.');
			}
		}
	}));
}
