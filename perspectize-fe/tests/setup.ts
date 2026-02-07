import '@testing-library/jest-dom/vitest';
import { vi } from 'vitest';

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
