# Perspectize Go Backend - Architecture Plan

## Project Overview

Migration of the Perspectize backend from C# ASP.NET Core to Go, maintaining all existing functionality while leveraging Go's simplicity and performance. This is a **hobby project** prioritizing developer experience, learning, and simplicity over maximum scale.

## Current C# System Analysis

### Technology Stack
- **Framework**: ASP.NET Core 9.0
- **Database**: PostgreSQL (currently using local testdb)
- **ORM/Data Access**: Entity Framework Core + Dapper (hybrid approach)
- **API Style**: RESTful with conventional routing
- **Configuration**: appsettings.json

### Current Feature Set

#### 1. Content Management
- **Models**: Content table with YouTube video metadata
- **Endpoints**:
  - `GET /content` - List all content
  - `GET /content/{name}` - Get specific content by name
- **Storage**: URL, name, content_type, length, length_units, JSONB response from YouTube API

#### 2. YouTube Integration
- **Service**: YouTubeService with API key from config
- **Endpoints**:
  - `GET /youtube/video?videoId={id}` - Fetch single video from YouTube API
  - `POST /youtube/videos` - Batch import videos (fetches from API, stores in DB)
  - `PUT /youtube/videos` - Update existing videos (currently reuses POST logic)
- **Features**:
  - Extract video ID from various YouTube URL formats
  - Convert ISO 8601 duration to seconds
  - Store full YouTube API response as JSONB
  - Upsert logic (find-then-update or insert)

#### 3. Perspectives System
- **Models**: Perspective table with user opinions/ratings on content
- **Endpoints**:
  - `GET /perspectives/{username}` - Get all perspectives for a user
  - `GET /perspectives/{id}` - Get single perspective by ID
  - `POST /perspectives` - Create multiple perspectives (batch)
  - `PUT /perspectives/{id}` - Update perspective (partial updates supported)
  - `DELETE /perspectives` - Delete multiple perspectives (batch)
- **Features**:
  - Rich rating system (quality, agreement, importance, confidence - all 0-10000 scale)
  - PostgreSQL custom domain `valid_integer_range` for validation
  - JSONB arrays for categorized_ratings
  - PostgreSQL arrays for parts and labels
  - Unique constraint on (claim, user_id)
  - Auto-updating timestamps via trigger
  - Foreign keys to content and users tables

#### 4. Database Schema
- **Tables**: content, perspectives, users
- **PostgreSQL Features Used**:
  - JSONB columns (response, categorized_ratings)
  - Array columns (parts, labels)
  - Custom domain types (valid_integer_range)
  - Triggers (update_updated_at)
  - Foreign keys with named constraints
  - Unique constraints

#### 5. Users (Minimal)
- **Table exists** with username and email fields
- **No authentication** - user_id referenced in perspectives
- **No user management endpoints** in current C# code

### Known TODOs from C# Code
- Content endpoint uses name instead of ID (names have spaces, can get long)
- YouTube PUT reuses POST logic (should be simpler update)
- YouTube uses find-then-update instead of ON CONFLICT upsert
- Browser window opens on startup (localhost:7253)

## Go Architecture Design

### Core Principles
1. **Simplicity First** - Standard library over frameworks where possible
2. **Learning-Friendly** - Clear, idiomatic Go code for someone new to the language
3. **Incremental Validation** - Small changes that can be tested independently
4. **Documentation** - Full OpenAPI/Swagger docs for API exploration
5. **Quality Tests** - Tests for each feature, but not excessive edge cases

### Technology Choices

#### 1. Web Framework: `net/http` + `chi` Router
**Choice**: Standard library + minimal router
```go
import (
    "net/http"
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
)
```
**Rationale**:
- `net/http` is production-ready and well-documented
- `chi` is lightweight, idiomatic, works with standard handlers
- Easier to understand than full frameworks (Fiber, Gin)
- Better for learning Go patterns

