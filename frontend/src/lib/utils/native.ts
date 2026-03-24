import { Capacitor } from '@capacitor/core';
import { Haptics, ImpactStyle } from '@capacitor/haptics';

/**
 * Returns true when running inside a native Capacitor WebView (iOS/Android).
 * Returns false in standard browser or PWA context.
 */
export function isNativePlatform(): boolean {
	return Capacitor.isNativePlatform();
}

/**
 * Returns the current platform: 'ios', 'android', or 'web'.
 */
export function getPlatform(): string {
	return Capacitor.getPlatform();
}

/**
 * Triggers light haptic feedback on native platforms.
 * No-op on web. Used for App Store compliance (Guideline 4.2)
 * and enhanced native UX feel.
 */
export async function nativeFeedback(): Promise<void> {
	if (Capacitor.isNativePlatform()) {
		await Haptics.impact({ style: ImpactStyle.Light });
	}
}
