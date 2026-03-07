import { defineConfig } from 'vitest/config';
import { svelte } from '@sveltejs/vite-plugin-svelte';
import { playwright } from '@vitest/browser-playwright';

export default defineConfig({
	plugins: [svelte()],
	resolve: {
		conditions: ['browser'],
		alias: {
			$lib: new URL('./src/lib', import.meta.url).pathname,
			'$app/environment': new URL('./tests/browser/mocks/app-environment.ts', import.meta.url).pathname,
			'$app/navigation': new URL('./tests/browser/mocks/app-navigation.ts', import.meta.url).pathname,
			'$app/stores': new URL('./tests/browser/mocks/app-stores.ts', import.meta.url).pathname,
		},
	},
	test: {
		name: 'browser',
		include: ['tests/browser/**/*.test.ts'],
		browser: {
			enabled: true,
			provider: playwright(),
			instances: [{ browser: 'chromium' }],
		},
	},
});