#### 2. Database: `database/sql` + `sqlx` + `pgx` Driver
**Choice**: Standard library + minimal extensions + best PostgreSQL driver
```go
import (
    "database/sql"
    _ "github.com/jackc/pgx/v5/stdlib"
    "github.com/jmoiron/sqlx"
)
```
**Rationale**:
- Direct SQL matches current Dapper usage in C# (familiar pattern)
- `pgx` is the best PostgreSQL driver (performance, features)
- `sqlx` adds struct scanning without heavy ORM
- Full control over queries and PostgreSQL-specific features
- **Why not GORM?**
  - Your C# code uses Dapper (direct SQL), not EF Core for queries - similar philosophy
  - GORM's AutoMigrate can't handle PostgreSQL-specific features (custom domains, triggers, JSONB arrays)
  - Complex queries with joins (like perspectives + content) are clearer in SQL
  - Harder to debug ORM-generated queries vs explicit SQL
  - Your current migrations use raw SQL for advanced features - ORMs fight this
  - **When GORM makes sense**: Simpler schemas, rapid prototyping, less PostgreSQL-specific features
  - **For this project**: Direct SQL gives you control over the advanced PostgreSQL features you're already using

#### 3. Migrations: `golang-migrate/migrate`
**Choice**: SQL-based migration tool
```go
import (
    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/postgres"
    "github.com/golang-migrate/migrate/v4/source/file"
)
```
**Rationale**:
- SQL migrations match EF Core approach
- Version control for schema changes
- Supports PostgreSQL-specific features
- Can port existing C# migrations directly

#### 4. Configuration: Standard Library `os` + `encoding/json`
**Choice**: Environment variables + JSON config file
```go
import (
    "os"
    "encoding/json"
)
```
**Rationale**:
- No external dependencies
- Matches appsettings.json approach
- 12-factor app compliant (env vars for secrets)

#### 5. HTTP Client: Standard Library `net/http`
**Choice**: Built-in HTTP client
```go
client := &http.Client{
    Timeout: 30 * time.Second,
}
```
**Rationale**:
- Sufficient for YouTube API calls
- No need for additional libraries

#### 6. Validation: `go-playground/validator/v10`
**Choice**: Struct tag-based validation
```go
import "github.com/go-playground/validator/v10"

type CreatePerspectiveRequest struct {
    UserId int    `json:"user_id" validate:"required"`
    Claim  string `json:"claim" validate:"required,max=255"`
    Quality *int  `json:"quality" validate:"omitempty,min=0,max=10000"`
}
```
**Rationale**:
- Similar to C# Data Annotations
- Standard in Go community
- Handles complex validation rules
- No good standard library alternative

#### 7. API Documentation: `swaggo/swag`
**Choice**: Comment-based OpenAPI generation
```go
import (
    "github.com/swaggo/swag"
    httpSwagger "github.com/swaggo/http-swagger"
)
```
**Rationale**:
- Generates Swagger/OpenAPI from code comments
- Documentation lives with code
- Provides Swagger UI for exploration
- Familiar to C# developers (similar to Swashbuckle)

#### 8. Logging: Standard Library `log/slog`
**Choice**: Structured logging (Go 1.21+)
```go
import "log/slog"

logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
```
**Rationale**:
- Modern structured logging in stdlib
- No external dependencies
- JSON output for production

#### 9. Testing: Hybrid Approach
**Choice**: Standard testing + assertion helpers + hybrid database testing strategy

**Dependencies**:
```go
import (
    "testing"                                    // Standard library
    "github.com/stretchr/testify/assert"        // Better assertions
    "github.com/DATA-DOG/go-sqlmock"            // Database mocking
)
```

**Testing Strategy - Hybrid Approach**:

**Unit Tests (Fast)** - Use `sqlmock` for repository tests
```go
// test/repositories/content_test.go
func TestContentRepository_GetByName(t *testing.T) {
    mockDB, mock, _ := sqlmock.New()
    defer mockDB.Close()

    db := sqlx.NewDb(mockDB, "sqlmock")
    repo := NewContentRepository(db)

    rows := sqlmock.NewRows([]string{"id", "name", "url"}).
        AddRow(1, "Test", "https://youtube.com/123")

    mock.ExpectQuery("SELECT (.+) FROM content WHERE name = ?").
        WithArgs("Test").
        WillReturnRows(rows)

    content, err := repo.GetByName("Test")

    assert.NoError(t, err)
    assert.Equal(t, "Test", content.Name)
    assert.NoError(t, mock.ExpectationsWereMet())
}
```

