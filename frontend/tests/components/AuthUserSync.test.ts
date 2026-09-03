import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render } from '@testing-library/svelte';
import AuthUserSync from '$lib/components/AuthUserSync.svelte';

const { mockSetSelectedUserId, mockClearUserSelection, mockQueryClientClear, mockClerkContext, mockQueryState } =
	vi.hoisted(() => ({
		mockSetSelectedUserId: vi.fn(),
		mockClearUserSelection: vi.fn(),
		mockQueryClientClear: vi.fn(),
		mockClerkContext: {
			isLoaded: true,
			auth: { userId: null as string | null },
		},
		mockQueryState: {
			isLoading: false,
			error: null as Error | null,
			data: null as { me: { id: string; username: string } } | null,
		},
	}));

let capturedQueryOptions: any;

vi.mock('@tanstack/svelte-query', () => ({
	createQuery: vi.fn((optionsFn: () => any) => {
		capturedQueryOptions = optionsFn();
		return mockQueryState;
	}),
	useQueryClient: vi.fn(() => ({
		clear: mockQueryClientClear,
	})),
}));

vi.mock('svelte-clerk', () => ({
	useClerkContext: vi.fn(() => mockClerkContext),
}));

vi.mock('$lib/queries/client', () => ({
	graphqlRequest: vi.fn(),
}));

vi.mock('$lib/queries/users', () => ({
	ME: 'mock-me-query',
}));

vi.mock('$lib/stores/userSelection.svelte', () => ({
	setSelectedUserId: (...args: unknown[]) => mockSetSelectedUserId(...args),
	clearUserSelection: () => mockClearUserSelection(),
}));

describe('AuthUserSync', () => {
	beforeEach(() => {
		vi.clearAllMocks();
		capturedQueryOptions = undefined;
		mockClerkContext.isLoaded = true;
		mockClerkContext.auth.userId = null;
		mockQueryState.isLoading = false;
		mockQueryState.error = null;
		mockQueryState.data = null;
	});

	it('does not sync while Clerk is still loading', () => {
		mockClerkContext.isLoaded = false;
		mockClerkContext.auth.userId = null;
		render(AuthUserSync);
		expect(mockQueryClientClear).not.toHaveBeenCalled();
		expect(mockSetSelectedUserId).not.toHaveBeenCalled();
		expect(mockClearUserSelection).not.toHaveBeenCalled();
	});

	it('on sign-in, clears query cache and syncs selected user from the me query', () => {
		mockClerkContext.isLoaded = true;
		mockClerkContext.auth.userId = 'clerk_user_123';
		mockQueryState.data = { me: { id: '5', username: 'alice' } };
		render(AuthUserSync);
		expect(mockQueryClientClear).toHaveBeenCalledTimes(1);
		expect(mockSetSelectedUserId).toHaveBeenCalledWith(5);
		expect(mockClearUserSelection).not.toHaveBeenCalled();
	});

	it('on sign-out, clears query cache and clears user selection', () => {
		mockClerkContext.isLoaded = true;
		mockClerkContext.auth.userId = null;
		render(AuthUserSync);
		expect(mockQueryClientClear).toHaveBeenCalledTimes(1);
		expect(mockClearUserSelection).toHaveBeenCalledTimes(1);
		expect(mockSetSelectedUserId).not.toHaveBeenCalled();
	});

	it('account switch: mounting as a different signed-in user clears cache and re-syncs', () => {
		// First mount simulates being signed in as user A. A fresh mount is used
		// here (rather than mutating context mid-test) because the mocked
		// useClerkContext() returns a plain object, not a Svelte $state-backed
		// one like the real ClerkProvider — see AuthUserSync.svelte for why a
		// real account switch re-fires the effect via genuine reactivity.
		mockClerkContext.isLoaded = true;
		mockClerkContext.auth.userId = 'clerk_user_A';
		mockQueryState.data = { me: { id: '1', username: 'alice' } };
		const first = render(AuthUserSync);
		expect(mockSetSelectedUserId).toHaveBeenCalledWith(1);
		first.unmount();

		vi.clearAllMocks();

		mockClerkContext.auth.userId = 'clerk_user_B';
		mockQueryState.data = { me: { id: '2', username: 'bob' } };
		render(AuthUserSync);
		expect(mockQueryClientClear).toHaveBeenCalledTimes(1);
		expect(mockSetSelectedUserId).toHaveBeenCalledWith(2);
	});

	it('constructs the me query with a clerk-scoped key, 5 minute staleTime, and enabled when signed in', () => {
		mockClerkContext.isLoaded = true;
		mockClerkContext.auth.userId = 'clerk_user_123';
		render(AuthUserSync);
		expect(capturedQueryOptions.queryKey).toEqual(['me', 'clerk_user_123']);
		expect(capturedQueryOptions.staleTime).toBe(5 * 60 * 1000);
		expect(capturedQueryOptions.enabled).toBe(true);
	});

	it('disables the me query while signed out', () => {
		mockClerkContext.isLoaded = true;
		mockClerkContext.auth.userId = null;
		render(AuthUserSync);
		expect(capturedQueryOptions.enabled).toBe(false);
	});
});
