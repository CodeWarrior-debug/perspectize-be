# Clerk-Derived User Identity Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make Clerk authentication the source of "who am I" throughout the app — expose a `me` GraphQL query backed by the existing Clerk middleware, close the `CreateContentFromYouTube` authorization gap, and add a frontend component that keeps `selectedUserId` in sync with the Clerk session (replacing the manual `UserSelector` dropdown for the "who am I" purpose).

**Architecture:** Backend exposes a new `@auth`-gated `me: User` query that reads the `AuthenticatedUser` the Clerk middleware already put in context and maps it through the existing `UserService.GetByID` → `userDomainToModel` path (same pattern `CreatePerspective` already uses to derive `userID` when unset). `CreateContentFromYouTube` gets the same zero-sentinel derivation `CreatePerspective` already has. On the frontend, a new non-rendering `AuthUserSync.svelte` component watches Clerk's `auth.userId` via `useClerkContext()`, clears the TanStack Query cache and re-syncs `selectedUserId` on sign-in/sign-out/account-switch, and `UserSelector`/`CreateUserPopover` are left mounted nowhere (already true on this branch — no code change needed for that part).

**Tech Stack:** Go 1.25+, gqlgen (schema-first), sqlx/GORM (no schema/DB changes needed), Svelte 5 runes, `@tanstack/svelte-query` v6, `svelte-clerk` v1.

**Spec:** `docs/superpowers/specs/2026-08-15-clerk-derived-user-identity-design.md` — read alongside this plan; this plan argues from that spec's decisions and doesn't re-justify them.

## Global Constraints

