# Coding Conventions

**Analysis Date:** 2026-02-16

## Naming Patterns

**Files (Backend Go):**
- Package-level: `{domain|service|repository}_{entity}.go` (e.g., `content_service.go`, `user_repository.go`)
- Repository implementations: `gorm_{entity}_repository.go` (e.g., `gorm_user_repository.go`)
- Mapper files: `gorm_mappers.go`, `gorm_models.go`
- Test files: `{source}_test.go` (e.g., `content_service_test.go`)
- Helper files: `helpers.go` for shared utilities

**Files (Frontend TypeScript/Svelte):**
- Components: PascalCase with `.svelte` extension (e.g., `ActivityTable.svelte`, `Header.svelte`)
- Utilities: camelCase with `.ts` extension (e.g., `formatting.ts`, `youtube.ts`)
- Queries: camelCase with `.ts` extension (e.g., `content.ts`, `users.ts`)
- Stores: camelCase with `.svelte.ts` extension (e.g., `userSelection.svelte.ts`)
- Test files: `{source}.test.ts`

**Functions (Go):**
- CamelCase, exported functions start with capital letter
- Constructor pattern: `New{TypeName}` (e.g., `NewContentService`, `NewGormUserRepository`)
- Receiver methods follow interface contract names
- Private functions: lowercase (e.g., `userDomainToModel`, `encodeCursor`)

**Functions (TypeScript/Svelte):**
- camelCase for all functions
- Export utility functions explicitly
- Component event handlers: `on{EventName}` (e.g., `onclick` in Svelte 5)

**Variables:**
- Go: camelCase, interface types descriptive (e.g., `repo repositories.ContentRepository`)
- TypeScript: camelCase for let/const, UPPERCASE for constants
- Svelte 5 state: `let value = $state(initial)` pattern

**Types:**
- Go: PascalCase (exported), interfaces in `ports/` packages
- Go enums: UPPERCASE constants (e.g., `ContentTypeYouTube = "YOUTUBE"`)
- TypeScript: PascalCase interfaces and types

**Packages/Modules:**
- Go: lowercase, single-word or compound (e.g., `services`, `repositories`, `domain`)
- TypeScript paths: `$lib/{feature}/{type}` (e.g., `$lib/queries/content`, `$lib/components/Header`)

## Code Style

**Formatting:**

Go:
- Enforced by `go fmt` (standard)
- Line length: no hard limit, readability-driven
- Imports organized: stdlib > third-party > internal

TypeScript/Svelte:
- No explicit formatter config (no `.prettierrc`)
- Code follows consistent spacing patterns
- svelte-check used during development (`npm run check`)

**Linting:**

Go:
- No explicit linting tool configured in Makefile
- Code follows standard Go conventions

TypeScript/Svelte:
- No ESLint/Prettier configured
- Tests via `vitest` (in `package.json`)
- svelte-check validates during `npm run check`

**Comments:**
- Explain "why" not "what"
- JSDoc/TSDoc for exported functions (e.g., `/** Format ISO date string to locale string. */`)
- Resolver comments document field purposes
- Single-line comments for inline logic

## Import Organization

**Order (Go):**
1. Standard library (`context`, `errors`, `fmt`, etc.)
2. Third-party packages (`gorm.io`, `github.com/...`)
3. Internal packages (full module path)

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

**Path Aliases (TypeScript):**
- `$lib` → `src/lib/` (SvelteKit built-in)
- No additional custom aliases
- Imports: `import { ActivityTable } from '$lib/components/ActivityTable.svelte'`

**Order (TypeScript/Svelte):**
1. Built-in/runtime imports (Svelte, SvelteKit)
2. Third-party packages (TanStack Query, graphql-request)
3. Local imports (components, queries, utilities)

## Error Handling

**Go Pattern:**

Domain errors in `internal/core/domain/errors.go`:
```go
var (
    ErrNotFound       = errors.New("resource not found")
    ErrAlreadyExists  = errors.New("resource already exists")
    ErrInvalidInput   = errors.New("invalid input")
    ErrInvalidURL     = errors.New("invalid URL")
)
```

Service layer wraps with context:
```go
func (s *ContentService) GetByID(ctx context.Context, id int) (*domain.Content, error) {
    if id <= 0 {
        return nil, fmt.Errorf("%w: content id must be positive", domain.ErrInvalidInput)
    }
    content, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("failed to get content: %w", err)
    }
    return content, nil
}
```

Repository translates database errors:
```go
func (r *GormUserRepository) GetByID(ctx context.Context, id int) (*domain.User, error) {
    var model UserModel
    err := r.db.WithContext(ctx).First(&model, id).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, domain.ErrNotFound
        }
        return nil, err
    }
    return userModelToDomain(&model), nil
}
```

