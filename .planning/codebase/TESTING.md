# Testing Patterns

**Analysis Date:** 2026-02-07

## Test Framework

### Go

**Runner:**
- Go standard `testing` package (built-in)
- Config: `Makefile` with targets for running tests
- No external test runner or discovery tool

**Assertion Library:**
- `github.com/stretchr/testify/assert` — Non-fatal assertions (`assert.Equal()`, `assert.Nil()`, `assert.Contains()`, `assert.True()`, `assert.False()`)
- `github.com/stretchr/testify/require` — Fatal assertions (`require.NoError()`, `require.Equal()`) that stop test on failure immediately

**Run Commands:**
```bash
make test                     # Run all tests with coverage
make test-coverage            # Generate HTML coverage report (→ coverage.html)
go test -v ./...              # Run all tests with verbose output
go test -v ./internal/core/services          # Run tests in specific package
go test -v -run TestGetByID ./internal/...   # Run specific test by name
```

### TypeScript/Svelte

**Runner:**
- Vitest (configured in `vite.config.ts`)
- Environment: jsdom
- Setup files: `tests/setup.ts`

**Assertion Library:**
- `vitest` built-in assertions (`expect()`, `toBe()`, `toEqual()`, `toContain()`, `toThrow()`)
- `@testing-library/jest-dom` for DOM matchers (optional)

**Run Commands:**
```bash
pnpm run test              # Run tests in watch mode
pnpm run test:run          # Run all tests once
pnpm run test:coverage     # Generate coverage report
pnpm run test:duplication  # Code duplication check (jscpd)
```

## Test File Organization

### Go

**Location:**
- Centralized in `test/` directory at project root
- Mirrored structure from source:
  ```
  backend/
  ├── internal/core/services/content_service.go    (source)
  ├── test/services/content_service_test.go        (test)
  ├── internal/core/domain/content.go              (source)
  └── test/domain/content_test.go                  (test)
  ```

**Naming:**
- Test files: `{subject}_test.go` (e.g., `content_service_test.go`, `errors_test.go`)
- Package: `{subject}_test` (e.g., `package services_test`, `package domain_test`)
- Test functions: `Test{FunctionName}` (e.g., `TestGetByID_Success`, `TestGetByID_NotFound`)

**Directory Structure:**
```
backend/test/
├── services/
│   ├── content_service_test.go
│   ├── perspective_service_test.go
│   └── user_service_test.go
├── domain/
│   ├── content_test.go
│   ├── errors_test.go
│   ├── perspective_test.go
│   └── user_test.go
├── resolvers/
│   └── content_resolver_test.go
├── config/
│   └── config_test.go
├── database/
│   └── postgres_test.go
└── youtube/
    └── parser_test.go
```

### TypeScript/Svelte

**Location:**
- Tests in `tests/` directory mirroring `src/` structure
- Co-located with source for unit tests or in `tests/{category}/`
  ```
  frontend/
  ├── src/lib/utils.ts
  ├── tests/unit/utils.test.ts
  ├── src/lib/stores/userSelection.svelte.ts
  └── tests/unit/stores-userSelection.test.ts
  ```

**Naming:**
- Test files: `*.test.ts` or `*.spec.ts`
- Kebab-case for multi-word tests: `queries-users.test.ts`, `stores-userSelection.test.ts`
- Test suites: `describe('name', () => { ... })`
- Test cases: `it('should do something', () => { ... })`

**Directory Structure:**
```
frontend/tests/
├── setup.ts
├── unit/
│   ├── utils.test.ts
│   ├── queries-client.test.ts
│   ├── queries-content.test.ts
│   ├── queries-users.test.ts
│   ├── stores-userSelection.test.ts
│   └── shadcn-barrel.test.ts
├── components/
│   └── (component tests here)
├── fixtures/
│   └── (mock data)
└── helpers/
    └── (test utilities)
```

## Test Structure

### Go

