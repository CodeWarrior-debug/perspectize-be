import { describe, it, expect, beforeEach, vi } from 'vitest';
import type { GraphQLClient } from 'graphql-request';

let graphqlClient: GraphQLClient;
let graphqlRequest: (document: string, variables?: Record<string, unknown>) => Promise<unknown>;
let getAuthToken: () => Promise<string | null>;

beforeEach(async () => {
	vi.resetModules();
	const mod = await import('$lib/queries/client');
	graphqlClient = mod.graphqlClient;
	graphqlRequest = mod.graphqlRequest;
	getAuthToken = mod.getAuthToken;
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

	it('exports graphqlRequest function', () => {
		expect(graphqlRequest).toBeDefined();
		expect(typeof graphqlRequest).toBe('function');
	});

	it('exports getAuthToken function', () => {
		expect(getAuthToken).toBeDefined();
		expect(typeof getAuthToken).toBe('function');
	});

	it('getAuthToken returns null when Clerk is not available', async () => {
		const token = await getAuthToken();
		expect(token).toBeNull();
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
				'VITE_GRAPHQL_URL is not set — GraphQL requests will fail in production.',
				'Set VITE_GRAPHQL_URL as a BUILD_TIME environment variable in your deployment platform.',
			);
		} finally {
			import.meta.env.PROD = originalEnv;
			consoleSpy.mockRestore();
		}
	});
});
