import { gql } from 'graphql-request';

export interface User {
	id: string;
	username: string;
}

export interface UsersResponse {
	users: User[];
}

export const LIST_USERS = gql`
	query ListUsers {
		users {
			id
			username
		}
	}
`;
