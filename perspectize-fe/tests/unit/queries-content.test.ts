import { describe, it, expect } from 'vitest';
import { LIST_CONTENT, GET_CONTENT } from '$lib/queries/content';

describe('GraphQL query definitions', () => {
	describe('LIST_CONTENT', () => {
		it('is defined and is a string', () => {
			expect(LIST_CONTENT).toBeDefined();
			expect(typeof LIST_CONTENT).toBe('string');
		});

		it('contains the ListContent operation name', () => {
			expect(LIST_CONTENT).toContain('ListContent');
		});

		it('queries content with pagination params', () => {
			expect(LIST_CONTENT).toContain('$first: Int');
			expect(LIST_CONTENT).toContain('$after: String');
		});

		it('requests expected content fields', () => {
			expect(LIST_CONTENT).toContain('id');
			expect(LIST_CONTENT).toContain('name');
			expect(LIST_CONTENT).toContain('url');
			expect(LIST_CONTENT).toContain('contentType');
			expect(LIST_CONTENT).toContain('createdAt');
			expect(LIST_CONTENT).toContain('updatedAt');
		});

		it('includes pageInfo for cursor pagination', () => {
			expect(LIST_CONTENT).toContain('pageInfo');
			expect(LIST_CONTENT).toContain('hasNextPage');
			expect(LIST_CONTENT).toContain('endCursor');
		});

		it('includes totalCount', () => {
			expect(LIST_CONTENT).toContain('totalCount');
		});
	});

	describe('GET_CONTENT', () => {
		it('is defined and is a string', () => {
			expect(GET_CONTENT).toBeDefined();
			expect(typeof GET_CONTENT).toBe('string');
		});

		it('contains the GetContent operation name', () => {
			expect(GET_CONTENT).toContain('GetContent');
		});

		it('takes an id parameter', () => {
			expect(GET_CONTENT).toContain('$id: ID!');
		});

		it('uses contentByID query', () => {
			expect(GET_CONTENT).toContain('contentByID');
		});

		it('requests detailed content fields', () => {
			expect(GET_CONTENT).toContain('length');
			expect(GET_CONTENT).toContain('viewCount');
			expect(GET_CONTENT).toContain('likeCount');
			expect(GET_CONTENT).toContain('commentCount');
			expect(GET_CONTENT).toContain('response');
		});
	});
});
