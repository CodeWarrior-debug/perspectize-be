# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**Perspectize Go Backend** - GraphQL API for storing, refining, and sharing perspectives on content (initially YouTube videos). Built with Go 1.25+, PostgreSQL 18, and GraphQL using gqlgen.

**Important:** This repository is migrating from C# ASP.NET Core to Go. The C# implementation in [perspectize-be/](perspectize-be/) is legacy code and should not be modified. All new development happens in [perspectize-go/](perspectize-go/).

## Architecture

This project follows **Hexagonal Architecture** (Ports and Adapters pattern):

```
perspectize-go/
├── cmd/server/              # Application entry point
├── internal/
│   ├── core/                # Domain layer (business logic)
│   │   ├── domain/          # Domain models and entities
│   │   ├── ports/           # Port interfaces (contracts)
│   │   │   ├── repositories/   # Repository interfaces
│   │   │   └── services/       # Service interfaces
│   │   └── services/        # Domain services (business logic)
│   ├── adapters/            # Adapters layer (infrastructure)
│   │   ├── graphql/         # GraphQL resolvers (primary adapter)
│   │   ├── repositories/    # Database implementations (secondary adapter)
│   │   └── youtube/         # External API clients (secondary adapter)
│   ├── config/              # Configuration loading
│   └── middleware/          # HTTP middleware
├── pkg/                     # Shared packages
│   ├── database/            # Database connection utilities
│   └── graphql/             # Custom GraphQL scalars (IntID)
└── migrations/              # SQL migration files
```

### Key Architectural Principles

**Hexagon Core (Domain Layer):**
- `internal/core/domain/` - Pure domain models, no external dependencies
- `internal/core/ports/` - Interfaces defining contracts (repositories, services)
- `internal/core/services/` - Business logic, depends only on ports

**Adapters (Infrastructure):**
- `internal/adapters/graphql/` - PRIMARY adapter: GraphQL API (gqlgen)
- `internal/adapters/repositories/` - SECONDARY adapter: PostgreSQL (sqlx + pgx)
- `internal/adapters/youtube/` - SECONDARY adapter: YouTube Data API

**Dependency Rule:** Dependencies point inward. Domain never depends on adapters. Adapters depend on domain ports.

## Development Commands

### Initial Setup

```bash
cd perspectize-go

# Install Go dependencies
go mod download

# Start PostgreSQL (Docker)
make docker-up

# Run database migrations
make migrate-up

# Copy and configure environment
cp .env.example .env
# Edit .env with your DATABASE_URL and YOUTUBE_API_KEY
```

### Daily Development

```bash
# Run the server (loads .env automatically)
make run
# Server runs on http://localhost:8080

# Run with hot-reload (uses air, auto-restarts on file changes)
make dev

# Run all tests
make test

# Run tests with coverage
make test-coverage
# Generates coverage.out and coverage.html
# Open the HTML report in browser:
open coverage.html   # macOS

# Run all tests (integration tests self-skip when DB is unavailable)
go test ./...

# Run single test
go test -v -run TestFunctionName ./path/to/package

# Format code
make fmt

# Lint code
make lint

# Generate GraphQL code (after schema changes)
make graphql-gen
```

### Database Migrations

```bash
# Run all pending migrations
make migrate-up

# Rollback last migration
make migrate-down

# Create new migration
make migrate-create
# Then enter migration name when prompted

# Check current migration version
make migrate-version

# Force migration to specific version (recovery)
make migrate-force
```

### GraphQL Development

```bash
# After modifying schema.graphql:
make graphql-gen

# View GraphQL Playground
# Start server, then browse to http://localhost:8080/graphql
```

### Docker

```bash
# Start PostgreSQL
make docker-up

# Stop and remove containers
make docker-down

# View PostgreSQL logs
make docker-logs
```

### GitHub & Repository Management

**Always use the `gh` CLI** for all GitHub operations. Do not use MCP plugins or other tools for GitHub interactions.

```bash
# Pull requests
gh pr create --title "Title" --body "Description"
gh pr list
gh pr view 123
gh pr merge 123

# Edit PR (use REST API - gh pr edit may fail with Projects Classic deprecation error)
gh api repos/{owner}/{repo}/pulls/123 -X PATCH -f body="New description"

# Issues (use API - gh issue view fails with Projects Classic deprecation error)
gh issue create --title "Title" --body "Description"
gh issue list
gh api repos/{owner}/{repo}/issues/123 --jq '.title, .html_url'

# Repository info
gh repo view

# API access (for anything not covered by commands)
gh api repos/{owner}/{repo}/pulls/123/comments
```

