import { browser } from '$app/environment';

const STORAGE_KEY = 'perspectize:selectedUserId';

function loadFromSession(): number | null {
	if (!browser) return null;
	const stored = sessionStorage.getItem(STORAGE_KEY);
	if (!stored) return null;
	const parsed = parseInt(stored, 10);
	return Number.isNaN(parsed) ? null : parsed;
}

function syncToSession(value: number | null): void {
	if (!browser) return;
	if (value !== null) {
		sessionStorage.setItem(STORAGE_KEY, String(value));
	} else {
		sessionStorage.removeItem(STORAGE_KEY);
	}
}

// Internal reactive state
let _selectedUserId = $state<number | null>(loadFromSession());

// Export getter/setter functions for external access
export function getSelectedUserId(): number | null {
	return _selectedUserId;
}

export function setSelectedUserId(value: number | null): void {
	_selectedUserId = value;
	syncToSession(value);
}

// Export object with value getter/setter for convenience
export const selectedUserId = {
	get value() {
		return _selectedUserId;
	},
	set value(newValue: number | null) {
		_selectedUserId = newValue;
		syncToSession(newValue);
	}
};

// Helper to clear selection
export function clearUserSelection(): void {
	_selectedUserId = null;
	syncToSession(null);
}
