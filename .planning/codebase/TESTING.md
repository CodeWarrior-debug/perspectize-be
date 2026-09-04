# Testing Patterns

**Analysis Date:** 2026-09-04

This monorepo has two independent test suites: **Go backend** (`backend/test/`, `testify` + stdlib `testing`) and **SvelteKit frontend** (`frontend/tests/`, Vitest + Testing Library + Playwright browser mode).

---

## Backend: Go Testing

### Test Framework

**Runner:** Go stdlib `testing` package, run via `make test` (wraps `go test ./...`).

**Assertion library:** `github.com/stretchr/testify` — `assert` (soft, continues on failure) and `require` (fatal, stops test on failure).

**Run commands:**
```bash
cd backend
make test              # All tests
make test-coverage     # Coverage → coverage.html / coverage.out
go test ./...          # Equivalent to make test
go test ./test/services/... -run TestGetByID_Success   # Single test
```

### Test File Organization

**Location:** Centralized under `backend/test/`, mirroring the package structure being tested — NOT co-located with source files.
```
backend/test/
├── config/config_test.go
├── database/postgres_test.go
├── domain/{content,errors,perspective,user}_test.go
├── graphql/intid_test.go
├── resolvers/{content_resolver,me_resolver,helpers}_test.go
├── services/{content_service,perspective_service,user_service}_test.go
├── services/{content_service_bench,perspective_service_bench}_test.go
└── youtube/{cache,client,parser}_test.go
```

**Naming:** `<source_area>_test.go`, e.g. `content_service.go` (source) → `test/services/content_service_test.go`. Benchmark files use `_bench_test.go` suffix. Package declared as `<area>_test` (e.g. `package services_test`), forcing tests to only use the package's exported API — no internal white-box access.

### Test Structure

**Flat `TestXxx` functions**, not `t.Run` table-driven by default (though `t.Run` is used where the project-conventions skill recommends it for validation-heavy cases). Naming pattern is `Test<Method>_<Scenario>`:
```go
func TestGetByID_Success(t *testing.T) { ... }
func TestGetByID_NotFound(t *testing.T) { ... }
func TestGetByID_InvalidID_Zero(t *testing.T) { ... }
func TestGetByID_InvalidID_Negative(t *testing.T) { ... }
func TestGetByID_RepositoryError(t *testing.T) { ... }
```
Tests are grouped in the file with `// --- MethodName Tests ---` comment banners.

**Assertion patterns:**
```go
require.NoError(t, err)             // fatal — stop test if setup/precondition fails
assert.Equal(t, expected, result)   // soft — collect all failures
assert.True(t, errors.Is(err, domain.ErrInvalidInput))
assert.Contains(t, err.Error(), "content id must be a positive integer")
```
`require` is used immediately after an operation whose failure would make later assertions meaningless (e.g. `require.NoError` before inspecting the result); `assert` is used for the actual behavioral checks.

### Mocking

**No mocking framework/codegen** (no gomock/mockery). Mocks are **hand-written structs with function fields**, implementing the port interface directly, defined at the top of the relevant `_test.go` file:
```go
type mockContentRepository struct {
    createFn           func(ctx context.Context, content *domain.Content) (*domain.Content, error)
    getByIDFn          func(ctx context.Context, id int) (*domain.Content, error)
    getByURLFn         func(ctx context.Context, url string) (*domain.Content, error)
    getOrCreateByURLFn func(ctx context.Context, content *domain.Content, refreshOnConflict bool) (*domain.Content, bool, error)
    updateMetadataFn   func(ctx context.Context, id int, name string, response json.RawMessage, length *int) (*domain.Content, error)
    listFn             func(ctx context.Context, params domain.ContentListParams) (*domain.PaginatedContent, error)
}

func (m *mockContentRepository) GetByID(ctx context.Context, id int) (*domain.Content, error) {
    if m.getByIDFn != nil {
        return m.getByIDFn(ctx, id)
    }
    return nil, domain.ErrNotFound   // sensible zero-value default
}
```
Each test only sets the `*Fn` fields it needs; unset fields fall back to a default (usually `domain.ErrNotFound` or a no-op). This lets call-site assertions run inside the mock function itself:
```go
getByIDFn: func(ctx context.Context, id int) (*domain.Content, error) {
    assert.Equal(t, 1, id)   // assert on the args the service passed in
    return expected, nil
},
```

