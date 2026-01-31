---
name: project-conventions
description: Project-specific conventions and patterns for the Perspectize backend. Use when working on any Perspectize Go code to ensure consistency with established patterns.
---

# Perspectize Project Conventions

This skill contains project-specific patterns and conventions for the Perspectize backend. Follow these guidelines for all code contributions.

## Project Philosophy

> "An app dedicated to empowering in-depth thought, whether deep in reach or in detail."

- **Contemplative over reactive**: Encourage thoughtful engagement
- **Explicit over implicit**: No magic, clear code flow
- **Standard library first**: Minimize dependencies
- **Developer delight**: Code should be enjoyable to work with

## Perspective Dimensions

The core data model uses multi-dimensional ratings:

| Dimension | Range | Storage | Display |
|-----------|-------|---------|---------|
| Quality | 0-10000 | INTEGER | 0.00% - 100.00% |
| Agreement | 0-10000 | INTEGER | 0.00% - 100.00% |

**Why integers?** Avoids floating-point precision issues while allowing 0.01% granularity.

```go
// Convert for display
func FormatAsPercent(value int) string {
    return fmt.Sprintf("%.2f%%", float64(value)/100)
}

// Validate range
func ValidateRating(value int) error {
    if value < 0 || value > 10000 {
        return fmt.Errorf("rating must be between 0 and 10000, got %d", value)
    }
    return nil
}
```

## Naming Conventions

### Package Names
```go
// Good
package handlers
package services
package repositories

// Bad
package handler      // Use plural
package svc          // Don't abbreviate
package persistence  // Use 'repositories'
```

### Type Names
```go
// Services
type PerspectiveService struct {}
type ContentService struct {}
type YouTubeService struct {}

// Repositories
type PerspectiveRepository struct {}
type ContentRepository struct {}

// Handlers
type PerspectiveHandler struct {}
type ContentHandler struct {}

// Models - singular, matches table
type Perspective struct {}
type Content struct {}
type User struct {}
```

### Method Names
```go
// Repository methods
GetByID(ctx, id)           // Single by ID
GetByUserID(ctx, userID)   // Single by foreign key
List(ctx, opts)            // Multiple with options
Create(ctx, model)         // Insert
Update(ctx, model)         // Update
Delete(ctx, id)            // Delete

// Service methods - may combine operations
GetPerspectiveWithContent(ctx, id)
CreatePerspectiveForUser(ctx, userID, input)
```

## Error Handling

### Error Types

```go
// internal/core/domain/errors.go

var (
    ErrNotFound      = errors.New("not found")
    ErrValidation    = errors.New("validation failed")
    ErrUnauthorized  = errors.New("unauthorized")
    ErrForbidden     = errors.New("forbidden")
    ErrConflict      = errors.New("resource conflict")
)

// Wrap with context
return fmt.Errorf("get perspective %d: %w", id, err)
```

### Error Handling Pattern

```go
// In resolvers/handlers
content, err := h.service.GetByID(ctx, id)
if err != nil {
    switch {
    case errors.Is(err, domain.ErrNotFound):
        return nil, fmt.Errorf("content not found")
    case errors.Is(err, domain.ErrValidation):
        return nil, err
    default:
        h.logger.Error("unexpected error", "error", err)
        return nil, fmt.Errorf("internal error")
    }
}
```

## HTTP Response Patterns

### Success Responses

```go
// JSON response helper
func respondJSON(w http.ResponseWriter, status int, data any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(data)
}

// Usage
respondJSON(w, http.StatusOK, content)           // 200 with data
respondJSON(w, http.StatusCreated, perspective)   // 201 with created
w.WriteHeader(http.StatusNoContent)               // 204 no body
```

### Error Responses

```go
type ErrorResponse struct {
    Code    string `json:"code"`
    Message string `json:"message"`
}

func respondError(w http.ResponseWriter, status int, message string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(ErrorResponse{
        Code:    http.StatusText(status),
        Message: message,
    })
}
```

## Database Patterns

### Query Style

