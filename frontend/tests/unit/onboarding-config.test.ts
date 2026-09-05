import { describe, it, expect } from 'vitest';
import { CURRENT_INTRO_VERSION, ONBOARDING_VIDEOS } from '$lib/onboarding/config';

describe('onboarding config', () => {
	it('exports CURRENT_INTRO_VERSION as 1', () => {
		expect(CURRENT_INTRO_VERSION).toBe(1);
	});

	it('exposes optional video URL slots', () => {
		expect(ONBOARDING_VIDEOS).toHaveProperty('guestProduct');
		expect(ONBOARDING_VIDEOS).toHaveProperty('howAddVideo');
		expect(ONBOARDING_VIDEOS).toHaveProperty('howPerspective');
		for (const key of ['guestProduct', 'howAddVideo', 'howPerspective'] as const) {
			const value = ONBOARDING_VIDEOS[key];
			expect(value === undefined || typeof value === 'string').toBe(true);
		}
	});
});
