import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vitest/config';
import tailwindcss from '@tailwindcss/vite';

export default defineConfig({
	plugins: [sveltekit(), tailwindcss()],
	resolve: {
		conditions: ['browser']
	},
	test: {
		include: ['src/**/*.{test,spec}.{js,ts}', 'tests/**/*.{test,spec}.{js,ts}'],
		environment: 'jsdom',
		globals: true,
		setupFiles: ['./tests/setup.ts'],
		coverage: {
			provider: 'v8',
			reporter: ['text', 'json', 'html'],
			exclude: [
				'node_modules/',
				'.svelte-kit/',
				'**/*.d.ts',
				'**/*.config.*',
				'**/setup.ts',
				'tests/helpers/**',
				'src/lib/components/shadcn/**',
				'src/routes/**',
				'src/lib/components/ActivityTable.svelte' // JSDOM limitation: AG Grid doesn't render in test environment
			],
			thresholds: {
				lines: 80,
				functions: 75, // 75% (vs 80%) due to bits-ui interaction handlers not testable in JSDOM
				branches: 75,
				statements: 80
			}
		}
	}
});
