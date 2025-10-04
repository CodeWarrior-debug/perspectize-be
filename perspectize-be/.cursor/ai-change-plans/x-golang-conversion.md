# X - Golang Backend Conversion

**Ticket**: X (Golang Conversion)
**Branch**: `x-golang-conversion`
**Architecture Doc**: [.cursor/docs/1-architecture-plan.md](../.cursor/docs/1-architecture-plan.md)

## Overview

Convert the Perspectize backend from C# ASP.NET Core to Go, maintaining all existing functionality while leveraging Go's simplicity and performance. This is a learning-focused migration with small, testable steps.

**Key Principles**:
- Small, incremental changes that can be validated independently
- Each step should be testable before moving to the next
- Focus on learning Go idioms and patterns
- Document as we go with full Swagger/OpenAPI

---

## Phase 1: Foundation & Setup

### 1.1 Project Initialization
- [x] 1.1.1 Create `perspectize-go/` folder in project root
- [x] 1.1.2 Initialize Go module: `go mod init github.com/yourorg/perspectize-go`
- [x] 1.1.3 Create project folder structure:
  ```
  perspectize-go/
  ├── cmd/server/
  ├── internal/
  │   ├── config/
  │   ├── models/
  │   ├── dto/
  │   ├── controllers/
  │   ├── services/
  │   ├── repositories/
  │   ├── middleware/
  │   └── validator/
  ├── pkg/database/
  ├── migrations/
  ├── test/
  ├── config/
  └── docs/
  ```
- [x] 1.1.4 Create basic `README.md` with setup instructions
- [x] 1.1.5 Create `.gitignore` for Go project

### 1.2 Configuration System
- [x] 1.2.1 Create `config/config.json` based on C# appsettings.json
- [x] 1.2.2 Create `config/config.example.json` (with placeholders for secrets)
- [x] 1.2.3 Implement `internal/config/config.go`:
  - Define Config struct matching JSON structure
  - Load from config.json
  - Override with environment variables (DATABASE_PASSWORD, YOUTUBE_API_KEY)
- [x] 1.2.4 Write tests for config loading in `test/config/config_test.go`

### 1.3 Database Connection
- [x] 1.3.1 Add dependencies to go.mod:
  - `github.com/jmoiron/sqlx`
  - `github.com/jackc/pgx/v5`
- [x] 1.3.2 Create `pkg/database/postgres.go`:
  - Database connection function using pgx driver
  - Connection pool configuration
  - Ping/health check function