- Branch: stay on `feature/discovery-page` (user decision — do not create a new branch).
- `CreateContentFromYouTubeInput.userId` stays `IntID!` (required-but-0-means-derive), matching `CreatePerspectiveInput.userID`'s existing convention. No GraphQL schema breaking change.
- Frontend Svelte 5 runes only (`$state`, `$derived`, `$effect`) — no Svelte 4 syntax (`export let`, `$:`).
- TanStack Query v5+/Svelte 5 function-wrapper pattern for `createQuery`/`createMutation` (per `frontend/CLAUDE.md`) — never pass an options object directly.
- After any `schema.graphql` edit: run `make graphql-gen` from `backend/` before touching resolver code.
- `Header.svelte` / `Header.test.ts` are already in their pre-stopgap state (no `UserSelector` import) as of this plan's start — confirmed via `git status` showing a clean tree. No task in this plan touches them.
- `UserSelector.svelte`, `CreateUserPopover.svelte`, `useCreateUser.ts`, and their tests are left untouched (dead code, kept for a future browse/admin flow per spec's Deferred Scope).

---

## File Structure

**Backend:**
- Modify: `backend/schema.graphql` — add `me: User @auth` query field
- Modify (generated): `backend/internal/adapters/graphql/generated/*.go`, `backend/internal/adapters/graphql/model/models_gen.go` — via `make graphql-gen`
- Modify: `backend/internal/adapters/graphql/resolvers/schema.resolvers.go` — implement `Me` resolver; fix `CreateContentFromYouTube` to derive `userID` when zero
- Create: `backend/test/resolvers/me_resolver_test.go` — `me` resolver tests + a no-auth test server helper
- Modify: `backend/test/resolvers/content_resolver_test.go` — add zero-sentinel derivation test for `CreateContentFromYouTube`

**Frontend:**
- Modify: `frontend/src/lib/queries/users.ts` — add `ME` query + `Me`/`MeResponse` types
- Create: `frontend/src/lib/components/AuthUserSync.svelte` — non-rendering Clerk↔selectedUserId sync component
- Modify: `frontend/src/routes/+layout.svelte` — mount `<AuthUserSync />` inside `<ClerkLoaded>`
- Create: `frontend/tests/components/AuthUserSync.test.ts`
- Modify: `frontend/src/lib/queries/hooks/useAddVideo.ts` — always send `userId: 0` (derive from session), drop the `getSelectedUserId`/"no user selected" throw
- Modify: `frontend/tests/unit/hooks-useAddVideo.test.ts` — drop the "no user selected" mapping test and its `mockGetSelectedUserId` scaffolding that's no longer exercised
- Modify: `frontend/tests/components/AddVideoPopover.test.ts` — drop "mutationFn throws when no user is selected" test, update the args-assertion test to expect `userId: 0`
- Modify: `frontend/tests/components/AddVideoDialog.test.ts` — same two changes as above

---

## Task 1: Backend — `me` query

**Files:**
- Modify: `backend/schema.graphql` (Query block, ~line 261-294)
- Modify: `backend/internal/adapters/graphql/resolvers/schema.resolvers.go` (new `Me` method, generated as a stub after codegen)
- Test: `backend/test/resolvers/me_resolver_test.go`

**Interfaces:**
- Consumes: `auth.RequireAuth(ctx) (*domain.AuthenticatedUser, error)` (`backend/internal/adapters/auth/context.go`), `r.UserService.GetByID(ctx, id int) (*domain.User, error)`, `userDomainToModel(u *domain.User) *model.User` (`backend/internal/adapters/graphql/resolvers/helpers.go:13`)
- Produces: GraphQL `me: User` query, resolvable by the frontend's `ME` query (Task 3)

- [ ] **Step 1: Add the schema field**

Edit `backend/schema.graphql`, in `type Query { ... }` (currently lines 261-294), add `me` as the first field:

```graphql
type Query {
  # Current authenticated user, derived from the Clerk session
  me: User @auth

  # Get single content by ID
  contentByID(id: ID!): Content
  ...
```

- [ ] **Step 2: Regenerate gqlgen code**

Run from `backend/`:

```bash
make graphql-gen
```

This adds a `Me(ctx context.Context) (*model.User, error)` stub (panicking with "not implemented") to the bottom of `schema.resolvers.go` and updates `generated/*.go`. Verify with:

```bash
grep -n "func (r \*queryResolver) Me(" backend/internal/adapters/graphql/resolvers/schema.resolvers.go
```

- [ ] **Step 3: Write the failing test**

Create `backend/test/resolvers/me_resolver_test.go`:

```go
package resolvers_test

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/CodeWarrior-debug/perspectize/backend/internal/adapters/graphql/directives"
	"github.com/CodeWarrior-debug/perspectize/backend/internal/adapters/graphql/generated"
	"github.com/CodeWarrior-debug/perspectize/backend/internal/adapters/graphql/resolvers"
	"github.com/CodeWarrior-debug/perspectize/backend/internal/core/domain"
	"github.com/CodeWarrior-debug/perspectize/backend/internal/core/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestServerWithUserRepo mirrors setupTestServer (content_resolver_test.go)
// but takes a caller-configured mockUserRepository, so tests can control what
// UserService.GetByID returns for the authenticated user's ID.
func setupTestServerWithUserRepo(userRepo *mockUserRepository) *httptest.Server {
	contentRepo := &mockContentRepository{}
	perspectiveRepo := &mockPerspectiveRepository{}
	contentService := services.NewContentService(contentRepo, &mockYouTubeClient{})
	userService := services.NewUserService(userRepo, contentRepo, perspectiveRepo)
	perspectiveService := services.NewPerspectiveService(perspectiveRepo, userRepo)
	resolver := resolvers.NewResolver(contentService, userService, perspectiveService)
	directiveRoot := directives.NewDirectiveRoot(contentService, perspectiveService)
	gqlConfig := generated.Config{
		Resolvers: resolver,
		Directives: generated.DirectiveRoot{
			Auth:  directiveRoot.Auth,
			Owner: directiveRoot.Owner,
		},
	}
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(gqlConfig))
	authHandler := injectAuthMiddleware(srv)
	return httptest.NewServer(authHandler)
}

// setupTestServerNoAuth mirrors setupTestServer but omits the auth-injecting
// middleware, so requests reach resolvers/directives as truly unauthenticated.
func setupTestServerNoAuth(userRepo *mockUserRepository) *httptest.Server {
	contentRepo := &mockContentRepository{}
	perspectiveRepo := &mockPerspectiveRepository{}
	contentService := services.NewContentService(contentRepo, &mockYouTubeClient{})
	userService := services.NewUserService(userRepo, contentRepo, perspectiveRepo)
	perspectiveService := services.NewPerspectiveService(perspectiveRepo, userRepo)
	resolver := resolvers.NewResolver(contentService, userService, perspectiveService)
	directiveRoot := directives.NewDirectiveRoot(contentService, perspectiveService)
	gqlConfig := generated.Config{
		Resolvers: resolver,
		Directives: generated.DirectiveRoot{
			Auth:  directiveRoot.Auth,
			Owner: directiveRoot.Owner,
		},
	}
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(gqlConfig))
	return httptest.NewServer(srv)
}

func TestMeQuery_Success(t *testing.T) {
	userRepo := &mockUserRepository{
		getByIDFn: func(ctx context.Context, id int) (*domain.User, error) {
			assert.Equal(t, 1, id) // injectAuthMiddleware authenticates as user ID 1
			return &domain.User{
				ID:       1,
				Username: "testuser",
				Email:    "test@example.com",
				Active:   true,
				Role:     domain.UserRoleDefault,
			}, nil
		},
	}

	server := setupTestServerWithUserRepo(userRepo)
	defer server.Close()

	result := executeGraphQL(t, server, `{ me { id username email } }`)

	assert.Empty(t, result.Errors)

	var data struct {
		Me struct {
			ID       string `json:"id"`
			Username string `json:"username"`
			Email    string `json:"email"`
		} `json:"me"`
	}
	err := json.Unmarshal(result.Data, &data)
	require.NoError(t, err)

	assert.Equal(t, "1", data.Me.ID)
	assert.Equal(t, "testuser", data.Me.Username)
	assert.Equal(t, "test@example.com", data.Me.Email)
}

func TestMeQuery_Unauthenticated_ReturnsError(t *testing.T) {
	server := setupTestServerNoAuth(&mockUserRepository{})
	defer server.Close()

	result := executeGraphQL(t, server, `{ me { id username } }`)

	require.NotEmpty(t, result.Errors)
	assert.Contains(t, result.Errors[0].Message, "access denied: authentication required")
}

func TestMeQuery_UserNotFound_ReturnsError(t *testing.T) {
	userRepo := &mockUserRepository{
		getByIDFn: func(ctx context.Context, id int) (*domain.User, error) {
			return nil, domain.ErrNotFound
		},
	}

	server := setupTestServerWithUserRepo(userRepo)
	defer server.Close()

	result := executeGraphQL(t, server, `{ me { id username } }`)

	require.NotEmpty(t, result.Errors)
	assert.Contains(t, result.Errors[0].Message, "user not found")
}
```

- [ ] **Step 4: Run tests to verify they fail**

```bash
cd backend && go test ./test/resolvers/... -run TestMeQuery -v
```

Expected: compile error or FAIL — `Me` resolver still panics with "not implemented".

- [ ] **Step 5: Implement the resolver**

In `backend/internal/adapters/graphql/resolvers/schema.resolvers.go`, find the `Me` stub gqlgen appended (bottom of file) and replace its body:

```go
// Me is the resolver for the me field.
func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
	authUser, err := auth.RequireAuth(ctx)
	if err != nil {
		return nil, fmt.Errorf("access denied: authentication required")
	}

	user, err := r.UserService.GetByID(ctx, authUser.ID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		slog.Error("getting current user failed", "userID", authUser.ID, "error", err)
		return nil, fmt.Errorf("failed to get current user")
	}

	return userDomainToModel(user), nil
}
```

(`auth`, `errors`, `fmt`, `slog`, `domain` are all already imported at the top of this file.)

- [ ] **Step 6: Run tests to verify they pass**

```bash
cd backend && go test ./test/resolvers/... -run TestMeQuery -v
```

Expected: PASS (all 3 subtests).

- [ ] **Step 7: Run the full backend test suite and build**

```bash
cd backend && go build ./...
cd backend && go test ./...
```

Expected: zero errors, all tests pass (this also catches any gqlgen-regeneration fallout in unrelated resolvers).

- [ ] **Step 8: Commit**

```bash
git add backend/schema.graphql backend/internal/adapters/graphql/generated backend/internal/adapters/graphql/model backend/internal/adapters/graphql/resolvers/schema.resolvers.go backend/test/resolvers/me_resolver_test.go
git commit -m "feat(auth): add me query to resolve current user from Clerk session"
```

---

## Task 2: Backend — close `CreateContentFromYouTube` security gap

**Files:**
- Modify: `backend/internal/adapters/graphql/resolvers/schema.resolvers.go:23-49`
- Test: `backend/test/resolvers/content_resolver_test.go`

**Interfaces:**
- Consumes: `auth.RequireAuth(ctx)` (same helper as Task 1); `domain.Content.AddedByUserID int` (`backend/internal/core/domain/content.go:22`) — the field the derived `userID` ends up on, via `ContentService.CreateFromYouTube(ctx, url, userID)`.
- Produces: no new exported surface — this is a behavior fix within the existing mutation.

- [ ] **Step 1: Write the failing test**

Add to `backend/test/resolvers/content_resolver_test.go`, near the other `CreateContentFromYouTube` tests (after `TestCreateContentFromYouTube_Success`, ~line 447):

```go
func TestCreateContentFromYouTube_DerivesUserIDFromSession_WhenZero(t *testing.T) {
	metadata := &portservices.VideoMetadata{
		Title:       "Amazing Video",
		Description: "Description",
		Duration:    600,
		ChannelName: "Channel",
		Response:    json.RawMessage(`{"items":[]}`),
	}

	var capturedUserID int
	repo := &mockContentRepository{
		getByURLFn: func(ctx context.Context, url string) (*domain.Content, error) {
			return nil, domain.ErrNotFound
		},
		getOrCreateByURLFn: func(ctx context.Context, content *domain.Content) (*domain.Content, bool, error) {
			capturedUserID = content.AddedByUserID
			content.ID = 42
			return content, false, nil
		},
	}

	ytClient := &mockYouTubeClient{
		getVideoMetadataFn: func(ctx context.Context, videoID string) (*portservices.VideoMetadata, error) {
			return metadata, nil
		},
	}

	server := setupTestServer(repo, ytClient)
	defer server.Close()

	// userId: 0 is the "derive from my session" sentinel. setupTestServer's
	// injectAuthMiddleware authenticates as user ID 1.
	result := executeGraphQL(t, server, `mutation { createContentFromYouTube(input: { url: "https://www.youtube.com/watch?v=dQw4w9WgXcQ", userId: 0 }) { content { id } alreadyExisted } }`)

	assert.Empty(t, result.Errors)
	assert.Equal(t, 1, capturedUserID)
}
```

- [ ] **Step 2: Run test to verify it fails**

```bash
cd backend && go test ./test/resolvers/... -run TestCreateContentFromYouTube_DerivesUserIDFromSession -v
```

Expected: FAIL — `capturedUserID` is `0`, not `1` (resolver currently trusts `input.UserID` verbatim).

- [ ] **Step 3: Fix the resolver**

In `backend/internal/adapters/graphql/resolvers/schema.resolvers.go`, replace lines 23-26:

```go
// CreateContentFromYouTube is the resolver for the createContentFromYouTube field.
func (r *mutationResolver) CreateContentFromYouTube(ctx context.Context, input model.CreateContentFromYouTubeInput) (*model.CreateContentResult, error) {
	content, err := r.ContentService.CreateFromYouTube(ctx, input.URL, input.UserID)
```

with:

```go
// CreateContentFromYouTube is the resolver for the createContentFromYouTube field.
func (r *mutationResolver) CreateContentFromYouTube(ctx context.Context, input model.CreateContentFromYouTubeInput) (*model.CreateContentResult, error) {
	// Use authenticated user when userID is not provided or zero (mirrors CreatePerspective)
	userID := input.UserID
	if userID == 0 {
		authUser, err := auth.RequireAuth(ctx)
		if err != nil {
			return nil, fmt.Errorf("access denied: authentication required")
		}
		userID = authUser.ID
	}

	content, err := r.ContentService.CreateFromYouTube(ctx, input.URL, userID)
```

- [ ] **Step 4: Run test to verify it passes**

```bash
cd backend && go test ./test/resolvers/... -run TestCreateContentFromYouTube -v
```

Expected: PASS for all `TestCreateContentFromYouTube_*` tests (the new one plus the 3 pre-existing ones, which still pass explicit non-zero `userId` and are unaffected).

- [ ] **Step 5: Run the full backend test suite and build**

```bash
cd backend && go build ./...
cd backend && go test ./...
```

Expected: zero errors, all tests pass.

- [ ] **Step 6: Commit**

```bash
git add backend/internal/adapters/graphql/resolvers/schema.resolvers.go backend/test/resolvers/content_resolver_test.go
git commit -m "fix(content): derive userID from Clerk session in createContentFromYouTube when unset"
```

---

## Task 3: Frontend — `ME` query

**Files:**
- Modify: `frontend/src/lib/queries/users.ts`

**Interfaces:**
- Produces: `ME` (gql tagged template), `Me` interface `{ id: string; username: string }`, `MeResponse` interface `{ me: Me | null }` — consumed by `AuthUserSync.svelte` (Task 4)

- [ ] **Step 1: Add the query and types**

Append to `frontend/src/lib/queries/users.ts`:

```ts
export interface Me {
	id: string;
	username: string;
}

export interface MeResponse {
	me: Me | null;
}

export const ME = gql`
	query Me {
		me {
			id
			username
		}
	}
`;
```

- [ ] **Step 2: Verify the frontend still type-checks**

```bash
cd frontend && pnpm run check
```

Expected: no new errors.

- [ ] **Step 3: Commit**

```bash
git add frontend/src/lib/queries/users.ts
git commit -m "feat(auth): add me query for frontend Clerk identity sync"
```

---

## Task 4: Frontend — `AuthUserSync.svelte`

**Files:**
- Create: `frontend/src/lib/components/AuthUserSync.svelte`
- Modify: `frontend/src/routes/+layout.svelte`
- Test: `frontend/tests/components/AuthUserSync.test.ts`

**Interfaces:**
- Consumes: `ME`, `MeResponse` (Task 3, `$lib/queries/users`); `graphqlRequest` (`$lib/queries/client`); `setSelectedUserId`, `clearUserSelection` (`$lib/stores/userSelection.svelte`, already exist — `frontend/src/lib/stores/userSelection.svelte.ts`); `useClerkContext` (`svelte-clerk`, returns `{ isLoaded: boolean; auth: { userId: string | null | undefined; ... }; ... }`)
- Produces: `AuthUserSync` component with no props, no rendered output — mounted once in `+layout.svelte`

- [ ] **Step 1: Write the failing test**

Create `frontend/tests/components/AuthUserSync.test.ts`:

```ts
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render } from '@testing-library/svelte';
import AuthUserSync from '$lib/components/AuthUserSync.svelte';

const { mockSetSelectedUserId, mockClearUserSelection, mockQueryClientClear, mockClerkContext, mockQueryState } =
	vi.hoisted(() => ({
		mockSetSelectedUserId: vi.fn(),
		mockClearUserSelection: vi.fn(),
		mockQueryClientClear: vi.fn(),
		mockClerkContext: {
			isLoaded: true,
			auth: { userId: null as string | null },
		},
		mockQueryState: {
			isLoading: false,
			error: null as Error | null,
			data: null as { me: { id: string; username: string } } | null,
		},
	}));

let capturedQueryOptions: any;

vi.mock('@tanstack/svelte-query', () => ({
	createQuery: vi.fn((optionsFn: () => any) => {
		capturedQueryOptions = optionsFn();
		return mockQueryState;
	}),
	useQueryClient: vi.fn(() => ({
		clear: mockQueryClientClear,
	})),
}));

vi.mock('svelte-clerk', () => ({
	useClerkContext: vi.fn(() => mockClerkContext),
}));

vi.mock('$lib/queries/client', () => ({
	graphqlRequest: vi.fn(),
}));

vi.mock('$lib/queries/users', () => ({
	ME: 'mock-me-query',
}));

vi.mock('$lib/stores/userSelection.svelte', () => ({
	setSelectedUserId: (...args: unknown[]) => mockSetSelectedUserId(...args),
	clearUserSelection: () => mockClearUserSelection(),
}));

describe('AuthUserSync', () => {
	beforeEach(() => {
		vi.clearAllMocks();
		capturedQueryOptions = undefined;
		mockClerkContext.isLoaded = true;
		mockClerkContext.auth.userId = null;
		mockQueryState.isLoading = false;
		mockQueryState.error = null;
		mockQueryState.data = null;
	});

	it('does not sync while Clerk is still loading', () => {
		mockClerkContext.isLoaded = false;
		mockClerkContext.auth.userId = null;
		render(AuthUserSync);
		expect(mockQueryClientClear).not.toHaveBeenCalled();
		expect(mockSetSelectedUserId).not.toHaveBeenCalled();
		expect(mockClearUserSelection).not.toHaveBeenCalled();
	});

	it('on sign-in, clears query cache and syncs selected user from the me query', () => {
		mockClerkContext.isLoaded = true;
		mockClerkContext.auth.userId = 'clerk_user_123';
		mockQueryState.data = { me: { id: '5', username: 'alice' } };
		render(AuthUserSync);
		expect(mockQueryClientClear).toHaveBeenCalledTimes(1);
		expect(mockSetSelectedUserId).toHaveBeenCalledWith(5);
		expect(mockClearUserSelection).not.toHaveBeenCalled();
	});

	it('on sign-out, clears query cache and clears user selection', () => {
		mockClerkContext.isLoaded = true;
		mockClerkContext.auth.userId = null;
		render(AuthUserSync);
		expect(mockQueryClientClear).toHaveBeenCalledTimes(1);
		expect(mockClearUserSelection).toHaveBeenCalledTimes(1);
		expect(mockSetSelectedUserId).not.toHaveBeenCalled();
	});

	it('account switch: mounting as a different signed-in user clears cache and re-syncs', () => {
		// First mount simulates being signed in as user A. A fresh mount is used
		// here (rather than mutating context mid-test) because the mocked
		// useClerkContext() returns a plain object, not a Svelte $state-backed
		// one like the real ClerkProvider — see AuthUserSync.svelte for why a
		// real account switch re-fires the effect via genuine reactivity.
		mockClerkContext.isLoaded = true;
		mockClerkContext.auth.userId = 'clerk_user_A';
		mockQueryState.data = { me: { id: '1', username: 'alice' } };
		const first = render(AuthUserSync);
		expect(mockSetSelectedUserId).toHaveBeenCalledWith(1);
		first.unmount();

		vi.clearAllMocks();

		mockClerkContext.auth.userId = 'clerk_user_B';
		mockQueryState.data = { me: { id: '2', username: 'bob' } };
		render(AuthUserSync);
		expect(mockQueryClientClear).toHaveBeenCalledTimes(1);
		expect(mockSetSelectedUserId).toHaveBeenCalledWith(2);
	});

	it('constructs the me query with a clerk-scoped key, 5 minute staleTime, and enabled when signed in', () => {
		mockClerkContext.isLoaded = true;
		mockClerkContext.auth.userId = 'clerk_user_123';
		render(AuthUserSync);
		expect(capturedQueryOptions.queryKey).toEqual(['me', 'clerk_user_123']);
		expect(capturedQueryOptions.staleTime).toBe(5 * 60 * 1000);
		expect(capturedQueryOptions.enabled).toBe(true);
	});

	it('disables the me query while signed out', () => {
		mockClerkContext.isLoaded = true;
		mockClerkContext.auth.userId = null;
		render(AuthUserSync);
		expect(capturedQueryOptions.enabled).toBe(false);
	});
});
```

- [ ] **Step 2: Run test to verify it fails**

```bash
cd frontend && pnpm exec vitest run tests/components/AuthUserSync.test.ts
```

Expected: FAIL — `frontend/src/lib/components/AuthUserSync.svelte` does not exist yet.

- [ ] **Step 3: Write the component**

Create `frontend/src/lib/components/AuthUserSync.svelte`:

```svelte
<script lang="ts">
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import { useClerkContext } from 'svelte-clerk';
	import { graphqlRequest } from '$lib/queries/client';
	import { ME, type MeResponse } from '$lib/queries/users';
	import { setSelectedUserId, clearUserSelection } from '$lib/stores/userSelection.svelte';

	const queryClient = useQueryClient();
	const clerk = useClerkContext();

	const clerkUserId = $derived(clerk.auth.userId);

	const meQuery = createQuery(() => ({
		queryKey: ['me', clerkUserId],
		queryFn: () => graphqlRequest<MeResponse>(ME),
		enabled: clerk.isLoaded && !!clerkUserId,
		staleTime: 5 * 60 * 1000,
	}));

	// Tracks the last clerkUserId we reacted to, so we only clear/resync on an
	// actual transition (sign-in, sign-out, or account switch on a shared
	// device) — not on every unrelated re-render.
	let lastSyncedClerkUserId: string | null | undefined = undefined;

	$effect(() => {
		if (!clerk.isLoaded) return;
		const currentId = clerkUserId ?? null;
		if (currentId === lastSyncedClerkUserId) return;

		lastSyncedClerkUserId = currentId;
		// content.lists(), perspectives.listByUser(), etc. are all user-scoped —
		// clear everything, not just the me query, so nothing leaks across
		// accounts on a shared device.
		queryClient.clear();

		if (currentId === null) {
			clearUserSelection();
		}
	});

	$effect(() => {
		if (clerk.isLoaded && clerkUserId && meQuery.data?.me) {
			setSelectedUserId(parseInt(meQuery.data.me.id, 10));
		}
	});
</script>
```

- [ ] **Step 4: Run test to verify it passes**

```bash
cd frontend && pnpm exec vitest run tests/components/AuthUserSync.test.ts
```

Expected: PASS (all 6 tests).

- [ ] **Step 5: Mount the component in the layout**

In `frontend/src/routes/+layout.svelte`, add the import:

```ts
import AuthUserSync from '$lib/components/AuthUserSync.svelte';
```

And mount it inside `<ClerkLoaded>`, alongside `<Header />`:

```svelte
<ClerkLoaded>
	<AuthUserSync />
	<div class="min-h-screen bg-background text-foreground">
		<Header />
		{@render children()}
	</div>
</ClerkLoaded>
```

- [ ] **Step 6: Type-check and run the full frontend test suite**

```bash
cd frontend && pnpm run check
cd frontend && pnpm run test:run
```

Expected: no new type errors, all tests pass.

- [ ] **Step 7: Commit**

```bash
git add frontend/src/lib/components/AuthUserSync.svelte frontend/src/routes/+layout.svelte frontend/tests/components/AuthUserSync.test.ts
git commit -m "feat(auth): sync selectedUserId to Clerk session via AuthUserSync"
```

---

## Task 5: Frontend — `useAddVideo.ts` zero-sentinel pass-through

**Files:**
- Modify: `frontend/src/lib/queries/hooks/useAddVideo.ts`
- Modify: `frontend/tests/unit/hooks-useAddVideo.test.ts`
- Modify: `frontend/tests/components/AddVideoPopover.test.ts`
- Modify: `frontend/tests/components/AddVideoDialog.test.ts`

**Interfaces:**
- Consumes: backend's zero-sentinel derivation in `createContentFromYouTube` (Task 2) — the mutation no longer needs a real `userId` from the client.
- Produces: `useAddVideo()`'s `mutationFn` always sends `userId: 0`; the "no user selected" error path is removed.

- [ ] **Step 1: Update the failing/obsolete tests first (red)**

In `frontend/tests/unit/hooks-useAddVideo.test.ts`:
- Remove `mockGetSelectedUserId` from the `vi.hoisted(...)` block and the `vi.mock('$lib/stores/userSelection.svelte', ...)` block entirely (lines 10, 18, 49-51).
- Remove the `it('shows "no user selected" message', ...)` test (lines 152-155) — there is no longer a code path that produces this error from `useAddVideo` itself.

In `frontend/tests/components/AddVideoPopover.test.ts`:
- Change the assertion in `mutationFn calls graphqlRequest with correct args` from `input: { url: '...', userId: 1 }` to `input: { url: '...', userId: 0 }`.
- Remove the `it('mutationFn throws when no user is selected', ...)` test (~lines 107-112).

In `frontend/tests/components/AddVideoDialog.test.ts`:
- Same two changes as `AddVideoPopover.test.ts` (matching assertion + test at the equivalent lines, ~87-97).

- [ ] **Step 2: Run tests to verify they fail**

```bash
cd frontend && pnpm exec vitest run tests/unit/hooks-useAddVideo.test.ts tests/components/AddVideoPopover.test.ts tests/components/AddVideoDialog.test.ts
```

Expected: FAIL on the `userId: 0` assertions (hook still sends `userId: 1` / throws on `null`).

- [ ] **Step 3: Update the hook**

Replace `frontend/src/lib/queries/hooks/useAddVideo.ts` in full:

```ts
import { createMutation, useQueryClient } from '@tanstack/svelte-query';
import { toast } from 'svelte-sonner';
import { graphqlRequest } from '../client';
import { CREATE_CONTENT_FROM_YOUTUBE, type CreateContentResponse, type ContentResponse } from '../content';
import { queryKeys } from '../keys';

export function useAddVideo() {
	const queryClient = useQueryClient();

	return createMutation(() => ({
		mutationFn: async (url: string) => {
			// userId: 0 is the "derive from my Clerk session" sentinel — the
			// backend resolves it via auth.RequireAuth when zero.
			return graphqlRequest<CreateContentResponse>(CREATE_CONTENT_FROM_YOUTUBE, {
				input: { url, userId: 0 },
			});
		},
		onSuccess: (data: CreateContentResponse) => {
			const result = data?.createContentFromYouTube;
			const newItem = result?.content;

			if (result?.alreadyExisted) {
				// VIDEO-05: Warn user that video already exists
				toast.warning('This video has already been added');
			} else {
				toast.success(`Added: ${newItem?.name ?? 'video'}`);
			}

			if (newItem) {
				// Only insert into cache if it's a genuinely new item
				if (!result?.alreadyExisted) {
					queryClient.setQueriesData<ContentResponse>({ queryKey: queryKeys.content.lists() }, (oldData) => {
						if (!oldData) return oldData;
						return {
							content: {
								...oldData.content,
								items: [newItem, ...oldData.content.items],
								totalCount: (oldData.content.totalCount ?? 0) + 1,
							},
						};
					});
				}

				// Mark stale for eventual consistency
				queryClient.invalidateQueries({
					queryKey: queryKeys.content.lists(),
					refetchType: 'none',
				});
			} else {
				// Fallback: full refetch if response shape is unexpected
				queryClient.invalidateQueries({ queryKey: queryKeys.content.lists() });
			}
		},
		onError: (err: Error) => {
			console.error('[AddVideo] mutation failed:', err);
			const message = err.message.toLowerCase();
			if (message.includes('invalid youtube url') || message.includes('video not found')) {
				toast.error('Invalid YouTube URL or video not found');
			} else if (message.includes('load failed') || message.includes('failed to fetch')) {
				toast.error('Cannot reach the server. Check your connection and try again.');
			} else if (message.includes('access denied') || message.includes('authentication required')) {
				toast.error('Please sign in to add a video');
			} else {
				toast.error('Failed to add video. Please try again.');
			}
		},
	}));
}
```

(Adds an `access denied`/`authentication required` branch to `onError` — the backend's new `@auth` derivation failure path is a realistic error this hook can now receive if a session is invalid; not tested explicitly by this task's test changes, but keep it since the message would otherwise fall through to the generic "Failed to add video" toast, which is misleading for a sign-out-mid-session case.)

- [ ] **Step 4: Run tests to verify they pass**

```bash
cd frontend && pnpm exec vitest run tests/unit/hooks-useAddVideo.test.ts tests/components/AddVideoPopover.test.ts tests/components/AddVideoDialog.test.ts
```

Expected: PASS.

- [ ] **Step 5: Run the full frontend test suite and type-check**

```bash
cd frontend && pnpm run check
cd frontend && pnpm run test:run
```

Expected: no new type errors, all tests pass.

- [ ] **Step 6: Commit**

```bash
git add frontend/src/lib/queries/hooks/useAddVideo.ts frontend/tests/unit/hooks-useAddVideo.test.ts frontend/tests/components/AddVideoPopover.test.ts frontend/tests/components/AddVideoDialog.test.ts
git commit -m "fix(video): always derive userId from Clerk session in useAddVideo"
```

---

## Self-Review Notes

- **Spec coverage:** §1 `me` query → Task 1. §2 `CreateContentFromYouTube` gap → Task 2. §3 `AuthUserSync.svelte` → Task 4. §4 "unmount, don't delete" → already true on this branch (documented in Global Constraints, no task needed). §5 Testing → covered per-task (backend `me` resolver + `CreateContentFromYouTube` zero-sentinel tests in Tasks 1-2; `AuthUserSync.test.ts` in Task 4; `useAddVideo`'s throw removed + replaced with 0-sentinel pass-through, and dependent component tests updated, in Task 5). Deferred Scope and Existing Data Note are explicitly out of scope — no tasks. Implementation Notes' branch question — resolved (stay on `feature/discovery-page`, in Global Constraints).
- **Placeholder scan:** no TBD/TODO, no "add appropriate error handling," no unshown code — clean.
- **Type consistency:** `Me`/`MeResponse` (Task 3) match the shape `AuthUserSync.svelte` (Task 4) destructures (`meQuery.data?.me.id`). `setSelectedUserId`/`clearUserSelection` signatures match the existing store (`frontend/src/lib/stores/userSelection.svelte.ts`) — no changes needed there. `CreateContentFromYouTubeInput.userId: IntID!` unchanged; frontend continues sending a number (now always `0`).