**Integration Tests (Optional, Slower)** - Use build tags for real database tests
```go
// test/integration/content_integration_test.go
// +build integration

func TestContentIntegration(t *testing.T) {
    // Use real PostgreSQL (via docker-compose or testcontainers)
    db := setupTestPostgres(t)
    defer db.Close()

    // Run actual migrations
    runMigrations(db)

    // Test complete workflows with real DB
    // Tests PostgreSQL-specific features (JSONB, arrays, triggers)
}
```

**Running Tests**:
```bash
# Fast unit tests only (default)
go test ./...

# Include integration tests (when needed)
go test -tags=integration ./...

# With coverage
go test -cover ./...
```

**Rationale**:
- **Unit tests with sqlmock**: Fast, run on every save, no dependencies
- **Integration tests optional**: Test PostgreSQL-specific features when needed
- **Best of both worlds**: Speed for TDD, accuracy for PostgreSQL features
- `testify/assert` improves test readability over plain Go assertions
- Similar to C# unit tests (mocking) + integration tests (real DB)
- **No in-memory SQLite**: PostgreSQL syntax differences would cause issues (JSONB, arrays, custom domains)

#### 10. Dependency Injection: Manual (Constructor Injection)
**Choice**: Explicit wiring in main.go
```go
// No framework - just pass dependencies
func main() {
    db := setupDatabase()
    repo := repositories.NewPerspectiveRepository(db)
    service := services.NewPerspectiveService(repo)
    controller := controllers.NewPerspectiveController(service)
    // ...
}
```
**Rationale**:
- Go idiom: explicit over implicit
- Easy to understand and debug
- No magic reflection
- Similar to manual DI (just no container)

### Deferred/Skipped Features

#### Not Implementing Initially
1. **Authentication/Authorization** - Skip for now
   - Future: May use Clerk for auth
   - For now: Use hardcoded user_id or accept from request

2. **Background Jobs/Cron Scheduler** - Not needed yet
   - No scheduled tasks in current implementation
   - Can add `robfig/cron` later if needed

3. **Extensive Edge Case Testing** - Focus on happy path + basic error cases
   - Quality tests for each feature
   - Not exhaustive test coverage initially

## Project Structure

```
perspectize-go/
├── cmd/
│   └── server/
│       └── main.go                 # Entry point, wiring, startup
├── internal/
│   ├── config/
│   │   └── config.go               # Configuration loading
│   ├── models/
│   │   ├── content.go              # Domain models
│   │   ├── perspective.go
│   │   └── user.go
│   ├── dto/
│   │   ├── content.go              # Request/Response DTOs
│   │   ├── perspective.go
│   │   └── youtube.go
│   ├── controllers/                # HTTP controllers (like ASP.NET)
│   │   ├── content.go
│   │   ├── perspective.go
│   │   └── youtube.go
│   ├── services/                   # Business logic
│   │   ├── perspective.go
│   │   └── youtube.go
│   ├── repositories/               # Database access
│   │   ├── content.go
│   │   └── perspective.go
│   ├── middleware/
│   │   ├── logger.go               # Request logging
│   │   └── recovery.go             # Panic recovery
│   └── validator/
│       └── validator.go            # Validation setup
├── pkg/
│   └── database/
│       └── postgres.go             # DB connection management
├── migrations/
│   ├── 000001_initial_create.up.sql
│   ├── 000001_initial_create.down.sql
│   ├── 000002_add_perspectives.up.sql
│   ├── 000002_add_perspectives.down.sql
│   └── ...
├── test/
│   ├── controllers/                # Controller tests
│   ├── services/                   # Service tests
│   └── repositories/               # Repository tests
├── docs/                           # Swagger generated docs
├── config/
│   ├── config.json                 # Default configuration
│   └── config.example.json
├── docker-compose.yml              # PostgreSQL for development
├── Makefile                        # Build/test/run commands
├── go.mod
├── go.sum
└── README.md
```

### Folder Conventions

**`cmd/`** - Application entry points
- Each application gets a subfolder
- `main.go` handles wiring and startup

