# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**Perspectize** — Platform for storing, refining, and sharing perspectives on content (initially YouTube videos).

Monorepo with two stacks:
- **Backend:** `perspectize-go/` — Go GraphQL API (see `perspectize-go/CLAUDE.md`)
- **Frontend:** `perspectize-fe/` — SvelteKit web app (see `perspectize-fe/CLAUDE.md`)

**Important:** `perspectize-be/` contains legacy C# code. **Do not modify, except to delete.** All backend work happens in `perspectize-go/`.

**CLAUDE.md structure:** Root file (this) contains shared concerns. Package-level files contain stack-specific instructions. Claude loads root + the relevant package file per session.

## GitHub & Repository Management

**Always use `gh` CLI** for GitHub operations. Do not use MCP plugins.

```bash
# Pull requests
gh pr create --title "Title" --body "Description"
gh pr list
gh pr view 123
gh pr merge 123

# Edit PR (use API — gh pr edit fails with Projects Classic deprecation)
gh api repos/CodeWarrior-debug/perspectize-be/pulls/123 -X PATCH -f body="New description"

# Issues (use API — gh issue view fails with Projects Classic deprecation)
gh issue create --title "Title" --body "Description"
gh issue list
gh api repos/CodeWarrior-debug/perspectize-be/issues/123 --jq '.title, .html_url'

# API access
gh api repos/CodeWarrior-debug/perspectize-be/pulls/123/comments
```

GitHub Projects v2: See [docs/GITHUB_PROJECTS.md](docs/GITHUB_PROJECTS.md).

## Branch Naming

**Always branch from updated `main`:** `git checkout main && git pull origin main && git checkout -b <name>`

**Format:** `type/initiativePrefix-issueNumber-description-in-kebab-case`

| Component | Values |
|-----------|--------|
| **type** | `feature`, `bugfix`, `chore` |
| **initiativePrefix** | `INI` (Initialization Phase) |
| **issueNumber** | GitHub issue number |

Example: `feature/INI-16-youtube-post-graphql`

### GitHub Issues with GSD Plans

Include: GSD Plan Reference (`.planning/phases/{phase}/{plan}-PLAN.md`), acceptance criteria from `must_haves.truths`, dependencies if present.

## Agent Delegation

| Task | Model | Subagent |
|------|-------|----------|
| Architecture decisions | Opus | — |
| Go implementation | Sonnet | `go-backend` |
| GraphQL schema | Sonnet | `graphql-designer` |
| DB migrations | Sonnet | `db-migration` |
| Code review | Haiku | `code-reviewer` |
| Test generation | Haiku | `test-writer` |

## GSD Workflow

Planning and execution artifacts in `.planning/`: `PROJECT.md`, `ROADMAP.md`, `STATE.md`, `phases/`. Branching: see [docs/GSD_BRANCHING.md](docs/GSD_BRANCHING.md).

## Self-Verification

Before marking work complete, verify against plan `must_haves` and capture evidence. See [docs/VERIFICATION.md](docs/VERIFICATION.md) for checklist and evidence capture workflow.

## Code Search with qmd

**Prefer qmd MCP tools over Read/Glob for exploration.** Two collections:

| Collection | Scope | Use for |
|------------|-------|---------|
| `perspectize` | Source code, docs, configs | Codebase understanding, pattern discovery |
| `planning` | `.planning/` files | Project context, roadmap, research, completed phases |

| Tool | When to use |
|------|-------------|
| `qmd_search` | Quick keyword lookup |
| `qmd_vsearch` | Semantic/concept search |
| `qmd_query` | Complex questions (BM25 + vector + reranking) |
| `qmd_get` / `qmd_multi_get` | Retrieve specific files from search results |

**Stable vs live files:** Use qmd for stable reference (PROJECT.md, ROADMAP.md, research/*, completed SUMMARYs). **Always `Read` fresh:** STATE.md and the current phase PLAN.md — qmd index may be stale.

**For GSD agents:** Start with `qmd_query` for context, `Read` STATE.md and current PLAN.md fresh, avoid broad `Glob`/`Read` sweeps.

**Re-index:** `qmd update && qmd embed`

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

### qmd Pre-Exploration for GSD Planning

Before spawning `gsd-phase-researcher` or `gsd-planner`, use qmd to gather codebase context efficiently:

1. **Update index first:** `qmd update && qmd embed` (ensures index matches current branch)
2. **Broad sweep:** `mcp__qmd__query` — "what source files exist in {directory}?" / "what patterns are established?"
3. **Targeted detail:** `mcp__qmd__get` or `mcp__qmd__vsearch` for specific files or concepts
4. **Feed digests into planner prompt** as `<codebase_context>` instead of inlining raw file contents

Only inline raw code for **files being directly modified** by the plan. Use qmd summaries for everything else. This reduces planner token usage by 40-60%.

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
Also test any frontend GraphQL queries (`src/lib/queries/*.ts`) against the live backend to catch schema drift.

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

Save screenshots to `~/Downloads/screenshots/` with naming convention:
- **Prefix:** `ccsv-` (Claude Code Self Verification)
- **Format:** `ccsv-{plan}-{description}-{width}.png`
- **Example:** `ccsv-01-02-mobile-375px.png`, `ccsv-01-04-ag-grid-desktop-1280px.png`
- Use `filePath` parameter on `take_screenshot` to save directly
- Take full-page screenshots (`fullPage: true`) at mobile (375px), tablet (768px), desktop (1280px)

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

## Resources

- [Architecture](docs/ARCHITECTURE.md) — System design and hexagonal architecture
- [Local Development](docs/LOCAL_DEVELOPMENT.md) — Setup guide
- [Agent Routing](docs/AGENTS.md) — AI agent navigation guide
- [gqlgen](https://gqlgen.com/) | [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/) | [Effective Go](https://go.dev/doc/effective_go) | [PostgreSQL 18](https://www.postgresql.org/docs/18/)
