/**
 * Centralized query key factory for type-safe, hierarchical cache invalidation.
 */
export const queryKeys = {
	all: ['app'] as const,

	content: {
		all: () => [...queryKeys.all, 'content'] as const,
		lists: () => [...queryKeys.content.all(), 'list'] as const,
		list: (filters: {
			sortBy?: string;
			sortOrder?: string;
			search?: string;
			first?: number;
			after?: string | null;
		}) => [...queryKeys.content.lists(), filters] as const,
		details: () => [...queryKeys.content.all(), 'detail'] as const,
		detail: (id: string) => [...queryKeys.content.details(), id] as const,
	},

	users: {
		all: () => [...queryKeys.all, 'users'] as const,
		lists: () => [...queryKeys.users.all(), 'list'] as const,
		list: () => [...queryKeys.users.lists()] as const,
		details: () => [...queryKeys.users.all(), 'detail'] as const,
		detail: (id: string) => [...queryKeys.users.details(), id] as const,
	},
} as const;
