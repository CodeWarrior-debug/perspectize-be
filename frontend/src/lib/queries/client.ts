import { GraphQLClient } from 'graphql-request';

const GRAPHQL_ENDPOINT = import.meta.env.VITE_GRAPHQL_URL || 'http://localhost:8080/graphql';

if (!import.meta.env.VITE_GRAPHQL_URL && import.meta.env.PROD) {
	console.error('VITE_GRAPHQL_URL is not set — GraphQL requests will fail in production');
}

export const graphqlClient = new GraphQLClient(GRAPHQL_ENDPOINT, {
	headers: {}
});
