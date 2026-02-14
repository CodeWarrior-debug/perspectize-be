# Coding Conventions

**Analysis Date:** 2026-02-13

## Naming Patterns

**Go Files:**
- Package-level descriptors: `*_repository.go`, `*_service.go`, `*_resolver.go`, `*_test.go`
- Example: `user_repository.go`, `content_service.go`, `gorm_models.go`, `gorm_content_repository.go`
- Domain models: single nouns — `user.go`, `content.go`, `perspective.go`, `errors.go`

**Go Functions/Methods:**
- PascalCase for exported functions: `NewUserService()`, `Create()`, `GetByID()`, `ListAll()`
- Descriptive verb-noun pattern: `CreateFromYouTube()`, `GetByUsername()`, `List()`
- Receiver variable single lowercase letter: `func (s *UserService)`, `func (r *GormContentRepository)`
- Repository initializers prefix type: `NewGormUserRepository()`, `NewGormContentRepository()`

**Go Types/Interfaces:**
- PascalCase: `User`, `UserService`, `UserRepository`, `GormContentRepository`
- Interface names do not end in "Interface" — use bare name: `ContentRepository`, `UserService`
- Compile-time interface checks with underscore: `var _ repositories.ContentRepository = (*GormContentRepository)(nil)`
- Model types for GORM: `UserModel`, `ContentModel`, `PerspectiveModel` (append Model suffix)

**Go Variables:**
- camelCase for package vars: `emailRegex`, `STORAGE_KEY`
- Struct field names PascalCase: `Username`, `CreatedAt`, `MaxOpenConns`
- Short receiver names: `r`, `s`, `ctx` (context always abbreviated as `ctx`)
- Abbreviated but clear: `db`, `repo`, `err`
- Domain error constants UPPERCASE: `ErrNotFound`, `ErrAlreadyExists`, `ErrInvalidInput`

**TypeScript/Svelte Files:**
- Components: PascalCase `.svelte` files: `AddVideoDialog.svelte`, `Header.svelte`, `UserSelector.svelte`
- Query definitions: UPPERCASE_WITH_UNDERSCORE: `LIST_CONTENT`, `GET_CONTENT`, `CREATE_CONTENT_FROM_YOUTUBE`
- Utility functions: camelCase: `validateYouTubeUrl()`, `cn()`, `getSelectedUserId()`, `setSelectedUserId()`
- Store functions: camelCase with clear intent: `loadFromSession()`, `syncToSession()`, `clearUserSelection()`
- Type interfaces: PascalCase: `ContentItem`, `ContentResponse`, `CreateUserInput`
- Event handlers: camelCase with handle prefix or direct: `handleSubmit()`, `onclick={handler}`

**Environment Variables:**
- UPPERCASE_WITH_UNDERSCORES: `DATABASE_URL`, `YOUTUBE_API_KEY`, `APP_ENV`, `VITE_GRAPHQL_URL`
- Database pool config: `DB_MAX_OPEN_CONNS`, `DB_MAX_IDLE_CONNS`, `DB_CONN_MAX_LIFETIME`
- Storage keys: snake_case with namespace: `perspectize:selectedUserId`

## Code Style

**Go Formatting:**
- Tool: `go fmt` (built-in, enforced)
- Run: `make fmt`
- Indentation: tabs (Go standard)
- Imports organized in three groups:
  1. Standard library
  2. Third-party packages
  3. Local project imports
- Line length: no strict limit but prefer readable wrapping

**Go Linting:**
- Tool: `golangci-lint`
- Run: `make lint`
- Configuration: uses default ruleset (no config file checked in)
- Pre-commit hooks available via `make install-hooks`

**TypeScript/Svelte Formatting:**
- No explicit Prettier/ESLint config — project uses IDE defaults
- Indentation: 2 spaces (SvelteKit standard)
- Max line length: soft limit 120 characters

**Svelte 5 Runes (MANDATORY):**
- State: `let items = $state<Item[]>([])`
- Derived values: `let total = $derived(items.length)` (never use `$effect` for derivation)
- Props with defaults: `let { optional = 'default', required } = $props()`
- Bindable props: `let { open = $bindable(false) } = $props()`
- Side effects only: `$effect(() => { ... })` for DOM updates, subscriptions
- Event handlers: `onclick={handler}` not `on:click={handler}`
- Render children via snippet: `{@render children()}` not `<slot />`

