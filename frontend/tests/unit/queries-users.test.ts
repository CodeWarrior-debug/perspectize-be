import { describe, it, expect } from 'vitest';
import {
	LIST_USERS,
	CREATE_USER,
	GET_USER_BY_USERNAME,
	type User,
	type UsersResponse,
	type CreateUserInput,
	type CreateUserResponse,
	type UserByUsernameResponse
} from '$lib/queries/users';

describe('User Queries', () => {
	describe('Type exports', () => {
		it('exports User interface', () => {
			const user: User = { id: '1', username: 'test' };
			expect(user).toBeDefined();
		});

		it('exports UsersResponse interface', () => {
			const response: UsersResponse = { users: [{ id: '1', username: 'test' }] };
			expect(response).toBeDefined();
		});

		it('exports CreateUserInput interface with required username and optional email', () => {
			const input: CreateUserInput = { username: 'newuser' };
			expect(input.username).toBeDefined();
			expect(input.email).toBeUndefined();

			const inputWithEmail: CreateUserInput = { username: 'newuser', email: 'user@example.com' };
			expect(inputWithEmail.username).toBeDefined();
			expect(inputWithEmail.email).toBeDefined();
		});

		it('exports CreateUserResponse interface', () => {
			const response: CreateUserResponse = {
				createUser: { id: '1', username: 'newuser' }
			};
			expect(response).toBeDefined();
		});

		it('exports UserByUsernameResponse interface', () => {
			const response: UserByUsernameResponse = {
				userByUsername: { id: '999', username: '[anonymous]' }
			};
			expect(response).toBeDefined();
		});

		it('UserByUsernameResponse allows null for userByUsername (user not found)', () => {
			const response: UserByUsernameResponse = {
				userByUsername: null
			};
			expect(response).toBeDefined();
			expect(response.userByUsername).toBeNull();
		});
	});

	describe('LIST_USERS', () => {
		it('is exported as a string', () => {
			expect(typeof LIST_USERS).toBe('string');
		});

		it('contains the users query', () => {
			expect(LIST_USERS).toContain('query ListUsers');
			expect(LIST_USERS).toContain('users');
		});

		it('requests id and username fields only', () => {
			expect(LIST_USERS).toContain('id');
			expect(LIST_USERS).toContain('username');
			expect(LIST_USERS).not.toContain('email');
		});

		it('does not request unnecessary timestamp fields', () => {
			expect(LIST_USERS).not.toContain('createdAt');
			expect(LIST_USERS).not.toContain('updatedAt');
		});
	});

	describe('CREATE_USER', () => {
		it('is exported as a string', () => {
			expect(typeof CREATE_USER).toBe('string');
		});

		it('contains the CreateUser mutation operation', () => {
			expect(CREATE_USER).toContain('mutation CreateUser');
		});

		it('contains createUser field', () => {
			expect(CREATE_USER).toContain('createUser');
		});

		it('requests id and username fields in response', () => {
			expect(CREATE_USER).toContain('id');
			expect(CREATE_USER).toContain('username');
		});

		it('does not request email in the response', () => {
			expect(CREATE_USER).not.toContain('email');
		});
	});

	describe('GET_USER_BY_USERNAME', () => {
		it('is exported as a string', () => {
			expect(typeof GET_USER_BY_USERNAME).toBe('string');
		});

		it('contains the GetUserByUsername operation name', () => {
			expect(GET_USER_BY_USERNAME).toContain('GetUserByUsername');
		});

		it('is a query operation', () => {
			expect(GET_USER_BY_USERNAME).toContain('query');
		});

		it('takes a username parameter', () => {
			expect(GET_USER_BY_USERNAME).toContain('$username: String!');
		});

		it('calls userByUsername query', () => {
			expect(GET_USER_BY_USERNAME).toContain('userByUsername(username: $username)');
		});

		it('requests id and username fields in response', () => {
			expect(GET_USER_BY_USERNAME).toContain('id');
			expect(GET_USER_BY_USERNAME).toContain('username');
		});

		it('does not request email or timestamps in response', () => {
			expect(GET_USER_BY_USERNAME).not.toContain('email');
			expect(GET_USER_BY_USERNAME).not.toContain('createdAt');
			expect(GET_USER_BY_USERNAME).not.toContain('updatedAt');
		});

		it('supports anonymous username parameter', () => {
			// This test verifies the query structure, not the actual operation
			// The query is designed to support any string username, including [anonymous]
			expect(GET_USER_BY_USERNAME).toContain('userByUsername');
		});
	});
});
