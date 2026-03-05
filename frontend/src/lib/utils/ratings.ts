/**
 * Rating conversion utilities for the Perspectize perspective form.
 *
 * Storage: integer 0–10000 in PostgreSQL (thousandths precision)
 * Display: 0.000 to 10.000 (3 decimal places)
 * Conversion: display * 1000 = storage (e.g., 9.234 → 9234)
 */

/** Convert storage value (0-10000) to display string (0.000-10.000, 3 decimal places) */
export function ratingToDisplay(value: number): string {
	return (value / 1000).toFixed(3);
}

/** Convert display value (0.000-10.000) to storage value (integer 0-10000) */
export function displayToRating(display: number): number {
	return Math.round(display * 1000);
}

/** Validate that a storage rating is within range and is an integer */
export function isValidRating(value: number): boolean {
	return Number.isInteger(value) && value >= 0 && value <= 10000;
}

/** Step amount for rating inputs — 0.250 display = 250 storage units per click */
export const RATING_STEP = 0.25;

/** Minimum display value */
export const RATING_MIN = 0;

/** Maximum display value */
export const RATING_MAX = 10;

/** Default display value shown in create mode (gray/muted until user interacts) */
export const RATING_DEFAULT_DISPLAY = 5.0;

/** Default storage value corresponding to RATING_DEFAULT_DISPLAY */
export const RATING_DEFAULT_STORAGE = 5000;