**GraphQL Query Organization:**
- Define in separate resource files: `lib/queries/content.ts`, `lib/queries/users.ts`
- Use `gql` tagged template literals
- Export query/mutation definitions as named UPPERCASE constants
- Export TypeScript interfaces alongside queries: `ContentItem`, `ContentResponse`

## Import Organization

**Go:**
1. Standard library (`fmt`, `context`, `errors`)
2. External packages (`github.com/*`, `gorm.io/*`)
3. Local imports (`github.com/CodeWarrior-debug/perspectize/backend/*`)

Example from `content_service.go`:
```go
import (
	"context"
	"errors"
	"fmt"

	"github.com/CodeWarrior-debug/perspectize/backend/internal/core/domain"
	"github.com/CodeWarrior-debug/perspectize/backend/internal/core/ports/repositories"
	portservices "github.com/CodeWarrior-debug/perspectize/backend/internal/core/ports/services"
)
```

**TypeScript/Svelte:**
1. Framework imports: `svelte/*`, `$app/*`
2. External packages: `@tanstack/svelte-query`, `graphql-request`
3. Local imports: `$lib/*`
4. Type imports: `import type { Type } from '...'`

Example from `AddVideoDialog.svelte`:
```svelte
import { createMutation, useQueryClient } from '@tanstack/svelte-query';
import { toast } from 'svelte-sonner';
import { Dialog, DialogContent, /* ... */ } from '$lib/components/shadcn';
import { graphqlClient } from '$lib/queries/client';
import { CREATE_CONTENT_FROM_YOUTUBE } from '$lib/queries/content';
import { validateYouTubeUrl } from '$lib/utils/youtube';
```

**Path Aliases:**
- Go: Full explicit paths `github.com/CodeWarrior-debug/perspectize/backend/*`
- TypeScript: `$lib/` for `frontend/src/lib`, `$app/` for SvelteKit internals

## Error Handling

**Go Domain Errors:**
- Define as package variables in `core/domain/errors.go`
- Examples: `var ErrNotFound = errors.New("resource not found")`
- Use sentinel errors for specific cases: `ErrAlreadyExists`, `ErrInvalidInput`, `ErrInvalidURL`, `ErrYouTubeAPI`, `ErrInvalidRating`, `ErrDuplicateClaim`
- Wrap with context: `fmt.Errorf("failed to create user: %w", err)`
- Check via `errors.Is()`: `if errors.Is(err, domain.ErrNotFound) { ... }`

Example from `user_service.go`:
```go
if err == nil && existing != nil {
    return nil, fmt.Errorf("%w: username already taken", domain.ErrAlreadyExists)
}
if err != nil && !errors.Is(err, domain.ErrNotFound) {
    return nil, fmt.Errorf("failed to check username: %w", err)
}
```

**Go Repository Errors:**
- Convert GORM errors to domain errors: `if errors.Is(err, gorm.ErrRecordNotFound) { return nil, domain.ErrNotFound }`
- Wrap all database errors with context: `fmt.Errorf("failed to get content by id: %w", err)`

Example from `gorm_content_repository.go`:
```go
err := r.db.WithContext(ctx).First(&model, id).Error
if err != nil {
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, domain.ErrNotFound
    }
    return nil, fmt.Errorf("failed to get content by id: %w", err)
}
```

**Go Resolver Errors:**
- Translate domain errors to GraphQL error messages
- Log via `slog.Error()` for unexpected errors
- Return generic message to client for security

Example from `schema.resolvers.go`:
```go
if errors.Is(err, domain.ErrAlreadyExists) {
    return nil, fmt.Errorf("content already exists for this URL")
}
if errors.Is(err, domain.ErrInvalidURL) {
    return nil, fmt.Errorf("invalid YouTube URL")
}
slog.Error("creating content failed", "error", err)
return nil, fmt.Errorf("failed to create content")
```

**Svelte Error Handling:**
- Mutation error callbacks via TanStack Query: `onError: (err: Error) => { ... }`
- Parse error message for user-facing display
- Toast notifications for user feedback via `svelte-sonner`

