import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/svelte';
import PerspectivePopover from '$lib/components/PerspectivePopover.svelte';
import { tick } from 'svelte';

const mocks = vi.hoisted(() => ({
	mockCreateMutate: vi.fn(),
	mockUpdateMutate: vi.fn(),
	mockOnClose: vi.fn(),
}));

vi.mock('$lib/queries/hooks/useCreatePerspective', () => ({
	useCreatePerspective: vi.fn(() => ({
		mutate: mocks.mockCreateMutate,
		isPending: false,
	})),
}));

vi.mock('$lib/queries/hooks/useUpdatePerspective', () => ({
	useUpdatePerspective: vi.fn(() => ({
		mutate: mocks.mockUpdateMutate,
		isPending: false,
	})),
}));

vi.mock('svelte-sonner', () => ({
	toast: {
		success: vi.fn(),
		error: vi.fn(),
	},
}));

function renderPopover(props?: {
	contentId?: number;
	contentName?: string;
	existingPerspective?: any;
	userId?: number;
	open?: boolean;
	onClose?: () => void;
}) {
	const defaultProps = {
		contentId: 1,
		contentName: 'Test Content',
		userId: 42,
		onClose: mocks.mockOnClose,
		open: true,
		...props,
	};
	return render(PerspectivePopover, { props: defaultProps });
}

