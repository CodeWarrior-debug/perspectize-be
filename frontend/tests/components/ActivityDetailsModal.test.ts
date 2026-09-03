import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/svelte';
import ActivityDetailsModal from '$lib/components/ActivityDetailsModal.svelte';

// Hoisted mocks — shared pattern for useUpdateSourceData components (see AddVideoPopover.test.ts)
const mocks = vi.hoisted(() => ({
	mockMutate: vi.fn(),
	mockInvalidateQueries: vi.fn(),
	mockSetQueriesData: vi.fn(),
	mockToastSuccess: vi.fn(),
	mockToastError: vi.fn(),
	mockMutationState: { mutate: null as any, isPending: false },
}));

vi.mock('@tanstack/svelte-query', () => ({
	createMutation: vi.fn((optionsFn: () => any) => {
		optionsFn();
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

vi.mock('$lib/queries/client', () => ({ graphqlRequest: vi.fn() }));

function reset() {
	vi.clearAllMocks();
	mocks.mockMutationState.isPending = false;
}

const content = {
	id: '42',
	name: 'Stephen Paea breaking bench',
	url: 'https://youtube.com/watch?v=abc123',
	channelTitle: 'TBD tribute',
	viewCount: 1300000,
	likeCount: 26500,
	length: 59,
	lengthUnits: 'seconds',
	publishedAt: '2026-02-24T00:00:00Z',
	updatedAt: '2026-03-01T00:00:00Z',
	tags: ['tom brady', 'tom brady goat'],
};

describe('ActivityDetailsModal', () => {
	beforeEach(reset);

	it('renders nothing when closed', () => {
		render(ActivityDetailsModal, { props: { content, open: false, onClose: vi.fn() } });
		expect(screen.queryByText(content.name)).not.toBeInTheDocument();
	});

	it('renders the video title, channel, link, and stats when open', () => {
		render(ActivityDetailsModal, { props: { content, open: true, onClose: vi.fn() } });

		expect(screen.getByText(content.name)).toBeInTheDocument();
		expect(screen.getByText('TBD tribute')).toBeInTheDocument();
		expect(screen.getByText(content.url)).toBeInTheDocument();
		expect(screen.getByText('1.3 M')).toBeInTheDocument(); // views
		expect(screen.getByText('26.5 K')).toBeInTheDocument(); // likes
		expect(screen.getByText('0:59')).toBeInTheDocument(); // duration
	});

	it('shows placeholder stats for perspectives and avg rating', () => {
		render(ActivityDetailsModal, { props: { content, open: true, onClose: vi.fn() } });

		expect(screen.getByText('Perspectives')).toBeInTheDocument();
		expect(screen.getByText('Avg. Rating')).toBeInTheDocument();
		expect(screen.getAllByText('—').length).toBeGreaterThanOrEqual(2);
	});

	it('renders tags when present', () => {
		render(ActivityDetailsModal, { props: { content, open: true, onClose: vi.fn() } });
		expect(screen.getByText('tom brady, tom brady goat')).toBeInTheDocument();
	});

	it('omits the tags section when there are none', () => {
		render(ActivityDetailsModal, {
			props: { content: { ...content, tags: null }, open: true, onClose: vi.fn() },
		});
		expect(screen.queryByText('Tags')).not.toBeInTheDocument();
	});

	it('renders an "Update source data" button that triggers the update mutation for this item', async () => {
		render(ActivityDetailsModal, { props: { content, open: true, onClose: vi.fn() } });
		const button = screen.getByRole('button', { name: 'Update source data' });
		expect(button).toBeInTheDocument();

		await fireEvent.click(button);

		expect(mocks.mockMutate).toHaveBeenCalledWith(content.id);
	});

	it('disables the Update source data button while the mutation is pending', () => {
		mocks.mockMutationState.isPending = true;
		render(ActivityDetailsModal, { props: { content, open: true, onClose: vi.fn() } });

		const button = screen.getByRole('button', { name: /update source data/i });
		expect(button).toBeDisabled();
	});

	it('calls onClose when the close button is clicked', async () => {
		const onClose = vi.fn();
		render(ActivityDetailsModal, { props: { content, open: true, onClose } });

		await fireEvent.click(screen.getByRole('button', { name: /close/i }));
		expect(onClose).toHaveBeenCalled();
	});
});
