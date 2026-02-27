/**
 * Resolve @this/@here tokens in text to the parent content name.
 *
 * Storage format: "@this ran 22.3 mph"
 * Display format: "Bo Jackson ran 22.3 mph"
 *
 * Tokens are stored raw in the database and resolved at display time so that
 * the parent content's name can change without updating claim text.
 */
export function resolveAtReference(text: string, parentContentName: string): string {
	return text.replace(/@this|@here/gi, parentContentName);
}

/**
 * Check if text contains @this or @here tokens.
 */
export function hasAtReference(text: string): boolean {
	return /@this|@here/i.test(text);
}
