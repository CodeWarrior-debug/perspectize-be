# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**Perspectize** — Platform for storing, refining, and sharing perspectives on content (initially YouTube videos).

Monorepo with two stacks:
- **Backend:** `backend/` — Go GraphQL API (see `backend/CLAUDE.md`)
- **Frontend:** `frontend/` — SvelteKit web app (see `frontend/CLAUDE.md`)

**CLAUDE.md structure:** Root file (this) contains shared concerns. Package-level files contain stack-specific instructions. Claude loads root + the relevant package file per session.

## Context Lookup (qmd retired — graphify coming soon)

**graphify coming soon.** qmd is no longer the mandated pre-search step — use Read/Glob/Grep directly for now. The section below is kept commented out for reference until graphify lands.

<!--
## Context Lookup with qmd (MANDATORY — Use Before Reading Files)

**ALWAYS use qmd bash commands** to search for code before using Read/Glob/Grep. This applies to ALL agents including GSD subagents.

⚠️ **DO NOT use MCP qmd tools** — use Bash commands only. MCP is not available in all contexts.

⚠️ **Cloud sessions (claude.ai/code web):** qmd is NOT available. Skip qmd entirely and use Read/Glob/Grep directly.

**Allowed commands (pre-approved):**
- `qmd search *` — keyword search
- `qmd vsearch *` — semantic search
- `qmd query *` — hybrid search with reranking
- `qmd get *` — retrieve file content
- `qmd ls *` — list files in collection
- `qmd update` — refresh index after changes
- `qmd status` — check index status
- `qmd embed` — generate embeddings after update

```bash
# Quick keyword search (BM25) — use 80% of the time
qmd search "auth middleware" -c perspectize

# Semantic search — finds related code even without exact keywords
qmd vsearch "how does error handling work" -c perspectize

# Hybrid search with reranking — best quality for complex questions
qmd query "checkout flow validation" -c perspectize

# Get specific file (optionally from line N, max L lines)
qmd get qmd://perspectize/backend/internal/domain/content.go
qmd get qmd://perspectize/backend/internal/domain/content.go:45 -l 30

# List files in collection
qmd ls perspectize
```

**Workflow:**
1. Use `qmd search` or `qmd query` first to find relevant code
2. Use `qmd get` to retrieve targeted snippets (not full files)
3. Fall back to Read/Glob only if qmd doesn't return enough results

**Available collections:** `perspectize` (code), `planning` (GSD docs)

**Update index after major changes:** `qmd update && qmd embed`
-->

## GitHub & Repository Management

**Git branch gotcha:** Local default branch is `master`, remote is `main`. Use `origin/main` (not `main`) for diff/log comparisons: `git diff origin/main...HEAD`.

**Always use `gh` CLI** for GitHub operations. Do not use MCP plugins.

**Note:** In Claude Code web sessions, `gh` CLI may not be authenticated. If `gh` auth fails:
- **Creating a PR:** Push the branch with `git push -u origin <branch>` and let the user create the PR via the GitHub UI button. Prepare the PR title and body as copyable text for the user.
- **Updating a PR:** Output the updated title/body as copyable text so the user can paste it into the GitHub UI.

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

GitHub Projects v2: See [.docs/GITHUB_PROJECTS.md](.docs/GITHUB_PROJECTS.md).

### PR Merge Preferences

```bash
gh pr merge 123 --squash --delete-branch --admin
```

- `--squash` — Single commit (cleaner history)
- `--delete-branch` — Auto-delete branch after merge
- `--admin` — Bypass branch protection when needed

## Branch Naming

**Always branch from updated `main`:** `git checkout main && git pull origin main && git checkout -b <name>`

**Format:** `type/initiativePrefix-issueNumber-description-in-kebab-case`

| Component | Values |
|-----------|--------|
| **type** | `feature`, `bugfix`, `chore` |
| **initiativePrefix** | `INI` (Initialization Phase) |
| **issueNumber** | GitHub issue number |

Example: `feature/INI-16-youtube-post-graphql`

### GitHub Issues with Plans

Include a plan reference and dependencies if present: for new work, the superpowers plan/spec path (`docs/superpowers/plans/{name}-plan.md`); for legacy in-flight work, the GSD plan reference (`.planning/phases/{phase}/{plan}-PLAN.md`) and acceptance criteria from `must_haves.truths`.

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

**No chained bash commands:** Do not use `&&` to chain shell commands. Run each command as a separate Bash tool call. Chained commands don't match permission allow-list patterns and block on approval prompts. This applies to all agents and subagents.

**Migration numbering:** Always check existing migration files before creating new ones. Plan-specified numbers may be stale — use `ls backend/migrations/ | tail -5` to find the next available number.

**Commit messages:** Conventional commit format (`feat`, `fix`, `refactor`, `chore`, `docs`, `test`). One logical change per commit. GSD planning work (PLAN.md, CONTEXT.md, RESEARCH.md, ROADMAP.md) uses the `docs` tag — e.g., `docs(11,13): create execution plans`.

## Planning & Execution Workflow

**Primary workflow: obra/superpowers** (plugin enabled in `.claude/settings.json`). Use `superpowers:writing-plans` (or its brainstorming/spec-writing counterparts) for planning, and `superpowers:executing-plans` / `superpowers:subagent-driven-development` for execution. Plans and specs live in `docs/superpowers/plans/` and `docs/superpowers/specs/` — see `docs/superpowers/plans/2026-08-15-clerk-derived-user-identity-plan.md` for the established format (plan header names the required execution sub-skill, links its spec, checkbox-tracked (`- [ ]`) tasks).

