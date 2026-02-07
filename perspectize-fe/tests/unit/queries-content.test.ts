import { describe, it, expect } from 'vitest';
import { LIST_CONTENT, GET_CONTENT, CREATE_CONTENT_FROM_YOUTUBE } from '$lib/queries/content';

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

		it('includes sorting parameters', () => {
			expect(LIST_CONTENT).toContain('$sortBy: ContentSortBy');
			expect(LIST_CONTENT).toContain('$sortOrder: SortOrder');
		});

		it('includes includeTotalCount parameter', () => {
			expect(LIST_CONTENT).toContain('$includeTotalCount: Boolean');
		});

		it('requests expected content fields', () => {
			expect(LIST_CONTENT).toContain('id');
			expect(LIST_CONTENT).toContain('name');
			expect(LIST_CONTENT).toContain('url');
			expect(LIST_CONTENT).toContain('contentType');
			expect(LIST_CONTENT).toContain('createdAt');
			expect(LIST_CONTENT).toContain('updatedAt');
		});

		it('requests length and lengthUnits fields', () => {
			expect(LIST_CONTENT).toContain('length');
			expect(LIST_CONTENT).toContain('lengthUnits');
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

	describe('CREATE_CONTENT_FROM_YOUTUBE', () => {
		it('is defined and is a string', () => {
			expect(CREATE_CONTENT_FROM_YOUTUBE).toBeDefined();
			expect(typeof CREATE_CONTENT_FROM_YOUTUBE).toBe('string');
		});

		it('contains the CreateContentFromYouTube operation name', () => {
			expect(CREATE_CONTENT_FROM_YOUTUBE).toContain('CreateContentFromYouTube');
		});

		it('is a mutation operation', () => {
			expect(CREATE_CONTENT_FROM_YOUTUBE).toContain('mutation');
		});

		it('takes CreateContentFromYouTubeInput input type', () => {
			expect(CREATE_CONTENT_FROM_YOUTUBE).toContain('$input: CreateContentFromYouTubeInput!');
		});

		it('calls createContentFromYouTube mutation', () => {
			expect(CREATE_CONTENT_FROM_YOUTUBE).toContain('createContentFromYouTube(input: $input)');
		});

		it('requests essential content fields', () => {
			expect(CREATE_CONTENT_FROM_YOUTUBE).toContain('id');
			expect(CREATE_CONTENT_FROM_YOUTUBE).toContain('name');
			expect(CREATE_CONTENT_FROM_YOUTUBE).toContain('url');
			expect(CREATE_CONTENT_FROM_YOUTUBE).toContain('contentType');
		});

		it('requests YouTube-specific metadata fields', () => {
			expect(CREATE_CONTENT_FROM_YOUTUBE).toContain('viewCount');
			expect(CREATE_CONTENT_FROM_YOUTUBE).toContain('likeCount');
			expect(CREATE_CONTENT_FROM_YOUTUBE).toContain('commentCount');
		});

		it('requests length fields', () => {
			expect(CREATE_CONTENT_FROM_YOUTUBE).toContain('length');
			expect(CREATE_CONTENT_FROM_YOUTUBE).toContain('lengthUnits');
		});

		it('requests createdAt timestamp', () => {
			expect(CREATE_CONTENT_FROM_YOUTUBE).toContain('createdAt');
		});
	});
});
