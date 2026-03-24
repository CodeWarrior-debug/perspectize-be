# Codebase Structure

**Analysis Date:** 2026-02-16

## Directory Layout

```
perspectize/
├── backend/
│   ├── cmd/
│   │   └── server/
│   │       └── main.go                 # Server entry point
│   ├── internal/
│   │   ├── adapters/
│   │   │   ├── graphql/
│   │   │   │   ├── generated/          # gqlgen auto-generated code
│   │   │   │   ├── model/              # GraphQL input/output types
│   │   │   │   └── resolvers/          # GraphQL resolver implementations
│   │   │   ├── repositories/
│   │   │   │   └── postgres/           # PostgreSQL GORM repository implementations
│   │   │   └── youtube/                # YouTube API v3 client
│   │   ├── config/
│   │   │   ├── config.go               # Configuration loading
│   │   │   └── validation.go           # Config validation
│   │   └── core/
│   │       ├── domain/                 # Domain models, no external deps
│   │       ├── ports/                  # Interfaces for repositories/services
│   │       │   ├── repositories/       # Repository port interfaces
│   │       │   └── services/           # Service port interfaces
│   │       └── services/               # Business logic implementations
│   ├── pkg/
│   │   ├── database/                   # Database connection, pooling
│   │   ├── graphql/                    # GraphQL utilities (IntID scalar, timing)
│   │   ├── logger/                     # Structured logging setup
│   │   └── middleware/                 # HTTP middleware
│   ├── migrations/                     # SQL migration files (.up.sql / .down.sql)
│   ├── schema.graphql                  # GraphQL schema definition
│   ├── gqlgen.yml                      # GraphQL code generation config
│   ├── Makefile                        # Build and development tasks
│   ├── go.mod                          # Go module definition
│   ├── go.sum                          # Go dependencies lockfile
│   ├── .env.example                    # Example environment variables
│   └── docker-compose.yml              # Local PostgreSQL container config
│
├── frontend/
│   ├── src/
│   │   ├── routes/
│   │   │   ├── +layout.svelte          # Root layout (QueryClientProvider, Header)
│   │   │   ├── +layout.ts              # Layout config (prerender)
│   │   │   └── +page.svelte            # Home page (ActivityTable)
│   │   ├── lib/
│   │   │   ├── components/
│   │   │   │   ├── shadcn/             # shadcn-svelte UI primitives
│   │   │   │   ├── ActivityTable.svelte   # AG Grid data table wrapper
│   │   │   │   ├── AddVideoDialog.svelte  # Modal for URL input
│   │   │   │   ├── AddVideoPopover.svelte # Lightweight URL input
│   │   │   │   ├── CreateUserPopover.svelte # User creation
│   │   │   │   ├── FormPopover.svelte     # Generic form popover
│   │   │   │   ├── Header.svelte          # Navigation header
│   │   │   │   ├── PageWrapper.svelte     # Layout container
│   │   │   │   ├── UserSelector.svelte    # User dropdown
│   │   │   │   ├── DescriptionTooltip.ts  # AG Grid cell renderer
│   │   │   │   └── TagsTooltip.ts         # AG Grid cell renderer
│   │   │   ├── queries/
│   │   │   │   ├── client.ts           # GraphQL client setup
│   │   │   │   ├── content.ts          # Content GraphQL queries
│   │   │   │   ├── users.ts            # User GraphQL queries
│   │   │   │   ├── keys.ts             # TanStack Query key factories
│   │   │   │   └── hooks/
│   │   │   │       ├── useAddVideo.ts  # Add video mutation hook
│   │   │   │       └── useCreateUser.ts # Create user mutation hook
│   │   │   ├── stores/
│   │   │   │   └── userSelection.svelte.ts # User selection state
│   │   │   ├── utils/                  # Helper utilities
│   │   │   └── assets/                 # Static assets (favicon, fonts)
│   │   ├── app.css                     # Global styles, Tailwind config, design tokens
│   │   └── app.html                    # HTML shell
│   ├── tests/
│   │   ├── unit/                       # Unit tests
│   │   ├── components/                 # Component tests
│   │   ├── fixtures/                   # Test data
│   │   └── helpers/                    # Test utilities
│   ├── package.json                    # Dependencies, pnpm scripts
│   ├── tsconfig.json                   # TypeScript configuration
│   ├── vite.config.ts                  # Vite build configuration
│   ├── vitest.config.ts                # Vitest test configuration
│   ├── svelte.config.js                # SvelteKit configuration
│   ├── tailwind.config.ts              # Tailwind configuration
│   └── .env.example                    # Example build-time env vars
│
├── .docs/
│   ├── ARCHITECTURE.md                 # System design and hexagonal architecture
│   ├── LOCAL_DEVELOPMENT.md            # Development environment setup
│   ├── DOMAIN_GUIDE.md                 # Domain layer rules
│   ├── GO_PATTERNS.md                  # Go error handling and DB patterns
│   ├── AGENTS.md                       # AI agent routing guide
│   └── GITHUB_PROJECTS.md              # GitHub Projects v2 workflow
│
├── .planning/
│   ├── PROJECT.md                      # Project overview
│   ├── ROADMAP.md                      # Feature roadmap
│   ├── STATE.md                        # Current state and blockers
│   ├── codebase/
│   │   ├── STACK.md                    # Technology stack analysis
│   │   ├── INTEGRATIONS.md             # External APIs and services
│   │   ├── ARCHITECTURE.md             # Architecture and layers
│   │   ├── STRUCTURE.md                # This file - directory guide
│   │   ├── CONVENTIONS.md              # Coding standards
│   │   ├── TESTING.md                  # Testing patterns
│   │   └── CONCERNS.md                 # Technical debt and issues (gitignored — private)
│   ├── phases/                         # GSD workflow phases
│   │   ├── bugs/                       # Persistent bug tracking phase (gitignored — private)
│   └── research/                       # Research and investigation
│
├── .claude/
│   ├── docs/                           # How-to guides
│   ├── agents/                         # Subagent definitions
│   └── skills/                         # Reusable skill definitions
│
├── README.md                           # Project overview and setup
├── FEATURE_BACKLOG.md                  # Future features not tied to phases
├── CLAUDE.md                           # Root Claude instructions
├── LICENSE
└── .gitignore
```

