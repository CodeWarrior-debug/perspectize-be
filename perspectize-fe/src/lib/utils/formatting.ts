/**
 * Convert length + lengthUnits to display format.
 */
export function formatDuration(length: number | null, lengthUnits: string | null): string {
	if (length === null) return '—';

	if (lengthUnits === 'seconds') {
		const minutes = Math.floor(length / 60);
		const seconds = length % 60;
		return `${minutes}:${seconds.toString().padStart(2, '0')}`;
	}

	return `${length} ${lengthUnits}`;
}

/**
 * Format ISO date string to locale string.
 */
export function formatDate(isoString: string): string {
	const date = new Date(isoString);
	if (isNaN(date.getTime())) return '—';
	return date.toLocaleDateString('en-US', {
		year: 'numeric',
		month: 'short',
		day: 'numeric'
	});
}

/**
 * AG Grid value getter for duration column.
 */
export function durationValueGetter(params: { data?: { length: number | null; lengthUnits: string | null } }): string {
	if (!params.data) return '—';
	return formatDuration(params.data.length, params.data.lengthUnits);
}

/**
 * AG Grid value formatter for date columns.
 */
export function dateValueFormatter(params: { value?: string }): string {
	return params.value ? formatDate(params.value) : '—';
}

/**
 * AG Grid row ID getter for content rows.
 */
export function contentRowId(params: { data?: { id: string | number } }): string {
	return String(params.data?.id ?? '');
}

/**
 * AG Grid cell renderer for content name column.
 * Returns an anchor if URL exists, otherwise a span.
 */
export function nameCellRenderer(params: { data?: { name: string; url: string | null } }): HTMLElement | string {
	if (!params.data) return '';
	if (params.data.url) {
		const a = document.createElement('a');
		a.href = params.data.url;
		a.target = '_blank';
		a.rel = 'noopener noreferrer';
		a.className = 'text-primary hover:underline';
		a.textContent = params.data.name;
		return a;
	}
	const span = document.createElement('span');
	span.textContent = params.data.name;
	return span;
}