**Gotcha (documented in backend `CLAUDE.md`):** adding a method to a port interface (e.g. `ListAll` on `UserRepository`) breaks compilation of every mock implementing that interface — grep `test/` and update all mocks when extending a port.

**What's mocked:** repository ports (`repositories.ContentRepository`), external service ports (`portservices.YouTubeClient`). What's NOT mocked: the service under test itself (real `services.ContentService`), domain logic, error values.

### Fixtures and Factories

No dedicated fixture/factory package — test data is constructed inline as literal `&domain.Content{...}` / `&portservices.VideoMetadata{...}` structs per test. Shared helpers (e.g. `clearConfigEnvVars`) live alongside the tests that need them, e.g. `test/config/config_test.go`.

### Coverage

No enforced coverage threshold for backend (unlike frontend). View coverage:
```bash
make test-coverage   # generates backend/coverage.html and coverage.out
```

### Test Types

**Unit tests** (`test/services/`, `test/domain/`, `test/graphql/`): mock all dependencies, no real DB/network. This is the large majority of the suite.

**Integration tests** (`test/database/postgres_test.go`): connect to a real PostgreSQL instance; **auto-skip via `t.Skip()`** when the DB is unavailable rather than failing:
```go
db, err := database.ConnectGORM(dsn, poolCfg)
if err != nil {
    t.Skip("Skipping test - PostgreSQL not available. Run 'make docker-up' to start database.")
}
```

**Env isolation:** tests that load config must explicitly clear relevant env vars with `t.Setenv("KEY", "")` to avoid bleed-through from the developer's shell/CI env — see `clearConfigEnvVars` helper in `test/config/config_test.go`.

**Benchmarks:** `test/services/content_service_bench_test.go`, `perspective_service_bench_test.go` — standard `func BenchmarkXxx(b *testing.B)`, run via `go test -bench=.`.

### Common Patterns

**Error-path testing via `errors.Is` against domain sentinels:**
```go
result, err := svc.GetByID(context.Background(), 0)
assert.Nil(t, result)
require.Error(t, err)
assert.True(t, errors.Is(err, domain.ErrInvalidInput))
assert.Contains(t, err.Error(), "content id must be a positive integer")
```

**Dual-return semantics testing** (e.g. "already exists" returns both a non-nil result AND an error):
```go
result, err := svc.CreateFromYouTube(context.Background(), canonicalURL, 1)
require.NotNil(t, result)
require.Error(t, err)
assert.True(t, errors.Is(err, domain.ErrAlreadyExists))
```

**Capturing call arguments via closures** to assert what the service passed downstream (see `capturedID`/`capturedName`/`capturedLength` pattern in `content_service_test.go` `TestUpdateSourceData_Success`).

---

## Frontend: Vitest Testing

### Test Framework

**Runner:** Vitest, configured as **two projects** inside `frontend/vite.config.ts` (`test.projects`):
- `unit` project — jsdom environment, `tests/**/*.{test,spec}.{js,ts}` excluding `tests/browser/**`
- `browser` project — real Chromium via `@vitest/browser-playwright`, defined separately in `frontend/vitest.config.browser.ts`, scope `tests/browser/**/*.test.ts`

**Component testing:** `@testing-library/svelte` (+ `@testing-library/jest-dom/vitest` matchers).

**Run commands** (from `frontend/`, or `pnpm --dir frontend`):
```bash
pnpm run test            # unit project, watch mode
pnpm run test:run        # unit project, single run (CI/verification)
pnpm run test:browser    # browser project, single run (real Chromium, AG Grid etc.)
pnpm run test:browser:watch
pnpm run test:all        # both projects
pnpm run test:coverage   # unit project + v8 coverage
```
`pnpm exec`/`pnpm run` must be invoked from `frontend/` (or with `--dir frontend`) — running from repo root fails with `ERR_PNPM_RECURSIVE_EXEC_NO_PACKAGE`.

### Test File Organization

