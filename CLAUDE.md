# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**Perspectize** — Platform for storing, refining, and sharing perspectives on content (initially YouTube videos).

Monorepo with two stacks:
- **Backend:** `backend/` — Go GraphQL API (see `backend/CLAUDE.md`)
- **Frontend:** `frontend/` — SvelteKit web app (see `frontend/CLAUDE.md`)

**CLAUDE.md structure:** Root file (this) contains shared concerns. Package-level files contain stack-specific instructions. Claude loads root + the relevant package file per session.

## GitHub & Repository Management

**Always use `gh` CLI** for GitHub operations. Do not use MCP plugins.

```bash
# Pull requests
gh pr create --title "Title" --body "..."  # Use PR template (see below)
gh pr list
gh pr view 123
gh pr merge 123

# Edit PR (use API — gh pr edit fails with Projects Classic deprecation)
gh api repos/CodeWarrior-debug/perspectize/pulls/123 -X PATCH -f body="New description"

# Issues (use API — gh issue view fails with Projects Classic deprecation)
gh issue create --title "Title" --body "..."  # Use issue templates (see below)
gh issue list
gh api repos/CodeWarrior-debug/perspectize/issues/123 --jq '.title, .html_url'

# API access
gh api repos/CodeWarrior-debug/perspectize/pulls/123/comments
```

### GitHub Templates

**Always use the repository templates** in `.github/` when creating PRs and issues.

**Pull Requests** — follow `.github/pull_request_template.md` (Summary, Problem, Solution, Technical Changes, Testing, Checklist, Related Issues).

**Issues** — use templates from `.github/ISSUE_TEMPLATE/` (feature_request.md or bug_report.md).

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

## Agent Delegation Strategy

| Task Type | Model | Subagent | Rationale |
|-----------|-------|----------|-----------|
| Architecture decisions | Opus | - | Complex multi-file reasoning |
| Go implementation | Sonnet | `go-backend` | Balanced quality/cost |
| GraphQL schema design | Sonnet | `graphql-designer` | Schema patterns |
| Database migrations | Sonnet | `db-migration` | SQL generation |
| Code review | Haiku | `code-reviewer` | Fast pattern matching |
| Test generation | Haiku | `test-writer` | Boilerplate generation |

## Coding Conventions

**Learning comments:** Mark temporary explanatory comments with `*TEMP*` for easy grep/removal:
```go
// *TEMP* - defer runs after function returns, ensures cleanup
defer db.Close()
```

**Commit messages:** Conventional commit format (`feat`, `fix`, `refactor`, `chore`, `docs`, `test`). One logical change per commit.

## GSD Workflow

Planning and execution artifacts in `.planning/`: `PROJECT.md`, `ROADMAP.md`, `STATE.md`, `phases/`. Branching: see [docs/GSD_BRANCHING.md](docs/GSD_BRANCHING.md).

## Self-Verification (MANDATORY)

**Before claiming work is complete, pushing, or creating a PR**, you MUST run verification. No exceptions.

### Verification checklist

1. **Build**: `go build ./...` in `backend/` — must compile with zero errors
2. **Backend tests**: `go test ./...` in `backend/` — all must pass
3. **Frontend tests**: `pnpm run test:run` in `frontend/` — all must pass
4. **Stale references**: If renaming/moving files or paths, grep the entire repo for old names
5. **Plan must_haves**: If executing a GSD plan, verify each `must_haves.truths` item

Run the relevant subset (e.g., backend-only changes skip step 3). Report results explicitly — don't just say "tests pass", show the output summary.

See [docs/VERIFICATION.md](docs/VERIFICATION.md) for evidence capture workflow.

## Resources

**Monorepo docs:**
- [Architecture](docs/ARCHITECTURE.md) — System design and hexagonal architecture
- [Local Development](docs/LOCAL_DEVELOPMENT.md) — Setup guide
- [Agent Routing](docs/AGENTS.md) — AI agent navigation guide
- [Domain Guide](docs/DOMAIN_GUIDE.md) — Domain layer rules and patterns
- [Go Patterns](docs/GO_PATTERNS.md) — Error handling and DB query patterns

**Frontend docs:**
- [Frontend CLAUDE.md](frontend/CLAUDE.md) — SvelteKit, Svelte 5, TanStack Query patterns
- [Design Spec](frontend/docs/DESIGN_SPEC.md) — Figma design system, color tokens, typography, component specs
- [Figma Reference](frontend/docs/FIGMA.md) — File keys, pages, variables, code↔Figma mapping

**Planning & backlog:**
- [Feature Backlog](FEATURE_BACKLOG.md) — Future ideas and enhancements not tied to any milestone. Capture ideas here during development; evaluate when planning new work.

**External references:**
- [gqlgen](https://gqlgen.com/) | [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/) | [Effective Go](https://go.dev/doc/effective_go) | [PostgreSQL 17](https://www.postgresql.org/docs/17/)
