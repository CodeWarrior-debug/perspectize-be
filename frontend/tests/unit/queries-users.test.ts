import { describe, it, expect } from 'vitest';
import {
	LIST_USERS,
	CREATE_USER,
	type User,
	type UsersResponse,
	type CreateUserInput,
	type CreateUserResponse
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
});
