import { GraphQLClient } from 'graphql-request';

const GRAPHQL_ENDPOINT = import.meta.env.VITE_GRAPHQL_URL || 'http://localhost:8080/graphql';

export const graphqlClient = new GraphQLClient(GRAPHQL_ENDPOINT, {
	headers: {}
});
