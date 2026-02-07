# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**Perspectize** - Platform for storing, refining, and sharing perspectives on content (initially YouTube videos).

This is a monorepo with two stacks:
- **Backend:** `perspectize-go/` - Go GraphQL API (see `perspectize-go/CLAUDE.md` for backend-specific instructions)
- **Frontend:** `perspectize-fe/` - SvelteKit web app (see `perspectize-fe/CLAUDE.md` for frontend-specific instructions)

**Important:** The `perspectize-be/` directory contains legacy C# ASP.NET Core code. **Do not modify, except to delete.** All backend development happens in `perspectize-go/`.

## GitHub & Repository Management

**Always use the `gh` CLI** for all GitHub operations. Do not use MCP plugins or other tools for GitHub interactions.

```bash
# Pull requests
gh pr create --title "Title" --body "Description"
gh pr list
gh pr view 123
gh pr merge 123

# Edit PR (use REST API - gh pr edit may fail with Projects Classic deprecation error)
gh api repos/CodeWarrior-debug/perspectize-be/pulls/123 -X PATCH -f body="New description"

# Issues (use API - gh issue view fails with Projects Classic deprecation error)
gh issue create --title "Title" --body "Description"
gh issue list
gh api repos/CodeWarrior-debug/perspectize-be/issues/123 --jq '.title, .html_url'

# Repository info
gh repo view

# API access (for anything not covered by commands)
gh api repos/CodeWarrior-debug/perspectize-be/pulls/123/comments
```

### GitHub Projects (v2)

See [docs/GITHUB_PROJECTS.md](docs/GITHUB_PROJECTS.md) for GraphQL API queries and token scope setup.

## Branch Naming Convention

**Always create branches from an updated `main` branch.**

```bash
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

## Agent Delegation Strategy

| Task Type | Model | Subagent | Rationale |
|-----------|-------|----------|-----------|
| Architecture decisions | Opus | - | Complex multi-file reasoning |
| Go implementation | Sonnet | `go-backend` | Balanced quality/cost |
| GraphQL schema design | Sonnet | `graphql-designer` | Schema patterns |
| Database migrations | Sonnet | `db-migration` | SQL generation |
| Code review | Haiku | `code-reviewer` | Fast pattern matching |
| Test generation | Haiku | `test-writer` | Boilerplate generation |

## Workflow Integration

This project uses **GSD workflow** for planning and execution. See `.planning/` for:
- `PROJECT.md` - Project definition and requirements
- `ROADMAP.md` - Phase-based milestone planning
- `STATE.md` - Current position and accumulated context
- `phases/` - Detailed execution plans

### GSD Workflow Branching

See [docs/GSD_BRANCHING.md](docs/GSD_BRANCHING.md) for stacked PR workflow and `.planning/config.json` branching config.

## Self-Verification Workflow

Before marking any work complete, run interactive verification:

### GSD Plan Verification

For each plan's `must_haves`:

| Check | Command |
|-------|---------|
| `truths` | Run actual command, verify output |
| `artifacts.path` | `test -f {path} && echo "exists"` |
| `artifacts.contains` | `grep -q "{pattern}" {path}` |
| `artifacts.min_lines` | `wc -l < {path}` ≥ N |
| `key_links.pattern` | `grep -q "{pattern}" {from}` |

### Evidence Capture

Before creating PR:
- Screenshot at mobile (375px), tablet (768px), desktop (1024px+)
- Console output showing no errors
- Verification commands output

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
