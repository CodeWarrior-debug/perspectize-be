import { describe, it, expect, vi, beforeEach } from 'vitest';

// Hoisted mocks — these are referenced inside vi.mock factories
const { mockMutate, mockInvalidateQueries, mockToastSuccess, mockToastError } =
	vi.hoisted(() => ({
		mockMutate: vi.fn(),
		mockInvalidateQueries: vi.fn(),
		mockToastSuccess: vi.fn(),
		mockToastError: vi.fn(),
	}));

// Capture mutation options for behavioral testing
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
	graphqlClient: {
		request: vi.fn(),
	},
}));

vi.mock('$lib/stores/userSelection.svelte', () => ({
	getSelectedUserId: vi.fn(() => 1),
}));

describe('useAddVideo hook', () => {
	beforeEach(() => {
		vi.clearAllMocks();
		capturedMutationOptions = undefined;
	});

	describe('hook initialization', () => {
		it('returns a mutation object with mutate method', async () => {
			const { useAddVideo } = await import('$lib/queries/hooks/useAddVideo');
			const mutation = useAddVideo();

			expect(mutation).toBeDefined();
			expect(mutation.mutate).toBeDefined();
			expect(typeof mutation.mutate).toBe('function');
		});

		it('calls createMutation with a function that returns options', async () => {
			const { createMutation } = await import('@tanstack/svelte-query');
			const { useAddVideo } = await import('$lib/queries/hooks/useAddVideo');

			useAddVideo();

			expect(createMutation).toHaveBeenCalled();
			// Verify createMutation receives a function
			const call = (createMutation as any).mock.calls[0];
			expect(typeof call[0]).toBe('function');
		});

		it('captures mutation options via createMutation factory', async () => {
			const { useAddVideo } = await import('$lib/queries/hooks/useAddVideo');
			useAddVideo();

			expect(capturedMutationOptions).toBeDefined();
		});
	});

	describe('mutationFn', () => {
		beforeEach(async () => {
			const { useAddVideo } = await import('$lib/queries/hooks/useAddVideo');
			useAddVideo();
		});

		it('is defined and callable', () => {
			expect(capturedMutationOptions).toBeDefined();
			expect(capturedMutationOptions.mutationFn).toBeDefined();
			expect(typeof capturedMutationOptions.mutationFn).toBe('function');
		});

		it('accepts AddVideoParams object with url and optional userIdOverride', async () => {
			const { graphqlClient } = await import('$lib/queries/client');
			(graphqlClient.request as any).mockResolvedValue({
				createContentFromYouTube: { name: 'Test' }
			});

			const params = {
				url: 'https://youtube.com/watch?v=abc123',
				userIdOverride: undefined
			};

			await capturedMutationOptions.mutationFn(params);

			expect(graphqlClient.request).toHaveBeenCalled();
		});

		it('uses getSelectedUserId when userIdOverride is not provided', async () => {
			const { graphqlClient } = await import('$lib/queries/client');
			const { getSelectedUserId } = await import('$lib/stores/userSelection.svelte');
			(graphqlClient.request as any).mockResolvedValue({
				createContentFromYouTube: { name: 'Test' }
			});

			await capturedMutationOptions.mutationFn({
				url: 'https://youtube.com/watch?v=abc123',
				userIdOverride: undefined
			});

			expect(graphqlClient.request).toHaveBeenCalledWith(
				expect.anything(),
				{ input: { url: 'https://youtube.com/watch?v=abc123', userId: 1 } }
			);
		});

		it('uses userIdOverride when provided (anonymous case)', async () => {
			const { graphqlClient } = await import('$lib/queries/client');
			(graphqlClient.request as any).mockResolvedValue({
				createContentFromYouTube: { name: 'Test' }
			});

			const anonUserId = 999;
			await capturedMutationOptions.mutationFn({
				url: 'https://youtube.com/watch?v=abc123',
				userIdOverride: anonUserId
			});

			expect(graphqlClient.request).toHaveBeenCalledWith(
				expect.anything(),
				{ input: { url: 'https://youtube.com/watch?v=abc123', userId: anonUserId } }
			);
		});

		it('prefers userIdOverride over selected user ID', async () => {
			const { graphqlClient } = await import('$lib/queries/client');
			(graphqlClient.request as any).mockResolvedValue({
				createContentFromYouTube: { name: 'Test' }
			});

			// Selected user is 1, but override with 999
			await capturedMutationOptions.mutationFn({
				url: 'https://youtube.com/watch?v=abc123',
				userIdOverride: 999
			});

			expect(graphqlClient.request).toHaveBeenCalledWith(
				expect.anything(),
				{ input: { url: 'https://youtube.com/watch?v=abc123', userId: 999 } }
			);
		});

		it('calls graphqlClient.request with CREATE_CONTENT_FROM_YOUTUBE mutation', async () => {
			const { graphqlClient } = await import('$lib/queries/client');
			const { CREATE_CONTENT_FROM_YOUTUBE } = await import('$lib/queries/content');
			(graphqlClient.request as any).mockResolvedValue({
				createContentFromYouTube: { name: 'Test' }
			});

			await capturedMutationOptions.mutationFn({
				url: 'https://youtube.com/watch?v=abc123',
				userIdOverride: undefined
			});

			const firstArg = (graphqlClient.request as any).mock.calls[0][0];
			expect(firstArg).toBe(CREATE_CONTENT_FROM_YOUTUBE);
		});

		it('passes correct input variables to GraphQL mutation', async () => {
			const { graphqlClient } = await import('$lib/queries/client');
			(graphqlClient.request as any).mockResolvedValue({
				createContentFromYouTube: { name: 'Test Video' }
			});

			const testUrl = 'https://www.youtube.com/watch?v=dQw4w9WgXcQ';
			const testUserId = 42;

			await capturedMutationOptions.mutationFn({
				url: testUrl,
				userIdOverride: testUserId
			});

			expect(graphqlClient.request).toHaveBeenCalledWith(
				expect.anything(),
				{ input: { url: testUrl, userId: testUserId } }
			);
		});

		it('returns the response from graphqlClient.request', async () => {
			const { graphqlClient } = await import('$lib/queries/client');
			const mockResponse = {
				createContentFromYouTube: {
					id: '123',
					name: 'My Video',
					url: 'https://youtube.com/watch?v=test'
				}
			};
			(graphqlClient.request as any).mockResolvedValue(mockResponse);

			const result = await capturedMutationOptions.mutationFn({
				url: 'https://youtube.com/watch?v=test',
				userIdOverride: undefined
			});

			expect(result).toEqual(mockResponse);
		});

		it('throws "No user selected" error when no user is selected and no override', async () => {
			const { getSelectedUserId } = await import('$lib/stores/userSelection.svelte');
			(getSelectedUserId as any).mockReturnValueOnce(null);

			await expect(capturedMutationOptions.mutationFn({
				url: 'https://youtube.com/watch?v=abc123',
				userIdOverride: undefined
			})).rejects.toThrow('No user selected');
		});

		it('does not throw error when no selected user but userIdOverride is provided', async () => {
			const { graphqlClient } = await import('$lib/queries/client');
			const { getSelectedUserId } = await import('$lib/stores/userSelection.svelte');
			(getSelectedUserId as any).mockReturnValueOnce(null);
			(graphqlClient.request as any).mockResolvedValue({
				createContentFromYouTube: { name: 'Test' }
			});

			const result = await capturedMutationOptions.mutationFn({
				url: 'https://youtube.com/watch?v=abc123',
				userIdOverride: 999
			});

			expect(result).toBeDefined();
			expect(graphqlClient.request).toHaveBeenCalled();
		});

		it('handles network errors gracefully', async () => {
			const { graphqlClient } = await import('$lib/queries/client');
			const networkError = new Error('Network request failed');
			(graphqlClient.request as any).mockRejectedValue(networkError);

			await expect(capturedMutationOptions.mutationFn({
				url: 'https://youtube.com/watch?v=abc123',
				userIdOverride: 1
			})).rejects.toThrow('Network request failed');
		});

		it('handles GraphQL errors gracefully', async () => {
			const { graphqlClient } = await import('$lib/queries/client');
			const graphqlError = new Error('invalid YouTube URL');
			(graphqlClient.request as any).mockRejectedValue(graphqlError);

			await expect(capturedMutationOptions.mutationFn({
				url: 'https://invalid-url',
				userIdOverride: 1
			})).rejects.toThrow('invalid YouTube URL');
		});
	});

	describe('onSuccess callback', () => {
		beforeEach(async () => {
			const { useAddVideo } = await import('$lib/queries/hooks/useAddVideo');
			useAddVideo();
		});

		it('is defined and callable', () => {
			expect(capturedMutationOptions).toBeDefined();
			expect(capturedMutationOptions.onSuccess).toBeDefined();
			expect(typeof capturedMutationOptions.onSuccess).toBe('function');
		});

		it('shows success toast with video name from response', () => {
			const mockData = {
				createContentFromYouTube: {
					id: '1',
					name: 'Amazing YouTube Video',
					url: 'https://youtube.com/watch?v=test'
				}
			};

			capturedMutationOptions.onSuccess(mockData);

			expect(mockToastSuccess).toHaveBeenCalledWith('Added: Amazing YouTube Video');
		});

		it('shows success toast with default "video" text when name is missing', () => {
			const mockData = {
				createContentFromYouTube: {
					id: '1',
					url: 'https://youtube.com/watch?v=test'
				}
			};

			capturedMutationOptions.onSuccess(mockData);

			expect(mockToastSuccess).toHaveBeenCalledWith('Added: video');
		});

		it('shows success toast with default "video" text when response is null', () => {
			capturedMutationOptions.onSuccess(null);

			expect(mockToastSuccess).toHaveBeenCalledWith('Added: video');
		});

		it('invalidates content list queryKey after successful creation', () => {
			const mockData = {
				createContentFromYouTube: {
					id: '1',
					name: 'Test Video',
					url: 'https://youtube.com/watch?v=test'
				}
			};

			capturedMutationOptions.onSuccess(mockData);

			expect(mockInvalidateQueries).toHaveBeenCalledWith({
				queryKey: ['app', 'content', 'list']
			});
		});

		it('calls both toast and invalidateQueries on success', () => {
			const mockData = {
				createContentFromYouTube: {
					id: '1',
					name: 'Test Video',
					url: 'https://youtube.com/watch?v=test'
				}
			};

			capturedMutationOptions.onSuccess(mockData);

			expect(mockToastSuccess).toHaveBeenCalled();
			expect(mockInvalidateQueries).toHaveBeenCalled();
		});

		it('handles response with empty name string', () => {
			const mockData = {
				createContentFromYouTube: {
					id: '1',
					name: '',
					url: 'https://youtube.com/watch?v=test'
				}
			};

			capturedMutationOptions.onSuccess(mockData);

			// Empty string should be used as is (falsy check with ?? operator)
			expect(mockToastSuccess).toHaveBeenCalledWith('Added: ');
		});
	});

	describe('onError callback', () => {
		beforeEach(async () => {
			const { useAddVideo } = await import('$lib/queries/hooks/useAddVideo');
			useAddVideo();
		});

		it('is defined and callable', () => {
			expect(capturedMutationOptions).toBeDefined();
			expect(capturedMutationOptions.onError).toBeDefined();
			expect(typeof capturedMutationOptions.onError).toBe('function');
		});

		it('shows "no user selected" message for no user error', () => {
			const error = new Error('No user selected');

			capturedMutationOptions.onError(error);

			expect(mockToastError).toHaveBeenCalledWith('Please select a user first');
		});

		it('shows "already exists" message for duplicate content errors', () => {
			const error = new Error('content already exists for this URL');

			capturedMutationOptions.onError(error);

			expect(mockToastError).toHaveBeenCalledWith('This video has already been added');
		});

		it('matches "already exists" case-insensitively', () => {
			const error = new Error('Content ALREADY EXISTS in the system');

			capturedMutationOptions.onError(error);

			expect(mockToastError).toHaveBeenCalledWith('This video has already been added');
		});

		it('shows invalid URL message for "invalid youtube url" errors', () => {
			const error = new Error('invalid YouTube URL');

			capturedMutationOptions.onError(error);

			expect(mockToastError).toHaveBeenCalledWith('Invalid YouTube URL or video not found');
		});

		it('shows invalid URL message for "video not found" errors', () => {
			const error = new Error('video not found: abc123');

			capturedMutationOptions.onError(error);

			expect(mockToastError).toHaveBeenCalledWith('Invalid YouTube URL or video not found');
		});

		it('matches "invalid youtube url" case-insensitively', () => {
			const error = new Error('INVALID youtube url provided');

			capturedMutationOptions.onError(error);

			expect(mockToastError).toHaveBeenCalledWith('Invalid YouTube URL or video not found');
		});

		it('matches "video not found" case-insensitively', () => {
			const error = new Error('VIDEO NOT FOUND');

			capturedMutationOptions.onError(error);

			expect(mockToastError).toHaveBeenCalledWith('Invalid YouTube URL or video not found');
		});

		it('shows generic message for unknown errors', () => {
			const error = new Error('connection timeout');

			capturedMutationOptions.onError(error);

			expect(mockToastError).toHaveBeenCalledWith('Failed to add video. Please try again.');
		});

		it('shows generic message for validation errors', () => {
			const error = new Error('invalid input format');

			capturedMutationOptions.onError(error);

			expect(mockToastError).toHaveBeenCalledWith('Failed to add video. Please try again.');
		});

		it('shows generic message for network errors', () => {
			const error = new Error('network error occurred');

			capturedMutationOptions.onError(error);

			expect(mockToastError).toHaveBeenCalledWith('Failed to add video. Please try again.');
		});

		it('handles error messages with extra whitespace', () => {
			const error = new Error('   no user selected   ');

			capturedMutationOptions.onError(error);

			expect(mockToastError).toHaveBeenCalledWith('Please select a user first');
		});

		it('only shows error toast on error (does not call invalidateQueries)', () => {
			const error = new Error('some error');

			capturedMutationOptions.onError(error);

			expect(mockToastError).toHaveBeenCalled();
			expect(mockInvalidateQueries).not.toHaveBeenCalled();
		});

		it('does not misclassify errors containing "invalid" in unrelated context', () => {
			const error = new Error('invalid connection state');

			capturedMutationOptions.onError(error);

			expect(mockToastError).toHaveBeenCalledWith('Failed to add video. Please try again.');
		});

		it('handles errors with multiple matching keywords (first match wins)', () => {
			// If error contains both "already exists" and "invalid", "already exists" should match first
			const error = new Error('already exists and invalid youtube url');

			capturedMutationOptions.onError(error);

			expect(mockToastError).toHaveBeenCalledWith('This video has already been added');
		});
	});

	describe('integration with queryClient and toasts', () => {
		beforeEach(async () => {
			const { useAddVideo } = await import('$lib/queries/hooks/useAddVideo');
			useAddVideo();
		});

		it('calls both toast and invalidateQueries on success', () => {
			const mockData = {
				createContentFromYouTube: {
					id: '1',
					name: 'Test Video',
					url: 'https://youtube.com/watch?v=test'
				}
			};

			capturedMutationOptions.onSuccess(mockData);

			expect(mockToastSuccess).toHaveBeenCalled();
			expect(mockInvalidateQueries).toHaveBeenCalled();
		});

		it('only shows error toast on error (does not call invalidateQueries)', () => {
			const error = new Error('some error');

			capturedMutationOptions.onError(error);

			expect(mockToastError).toHaveBeenCalled();
			expect(mockInvalidateQueries).not.toHaveBeenCalled();
		});

		it('useQueryClient is called to get queryClient instance', async () => {
			const { useQueryClient } = await import('@tanstack/svelte-query');

			expect(useQueryClient).toHaveBeenCalled();
		});
	});

	describe('AddVideoParams input type', () => {
		beforeEach(async () => {
			const { useAddVideo } = await import('$lib/queries/hooks/useAddVideo');
			useAddVideo();
		});

		it('accepts object with url and falls back to selected user', async () => {
			const { graphqlClient } = await import('$lib/queries/client');
			(graphqlClient.request as any).mockResolvedValue({
				createContentFromYouTube: { name: 'Test' }
			});

			const params = {
				url: 'https://youtube.com/watch?v=abc123',
				userIdOverride: 1
			};

			await capturedMutationOptions.mutationFn(params);

			expect(graphqlClient.request).toHaveBeenCalledWith(
				expect.anything(),
				{ input: { url: 'https://youtube.com/watch?v=abc123', userId: 1 } }
			);
		});

		it('accepts object with url and defined userIdOverride', async () => {
			const { graphqlClient } = await import('$lib/queries/client');
			(graphqlClient.request as any).mockResolvedValue({
				createContentFromYouTube: { name: 'Test' }
			});

			const params = {
				url: 'https://youtube.com/watch?v=abc123',
				userIdOverride: 999
			};

			await capturedMutationOptions.mutationFn(params);

			expect(graphqlClient.request).toHaveBeenCalledWith(
				expect.anything(),
				{ input: { url: 'https://youtube.com/watch?v=abc123', userId: 999 } }
			);
		});

		it('uses nullish coalescing to prefer override over selected user', async () => {
			const { graphqlClient } = await import('$lib/queries/client');
			(graphqlClient.request as any).mockResolvedValue({
				createContentFromYouTube: { name: 'Test' }
			});

			// Test with 0 override (falsy but valid)
			const params = {
				url: 'https://youtube.com/watch?v=abc123',
				userIdOverride: 0
			};

			await capturedMutationOptions.mutationFn(params);

			// 0 is falsy but should be preferred over selected user 1
			expect(graphqlClient.request).toHaveBeenCalledWith(
				expect.anything(),
				{ input: { url: 'https://youtube.com/watch?v=abc123', userId: 0 } }
			);
		});
	});
});
