# Coding Conventions

**Analysis Date:** 2026-02-07

## Naming Patterns

### Go Conventions

**Files:**
- `snake_case.go` — All Go files use lowercase with underscores
- Test files: `*_test.go` (e.g., `content_service_test.go`) placed in `test/` directory
- Package organization mirrors domain/adapter structure:
  - `internal/core/domain/` — Domain models
  - `internal/core/ports/` — Interfaces (repositories, services)
  - `internal/core/services/` — Business logic
  - `internal/adapters/` — Infrastructure (GraphQL, repositories)
  - `test/` — Test files mirroring source structure

**Functions and Methods:**
- PascalCase for exported functions: `GetByID`, `NewContentService`, `ListContent`, `Create`
- camelCase for unexported: `contentTypeToDBValue`, `toNullString`, `convertToRow`
- Constructor pattern: `New<Type>(...) *<Type>` (e.g., `NewContentService`, `NewContentRepository`)
- Receiver names: single lowercase letter (e.g., `r *ContentRepository`, `s *ContentService`)

**Variables and Constants:**
- camelCase for local variables: `userID`, `existing`, `createdAt`
- UPPERCASE for domain error constants: `ErrNotFound`, `ErrAlreadyExists`, `ErrInvalidInput`, `ErrInvalidURL`, `ErrYouTubeAPI`, `ErrInvalidRating`, `ErrDuplicateClaim`
- UPPERCASE for enum type constants: `ContentTypeYouTube`, `SortOrderASC`, `SortOrderDESC`, `PrivacyPublic`, `PrivacyPrivate`
- Type aliases for enums use PascalCase: `ContentType`, `SortOrder`, `Privacy`, `ReviewStatus`, `ContentSortBy`, `PerspectiveSortBy`
- Pointers for optional fields: `*string`, `*int`, `*time.Time`

**Struct Types:**
- PascalCase: `ContentRepository`, `ContentService`, `PaginatedContent`, `ContentListParams`
- Database row structs: lowercase: `contentRow`, `userRow`, `perspectiveRow`
- Mock implementations: `mock<Type>` (e.g., `mockContentRepository`, `mockYouTubeClient`)

**Domain Enums Pattern:**
- Type alias in PascalCase: `type ContentType string`
- Constants in UPPERCASE: `ContentTypeYouTube = "YOUTUBE"`, `ContentTypeArticle = "ARTICLE"`
- Bound to GraphQL in `gqlgen.yml` to prevent manual conversion
- Stored enums have DB converters (lowercase ↔ UPPERCASE)

### TypeScript/Svelte Conventions

**Files:**
- camelCase: `utils.ts`, `userSelection.svelte.ts`, `client.ts`, `content.ts`
- Components: PascalCase: `UserSelector.svelte`, `ActivityTable.svelte`, `PageWrapper.svelte`, `Header.svelte`
- Test files: `*.test.ts` (co-located with source or in `tests/` directory)

**Functions and Variables:**
- camelCase for all functions and variables: `formatDuration`, `handleChange`, `getSelectedUserId`, `setSelectedUserId`
- React/Svelte components: PascalCase
- Constants: UPPERCASE_SNAKE_CASE: `STORAGE_KEY = 'perspectize:selectedUserId'`, `GRAPHQL_ENDPOINT`

**Svelte 5 Reactive Patterns (MANDATORY):**
- State: `let variable = $state(initialValue)` — NOT `let variable = initialValue`
- Derived: `let derived = $derived(expression)` — NOT `$: derived = expression`
- Props via destructuring: `let { prop1, prop2 = default } = $props()`
- Render children via snippet: `let { children } = $props()` then `{@render children()}`
- Event handlers: `onclick={handler}` — NOT `on:click={handler}`

**Import Paths:**
- Absolute aliases: `$lib` → `src/lib`, `$app` → SvelteKit internals
- Relative imports for local references

## Code Style

### Go Formatting

**Tool:** `go fmt` (built-in)
- Enforced via `make fmt`
- All files must be gofmt-compliant

**Conventions:**
- Max line length: No strict limit, but keep under 120 characters when reasonable
- Imports: Organized in three groups separated by blank lines:
  1. Standard library: `"context"`, `"encoding/json"`, `"errors"`, `"fmt"`
  2. Third-party: `"github.com/stretchr/testify/assert"`
  3. Local project: `"github.com/CodeWarrior-debug/perspectize-be/perspectize-go/internal/..."`
- Example from `internal/core/services/content_service.go`:
  ```go
  import (
      "context"
      "errors"
      "fmt"

      "github.com/CodeWarrior-debug/perspectize-be/perspectize-go/internal/core/domain"
      "github.com/CodeWarrior-debug/perspectize-be/perspectize-go/internal/core/ports/repositories"
      portservices "github.com/CodeWarrior-debug/perspectize-be/perspectize-go/internal/core/ports/services"
  )
  ```

**Receivers:**
- Single lowercase letter abbreviations (e.g., `r *ContentRepository`, `s *ContentService`)
- Pointer receivers for methods that modify state or operate on large structs