Example from `AddVideoDialog.svelte`:
```typescript
onError: (err: Error) => {
    const message = err.message.toLowerCase();
    if (message.includes('already exists')) {
        toast.error('This video has already been added');
    } else if (message.includes('invalid youtube url')) {
        toast.error('Invalid YouTube URL or video not found');
    } else {
        toast.error('Failed to add video. Please try again.');
    }
}
```

## Logging

**Go Framework:** `log/slog` (structured logging)

**Patterns:**
- Info level for startup/lifecycle: `slog.Info("server running", "addr", addr)`
- Warn level for configuration issues: `slog.Warn(".env file not found", "hint", "set APP_ENV=production to suppress")`
- Error level for actual errors: `slog.Error("shutdown failed", "error", err)`
- Use key-value pairs for context: `"host", cfg.Database.Host, "port", cfg.Database.Port`
- Mask credentials in output (avoid logging DATABASE_URL values directly)
- Check examples in `cmd/server/main.go` for lifecycle logging pattern

**Svelte:** `console.error()` for development, `toast` notifications for user messages

## Comments

**When to Comment:**
- Explain non-obvious decisions, not obvious code
- Document why not what — code should be self-describing
- Mark temporary explanatory comments with `*TEMP*` prefix for easy grep/removal

Example from `cmd/server/main.go`:
```go
// M-06: request logging
r.Use(middleware.Logger)
// panic recovery
r.Use(middleware.Recoverer)
```

**Avoid:**
- Narrating what code does: "this function creates a user" (function name says this)
- Comments that duplicate code: `db.Close() // close database` (obvious from name)

**Function/Type Documentation:**
- Go requires doc comments for exported functions/types starting with the name
- Example: `// NewUserService creates a new user service`

## Function Design

**Size:**
- Aim for single responsibility
- Example: `validateYouTubeUrl()` — only validates, doesn't fetch or create
- Go typical: 20-50 lines for clarity
- Svelte components: 50-200 lines

**Parameters:**
- Go: Context first (for cancellation/timeout): `func (s *Service) GetByID(ctx context.Context, id int)`
- Go: Pointer receivers for methods that modify state
- Go: Input structs for multiple parameters: `func List(ctx context.Context, params domain.ContentListParams)`
- TypeScript: Destructure props: `let { optional = 'default', required } = $props()`

**Return Values:**
- Go: Multiple returns with error last: `(*domain.User, error)`
- Go: Use domain types for return values (not GORM models)
- TypeScript: Null/optional via union: `Promise<ContentResponse | null>`

**Interface Compliance:**
- Use compile-time check: `var _ repositories.ContentRepository = (*GormContentRepository)(nil)`
- All interface methods must be implemented or Go compiler will fail

## Module Design

**Go Exports:**
- Package-level factory functions with `New` prefix: `NewUserService()`, `NewGormContentRepository()`
- Explicit interface definitions in `ports/` — all implementations implement these interfaces
- Unexported types/functions by default, export only for cross-package use

**TypeScript Exports:**
- Named exports for functions/types: `export function cn(...)`, `export interface ContentItem`
- Barrel files for shadcn components: `lib/components/shadcn/index.ts` re-exports all primitives
- Query definitions as named exports: `export const LIST_CONTENT = gql`...``

**Barrel Files:**
- Location: `src/lib/components/shadcn/index.ts`
- Pattern: `export { default as Button } from './button'`
- Enable concise imports: `import { Button } from '$lib/components/shadcn'`

## Dependency Injection

**Go Pattern:**
- Constructor injection via factory functions: `NewUserService(repo UserRepository)`
- All dependencies are interfaces from `ports/`, not concrete types (enables testing/mocking)
- Wiring happens in `cmd/server/main.go`

Example from `main.go`:
```go
userRepo := postgres.NewGormUserRepository(db)
userService := services.NewUserService(userRepo)
resolver := resolvers.NewResolver(contentService, userService, perspectiveService)
```

**Svelte Pattern:**
- Props for component composition
- Store functions for shared state: `getSelectedUserId()`, `setSelectedUserId(value)`, `clearUserSelection()`
- TanStack Query client via `useQueryClient()` in mutation callbacks

## Database Naming Conventions

**Table Names:** snake_case, lowercase: `users`, `content`, `perspectives`

**Column Names:** snake_case, lowercase: `user_id`, `created_at`, `view_count`, `length_seconds`