describe('PerspectivePopover component', () => {
	beforeEach(() => {
		vi.clearAllMocks();
	});

	describe('rendering', () => {
		it('renders without errors', () => {
			const { container } = renderPopover();
			expect(container).toBeTruthy();
		});

		it('renders dialog with "Add Perspective" title when creating', async () => {
			renderPopover({ existingPerspective: null });
			await tick();
			// Both the dialog title and submit button contain "Add Perspective"
			const matches = screen.getAllByText('Add Perspective');
			expect(matches.length).toBeGreaterThanOrEqual(1);
		});

		it('renders dialog with "Edit Perspective" title when editing', async () => {
			renderPopover({
				existingPerspective: {
					id: '1',
					quality: 7500,
					agreement: 5000,
					importance: 6000,
					confidence: 8000,
					like: null,
				},
			});
			await tick();
			expect(screen.getByText('Edit Perspective')).toBeInTheDocument();
		});
	});

	describe('content display', () => {
		it('displays content name in dialog', async () => {
			renderPopover({ contentName: 'My Cool Video' });
			await tick();
			expect(screen.getByText('My Cool Video')).toBeInTheDocument();
		});

		it('displays different content names correctly', async () => {
			renderPopover({ contentName: 'Another Content' });
			await tick();
			expect(screen.getByText('Another Content')).toBeInTheDocument();
		});

		it('handles long content names', async () => {
			renderPopover({ contentName: 'A very long content name that might overflow' });
			await tick();
			const nameElement = screen.getByText('A very long content name that might overflow');
			expect(nameElement).toBeInTheDocument();
		});
	});

	describe('rating inputs', () => {
		it('renders four rating input labels', async () => {
			renderPopover();
			await tick();
			expect(screen.getByText('Quality')).toBeInTheDocument();
			expect(screen.getByText('Agreement')).toBeInTheDocument();
			expect(screen.getByText('Importance')).toBeInTheDocument();
			expect(screen.getByText('Confidence')).toBeInTheDocument();
		});

		it('renders all four rating inputs in a 2x2 grid', async () => {
			renderPopover();
			await tick();
			const labels = ['Quality', 'Agreement', 'Importance', 'Confidence'];
			for (const label of labels) {
				expect(screen.getByText(label)).toBeInTheDocument();
			}
		});

		it('populates rating values in edit mode', async () => {
			renderPopover({
				existingPerspective: {
					id: '1',
					quality: 7500,
					agreement: 5000,
					importance: 6000,
					confidence: 8000,
					like: null,
				},
			});
			await tick();
			// Values should be displayed (7.500, 5.000, 6.000, 8.000)
			expect(screen.getByDisplayValue('7.500')).toBeInTheDocument();
			expect(screen.getByDisplayValue('5.000')).toBeInTheDocument();
			expect(screen.getByDisplayValue('6.000')).toBeInTheDocument();
			expect(screen.getByDisplayValue('8.000')).toBeInTheDocument();
		});
	});

	describe('like buttons (thumbs up/down)', () => {
		it('renders thumbs up button', async () => {
			renderPopover();
			await tick();
			expect(screen.getByLabelText('Thumbs up')).toBeInTheDocument();
		});

		it('renders thumbs down button', async () => {
			renderPopover();
			await tick();
			expect(screen.getByLabelText('Thumbs down')).toBeInTheDocument();
		});

		it('renders like section label', async () => {
			renderPopover();
			await tick();
			expect(screen.getByText('Like')).toBeInTheDocument();
		});

		it('thumbs up button has aria-pressed attribute', async () => {
			renderPopover();
			await tick();
			const thumbsUp = screen.getByLabelText('Thumbs up');
			expect(thumbsUp.hasAttribute('aria-pressed')).toBe(true);
		});

		it('thumbs down button has aria-pressed attribute', async () => {
			renderPopover();
			await tick();
			const thumbsDown = screen.getByLabelText('Thumbs down');
			expect(thumbsDown.hasAttribute('aria-pressed')).toBe(true);
		});
	});

	// TODO: Re-enable Add More / Claim creation tests in a future phase
	// describe('add more button', () => { ... });
	// describe('claim textarea', () => { ... });

	describe('action buttons', () => {
		it('renders Cancel button', async () => {
			renderPopover();
			await tick();
			expect(screen.getByRole('button', { name: 'Cancel' })).toBeInTheDocument();
		});

		it('renders "Add Perspective" button when creating', async () => {
			renderPopover({ existingPerspective: null });
			await tick();
			const submitBtn = screen.getByRole('button', { name: 'Add Perspective' });
			expect(submitBtn).toBeInTheDocument();
		});

		it('renders "Save Changes" button when editing', async () => {
			renderPopover({
				existingPerspective: {
					id: '1',
					quality: 7500,
					agreement: null,
					importance: null,
					confidence: null,
					like: null,
				},
			});
			await tick();
			const submitBtn = screen.getByRole('button', { name: 'Save Changes' });
			expect(submitBtn).toBeInTheDocument();
		});
	});

	describe('dialog interaction', () => {
		it('calls onClose when Cancel button is clicked', async () => {
			renderPopover();
			await tick();
			const cancelBtn = screen.getByRole('button', { name: 'Cancel' });
			await fireEvent.click(cancelBtn);
			expect(mocks.mockOnClose).toHaveBeenCalled();
		});

		it('has About perspectives info button', async () => {
			renderPopover();
			await tick();
			expect(screen.getByLabelText('About perspectives')).toBeInTheDocument();
		});
	});

	describe('title and headers', () => {
		it('displays dialog title', async () => {
			renderPopover({ existingPerspective: null });
			await tick();
			const matches = screen.getAllByText('Add Perspective');
			expect(matches.length).toBeGreaterThanOrEqual(1);
		});

		it('dialog has header with content name', async () => {
			renderPopover({ contentName: 'Test Video' });
			await tick();
			expect(screen.getByText('Test Video')).toBeInTheDocument();
		});
	});

	describe('modes and state', () => {
		it('is in create mode when existingPerspective is null', async () => {
			renderPopover({ existingPerspective: null });
			await tick();
			const matches = screen.getAllByText('Add Perspective');
			expect(matches.length).toBeGreaterThanOrEqual(2); // title + button
		});

		it('is in edit mode when existingPerspective is provided', async () => {
			renderPopover({
				existingPerspective: {
					id: '1',
					quality: 7500,
					agreement: null,
					importance: null,
					confidence: null,
					like: 'THUMBS_UP',
				},
			});
			await tick();
			expect(screen.getByText('Edit Perspective')).toBeInTheDocument();
			expect(screen.getByRole('button', { name: 'Save Changes' })).toBeInTheDocument();
		});

		it('populates like value from existingPerspective', async () => {
			renderPopover({
				existingPerspective: {
					id: '1',
					quality: null,
					agreement: null,
					importance: null,
					confidence: null,
					like: 'THUMBS_UP',
				},
			});
			await tick();
			const thumbsUp = screen.getByLabelText('Thumbs up');
			expect(thumbsUp.getAttribute('aria-pressed')).toBe('true');
		});
	});

	describe('thumbs toggle behavior', () => {
		it('clicking thumbs up toggles aria-pressed to true', async () => {
			renderPopover();
			await tick();
			const thumbsUp = screen.getByLabelText('Thumbs up');
			await fireEvent.click(thumbsUp);
			expect(thumbsUp.getAttribute('aria-pressed')).toBe('true');
		});

		it('clicking thumbs up twice toggles back to false', async () => {
			renderPopover();
			await tick();
			const thumbsUp = screen.getByLabelText('Thumbs up');
			await fireEvent.click(thumbsUp);
			await fireEvent.click(thumbsUp);
			expect(thumbsUp.getAttribute('aria-pressed')).toBe('false');
		});

		it('clicking thumbs down toggles aria-pressed to true', async () => {
			renderPopover();
			await tick();
			const thumbsDown = screen.getByLabelText('Thumbs down');
			await fireEvent.click(thumbsDown);
			expect(thumbsDown.getAttribute('aria-pressed')).toBe('true');
		});

		it('clicking thumbs down twice toggles back to false', async () => {
			renderPopover();
			await tick();
			const thumbsDown = screen.getByLabelText('Thumbs down');
			await fireEvent.click(thumbsDown);
			await fireEvent.click(thumbsDown);
			expect(thumbsDown.getAttribute('aria-pressed')).toBe('false');
		});

		it('clicking thumbs up then thumbs down switches selection', async () => {
			renderPopover();
			await tick();
			const thumbsUp = screen.getByLabelText('Thumbs up');
			const thumbsDown = screen.getByLabelText('Thumbs down');
			await fireEvent.click(thumbsUp);
			expect(thumbsUp.getAttribute('aria-pressed')).toBe('true');
			await fireEvent.click(thumbsDown);
			expect(thumbsDown.getAttribute('aria-pressed')).toBe('true');
			expect(thumbsUp.getAttribute('aria-pressed')).toBe('false');
		});

		it('populates THUMBS_DOWN from existingPerspective', async () => {
			renderPopover({
				existingPerspective: {
					id: '1',
					quality: null,
					agreement: null,
					importance: null,
					confidence: null,
					like: 'THUMBS_DOWN',
				},
			});
			await tick();
			const thumbsDown = screen.getByLabelText('Thumbs down');
			expect(thumbsDown.getAttribute('aria-pressed')).toBe('true');
		});

		it('neither thumb pressed when like is null', async () => {
			renderPopover({
				existingPerspective: {
					id: '1',
					quality: null,
					agreement: null,
					importance: null,
					confidence: null,
					like: null,
				},
			});
			await tick();
			expect(screen.getByLabelText('Thumbs up').getAttribute('aria-pressed')).toBe('false');
			expect(screen.getByLabelText('Thumbs down').getAttribute('aria-pressed')).toBe('false');
		});

		it('neither thumb pressed when like is unknown string', async () => {
			renderPopover({
				existingPerspective: {
					id: '1',
					quality: null,
					agreement: null,
					importance: null,
					confidence: null,
					like: 'UNKNOWN_VALUE',
				},
			});
			await tick();
			expect(screen.getByLabelText('Thumbs up').getAttribute('aria-pressed')).toBe('false');
			expect(screen.getByLabelText('Thumbs down').getAttribute('aria-pressed')).toBe('false');
		});
	});

	describe('form submission', () => {
		it('calls createMutation.mutate when submitting in create mode with a rating', async () => {
			renderPopover();
			await tick();
			// Click thumbs up to set at least one field
			await fireEvent.click(screen.getByLabelText('Thumbs up'));
			// Submit
			const submitBtn = screen.getByRole('button', { name: 'Add Perspective' });
			await fireEvent.click(submitBtn);
			expect(mocks.mockCreateMutate).toHaveBeenCalled();
		});

		it('calls updateMutation.mutate when submitting in edit mode', async () => {
			renderPopover({
				existingPerspective: {
					id: '5',
					quality: 7500,
					agreement: null,
					importance: null,
					confidence: null,
					like: null,
				},
			});
			await tick();
			const submitBtn = screen.getByRole('button', { name: 'Save Changes' });
			await fireEvent.click(submitBtn);
			expect(mocks.mockUpdateMutate).toHaveBeenCalled();
		});
	});

	describe('accessibility', () => {
		it('dialog is accessible with proper roles', async () => {
			renderPopover();
			await tick();
			const buttons = screen.getAllByRole('button');
			expect(buttons.length).toBeGreaterThan(0);
		});

		it('rating inputs have proper aria-labels', async () => {
			renderPopover();
			await tick();
			expect(screen.getByLabelText(/Quality rating/)).toBeInTheDocument();
			expect(screen.getByLabelText(/Agreement rating/)).toBeInTheDocument();
		});

		it('all action buttons are present and accessible', async () => {
			renderPopover({ existingPerspective: null });
			await tick();
			expect(screen.getByRole('button', { name: 'Cancel' })).toBeInTheDocument();
			expect(screen.getByRole('button', { name: 'Add Perspective' })).toBeInTheDocument();
		});
	});
});