### TypeScript/Svelte Formatting

**Tool:** No formatter configured (rely on IDE)
- Tab width: 2 spaces (SvelteKit default)
- Strict TypeScript: enabled in `tsconfig.json`

**Conventions:**
- Max line length: 120 characters (soft limit)
- Indentation: 2 spaces throughout
- CSS: Tailwind utilities via class attributes

## Import Organization

### Go

**Order (enforced by go fmt):**
1. Standard library: `"context"`, `"fmt"`, `"errors"`
2. Third-party: `"github.com/stretchr/testify/assert"`
3. Local: `"github.com/CodeWarrior-debug/perspectize-be/perspectize-go/internal/..."`

**Aliases:** Use for clarity when importing multiple ports packages
```go
portservices "github.com/CodeWarrior-debug/perspectize-be/perspectize-go/internal/core/ports/services"
```

### TypeScript/Svelte

**Order:**
1. Framework imports: `import { Component } from 'svelte'`
2. Third-party: `import { createQuery } from '@tanstack/svelte-query'`
3. Local imports: `import { graphqlClient } from '$lib/queries/client'`
4. Type imports: `import type { GridOptions } from '@ag-grid-community/core'`

**Aliases (via SvelteKit):**
- `$lib` → `src/lib`
- `$app/environment`, `$app/stores`, `$app/navigation`

**Example from `UserSelector.svelte`:**
```typescript
import { createQuery } from '@tanstack/svelte-query';
import { graphqlClient } from '$lib/queries/client';
import { LIST_USERS } from '$lib/queries/users';
import { setSelectedUserId, getSelectedUserId } from '$lib/stores/userSelection.svelte';
```

## Error Handling

### Go Patterns

**Domain Errors:** Defined in `internal/core/domain/errors.go`
```go
var (
    ErrNotFound       = errors.New("resource not found")
    ErrAlreadyExists  = errors.New("resource already exists")
    ErrInvalidInput   = errors.New("invalid input")
    ErrInvalidURL     = errors.New("invalid URL")
    ErrYouTubeAPI     = errors.New("youtube API error")
    ErrInvalidRating  = errors.New("rating must be between 0 and 10000")
    ErrDuplicateClaim = errors.New("claim already exists for this user")
)
```

**Wrapping Pattern:**
```go
// Repository detects domain error and wraps with context
if err != nil && !errors.Is(err, domain.ErrNotFound) {
    return nil, fmt.Errorf("failed to check existing content: %w", err)
}

// Service validates input and wraps domain error
if errors.Is(err, domain.ErrNotFound) {
    return nil, fmt.Errorf("resource not found")
}

// Never use switch on error types - use errors.Is()
// Example of what NOT to do:
// if err == domain.ErrNotFound { } // WRONG

// Correct approach:
// if errors.Is(err, domain.ErrNotFound) { } // CORRECT
```

**Detection:** Always use `errors.Is(err, domain.ErrNotFound)` to detect wrapped errors

**GraphQL Resolution:** Resolvers translate domain errors to HTTP responses
```go
result, err := r.service.GetById(id)
if errors.Is(err, domain.ErrNotFound) {
    return nil, fmt.Errorf("resource not found")
}
return result, err
```

### TypeScript Conventions

**Async/await with error handling:**
```typescript
try {
    const result = await graphqlClient.request(QUERY);
    return result;
} catch (error) {
    console.error('Query failed:', error);
    throw error;
}
```

**TanStack Query error states:**
```typescript
if (query.error) {
    console.error('Failed to load:', query.error);
}
```

## Logging

### Go

**Framework:** Standard library `log/slog`
- Usage: `log.Println("message")`, `log.Printf("format: %v", val)`, `log.Fatalf("error: %v", err)`
- Structured logging available: `slog.Info()`, `slog.Error()`
- Used in `cmd/server/main.go` for lifecycle events
- Example: `log.Println("Successfully connected to database!")`

### TypeScript/Svelte

**Framework:** Browser `console`
- Development: `console.log()`, `console.error()`, `console.warn()`, `console.debug()`
- No centralized logging library configured
- Example: `console.error('Failed to load:', error)`

## Comments

### Go

**When to Comment:**
- Exported functions and types: Required comment starting with name
  ```go
  // ContentService implements business logic for content operations
  type ContentService struct { ... }

  // GetByID retrieves content by ID
  func (s *ContentService) GetByID(ctx context.Context, id int) (*domain.Content, error) { ... }
  ```
- Complex logic: Explain WHY, not WHAT
- Non-obvious error handling: Comment domain error checks
- Example: `// Check if content already exists (ErrNotFound is expected if not found)`

**Style:** `// Comment` format (not `/* */` for single-line comments)

### TypeScript/Svelte

**When to Comment:**
- Complex business logic requiring explanation
- Non-obvious GraphQL or state management
- Component behavior that needs context