- [x] 1.3.3 Write connection tests in `test/database/postgres_test.go`
- [x] 1.3.4 Create `docker-compose.yml` for PostgreSQL 18
- [x] 1.3.5 Create `Makefile` with docker-up, docker-down, test commands
  ```yaml
  version: '3.8'
  services:
    postgres:
      image: postgres:18
      container_name: perspectize-postgres-go
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

### 1.4 Database Migrations
- [ ] 1.4.1 Add migration dependency: `github.com/golang-migrate/migrate/v4`
- [ ] 1.4.2 Port C# migration: Initial Content table
  - Create `migrations/000001_create_content.up.sql`
  - Create `migrations/000001_create_content.down.sql`
  - Match exact schema from C# EF Core migration
- [ ] 1.4.3 Port C# migration: Update Response to JSONB
  - Create `migrations/000002_update_response_jsonb.up.sql`
  - Create `migrations/000002_update_response_jsonb.down.sql`
- [ ] 1.4.4 Port C# migration: Update Length to Numeric
  - Create `migrations/000003_update_length_numeric.up.sql`
  - Create `migrations/000003_update_length_numeric.down.sql`
- [ ] 1.4.5 Port C# migration: Add Perspectives, Users, Domain
  - Create `migrations/000004_add_perspectives_users.up.sql`
  - Create `migrations/000004_add_perspectives_users.down.sql`
  - Include custom domain `valid_integer_range`
  - Include `update_updated_at` trigger function
  - Include all constraints and foreign keys
- [ ] 1.4.6 Create `Makefile` with migration commands:
  ```makefile
  migrate-up:
    migrate -path migrations -database "postgres://..." up
  migrate-down:
    migrate -path migrations -database "postgres://..." down 1
  ```
- [ ] 1.4.7 Test migrations: up and down

### 1.5 Basic HTTP Server
- [ ] 1.5.1 Add dependencies:
  - `github.com/go-chi/chi/v5`
- [ ] 1.5.2 Create `cmd/server/main.go`:
  - Load configuration
  - Initialize database connection
  - Setup chi router
  - Add health check endpoint: `GET /health`
  - Start HTTP server
- [ ] 1.5.3 Add logging with `log/slog`:
  - Structured JSON logging
  - Log server startup, shutdown, errors
- [ ] 1.5.4 Create `internal/middleware/logger.go`:
  - HTTP request logging middleware
  - Log method, path, status code, duration
- [ ] 1.5.5 Create `internal/middleware/recovery.go`:
  - Panic recovery middleware
  - Log panics and return 500 error
- [ ] 1.5.6 Update Makefile:
  ```makefile
  build:
    go build -o bin/perspectize-server cmd/server/main.go
  run:
    go run cmd/server/main.go
  test:
    go test -v -cover ./...
  ```
- [ ] 1.5.7 Test: Start server, hit health check endpoint
- [ ] 1.5.8 Test: Verify logging output
- [ ] 1.5.9 Test: Trigger panic and verify recovery

---

## Phase 2: Content Read Operations

### 2.1 Content Models and DTOs
- [ ] 2.1.1 Create `internal/models/content.go`:
  - Port Content model from C# with proper struct tags
  - Use `db:` tags for sqlx mapping
  - Use `json:` tags for JSON serialization
  - Handle nullable fields appropriately
- [ ] 2.1.2 Create `internal/dto/content.go`:
  - ContentResponse DTO (if needed)
  - Match C# DTO structure

### 2.2 Content Repository (Read Only)
- [ ] 2.2.1 Create `internal/repositories/content.go`:
  - Define ContentRepository interface:
    - `GetAll(ctx context.Context) ([]models.Content, error)`
    - `GetByName(ctx context.Context, name string) (*models.Content, error)`
  - Implement repository struct with sqlx.DB
  - Implement GetAll using sqlx.Select
  - Implement GetByName using sqlx.Get
- [ ] 2.2.2 Create `test/repositories/content_test.go`:
  - Test GetAll with mock data (using sqlmock)
  - Test GetByName with existing content
  - Test GetByName with non-existent content (should return nil)
- [ ] 2.2.3 Add dependency: `github.com/DATA-DOG/go-sqlmock`

### 2.3 Content Controller
- [ ] 2.3.1 Create `internal/controllers/content.go`:
  - Define ContentController struct with repository dependency
  - Implement `GetAll(w http.ResponseWriter, r *http.Request)`:
    - Call repository.GetAll
    - Return JSON array of content
    - Handle errors (500 for DB errors)
  - Implement `GetByName(w http.ResponseWriter, r *http.Request)`:
    - Extract name from URL params (chi.URLParam)
    - Validate name is not empty (400 if empty)
    - Call repository.GetByName
    - Return 404 if not found
    - Return JSON content if found
- [ ] 2.3.2 Add routes in `cmd/server/main.go`:
  - `GET /content` → ContentController.GetAll
  - `GET /content/{name}` → ContentController.GetByName
- [ ] 2.3.3 Create `test/controllers/content_test.go`:
  - Test GetAll endpoint
  - Test GetByName with valid name
  - Test GetByName with empty name (400)
  - Test GetByName with non-existent name (404)

### 2.4 Manual Testing & Validation
- [ ] 2.4.1 Start server: `make run`
- [ ] 2.4.2 Test GET /content (should return empty array initially)
- [ ] 2.4.3 Insert test data via C# app or psql
- [ ] 2.4.4 Test GET /content (should return test data)
- [ ] 2.4.5 Test GET /content/{name} with existing content
- [ ] 2.4.6 Test GET /content/nonexistent (should return 404)

---

## Phase 3: YouTube Integration

### 3.1 YouTube Service
- [ ] 3.1.1 Create `internal/services/youtube.go`:
  - Define YouTubeService struct with:
    - HTTP client
    - API key from config
  - Implement `ExtractVideoID(url string) (string, error)`:
    - Port regex logic from C# YouTubeService
    - Handle various YouTube URL formats
    - Return error for invalid URLs
  - Implement `ConvertDurationToSeconds(isoDuration string) (int, error)`:
    - Parse ISO 8601 duration (PT format)
    - Convert to total seconds
    - Return error for invalid format
- [ ] 3.1.2 Create `test/services/youtube_test.go`:
  - Test ExtractVideoID with various URL formats:
    - `https://www.youtube.com/watch?v=dQw4w9WgXcQ`
    - `https://youtu.be/dQw4w9WgXcQ`
    - `https://www.youtube.com/embed/dQw4w9WgXcQ`
  - Test ExtractVideoID with invalid URLs (should error)
  - Test ConvertDurationToSeconds with various ISO durations:
    - `PT1M30S` → 90
    - `PT1H2M3S` → 3723
    - `PT30S` → 30
  - Test ConvertDurationToSeconds with invalid format

