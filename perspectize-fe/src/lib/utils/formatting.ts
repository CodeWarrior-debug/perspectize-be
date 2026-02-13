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
 * Format count numbers with K/M suffixes.
 */
export function formatCount(count: number | null): string {
	if (count === null) return '--';
	if (count < 1000) return String(count);
	if (count < 1000000) return `${(count / 1000).toFixed(1)}K`;
	return `${(count / 1000000).toFixed(1)}M`;
}

/**
 * Format YouTube publish date (ISO string).
 */
export function formatPublishDate(isoString: string | null): string {
	if (!isoString) return '--';
	return formatDate(isoString);
}

/**
 * Format tags array to comma-separated string.
 */
export function formatTags(tags: string[] | null): string {
	if (!tags || tags.length === 0) return '--';
	return tags.join(', ');
}

/**
 * Truncate description with ellipsis.
 */
export function truncateDescription(desc: string | null, maxLength = 100): string {
	if (!desc) return '--';
	if (desc.length <= maxLength) return desc;
	return desc.substring(0, maxLength) + '...';
}

/**
 * Extract video ID from YouTube URL.
 */
export function extractVideoIdFromUrl(url: string | null): string | null {
	if (!url) return null;
	try {
		const urlObj = new URL(url);
		// youtube.com/watch?v=ID
		if (urlObj.hostname.includes('youtube.com') && urlObj.pathname === '/watch') {
			return urlObj.searchParams.get('v');
		}
		// youtu.be/ID
		if (urlObj.hostname === 'youtu.be') {
			return urlObj.pathname.slice(1);
		}
		return null;
	} catch {
		return null;
	}
}

/**
 * AG Grid cell renderer for item column with thumbnail and clickable title.
 */
export function itemCellRenderer(params: { data?: { name: string; url: string | null } }): HTMLElement | string {
	if (!params.data) return '';

	const container = document.createElement('div');
	container.className = 'flex items-center gap-2';

	// Thumbnail
	const videoId = extractVideoIdFromUrl(params.data.url);
	if (videoId) {
		const img = document.createElement('img');
		img.src = `https://i.ytimg.com/vi/${videoId}/default.jpg`;
		img.alt = '';
		img.className = 'w-10 h-8 object-cover rounded';
		container.appendChild(img);
	}

	// Title link
	if (params.data.url) {
		const a = document.createElement('a');
		a.href = params.data.url;
		a.target = '_blank';
		a.rel = 'noopener noreferrer';
		a.className = 'text-primary hover:underline';
		a.textContent = params.data.name;
		container.appendChild(a);
	} else {
		const span = document.createElement('span');
		span.textContent = params.data.name;
		container.appendChild(span);
	}

	return container;
}

/**
 * AG Grid cell renderer for type column with YouTube icon.
 */
export function typeCellRenderer(params: { data?: { contentType: string } }): HTMLElement | string {
	if (!params.data) return '';

	const container = document.createElement('div');
	container.className = 'flex items-center justify-center';

	// YouTube play button icon (red)
	const svg = document.createElementNS('http://www.w3.org/2000/svg', 'svg');
	svg.setAttribute('width', '20');
	svg.setAttribute('height', '20');
	svg.setAttribute('viewBox', '0 0 24 24');
	svg.setAttribute('fill', '#FF0000');

	const path = document.createElementNS('http://www.w3.org/2000/svg', 'path');
	path.setAttribute('d', 'M10 16.5l6-4.5-6-4.5v9zM12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm0 18c-4.41 0-8-3.59-8-8s3.59-8 8-8 8 3.59 8 8-3.59 8-8 8z');

	svg.appendChild(path);
	container.appendChild(svg);

	return container;
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
