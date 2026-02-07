import { gql } from 'graphql-request';

export const LIST_CONTENT = gql`
	query ListContent(
		$first: Int
		$after: String
		$sortBy: ContentSortBy = UPDATED_AT
		$sortOrder: SortOrder = DESC
		$includeTotalCount: Boolean = true
	) {
		content(
			first: $first
			after: $after
			sortBy: $sortBy
			sortOrder: $sortOrder
			includeTotalCount: $includeTotalCount
		) {
			items {
				id
				name
				url
				contentType
				length
				lengthUnits
				createdAt
				updatedAt
			}
			pageInfo {
				hasNextPage
				endCursor
			}
			totalCount
		}
	}
`;

export const GET_CONTENT = gql`
	query GetContent($id: ID!) {
		contentByID(id: $id) {
			id
			name
			url
			contentType
			length
			lengthUnits
			viewCount
			likeCount
			commentCount
			response
			createdAt
			updatedAt
		}
	}
`;
