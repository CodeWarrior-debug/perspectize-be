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

	it('exports getSelectedUserId function', async () => {
		const store = await import('$lib/stores/userSelection.svelte');
		expect(typeof store.getSelectedUserId).toBe('function');
	});

	it('exports setSelectedUserId function', async () => {
		const store = await import('$lib/stores/userSelection.svelte');
		expect(typeof store.setSelectedUserId).toBe('function');
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

	it('getSelectedUserId returns the current value', async () => {
		sessionStorage.setItem('perspectize:selectedUserId', '99');
		const store = await import('$lib/stores/userSelection.svelte');
		expect(store.getSelectedUserId()).toBe(99);
	});

	it('getSelectedUserId returns null when no value set', async () => {
		const store = await import('$lib/stores/userSelection.svelte');
		expect(store.getSelectedUserId()).toBeNull();
	});

	it('setSelectedUserId updates the value and syncs to session storage', async () => {
		const store = await import('$lib/stores/userSelection.svelte');
		store.setSelectedUserId(123);
		expect(store.getSelectedUserId()).toBe(123);
		expect(sessionStorage.getItem('perspectize:selectedUserId')).toBe('123');
	});

	it('setSelectedUserId with null removes from session storage', async () => {
		sessionStorage.setItem('perspectize:selectedUserId', '42');
		const store = await import('$lib/stores/userSelection.svelte');
		store.setSelectedUserId(null);
		expect(store.getSelectedUserId()).toBeNull();
		expect(sessionStorage.getItem('perspectize:selectedUserId')).toBeNull();
	});

	it('selectedUserId.value setter updates and syncs to session storage', async () => {
		const store = await import('$lib/stores/userSelection.svelte');
		store.selectedUserId.value = 55;
		expect(store.selectedUserId.value).toBe(55);
		expect(sessionStorage.getItem('perspectize:selectedUserId')).toBe('55');
	});

	it('selectedUserId.value setter with null removes from session storage', async () => {
		const store = await import('$lib/stores/userSelection.svelte');
		store.selectedUserId.value = 77;
		expect(sessionStorage.getItem('perspectize:selectedUserId')).toBe('77');
		store.selectedUserId.value = null;
		expect(store.selectedUserId.value).toBeNull();
		expect(sessionStorage.getItem('perspectize:selectedUserId')).toBeNull();
	});

	it('clearUserSelection sets value to null and clears session storage', async () => {
		sessionStorage.setItem('perspectize:selectedUserId', '42');
		const store = await import('$lib/stores/userSelection.svelte');
		expect(store.getSelectedUserId()).toBe(42);
		store.clearUserSelection();
		expect(store.getSelectedUserId()).toBeNull();
		expect(sessionStorage.getItem('perspectize:selectedUserId')).toBeNull();
	});

	it('clearUserSelection works when no value was set', async () => {
		const store = await import('$lib/stores/userSelection.svelte');
		store.clearUserSelection();
		expect(store.getSelectedUserId()).toBeNull();
		expect(sessionStorage.getItem('perspectize:selectedUserId')).toBeNull();
	});
});