**GSD is legacy — do NOT start new work with it.** Some milestones still have unfinished work tracked under the old workflow in `.planning/phases/` (`PROJECT.md`, `ROADMAP.md`, `STATE.md`, phase `PLAN.md`/`must_haves.truths` files). Finish those specific in-flight phases using their existing GSD plan files/commands rather than replanning them from scratch under superpowers — don't discard partially-done GSD work. All new planning and execution goes through superpowers. Branching for legacy GSD phases: see [.docs/GSD_BRANCHING.md](.docs/GSD_BRANCHING.md).

## Self-Verification (MANDATORY)

**Before claiming work is complete, pushing, or creating a PR**, you MUST run verification. No exceptions.

### Verification checklist

1. **Build**: `go build ./...` in `backend/` — must compile with zero errors
2. **Backend tests**: `go test ./...` in `backend/` — all must pass
3. **Frontend tests**: `pnpm run test:run` in `frontend/` — all must pass
4. **Stale references**: If renaming/moving files or paths, grep the entire repo for old names
5. **Plan must_haves**: If executing a GSD plan, verify each `must_haves.truths` item

Run the relevant subset (e.g., backend-only changes skip step 3). Report results explicitly — don't just say "tests pass", show the output summary.

**GSD verification is not self-verification.** The GSD verifier checks must_haves against codebase structure. It does NOT run builds or tests. Always run the full checklist (build, backend tests, frontend tests) before creating a PR, even after GSD verification passes.

See [.docs/VERIFICATION.md](.docs/VERIFICATION.md) for evidence capture workflow.

## Resources

**Monorepo docs:**
- [Architecture](.docs/ARCHITECTURE.md) — System design and hexagonal architecture
- [Local Development](.docs/LOCAL_DEVELOPMENT.md) — Setup guide
- [Agent Routing](.docs/AGENTS.md) — AI agent navigation guide
- [Domain Guide](.docs/DOMAIN_GUIDE.md) — Domain layer rules and patterns
- [Go Patterns](.docs/GO_PATTERNS.md) — Error handling and DB query patterns
- [Security](.docs/SECURITY.md) — Secret management, rotation procedures, incident response
- [Dependency Security](.docs/DEPENDENCY_SECURITY.md) — Trivy/pnpm-audit scanning, CVE remediation workflow, CI gotchas

**Frontend docs:**
- [Frontend CLAUDE.md](frontend/CLAUDE.md) — SvelteKit, Svelte 5, TanStack Query patterns
- [Design Spec](frontend/docs/DESIGN_SPEC.md) — Figma design system, color tokens, typography, component specs
- [Figma Reference](frontend/docs/FIGMA.md) — File keys, pages, variables, code↔Figma mapping

**How-to guides:**
- [Adding an AG Grid Column](.claude/docs/ADDING_AG_GRID_COLUMN.md) — Decision checklist for adding columns to the ActivityTable
- [Adding a Content Type](.claude/docs/ADDING_CONTENT_TYPE.md) — End-to-end guide for new content types (backend + frontend)
- [Code to Figma Canvas](.claude/docs/CODE_TO_FIGMA_CANVAS.md) — Capture running app into Figma to keep designs in sync

**Planning & backlog:**
- [Feature Backlog](FEATURE_BACKLOG.md) — Future ideas and enhancements not tied to any milestone. Capture ideas here during development; evaluate when planning new work.
- [Bug Tracking](.docs/BUG_TRACKING.md) — How known bugs are tracked privately (gitignored files, persistent bugs phase)

**Bug logging (MANDATORY):** When you discover a bug during development, review, or testing, log it in `.planning/phases/bugs/BACKLOG.md` with severity and location. Also create a GitHub issue using the bug report template — keep sensitive details (exact paths, line numbers, security specifics) in the backlog only. When a bug is fixed, move it to `.planning/phases/bugs/CLOSED.md` with the PR reference. These files are gitignored — never commit them.

**Hookify rules (check `.claude/hookify.*.local.md` for all rules):**
- **qmd first:** Spawning an Explore agent triggers a warning to use `qmd search`/`qmd get` first. Only spawn an Explore agent if qmd doesn't return what you need.
- **Pre-PR:** `gh pr create` is blocked until `claude-md-management:revise-claude-md` has been run. The block can't detect completion, so use `gh api` to create the PR after running the skill. Example: `gh api repos/CodeWarrior-debug/perspectize/pulls -f title="..." -f body="..." -f head="branch" -f base="main"`
- **Pre-commit tests:** `git commit` triggers a warning to verify test coverage for new/modified frontend `src/` files. Config, styles, docs, and test files are exempt.
- **Pre-commit prettier:** `git commit` triggers a warning to run `pnpm exec prettier --write` on staged frontend files.
- **HEREDOC gotcha:** `git commit -m "$(cat <<'EOF'...)"` syntax breaks hookify's regex parser. Use simple quoted strings for commit messages instead.

**Cowork session cleanup:** Claude cowork (claude.ai web) sessions leave `_tmp_*` files and conversation transcript `.txt` files in the repo root and `frontend/`. Delete these before committing.

## Merge Conflict Patterns

**pnpm-lock.yaml conflicts:** Accept either side (`git checkout --theirs frontend/pnpm-lock.yaml`), then regenerate: `pnpm install --dir frontend`. Always use `--dir` instead of `cd` to avoid hook/shell side effects that can switch branches mid-operation.

**External references:**
- [gqlgen](https://gqlgen.com/) | [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/) | [Effective Go](https://go.dev/doc/effective_go) | [PostgreSQL 17](https://www.postgresql.org/docs/17/)
