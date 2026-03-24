import { createMutation, useQueryClient } from '@tanstack/svelte-query';
import { toast } from 'svelte-sonner';
import { graphqlRequest } from '../client';
import { UPDATE_PERSPECTIVE, type UpdatePerspectiveResponse } from '../perspectives';
import { queryKeys } from '../keys';

export interface UpdatePerspectiveInput {
	id: number;
	quality?: number;
	agreement?: number;
	importance?: number;
	confidence?: number;
	like?: string;
	review?: string;
}

export function useUpdatePerspective() {
	const queryClient = useQueryClient();

	return createMutation(() => ({
		mutationFn: async (input: UpdatePerspectiveInput) => {
			return graphqlRequest<UpdatePerspectiveResponse>(UPDATE_PERSPECTIVE, { input });
		},
		onSuccess: () => {
			toast.success('Perspective updated');
			queryClient.invalidateQueries({ queryKey: queryKeys.perspectives.lists() });
		},
		onError: (err: Error) => {
			const message = err.message.toLowerCase();
			if (message.includes('no user selected') || message.includes('user not found')) {
				toast.error('No user selected');
			} else if (message.includes('invalid rating')) {
				toast.error('Invalid rating value');
			} else if (message.includes('at least one field')) {
				toast.error('Please fill in at least one field');
			} else {
				toast.error('Failed to update perspective. Please try again.');
			}
		},
	}));
}
