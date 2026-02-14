import { describe, it, expect } from 'vitest';
import {
	formatDuration,
	formatDate,
	formatCount,
	formatPublishDate,
	formatTags,
	truncateDescription,
	extractVideoIdFromUrl,
	itemCellRenderer,
	typeCellRenderer,
	nameCellRenderer,
	durationValueGetter,
	dateValueFormatter,
	contentRowId,
} from '$lib/utils/formatting';

describe('formatDuration', () => {
	it('returns dash for null length', () => {
		expect(formatDuration(null, null)).toBe('—');
	});

	it('returns dash for null length with units', () => {
		expect(formatDuration(null, 'seconds')).toBe('—');
	});

	it('formats seconds as minutes:seconds', () => {
		expect(formatDuration(300, 'seconds')).toBe('5:00');
	});

	it('formats seconds with padded seconds part', () => {
		expect(formatDuration(65, 'seconds')).toBe('1:05');
	});

	it('formats zero seconds', () => {
		expect(formatDuration(0, 'seconds')).toBe('0:00');
	});

	it('formats seconds under a minute', () => {
		expect(formatDuration(45, 'seconds')).toBe('0:45');
	});

	it('formats large durations', () => {
		expect(formatDuration(3661, 'seconds')).toBe('61:01');
	});

	it('formats non-seconds units with value and unit', () => {
		expect(formatDuration(10, 'minutes')).toBe('10 minutes');
	});

	it('formats with null units', () => {
		expect(formatDuration(10, null)).toBe('10 null');
	});
});

describe('formatDate', () => {
	it('formats a valid ISO date string', () => {
		expect(formatDate('2026-01-15T12:00:00Z')).toMatch(/Jan 15, 2026/);
	});

	it('formats another valid date', () => {
		expect(formatDate('2026-06-20T12:00:00Z')).toMatch(/Jun 20, 2026/);
	});

	it('returns dash for invalid date string', () => {
		expect(formatDate('not-a-date')).toBe('—');
	});

	it('returns dash for empty string', () => {
		expect(formatDate('')).toBe('—');
	});
});

describe('durationValueGetter', () => {
	it('returns dash when data is missing', () => {
		expect(durationValueGetter({ data: undefined })).toBe('—');
	});

	it('returns formatted duration for seconds', () => {
		expect(durationValueGetter({ data: { length: 300, lengthUnits: 'seconds' } })).toBe('5:00');
	});

	it('returns dash for null length', () => {
		expect(durationValueGetter({ data: { length: null, lengthUnits: null } })).toBe('—');
	});

	it('returns formatted duration for non-seconds', () => {
		expect(durationValueGetter({ data: { length: 10, lengthUnits: 'minutes' } })).toBe('10 minutes');
	});
});

describe('dateValueFormatter', () => {
	it('formats a valid date value', () => {
		expect(dateValueFormatter({ value: '2026-06-20T12:00:00Z' })).toMatch(/Jun 20, 2026/);
	});

	it('returns dash for null value', () => {
		expect(dateValueFormatter({ value: undefined })).toBe('—');
	});

	it('returns dash for empty string value', () => {
		expect(dateValueFormatter({ value: '' })).toBe('—');
	});

	it('returns dash for invalid date', () => {
		expect(dateValueFormatter({ value: 'not-a-date' })).toBe('—');
	});
});

describe('contentRowId', () => {
	it('returns string ID from numeric data', () => {
		expect(contentRowId({ data: { id: 42 } })).toBe('42');
	});

	it('returns string ID from string data', () => {
		expect(contentRowId({ data: { id: 'abc' } })).toBe('abc');
	});

	it('returns empty string when data is undefined', () => {
		expect(contentRowId({ data: undefined })).toBe('');
	});
});

describe('nameCellRenderer', () => {
	it('returns empty string when data is missing', () => {
		expect(nameCellRenderer({ data: undefined })).toBe('');
	});

	it('returns anchor element when URL is present', () => {
		const result = nameCellRenderer({
			data: { name: 'My Video', url: 'https://youtube.com/watch?v=abc' }
		});
		expect(result).toBeInstanceOf(HTMLAnchorElement);
		const anchor = result as HTMLAnchorElement;
		expect(anchor.href).toBe('https://youtube.com/watch?v=abc');
		expect(anchor.textContent).toBe('My Video');
		expect(anchor.target).toBe('_blank');
		expect(anchor.rel).toBe('noopener noreferrer');
		expect(anchor.className).toBe('text-primary hover:underline');
	});

	it('returns span element when URL is null', () => {
		const result = nameCellRenderer({
			data: { name: 'No URL Video', url: null }
		});
		expect(result).toBeInstanceOf(HTMLSpanElement);
		expect((result as HTMLSpanElement).textContent).toBe('No URL Video');
	});

	it('returns span element when URL is empty string', () => {
		const result = nameCellRenderer({
			data: { name: 'Empty URL', url: '' as unknown as null }
		});
		expect(result).toBeInstanceOf(HTMLSpanElement);
		expect((result as HTMLSpanElement).textContent).toBe('Empty URL');
	});
});

