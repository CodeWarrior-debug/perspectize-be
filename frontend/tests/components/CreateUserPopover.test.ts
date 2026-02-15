import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/svelte';
import CreateUserPopover from '$lib/components/CreateUserPopover.svelte';
import { tick } from 'svelte';

// Hoisted mocks — these are referenced inside vi.mock factories
const { mockMutate, mockInvalidateQueries, mockToastSuccess, mockToastError } = vi.hoisted(
	() => ({
		mockMutate: vi.fn(),
		mockInvalidateQueries: vi.fn(),
		mockToastSuccess: vi.fn(),
		mockToastError: vi.fn(),
	})
);

// Capture mutation options for behavioral testing
let capturedMutationOptions: any;

vi.mock('@tanstack/svelte-query', () => ({
	createMutation: vi.fn((optionsFn: () => any) => {
		capturedMutationOptions = optionsFn();
		return {
			mutate: mockMutate,
			isPending: false,
			isSuccess: false,
			data: undefined,
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

vi.mock('$lib/queries/users', () => ({
	CREATE_USER: 'CREATE_USER_QUERY_STRING',
}));

describe('CreateUserPopover component', () => {
	const mockOnUserCreated = vi.fn();

	beforeEach(() => {
		vi.clearAllMocks();
		capturedMutationOptions = undefined;
	});

	it('renders without errors', () => {
		const result = render(CreateUserPopover, {
			props: { onUserCreated: mockOnUserCreated },
		});
		expect(result.container).toBeTruthy();
	});

	it('renders "New User" trigger button', () => {
		render(CreateUserPopover, { props: { onUserCreated: mockOnUserCreated } });
		expect(screen.getByRole('button', { name: /new user/i })).toBeInTheDocument();
	});

	it('button has UserPlus icon', () => {
		render(CreateUserPopover, { props: { onUserCreated: mockOnUserCreated } });
		const button = screen.getByRole('button', { name: /new user/i });
		const svg = button.querySelector('svg');
		expect(svg).toBeInTheDocument();
	});

	it('mutationFn is defined and callable', () => {
		render(CreateUserPopover, { props: { onUserCreated: mockOnUserCreated } });
		expect(capturedMutationOptions).toBeDefined();
		expect(capturedMutationOptions.mutationFn).toBeDefined();
		expect(typeof capturedMutationOptions.mutationFn).toBe('function');
	});

	it('mutationFn calls graphqlClient.request with correct args', async () => {
		render(CreateUserPopover, { props: { onUserCreated: mockOnUserCreated } });
		expect(capturedMutationOptions).toBeDefined();
		const { graphqlClient } = await import('$lib/queries/client');
		(graphqlClient.request as any).mockResolvedValue({
			createUser: { id: '123', username: 'testuser' },
		});

		await capturedMutationOptions.mutationFn({ username: 'testuser' });

		expect(graphqlClient.request).toHaveBeenCalledWith(expect.anything(), {
			input: { username: 'testuser' },
		});
	});
});

describe('CreateUserPopover mutation callbacks', () => {
	const mockOnUserCreated = vi.fn();

	beforeEach(() => {
		vi.clearAllMocks();
		capturedMutationOptions = undefined;
		render(CreateUserPopover, { props: { onUserCreated: mockOnUserCreated } });
	});

	it('onSuccess shows toast with username', () => {
		expect(capturedMutationOptions).toBeDefined();

		capturedMutationOptions.onSuccess({
			createUser: { id: '123', username: 'testuser' },
		});

		expect(mockToastSuccess).toHaveBeenCalledWith('Created user: testuser');
	});

	it('onSuccess invalidates users.list query cache', () => {
		expect(capturedMutationOptions).toBeDefined();

		capturedMutationOptions.onSuccess({
			createUser: { id: '123', username: 'testuser' },
		});

		expect(mockInvalidateQueries).toHaveBeenCalledWith({
			queryKey: expect.arrayContaining(['users', 'list']),
		});
	});

	it('onError shows "already exists" toast for duplicate errors', () => {
		expect(capturedMutationOptions).toBeDefined();

		capturedMutationOptions.onError(new Error('A user with that username already exists'));

		expect(mockToastError).toHaveBeenCalledWith('A user with that username already exists');
	});

	it('onError shows generic toast for unknown errors', () => {
		expect(capturedMutationOptions).toBeDefined();

		capturedMutationOptions.onError(new Error('failed to create user'));

		expect(mockToastError).toHaveBeenCalledWith('Failed to create user. Please try again.');
	});
});