**Suite Organization Pattern (from `test/services/content_service_test.go`):**
```go
package services_test

import (
    "context"
    "errors"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/CodeWarrior-debug/perspectize/backend/internal/core/domain"
    "github.com/CodeWarrior-debug/perspectize/backend/internal/core/services"
)

// Mock implementation of interface
type mockContentRepository struct {
    getByIDFn func(ctx context.Context, id int) (*domain.Content, error)
}

func (m *mockContentRepository) GetByID(ctx context.Context, id int) (*domain.Content, error) {
    if m.getByIDFn != nil {
        return m.getByIDFn(ctx, id)
    }
    return nil, domain.ErrNotFound
}

// Test cases grouped by method
func TestGetByID_Success(t *testing.T) {
    expected := &domain.Content{ID: 1, Name: "Test"}
    repo := &mockContentRepository{
        getByIDFn: func(ctx context.Context, id int) (*domain.Content, error) {
            assert.Equal(t, 1, id)  // Assert input
            return expected, nil
        },
    }
    svc := services.NewContentService(repo)

    result, err := svc.GetByID(context.Background(), 1)

    require.NoError(t, err)  // Fatal if error
    assert.Equal(t, expected, result)  // Non-fatal assertion
}

func TestGetByID_NotFound(t *testing.T) {
    repo := &mockContentRepository{
        getByIDFn: func(ctx context.Context, id int) (*domain.Content, error) {
            return nil, domain.ErrNotFound
        },
    }
    svc := services.NewContentService(repo)

    result, err := svc.GetByID(context.Background(), 999)

    assert.Nil(t, result)
    require.Error(t, err)
    assert.True(t, errors.Is(err, domain.ErrNotFound))
}
```

**Patterns:**
- Mock interfaces implement methods with `Fn` field callbacks
- Table-driven tests for multiple input combinations
- One assertion per test (or grouped logically)
- Use `require.NoError()` for setup errors that should stop test
- Use `assert.Equal()` for business logic assertions
- Context always passed: `context.Background()` for unit tests
- Subtest naming: `TestFunctionName_Scenario` or `TestFunctionName_ErrorCondition`

### TypeScript/Svelte

**Suite Organization Pattern (from `tests/unit/utils.test.ts`):**
```typescript
import { describe, it, expect } from 'vitest';
import { cn } from '$lib/utils';

describe('cn() utility', () => {
    it('returns empty string for no arguments', () => {
        expect(cn()).toBe('');
    });

    it('passes through a single class', () => {
        expect(cn('text-red-500')).toBe('text-red-500');
    });

    it('merges multiple classes', () => {
        const result = cn('px-4', 'py-2', 'text-sm');
        expect(result).toContain('px-4');
        expect(result).toContain('py-2');
    });

    it('handles conditional classes via clsx', () => {
        const isActive = true;
        const result = cn('base', isActive && 'active');
        expect(result).toContain('active');
    });
});
```

**Patterns:**
- Descriptive test names with `it('should...')`
- One logical assertion per test
- Setup shared state with `beforeEach()`
- Clear arrange-act-assert structure within test body
- `describe()` for test grouping by feature/function

## Mocking

### Go

**Pattern:** Hand-written mocks implementing interfaces
```go
type mockContentRepository struct {
    createFn   func(ctx context.Context, content *domain.Content) (*domain.Content, error)
    getByIDFn  func(ctx context.Context, id int) (*domain.Content, error)
    getByURLFn func(ctx context.Context, url string) (*domain.Content, error)
}

func (m *mockContentRepository) Create(ctx context.Context, content *domain.Content) (*domain.Content, error) {
    if m.createFn != nil {
        return m.createFn(ctx, content)
    }
    return content, nil  // Default behavior
}
```

**What to Mock:**
- Repository interfaces (use function fields for behavior)
- External service clients (YouTube API, etc.)
- Database connections (use sqlmock for integration tests)

**What NOT to Mock:**
- Domain models — construct real instances
- Error wrapping logic — test with real errors
- Validation rules — test with real inputs