**Example from `ActivityTable.svelte`:**
```typescript
// Convert length + lengthUnits to display format
function formatDuration(length: number | null, lengthUnits: string | null): string { ... }

// Format dates to locale string
function formatDate(isoString: string): string { ... }
```

## Function Design

### Go

**Size:** 20-60 lines typical (service and repository methods)
- Prefer small focused functions
- Example methods: `GetByID`, `Create`, `List` are 10-30 lines each

**Parameters:**
- Context as first parameter (Go convention): `func (s *Service) Method(ctx context.Context, ...)`
- Domain models by pointer for mutations: `*domain.Content`
- Scalars by value: `id int`, `url string`
- Input structs for multiple parameters: `params domain.ContentListParams`

**Return Values:**
- Pointer + error tuple: `(*Model, error)`
- Pagination: `(*PaginatedModel, error)`
- Validation result with wrapped error: `fmt.Errorf("%w: message", domain.ErrInvalidInput)`

**Example from `internal/core/services/content_service.go`:**
```go
func (s *ContentService) GetByID(ctx context.Context, id int) (*domain.Content, error) {
    if id <= 0 {
        return nil, fmt.Errorf("%w: content id must be a positive integer", domain.ErrInvalidInput)
    }
    content, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("failed to get content: %w", err)
    }
    return content, nil
}
```

### TypeScript/Svelte

**Size:** Utilities 10-50 lines, components 50-200 lines

**Parameters:**
- Object when more than 2 parameters
- Svelte store getters: `export function getSelectedUserId(): number | null`
- Svelte store setters: `export function setSelectedUserId(value: number | null): void`

**Return Values:**
- Typed via TypeScript: `function getValue(): string`
- Nullable: `| null | undefined`

**Svelte 5 Store Patterns:**
- Getter: `export function getValue(): Type { return _internal; }`
- Setter: `export function setValue(val: Type): void { _internal = val; syncStorage(); }`
- Property: `export const value = { get value() { ... }, set value(...) { ... } }`

## Module Design

### Go Packages

**Exports:** Follow Effective Go
- Unexported (lowercase) by default
- Export only for cross-package use
- Example: Repository types exported as `*ContentRepository`, methods as `GetByID()`, `Create()`

**Package Structure:**
- `internal/core/domain/` — Domain models, NO external dependencies
- `internal/core/ports/` — Interfaces only (repositories, services)
- `internal/core/services/` — Business logic using ports
- `internal/adapters/` — Infrastructure (implements ports, depends on domain)
- `internal/adapters/repositories/postgres/` — Database implementations
- `internal/adapters/graphql/resolvers/` — GraphQL resolver implementations

**Dependency Rule:** Always points inward
```
adapters → ports → domain
services → ports → domain
Domain never imports adapters or services
```

### TypeScript/Svelte

**Exports:** ES modules
- Default for single exports: `export default Component`
- Named for multiple: `export function cn(...)`
- Re-export via barrel files: `export { Button } from './button/button.svelte'`

**Module Structure:**
- `src/lib/queries/` — GraphQL client and query definitions
- `src/lib/stores/` — Reactive state (Svelte 5 runes)
- `src/lib/components/` — Reusable components
- `src/lib/utils/` — Helper functions
- `src/routes/` — SvelteKit file-based routing

**Barrel Files:**
- Location: `src/lib/components/shadcn/index.ts`
- Pattern: `export { default as Button } from './button/button.svelte'`
- Enable: `import { Button } from '$lib/components/shadcn'`

## Svelte 5 Specific Conventions

**Runes (MANDATORY - violating these breaks reactivity):**

| Use This | NOT This | Reason |
|----------|----------|--------|
| `let count = $state(0)` | `let count = 0` | State requires rune for reactivity |
| `let doubled = $derived(count * 2)` | `$: doubled = count * 2` | Derived values use runes in v5 |
| `let { prop } = $props()` | `export let prop` | Props use runes in v5 |
| `$effect(() => { ... })` | `onMount(() => { ... })` | Effects replace lifecycle hooks |
| `{@render children()}` | `<slot />` | Render snippets in v5 |
| `onclick={handler}` | `on:click={handler}` | Event directive syntax changed |

**Anti-patterns (will break):**
- Using `$effect` for derivation → Use `$derived` instead
- Using `onMount`, `beforeUpdate`, etc. → Use `$effect` instead
- Using `$:` reactive statements → Use `$state` + `$derived`
- Using store `$store` syntax → TanStack Query v5 returns reactive objects (no `$`)

**Query Pattern (CRITICAL for v5):**
```typescript
// Function wrapper REQUIRED for reactivity
const query = createQuery(() => ({
    queryKey: ['key'],
    queryFn: () => graphqlClient.request(QUERY),
    staleTime: 5 * 60 * 1000  // optional
}));

// Access as reactive object (NO $ prefix)
{#if query.isLoading}Loading...{/if}
{#if query.data}{@render renderData(query.data)}{/if}
{#if query.error}Error: {query.error}{/if}
```

---

*Convention analysis: 2026-02-07*
