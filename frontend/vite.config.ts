import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vitest/config';
import tailwindcss from '@tailwindcss/vite';
import { SvelteKitPWA } from '@vite-pwa/sveltekit';

export default defineConfig({
	plugins: [
		sveltekit(),
		tailwindcss(),
		SvelteKitPWA({
			registerType: 'autoUpdate',
			workbox: {
				globPatterns: ['**/*.{js,css,html,ico,png,svg,woff2}'],
			},
			manifest: {
				name: 'Perspectize',
				short_name: 'Perspectize',
				description: 'Store, refine, and share perspectives on content',
				theme_color: '#1a365d',
				background_color: '#1a365d',
				display: 'standalone',
				scope: '/',
				start_url: '/',
				icons: [
					{ src: 'icons/icon-192.png', sizes: '192x192', type: 'image/png' },
					{ src: 'icons/icon-512.png', sizes: '512x512', type: 'image/png' },
					{
						src: 'icons/icon-512-maskable.png',
						sizes: '512x512',
						type: 'image/png',
						purpose: 'maskable',
					},
				],
			},
		}),
	],
	resolve: {
		conditions: ['browser'],
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
				'src/lib/components/ActivityTable.svelte', // JSDOM limitation: AG Grid doesn't render in test environment
			],
			thresholds: {
				lines: 80,
				functions: 75, // 75% (vs 80%) due to bits-ui interaction handlers not testable in JSDOM
				branches: 75,
				statements: 80,
			},
		},
	},
});
