# Testing Patterns

**Analysis Date:** 2026-02-04

## Test Framework

**Runner:**
- Go standard `testing` package
- Config: `Makefile` target `test` runs all tests with coverage
- No external test runner (e.g., no testify for test discovery, but used for assertions)

**Assertion Library:**
- `github.com/stretchr/testify/assert` - non-fatal assertions (`assert.Equal()`, `assert.Nil()`, `assert.Contains()`)
- `github.com/stretchr/testify/require` - fatal assertions (`require.NoError()`, `require.Equal()`) that stop test on failure

**Run Commands:**
```bash
make test                     # Run all tests with coverage report
make test-coverage            # Run tests with HTML coverage output (generates coverage.html)
go test ./...                 # Run all tests manually
go test -v -run TestFunctionName ./path/to/package  # Run single test
go test ./internal/core/services -v  # Run tests in specific package
```

## Test File Organization

**Location:**
- Tests are co-located in `test/` directory at project root (not in same directory as source code)
- Structure mirrors source: `test/services/`, `test/domain/`, `test/repositories/`, `test/resolvers/`, `test/config/`, `test/database/`

**Naming:**
- Test files: `{subject}_test.go` (e.g., `user_service_test.go`, `content_resolver_test.go`)
- Package: `{subject}_test` (e.g., `package services_test`, `package domain_test`, `package resolvers_test`)

**Structure:**
```
perspectize-go/
├── test/
│   ├── services/
│   │   ├── user_service_test.go
│   │   ├── content_service_test.go
│   │   └── perspective_service_test.go
│   ├── domain/
│   │   ├── user_test.go
│   │   ├── content_test.go
│   │   ├── perspective_test.go
│   │   └── errors_test.go
│   ├── resolvers/
│   │   └── content_resolver_test.go
│   ├── config/
│   │   └── config_test.go
│   ├── database/
│   │   └── postgres_test.go
│   └── youtube/
│       └── parser_test.go
├── internal/
│   ├── core/
│   │   ├── domain/
│   │   ├── ports/
│   │   └── services/
│   └── adapters/
└── cmd/
```

## Test Structure

**Suite Organization:**
```go
package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yourorg/perspectize-go/internal/core/domain"
	"github.com/yourorg/perspectize-go/internal/core/services"
)

// Mock implementations (defined per test file)
type mockUserRepository struct {
	createFn func(ctx context.Context, user *domain.User) (*domain.User, error)
	// ... other mock methods
}

func (m *mockUserRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	if m.createFn != nil {
		return m.createFn(ctx, user)
	}
	user.ID = 1
	return user, nil
}

// --- Test Groups (with comment separators) ---

func TestCreate_Success(t *testing.T) {
	repo := &mockUserRepository{
		createFn: func(ctx context.Context, user *domain.User) (*domain.User, error) {
			user.ID = 1
			return user, nil
		},
	}

	svc := services.NewUserService(repo)
	result, err := svc.Create(context.Background(), "testuser", "test@example.com")

	require.NoError(t, err)
	assert.Equal(t, 1, result.ID)
	assert.Equal(t, "testuser", result.Username)
}

func TestCreate_UsernameEmpty(t *testing.T) {
	repo := &mockUserRepository{}
	svc := services.NewUserService(repo)

	result, err := svc.Create(context.Background(), "", "test@example.com")

	assert.Nil(t, result)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrInvalidInput))
	assert.Contains(t, err.Error(), "username is required")
}

// --- Subtests for variations ---

func TestCreate_EmailInvalid(t *testing.T) {
	testCases := []string{
		"notanemail",
		"@example.com",
		"test@",
	}

	repo := &mockUserRepository{}
	svc := services.NewUserService(repo)

	for _, email := range testCases {
		t.Run(email, func(t *testing.T) {
			result, err := svc.Create(context.Background(), "testuser", email)

			assert.Nil(t, result)
			require.Error(t, err)
			assert.True(t, errors.Is(err, domain.ErrInvalidInput))
		})
	}
}
```

**Patterns:**
- Mock implementations defined at top of test file with function pointers for each method
- Mock functions accept `func(...)` as method field to customize behavior per test
- Test names: `Test{FunctionName}_{Scenario}` (e.g., `TestCreate_Success`, `TestCreate_UsernameEmpty`, `TestUserGetByID_NotFound`)
- Logical groupings with comment separators: `// --- Create Tests ---`
- Sub-tests for variations: `t.Run(scenario, func(t *testing.T) { ... })`
- Setup in test: create repo mock with specific behavior, inject into service, call method, assert
- Fatal assertions (`require.NoError()`) used first to stop on unexpected errors, then non-fatal checks (`assert.Equal()`) for specific values

## Mocking

**Framework:** Hand-written mocks (no mock generation library used)

**Patterns:**
```go
// Mock struct with function pointers for each method
type mockUserRepository struct {
	createFn        func(ctx context.Context, user *domain.User) (*domain.User, error)
	getByIDFn       func(ctx context.Context, id int) (*domain.User, error)
	getByUsernameFn func(ctx context.Context, username string) (*domain.User, error)
}

// Implement interface methods
func (m *mockUserRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	if m.createFn != nil {
		return m.createFn(ctx, user)
	}
	// Default behavior (needed for tests that don't set custom function)
	user.ID = 1
	return user, nil
}

// In test: inject mock with custom behavior
repo := &mockUserRepository{
	createFn: func(ctx context.Context, user *domain.User) (*domain.User, error) {
		user.ID = 1
		return user, nil
	},
}
```

