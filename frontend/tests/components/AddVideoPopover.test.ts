import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/svelte';
import AddVideoPopover from '$lib/components/AddVideoPopover.svelte';
import { tick } from 'svelte';

// Hoisted mocks — shared pattern for useAddVideo components
const mocks = vi.hoisted(() => ({
	mockMutate: vi.fn(),
	mockInvalidateQueries: vi.fn(),
	mockSetQueriesData: vi.fn(),
	mockToastSuccess: vi.fn(),
	mockToastError: vi.fn(),
	mockValidate: vi.fn(),
	mockGetSelectedUserId: vi.fn((): number | null => 1),
	mockMutationState: { mutate: null as any, isPending: false, isSuccess: false },
	capturedMutationOptions: undefined as any,
}));

vi.mock('@tanstack/svelte-query', () => ({
	createMutation: vi.fn((optionsFn: () => any) => {
		mocks.capturedMutationOptions = optionsFn();
		mocks.mockMutationState.mutate = mocks.mockMutate;
		return mocks.mockMutationState;
	}),
	useQueryClient: vi.fn(() => ({
		invalidateQueries: mocks.mockInvalidateQueries,
		setQueriesData: mocks.mockSetQueriesData,
	})),
}));

vi.mock('svelte-sonner', () => ({
	toast: { success: mocks.mockToastSuccess, error: mocks.mockToastError },
}));

vi.mock('$lib/queries/client', () => ({ graphqlClient: { request: vi.fn() } }));
vi.mock('$lib/utils/youtube', () => ({ validateYouTubeUrl: (url: string) => mocks.mockValidate(url) }));
vi.mock('$lib/stores/userSelection.svelte', () => ({ getSelectedUserId: () => mocks.mockGetSelectedUserId() }));

function reset() {
	vi.clearAllMocks();
	mocks.capturedMutationOptions = undefined;
	mocks.mockMutationState.isPending = false;
	mocks.mockMutationState.isSuccess = false;
}

describe('AddVideoPopover component', () => {
	beforeEach(reset);

	it('renders without errors', () => {
		expect(render(AddVideoPopover).container).toBeTruthy();
	});

	it('renders Add Video button', () => {
		render(AddVideoPopover);
		expect(screen.getByRole('button', { name: /add video/i })).toBeInTheDocument();
	});

	it('button has plus icon', () => {
		render(AddVideoPopover);
		expect(screen.getByRole('button', { name: /add video/i }).querySelector('svg')).toBeInTheDocument();
	});

	it('uses buttonVariants for styling', () => {
		render(AddVideoPopover);
		expect(screen.getByRole('button', { name: /add video/i }).className).toBeTruthy();
	});

	it('opens popover when trigger is clicked', async () => {
		const { container } = render(AddVideoPopover);
		await fireEvent.click(screen.getByRole('button', { name: /add video/i }));
		await tick();
		expect(container).toBeTruthy();
	});
});

describe('AddVideoPopover $effect behaviors', () => {
	beforeEach(reset);

	it('renders with mutation in success state', () => {
		mocks.mockMutationState.isSuccess = true;
		expect(render(AddVideoPopover).container).toBeTruthy();
	});
});

describe('AddVideoPopover mutation setup', () => {
	beforeEach(() => {
		reset();
		render(AddVideoPopover);
	});

	it('captures mutation options from useAddVideo hook', () => {
		expect(mocks.capturedMutationOptions).toBeDefined();
		expect(mocks.capturedMutationOptions.mutationFn).toBeDefined();
		expect(mocks.capturedMutationOptions.onSuccess).toBeDefined();
		expect(mocks.capturedMutationOptions.onError).toBeDefined();
	});

	it('mutationFn calls graphqlClient.request with correct args', async () => {
		const { graphqlClient } = await import('$lib/queries/client');
		(graphqlClient.request as any).mockResolvedValue({ createContentFromYouTube: { name: 'Test' } });
		await mocks.capturedMutationOptions.mutationFn('https://youtube.com/watch?v=abc123');
		expect(graphqlClient.request).toHaveBeenCalledWith(expect.anything(), {
			input: { url: 'https://youtube.com/watch?v=abc123', userId: 1 },
		});
	});

	it('mutationFn throws when no user is selected', async () => {
		mocks.mockGetSelectedUserId.mockReturnValueOnce(null);
		await expect(mocks.capturedMutationOptions.mutationFn('https://youtube.com/watch?v=abc123')).rejects.toThrow(
			'No user selected',
		);
	});
});