### TypeScript/Svelte

**Pattern:** Vitest vi.mock() for module mocking in setup
```typescript
// tests/setup.ts
vi.mock('$app/environment', () => ({
    browser: true,
    dev: true,
    building: false
}));

vi.mock('$app/stores', () => ({
    page: readable({ url: new URL('http://localhost'), ... })
}));
```

**What to Mock:**
- SvelteKit internals (`$app/environment`, `$app/stores`, `$app/navigation`)
- Static assets (`$lib/assets/favicon.svg`)
- External services (GraphQL client)

**What NOT to Mock:**
- Utility functions — import and test directly
- Components — test with real instances
- Store behavior — test actual store code

### Module Reset Pattern (TypeScript)

```typescript
// tests/unit/stores-userSelection.test.ts
import { beforeEach, vi } from 'vitest';

describe('userSelection store', () => {
    beforeEach(() => {
        sessionStorage.clear();
        vi.resetModules();  // Re-import store with fresh state
    });

    it('loads stored user ID', async () => {
        sessionStorage.setItem('perspectize:selectedUserId', '42');
        const store = await import('$lib/stores/userSelection.svelte');
        expect(store.selectedUserId.value).toBe(42);
    });
});
```

## Fixtures and Factories

### Go

**Test Data Creation Pattern:**
- Inline construction in test functions
- Domain models created directly: `&domain.Content{ID: 1, Name: "Test"}`
- No factory functions (tests self-contained)

**Example from `test/services/content_service_test.go`:**
```go
func TestCreateFromYouTube_Success(t *testing.T) {
    metadata := &portservices.VideoMetadata{
        Title:       "Test Video Title",
        Description: "A great video",
        Duration:    300,
        ChannelName: "Test Channel",
        Response:    json.RawMessage(`{"items":[]}`),
    }
    // ... use metadata in test
}
```

### TypeScript/Svelte

**Fixtures Directory:** `tests/fixtures/`
- Reusable mock data
- Example: `tests/fixtures/users.ts` with mock user objects

**Example from `tests/fixtures/`:**
```typescript
// tests/fixtures/users.ts
export const mockUsers = [
    { id: '1', username: 'alice', email: 'alice@example.com' },
    { id: '2', username: 'bob', email: 'bob@example.com' }
];
```

**Usage in Tests:**
```typescript
import { mockUsers } from '../fixtures/users';

it('displays user list', () => {
    const query = createQuery(() => ({
        queryKey: ['users'],
        queryFn: () => Promise.resolve({ users: mockUsers })
    }));
    // ... test with mockUsers
});
```

## Coverage

### Go

**Requirements:** None enforced
- Targets: Check actual coverage via HTML report
- View coverage: `make test-coverage` → opens `coverage.html` in browser

**Typical Coverage:**
- Core domain/services: 80-90%
- Repositories: 85%+
- Resolvers: 70-80%

### TypeScript/Svelte

**Requirements:** Configured thresholds in `vite.config.ts`
```typescript
coverage: {
    provider: 'v8',
    thresholds: {
        lines: 80,
        functions: 80,
        branches: 75,
        statements: 80
    }
}
```

**View Coverage:**
```bash
pnpm run test:coverage
```

**Exclusions:**
- `node_modules/`, `.svelte-kit/`
- `**/*.d.ts`
- Config files, setup files
- `src/lib/components/shadcn/**` (third-party)
- `src/routes/**` (SvelteKit routes)

## Test Types

### Go

**Unit Tests:**
- Location: `test/services/`, `test/domain/`, `test/repositories/`
- Scope: Single function/method
- Mocks: All external dependencies
- Database: Not used (mocked repositories)
- Example: `TestContentService.GetByID()` tests only service logic

