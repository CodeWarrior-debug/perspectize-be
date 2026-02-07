import '@testing-library/jest-dom/vitest';
import { vi } from 'vitest';
import { readable } from 'svelte/store';

// Mock $app/environment for tests
vi.mock('$app/environment', () => ({
	browser: true,
	dev: true,
	building: false
}));

// Mock $app/navigation if needed
vi.mock('$app/navigation', () => ({
	goto: vi.fn(),
	invalidate: vi.fn(),
	invalidateAll: vi.fn(),
	preloadData: vi.fn(),
	preloadCode: vi.fn(),
	afterNavigate: vi.fn(),
	beforeNavigate: vi.fn()
}));

// Mock $app/stores for components that use page store
vi.mock('$app/stores', () => {
	return {
		page: readable({
			url: new URL('http://localhost'),
			params: {},
			route: { id: '/' },
			status: 200,
			error: null,
			data: {},
			form: null
		}),
		navigating: readable(null),
		updated: { check: vi.fn(), subscribe: readable(false).subscribe }
	};
});

// Mock $lib/assets/favicon.svg
vi.mock('$lib/assets/favicon.svg', () => ({
	default: '/favicon.svg'
}));
