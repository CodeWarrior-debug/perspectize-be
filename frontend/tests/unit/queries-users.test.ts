import { describe, it, expect } from 'vitest';
import { LIST_USERS, type User, type UsersResponse } from '$lib/queries/users';

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
});
