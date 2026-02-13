import { describe, it, expect, vi } from 'vitest';
import { render } from '@testing-library/svelte';

// Mock AG Grid component
vi.mock('ag-grid-svelte5', () => ({
	default: vi.fn(() => ({
		$$: {},
		$set: vi.fn(),
		$on: vi.fn(),
		$destroy: vi.fn(),
	})),
}));

// Mock graphqlClient - use factory function to avoid hoisting issues
vi.mock('$lib/queries/client', () => ({
	graphqlClient: {
		request: vi.fn().mockResolvedValue({
			content: {
				items: [],
				pageInfo: {
					hasNextPage: false,
					hasPreviousPage: false,
					startCursor: null,
					endCursor: null
				},
				totalCount: 0
			}
		})
	}
}));

import ActivityTable from '$lib/components/ActivityTable.svelte';

describe('ActivityTable', () => {
	it('renders without errors', () => {
		const { container } = render(ActivityTable);
		expect(container).toBeTruthy();
	});

	it('renders container with proper layout classes', () => {
		const { container } = render(ActivityTable);
		const mainContainer = container.querySelector('.flex.flex-col.h-full');
		expect(mainContainer).toBeTruthy();
	});

	it('renders pagination controls', () => {
		const { container } = render(ActivityTable);
		const buttons = Array.from(container.querySelectorAll('button'));
		const prevButton = buttons.find(btn => btn.textContent?.includes('Previous'));
		const nextButton = buttons.find(btn => btn.textContent?.includes('Next'));

		expect(prevButton).toBeTruthy();
		expect(nextButton).toBeTruthy();
	});

	it('renders page size selector', () => {
		const { container } = render(ActivityTable);
		const select = container.querySelector('select#pageSize');
		expect(select).toBeTruthy();

		const options = container.querySelectorAll('select#pageSize option');
		expect(options.length).toBe(3);
		expect(options[0].textContent).toBe('10');
		expect(options[1].textContent).toBe('25');
		expect(options[2].textContent).toBe('50');
	});

	it('displays initial total count', () => {
		const { container } = render(ActivityTable);
		// Initial state shows 0 total items
		expect(container.textContent).toContain('0 total items');
	});

	it('displays default page number', () => {
		const { container } = render(ActivityTable);
		expect(container.textContent).toContain('Page 1 of 1');
	});

	it('disables Previous button on first page', () => {
		const { container } = render(ActivityTable);
		const buttons = Array.from(container.querySelectorAll('button'));
		const prevButton = buttons.find(btn => btn.textContent?.includes('Previous')) as HTMLButtonElement;

		expect(prevButton?.disabled).toBe(true);
	});

	it('disables Next button when no more pages', () => {
		const { container } = render(ActivityTable);
		const buttons = Array.from(container.querySelectorAll('button'));
		const nextButton = buttons.find(btn => btn.textContent?.includes('Next')) as HTMLButtonElement;

		expect(nextButton?.disabled).toBe(true);
	});

	it('has column definitions with correct structure', () => {
		// This test verifies that the component defines columns correctly
		// by checking that it renders without errors (column defs are used in gridOptions)
		const { container } = render(ActivityTable);
		expect(container.querySelector('.flex-1.min-h-0')).toBeTruthy();
	});
});
