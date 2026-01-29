# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with the Perspectize monorepo.

## Project Overview

**Perspectize** - A platform for storing, refining, and sharing perspectives on content. This monorepo contains:

- `apps/backend/` - Go GraphQL API
- `apps/frontend/` - Svelte web application (placeholder)
- `packages/` - Shared packages (if needed)

## Repository Structure

```
perspectize/
├── apps/
│   ├── backend/           # Go GraphQL API
│   │   ├── cmd/server/    # Entry point
│   │   ├── internal/      # Application code
│   │   ├── migrations/    # Database migrations
│   │   ├── schema.graphql # GraphQL schema
│   │   └── CLAUDE.md      # Backend-specific guidance
│   └── frontend/          # Svelte app (placeholder)
│       └── CLAUDE.md      # Frontend-specific guidance
├── packages/              # Shared packages
├── .github/               # GitHub Actions, templates
├── docker-compose.yml     # Development stack
├── Makefile               # Root orchestration
└── CLAUDE.md              # This file
```

## Quick Start

```bash
# Start all services (PostgreSQL)
make docker-up

# Run backend
make backend-run

# Run tests
make backend-test
```

## App-Specific Guidance

Each app has its own `CLAUDE.md` with detailed instructions:
- **Backend:** See `apps/backend/CLAUDE.md`
- **Frontend:** See `apps/frontend/CLAUDE.md`

## GitHub Workflow

**Always use the `gh` CLI** for GitHub operations:

```bash
gh pr create --title "Title" --body "Description"
gh issue create --title "Title" --body "Description"
gh pr list
gh issue list
```

### Branch Naming

Format: `type/INI-issueNumber-description-kebab-case`

| Type | Use For |
|------|---------|
| `feature` | New features |
| `bugfix` | Bug fixes |
| `chore` | Maintenance tasks |

Examples:
- `feature/INI-23-users-domain-layer`
- `bugfix/INI-30-fix-pagination`
- `chore/INI-25-cicd-infrastructure`

## Technology Stack

| Component | Technology |
|-----------|------------|
| Backend | Go 1.23+, gqlgen, PostgreSQL 16+ |
| Frontend | Svelte, TanStack, AG Grid (planned) |
| CI/CD | GitHub Actions |
| Database | PostgreSQL with migrations |

## Development Workflow

1. Create GitHub issue for the task
2. Create branch: `type/INI-XX-description`
3. Implement changes
4. Run tests: `make backend-test`
5. Create PR with standardized template
6. @ mention reviewer in GitHub issue for notifications

## Learnings & Best Practices

### Monorepo Workflow
- When chaining PRs, base new branches on previous feature branches
- Update CLAUDE.md files with learnings before creating PRs
- Use @ mentions in GitHub comments for mobile notifications
- When restructuring directories, use `git mv` to preserve history
- Update Go module paths in go.mod AND all import statements when moving code

### Go Backend
- Follow hexagonal architecture strictly
- Domain layer has NO external dependencies
- Use pointers for optional fields in domain structs
- Rating validation: 0-10000 range

### Testing
- Unit tests mock repository interfaces
- Integration tests require PostgreSQL (auto-skip when unavailable)
- Tests in same package use `_test` suffix
- Test function names must be unique across ALL test files in same package (e.g., prefix with entity name: `TestUserGetByID_Success`)

### GitHub & Notifications
- The `gh` CLI uses the user's credentials, so @ mentions won't trigger notifications for the user themselves
- To notify the repository owner, use explicit @ mentions in PR/issue comments
- Use `gh api` for PR edits that may fail with regular `gh pr edit` commands

## Resources

- [Backend README](apps/backend/README.md)
- [GitHub Actions CI](.github/workflows/ci.yml)
- [Contributing Guidelines](.github/pull_request_template.md)
