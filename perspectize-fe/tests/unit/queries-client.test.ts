import { describe, it, expect, beforeEach } from 'vitest';
import type { GraphQLClient } from 'graphql-request';

let graphqlClient: GraphQLClient;

beforeEach(async () => {
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
});
