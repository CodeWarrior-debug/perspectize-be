---
created: 2026-03-11T01:22:12.834Z
title: Remove UserSelector and use authenticated Clerk user
area: auth
files:
  - frontend/src/lib/components/UserSelector.svelte
  - frontend/src/lib/stores/userSelection.svelte.ts
  - frontend/src/lib/components/ActivityTable.svelte:70
  - frontend/src/lib/components/ActivityTable.svelte:586
  - frontend/src/lib/queries/hooks/useCreatePerspective.ts
  - frontend/src/lib/queries/hooks/useAddVideo.ts
  - frontend/src/lib/components/PerspectivePopover.svelte:41
---

## Problem

The UserSelector dropdown in the header is a pre-auth artifact from before Clerk authentication was integrated. It requires users to manually pick themselves from a list of all users before they can add perspectives or videos. With Clerk auth now working, `selectedUserId` is often `null` (no selection), which falls back to `0` and causes "user_id must be a positive integer" errors.

A backend workaround was added (commit a27081e) to use the authenticated user when `userID` is 0 in `createPerspective`, but the frontend still has the unnecessary dropdown and store.

## Solution

1. Remove `UserSelector.svelte` component from the header
2. Remove `userSelection.svelte.ts` store (selectedUserId, getSelectedUserId, setSelectedUserId)
3. Update `ActivityTable.svelte` to get the authenticated user's local ID from the backend (e.g., a `me` query) instead of from the dropdown store
4. Update `PerspectivePopover` — either pass `userId: 0` (backend falls back to auth user) or fetch the local user ID via a `me` query
5. Update `useAddVideo` hook similarly
6. Consider adding a `me` GraphQL query that returns the authenticated user's profile
7. Remove `CreateUserPopover` component (manual user creation is no longer needed)