**What to Mock:**
- Repository interfaces (core/ports/repositories)
- External service clients (e.g., YouTube API client)
- Validation is tested, not mocked

**What NOT to Mock:**
- Domain models
- Service business logic (test through service methods)
- Standard library functions

**Multiple Mock Definitions:**
- Each test file that uses mocks defines its own mock types
- Mocks are specific to test scenario (e.g., `services_test.go` has `mockUserRepository`, `resolvers_test.go` has same-named `mockUserRepository`)
- No shared mock file (patterns may vary between service and resolver tests)

## Fixtures and Factories

**Test Data:**
- Fixtures are created inline in tests, not in separate files
- Common test data: hardcoded values (e.g., `"testuser"`, `"test@example.com"`)
- No factory functions (helpers like `mockUserRepository` used instead)

**Location:**
- Mock implementations and helpers at top of test file
- Test data embedded in test functions or in mock callback functions

## Coverage

**Requirements:** No enforced coverage target in configuration
- Coverage report generated: `go tool cover -html=coverage.out -o coverage.html`
- Coverage tracked at package level: `-coverpkg=./internal/...,./pkg/...`

**View Coverage:**
```bash
make test-coverage
# Generates coverage.out and coverage.html
# Open with: open coverage.html  (macOS)
```

## Test Types

**Unit Tests:**
- Scope: Individual service/repository methods
- Approach: Mocked dependencies, fast (no database)
- Files: `test/services/`, `test/domain/`
- Examples: `user_service_test.go` tests UserService methods with mocked repositories

**Integration Tests:**
- Scope: Adapters against real database
- Approach: Docker Compose PostgreSQL or testcontainers
- Guarded with `t.Skip()` when prerequisites unavailable
- Files: `test/database/`, `test/repositories/` (if added)
- Example: `config_test.go` loads actual config file and verifies DSN generation

**E2E Tests:**
- Framework: Not used in current codebase
- Future: Could use `httptest` for GraphQL resolver testing (see `content_resolver_test.go` pattern with `httptest.NewRequest()`)

## Environment Isolation (IMPORTANT)

**Pattern:** Tests that load configuration must clear environment variables to avoid leakage from Makefile.

**Implementation:**
```go
// clearConfigEnvVars ensures config-relevant env vars are empty so tests
// run against config file values only. t.Setenv restores originals on cleanup.
func clearConfigEnvVars(t *testing.T) {
	t.Helper()
	for _, key := range []string{"DATABASE_URL", "DATABASE_PASSWORD", "YOUTUBE_API_KEY"} {
		t.Setenv(key, "")
	}
}

// In test:
func TestLoad_RealConfigFile(t *testing.T) {
	clearConfigEnvVars(t)
	// Now env vars are cleared, test loads config from file only
	cfg, err := config.Load("../../config/config.example.json")
	assert.NoError(t, err)
}
```

**Why needed:** Makefile exports `DATABASE_URL` etc. during `make test` - if not cleared, tests that verify config file loading will pass but shouldn't (env var shadow config file).

**Auto-cleanup:** `t.Setenv()` automatically restores original values after test completes.

## Common Patterns

**Async Testing:**
- Most service tests are synchronous (accept `context.Background()`)
- Context passed to methods but not explicitly tested for cancellation
- Example: `svc.Create(context.Background(), "testuser", "test@example.com")`

**Error Testing:**
```go
func TestCreate_EmailInvalid(t *testing.T) {
	repo := &mockUserRepository{}
	svc := services.NewUserService(repo)

	result, err := svc.Create(context.Background(), "testuser", "notanemail")

	assert.Nil(t, result)  // Resource should be nil on error
	require.Error(t, err)  // Error must be present
	assert.True(t, errors.Is(err, domain.ErrInvalidInput))  // Check error type
	assert.Contains(t, err.Error(), "invalid email format")  // Check error message
}
```

**Domain Validation:**
- Domain logic tested: `domain_test.go` tests struct field values and zero-value behavior
- Example: `test/domain/user_test.go` tests User struct construction and field initialization

**Test Naming for Readability:**
- Format: `Test{Method}_{Scenario}` makes test purpose clear
- Examples: `TestCreate_Success`, `TestCreate_UsernameEmpty`, `TestUserGetByID_InvalidID_Zero`
- Grouping with `// --- Category ---` comments separates logical test suites

## Test Status & Adoption

**Fully Adopted:**
- Unit test structure with mocks in services
- Error handling assertions with `errors.Is()`
- Mock repository pattern with function pointers
- Test naming convention `Test{Function}_{Scenario}`

**Partially Adopted:**
- Integration tests: Only `config_test.go` tests actual file loading; database adapter tests are sparse (no direct repository integration tests against real DB yet)
- GraphQL resolver testing: Pattern exists (`content_resolver_test.go`) with `httptest` but minimal coverage

**Not Yet Adopted:**
- E2E tests (pattern mentioned in CLAUDE.md, not implemented)
- Structured logging in tests (logging not tested)
- Table-driven tests (some tests use sub-tests with `t.Run()` but could be more systematic)

---

*Testing analysis: 2026-02-04*
