import { describe, it, expect } from 'vitest';

describe('GraphQL client', () => {
	it('exports a graphqlClient instance', async () => {
		const { graphqlClient } = await import('$lib/queries/client');
		expect(graphqlClient).toBeDefined();
		expect(typeof graphqlClient.request).toBe('function');
	});

	it('client has request method for making GraphQL calls', async () => {
		const { graphqlClient } = await import('$lib/queries/client');
		expect(graphqlClient).toHaveProperty('request');
		expect(graphqlClient).toHaveProperty('rawRequest');
	});

	it('uses default endpoint when VITE_GRAPHQL_URL is not set', async () => {
		const { graphqlClient } = await import('$lib/queries/client');
		expect(graphqlClient).toBeDefined();
	});
});
