# Root Makefile - Perspectize Monorepo
# Orchestrates commands across all apps

.PHONY: help backend-run backend-dev backend-test frontend-dev docker-up docker-down

# Default target
help:
	@echo "Perspectize Monorepo Commands"
	@echo ""
	@echo "Backend (Go GraphQL API):"
	@echo "  make backend-run        - Run the backend server"
	@echo "  make backend-dev        - Run backend with hot-reload"
	@echo "  make backend-test       - Run backend tests"
	@echo "  make backend-lint       - Lint backend code"
	@echo "  make backend-graphql    - Generate GraphQL code"
	@echo ""
	@echo "Database:"
	@echo "  make docker-up          - Start PostgreSQL"
	@echo "  make docker-down        - Stop PostgreSQL"
	@echo "  make migrate-up         - Run database migrations"
	@echo "  make migrate-down       - Rollback last migration"
	@echo ""
	@echo "Frontend (Svelte - not yet implemented):"
	@echo "  make frontend-dev       - Start frontend dev server"
	@echo "  make frontend-build     - Build frontend for production"
	@echo ""
	@echo "Full Stack:"
	@echo "  make dev                - Start all services for development"
	@echo "  make test               - Run all tests"

# Backend commands (delegate to apps/backend/Makefile)
backend-run:
	$(MAKE) -C apps/backend run

backend-dev:
	$(MAKE) -C apps/backend dev

backend-test:
	$(MAKE) -C apps/backend test

backend-lint:
	$(MAKE) -C apps/backend lint

backend-graphql:
	$(MAKE) -C apps/backend graphql-gen

# Database commands
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

migrate-up:
	$(MAKE) -C apps/backend migrate-up

migrate-down:
	$(MAKE) -C apps/backend migrate-down

# Frontend commands (placeholder)
frontend-dev:
	@echo "Frontend not yet implemented. See apps/frontend/README.md"

frontend-build:
	@echo "Frontend not yet implemented. See apps/frontend/README.md"

# Full stack development
dev: docker-up
	@echo "Starting backend..."
	$(MAKE) -C apps/backend dev

# Run all tests
test: backend-test
	@echo "All tests completed"