**Integration Tests:**
- Location: `test/resolvers/`, `test/database/`
- Scope: Multiple components working together
- Mocks: External services only (YouTube API)
- Database: Real PostgreSQL connection (auto-skip if unavailable)
- Pattern: Auto-skip when DB not available
  ```go
  func TestGetByID_WithDB(t *testing.T) {
      db := setupTestDB(t)  // Skip test if DB unavailable
      // ... test with real DB
  }
  ```

### TypeScript/Svelte

**Unit Tests:**
- Location: `tests/unit/`
- Scope: Single function, utility, or store logic
- Mocks: SvelteKit internals
- Example: `utils.test.ts` tests `cn()` utility in isolation

**Component Tests (not fully implemented):**
- Would use `@testing-library/svelte`
- Located in `tests/components/`
- Test component behavior, not implementation

**No E2E Tests Currently:**
- Would use Playwright or Cypress
- Not configured in this project

## Common Patterns

### Go Async Testing

```go
func TestGetContentAsync(t *testing.T) {
    // Create channel for async result
    result := make(chan *domain.Content)
    go func() {
        content, _ := service.GetByID(context.Background(), 1)
        result <- content
    }()

    // Wait for result with timeout
    select {
    case content := <-result:
        assert.NotNil(t, content)
    case <-time.After(time.Second):
        t.Fatal("operation timed out")
    }
}
```

**OR use context with timeout:**
```go
func TestGetByIDWithTimeout(t *testing.T) {
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()

    result, err := service.GetByID(ctx, 1)
    require.NoError(t, err)
    assert.NotNil(t, result)
}
```

### Go Error Testing

```go
func TestGetByID_InvalidID_Zero(t *testing.T) {
    svc := services.NewContentService(&mockContentRepository{})

    result, err := svc.GetByID(context.Background(), 0)

    assert.Nil(t, result)
    require.Error(t, err)
    assert.True(t, errors.Is(err, domain.ErrInvalidInput))  // Check error type
    assert.Contains(t, err.Error(), "content id must be a positive integer")  // Check message
}

func TestCreateFromYouTube_AlreadyExists(t *testing.T) {
    repo := &mockContentRepository{
        getByURLFn: func(ctx context.Context, url string) (*domain.Content, error) {
            return &domain.Content{ID: 1}, nil  // Already exists
        },
    }

    result, err := svc.CreateFromYouTube(context.Background(), url, idExtractor)

    assert.Nil(t, result)
    require.Error(t, err)
    assert.True(t, errors.Is(err, domain.ErrAlreadyExists))
}
```

### TypeScript Async Testing

```typescript
it('loads data from GraphQL', async () => {
    const { result } = renderHook(() =>
        createQuery(() => ({
            queryKey: ['data'],
            queryFn: () => graphqlClient.request(QUERY)
        }))
    );

    // Wait for query to resolve
    await waitFor(() => {
        expect(result.current.isLoading).toBe(false);
    });

    expect(result.current.data).toBeDefined();
});
```

### TypeScript Error Testing

```typescript
it('handles query errors', async () => {
    vi.spyOn(graphqlClient, 'request').mockRejectedValueOnce(
        new Error('Network error')
    );

    const { result } = renderHook(() =>
        createQuery(() => ({
            queryKey: ['data'],
            queryFn: () => graphqlClient.request(QUERY)
        }))
    );

    await waitFor(() => {
        expect(result.current.error).toBeDefined();
    });
});
```

### TypeScript Store Testing with Reset

```typescript
describe('userSelection store', () => {
    beforeEach(() => {
        sessionStorage.clear();
        vi.resetModules();  // Force re-import with clean state
    });

    it('persists to session storage', async () => {
        const store = await import('$lib/stores/userSelection.svelte');
        store.setSelectedUserId(42);

        expect(sessionStorage.getItem('perspectize:selectedUserId')).toBe('42');
    });

    it('loads from session storage on init', async () => {
        sessionStorage.setItem('perspectize:selectedUserId', '42');
        const store = await import('$lib/stores/userSelection.svelte');

        expect(store.selectedUserId.value).toBe(42);
    });
});
```

---

*Testing analysis: 2026-02-07*
