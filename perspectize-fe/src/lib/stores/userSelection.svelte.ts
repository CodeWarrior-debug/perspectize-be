import { browser } from '$app/environment';

const STORAGE_KEY = 'perspectize:selectedUserId';

function loadFromSession(): number | null {
	if (!browser) return null;
	const stored = sessionStorage.getItem(STORAGE_KEY);
	if (!stored) return null;
	const parsed = parseInt(stored, 10);
	return Number.isNaN(parsed) ? null : parsed;
}

// Exported reactive state — import and use directly in components
export let selectedUserId = $state<number | null>(loadFromSession());

// Auto-sync to session storage when value changes
if (browser) {
	$effect(() => {
		if (selectedUserId !== null) {
			sessionStorage.setItem(STORAGE_KEY, String(selectedUserId));
		} else {
			sessionStorage.removeItem(STORAGE_KEY);
		}
	});
}

// Helper to clear selection
export function clearUserSelection(): void {
	selectedUserId = null;
}