**`internal/`** - Private application code
- Cannot be imported by other projects
- Core business logic lives here

**`pkg/`** - Public library code
- Can be imported by other projects
- Shared utilities

**`migrations/`** - Database migrations
- Sequential numbering: `000001_`, `000002_`, etc.
- `.up.sql` for forward migration
- `.down.sql` for rollback

## API Endpoints (Matching C# Implementation)

### Content Endpoints
```
GET    /content           - List all content
GET    /content/{name}    - Get content by name (TODO: change to ID later)
```

### YouTube Endpoints
```
GET    /youtube/video?videoId={id}  - Fetch video from YouTube API
POST   /youtube/videos              - Import multiple videos
PUT    /youtube/videos              - Update videos (currently mirrors POST)
```

### Perspectives Endpoints
```
GET    /perspectives/{username}  - Get user's perspectives
GET    /perspectives/{id}        - Get single perspective
POST   /perspectives             - Create perspectives (batch)
PUT    /perspectives/{id}        - Update perspective
DELETE /perspectives             - Delete perspectives (batch)
```

### Future Endpoints (Not in C# version yet)
```
# May add later:
GET    /content/{id}      - Get content by ID (cleaner than by name)
PUT    /youtube/video     - Simple single video update
```

## Database Strategy

### Approach: Shared Database
- **Use existing PostgreSQL database** from C# version
- **Same schema** initially (can refine later)
- **Port EF Core migrations** to SQL format
- **Both apps can coexist** during development

### Migration Porting Strategy
1. Convert C# migrations to SQL files
2. Keep same schema, constraints, triggers
3. Support for PostgreSQL-specific features:
   - JSONB columns
   - Array types
   - Custom domains (`valid_integer_range`)
   - Triggers (`update_updated_at`)

### Example Migration (from C#)
```sql
-- migrations/000001_create_content.up.sql
CREATE TABLE IF NOT EXISTS content (
    id SERIAL PRIMARY KEY,
    url VARCHAR UNIQUE,
    length INTEGER,
    length_units VARCHAR,
    response JSONB,
    content_type VARCHAR NOT NULL,
    name VARCHAR NOT NULL UNIQUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_content_url ON content(url);
CREATE INDEX idx_content_name ON content(name);
```

## Configuration Format

### config.json (matches appsettings.json structure)
```json
{
  "server": {
    "port": 8080,
    "host": "localhost"
  },
  "database": {
    "host": "localhost",
    "port": 5432,
    "name": "testdb",
    "user": "testuser",
    "sslmode": "disable"
  },
  "youtube": {
    "api_key": ""
  },
  "logging": {
    "level": "info",
    "format": "json"
  }
}
```

### Environment Variables (for secrets)
```bash
DATABASE_PASSWORD=testpass
YOUTUBE_API_KEY=your_key_here
```

## Development Workflow

### Makefile Commands
```makefile
.PHONY: build run test migrate-up migrate-down swagger

build:
	go build -o bin/perspectize-server cmd/server/main.go

run:
	go run cmd/server/main.go

test:
	go test -v -cover ./...

migrate-up:
	migrate -path migrations -database "postgres://testuser:testpass@localhost:5432/testdb?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgres://testuser:testpass@localhost:5432/testdb?sslmode=disable" down 1

swagger:
	swag init -g cmd/server/main.go -o docs

lint:
	golangci-lint run

fmt:
	go fmt ./...

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down
```

### Docker Compose (for local PostgreSQL)
```yaml
version: '3.8'
services:
  postgres:
    image: postgres:18
    container_name: perspectize-postgres
    environment:
      POSTGRES_DB: testdb
      POSTGRES_USER: testuser
      POSTGRES_PASSWORD: testpass
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```

## Known Improvements from C# Code

These will be implemented in the Go version:

1. **Use ON CONFLICT for upserts** instead of find-then-update
   ```sql
   INSERT INTO content (...) VALUES (...)
   ON CONFLICT (url) DO UPDATE SET ...
   ```

2. **Content lookup by ID** in addition to name
   - Add `GET /content/{id}` endpoint
   - Keep name-based lookup for backwards compatibility

