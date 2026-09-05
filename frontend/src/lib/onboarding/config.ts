/**
 * Intro / onboarding config shared by guest landing and checklist coach.
 * Bump CURRENT_INTRO_VERSION only on material flow changes that warrant a forced re-intro.
 */

export const CURRENT_INTRO_VERSION = 1;

function optionalEnvUrl(value: string | undefined): string | undefined {
	const trimmed = value?.trim();
	return trimmed ? trimmed : undefined;
}

/**
 * Optional mp4 slots. Unset URLs hide Watch / player UI (copy-only coach still works).
 * Override via Vite env or drop files under `static/onboarding/`.
 */
export const ONBOARDING_VIDEOS = {
	guestProduct: optionalEnvUrl(import.meta.env.VITE_ONBOARDING_VIDEO_GUEST_PRODUCT),
	howAddVideo: optionalEnvUrl(import.meta.env.VITE_ONBOARDING_VIDEO_HOW_ADD_VIDEO),
	howPerspective: optionalEnvUrl(import.meta.env.VITE_ONBOARDING_VIDEO_HOW_PERSPECTIVE),
} as const;

export type OnboardingVideoKey = keyof typeof ONBOARDING_VIDEOS;