### GitHub Projects (v2)

The `gh project` command may not be available. Use GraphQL API directly:

```bash
# List user's projects
gh api graphql -f query='{ viewer { projectsV2(first: 10) { nodes { id title number } } } }'

# Get issue node ID
gh api graphql -f query='query { repository(owner: "OWNER", name: "REPO") { issue(number: 35) { id } } }'

# Add issue to project
gh api graphql -f query='mutation { addProjectV2ItemById(input: { projectId: "PROJECT_ID", contentId: "ISSUE_NODE_ID" }) { item { id } } }'
```

**Token scopes:** If you get "INSUFFICIENT_SCOPES" errors for project operations, refresh auth:
```bash
gh auth refresh -h github.com -s read:project -s project
```
This opens a browser flow to authorize additional scopes.

### Branch Naming Convention

**Always create branches from an updated `main` branch.**

```bash
# Before creating a new branch
git checkout main
git pull origin main
git checkout -b <branch-name>
```

**Branch name format:** `type/initiativePrefix-issueNumber-description-in-kebab-case`

| Component | Values |
|-----------|--------|
| **type** | `feature`, `bugfix`, `chore` |
| **initiativePrefix** | `INI` (Initialization Phase) |
| **issueNumber** | GitHub issue number (e.g., `16`) |
| **description** | Brief kebab-case description |

**Examples:**
- `feature/INI-16-youtube-post-graphql`
- `bugfix/INI-23-fix-auth-middleware`
- `chore/INI-8-update-dependencies`

### GitHub Issues with GSD Plans

When creating issues that correspond to GSD plans, include:

1. **GSD Plan Reference** section with path: `.planning/phases/{phase}/{plan}-PLAN.md`
2. **Acceptance criteria** matching the plan's `must_haves.truths`
3. **Dependencies** section if the plan has `depends_on`

This links GitHub tracking to the detailed execution plans in `.planning/`.

## Configuration

Configuration is loaded from two sources (in order of precedence):

1. **Environment variables** (highest priority) - secrets like `DATABASE_URL`, `YOUTUBE_API_KEY`
2. **config/config.json** - non-secret configuration

### Environment Variables

Required:
- `DATABASE_URL` - Full PostgreSQL connection string (e.g., `postgres://user:pass@host:5432/dbname?sslmode=disable`)

Optional:
- `YOUTUBE_API_KEY` - YouTube Data API v3 key (for content fetching)
- `DATABASE_PASSWORD` - Alternative to full DATABASE_URL

See [.env.example](perspectize-go/.env.example) for details.

### Local Development Setup

```bash
# Use Docker Compose database
export DATABASE_URL="postgres://testuser:testpass@localhost:5432/testdb?sslmode=disable"

# Or override just password (uses config.json for other settings)
export DATABASE_PASSWORD="testpass"
```

### Production Setup (Sevalla/Fly.io)

Use `DATABASE_URL` with external endpoint from hosting provider. Note: Sevalla connections may require `?sslmode=disable` and may succeed on second attempt.

## Technology Stack

- **Language:** Go 1.25+
- **GraphQL:** gqlgen (code generation, schema-first)
- **Database:** PostgreSQL 18 with sqlx + pgx/v5 driver
- **Migrations:** golang-migrate
- **Validation:** go-playground/validator
- **Testing:** testing + testify + sqlmock
- **Logging:** log/slog (structured logging)
- **Environment:** godotenv (.env file loading)

## Agent Delegation Strategy

| Task Type | Model | Subagent | Rationale |
|-----------|-------|----------|-----------|
| Architecture decisions | Opus | - | Complex multi-file reasoning |
| Go implementation | Sonnet | `go-backend` | Balanced quality/cost |
| GraphQL schema design | Sonnet | `graphql-designer` | Schema patterns |
| Database migrations | Sonnet | `db-migration` | SQL generation |
| Code review | Haiku | `code-reviewer` | Fast pattern matching |
| Test generation | Haiku | `test-writer` | Boilerplate generation |

## Database

PostgreSQL 18 with advanced features:
- **JSONB columns** - Structured data (e.g., YouTube API responses)
- **Array types** - Collections (tags, categories)
- **Custom domains** - Type validation at database level
- **Triggers** - Automatic timestamp updates

### Migration Files

Located in [migrations/](perspectize-go/migrations/):
- `000001_*.up.sql` - Apply migration
- `000001_*.down.sql` - Rollback migration

Naming: `{sequence}_{description}.{up|down}.sql`

## GraphQL Schema

GraphQL schema is defined in `schema.graphql` (schema-first approach). After modifying the schema:

