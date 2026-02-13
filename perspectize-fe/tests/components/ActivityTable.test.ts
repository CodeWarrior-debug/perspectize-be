import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, fireEvent, waitFor } from '@testing-library/svelte';
import { tick } from 'svelte';

const { mockRequest } = vi.hoisted(() => ({
	mockRequest: vi.fn()
}));

// Mock AG Grid component
vi.mock('ag-grid-svelte5', () => ({
	default: vi.fn(() => ({
		$$: {},
		$set: vi.fn(),
		$on: vi.fn(),
		$destroy: vi.fn(),
	})),
}));

// Mock graphqlClient
vi.mock('$lib/queries/client', () => ({
	graphqlClient: {
		request: mockRequest
	}
}));

import ActivityTable from '$lib/components/ActivityTable.svelte';

const mockEmptyResponse = {
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
};

const mockDataResponse = {
	content: {
		items: [
			{
				id: '1', name: 'Test Video', url: 'https://youtube.com/watch?v=abc',
				contentType: 'YOUTUBE', length: 300, lengthUnits: 'seconds',
				viewCount: 1500, likeCount: 100, channelTitle: 'Test Channel',
				publishedAt: '2024-01-15T00:00:00Z', tags: ['test'],
				description: 'A test video', createdAt: '2024-01-20', updatedAt: '2024-01-25'
			}
		],
		pageInfo: {
			hasNextPage: true,
			hasPreviousPage: false,
			startCursor: 'cursor1',
			endCursor: 'cursor2'
		},
		totalCount: 25
	}
};

describe('ActivityTable', () => {
	beforeEach(() => {
		vi.clearAllMocks();
		mockRequest.mockResolvedValue(mockEmptyResponse);
	});

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

	it('renders page size selector with options 10/25/50', () => {
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

	it('has AG Grid container for sticky headers', () => {
		const { container } = render(ActivityTable);
		expect(container.querySelector('.flex-1.min-h-0')).toBeTruthy();
	});

	it('clicking Previous button when on first page does nothing', async () => {
		const { container } = render(ActivityTable);
		const buttons = Array.from(container.querySelectorAll('button'));
		const prevButton = buttons.find(btn => btn.textContent?.includes('Previous')) as HTMLButtonElement;

		await fireEvent.click(prevButton);
		await tick();

		// Should still be on page 1
		expect(container.textContent).toContain('Page 1');
	});

	it('clicking Next button when disabled does nothing', async () => {
		const { container } = render(ActivityTable);
		const buttons = Array.from(container.querySelectorAll('button'));
		const nextButton = buttons.find(btn => btn.textContent?.includes('Next')) as HTMLButtonElement;

		await fireEvent.click(nextButton);
		await tick();

		// Should still be on page 1
		expect(container.textContent).toContain('Page 1');
	});

	it('changing page size triggers data fetch', async () => {
		mockRequest.mockResolvedValue(mockDataResponse);
		const { container } = render(ActivityTable);
		const select = container.querySelector('select#pageSize') as HTMLSelectElement;

		await fireEvent.change(select, { target: { value: '25' } });
		await tick();

		// Should trigger a request with new page size
		await waitFor(() => {
			const calls = mockRequest.mock.calls;
			const hasPageSize25 = calls.some((call: any[]) => call[1]?.first === 25);
			expect(hasPageSize25).toBe(true);
		});
	});

	it('renders with data from GraphQL response after page size change', async () => {
		mockRequest.mockResolvedValue(mockDataResponse);
		const { container } = render(ActivityTable);
		const select = container.querySelector('select#pageSize') as HTMLSelectElement;

		// Trigger fetchData via page size change (onGridReady doesn't fire with mocked AG Grid)
		await fireEvent.change(select, { target: { value: '10' } });

		await waitFor(() => {
			expect(container.textContent).toContain('25 total items');
		});
	});

	it('shows correct pagination for multi-page data', async () => {
		mockRequest.mockResolvedValue(mockDataResponse);
		const { container } = render(ActivityTable);
		const select = container.querySelector('select#pageSize') as HTMLSelectElement;

		// Trigger fetchData via page size change
		await fireEvent.change(select, { target: { value: '10' } });

		await waitFor(() => {
			expect(container.textContent).toContain('Page 1 of 3');
		});
	});

	it('handles GraphQL request error gracefully', async () => {
		mockRequest.mockRejectedValue(new Error('Network error'));
		const { container } = render(ActivityTable);
		const select = container.querySelector('select#pageSize') as HTMLSelectElement;

		// Trigger fetchData via page size change
		await fireEvent.change(select, { target: { value: '10' } });

		await waitFor(() => {
			expect(container.textContent).toContain('0 total items');
		});
	});

	it('has border-t on pagination bar', () => {
		const { container } = render(ActivityTable);
		const paginationBar = container.querySelector('.border-t.border-border');
		expect(paginationBar).toBeTruthy();
	});

	it('has label for page size selector', () => {
		const { container } = render(ActivityTable);
		const label = container.querySelector('label[for="pageSize"]');
		expect(label).toBeTruthy();
		expect(label?.textContent).toContain('Page size');
	});
});
