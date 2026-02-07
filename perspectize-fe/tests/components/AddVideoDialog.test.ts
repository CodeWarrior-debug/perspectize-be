import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render } from '@testing-library/svelte';
import AddVideoDialog from '$lib/components/AddVideoDialog.svelte';

// Mock TanStack Query
const mockMutate = vi.fn();
const mockInvalidateQueries = vi.fn();

vi.mock('@tanstack/svelte-query', () => ({
	createMutation: vi.fn(() => ({
		mutate: mockMutate,
		isPending: false,
	})),
	useQueryClient: vi.fn(() => ({
		invalidateQueries: mockInvalidateQueries,
	})),
}));

// Mock svelte-sonner
vi.mock('svelte-sonner', () => ({
	toast: {
		success: vi.fn(),
		error: vi.fn(),
	},
}));

// Mock GraphQL client
vi.mock('$lib/queries/client', () => ({
	graphqlClient: {
		request: vi.fn(),
	},
}));

// Mock YouTube validator
vi.mock('$lib/utils/youtube', () => ({
	validateYouTubeUrl: vi.fn(),
}));

// Mock shadcn components
vi.mock('$lib/components/shadcn', () => ({
	Dialog: vi.fn(() => ({ $$: {}, $set: vi.fn(), $on: vi.fn(), $destroy: vi.fn() })),
	DialogContent: vi.fn(() => ({ $$: {}, $set: vi.fn(), $on: vi.fn(), $destroy: vi.fn() })),
	DialogHeader: vi.fn(() => ({ $$: {}, $set: vi.fn(), $on: vi.fn(), $destroy: vi.fn() })),
	DialogTitle: vi.fn(() => ({ $$: {}, $set: vi.fn(), $on: vi.fn(), $destroy: vi.fn() })),
	DialogFooter: vi.fn(() => ({ $$: {}, $set: vi.fn(), $on: vi.fn(), $destroy: vi.fn() })),
	Button: vi.fn(() => ({ $$: {}, $set: vi.fn(), $on: vi.fn(), $destroy: vi.fn() })),
	Input: vi.fn(() => ({ $$: {}, $set: vi.fn(), $on: vi.fn(), $destroy: vi.fn() })),
	Label: vi.fn(() => ({ $$: {}, $set: vi.fn(), $on: vi.fn(), $destroy: vi.fn() })),
}));

describe('AddVideoDialog component', () => {
	beforeEach(() => {
		vi.clearAllMocks();
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
