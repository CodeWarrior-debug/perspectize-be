/**
 * Validates if a string is a YouTube URL.
 *
 * Supports:
 * - youtube.com, www.youtube.com, m.youtube.com (with /watch, /embed, /shorts)
 * - youtu.be (short URLs)
 *
 * Uses URL constructor to avoid catastrophic backtracking issues with complex regex.
 *
 * @param url - The URL string to validate
 * @returns true if the URL is a valid YouTube URL, false otherwise
 */
export function validateYouTubeUrl(url: string): boolean {
	// Empty or whitespace-only strings are invalid
	if (!url.trim()) {
		return false;
	}

	try {
		const urlObj = new URL(url);

		// Valid YouTube hosts
		const validHosts = [
			'youtube.com',
			'www.youtube.com',
			'youtu.be',
			'm.youtube.com'
		];

		if (!validHosts.includes(urlObj.hostname)) {
			return false;
		}

		// For youtube.com hosts, pathname must include /watch, /embed, or /shorts
		if (urlObj.hostname.includes('youtube.com')) {
			return (
				urlObj.pathname.includes('/watch') ||
				urlObj.pathname.includes('/embed') ||
				urlObj.pathname.includes('/shorts')
			);
		}

		// For youtu.be, pathname must have video ID (length > 1)
		if (urlObj.hostname === 'youtu.be') {
			return urlObj.pathname.length > 1;
		}

		return true;
	} catch {
		// Malformed URLs throw, which we treat as invalid
		return false;
	}
}
