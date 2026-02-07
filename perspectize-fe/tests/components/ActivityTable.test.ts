import { describe, it, expect, vi } from 'vitest';
import { render } from '@testing-library/svelte';
import ActivityTable from '$lib/components/ActivityTable.svelte';
import { readFileSync } from 'fs';
import { fileURLToPath } from 'url';
import path from 'path';

// Mock AG Grid component
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

		const { container } = render(ActivityTable, {
			props: { rowData }
		});

		expect(container).toBeTruthy();
	});

	it('accepts searchText prop', () => {
		const { container } = render(ActivityTable, {
			props: {
				rowData: [],
				searchText: 'test search'
			}
		});

		expect(container).toBeTruthy();
	});

	it('accepts loading prop', () => {
		const { container } = render(ActivityTable, {
			props: {
				rowData: [],
				loading: true
			}
		});

		expect(container).toBeTruthy();
	});

	it('renders wrapper div with w-full class', () => {
		const { container } = render(ActivityTable, {
			props: { rowData: [] }
		});

		const wrapper = container.querySelector('.w-full');
		expect(wrapper).toBeTruthy();
	});

	it('component includes updateColumnVisibility function', () => {
		// Read the source file to verify updateColumnVisibility exists
		const __dirname = path.dirname(fileURLToPath(import.meta.url));
		const componentPath = path.resolve(__dirname, '../../src/lib/components/ActivityTable.svelte');
		const source = readFileSync(componentPath, 'utf-8');

		expect(source).toContain('function updateColumnVisibility()');
		expect(source).toContain('gridApi.setColumnsVisible');
	});

	it('component includes onGridSizeChanged handler', () => {
		const __dirname = path.dirname(fileURLToPath(import.meta.url));
		const componentPath = path.resolve(__dirname, '../../src/lib/components/ActivityTable.svelte');
		const source = readFileSync(componentPath, 'utf-8');

		expect(source).toContain('onGridSizeChanged:');
		expect(source).toContain('updateColumnVisibility()');
	});

	it('component includes onFirstDataRendered handler that calls updateColumnVisibility', () => {
		const __dirname = path.dirname(fileURLToPath(import.meta.url));
		const componentPath = path.resolve(__dirname, '../../src/lib/components/ActivityTable.svelte');
		const source = readFileSync(componentPath, 'utf-8');

		expect(source).toContain('onFirstDataRendered:');
		// Verify onFirstDataRendered calls updateColumnVisibility (deferred to avoid AG Grid bean init race)
		const onFirstDataRenderedMatch = source.match(/onFirstDataRendered:\s*\(\)\s*=>\s*{([^}]+)}/);
		expect(onFirstDataRenderedMatch).toBeTruthy();
		expect(onFirstDataRenderedMatch![1]).toContain('updateColumnVisibility()');
	});

	it('duration column includes colId for responsive hiding', () => {
		const __dirname = path.dirname(fileURLToPath(import.meta.url));
		const componentPath = path.resolve(__dirname, '../../src/lib/components/ActivityTable.svelte');
		const source = readFileSync(componentPath, 'utf-8');

		// Verify duration column has colId: 'duration'
		expect(source).toContain("colId: 'duration'");
	});
});
