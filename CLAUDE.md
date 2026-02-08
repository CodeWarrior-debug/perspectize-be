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

## Coding Conventions

**Naming:** Balance brevity with clarity — names should be self-describing without needing comments. Avoid both cryptic abbreviations and excessive verbosity.

**Comments:** Explain **why**, not **what**. Use comments for non-obvious decisions, not to narrate code.

**Learning comments:** Mark temporary explanatory comments with `*TEMP*` for easy grep/removal:
```go
// *TEMP* - defer runs after function returns, ensures cleanup
defer db.Close()
```

**Commit messages:** Use conventional commit format. Keep commits focused and atomic — one logical change per commit.
```
type(scope): short description

Optional body explaining why, not what.
```
Types: `feat`, `fix`, `refactor`, `chore`, `docs`, `test`

## Approved Tools & Permissions

Claude Code may use the following tools **without prompting** for this project:

### File Operations
- **Read**: Any file in `.planning/`, `perspectize-fe/src/`, `perspectize-go/`, `docs/`
- **Edit/Write**: Any file in `.planning/` (design specs, roadmaps, plans), `perspectize-fe/src/`, `perspectize-go/` (excluding sensitive files)
- **Glob/Grep**: Unlimited search across the codebase

### Bash Commands
- **Safe navigation**: `git`, `pnpm`, `npm`, `ls`, `pwd`, `cd`, `cat`, `head`, `tail`
- **Git operations**: All git commands (clone, checkout, pull, push, commit, log, diff)
- **Package management**: `pnpm install`, `pnpm run`
- **Forbidden**: `rm -rf`, `sudo`, `chmod 777` (these require explicit permission)

### Figma MCP Tools
- **Design context**: `get_screenshot`, `get_design_context`, `get_metadata` for fileKey `SyvrP9yYbrmCorofJK4Co8` (Perspectize Figma file)
- **Design system**: `get_variable_defs`

### Task & Execution Tools
- **TaskCreate/TaskUpdate**: Creating and updating tasks from conversations
- **Agents**: Launching Task tool subagents for exploration, planning, code review

**Rationale:** These tools are safe for this workflow. The project uses GSD planning, TanStack + Svelte 5 stack, and Figma for design. File reads/edits focus on spec updates and source code. Bash is limited to safe operations.

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
