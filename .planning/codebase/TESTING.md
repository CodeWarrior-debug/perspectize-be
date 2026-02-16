# Testing Patterns

**Analysis Date:** 2026-02-16

## Backend: Go Testing

### Test Framework

**Runner:**
- Standard library `testing` package
- Execute: `make test` or `go test -v -cover -coverpkg=./internal/...,./pkg/... ./...`

**Assertion Library:**
- `github.com/stretchr/testify/assert` (assertions)
- `github.com/stretchr/testify/require` (assertions that fail test)

**Run Commands:**
```bash
make test              # Run all tests with coverage
make test-coverage     # Run tests, generate coverage report → coverage.html
go test ./...          # Run all tests
go test -v ./...       # Verbose output
go test -run TestName  # Run specific test
go test -short ./...   # Skip integration tests (use t.Skip())
```

### Test File Organization

**Location:**
- Separate `test/` directory at project root (not co-located with source)
- Mirror source structure: `test/{domain|services|repositories}/{name}_test.go`

**Directory Structure:**
```
backend/
├── internal/
│   ├── core/
│   │   ├── domain/
│   │   ├── services/
│   │   └── ports/
│   └── adapters/
│       └── repositories/postgres/
└── test/
    ├── config/
    ├── database/
    ├── domain/
    ├── services/
    ├── repositories/
    ├── resolvers/
    ├── youtube/
    └── graphql/
```

**Naming:**
- Test files: `{entity}_{test_scope}_test.go` (e.g., `content_service_test.go`)
- Benchmark files: `{entity}_bench_test.go` (e.g., `perspective_service_bench_test.go`)
- Package name: `{source_package}_test` (e.g., `services_test` for testing `services` package)

### Test Structure

**Suite Organization:**

Go uses table-driven tests. Each test is a function starting with `Test`:

```go
package services_test

import (
    "context"
    "testing"

    "github.com/CodeWarrior-debug/perspectize/backend/internal/core/domain"
    "github.com/CodeWarrior-debug/perspectize/backend/internal/core/services"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

// Test function: TestGetByID_Success
func TestGetByID_Success(t *testing.T) {
    url := "https://youtube.com/watch?v=abc123"
    expected := &domain.Content{
        ID:          1,
        Name:        "Test Video",
        URL:         &url,
        ContentType: domain.ContentTypeYouTube,
    }

    repo := &mockContentRepository{
        getByIDFn: func(ctx context.Context, id int) (*domain.Content, error) {
            assert.Equal(t, 1, id)
            return expected, nil
        },
    }

    svc := services.NewContentService(repo, &mockYouTubeClient{})
    result, err := svc.GetByID(context.Background(), 1)

    require.NoError(t, err)
    assert.Equal(t, expected, result)
}
```

**Patterns:**

1. **Test naming:** `Test{FunctionName}_{Scenario}` (e.g., `TestGetByID_Success`, `TestGetByID_InvalidID`)
2. **Setup:** Create mocks and service instances inline
3. **Act:** Call the function
4. **Assert:** Use `require` for critical assertions (fail test immediately), `assert` for checks
5. **Cleanup:** Implicit (no defer needed for simple mocks)

### Mocking

**Framework:** Manual mock structs (no external mocking library)

**Pattern:**

Mocks implement interfaces with function pointers:

```go
// mockContentRepository implements repositories.ContentRepository for testing
type mockContentRepository struct {
    createFn   func(ctx context.Context, content *domain.Content) (*domain.Content, error)
    getByIDFn  func(ctx context.Context, id int) (*domain.Content, error)
    getByURLFn func(ctx context.Context, url string) (*domain.Content, error)
    listFn     func(ctx context.Context, params domain.ContentListParams) (*domain.PaginatedContent, error)
}

func (m *mockContentRepository) Create(ctx context.Context, content *domain.Content) (*domain.Content, error) {
    if m.createFn != nil {
        return m.createFn(ctx, content)
    }
    return content, nil
}

func (m *mockContentRepository) GetByID(ctx context.Context, id int) (*domain.Content, error) {
    if m.getByIDFn != nil {
        return m.getByIDFn(ctx, id)
    }
    return nil, domain.ErrNotFound
}

// ... implement all interface methods
```

**Usage:**

```go
repo := &mockContentRepository{
    getByIDFn: func(ctx context.Context, id int) (*domain.Content, error) {
        assert.Equal(t, 1, id)
        return expected, nil
    },
}
svc := services.NewContentService(repo, &mockYouTubeClient{})
```

**What to Mock:**
- Repository interfaces (database access)
- Service interfaces (external APIs like YouTube)
- Any port interface

**What NOT to Mock:**
- Domain models (use real instances)
- Conversion functions (test with real data)
- Error creation (use real domain errors)

### Fixtures and Factories

**Test Data:**

No dedicated fixture files. Test data created inline:

```go
func TestGetByID_Success(t *testing.T) {
    url := "https://youtube.com/watch?v=abc123"
    expected := &domain.Content{
        ID:          1,
        Name:        "Test Video",
        URL:         &url,
        ContentType: domain.ContentTypeYouTube,
    }
    // ... test continues
}
```