Resolver translates to GraphQL errors:
```go
func (r *mutationResolver) CreateContentFromYouTube(ctx context.Context, input model.CreateContentFromYouTubeInput) (*model.Content, error) {
    content, err := r.ContentService.CreateFromYouTube(ctx, input.URL, input.UserID)
    if err != nil {
        if errors.Is(err, domain.ErrAlreadyExists) {
            return nil, fmt.Errorf("content already exists for this URL")
        }
        slog.Error("creating content failed", "error", err, "url", input.URL)
        return nil, fmt.Errorf("failed to create content")
    }
    return domainToModel(content), nil
}
```

**TypeScript Pattern:**

Errors via try-catch or TanStack Query state:
```typescript
if (query.isError) {
    // Handle error
}
```

## Logging

**Framework:** Go uses `log/slog` (structured logging)

**Pattern:**
```go
slog.Error("creating content failed",
    "error", err,
    "url", input.URL,
    "userID", input.UserID)

slog.Warn("failed to parse response JSON",
    "contentID", c.ID,
    "error", err)
```

**When to log:**
- Errors at resolver level (post domain/service)
- Warnings when data parsing fails but operation continues

TypeScript:
- No structured logging
- Frontend errors via TanStack Query error states or toast notifications

## Function Design

**Size (Go):**
- Service methods: 15-50 lines
- Repository methods: 10-30 lines
- Resolvers: 20-40 lines (includes error handling)

**Parameters:**
- Go: Receiver methods, context first param, domain entities as input
- Validation early with specific error wrapping
- Slices for lists, pointers for optional values

**Return Values:**
- Go: `(T, error)` pattern exclusively
- Domain model conversion in mappers, not method returns

**Svelte Components:**
- Use `let { prop } = $props()` for inputs
- Reactivity via `$state`, `$derived`, `$effect`
- Event handlers via `onchange={}`, `onclick={}` attributes

## Module Design

**Exports (Go):**
- Constructors: `New{Type}(...) *{Type}`
- Methods on exported types implementing service/repository interfaces
- Private helpers: lowercase

**Exports (TypeScript):**
- Named exports: `export function {name}() {}`
- Barrel files sparingly (shadcn components have `index.ts`)
- Components default exported

**Interfaces:**
- Defined in `internal/core/ports/`
- Compile-time check: `var _ Interface = (*ConcreteType)(nil)`
- Decouple adapters from domain

**Domain Models (Go):**
- Pure Go structs, no external dependencies
- In `internal/core/domain/`
- Converted via mappers (`gorm_mappers.go`)
- Ensures domain purity and testability

## Specific Patterns

**GraphQL Enum Handling (REQUIRED):**
Use gqlgen model binding; never switch statements:
```go
type ContentType string
const ContentTypeYouTube ContentType = "YOUTUBE"
```

Bind in `gqlgen.yml`:
```yaml
models:
  ContentType:
    model:
      - github.com/CodeWarrior-debug/perspectize/backend/internal/core/domain.ContentType
```

**ID Handling:**
- Top-level query/mutation IDs: `string` type, convert via `strconv.Atoi` in resolver
- Filter/input field IDs: use `IntID` scalar (`pkg/graphql/intid.go`)

**Repository Pagination (Go):**
- Cursor-based with opaque base64 (`cursor:{id}`)
- Fetch `limit+1` to determine `hasNextPage`
- Helpers in `helpers.go`: `encodeCursor`, `decodeCursor`
- Whitelist sort columns (prevents SQL injection)

**GORM Patterns:**
- Domain models: `core/domain/`
- GORM models: `adapters/repositories/postgres/gorm_models.go`
- Mappers: bidirectional conversion
- `gorm.ErrRecordNotFound` → domain errors at repository
- Context via `db.WithContext(ctx)` for all queries

**TanStack Query (Frontend):**
- Query functions return data or throw
- Keys in `lib/queries/keys.ts`
- Svelte 5: result is reactive object (access as `query.data`, not `$query.data`)
- Function wrapper: `createQuery(() => ({ ... }))`

**Svelte 5 Patterns (REQUIRED):**
- `let value = $state(initial)` for state
- `let derived = $derived(computed)` for reactive values
- `let { prop } = $props()` for props
- `$effect(() => { ... })` for side effects
- Event handlers: `on{EventName}` attribute (e.g., `onchange`, `onclick`)
- Slot rendering: `{@render children()}` with `let { children } = $props()`
- Never use Svelte 4: no `<slot/>`, no `$:` reactivity, no `export let`

---

*Convention analysis: 2026-02-16*