**GORM Model Tags:**
- Use `gorm:` tags for database mapping: `gorm:"column:user_id"`, `gorm:"primaryKey"`
- Domain models in `core/domain/` have NO ORM tags — only GORM models in adapters do

Example from `gorm_models.go`:
```go
type UserModel struct {
    ID        int       `gorm:"primaryKey"`
    Username  string    `gorm:"uniqueIndex"`
    Email     string    `gorm:"uniqueIndex"`
    CreatedAt time.Time `gorm:"autoCreateTime"`
}
```

**Mappers Pattern:**
- Bidirectional conversion functions: `userDomainToModel()` and `userModelToDomain()`
- Located in `gorm_mappers.go` alongside domain models
- Keep domain models pure, use mappers for ORM-specific conversions

## Pagination Patterns

**Go Repository:**
- Cursor-based pagination with base64 encoding
- Helper functions: `encodeCursor(id)`, `decodeCursor(cursor)`
- Fetch `limit + 1` to detect `hasNextPage`
- Keyset pagination: `WHERE id > cursor` (not OFFSET)

Example from `gorm_content_repository.go`:
```go
if params.After != nil {
    cursorID, err := decodeCursor(*params.After)
    if dir == "DESC" {
        query = query.Where("id < ?", cursorID)
    } else {
        query = query.Where("id > ?", cursorID)
    }
}
// Fetch limit+1 for hasNextPage detection
query = query.Limit(limit + 1)
```

**Domain Pagination Types:**
- `PaginatedContent` contains `Items`, `HasNext`, `HasPrev`, `TotalCount`, `StartCursor`, `EndCursor`
- `ContentListParams` contains filter, sort, pagination parameters
- GORM handles pagination, domain models represent results

**GraphQL Types:** `pageInfo` object with `hasNextPage`, `hasPreviousPage`, `startCursor`, `endCursor`, and optional `totalCount`

## Svelte 5 Component Patterns

**Props Binding:**
- Use `$bindable()` for two-way binding: `let { open = $bindable(false) } = $props()`
- Read-only props without `$bindable`: `let { data } = $props()`
- Provide default values in destructuring

Example from `AddVideoDialog.svelte`:
```svelte
let { open = $bindable(false) } = $props();
let url = $state('');
let error = $state('');

$effect(() => {
    if (open) {
        url = '';
        error = '';
    }
});
```

**Form Handling:**
- Inline form submission: `onsubmit={handleSubmit}`
- Input binding: `bind:value={url}`
- Disabled state via query mutation status: `disabled={mutation.isPending}`

**TanStack Query Mutations:**
- Function wrapper pattern: `createMutation(() => ({ ... }))`
- Success callback for cache invalidation: `queryClient.invalidateQueries({ queryKey: ['content'] })`
- Error callback for error messages via toast notifications
- Reactive mutation state: `mutation.isPending`, `mutation.error`, `mutation.data`

**TanStack Query Queries:**
- Function wrapper pattern for reactivity: `createQuery(() => ({ queryKey: ['key'], queryFn: () => ... }))`
- Access as reactive object properties: `query.data`, `query.isLoading`, `query.error` (no `$` prefix)
- Never pass options directly; wrap in function for Svelte 5 reactivity

## GORM Patterns (Phase 7.1 ORM Migration)

**Model Structure:**
- GORM models in `adapters/repositories/postgres/gorm_models.go`
- Domain models in `core/domain/` (pure Go, no ORM tags)
- Mappers for bidirectional conversion in `gorm_mappers.go`

**Repository Pattern:**
- Implements interface from `core/ports/repositories/`
- Compile-time check: `var _ repositories.ContentRepository = (*GormContentRepository)(nil)`
- Use GORM method chaining for dynamic queries
- Context propagation: `r.db.WithContext(ctx)`

Example from `gorm_content_repository.go`:
```go
func (r *GormContentRepository) GetByID(ctx context.Context, id int) (*domain.Content, error) {
	var model ContentModel
	err := r.db.WithContext(ctx).First(&model, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get content by id: %w", err)
	}
	return contentModelToDomain(&model), nil
}
```

**Filter Application:**
- Build queries conditionally via GORM chaining
- Never pass user input directly to `Where()` — use parameterized queries
- Example: `query.Where("length >= ?", *params.Filter.MinLengthSeconds)`

---

*Convention analysis: 2026-02-13*