## Directory Purposes

**backend/cmd/server/**
- Purpose: Application entry point
- Contains: `main.go` - database setup, service wiring, router initialization
- Key files: `main.go` (~193 lines) - wires entire backend system

**backend/internal/core/domain/**
- Purpose: Pure domain models, business entities
- Contains: Content, User, Perspective, Pagination models
- Constraint: Zero external imports except stdlib
- Key files:
  - `content.go` - Content entity (URL, type, metadata)
  - `user.go` - User entity (username, email, role)
  - `perspective.go` - Perspective entity (quality, agreement ratings)
  - `errors.go` - Domain error types (ErrNotFound, ErrValidation, etc.)

**backend/internal/core/ports/**
- Purpose: Interface definitions for hexagonal architecture
- Contains: Repository and Service interfaces
- Pattern: Defined here, implemented in adapters/
- Key files:
  - `repositories/content_repository.go` - Port for content persistence
  - `repositories/user_repository.go` - Port for user persistence
  - `repositories/perspective_repository.go` - Port for perspective persistence
  - `services/youtube_client.go` - Port for YouTube integration

**backend/internal/core/services/**
- Purpose: Business logic and orchestration
- Contains: Service implementations
- Pattern: Injected dependencies (repositories, external clients)
- Key files:
  - `content_service.go` - Create/list/filter content, YouTube integration
  - `user_service.go` - Create/list/filter users, statistics
  - `perspective_service.go` - Create/update perspectives, filtering

**backend/internal/adapters/graphql/**
- Purpose: PRIMARY adapter - GraphQL API endpoint
- Contains: Resolver implementations, generated code
- Pattern: Implements generated resolver interfaces
- Key files:
  - `generated/generated.go` - Auto-generated by gqlgen
  - `resolvers/*.go` - Resolver implementations
  - `model/models_gen.go` - GraphQL type definitions

**backend/internal/adapters/repositories/postgres/**
- Purpose: SECONDARY adapter - PostgreSQL persistence
- Contains: GORM models and repository implementations
- Pattern: Convert domain ↔ GORM models, implement repository ports
- Key files:
  - `gorm_models.go` - GORM model structs with database tags
  - `gorm_mappers.go` - Bidirectional conversion functions
  - `gorm_content_repository.go` - ContentRepository implementation
  - `gorm_user_repository.go` - UserRepository implementation
  - `gorm_perspective_repository.go` - PerspectiveRepository implementation
  - `helpers.go` - Cursor encoding, sort mapping, type converters

**backend/internal/adapters/youtube/**
- Purpose: SECONDARY adapter - YouTube API integration
- Contains: YouTube API client and response parsing
- Pattern: Implements YouTubeClient port interface
- Key files:
  - `client.go` - HTTP client for YouTube API v3
  - `parser.go` - Extract video ID from URLs, transform API responses

**backend/pkg/database/**
- Purpose: Database connection and pooling
- Contains: GORM setup, connection utilities
- Key files:
  - `postgres.go` - ConnectGORM(), PingGORM() functions
  - `postgres_test.go` - Integration tests
  - `stats.go` - DB connection pool stats
  - `plugins.go` - Slow query logger

**backend/pkg/middleware/**
- Purpose: HTTP request/response processing
- Key files:
  - `timing.go` - Request timing and latency logging
  - `recovery.go` - Panic recovery with JSON error responses

**backend/pkg/graphql/**
- Purpose: GraphQL utilities and extensions
- Key files:
  - `intid.go` - Custom scalar for integers in GraphQL
  - `timing.go` - GraphQL operation performance logging

**backend/migrations/**
- Purpose: Database schema evolution
- Format: `{sequence}_{name}.{up|down}.sql`
- Current migrations:
  - 000001_create_content.sql - Content table
  - 000002_update_response_jsonb.sql - JSONB for YouTube responses
  - 000003_update_length_numeric.sql - Duration field
  - 000004_add_perspectives_users.sql - Perspectives and users tables
  - 000005_add_user_timestamps.sql - Timestamps
  - 000006_user_mutations_sentinel.sql - Sentinel user system
  - 000007_remove_claim_add_system_user.sql - Remove claim, add system user
  - 000008_add_user_active.sql - Active status
  - 000009_add_user_role.sql - User roles (ADMIN, SENTINEL, DEFAULT)
  - 000010_drop_content_unique_name.sql - Allow duplicate names

**frontend/src/routes/**
- Purpose: SvelteKit file-based routing
- Pattern: File = URL route
- Key files:
  - `+layout.svelte` - Root layout wraps all pages with QueryClientProvider, Header, Toaster
  - `+page.svelte` - Home page with ActivityTable and controls

**frontend/src/lib/components/**
- Purpose: Reusable Svelte 5 components
- Pattern: Component per .svelte file
- Core components:
  - `ActivityTable.svelte` - Main data grid (AG Grid wrapper, ~400 lines)
  - `AddVideoPopover.svelte` - Quick video add interface
  - `UserSelector.svelte` - User selection dropdown
  - `Header.svelte` - Navigation and branding

**frontend/src/lib/queries/**
- Purpose: GraphQL queries and TanStack Query integration
- Pattern: gql-tagged query strings, TanStack Query hooks
- Key files:
  - `client.ts` - GraphQLClient with endpoint configuration
  - `content.ts` - Content queries (list, get, by URL)
  - `users.ts` - User queries (list, get, by username)
  - `keys.ts` - TanStack Query key factories
  - `hooks/useAddVideo.ts` - Add video mutation wrapper
  - `hooks/useCreateUser.ts` - Create user mutation wrapper

**frontend/src/lib/stores/**
- Purpose: Svelte state management (non-query state)
- Pattern: .svelte.ts files with $state() runes
- Key files:
  - `userSelection.svelte.ts` - Selected user context

**frontend/tests/**
- Purpose: Vitest test files
- Pattern: Mirrors src/ structure
- Key areas:
  - `unit/` - Pure function tests
  - `components/` - Svelte component tests
  - `fixtures/` - Test data
  - `helpers/` - Test utilities

## Key File Locations

**Entry Points:**
- Backend: `backend/cmd/server/main.go` - HTTP server startup, dependency injection
- Frontend: `frontend/src/routes/+layout.svelte` - Root layout, app wrapper
- Build (backend): `backend/Makefile` - Build targets and development commands
- Build (frontend): `frontend/package.json` - npm scripts for build/dev

**Configuration:**
- Backend config: `backend/internal/config/config.go` - Load and validate
- Frontend config: `frontend/vite.config.ts` - Vite/SvelteKit build setup
- Database: `backend/migrations/` - Schema definitions
- GraphQL: `backend/schema.graphql` - API schema, `backend/gqlgen.yml` - codegen config

**Core Logic:**
- Business logic: `backend/internal/core/services/` - Service implementations
- Data access: `backend/internal/adapters/repositories/postgres/` - Query execution
- API layer: `backend/internal/adapters/graphql/resolvers/` - Endpoint handlers
- UI: `frontend/src/lib/components/ActivityTable.svelte` - Main interface

**Testing:**
- Backend tests: `backend/test/` - Test mocks and fixtures
- Frontend tests: `frontend/tests/` - Vitest test files
- Integrations: `backend/test/database/` - DB connection tests

## Naming Conventions

**Files:**
- Domain models: singular lowercase, no prefix (`content.go`, `user.go`)
- Services: singular + `_service.go` (`content_service.go`)
- Repositories: singular + `_repository.go` or `gorm_` prefix (`gorm_content_repository.go`)
- Resolvers: singular + `_resolver.go` or `resolvers_gen.go` (gqlgen convention)
- Tests: `{file}_test.go` (Go), `{component}.test.ts` or `{file}.spec.ts` (Svelte)
- Migrations: `{sequence}_{description}.{up|down}.sql`
- Components: PascalCase (`.svelte` files)

**Directories:**
- Go packages: lowercase, no underscores (`internal/adapters/repositories/postgres`)
- Frontend components: lowercase, dash-separated (`shadcn`, `queries`, `stores`)
- Routes: kebab-case or reserved SvelteKit syntax (`+layout.svelte`, `+page.svelte`)
- Features: by domain (content, user, perspective)

**Functions/Methods:**
- Go: camelCase starting lowercase (`GetByID`, `CreateFromYouTube`)
- Svelte: camelCase or PascalCase (components are PascalCase)
- GraphQL: camelCase queries and mutations (`listContent`, `createContent`)

**Variables:**
- Go: camelCase (`userID`, `contentList`)
- TypeScript: camelCase (`selectedUser`, `isLoading`)
- GraphQL: camelCase (`createdAt`, `userId`)

## Where to Add New Code

**New Feature (Backend):**
1. Domain model: `backend/internal/core/domain/{feature}.go`
2. Repository port: `backend/internal/core/ports/repositories/{feature}_repository.go`
3. Service: `backend/internal/core/services/{feature}_service.go`
4. Repository impl: `backend/internal/adapters/repositories/postgres/gorm_{feature}_repository.go`
5. GraphQL schema: Add type and mutations to `backend/schema.graphql`
6. Generate: `make graphql-gen`
7. Resolvers: Implement in `backend/internal/adapters/graphql/resolvers/{feature}_resolver.go`
8. Wire: Add to `backend/cmd/server/main.go`
9. Tests: Add to `backend/test/{service|resolver}/{feature}_test.go`

**New Component/UI (Frontend):**
1. Component: `frontend/src/lib/components/{FeatureName}.svelte`
2. Queries (if needed): `frontend/src/lib/queries/{feature}.ts`
3. Hooks (if needed): `frontend/src/lib/queries/hooks/use{Feature}.ts`
4. Tests: `frontend/tests/components/{FeatureName}.test.ts`
5. Styles: Use Tailwind utilities in component `class=` attributes
6. Types: Add interfaces to query file or separate `types/` directory

**New Utility:**
- Shared helpers: `backend/internal/adapters/repositories/postgres/helpers.go` (DB helpers) or `backend/pkg/` (general utilities)
- Frontend helpers: `frontend/src/lib/utils/` (utility functions)

**New Endpoint:**
- Existing service: Add mutation/query to `backend/schema.graphql`, regenerate, implement resolver
- New service flow: Follow "New Feature" process above

## Special Directories

**backend/migrations/**
- Purpose: Version-controlled database schema
- Generated: No (hand-written)
- Committed: Yes
- Tool: golang-migrate
- Run: `make migrate-up` (apply), `make migrate-down` (rollback)
- New: `make migrate-create` prompts for name, creates .up.sql and .down.sql files

**backend/.svelte-kit/**
- Purpose: SvelteKit generated files (routes, types, etc.)
- Generated: Yes (auto-generated by SvelteKit)
- Committed: No (.gitignored)
- Regenerate: `pnpm run prepare` or automatic during dev

**frontend/.svelte-kit/**
- Purpose: SvelteKit build artifacts and generated types
- Generated: Yes
- Committed: No
- Regenerate: Automatic, or manual `npx svelte-kit sync`

**backend/tmp/**
- Purpose: Temporary build artifacts
- Generated: Yes
- Committed: No

**frontend/build/**
- Purpose: Production frontend build output
- Generated: Yes (by `vite build`)
- Committed: No
- Location: Static files ready for deployment

**frontend/coverage/**
- Purpose: Test coverage reports
- Generated: Yes (by `vitest run --coverage`)
- Committed: No (typically)
- View: `coverage/index.html` after `pnpm run test:coverage`

---

*Structure analysis: 2026-02-16*
