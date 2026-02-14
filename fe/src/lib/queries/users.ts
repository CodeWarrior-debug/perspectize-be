import { gql } from 'graphql-request';

export const LIST_USERS = gql`
	query ListUsers {
		users {
			id
			username
			email
		}
	}
`;
