import { describe, it, expect, vi } from 'vitest';
import { render } from '@testing-library/svelte';
import ActivityTable from '$lib/components/ActivityTable.svelte';

vi.mock('ag-grid-svelte5', () => ({
	default: vi.fn(() => ({
		$$: {},
		$set: vi.fn(),
		$on: vi.fn(),
		$destroy: vi.fn(),
	})),
}));

describe('ActivityTable', () => {
	it('renders without errors with empty rowData', () => {
		const { container } = render(ActivityTable, {
			props: { rowData: [] }
		});
		expect(container).toBeTruthy();
	});

	it('accepts rowData prop', () => {
		const rowData = [
			{
				id: '1',
				name: 'Test Video',
				url: 'https://example.com',
				contentType: 'video',
				length: 300,
				lengthUnits: 'seconds',
				createdAt: '2026-01-01T00:00:00Z',
				updatedAt: '2026-01-02T00:00:00Z',
			}
		];
		const { container } = render(ActivityTable, { props: { rowData } });
		expect(container).toBeTruthy();
	});

	it('accepts searchText prop', () => {
		const { container } = render(ActivityTable, {
			props: { rowData: [], searchText: 'test search' }
		});
		expect(container).toBeTruthy();
	});

	it('accepts loading prop', () => {
		const { container } = render(ActivityTable, {
			props: { rowData: [], loading: true }
		});
		expect(container).toBeTruthy();
	});

	it('renders wrapper div with overflow-x-auto', () => {
		const { container } = render(ActivityTable, {
			props: { rowData: [] }
		});
		const wrapper = container.querySelector('.overflow-x-auto');
		expect(wrapper).toBeTruthy();
	});

	it('renders with default props', () => {
		const { container } = render(ActivityTable, {
			props: { rowData: [] }
		});
		expect(container.querySelector('div')).toBeTruthy();
	});
});
