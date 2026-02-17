import { GraphQLClient } from 'graphql-request';

const GRAPHQL_ENDPOINT = import.meta.env.VITE_GRAPHQL_URL || 'http://localhost:8080/graphql';

if (!import.meta.env.VITE_GRAPHQL_URL && import.meta.env.PROD) {
	console.error(
		'VITE_GRAPHQL_URL is not set — GraphQL requests will fail in production.',
		'Set VITE_GRAPHQL_URL as a BUILD_TIME environment variable in your deployment platform.',
	);
}

console.debug('[GraphQL] endpoint:', GRAPHQL_ENDPOINT);

export const graphqlClient = new GraphQLClient(GRAPHQL_ENDPOINT, {
	headers: {},
});
