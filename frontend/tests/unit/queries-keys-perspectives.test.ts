import { describe, expect, it } from 'vitest';
import { queryKeys } from '$lib/queries/keys';

describe('queryKeys.perspectives namespace', () => {
	describe('all()', () => {
		it('returns base perspectives key', () => {
			expect(queryKeys.perspectives.all()).toEqual(['app', 'perspectives']);
		});
	});

	describe('lists()', () => {
		it('returns perspectives list prefix', () => {
			expect(queryKeys.perspectives.lists()).toEqual(['app', 'perspectives', 'list']);
		});

		it('builds on perspectives.all()', () => {
			const allKey = queryKeys.perspectives.all();
			expect(queryKeys.perspectives.lists()).toEqual([...allKey, 'list']);
		});
	});

	describe('listByUser(userId)', () => {
		it('returns key with userId object when passed a user ID', () => {
			const result = queryKeys.perspectives.listByUser(42);
			expect(result).toEqual(['app', 'perspectives', 'list', { userId: 42 }]);
		});

		it('builds on perspectives.lists()', () => {
			const listsKey = queryKeys.perspectives.lists();
			const result = queryKeys.perspectives.listByUser(123);
			expect(result).toEqual([...listsKey, { userId: 123 }]);
		});

		it('handles different user IDs independently', () => {
			const user1 = queryKeys.perspectives.listByUser(1);
			const user2 = queryKeys.perspectives.listByUser(2);
			expect(user1).toEqual(['app', 'perspectives', 'list', { userId: 1 }]);
			expect(user2).toEqual(['app', 'perspectives', 'list', { userId: 2 }]);
			expect(user1).not.toEqual(user2);
		});
	});

	describe('details()', () => {
		it('returns perspectives detail prefix', () => {
			expect(queryKeys.perspectives.details()).toEqual(['app', 'perspectives', 'detail']);
		});

		it('builds on perspectives.all()', () => {
			const allKey = queryKeys.perspectives.all();
			expect(queryKeys.perspectives.details()).toEqual([...allKey, 'detail']);
		});
	});

	describe('detail(id)', () => {
		it('returns key with perspective id', () => {
			const result = queryKeys.perspectives.detail('perspective-123');
			expect(result).toEqual(['app', 'perspectives', 'detail', 'perspective-123']);
		});

		it('builds on perspectives.details()', () => {
			const detailsKey = queryKeys.perspectives.details();
			const result = queryKeys.perspectives.detail('456');
			expect(result).toEqual([...detailsKey, '456']);
		});

		it('handles different perspective IDs independently', () => {
			const p1 = queryKeys.perspectives.detail('id-1');
			const p2 = queryKeys.perspectives.detail('id-2');
			expect(p1).toEqual(['app', 'perspectives', 'detail', 'id-1']);
			expect(p2).toEqual(['app', 'perspectives', 'detail', 'id-2']);
			expect(p1).not.toEqual(p2);
		});
	});

	describe('hierarchical prefix matching', () => {
		it('listByUser keys start with lists() prefix', () => {
			const listsPrefix = queryKeys.perspectives.lists();
			const userKey = queryKeys.perspectives.listByUser(99);
			expect(userKey.slice(0, listsPrefix.length)).toEqual(listsPrefix);
		});

		it('detail keys start with details() prefix', () => {
			const detailsPrefix = queryKeys.perspectives.details();
			const singleDetail = queryKeys.perspectives.detail('test-id');
			expect(singleDetail.slice(0, detailsPrefix.length)).toEqual(detailsPrefix);
		});

		it('all perspectives keys start with perspectives.all() prefix', () => {
			const perspectivesPrefix = queryKeys.perspectives.all();
			expect(queryKeys.perspectives.lists()[0]).toEqual(perspectivesPrefix[0]);
			expect(queryKeys.perspectives.lists()[1]).toEqual(perspectivesPrefix[1]);
		});
	});

	describe('type safety (as const)', () => {
		it('returns readonly arrays', () => {
			const key = queryKeys.perspectives.all();
			expect(Array.isArray(key)).toBe(true);
		});

		it('listByUser returns readonly array with object', () => {
			const key = queryKeys.perspectives.listByUser(1);
			expect(Array.isArray(key)).toBe(true);
			expect(typeof key[3]).toBe('object');
		});
	});
});
