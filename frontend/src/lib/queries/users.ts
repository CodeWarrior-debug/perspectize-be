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

export interface CreateUserInput {
	username: string;
	email?: string;
}

export interface CreateUserResponse {
	createUser: User;
}

export const CREATE_USER = gql`
	mutation CreateUser($input: CreateUserInput!) {
		createUser(input: $input) {
			id
			username
		}
	}
`;