**Location:** Parallel `tests/` tree, NOT co-located with `src/`:
```
frontend/tests/
├── setup.ts                      # global mocks (see below), loaded for the `unit` project
├── unit/                         # pure logic: utils, hooks, stores, formatting
│   ├── hooks-useCreateClaim.test.ts
│   ├── formatting.test.ts
│   ├── grid-config.test.ts
│   └── ...
├── components/                   # Svelte component rendering tests
│   ├── ActivityTable.test.ts
│   ├── FilterBar.test.ts
│   └── fixtures/SearchBarHost.svelte   # per-test host component when needed
├── browser/                      # real-browser (Playwright) tests, own project config
│   ├── ag-grid-integration.test.ts
│   ├── mocks/{app-environment,app-navigation,app-stores}.ts
│   └── fixtures/AGGridTestHarness.svelte
├── fixtures/                     # (currently empty except .gitkeep)
└── helpers/
    ├── render.ts                 # renderComponent() + expectClasses() helpers
    └── TestWrapper.svelte        # dynamic-component test wrapper
```

**Naming:** `<subject>.test.ts`, hook tests prefixed `hooks-` (`hooks-useCreateClaim.test.ts`, `hooks-useAddVideo.test.ts`), query-layer tests prefixed `queries-` (`queries-content.test.ts`, `queries-keys.test.ts`).

### Test Structure

Standard Vitest BDD-style `describe`/`it`, nested by scenario group:
```ts
import { describe, it, expect } from 'vitest';
import { formatDuration } from '$lib/utils/formatting';

describe('formatDuration', () => {
    it('returns dash for null length', () => {
        expect(formatDuration(null, null)).toBe('—');
    });
    it('formats seconds as minutes:seconds', () => {
        expect(formatDuration(300, 'seconds')).toBe('5:00');
    });
});
```
For hook tests, nested `describe` blocks group by concern (`'hook initialization'`, `'mutationFn'`, `'onSuccess callback'`, `'onError callback'`), each with its own `beforeEach` re-invoking the hook under test.

**Component render helper** (`tests/helpers/render.ts`):
```ts
export function renderComponent<T extends Record<string, any>>(
    component: Component<T>,
    props?: Partial<T>,
): RenderResult<Component<T>> {
    // @ts-expect-error - Testing Library type mismatch with Svelte 5
    return render(component, props ? { props } : {});
}

export function expectClasses(element: HTMLElement, ...classes: string[]) { ... }
```

### Mocking

**`vi.mock` + `vi.hoisted`** for module-level mocks that need to be referenced inside the mock factory (avoids the hoisting trap where `vi.mock` factories run before top-level `const` declarations):
```ts
const { mockMutate, mockInvalidateQueries, mockToastSuccess, mockToastError } = vi.hoisted(() => ({
    mockMutate: vi.fn(),
    mockInvalidateQueries: vi.fn(),
    mockToastSuccess: vi.fn(),
    mockToastError: vi.fn(),
}));

let capturedMutationOptions: any;

vi.mock('@tanstack/svelte-query', () => ({
    createMutation: vi.fn((optionsFn: () => any) => {
        capturedMutationOptions = optionsFn();   // capture the function-wrapper options for direct invocation
        return { mutate: mockMutate, isPending: false };
    }),
    useQueryClient: vi.fn(() => ({ invalidateQueries: mockInvalidateQueries })),
}));

vi.mock('svelte-sonner', () => ({ toast: { success: mockToastSuccess, error: mockToastError } }));
vi.mock('$lib/queries/client', () => ({ graphqlRequest: vi.fn() }));
```
Hooks are then dynamically imported inside each test (`await import('$lib/queries/hooks/useCreateClaim')`) after mocks are registered, and `capturedMutationOptions.onSuccess()` / `.onError(new Error(...))` are invoked directly to test callback branches without going through real TanStack Query internals.

**Global mocks** in `tests/setup.ts` (auto-loaded for the `unit` project via `setupFiles`):
- `$app/environment`, `$app/navigation`, `$app/state`, `$app/stores` — SvelteKit runtime mocked so components can render outside a real app shell. `mockPageState` is exported mutable so individual tests can override the current URL.
- `IntersectionObserver` — custom mock that fires `isIntersecting: true` synchronously on `observe()`, so lazy-loading components (e.g. `VideoCard` thumbnails) render real content immediately without simulating scroll.
- `localStorage` — hand-rolled `Map`-backed `Storage` implementation (Node's native webstorage global shadows jsdom's real one in this environment and is non-functional without a `--localstorage-file` flag).
- `window.matchMedia` — defaults `matches: false` (desktop) so responsive components don't need per-test setup unless testing mobile breakpoints.