```go
// Use named queries for readability
const getPerspectiveByID = `
    SELECT id, content_id, user_id, quality, agreement, claim, created_at
    FROM perspectives
    WHERE id = $1
`

func (r *Repository) GetByID(ctx context.Context, id int) (*domain.Perspective, error) {
    var p domain.Perspective
    err := r.db.GetContext(ctx, &p, getPerspectiveByID, id)
    if err == sql.ErrNoRows {
        return nil, domain.ErrNotFound
    }
    return &p, err
}
```

### Transaction Pattern

```go
func (r *Repository) CreateWithAudit(ctx context.Context, p *domain.Perspective) error {
    tx, err := r.db.BeginTxx(ctx, nil)
    if err != nil {
        return fmt.Errorf("begin transaction: %w", err)
    }
    defer tx.Rollback() // No-op if committed

    // Insert perspective
    _, err = tx.ExecContext(ctx, insertPerspective, p.ContentID, p.Quality)
    if err != nil {
        return fmt.Errorf("insert perspective: %w", err)
    }

    // Insert audit log
    _, err = tx.ExecContext(ctx, insertAuditLog, "perspective_created", p.ID)
    if err != nil {
        return fmt.Errorf("insert audit: %w", err)
    }

    return tx.Commit()
}
```

## Logging Conventions

```go
// Use structured logging with slog
logger := slog.Default()

// Info level - normal operations
logger.Info("perspective created",
    "perspective_id", p.ID,
    "user_id", p.UserID,
    "quality", p.Quality,
)

// Error level - include error
logger.Error("failed to create perspective",
    "error", err,
    "user_id", userID,
    "content_id", contentID,
)

// Debug level - verbose details
logger.Debug("database query executed",
    "query", "SELECT * FROM perspectives",
    "duration_ms", duration.Milliseconds(),
)
```

## Testing Conventions

### Test File Naming

```
perspective_service.go      -> perspective_service_test.go
perspective_service.go      -> perspective_service_integration_test.go (with build tag)
```

### Test Function Naming

```go
// Unit tests
func TestPerspectiveService_Create(t *testing.T)
func TestPerspectiveService_Create_ValidationError(t *testing.T)

// Table-driven - use t.Run
func TestPerspectiveService_Create(t *testing.T) {
    t.Run("valid input", func(t *testing.T) { ... })
    t.Run("invalid quality", func(t *testing.T) { ... })
}
```

### Assertions

```go
// Use testify/assert for soft assertions
assert.Equal(t, expected, actual)
assert.NoError(t, err)
assert.Contains(t, str, substring)

// Use testify/require for fatal assertions
require.NoError(t, err) // Test stops if fails
require.NotNil(t, result)
```

## Import Organization

```go
import (
    // Standard library
    "context"
    "encoding/json"
    "net/http"

    // Third-party
    "github.com/go-chi/chi/v5"
    "github.com/jmoiron/sqlx"

    // Internal
    "github.com/codewarrior-debug/perspectize-be/internal/core/domain"
    "github.com/codewarrior-debug/perspectize-be/internal/core/services"
)
```

## Configuration

```go
// Use environment variables with sensible defaults
type Config struct {
    Port        string
    DatabaseURL string
    LogLevel    string
    Env         string
}

func LoadConfig() *Config {
    return &Config{
        Port:        getEnv("PORT", "8080"),
        DatabaseURL: getEnv("DATABASE_URL", ""),
        LogLevel:    getEnv("LOG_LEVEL", "info"),
        Env:         getEnv("ENV", "development"),
    }
}

func getEnv(key, fallback string) string {
    if v := os.Getenv(key); v != "" {
        return v
    }
    return fallback
}
```

## Code Quality Checklist

Before submitting code:

- [ ] `make fmt` passes
- [ ] `make lint` passes with no errors
- [ ] `make test` passes
- [ ] New code has tests
- [ ] Error messages are helpful
- [ ] No hardcoded credentials
- [ ] Context is propagated
- [ ] Errors are wrapped with context
