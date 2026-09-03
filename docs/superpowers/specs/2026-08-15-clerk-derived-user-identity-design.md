# Clerk-Derived User Identity — Design Spec

**Date:** 2026-08-15
**Branch:** feature/discovery-page (or a follow-up branch — see Implementation Notes)
**Status:** Approved, pending implementation plan

## Problem

The app's notion of "the current user" — used to attribute new content (`Add to Library`) and to filter `ActivityTable`'s perspective list — is currently driven entirely by a manual dropdown (`UserSelector.svelte`), disconnected from Clerk authentication. This was discovered mid-session while verifying the Discover page: `UserSelector` was fully built and tested but never mounted anywhere in the app, so Add to Library was unconditionally broken (`Error: No user selected`).

Wiring `UserSelector` into the header as a stopgap made Add to Library testable again, but the user has since clarified: `UserSelector` was only ever meant as a temporary placeholder. The actual design intent is for Clerk itself — sign in via username/password or Google OAuth — to be the source of "who am I," with the app resolving that automatically.

## Current State (already built, discovered during investigation)

The backend has most of the necessary infrastructure already in place:

- `backend/internal/adapters/auth/clerk_middleware.go` verifies Clerk Bearer tokens on every request, resolves the Clerk user ID to a local `users` row via `GetByClerkID`, and auto-creates/links a local user on first authenticated request if a webhook hasn't fired yet (using Clerk's profile for username/email, with email-based linking to pre-existing rows).
- It injects `domain.AuthenticatedUser{ID, ClerkID, Username, Email, Role}` into request context, retrievable via `auth.ForContext(ctx)` / `auth.RequireAuth(ctx)` (`backend/internal/adapters/auth/context.go`).
- `frontend/src/lib/queries/client.ts` already attaches `Authorization: Bearer <token>` (via `window.Clerk?.session?.getToken()`) to every GraphQL request when signed in.
- `CreatePerspective`'s resolver already has the correct pattern for deriving a user ID from the authenticated session when the client doesn't supply one: `if userID == 0 { authUser, err := auth.RequireAuth(ctx); ...; userID = authUser.ID }`. `CreateContentFromYouTube` does not follow this pattern yet — it trusts client-supplied `input.UserID` outright, which is both the bug (breaks when no UI sets a userId) and a real security gap (any authenticated caller could currently attribute content to any other user's ID).
- One existing test user (`CodeWarrior-debug`, id=1) is already fully linked to a real Clerk account (has both `email` and `clerk_user_id` populated) — confirms the linkage path works end-to-end today.

What's missing: a way for the frontend to ask "who does my Clerk session map to" (no `me`/`currentUser` query exists), and a mechanism to keep the app's existing `selectedUserId` concept (`frontend/src/lib/stores/userSelection.svelte.ts`) in sync with that, instead of requiring manual selection.

## Design

### 1. Backend — `me` query

Add to `schema.graphql`, under `type Query`:

```graphql
me: User @auth
```

Resolver in `schema.resolvers.go` uses the existing `auth.RequireAuth(ctx)` helper to pull the request-scoped `AuthenticatedUser` and map it to `model.User`. No new domain/service logic — this exposes data the Clerk middleware already resolves on every authenticated request.

### 2. Backend — close the `CreateContentFromYouTube` security gap

Mirror `CreatePerspective`'s existing pattern exactly:

```go
userID := input.UserID
if userID == 0 {
    authUser, err := auth.RequireAuth(ctx)
    if err != nil {
        return nil, fmt.Errorf("access denied: authentication required")
    }
    userID = authUser.ID
}
```

Schema keeps `userId: IntID!` on `CreateContentFromYouTubeInput` (no breaking change — `CreatePerspectiveInput` already uses this same required-but-0-means-derive convention). Frontend sends `0` as the "derive from my session" sentinel.

### 3. Frontend — `AuthUserSync.svelte`

New, small, non-rendering component, mounted inside `<ClerkLoaded>` in `frontend/src/routes/+layout.svelte`, alongside `<Header />`.

- Uses `useClerkContext()` (existing `svelte-clerk` export) to reactively read `.auth.userId` (Clerk's ID — distinct from our local numeric user ID).
- A `$effect` watches that value:
  - **Transitions to a value (sign-in, or account switch on a shared device):** calls `queryClient.clear()` to wipe *all* cached queries (not just `me` — `content.lists()`, `perspectives.listByUser()`, etc. are also user-scoped and must not leak across accounts on a shared device), then runs the `me` query (`createQuery`, key `['me', clerkUserId]`, `staleTime` ~5 min) and calls `setSelectedUserId(me.id)` on success.
  - **Transitions to `null` (sign-out):** calls `queryClient.clear()` and `clearUserSelection()`.
- No loading UI for the gap between sign-in and `me` resolving — consumers (`ActivityTable`, `useAddVideo`) already guard on `selectedUserId !== null` / handle the unauthenticated case, so this is invisible in the common case. `localStorage`-backed `selectedUserId` also means a normal page reload doesn't show a "no user" flash — the last-known value is available synchronously before Clerk finishes loading, and gets reconciled (no-op if same account, corrected if not) once `AuthUserSync`'s effect runs.

### 4. Frontend — unmount, don't delete

- Remove `<UserSelector />` from `Header.svelte` (reverts this session's stopgap wiring).
- Revert the `Header.test.ts` mock addition for `UserSelector`.
- Leave `UserSelector.svelte`, `CreateUserPopover.svelte`, `useCreateUser.ts`, and their test files in place, untouched. Dead code for now — available if a future "browse others' perspectives" or admin flow wants them (see Deferred Scope below).
- Leave the `createUser` GraphQL mutation (schema field + resolver) in place — unused by any UI once `CreateUserPopover` is unmounted, but harmless, and removing+re-adding schema surface is more churn than it's worth for this pass.

### 5. Testing

- **Backend:** unit test for the `me` resolver (authenticated → returns mapped user; unauthenticated → error). Add/update a test for `CreateContentFromYouTube`'s new 0-sentinel derivation, following whatever pattern `CreatePerspective`'s equivalent test already uses.
- **Frontend:** new `AuthUserSync.test.ts` covering: sign-in → `me` fetched → `setSelectedUserId` called with the right ID → `queryClient.clear()` called; sign-out → `clearUserSelection()` called → `queryClient.clear()` called; no premature sync while Clerk is still loading (`isLoaded === false`); account switch (userId changes from one non-null value to another) → clears and re-syncs.
- `useAddVideo.ts`'s existing "no user selected" throw and its associated test cases get removed, replaced with the 0-sentinel pass-through (mutation always sends `userId: 0`, relying on backend derivation — matches `CreatePerspective`'s existing frontend convention if one exists, otherwise this establishes it for the content-creation path).

## Deferred Scope (explicitly out of this pass)

`ActivityTable` currently lets a user browse *any* other user's perspectives by picking them in `UserSelector` (`listByUser(selectedUserId)` for whoever's selected, not just "yourself"). Once `UserSelector` is unmounted, `selectedUserId` will always equal your own Clerk-derived identity — the ability to browse others' data is lost as a side effect of this change.

This is a known, accepted regression for this pass. A future "browse/view-as" feature — respecting the existing `Privacy` field so only non-private perspectives are visible to non-owners — would reintroduce this as its own, separately-designed capability. Not designed here.

## Existing Data Note

The four other pre-existing test users (`jimijordan`, `Kiraluvmusic`, `nana`, `ramulus`) have no `clerk_user_id` linkage yet. They'll link automatically the first time someone signs in via Clerk with a matching email on file (existing `clerk_middleware.go` logic), or remain unlinked/orphaned test data otherwise. Not a blocker — confirmed acceptable, no migration step needed.

## Implementation Notes

- This branches off `feature/discovery-page` context but is functionally unrelated to the Discover page — worth considering whether it lands as a separate branch/PR from the Discover page work once that's merged, to keep PR scope clean. Not decided here; a call for whoever runs `writing-plans` / executes.