**Helpers in Test File:**

Common assertions extracted to helpers:
```go
func clearConfigEnvVars(t *testing.T) {
    t.Setenv("DATABASE_URL", "")
    t.Setenv("YOUTUBE_API_KEY", "")
    t.Setenv("DATABASE_PASSWORD", "")
}
```

### Coverage

**Requirements:** Not enforced (no coverage gate in CI)

**View Coverage:**
```bash
make test-coverage  # Generates coverage.html
go test -cover ./...  # Command-line output
```

**Coverage tracked in:**
- `coverage.out` (raw)
- `coverage.html` (visual report)

## Frontend: Vitest Testing

### Test Framework

**Runner:**
- `vitest` (Vite-native test runner)
- Config: Defaults in `package.json`, no `vitest.config.ts`

**Assertion Library:**
- `vitest` assertions (expect API)
- `@testing-library/svelte` (component testing)
- `@testing-library/jest-dom` (DOM assertions)

**Run Commands:**
```bash
pnpm run test       # Watch mode
pnpm run test:run   # Run once (CI)
pnpm run test:coverage  # Coverage report
```

### Test File Organization

**Location:**
- Co-located with source in `tests/` directory
- Mirror structure: `tests/{unit|components}/{name}.test.ts`

**Directory Structure:**
```
frontend/
├── src/lib/
│   ├── components/
│   ├── queries/
│   ├── stores/
│   └── utils/
└── tests/
    ├── setup.ts           # Global test setup
    ├── helpers/
    │   └── TestWrapper.svelte
    ├── unit/              # Utility and query tests
    │   ├── formatting.test.ts
    │   ├── queries-content.test.ts
    │   ├── youtube.test.ts
    │   └── ...
    └── components/        # Component tests
        ├── ActivityTable.test.ts
        ├── Header.test.ts
        ├── UserSelector.test.ts
        └── ...
```

**Naming:**
- Test files: `{source}.test.ts`
- Setup file: `setup.ts` (auto-loaded by vitest)
- Example: `formatting.test.ts` tests `src/lib/utils/formatting.ts`

### Test Structure

**Suite Organization:**

```typescript
import { describe, it, expect, beforeEach, vi } from 'vitest';

describe('formatDuration', () => {
    it('returns dash for null length', () => {
        expect(formatDuration(null, null)).toBe('—');
    });

    it('formats seconds as minutes:seconds', () => {
        expect(formatDuration(300, 'seconds')).toBe('5:00');
    });

    it('formats large durations', () => {
        expect(formatDuration(3661, 'seconds')).toBe('61:01');
    });
});
```

**Patterns:**

1. **Suite:** `describe(description, () => { ... })`
2. **Test:** `it(description, () => { ... })`
3. **Setup:** `beforeEach(() => { ... })`
4. **Assertion:** `expect(value).toBe(expected)`
5. **Mocking:** `vi.mock(module)`, `vi.fn()`

### Mocking

**Framework:** `vitest` built-in (vi.mock, vi.fn, vi.spyOn)

**Pattern - Module Mocks:**

```typescript
// Mock entire module
vi.mock('$lib/queries/client', () => ({
    graphqlClient: {
        request: vi.fn()
    }
}));

// In hoisted scope
const { mockRequest } = vi.hoisted(() => ({
    mockRequest: vi.fn()
}));

// Use in test
beforeEach(() => {
    mockRequest.mockResolvedValue(mockDataResponse);
});
```

**Pattern - Component Mocks:**

```typescript
// Mock AG Grid component (complex library)
vi.mock('ag-grid-svelte5', () => ({
    default: vi.fn(() => ({
        $$: {},
        $set: vi.fn(),
        $on: vi.fn(),
        $destroy: vi.fn(),
    })),
}));
```

**What to Mock:**
- External API clients (GraphQL queries)
- Complex third-party components (AG Grid, query providers)
- SvelteKit built-ins (`$app/navigation`, `$app/stores`, `$app/environment`)

**What NOT to Mock:**
- Simple utility functions (format functions, extractors)
- Component logic (test real behavior)
- TanStack Query itself (mock the client instead)

### Fixtures and Factories

**Test Data:**

Defined inline in test file:

```typescript
const mockEmptyResponse = {
    content: {
        items: [],
        pageInfo: {
            hasNextPage: false,
            hasPreviousPage: false,
            startCursor: null,
            endCursor: null
        },
        totalCount: 0
    }
};

const mockDataResponse = {
    content: {
        items: [
            {
                id: '1',
                name: 'Test Video',
                url: 'https://youtube.com/watch?v=abc',
                contentType: 'YOUTUBE',
                // ... more fields
            }
        ],
        pageInfo: { /* ... */ },
        totalCount: 25
    }
};
```

**Test Wrapper Helper:**

Global wrapper for component tests (`tests/helpers/TestWrapper.svelte`):

```typescript
function renderWithQuery() {
    const queryClient = new QueryClient({
        defaultOptions: {
            queries: {
                retry: false,
                gcTime: 0,
                staleTime: 0
            },
            mutations: { retry: false }
        }
    });
    return render(TestWrapper, {
        props: {
            queryClient,
            component: ActivityTable,
            props: {}
        }
    });
}
```

