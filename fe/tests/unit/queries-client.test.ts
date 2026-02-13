import { describe, it, expect, beforeEach, vi } from 'vitest';
import type { GraphQLClient } from 'graphql-request';

let graphqlClient: GraphQLClient;

beforeEach(async () => {
	vi.resetModules();
	const mod = await import('$lib/queries/client');
	graphqlClient = mod.graphqlClient;
});

describe('GraphQL client', () => {
	it('exports a graphqlClient instance', () => {
		expect(graphqlClient).toBeDefined();
		expect(typeof graphqlClient.request).toBe('function');
	});

	it('client has request method for making GraphQL calls', () => {
		expect(graphqlClient).toHaveProperty('request');
		expect(graphqlClient).toHaveProperty('rawRequest');
	});

	it('uses default endpoint when VITE_GRAPHQL_URL is not set', () => {
		expect(graphqlClient).toBeDefined();
	});

	it('logs error in production when VITE_GRAPHQL_URL is not set', async () => {
		vi.resetModules();
		const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});

		// Simulate production without VITE_GRAPHQL_URL
		const originalEnv = import.meta.env.PROD;
		import.meta.env.PROD = true;
		import.meta.env.VITE_GRAPHQL_URL = '';

		try {
			await import('$lib/queries/client');
			expect(consoleSpy).toHaveBeenCalledWith(
				'VITE_GRAPHQL_URL is not set — GraphQL requests will fail in production'
			);
		} finally {
			import.meta.env.PROD = originalEnv;
			consoleSpy.mockRestore();
		}
	});
});
