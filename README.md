# Perspectize

A platform for storing, refining, and sharing perspectives on content.

## Overview

Perspectize allows users to create structured perspectives on media content (initially YouTube videos). Each perspective includes claims, ratings (quality, agreement, importance, confidence), and categorized feedback.

## Repository Structure

```
perspectize/
├── apps/
│   ├── backend/          # Go GraphQL API
│   │   ├── cmd/server/   # Application entry point
│   │   ├── internal/     # Core business logic (hexagonal architecture)
│   │   ├── migrations/   # Database migrations
│   │   └── schema.graphql
│   └── frontend/         # Svelte web app (planned)
├── packages/             # Shared packages (future)
├── Makefile              # Root orchestration
└── docker-compose.yml    # Development services
```

## Quick Start

### Prerequisites

- Go 1.23+
- Docker and Docker Compose
- PostgreSQL 16+ (via Docker or local install)

### Setup

```bash
# Clone the repository
git clone https://github.com/CodeWarrior-debug/perspectize-be.git
cd perspectize-be

# Start PostgreSQL
make docker-up

# Run database migrations
make migrate-up

# Start the backend server
make backend-run
```

The GraphQL playground will be available at http://localhost:8080/graphql

### Development

```bash
# Run backend with hot-reload
make backend-dev

# Run all tests
make test

# Generate GraphQL code after schema changes
make backend-graphql
```

## Architecture

The backend follows **Hexagonal Architecture** (Ports and Adapters):

- **Domain Layer** (`internal/core/domain/`) - Pure business entities
- **Ports** (`internal/core/ports/`) - Interface contracts
- **Services** (`internal/core/services/`) - Business logic
- **Adapters** (`internal/adapters/`) - Infrastructure implementations

See [apps/backend/CLAUDE.md](apps/backend/CLAUDE.md) for detailed backend documentation.

## Tech Stack

### Backend
- **Language:** Go 1.23+
- **API:** GraphQL (gqlgen)
- **Database:** PostgreSQL 16+
- **ORM:** sqlx + pgx/v5

### Frontend (Planned)
- **Framework:** Svelte/SvelteKit
- **State:** TanStack Query
- **Data Grid:** AG Grid

## API

### GraphQL Queries

```graphql
# Get content by ID
query {
  contentByID(id: 1) {
    id
    title
    contentType
    sourceUrl
  }
}

# List perspectives with pagination
query {
  perspectives(first: 10) {
    edges {
      node {
        id
        claim
        quality
        agreement
      }
    }
    pageInfo {
      hasNextPage
    }
  }
}
```

### GraphQL Mutations

```graphql
# Create a perspective
mutation {
  createPerspective(input: {
    claim: "This video explains the topic well"
    userID: 1
    contentID: 1
    quality: 8500
    agreement: 9000
  }) {
    id
    claim
  }
}
```

## Contributing

1. Create a feature branch: `git checkout -b feature/INI-XX-description`
2. Make your changes
3. Run tests: `make test`
4. Create a pull request

## License

Private repository - All rights reserved.
