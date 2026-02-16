# Design: Optimize YouTube Add Flow (Phase 07.5)

**Date:** 2026-02-16
**Status:** Approved

## Problem

Two issues with the current "create from YouTube + refresh" flow:

1. **Mutation returns wrong fields.** `CREATE_CONTENT_FROM_YOUTUBE` returns 10 fields (`id`, `name`, `url`, `contentType`, `length`, `lengthUnits`, `viewCount`, `likeCount`, `commentCount`, `createdAt`) but only `name` is used (for the toast). The response includes `commentCount` (not shown in the table) but is missing `channelTitle`, `publishedAt`, `tags`, `description`, `updatedAt` (all displayed in ActivityTable).

2. **Full list refetch after adding one video.** After mutation succeeds, `useAddVideo` calls `queryClient.invalidateQueries({ queryKey: queryKeys.content.lists() })` which re-fetches the entire paginated page (10-50 items) from the backend. This is a redundant round-trip — we already have the new item's data from the mutation response.

## Approach

**Direct cache insertion via `setQueryData`** — the mutation response is shaped to match `ContentItem`, then prepended into the active list cache. The list query is marked stale for eventual consistency but not immediately refetched.

### Why this approach

- **No extra network request.** The new row appears instantly from mutation response data.
- **Eventual consistency.** Background stale marking ensures the cache syncs on next natural refetch (focus, interval, navigation).
- **Simple.** No optimistic rollback logic needed — we insert only after mutation succeeds.
- **Backend unchanged.** The resolver already returns the full `Content` type; we just need to request the right fields in the GraphQL selection set.

### Rejected alternatives

| Alternative | Why rejected |
|-------------|-------------|
| Optimistic update (pre-mutation) | Adds rollback complexity for a flow that takes <2s. Overkill. |
| Keep full refetch, just trim response | Doesn't solve the core problem (unnecessary round-trip). |
| Insert without stale marking | Cache could drift if server enriches data beyond what mutation returns. |

## Design

### 1. Reshape mutation selection set

Align `CREATE_CONTENT_FROM_YOUTUBE` fields with `ContentItem` interface:

**Remove:** `commentCount` (not in ActivityTable)
**Add:** `channelTitle`, `publishedAt`, `tags`, `description`, `updatedAt`

```graphql
mutation CreateContentFromYouTube($input: CreateContentFromYouTubeInput!) {
  createContentFromYouTube(input: $input) {
    id
    name
    url
    contentType
    length
    lengthUnits
    viewCount
    likeCount
    channelTitle
    publishedAt
    tags
    description
    createdAt
    updatedAt
  }
}
```

### 2. Align `CreateContentResponse` type with `ContentItem`

```typescript
export interface CreateContentResponse {
  createContentFromYouTube: ContentItem;
}
```

This eliminates a separate interface that drifts from `ContentItem`.

### 3. Update `useAddVideo` hook — cache insertion

Replace `invalidateQueries` with `setQueryData` + stale marking:

```typescript
onSuccess: (data: CreateContentResponse) => {
  const newItem = data.createContentFromYouTube;
  toast.success(`Added: ${newItem.name ?? 'video'}`);

  // Insert new item into all active list caches
  queryClient.setQueriesData<ContentResponse>(
    { queryKey: queryKeys.content.lists() },
    (oldData) => {
      if (!oldData) return oldData;
      return {
        content: {
          ...oldData.content,
          items: [newItem, ...oldData.content.items],
          totalCount: (oldData.content.totalCount ?? 0) + 1,
        },
      };
    }
  );

  // Mark stale for background sync (no immediate refetch)
  queryClient.invalidateQueries({
    queryKey: queryKeys.content.lists(),
    refetchType: 'none',
  });
}
```

### 4. What stays the same

- Backend resolver — no changes needed
- ActivityTable query — unchanged, continues to observe the cache
- Error handling in `useAddVideo` — unchanged
- AddVideoPopover/AddVideoDialog — unchanged (they call the hook)

## Files Modified

| File | Change |
|------|--------|
| `frontend/src/lib/queries/content.ts` | Mutation GQL + `CreateContentResponse` type |
| `frontend/src/lib/queries/hooks/useAddVideo.ts` | `setQueriesData` instead of `invalidateQueries` |

## Testing

- Existing `useAddVideo` tests updated for new behavior
- Build and type-check pass
- Frontend test suite passes