3. **Simplified YouTube PUT**
   - Implement proper single-video update
   - Don't reuse POST logic

4. **Better error handling**
   - Structured error responses
   - Consistent error format across all endpoints

## Dependencies (go.mod)

```go
module github.com/yourorg/perspectize-go

go 1.25  // or latest stable

require (
    github.com/go-chi/chi/v5 v5.0.12
    github.com/jmoiron/sqlx v1.4.0
    github.com/jackc/pgx/v5 v5.6.0
    github.com/golang-migrate/migrate/v4 v4.17.1
    github.com/go-playground/validator/v10 v10.22.0
    github.com/swaggo/swag v1.16.3
    github.com/swaggo/http-swagger v1.3.4
    github.com/stretchr/testify v1.9.0
    github.com/DATA-DOG/go-sqlmock v1.5.2
)
```

## Migration from C# - Phase Approach

### Phase 1: Foundation (Small, Testable Steps)
1. Project setup (go.mod, folder structure)
2. Database connection and migrations
3. Configuration loading
4. Basic HTTP server with health check

### Phase 2: Read Operations (Low Risk)
1. Content repository (GET operations)
2. Content controllers and routes
3. Tests for content retrieval

### Phase 3: YouTube Integration
1. YouTube service (video ID extraction, API calls)
2. YouTube controllers
3. Content creation via YouTube import
4. Tests for YouTube integration

### Phase 4: Perspectives (Full CRUD)
1. Perspective repository (all operations)
2. Perspective service (business logic, validation)
3. Perspective controllers
4. Tests for perspective CRUD

### Phase 5: Polish
1. Swagger documentation
2. Error handling improvements
3. Logging middleware
4. Performance testing

## Learning Go - Helpful Patterns

### Error Handling
```go
// Go style: explicit error checking
content, err := repo.GetContentByName(name)
if err != nil {
    return nil, fmt.Errorf("failed to get content: %w", err)
}
```

### Struct Tags (familiar from C#)
```go
type Content struct {
    ID          int       `json:"id" db:"id"`
    Name        string    `json:"name" db:"name" validate:"required"`
    ContentType string    `json:"content_type" db:"content_type"`
}
```

### Controller Pattern
```go
func (c *ContentController) GetAll(w http.ResponseWriter, r *http.Request) {
    contents, err := c.service.GetAll(r.Context())
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(contents)
}
```

### Dependency Injection (Constructor Pattern)
```go
type PerspectiveService struct {
    repo PerspectiveRepository
}

func NewPerspectiveService(repo PerspectiveRepository) *PerspectiveService {
    return &PerspectiveService{repo: repo}
}
```

## Future Considerations

### Authentication (Deferred)
- **Planned**: Clerk integration
- **Alternative**: JWT tokens with standard library
- **For now**: Skip entirely or use hardcoded user

### Background Jobs (Deferred)
- **If needed later**: `robfig/cron/v3`
- **Use case**: Periodic YouTube metadata refresh
- **For now**: Manual updates via API

### PostgreSQL 18 Features
- Exploring latest features in major version 18
- Can leverage new JSON functions, performance improvements
- Migration tool supports latest PostgreSQL

## Success Criteria

### Functional
- ✓ All C# endpoints working in Go
- ✓ Same database schema and data
- ✓ Swagger documentation available
- ✓ Tests passing for each feature

### Learning Goals
- ✓ Understand Go project structure
- ✓ Learn Go HTTP patterns
- ✓ Database operations in Go
- ✓ Testing in Go

### Quality
- ✓ Small, reviewable commits
- ✓ Each step independently testable
- ✓ Clear error messages
- ✓ Code comments for learning

## References

- [Effective Go](https://go.dev/doc/effective_go)
- [Go Standard Library](https://pkg.go.dev/std)
- [pgx Documentation](https://pkg.go.dev/github.com/jackc/pgx/v5)
- [chi Router](https://go-chi.io/)
- [golang-migrate](https://github.com/golang-migrate/migrate)
- [Swaggo](https://github.com/swaggo/swag)

---

**Document Version**: 1.0
**Created**: October 2025
**Status**: Ready for Review
