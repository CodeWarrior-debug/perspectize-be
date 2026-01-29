# CLAUDE.md - Backend

This file provides guidance for Claude Code when working with the Perspectize Go backend.

## Overview

**Perspectize Go Backend** - GraphQL API for storing, refining, and sharing perspectives on content (initially YouTube videos). Built with Go 1.23+, PostgreSQL 16+, and GraphQL using gqlgen.

## Architecture

This project follows **Hexagonal Architecture** (Ports and Adapters pattern):

```
apps/backend/
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
│   └── database/            # Database connection utilities
├── migrations/              # SQL migration files
├── test/                    # Test files
└── schema.graphql           # GraphQL schema definition
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

**Dependency Rule:** Dependencies point inward. Domain never depends on adapters.

## Development Commands

All commands should be run from the `apps/backend/` directory.

```bash
# Run the server
make run

# Run with hot-reload
make dev

# Run tests
make test

# Run tests with coverage
make test-coverage

# Format code
make fmt

# Generate GraphQL code (after schema changes)
make graphql-gen

# Database migrations
make migrate-up
make migrate-down
make migrate-create

# Docker (PostgreSQL)
make docker-up
make docker-down
```

## Configuration

Required environment variables:
- `DATABASE_URL` - PostgreSQL connection string

Optional:
- `YOUTUBE_API_KEY` - YouTube Data API v3 key

## Technology Stack

- **Language:** Go 1.23+
- **GraphQL:** gqlgen (schema-first)
- **Database:** PostgreSQL 16+ with sqlx + pgx/v5
- **Migrations:** golang-migrate
- **Testing:** testing + testify

## Domain Entities

- **Content** - Media that users create perspectives on (YouTube videos)
- **User** - Users who create perspectives
- **Perspective** - A user's viewpoint/rating on content

## Adding New Features

1. Define domain model: `internal/core/domain/`
2. Define repository interface: `internal/core/ports/repositories/`
3. Implement business logic: `internal/core/services/`
4. Implement repository: `internal/adapters/repositories/postgres/`
5. Update GraphQL schema: `schema.graphql`
6. Generate GraphQL code: `make graphql-gen`
7. Implement resolver: `internal/adapters/graphql/resolvers/`
8. Wire in main: `cmd/server/main.go`
9. Write tests: `test/`

## Patterns & Gotchas

### GraphQL Schema Defaults
When a GraphQL field has a default value (e.g., `first: Int = 10`), gqlgen passes the default to the resolver as a non-nil pointer, not `nil`.

### Rating Validation
Ratings are validated at 0-10000 range both in service layer and database (custom domain type).

### Claim Uniqueness
Perspectives have unique (claim, user_id) constraint - each user can only have one perspective with a given claim.

## Resources

- [gqlgen Documentation](https://gqlgen.com/)
- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [Effective Go](https://go.dev/doc/effective_go)
