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

## Resources

- [Architecture](docs/ARCHITECTURE.md) — System design and hexagonal architecture
- [Local Development](docs/LOCAL_DEVELOPMENT.md) — Setup guide
- [Agent Routing](docs/AGENTS.md) — AI agent navigation guide
- [gqlgen](https://gqlgen.com/) | [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/) | [Effective Go](https://go.dev/doc/effective_go) | [PostgreSQL 18](https://www.postgresql.org/docs/18/)
