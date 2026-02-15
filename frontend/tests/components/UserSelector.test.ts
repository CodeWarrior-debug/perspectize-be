import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/svelte';
import UserSelector from '$lib/components/UserSelector.svelte';

// Hoisted mocks for store functions and query state
const {
	mockSetSelectedUserId,
	mockGetSelectedUserId,
	mockQueryState,
} = vi.hoisted(() => ({
	mockSetSelectedUserId: vi.fn(),
	mockGetSelectedUserId: vi.fn((): number | null => null),
	mockQueryState: {
		isLoading: false,
		error: null as Error | null,
		data: {
			users: [
				{ id: '1', username: 'alice' },
				{ id: '2', username: 'bob' },
			]
		} as { users: { id: string; username: string }[] } | null,
	},
}));

let capturedQueryOptions: any;

vi.mock('@tanstack/svelte-query', () => ({
	createQuery: vi.fn((optionsFn: () => any) => {
		capturedQueryOptions = optionsFn();
		return mockQueryState;
	}),
	createMutation: vi.fn(() => ({
		mutate: vi.fn(),
		isPending: false,
		isSuccess: false,
		data: undefined,
	})),
	useQueryClient: vi.fn(() => ({
		invalidateQueries: vi.fn(),
	})),
}));

vi.mock('$lib/queries/client', () => ({
	graphqlClient: { request: vi.fn() },
}));

vi.mock('$lib/queries/users', () => ({
	LIST_USERS: 'mock-list-users-query',
}));

vi.mock('svelte-sonner', () => ({
	toast: {
		success: vi.fn(),
		error: vi.fn(),
	},
}));

vi.mock('$lib/stores/userSelection.svelte', () => ({
	setSelectedUserId: (...args: any[]) => mockSetSelectedUserId(...args),
	getSelectedUserId: () => mockGetSelectedUserId(),
}));

function resetQueryState() {
	mockQueryState.isLoading = false;
	mockQueryState.error = null;
	mockQueryState.data = {
		users: [
			{ id: '1', username: 'alice' },
			{ id: '2', username: 'bob' },
		]
	};
}

describe('UserSelector with data', () => {
	beforeEach(() => {
		vi.clearAllMocks();
		capturedQueryOptions = undefined;
		resetQueryState();
		mockGetSelectedUserId.mockReturnValue(null);
	});

	it('renders without errors', () => {
		const { container } = render(UserSelector);
		expect(container).toBeTruthy();
	});

	it('renders a select element with user options when data is loaded', () => {
		render(UserSelector);
		const select = screen.getByRole('combobox');
		expect(select).toBeInTheDocument();
		expect(screen.getByText('Select user...')).toBeInTheDocument();
		expect(screen.getByText('alice')).toBeInTheDocument();
		expect(screen.getByText('bob')).toBeInTheDocument();
	});

	it('creates query with correct queryKey and staleTime', () => {
		render(UserSelector);
		expect(capturedQueryOptions).toBeDefined();
		expect(capturedQueryOptions.queryKey).toEqual(['app', 'users', 'list']);
		expect(capturedQueryOptions.staleTime).toBe(5 * 60 * 1000);
	});

	it('queryFn is defined', () => {
		render(UserSelector);
		expect(capturedQueryOptions.queryFn).toBeDefined();
		expect(typeof capturedQueryOptions.queryFn).toBe('function');
	});

	it('calls setSelectedUserId with parsed int on change', async () => {
		render(UserSelector);
		const select = screen.getByRole('combobox');
		await fireEvent.change(select, { target: { value: '2' } });
		expect(mockSetSelectedUserId).toHaveBeenCalledWith(2);
	});

	it('calls setSelectedUserId with null when empty option selected', async () => {
		render(UserSelector);
		const select = screen.getByRole('combobox');
		await fireEvent.change(select, { target: { value: '' } });
		expect(mockSetSelectedUserId).toHaveBeenCalledWith(null);
	});

	it('reflects currentUserId in select value', () => {
		mockGetSelectedUserId.mockReturnValue(1);
		render(UserSelector);
		const select = screen.getByRole('combobox') as HTMLSelectElement;
		expect(select.value).toBe('1');
	});

	it('shows empty value when currentUserId is null', () => {
		mockGetSelectedUserId.mockReturnValue(null);
		render(UserSelector);
		const select = screen.getByRole('combobox') as HTMLSelectElement;
		expect(select.value).toBe('');
	});

	it('renders CreateUserPopover adjacent to select', () => {
		render(UserSelector);
		const newUserButton = screen.getByRole('button', { name: /new user/i });
		expect(newUserButton).toBeInTheDocument();
	});
});

describe('UserSelector loading state', () => {
	beforeEach(() => {
		vi.clearAllMocks();
		mockQueryState.isLoading = true;
		mockQueryState.error = null;
		mockQueryState.data = null;
	});

	it('shows disabled select with loading text', () => {
		render(UserSelector);
		const select = screen.getByRole('combobox');
		expect(select).toBeDisabled();
		expect(screen.getByText('Loading users...')).toBeInTheDocument();
	});
});

describe('UserSelector error state', () => {
	beforeEach(() => {
		vi.clearAllMocks();
		mockQueryState.isLoading = false;
		mockQueryState.error = new Error('Network error');
		mockQueryState.data = null;
	});

	it('shows disabled select with error text', () => {
		render(UserSelector);
		const select = screen.getByRole('combobox');
		expect(select).toBeDisabled();
		expect(screen.getByText('Error loading users')).toBeInTheDocument();
	});
});
