import { describe, it, expect, beforeEach, vi } from 'vitest';

describe('userSelection store', () => {
	beforeEach(() => {
		// Clear session storage before each test
		sessionStorage.clear();
		// Reset modules to re-initialize store state
		vi.resetModules();
	});

	it('exports selectedUserId', async () => {
		const store = await import('$lib/stores/userSelection.svelte');
		// Initial value should be null (no session storage value)
		expect(store.selectedUserId.value).toBeNull();
	});

	it('exports clearUserSelection function', async () => {
		const store = await import('$lib/stores/userSelection.svelte');
		expect(typeof store.clearUserSelection).toBe('function');
	});

	it('loads stored user ID from session storage', async () => {
		sessionStorage.setItem('perspectize:selectedUserId', '42');
		const store = await import('$lib/stores/userSelection.svelte');
		expect(store.selectedUserId.value).toBe(42);
	});

	it('returns null for invalid session storage value', async () => {
		sessionStorage.setItem('perspectize:selectedUserId', 'not-a-number');
		const store = await import('$lib/stores/userSelection.svelte');
		expect(store.selectedUserId.value).toBeNull();
	});

	it('returns null when session storage is empty', async () => {
		const store = await import('$lib/stores/userSelection.svelte');
		expect(store.selectedUserId.value).toBeNull();
	});

	it('returns null for empty string in session storage', async () => {
		sessionStorage.setItem('perspectize:selectedUserId', '');
		const store = await import('$lib/stores/userSelection.svelte');
		expect(store.selectedUserId.value).toBeNull();
	});
});
