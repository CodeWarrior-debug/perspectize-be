import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@testing-library/svelte';
import ColumnPickerDialog from '$lib/components/ColumnPickerDialog.svelte';
import { DATA_COLUMNS, INTERNAL_COLUMNS } from '$lib/utils/grid-config';

const noop = () => {};

describe('ColumnPickerDialog', () => {
	beforeEach(() => vi.clearAllMocks());

	it('renders a checkbox row for every data column when open', () => {
		render(ColumnPickerDialog, { props: { open: true, onToggle: noop } });
		for (const col of DATA_COLUMNS) {
			expect(screen.getByLabelText(col.label)).toBeTruthy();
		}
	});

	it('hides the Internal group for non-admins', () => {
		render(ColumnPickerDialog, { props: { open: true, isAdmin: false, onToggle: noop } });
		expect(screen.queryByText('Internal')).toBeNull();
		expect(screen.queryByLabelText('Content ID')).toBeNull();
	});

	it('shows the Internal group for admins', () => {
		render(ColumnPickerDialog, { props: { open: true, isAdmin: true, onToggle: noop } });
		expect(screen.getByText('Internal')).toBeTruthy();
		for (const col of INTERNAL_COLUMNS) {
			expect(screen.getByLabelText(col.label)).toBeTruthy();
		}
	});

	it('reflects the visibility map in checkbox checked state', () => {
		render(ColumnPickerDialog, {
			props: { open: true, visibility: { tags: true, views: false }, onToggle: noop },
		});
		expect((screen.getByLabelText('Tags') as HTMLInputElement).checked).toBe(true);
		expect((screen.getByLabelText('Views') as HTMLInputElement).checked).toBe(false);
	});

	it('calls onToggle with (colId, next) when a checkbox changes', async () => {
		const onToggle = vi.fn();
		render(ColumnPickerDialog, { props: { open: true, visibility: { tags: true }, onToggle } });
		(screen.getByLabelText('Tags') as HTMLInputElement).click();
		expect(onToggle).toHaveBeenCalledWith('tags', false);
	});

	it('shows the manual-mode notice only when overrideActive', () => {
		const { rerender } = render(ColumnPickerDialog, {
			props: { open: true, overrideActive: false, onToggle: noop },
		});
		expect(screen.queryByTestId('manual-mode-notice')).toBeNull();
		return rerender({ open: true, overrideActive: true, onToggle: noop }).then(() => {
			expect(screen.getByTestId('manual-mode-notice')).toBeTruthy();
		});
	});
});
