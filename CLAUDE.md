# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**Perspectize** — Platform for storing, refining, and sharing perspectives on content (initially YouTube videos).

Monorepo with two stacks:
- **Backend:** `perspectize-go/` — Go GraphQL API (see `perspectize-go/CLAUDE.md`)
- **Frontend:** `perspectize-fe/` — SvelteKit web app (see `perspectize-fe/CLAUDE.md`)

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

## Agent Delegation Strategy

| Task Type | Model | Subagent | Rationale |
|-----------|-------|----------|-----------|
| Architecture decisions | Opus | - | Complex multi-file reasoning |
| Go implementation | Sonnet | `go-backend` | Balanced quality/cost |
| GraphQL schema design | Sonnet | `graphql-designer` | Schema patterns |
| Database migrations | Sonnet | `db-migration` | SQL generation |
| Code review | Haiku | `code-reviewer` | Fast pattern matching |
| Test generation | Haiku | `test-writer` | Boilerplate generation |

## GSD Workflow

Planning and execution artifacts in `.planning/`: `PROJECT.md`, `ROADMAP.md`, `STATE.md`, `phases/`. Branching: see [docs/GSD_BRANCHING.md](docs/GSD_BRANCHING.md).

## Self-Verification

Before marking work complete, verify against plan `must_haves` and capture evidence. See [docs/VERIFICATION.md](docs/VERIFICATION.md) for full checklist and evidence capture workflow.

### Production Setup (Sevalla/Fly.io)

Use `DATABASE_URL` with external endpoint from hosting provider. Note: Sevalla connections may require `?sslmode=disable` and may succeed on second attempt.

## Resources

- [Architecture](docs/ARCHITECTURE.md) — System design and hexagonal architecture
- [Local Development](docs/LOCAL_DEVELOPMENT.md) — Setup guide
- [Agent Routing](docs/AGENTS.md) — AI agent navigation guide
- [Domain Guide](docs/DOMAIN_GUIDE.md) — Domain layer rules and patterns
- [Go Patterns](docs/GO_PATTERNS.md) — Error handling and DB query patterns
- [gqlgen](https://gqlgen.com/) | [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/) | [Effective Go](https://go.dev/doc/effective_go) | [PostgreSQL 17](https://www.postgresql.org/docs/17/)
