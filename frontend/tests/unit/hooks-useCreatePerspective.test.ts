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

describe('useCreatePerspective hook', () => {
	beforeEach(async () => {
		vi.clearAllMocks();
		capturedMutationOptions = undefined;
		const { useCreatePerspective } = await import('$lib/queries/hooks/useCreatePerspective');
		useCreatePerspective();
	});

	describe('hook initialization', () => {
		it('returns a mutation object with mutate method', async () => {
			const { useCreatePerspective } = await import('$lib/queries/hooks/useCreatePerspective');
			const mutation = useCreatePerspective();
			expect(mutation).toBeDefined();
			expect(mutation.mutate).toBeDefined();
		});
	});

	describe('mutationFn', () => {
		it('is defined and callable', () => {
			expect(capturedMutationOptions).toBeDefined();
			expect(typeof capturedMutationOptions.mutationFn).toBe('function');
		});

		it('calls graphqlRequest with CREATE_PERSPECTIVE and input', async () => {
			const { graphqlRequest } = await import('$lib/queries/client');
			const { CREATE_PERSPECTIVE } = await import('$lib/queries/perspectives');
			(graphqlRequest as any).mockResolvedValue({
				createPerspective: { id: '1', userID: '42', quality: 7500, privacy: 'PUBLIC', createdAt: '', updatedAt: '' },
			});

			const input = { userID: 42, contentID: 10, quality: 7500 };
			await capturedMutationOptions.mutationFn(input);

			expect(graphqlRequest).toHaveBeenCalledWith(CREATE_PERSPECTIVE, { input });
		});

		it('passes partial input (only quality)', async () => {
			const { graphqlRequest } = await import('$lib/queries/client');
			(graphqlRequest as any).mockResolvedValue({
				createPerspective: { id: '1', userID: '42', quality: 5000, privacy: 'PUBLIC', createdAt: '', updatedAt: '' },
			});

			const input = { userID: 42, quality: 5000 };
			await capturedMutationOptions.mutationFn(input);

			expect(graphqlRequest).toHaveBeenCalledWith(expect.anything(), { input });
		});
	});

	describe('onSuccess callback', () => {
		it('shows success toast "Perspective added"', () => {
			capturedMutationOptions.onSuccess();
			expect(mockToastSuccess).toHaveBeenCalledWith('Perspective added');
		});

		it('invalidates perspectives lists cache', () => {
			capturedMutationOptions.onSuccess();
			expect(mockInvalidateQueries).toHaveBeenCalledWith(
				expect.objectContaining({ queryKey: expect.arrayContaining(['perspectives', 'list']) })
			);
		});

		it('does not invalidate content lists cache (content unchanged by perspective)', () => {
			capturedMutationOptions.onSuccess();
			expect(mockInvalidateQueries).not.toHaveBeenCalledWith(
				expect.objectContaining({ queryKey: expect.arrayContaining(['content', 'list']) })
			);
		});
	});

	describe('onError callback', () => {
		it('shows "No user selected" message for user not found errors', () => {
			capturedMutationOptions.onError(new Error('user not found'));
			expect(mockToastError).toHaveBeenCalledWith('No user selected');
		});

		it('shows "Invalid rating" message for rating errors', () => {
			capturedMutationOptions.onError(new Error('invalid rating: quality 99999'));
			expect(mockToastError).toHaveBeenCalledWith('Invalid rating value');
		});

		it('shows "at least one field" message for empty-input errors', () => {
			capturedMutationOptions.onError(new Error('at least one field must be provided'));
			expect(mockToastError).toHaveBeenCalledWith('Please fill in at least one field');
		});

		it('shows generic message for unknown errors', () => {
			capturedMutationOptions.onError(new Error('connection timeout'));
			expect(mockToastError).toHaveBeenCalledWith('Failed to add perspective. Please try again.');
		});
	});
});