**What NOT to mock:** pure utility functions under test (`$lib/utils/*`) are called directly/unmocked; only their external dependencies (network, toast, SvelteKit runtime) are mocked.

### Fixtures and Factories

No dedicated factory library — inline object literals per test, following the same GraphQL-shaped data as backend responses (e.g. `{ createClaim: { id: '1', text: 'Test claim', userID: '42' } }`). Component-specific host/wrapper fixtures live beside their tests: `tests/components/fixtures/SearchBarHost.svelte`, `tests/browser/fixtures/AGGridTestHarness.svelte`.

**Dynamic component test wrapper gotcha (`tests/helpers/TestWrapper.svelte`):** pass the component under test as a dotted-member expression, `<wrapped.Comp {...props} />`. A bare `<component>` or `.default` renders nothing (empty comment placeholder) and can make assertions pass vacuously — see `frontend/CLAUDE.md` "Testing Gotchas" and prior regression in `tests/components/ActivityTable.test.ts`.

### Coverage

**Enforced thresholds** (`vite.config.ts`, `unit` project, v8 provider):
```ts
coverage: {
    provider: 'v8',
    reporter: ['text', 'json', 'html'],
    exclude: [
        'node_modules/', '.svelte-kit/', '**/*.d.ts', '**/*.config.*', '**/setup.ts',
        'tests/helpers/**', 'src/lib/components/shadcn/**', 'src/routes/**',
        'src/lib/components/ActivityTable.svelte',   // excluded — see AG Grid testing strategy below
    ],
    thresholds: { lines: 80, functions: 75, branches: 75, statements: 80 },
}
```
View coverage: `pnpm run test:coverage` → HTML report in `frontend/coverage/`.

### Test Types

**Unit tests** (`tests/unit/`): pure functions (formatting, ratings, grid-config, URL state), hooks (TanStack Query mutations), stores, GraphQL query-key builders. Majority of the suite; run in jsdom, fast.

**Component tests** (`tests/components/`): render real Svelte 5 components via Testing Library + jsdom, assert DOM output/classes/interaction.

**Browser tests** (`tests/browser/`, separate `vitest.config.browser.ts` project): real Chromium via Playwright provider — needed specifically for AG Grid, which does not render/initialize in jsdom (no Grid API, no lifecycle hooks, no cell rendering). AG Grid logic is instead extracted into plain TS in `$lib/utils/grid-config.ts` and `$lib/utils/formatting.ts` and unit-tested there; only true grid integration (filter UI, sort clicks, responsive `$effect` column visibility) goes through browser mode or (future) Playwright E2E.

**E2E:** Not present as a separate framework — Vitest Browser Mode fills this role currently (see AG Grid testing strategy note above).

### Known Limitations / Gotchas

- **Date/timezone:** `formatDate`/`formatDateCompact` use `toLocaleDateString` (local TZ) — always seed test dates at midday UTC (`T12:00:00Z`), never midnight, to avoid the date shifting to the previous day in US timezones. See `formatDateTime` test comment referencing this.
- **AG Grid cell renderers** inherit `white-space: nowrap` from `.ag-cell` — DOM assertions on wrapping/clamped text need to account for an explicit `whitespace-normal` override on the element under test.
- Coverage thresholds exclude `ActivityTable.svelte` and all `shadcn/` primitives and `src/routes/**` — do not expect these to move the coverage numbers; new business logic should go in `$lib/utils/`, `$lib/queries/`, or extracted, coverage-counted components instead.

### Common Patterns

**Async/callback testing (mutation `onSuccess`/`onError`):**
```ts
capturedMutationOptions.onSuccess();
expect(mockToastSuccess).toHaveBeenCalledWith('Claim created');
expect(mockInvalidateQueries).toHaveBeenCalledWith(
    expect.objectContaining({ queryKey: expect.arrayContaining(['content', 'list']) })
);
```

**Error-message-based branching test (case-insensitive substring match on caught error):**
```ts
capturedMutationOptions.onError(new Error('PARENT CONTENT NOT FOUND in the database'));
expect(mockToastError).toHaveBeenCalledWith('Parent content not found');
```

**Pure-function edge-case sweep (null/empty/boundary values are always tested explicitly)** — see `formatting.test.ts` `describe('formatCount', ...)`, `describe('getSourceDataCooldown', ...)` for the TTL-boundary-exactly pattern.

---

*Testing analysis: 2026-09-04*
