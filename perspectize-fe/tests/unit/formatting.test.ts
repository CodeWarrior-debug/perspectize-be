import { describe, it, expect } from 'vitest';
import {
	formatDuration,
	formatDate,
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
