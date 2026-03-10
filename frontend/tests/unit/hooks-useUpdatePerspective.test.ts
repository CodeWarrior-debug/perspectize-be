import { describe, it, expect, vi, beforeEach } from 'vitest';

const { mockMutate, mockInvalidateQueries, mockToastSuccess, mockToastError } = vi.hoisted(() => ({
	mockMutate: vi.fn(),
	mockInvalidateQueries: vi.fn(),
	mockToastSuccess: vi.fn(),
	mockToastError: vi.fn(),
}));

let capturedMutationOptions: any;

vi.mock('@tanstack/svelte-query', () => ({
	createMutation: vi.fn((optionsFn: () => any) => {
		capturedMutationOptions = optionsFn();
		return {
			mutate: mockMutate,
			isPending: false,
		};
	}),
	useQueryClient: vi.fn(() => ({
		invalidateQueries: mockInvalidateQueries,
	})),
}));

vi.mock('svelte-sonner', () => ({
	toast: {
		success: mockToastSuccess,
		error: mockToastError,
	},
}));

vi.mock('$lib/queries/client', () => ({
	graphqlRequest: vi.fn(),
}));

describe('useUpdatePerspective hook', () => {
	beforeEach(async () => {
		vi.clearAllMocks();
		capturedMutationOptions = undefined;
		const { useUpdatePerspective } = await import('$lib/queries/hooks/useUpdatePerspective');
		useUpdatePerspective();
	});

	describe('hook initialization', () => {
		it('returns a mutation object with mutate method', async () => {
			const { useUpdatePerspective } = await import('$lib/queries/hooks/useUpdatePerspective');
			const mutation = useUpdatePerspective();
			expect(mutation).toBeDefined();
			expect(mutation.mutate).toBeDefined();
		});
	});

	describe('mutationFn', () => {
		it('is defined and callable', () => {
			expect(capturedMutationOptions).toBeDefined();
			expect(typeof capturedMutationOptions.mutationFn).toBe('function');
		});

		it('calls graphqlRequest with UPDATE_PERSPECTIVE and input', async () => {
			const { graphqlRequest } = await import('$lib/queries/client');
			const { UPDATE_PERSPECTIVE } = await import('$lib/queries/perspectives');
			(graphqlRequest as any).mockResolvedValue({
				updatePerspective: { id: '1', userID: '42', quality: 8000, privacy: 'PUBLIC', createdAt: '', updatedAt: '' },
			});

			const input = { id: 1, quality: 8000 };
			await capturedMutationOptions.mutationFn(input);

			expect(graphqlRequest).toHaveBeenCalledWith(UPDATE_PERSPECTIVE, { input });
		});

		it('supports partial update with only quality', async () => {
			const { graphqlRequest } = await import('$lib/queries/client');
			(graphqlRequest as any).mockResolvedValue({
				updatePerspective: { id: '5', userID: '42', quality: 6000, privacy: 'PUBLIC', createdAt: '', updatedAt: '' },
			});

			const input = { id: 5, quality: 6000 };
			await capturedMutationOptions.mutationFn(input);

			expect(graphqlRequest).toHaveBeenCalledWith(expect.anything(), { input: { id: 5, quality: 6000 } });
		});

		it('supports partial update with only like', async () => {
			const { graphqlRequest } = await import('$lib/queries/client');
			(graphqlRequest as any).mockResolvedValue({
				updatePerspective: { id: '3', userID: '42', like: 'THUMBS_UP', privacy: 'PUBLIC', createdAt: '', updatedAt: '' },
			});

			const input = { id: 3, like: 'THUMBS_UP' };
			await capturedMutationOptions.mutationFn(input);

			expect(graphqlRequest).toHaveBeenCalledWith(expect.anything(), { input: { id: 3, like: 'THUMBS_UP' } });
		});
	});

	describe('onSuccess callback', () => {
		it('shows success toast "Perspective updated"', () => {
			capturedMutationOptions.onSuccess();
			expect(mockToastSuccess).toHaveBeenCalledWith('Perspective updated');
		});

		it('invalidates perspectives lists cache', () => {
			capturedMutationOptions.onSuccess();
			expect(mockInvalidateQueries).toHaveBeenCalledWith(
				expect.objectContaining({ queryKey: expect.arrayContaining(['perspectives', 'list']) })
			);
		});

		it('does not invalidate content cache on update', () => {
			capturedMutationOptions.onSuccess();
			// Update only needs to invalidate perspectives, not content
			const calls = mockInvalidateQueries.mock.calls;
			const contentCall = calls.find((call: any[]) =>
				call[0]?.queryKey?.includes('content')
			);
			expect(contentCall).toBeUndefined();
		});
	});

	describe('onError callback', () => {
		it('shows "Invalid rating" message for rating errors', () => {
			capturedMutationOptions.onError(new Error('invalid rating: quality 99999'));
			expect(mockToastError).toHaveBeenCalledWith('Invalid rating value');
		});

		it('shows "at least one field" message for empty-input errors', () => {
			capturedMutationOptions.onError(new Error('at least one field must be provided'));
			expect(mockToastError).toHaveBeenCalledWith('Please fill in at least one field');
		});

		it('shows generic message for unknown errors', () => {
			capturedMutationOptions.onError(new Error('network failure'));
			expect(mockToastError).toHaveBeenCalledWith('Failed to update perspective. Please try again.');
		});

		it('shows "No user selected" for user-related errors', () => {
			capturedMutationOptions.onError(new Error('user not found'));
			expect(mockToastError).toHaveBeenCalledWith('No user selected');
		});
	});
});