### Coverage

**Requirements:** Not enforced in CI (no coverage gate)

**View Coverage:**
```bash
pnpm run test:coverage  # Generates coverage/ directory
```

**Detection of Gaps:** Test results show which lines/branches covered

## Global Test Setup

**Backend:** No global setup (each test is independent)

**Frontend:** `tests/setup.ts`

Contains:
- SvelteKit environment mocks (`$app/environment`, `$app/navigation`, `$app/stores`)
- `window.matchMedia` mock for responsive tests
- Global Jest-DOM assertions

```typescript
import '@testing-library/jest-dom/vitest';
import { vi } from 'vitest';

vi.mock('$app/environment', () => ({
    browser: true,
    dev: true,
    building: false
}));

vi.mock('$app/stores', () => ({
    page: readable({
        url: new URL('http://localhost'),
        params: {},
        // ... more page store properties
    }),
    // ... other stores
}));

Object.defineProperty(window, 'matchMedia', {
    writable: true,
    value: vi.fn().mockImplementation((query: string) => ({
        matches: false,
        media: query,
        // ... match media implementation
    }))
});
```

## Integration vs Unit Tests

**Backend:**

Unit tests (most common):
- No database dependency
- Mock repositories
- Fast execution
- Located in `test/{services,domain}/`

Integration tests (when DB needed):
- Marked with database setup
- Auto-skip via `t.Skip()` if database unavailable
- Test against real PostgreSQL (via docker-compose)
- Located in `test/{repositories}/`

**Frontend:**

Unit tests (most common):
- Test utility functions and formatting
- Mock external API clients
- No browser rendering (JSDOM)

Component tests (limited):
- Test component instantiation and basic rendering
- Mock complex child components (AG Grid)
- TanStack Query and AG Grid have known JSDOM limitations
- Full integration testing via manual verification or Playwright

## Known Limitations

**Frontend Component Testing:**

AG Grid + TanStack Query rendering in JSDOM has limitations:
- AG Grid mocked component doesn't trigger lifecycle (onGridReady)
- TanStack Query queries don't execute in JSDOM
- Tests verify component instantiation, not full behavior

Per test file comment in `tests/components/ActivityTable.test.ts`:
```typescript
/**
 * ActivityTable Tests
 *
 * LIMITATION: AG Grid + TanStack Query rendering in JSDOM has known limitations.
 * AG Grid's mocked component doesn't trigger lifecycle hooks (onGridReady), and
 * TanStack Query queries don't execute in this test environment.
 *
 * These tests verify component instantiation. Full integration testing requires
 * browser environment (manual verification or Playwright E2E tests).
 */
```

## Test Types

**Go Unit Tests:**
- Scope: Single function/method
- Dependencies: Mocked
- Execution: <100ms typically
- Example: `TestGetByID_Success` tests only service GetByID logic

**Go Integration Tests:**
- Scope: Service + repository + database
- Dependencies: Real PostgreSQL (skipped if unavailable)
- Execution: 100ms - 1s
- Example: Tests in `test/repositories/` exercise full CRUD with DB

**TypeScript Unit Tests:**
- Scope: Utility functions
- Dependencies: Mocked or none
- Execution: <50ms
- Example: `formatDuration`, `formatDate` tests

**TypeScript Component Tests:**
- Scope: Component instantiation
- Dependencies: Child components mocked
- Execution: 50-200ms
- Example: `Header.test.ts` verifies rendering and props
- Note: Limited by JSDOM (see Known Limitations)

## Common Patterns

**Async Testing (Go):**

```go
func TestCreateFromYouTube_Success(t *testing.T) {
    // Mocks return results immediately (no goroutines)
    repo := &mockContentRepository{
        createFn: func(ctx context.Context, content *domain.Content) (*domain.Content, error) {
            return content, nil
        },
    }
    svc := services.NewContentService(repo, &mockYouTubeClient{})
    
    result, err := svc.CreateFromYouTube(context.Background(), "https://...", 1)
    require.NoError(t, err)
    assert.NotNil(t, result)
}
```

**Async Testing (TypeScript):**

```typescript
it('loads data on mount', async () => {
    mockRequest.mockResolvedValue(mockDataResponse);
    
    render(ActivityTable);
    
    // Wait for query to resolve
    await waitFor(() => {
        expect(mockRequest).toHaveBeenCalled();
    });
});
```

**Error Testing (Go):**

```go
func TestGetByID_InvalidID(t *testing.T) {
    svc := services.NewContentService(&mockContentRepository{}, &mockYouTubeClient{})
    
    _, err := svc.GetByID(context.Background(), -1)
    
    require.Error(t, err)
    assert.ErrorIs(t, err, domain.ErrInvalidInput)
}
```

**Error Testing (TypeScript):**

```typescript
it('returns error message on failure', async () => {
    mockRequest.mockRejectedValue(new Error('Network error'));
    
    render(AddVideoPopover);
    
    await waitFor(() => {
        expect(screen.getByText(/error/i)).toBeInTheDocument();
    });
});
```

---

*Testing analysis: 2026-02-16*