### 3.2 YouTube DTOs
- [ ] 3.2.1 Create `internal/dto/youtube.go`:
  - VideosRequest struct (array of video URLs)
  - VideoResponse struct (status, videoId, name, url)
  - Match C# DTO structure

### 3.3 Content Repository (Write Operations)
- [ ] 3.3.1 Update `internal/repositories/content.go`:
  - Add to ContentRepository interface:
    - `Upsert(ctx context.Context, content *models.Content) error`
  - Implement Upsert using PostgreSQL ON CONFLICT:
    ```sql
    INSERT INTO content (...) VALUES (...)
    ON CONFLICT (url) DO UPDATE SET
      length = EXCLUDED.length,
      length_units = EXCLUDED.length_units,
      response = EXCLUDED.response,
      name = EXCLUDED.name,
      updated_at = NOW()
    ```
  - This improves on C# find-then-update approach
- [ ] 3.3.2 Update `test/repositories/content_test.go`:
  - Test Upsert with new content (INSERT)
  - Test Upsert with existing content (UPDATE)
  - Test Upsert with conflict on URL

### 3.4 YouTube Controller
- [ ] 3.4.1 Create `internal/controllers/youtube.go`:
  - Define YouTubeController struct with:
    - YouTubeService
    - ContentRepository
  - Implement `GetVideo(w http.ResponseWriter, r *http.Request)`:
    - Get videoId query param
    - Call YouTube API
    - Return raw JSON response
    - Handle errors (400, 500, YouTube API errors)
  - Implement `PostVideos(w http.ResponseWriter, r *http.Request)`:
    - Parse VideosRequest from body
    - Validate at least one URL
    - For each URL:
      - Extract video ID
      - Fetch from YouTube API
      - Parse response (title, duration)
      - Upsert to database via repository
      - Track result (created/updated/error)
    - Return array of VideoResponse
