import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render } from '@testing-library/svelte';
import AddVideoDialog from '$lib/components/AddVideoDialog.svelte';

// Hoisted mocks — these are referenced inside vi.mock factories
const {
	mockMutate,
	mockInvalidateQueries,
	mockToastSuccess,
	mockToastError,
	mockValidate,
} = vi.hoisted(() => ({
	mockMutate: vi.fn(),
	mockInvalidateQueries: vi.fn(),
	mockToastSuccess: vi.fn(),
	mockToastError: vi.fn(),
	mockValidate: vi.fn(),
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

vi.mock('$lib/utils/youtube', () => ({
	validateYouTubeUrl: (...args: any[]) => mockValidate(...args),
}));


describe('AddVideoDialog component', () => {
	beforeEach(() => {
		vi.clearAllMocks();
		capturedMutationOptions = undefined;
	});

	it('renders without errors when open', () => {
		const result = render(AddVideoDialog, { props: { open: true } });
		expect(result.container).toBeTruthy();
	});

	it('renders without errors when closed', () => {
		const result = render(AddVideoDialog, { props: { open: false } });
		expect(result.container).toBeTruthy();
	});

	it('accepts open prop', () => {
		const result = render(AddVideoDialog, { props: { open: true } });
		expect(result).toBeTruthy();
	});

	it('accepts no props', () => {
		const result = render(AddVideoDialog);
		expect(result).toBeTruthy();
	});
});

describe('AddVideoDialog mutation callbacks', () => {
	beforeEach(() => {
		vi.clearAllMocks();
		capturedMutationOptions = undefined;
		render(AddVideoDialog, { props: { open: true } });
	});

	it('onSuccess shows toast with video name and invalidates cache', () => {
		expect(capturedMutationOptions).toBeDefined();

		capturedMutationOptions.onSuccess({
			createContentFromYouTube: { name: 'Test Video' }
		});

		expect(mockToastSuccess).toHaveBeenCalledWith('Added: Test Video');
		expect(mockInvalidateQueries).toHaveBeenCalledWith({ queryKey: ['app', 'content', 'list'] });
	});

	it('onSuccess handles null response gracefully', () => {
		expect(capturedMutationOptions).toBeDefined();

		capturedMutationOptions.onSuccess(null);

		expect(mockToastSuccess).toHaveBeenCalledWith('Added: video');
		expect(mockInvalidateQueries).toHaveBeenCalled();
	});

	it('onSuccess handles missing nested properties', () => {
		expect(capturedMutationOptions).toBeDefined();

		capturedMutationOptions.onSuccess({ createContentFromYouTube: {} });

		expect(mockToastSuccess).toHaveBeenCalledWith('Added: video');
	});

	it('onError shows duplicate message for "already exists" errors', () => {
		expect(capturedMutationOptions).toBeDefined();

		capturedMutationOptions.onError(new Error('content already exists for this URL'));

		expect(mockToastError).toHaveBeenCalledWith('This video has already been added');
	});

	it('onError shows invalid URL message for "invalid youtube url" errors', () => {
		expect(capturedMutationOptions).toBeDefined();

		capturedMutationOptions.onError(new Error('invalid YouTube URL'));

		expect(mockToastError).toHaveBeenCalledWith('Invalid YouTube URL or video not found');
	});

	it('onError shows invalid URL message for "video not found" errors', () => {
		expect(capturedMutationOptions).toBeDefined();

		capturedMutationOptions.onError(new Error('video not found: abc123'));

		expect(mockToastError).toHaveBeenCalledWith('Invalid YouTube URL or video not found');
	});

	it('onError shows generic message for unknown errors', () => {
		expect(capturedMutationOptions).toBeDefined();

		capturedMutationOptions.onError(new Error('failed to create content'));

		expect(mockToastError).toHaveBeenCalledWith('Failed to add video. Please try again.');
	});

	it('onError does not misclassify errors containing "invalid" in unrelated context', () => {
		expect(capturedMutationOptions).toBeDefined();

		capturedMutationOptions.onError(new Error('invalid connection state'));

		expect(mockToastError).toHaveBeenCalledWith('Failed to add video. Please try again.');
	});

	it('mutationFn is defined and callable', () => {
		expect(capturedMutationOptions).toBeDefined();
		expect(capturedMutationOptions.mutationFn).toBeDefined();
		expect(typeof capturedMutationOptions.mutationFn).toBe('function');
	});

	it('mutationFn calls graphqlClient.request with correct args', async () => {
		expect(capturedMutationOptions).toBeDefined();
		const { graphqlClient } = await import('$lib/queries/client');
		(graphqlClient.request as any).mockResolvedValue({ createContentFromYouTube: { name: 'Test' } });
		await capturedMutationOptions.mutationFn('https://youtube.com/watch?v=abc123');
		expect(graphqlClient.request).toHaveBeenCalledWith(
			expect.anything(),
			{ input: { url: 'https://youtube.com/watch?v=abc123' } }
		);
	});
});