1. Run `make graphql-gen` to regenerate Go types and resolvers
2. Implement resolver logic in `internal/adapters/graphql/resolvers/`
3. Wire resolvers to domain services (follow hexagonal architecture)

## Testing Strategy

### Unit Tests
- Test domain logic in isolation
- Mock external dependencies (repositories, external APIs)
- Fast, no database required
- Run with: `make test` or `go test ./...`

### Integration Tests
- Test adapters against real database
- Use `testcontainers` or Docker Compose for PostgreSQL
- Guarded with `t.Skip()` when prerequisites (e.g., database) are unavailable
- Run with: `make test` (skipped tests are reported automatically)

### Environment Isolation
- Tests that load config must clear env vars leaked by the Makefile (`DATABASE_URL`, `DATABASE_PASSWORD`, `YOUTUBE_API_KEY`)
- Use `t.Setenv("KEY", "")` to clear vars — it auto-restores on test cleanup
- See `clearConfigEnvVars` helper in `test/config/config_test.go` for the pattern

## Code Style

- Follow [Effective Go](https://go.dev/doc/effective_go)
- Use `gofmt` for formatting (enforced by `make fmt`)
- Explicit error handling (no exceptions)
- Structured logging with `slog`
- Dependency injection via interfaces (ports)

## Hexagonal Architecture Guidelines

When implementing features:

1. **Start with domain** - Define models in `core/domain/`
2. **Define ports** - Create interfaces in `core/ports/`
3. **Implement business logic** - Write services in `core/services/`
4. **Add adapters** - Implement infrastructure in `adapters/`
5. **Wire dependencies** - Connect adapters to core in `cmd/server/main.go`

### Domain Layer (`core/domain/`)

The domain layer contains pure Go structs with **no external dependencies** - no database tags, no framework imports, no HTTP/GraphQL code. You should be able to copy domain files to another project and compile them with only the standard library.

**Core entities for this project:**
- `Content` - Media that users create perspectives on (YouTube videos, articles)
- `Perspective` - A user's viewpoint/rating on content (claim, quality, agreement, etc.)

**What belongs in domain:**

| Include | Do NOT Include |
|---------|----------------|
| Business entities (structs) | Database tags (`db:"column"`) |
| Constants/enums (`ContentType`, `Privacy`) | SQL queries |
| Domain errors (`ErrNotFound`, `ErrInvalidRating`) | HTTP/GraphQL code |
| Validation methods | External API calls |

**Optional fields pattern:** Use pointers for nullable/optional fields:

```go
type Perspective struct {
    Claim   string  // Required - always has a value
    Quality *int    // Optional - nil means "not provided"
}

// Check if optional field is set
if p.Quality != nil {
    fmt.Println(*p.Quality)  // Dereference to get value
}

// Set an optional field
quality := 85
p.Quality = &quality
```

**Example flow:**
```
GraphQL Request
  → GraphQL Resolver (adapter)
  → Domain Service (core, uses port interfaces)
  → Repository Interface (port)
  → PostgreSQL Repository (adapter)
```

## Workflow Integration

This project uses **GSD workflow** for planning and execution. See `.planning/` for:
- `PROJECT.md` - Project definition and requirements
- `ROADMAP.md` - Phase-based milestone planning
- `STATE.md` - Current position and accumulated context
- `phases/` - Detailed execution plans

## Self-Verification Workflow

Before marking any work complete, run interactive verification:

### 1. Start Services
```bash
# Terminal 1: Backend
cd perspectize-go && make run

# Terminal 2: Frontend
cd perspectize-fe && pnpm run dev
```

### 2. Verify Backend
```bash
curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -d '{"query": "{ __typename }"}'
# Expect: {"data":{"__typename":"Query"}}
```

### 3. Verify Frontend (Chrome DevTools MCP)

| Step | MCP Tool | Purpose |
|------|----------|---------|
| Navigate | `mcp__chrome-devtools__navigate_page` | Load frontend URL |
| Screenshot | `mcp__chrome-devtools__take_screenshot` | Visual verification |
| Snapshot | `mcp__chrome-devtools__take_snapshot` | DOM/component structure |
| Resize | `mcp__chrome-devtools__resize_page` | Responsive check (375px, 768px, 1024px) |
| Console | `mcp__chrome-devtools__list_console_messages` | Check for JS errors |
| Interact | `mcp__chrome-devtools__click` | Test buttons, toasts, navigation |

### 4. GSD Plan Verification

For each plan's `must_haves`:

| Check | Command |
|-------|---------|
| `truths` | Run actual command, verify output |
| `artifacts.path` | `test -f {path} && echo "exists"` |
| `artifacts.contains` | `grep -q "{pattern}" {path}` |
| `artifacts.min_lines` | `wc -l < {path}` ≥ N |
| `key_links.pattern` | `grep -q "{pattern}" {from}` |

### 5. Evidence Capture

Before creating PR:
- Screenshot at mobile (375px), tablet (768px), desktop (1024px+)
- Console output showing no errors
- Verification commands output

## Legacy C# Code

The `perspectize-be/` directory contains legacy C# ASP.NET Core code. **Do not modify, except to delete.** All development is in `perspectize-go/`.

## Common Patterns

### Adding a New Feature

1. Define domain model: `internal/core/domain/feature.go`
2. Define repository interface: `internal/core/ports/repositories/feature_repository.go`
3. Implement business logic: `internal/core/services/feature_service.go`
4. Implement repository: `internal/adapters/repositories/postgres/feature_repository.go`
5. Update GraphQL schema: `schema.graphql`
6. Generate GraphQL code: `make graphql-gen`
7. Implement resolver: `internal/adapters/graphql/resolvers/feature_resolver.go`
8. Wire in main: `cmd/server/main.go`
9. Write tests: `test/services/feature_service_test.go`, `test/repositories/feature_repository_test.go`

### Error Handling

```go
// Domain errors (core/domain/errors.go)
var ErrNotFound = errors.New("resource not found")

// Services return domain errors
func (s *Service) GetById(id int) (*Model, error) {
    result, err := s.repo.FindById(id)
    if err != nil {
        return nil, fmt.Errorf("failed to get resource: %w", err)
    }
    return result, nil
}

// Resolvers translate to GraphQL errors
func (r *Resolver) GetById(ctx context.Context, id int) (*Model, error) {
    result, err := r.service.GetById(id)
    if errors.Is(err, domain.ErrNotFound) {
        return nil, fmt.Errorf("resource not found")
    }
    return result, err
}
```

### Database Queries

```go
// Use sqlx for queries
var content Content
err := db.Get(&content, "SELECT * FROM content WHERE id = $1", id)

// Use transactions for multi-step operations
tx, err := db.Beginx()
defer tx.Rollback() // Safe to call after commit

// ... perform operations
if err := tx.Commit(); err != nil {
    return err
}
```

## Patterns & Gotchas

### GraphQL Schema Defaults
When a GraphQL field has a default value (e.g., `first: Int = 10`), gqlgen passes the default to the resolver as a non-nil pointer, not `nil`. Tests should expect the default value, not nil.

### JSON Scalar Type
For exposing JSONB data via GraphQL, use gqlgen's built-in `graphql.Map` scalar (configured in `gqlgen.yml` as `JSON`). This avoids string serialization overhead compared to exposing as `String`.

### Cursor-Based Pagination
- Cursors are opaque base64-encoded strings (format: `cursor:<id>`)
- Use keyset pagination in SQL for performance (not OFFSET)
- Fetch `limit+1` rows to determine `hasNextPage` without extra query
- Whitelist sort columns to prevent SQL injection

### GraphQL Enum & ID Handling (REQUIRED)

**Always use gqlgen model binding** to eliminate manual enum mapping code. Never write switch statements to convert between GraphQL and domain enums.

**For enums (SortOrder, Privacy, ContentType, etc.):**
1. Define domain enums with UPPERCASE values to match GraphQL conventions:
   ```go
   // internal/core/domain/pagination.go
   type SortOrder string
   const (
       SortOrderASC  SortOrder = "ASC"
       SortOrderDESC SortOrder = "DESC"
   )
   ```

2. Bind in `gqlgen.yml`:
   ```yaml
   models:
     SortOrder:
       model:
         - github.com/yourorg/perspectize-go/internal/core/domain.SortOrder
   ```

3. For DB-stored enums (Privacy, ContentType, ReviewStatus), add repository converters:
   ```go
   func privacyToDBValue(p domain.Privacy) string {
       return strings.ToLower(string(p))
   }
   func privacyFromDBValue(s string) domain.Privacy {
       return domain.Privacy(strings.ToUpper(s))
   }
   ```

**For ID fields in filters/inputs:**
Use the `IntID` custom scalar (`pkg/graphql/intid.go`) instead of `ID` with manual `strconv.Atoi`:
```graphql
# schema.graphql
scalar IntID

input PerspectiveFilter {
  userID: IntID      # Auto-converts string to int
  contentID: IntID
}
```

**What this eliminates:**
- ❌ No switch statements mapping GraphQL enums to domain enums
- ❌ No `strconv.Atoi` calls for ID filters
- ❌ No duplicate enum definitions in model and domain

**When adding new enums:**
1. Add UPPERCASE constants in `internal/core/domain/`
2. Add binding in `gqlgen.yml`
3. Add DB converter if stored in PostgreSQL
4. Run `make graphql-gen`

### AG Grid Svelte 5 Setup (CRITICAL)

The `ag-grid-svelte5` wrapper bundles AG Grid v32.x internally. **Do NOT install `ag-grid-community` separately** — it causes version conflicts.

**Correct setup:**
```bash
pnpm add ag-grid-svelte5 @ag-grid-community/core@32.2.1 @ag-grid-community/client-side-row-model@32.2.1 @ag-grid-community/theming@32.2.0
```

**Correct imports:**
```svelte
<script lang="ts">
  import AgGridSvelte5Component from 'ag-grid-svelte5';
  import { ClientSideRowModelModule } from '@ag-grid-community/client-side-row-model';
  import { themeQuartz } from '@ag-grid-community/theming';
  import type { GridOptions } from '@ag-grid-community/core';

  const modules = [ClientSideRowModelModule];
  const theme = themeQuartz.withParams({ fontFamily: 'Inter, sans-serif' });
</script>

<AgGridSvelte5Component {gridOptions} {rowData} {theme} {modules} />
```

**Do NOT:**
- ❌ Import from `ag-grid-community` (use `@ag-grid-community/*` packages)
- ❌ Import AG Grid CSS in app.css (use `themeQuartz.withParams()`)
- ❌ Use `AgGridSvelte` (use `AgGridSvelte5Component`)

### GSD Workflow Branching

By default, GSD executes all plans on the current branch. For stacked PRs:

1. **Configure branching in `.planning/config.json`:**
   ```json
   {
     "branching_strategy": "phase",
     "phase_branch_template": "feature/{issue}-plan-{phase}-{slug}"
   }
   ```

2. **Or create branches manually after execution:**
   ```bash
   # Create branch at each plan's completion commit
   git branch feature/plan-01-01 <commit-hash>
   git branch feature/plan-01-02 <commit-hash>
   ```

3. **Create stacked PRs:**
   - PR 1: plan-01-01 → main
   - PR 2: plan-01-02 → plan-01-01
   - PR 3: plan-01-03 → plan-01-02
   - Merge sequentially

## Code Search with qmd

This project has qmd indexing enabled. **Prefer qmd over Read/Glob for exploration.**

| Task | Tool | Example |
|------|------|---------|
| Quick keyword lookup | `qmd_search` | Find files mentioning "GraphQL" |
| Semantic/concept search | `qmd_vsearch` | Find "authentication patterns" |
| Complex questions | `qmd_query` | "How does pagination work?" |
| Get specific file | `qmd_get` | Retrieve by path after search |
| Batch retrieve | `qmd_multi_get` | Get multiple files by glob |

**Workflow:**
1. Use `qmd_search` or `qmd_query` first for exploration
2. Use `qmd_get` to retrieve specific files from search results
3. Fall back to `Read`/`Glob` only if qmd doesn't return enough context

**Re-index after major changes:**
```bash
qmd update  # Re-index modified files
qmd embed   # Update embeddings (run periodically)
```

**GSD Workflow Integration:**

| GSD Phase | qmd Usage |
|-----------|-----------|
| `/gsd:new-project` | Use `qmd_query` to understand existing patterns before defining project |
| `/gsd:map-codebase` | Use `qmd_search` to find all files matching architectural layers |
| `/gsd:plan-phase` | Use `qmd_vsearch` for semantic search of similar implementations |
| `/gsd:execute-phase` | Use `qmd_get` for targeted file retrieval from plan references |
| Context restoration | Use `qmd_multi_get` to batch-load `.planning/` files |

**For GSD agents:** When spawning gsd-executor or gsd-planner subagents, they should:
1. Start with `qmd_query` to understand relevant codebase context
2. Use `qmd_get` to retrieve files referenced in PLAN.md
3. Avoid broad `Glob`/`Read` sweeps that consume tokens

## Resources

**Project Documentation:**
- [Architecture](docs/ARCHITECTURE.md) - System design and hexagonal architecture
- [Local Development](docs/LOCAL_DEVELOPMENT.md) - Setup guide
- [Agent Routing](docs/AGENTS.md) - AI agent navigation guide

**External References:**
- [gqlgen Documentation](https://gqlgen.com/)
- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [Effective Go](https://go.dev/doc/effective_go)
- [PostgreSQL 18 Documentation](https://www.postgresql.org/docs/18/)