- [ ] 3.4.2 Add routes in `cmd/server/main.go`:
  - `GET /youtube/video` → YouTubeController.GetVideo
  - `POST /youtube/videos` → YouTubeController.PostVideos
  - `PUT /youtube/videos` → YouTubeController.PostVideos (reuse for now, like C#)
- [ ] 3.4.3 Create `test/controllers/youtube_test.go`:
  - Test GetVideo with valid videoId
  - Test GetVideo without videoId (400)
  - Test PostVideos with valid URLs
  - Test PostVideos with empty array (400)
  - Test PostVideos with invalid URL
  - Mock YouTube API responses

### 3.5 Manual Testing & Validation
- [ ] 3.5.1 Set YOUTUBE_API_KEY environment variable
- [ ] 3.5.2 Test GET /youtube/video?videoId=dQw4w9WgXcQ
- [ ] 3.5.3 Test POST /youtube/videos with array of URLs
- [ ] 3.5.4 Verify content created in database
- [ ] 3.5.5 Test POST /youtube/videos with same URLs (should update)
- [ ] 3.5.6 Verify updated_at timestamp changed

---

## Phase 4: Perspectives System

### 4.1 Perspective Models and DTOs
- [ ] 4.1.1 Create `internal/models/perspective.go`:
  - Port Perspective model from C#
  - Handle nullable integer pointers for ratings
  - Handle array types (parts, labels)
  - Handle JSONB arrays (categorized_ratings)
  - Use proper struct tags
- [ ] 4.1.2 Create `internal/models/user.go`:
  - Basic User model (id, username, email)
  - No auth functionality yet
- [ ] 4.1.3 Create `internal/dto/perspective.go`:
  - CreatePerspectiveRequest
  - UpdatePerspectiveRequest
  - PerspectiveResponse (with joined Content data)
  - Add validation tags using `validate:` from validator/v10
  - Port validation logic from C# (required, max lengths, ranges)

### 4.2 Validation Setup
- [ ] 4.2.1 Add dependency: `github.com/go-playground/validator/v10`
- [ ] 4.2.2 Create `internal/validator/validator.go`:
  - Initialize validator instance
  - Helper function to validate structs
  - Format validation errors into user-friendly messages
- [ ] 4.2.3 Create `test/validator/validator_test.go`:
  - Test validation on CreatePerspectiveRequest
  - Test required fields
  - Test max length constraints
  - Test range constraints (0-10000)

### 4.3 Perspective Repository
- [ ] 4.3.1 Create `internal/repositories/perspective.go`:
  - Define PerspectiveRepository interface:
    - `GetByUsername(ctx, username) ([]dto.PerspectiveResponse, error)`
    - `GetByID(ctx, id) (*dto.PerspectiveResponse, error)`
    - `Create(ctx, requests) (int, error)` - batch create
    - `Update(ctx, id, request) (int, error)` - returns affected rows
    - `Delete(ctx, ids) (int, error)` - batch delete
    - `Exists(ctx, id) (bool, error)`
  - Implement repository with sqlx.DB
  - Port SQL queries from C# PerspectiveRepository:
    - Use same JOIN logic for GetByUsername and GetByID
    - Handle PostgreSQL array parameters (ANY clause)
    - Handle JSONB serialization for categorized_ratings
    - Dynamic UPDATE query building for partial updates
- [ ] 4.3.2 Create `test/repositories/perspective_test.go`:
  - Test GetByUsername
  - Test GetByID
  - Test Create with valid data
  - Test Create with duplicate claim (should error)
  - Test Update with valid data
  - Test Update with non-existent ID
  - Test Delete batch operation
  - Test Exists

### 4.4 Perspective Service
- [ ] 4.4.1 Create `internal/services/perspective.go`:
  - Define PerspectiveService interface
  - Implement service with PerspectiveRepository dependency
  - Port business logic from C# PerspectiveService
  - Implement error handling:
    - Catch unique constraint violations (perspectives_unique_user_claims)
    - Catch domain constraint violations (valid_integer_range)
    - Translate DB errors to business errors
- [ ] 4.4.2 Create `test/services/perspective_test.go`:
  - Test service methods with mocked repository
  - Test error translation logic
  - Test validation integration

### 4.5 Perspective Controller
- [ ] 4.5.1 Create `internal/controllers/perspective.go`:
  - Define PerspectiveController struct with:
    - PerspectiveService
    - Validator
  - Implement `GetByUsername(w, r)`:
    - Extract username from URL params
    - Call service.GetByUsername
    - Return 200 with array (empty if no perspectives)
  - Implement `GetByID(w, r)`:
    - Extract id from URL params
    - Parse id as integer
    - Call service.GetByID
    - Return 404 if not found
    - Return 200 with perspective
  - Implement `Create(w, r)`:
    - Parse []CreatePerspectiveRequest from body
    - Validate each request
    - Return 400 if validation fails
    - Return 400 if empty array
    - Call service.Create
    - Handle conflicts (409)
    - Handle validation errors (400)
    - Return 201 with created count
  - Implement `Update(w, r)`:
    - Extract id from URL params
    - Parse UpdatePerspectiveRequest from body
    - Validate request
    - Call service.Update
    - Return 404 if not found
    - Handle conflicts (409)
    - Return 200 with updated perspective
  - Implement `Delete(w, r)`:
    - Parse []int from body
    - Return 400 if empty array
    - Call service.Delete
    - Return 204 with deleted count
- [ ] 4.5.2 Add routes in `cmd/server/main.go`:
  - `GET /perspectives/{username}` → PerspectiveController.GetByUsername
  - `GET /perspectives/{id}` → PerspectiveController.GetByID (id must be int)
  - `POST /perspectives` → PerspectiveController.Create
  - `PUT /perspectives/{id}` → PerspectiveController.Update
  - `DELETE /perspectives` → PerspectiveController.Delete
- [ ] 4.5.3 Create `test/controllers/perspective_test.go`:
  - Test all endpoints
  - Test validation errors
  - Test conflict errors
  - Test not found scenarios
  - Test batch operations

### 4.6 Manual Testing & Validation
- [ ] 4.6.1 Create test user in database (via psql)
- [ ] 4.6.2 Create test content in database
- [ ] 4.6.3 Test POST /perspectives with valid data
- [ ] 4.6.4 Test GET /perspectives/{username}
- [ ] 4.6.5 Test GET /perspectives/{id}
- [ ] 4.6.6 Test PUT /perspectives/{id} with partial update
- [ ] 4.6.7 Test POST /perspectives with duplicate claim (should 409)
- [ ] 4.6.8 Test POST /perspectives with invalid range (should 400)
- [ ] 4.6.9 Test DELETE /perspectives with array of IDs
- [ ] 4.6.10 Verify all database constraints work (unique, foreign key, domain)

---

## Phase 5: Documentation & Polish

### 5.1 Swagger Documentation
- [ ] 5.1.1 Add dependencies:
  - `github.com/swaggo/swag`
  - `github.com/swaggo/http-swagger`
- [ ] 5.1.2 Add Swagger annotations to all controller methods:
  - @Summary
  - @Description
  - @Tags
  - @Accept
  - @Produce
  - @Param
  - @Success
  - @Failure
  - @Router
- [ ] 5.1.3 Add general API info in `cmd/server/main.go`:
  - @title
  - @version
  - @description
  - @host
  - @BasePath
- [ ] 5.1.4 Update Makefile:
  ```makefile
  swagger:
    swag init -g cmd/server/main.go -o docs
  ```
- [ ] 5.1.5 Generate Swagger docs: `make swagger`
- [ ] 5.1.6 Add Swagger UI route: `GET /swagger/*`
- [ ] 5.1.7 Test: Browse to http://localhost:8080/swagger/index.html
- [ ] 5.1.8 Verify all endpoints documented
- [ ] 5.1.9 Test all endpoints via Swagger UI

### 5.2 Error Handling Improvements
- [ ] 5.2.1 Create `internal/dto/error.go`:
  - Standard error response structure:
    ```go
    type ErrorResponse struct {
        Message string `json:"message"`
        Code    string `json:"code,omitempty"`
        Details []string `json:"details,omitempty"`
    }
    ```
- [ ] 5.2.2 Update all controllers to use ErrorResponse
- [ ] 5.2.3 Create helper functions:
  - `respondWithError(w, status, message)`
  - `respondWithValidationError(w, errors)`
  - `respondWithJSON(w, status, data)`
- [ ] 5.2.4 Test error responses have consistent format

### 5.3 Logging Enhancements
- [ ] 5.3.1 Add request ID middleware:
  - Generate unique ID per request
  - Add to context
  - Include in all logs
  - Return in response header
- [ ] 5.3.2 Add structured logging to all controllers:
  - Log request start
  - Log request completion with status
  - Log errors with context
- [ ] 5.3.3 Test: Verify logs are JSON formatted
- [ ] 5.3.4 Test: Verify request IDs in logs

### 5.4 Code Quality & Linting
- [ ] 5.4.1 Add golangci-lint configuration (`.golangci.yml`)
- [ ] 5.4.2 Run `golangci-lint run` and fix issues
- [ ] 5.4.3 Run `go fmt ./...` on all code
- [ ] 5.4.4 Run `go vet ./...` and fix issues
- [ ] 5.4.5 Add code comments for exported functions
- [ ] 5.4.6 Add package documentation comments

### 5.5 README & Documentation
- [ ] 5.5.1 Update `perspectize-go/README.md`:
  - Project description
  - Prerequisites (Go version, PostgreSQL)
  - Setup instructions
  - Configuration guide
  - Running the application
  - Running tests
  - API documentation link
  - Makefile commands reference
- [ ] 5.5.2 Create `perspectize-go/DEVELOPMENT.md`:
  - Development workflow
  - Project structure explanation
  - Testing guide
  - Database migrations guide
  - Common issues and solutions

### 5.6 Final Testing
- [ ] 5.6.1 Run all tests: `make test`
- [ ] 5.6.2 Verify test coverage: `go test -cover ./...`
- [ ] 5.6.3 Run integration tests against real database
- [ ] 5.6.4 Test complete workflow:
  - Import YouTube videos
  - Create content
  - Create user (manual)
  - Create perspectives
  - Query perspectives by username
  - Update perspective
  - Delete perspective
- [ ] 5.6.5 Performance comparison with C# version:
  - Response times
  - Memory usage
  - CPU usage
- [ ] 5.6.6 Verify both C# and Go can use same database

---

## Phase 6: Optional Enhancements

### 6.1 Improvements Over C# Version
- [ ] 6.1.1 Add `GET /content/{id}` endpoint (lookup by ID, not name)
- [ ] 6.1.2 Implement proper `PUT /youtube/video` (single video update)
- [ ] 6.1.3 Add pagination to `GET /content` (query params: limit, offset)
- [ ] 6.1.4 Add pagination to `GET /perspectives/{username}`
- [ ] 6.1.5 Add filtering to perspectives (by privacy, category, etc.)
- [ ] 6.1.6 Add sorting options

### 6.2 Database Seed Data
- [ ] 6.2.1 Create `migrations/seed/` directory
- [ ] 6.2.2 Port SeedData.cs logic to SQL
- [ ] 6.2.3 Create seed script for development
- [ ] 6.2.4 Add `make seed` command

### 6.3 Docker Support
- [ ] 6.3.1 Create `Dockerfile` for Go application
- [ ] 6.3.2 Create multi-stage build (build + runtime)
- [ ] 6.3.3 Update docker-compose.yml to include Go service
- [ ] 6.3.4 Test: `docker-compose up` runs both PostgreSQL and Go service

---

## Success Criteria

### Functional Requirements
- ✅ All C# endpoints migrated and working
- ✅ Same database schema (can coexist with C#)
- ✅ All CRUD operations functional
- ✅ YouTube integration working
- ✅ Swagger documentation complete

### Quality Requirements
- ✅ Tests pass for all features
- ✅ Validation works for all inputs
- ✅ Error handling is consistent
- ✅ Logging is structured and helpful
- ✅ Code is formatted and linted

### Learning Goals
- ✅ Understand Go project structure
- ✅ Learn Go HTTP server patterns
- ✅ Understand sqlx for database operations
- ✅ Learn Go testing practices
- ✅ Understand Go error handling

### Documentation
- ✅ README with setup instructions
- ✅ API fully documented in Swagger
- ✅ Code comments for exported functions
- ✅ Development guide

---

## Notes

- **Small Steps**: Each checkbox should be completable in < 30 minutes
- **Test As You Go**: Don't move to next step until current step is tested
- **Ask Questions**: If unclear about Go patterns, ask for clarification
- **Commit Often**: Commit after each major checkpoint
- **Compare with C#**: Reference C# code frequently to ensure parity

## Key Files Reference

### C# Files (Reference)
- `perspectize-be/Program.cs` - Startup and DI
- `perspectize-be/Controllers/ContentController.cs`
- `perspectize-be/Controllers/PerspectivesController.cs`
- `perspectize-be/Controllers/YTController.cs`
- `perspectize-be/Services/PerspectiveService.cs`
- `perspectize-be/Services/YouTubeService.cs`
- `perspectize-be/Repositories/PerspectiveRepository.cs`
- `perspectize-be/Models/Content.cs`
- `perspectize-be/Models/Perspective.cs`
- `perspectize-be/Migrations/*.cs`

### Go Files (To Create)
- `perspectize-go/cmd/server/main.go`
- `perspectize-go/internal/controllers/*.go`
- `perspectize-go/internal/services/*.go`
- `perspectize-go/internal/repositories/*.go`
- `perspectize-go/internal/models/*.go`
- `perspectize-go/migrations/*.sql`

---

**Last Updated**: October 2025
**Status**: Ready to Begin
