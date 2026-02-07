import { gql } from 'graphql-request';

export const LIST_CONTENT = gql`
	query ListContent($first: Int, $after: String) {
		content(first: $first, after: $after) {
			edges {
				node {
					id
					title
					url
					contentType
					createdAt
					updatedAt
				}
				cursor
			}
			pageInfo {
				hasNextPage
				endCursor
			}
		}
	}
`;

export const GET_CONTENT = gql`
	query GetContent($id: IntID!) {
		content(id: $id) {
			id
			title
			url
			contentType
			description
			thumbnailUrl
			duration
		}
	}
`;