describe('formatCount', () => {
	it('returns -- for null', () => {
		expect(formatCount(null)).toBe('--');
	});

	it('returns number as string for counts under 1000', () => {
		expect(formatCount(0)).toBe('0');
		expect(formatCount(500)).toBe('500');
		expect(formatCount(999)).toBe('999');
	});

	it('formats thousands with K suffix', () => {
		expect(formatCount(1000)).toBe('1.0K');
		expect(formatCount(1234)).toBe('1.2K');
		expect(formatCount(5678)).toBe('5.7K');
		expect(formatCount(999999)).toBe('1000.0K');
	});

	it('formats millions with M suffix', () => {
		expect(formatCount(1000000)).toBe('1.0M');
		expect(formatCount(1234567)).toBe('1.2M');
		expect(formatCount(5678901)).toBe('5.7M');
	});
});

describe('formatPublishDate', () => {
	it('returns -- for null', () => {
		expect(formatPublishDate(null)).toBe('--');
	});

	it('formats valid ISO date string', () => {
		expect(formatPublishDate('2026-01-15T12:00:00Z')).toMatch(/Jan 15, 2026/);
	});
});

describe('formatTags', () => {
	it('returns -- for null', () => {
		expect(formatTags(null)).toBe('--');
	});

	it('returns -- for empty array', () => {
		expect(formatTags([])).toBe('--');
	});

	it('formats single tag', () => {
		expect(formatTags(['technology'])).toBe('technology');
	});

	it('formats multiple tags comma-separated', () => {
		expect(formatTags(['tech', 'science', 'ai'])).toBe('tech, science, ai');
	});
});

describe('truncateDescription', () => {
	it('returns -- for null', () => {
		expect(truncateDescription(null)).toBe('--');
	});

	it('returns -- for empty string', () => {
		expect(truncateDescription('')).toBe('--');
	});

	it('returns full description if under max length', () => {
		expect(truncateDescription('Short description', 100)).toBe('Short description');
	});

	it('truncates description with ellipsis if over max length', () => {
		const long = 'A'.repeat(150);
		const result = truncateDescription(long, 100);
		expect(result).toBe('A'.repeat(100) + '...');
		expect(result.length).toBe(103);
	});

	it('uses default max length of 100', () => {
		const long = 'A'.repeat(150);
		const result = truncateDescription(long);
		expect(result).toBe('A'.repeat(100) + '...');
	});
});

describe('extractVideoIdFromUrl', () => {
	it('returns null for null url', () => {
		expect(extractVideoIdFromUrl(null)).toBeNull();
	});

	it('extracts ID from youtube.com/watch URL', () => {
		expect(extractVideoIdFromUrl('https://www.youtube.com/watch?v=dQw4w9WgXcQ')).toBe('dQw4w9WgXcQ');
		expect(extractVideoIdFromUrl('https://youtube.com/watch?v=abc123')).toBe('abc123');
		expect(extractVideoIdFromUrl('https://m.youtube.com/watch?v=xyz789')).toBe('xyz789');
	});

	it('extracts ID from youtu.be URL', () => {
		expect(extractVideoIdFromUrl('https://youtu.be/dQw4w9WgXcQ')).toBe('dQw4w9WgXcQ');
		expect(extractVideoIdFromUrl('https://youtu.be/abc123')).toBe('abc123');
	});

	it('returns null for invalid URL', () => {
		expect(extractVideoIdFromUrl('not a url')).toBeNull();
	});

	it('returns null for non-YouTube URL', () => {
		expect(extractVideoIdFromUrl('https://example.com')).toBeNull();
	});

	it('returns null for YouTube URL without video ID', () => {
		expect(extractVideoIdFromUrl('https://www.youtube.com')).toBeNull();
		expect(extractVideoIdFromUrl('https://www.youtube.com/channel/UC123')).toBeNull();
	});
});

describe('itemCellRenderer', () => {
	it('returns empty string when data is missing', () => {
		expect(itemCellRenderer({ data: undefined })).toBe('');
	});

	it('renders container with thumbnail and title link', () => {
		const result = itemCellRenderer({
			data: { name: 'My Video', url: 'https://youtube.com/watch?v=abc123' }
		}) as HTMLElement;

		expect(result).toBeInstanceOf(HTMLDivElement);
		expect(result.className).toContain('flex');
		expect(result.className).toContain('items-center');

		const img = result.querySelector('img');
		expect(img).toBeTruthy();
		expect(img?.src).toBe('https://i.ytimg.com/vi/abc123/default.jpg');
		expect(img?.className).toContain('w-10');

		const link = result.querySelector('a');
		expect(link).toBeTruthy();
		expect(link?.href).toBe('https://youtube.com/watch?v=abc123');
		expect(link?.textContent).toBe('My Video');
		expect(link?.target).toBe('_blank');
	});

	it('renders without thumbnail when no video ID', () => {
		const result = itemCellRenderer({
			data: { name: 'No URL', url: null }
		}) as HTMLElement;

		expect(result.querySelector('img')).toBeNull();
		const span = result.querySelector('span');
		expect(span?.textContent).toBe('No URL');
	});
});

describe('typeCellRenderer', () => {
	it('returns empty string when data is missing', () => {
		expect(typeCellRenderer({ data: undefined })).toBe('');
	});

	it('renders YouTube play icon', () => {
		const result = typeCellRenderer({
			data: { contentType: 'youtube_video' }
		}) as HTMLElement;

		expect(result).toBeInstanceOf(HTMLDivElement);
		expect(result.className).toContain('flex');

		const svg = result.querySelector('svg');
		expect(svg).toBeTruthy();
		expect(svg?.getAttribute('fill')).toBe('#FF0000');
		expect(svg?.getAttribute('viewBox')).toBe('0 0 24 24');

		const path = svg?.querySelector('path');
		expect(path).toBeTruthy();
	});
});
