import DOMPurify from 'dompurify';

/**
 * Sanitize HTML from user-generated rich text (Tiptap review/comment content).
 * Allows safe formatting tags only — strips scripts, event handlers, and dangerous attributes.
 */
export function sanitizeHtml(dirty: string): string {
	return DOMPurify.sanitize(dirty, {
		ALLOWED_TAGS: [
			'p',
			'br',
			'strong',
			'b',
			'em',
			'i',
			'u',
			'ul',
			'ol',
			'li',
			'a',
			'span',
		],
		ALLOWED_ATTR: ['href', 'target', 'rel', 'class'],
	});
}
