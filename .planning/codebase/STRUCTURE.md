# Codebase Structure

**Analysis Date:** 2026-02-07

## Directory Layout

```
perspectize/ (repository root)
├── backend/                 # Go GraphQL backend (ACTIVE)
│   ├── cmd/server/
│   │   └── main.go                 # Server entry point, dependency wiring
│   │
│   ├── internal/
│   │   ├── core/                   # Hexagonal domain layer
│   │   │   ├── domain/             # Pure domain models, enums, errors
│   │   │   │   ├── content.go
│   │   │   │   ├── perspective.go
│   │   │   │   ├── user.go
│   │   │   │   ├── pagination.go
│   │   │   │   └── errors.go
│   │   │   │
│   │   │   ├── ports/              # Interface contracts
│   │   │   │   ├── repositories/
│   │   │   │   │   ├── content_repository.go
│   │   │   │   │   ├── perspective_repository.go
│   │   │   │   │   └── user_repository.go
│   │   │   │   └── services/
│   │   │   │       └── youtube_client.go
│   │   │   │
│   │   │   └── services/           # Business logic orchestration
│   │   │       ├── content_service.go
│   │   │       ├── perspective_service.go
│   │   │       └── user_service.go
│   │   │
│   │   ├── adapters/               # Infrastructure implementations
│   │   │   ├── graphql/            # PRIMARY: GraphQL HTTP API
│   │   │   │   ├── resolvers/
│   │   │   │   │   └── resolver.go (DI container + wiring)
│   │   │   │   ├── generated/      # Auto-generated (do NOT edit)
│   │   │   │   └── model/          # Auto-generated GraphQL types
│   │   │   │
│   │   │   ├── repositories/       # SECONDARY: Database access
│   │   │   │   └── postgres/
│   │   │   │       ├── content_repository.go
│   │   │   │       ├── perspective_repository.go
│   │   │   │       └── user_repository.go
│   │   │   │
│   │   │   └── youtube/            # SECONDARY: YouTube API client
│   │   │       ├── client.go
│   │   │       └── parser.go
│   │   │
│   │   ├── config/                 # Configuration loading
│   │   │   └── config.go
│   │   │
│   │   └── middleware/             # HTTP middleware (empty, future use)
│   │
│   ├── pkg/                        # Public shared packages
│   │   ├── database/               # DB connection/pooling
│   │   │   └── postgres.go
│   │   └── graphql/                # Custom scalars
│   │       └── intid.go
│   │
│   ├── test/                       # Test files (non-standard location)
│   │   ├── config/
│   │   ├── database/
│   │   ├── domain/
│   │   ├── resolvers/
│   │   ├── services/
│   │   └── youtube/
│   │
│   ├── migrations/                 # SQL migration files
│   │   ├── 000001_create_content.up.sql
│   │   ├── 000001_create_content.down.sql
│   │   └── ... (numbered sequence)
│   │
│   ├── config/                     # Configuration files
│   │   └── config.example.json
│   │
│   ├── schema.graphql              # GraphQL schema (schema-first)
│   ├── gqlgen.yml                  # gqlgen configuration
│   ├── Makefile                    # Build/test commands
│   ├── go.mod / go.sum             # Go dependencies
│   ├── .env.example                # Environment template
│   ├── .gitignore
│   ├── docker-compose.yml          # Local PostgreSQL setup
│   ├── CLAUDE.md                   # Backend instructions
│   └── README.md
│
│
├── frontend/                 # SvelteKit frontend (ACTIVE)
│   ├── src/
│   │   ├── routes/                 # SvelteKit file-based routing
│   │   │   ├── +layout.svelte      # Root layout (QueryClientProvider, Header, Toaster)
│   │   │   ├── +layout.ts          # Prerender config
│   │   │   └── +page.svelte        # Home page (activity feed with AG Grid)
│   │   │
│   │   ├── lib/
│   │   │   ├── components/         # Reusable Svelte 5 components
│   │   │   │   ├── Header.svelte
│   │   │   │   ├── ActivityTable.svelte (AG Grid content table)
│   │   │   │   ├── UserSelector.svelte
│   │   │   │   ├── PageWrapper.svelte
│   │   │   │   ├── AGGridTest.svelte
│   │   │   │   └── shadcn/         # shadcn-svelte UI primitives
│   │   │   │       └── button/
│   │   │   │
│   │   │   ├── queries/            # TanStack Query + GraphQL definitions
│   │   │   │   ├── client.ts       # GraphQLClient singleton (VITE_GRAPHQL_URL)
│   │   │   │   ├── content.ts      # Content queries (gql templates)
│   │   │   │   └── users.ts        # User queries (gql templates)
│   │   │   │
│   │   │   ├── stores/             # Svelte runes reactive state
│   │   │   │   └── userSelection.svelte.ts
│   │   │   │
│   │   │   ├── utils/              # Utility functions
│   │   │   │   └── utils.ts
│   │   │   │
│   │   │   ├── assets/             # Static bundled assets
│   │   │   │   └── favicon.svg
│   │   │   │
│   │   │   └── index.ts            # Barrel export
│   │   │
│   │   ├── app.css                 # Global styles (Tailwind v4 + custom theme)
│   │   └── app.html                # HTML shell
│   │
│   ├── static/                     # Public static files (not bundled)
│   │
│   ├── tests/                      # Vitest test files
│   │
│   ├── build/                      # Build output (generated)
│   │
│   ├── coverage/                   # Code coverage reports (generated)
│   │
│   ├── svelte.config.js            # SvelteKit configuration
│   ├── vite.config.ts              # Vite configuration
│   ├── tsconfig.json               # TypeScript configuration
│   ├── tailwind.config.ts          # Tailwind v4 theme config
│   ├── components.json             # shadcn-svelte configuration
│   ├── package.json                # npm dependencies
│   ├── pnpm-lock.yaml              # pnpm lockfile
│   ├── .npmrc                      # npm registry config
│   ├── .gitignore
│   ├── CLAUDE.md                   # Frontend instructions
│   └── README.md
│
│
├── perspectize/                 # Legacy C# (DEPRECATED - ignore)
│   ├── Services/
│   ├── controllers/
│   └── ... (legacy code)
│
│
├── .docs/                          # Shared documentation
│   ├── ARCHITECTURE.md
│   ├── LOCAL_DEVELOPMENT.md
│   ├── AGENTS.md
│   ├── DOMAIN_GUIDE.md
│   ├── GO_PATTERNS.md
│   ├── GITHUB_PROJECTS.md
│   └── VERIFICATION.md
│
├── .planning/                      # GSD workflow artifacts
│   ├── codebase/                   # Codebase analysis documents
│   │   ├── ARCHITECTURE.md         # This monorepo's patterns
│   │   ├── STRUCTURE.md            # This file
│   │   ├── CONVENTIONS.md          # Code style rules
│   │   ├── TESTING.md              # Testing patterns
│   │   ├── STACK.md                # Technology stack
│   │   ├── INTEGRATIONS.md         # External API integrations
│   │   └── CONCERNS.md             # Technical debt and issues
│   │
│   └── phases/                     # GSD phase planning
│       └── {phase-id}/
│           └── {plan}-PLAN.md
│
├── .git/                           # Git repository
├── .github/                        # GitHub workflows
├── CLAUDE.md                       # Root project instructions (YOU ARE HERE)
├── README.md                       # Root README
└── .gitignore
```

## Directory Purposes

### Backend (backend/)

**`cmd/server/`**
- Entry point for the application
- `main.go`: Loads config, initializes database, wires dependencies, starts HTTP server
- Startup sequence: config → DB connect → adapters → services → resolver wiring → listen

**`internal/core/domain/`**
- Pure domain models with zero framework dependencies
- `content.go` — Content entity with JSONB response field
- `perspective.go` — Perspective entity with 4-dimensional rating system
- `user.go` — User entity
- `pagination.go` — Pagination types and cursor abstractions
- `errors.go` — Domain-specific error constants

**`internal/core/ports/{repositories,services}/`**
- Interface contracts between core and adapters
- Repository ports: ContentRepository, UserRepository, PerspectiveRepository
- Service ports: YouTubeClient interface

**`internal/core/services/`**
- Business logic implementations using dependency injection
- `content_service.go` — ContentService (CreateFromYouTube, GetByID, ListContent)
- `user_service.go` — UserService (Create, GetByID, GetByUsername)
- `perspective_service.go` — PerspectiveService (Create, GetByID, Update, Delete, ListPerspectives)

**`internal/adapters/graphql/`**
- PRIMARY adapter: GraphQL HTTP API entry point
- `resolvers/resolver.go` — Dependency injection container
- `generated/` — Auto-generated by gqlgen (never edit)
- `model/` — Auto-generated GraphQL type definitions

**`internal/adapters/repositories/postgres/`**
- SECONDARY adapter: PostgreSQL database implementations
- Implements port interfaces using sqlx
- Cursor-based pagination pattern, JSONB handling

**`internal/adapters/youtube/`**
- SECONDARY adapter: YouTube Data API client
- Implements YouTubeClient port interface

**`migrations/`**
- SQL migration files (numbered sequence)
- Executed via golang-migrate: `make migrate-up`

**`config/`**
- Configuration files
- `config.example.json` — Configuration template

### Frontend (frontend/)

**`src/routes/`**
- SvelteKit file-based routing
- `+layout.svelte` — Root layout (QueryClientProvider, Header, Toaster wrapper)
- `+page.svelte` — Home page (activity feed with search and AG Grid table)

**`src/lib/components/`**
- Reusable Svelte 5 components using runes
- `Header.svelte` — Top navigation bar
- `ActivityTable.svelte` — AG Grid content table with sorting, pagination, search
- `UserSelector.svelte` — User filter dropdown
- `PageWrapper.svelte` — Page container with consistent styling
- `shadcn/` — shadcn-svelte UI primitives (button)

**`src/lib/queries/`**
- TanStack Query + GraphQL integration
- `client.ts` — GraphQLClient singleton configured with VITE_GRAPHQL_URL
- `content.ts` — Content query definitions using gql tagged templates
- `users.ts` — User query definitions

**`src/lib/stores/`**
- Svelte 5 runes reactive state modules
- `userSelection.svelte.ts` — User filter selection state

**`src/lib/utils/`**
- Helper/utility functions
- `utils.ts` — Date formatting, string utilities, etc.

**`src/lib/assets/`**
- Static assets bundled by Vite
- `favicon.svg` — App icon

**`tests/`**
- Vitest test files
- Structure mirrors src/ organization

## Key File Locations

### Backend Entry Points
- `backend/cmd/server/main.go` — Application startup
- `backend/schema.graphql` — GraphQL schema definition
- `backend/gqlgen.yml` — gqlgen code generation config

### Backend Configuration
- `backend/config/config.example.json` — Config template
- `backend/.env.example` — Environment variables template
- `backend/Makefile` — Build commands and targets

### Backend Domain Models
- `backend/internal/core/domain/content.go`
- `backend/internal/core/domain/perspective.go`
- `backend/internal/core/domain/user.go`
- `backend/internal/core/domain/pagination.go`
- `backend/internal/core/domain/errors.go`

### Backend Services
- `backend/internal/core/services/content_service.go`
- `backend/internal/core/services/perspective_service.go`
- `backend/internal/core/services/user_service.go`

### Backend Repositories
- `backend/internal/adapters/repositories/postgres/content_repository.go`
- `backend/internal/adapters/repositories/postgres/perspective_repository.go`
- `backend/internal/adapters/repositories/postgres/user_repository.go`

### Backend GraphQL
- `backend/schema.graphql` — Schema definition
- `backend/internal/adapters/graphql/resolvers/resolver.go` — Resolver DI container
- `backend/internal/adapters/graphql/generated/generated.go` — Auto-generated (do not edit)

### Frontend Entry Points
- `frontend/src/routes/+layout.svelte` — Root layout
- `frontend/src/routes/+page.svelte` — Home page
- `frontend/src/app.html` — HTML shell

### Frontend Configuration
- `frontend/svelte.config.js` — SvelteKit config
- `frontend/vite.config.ts` — Vite config
- `frontend/tailwind.config.ts` — Tailwind v4 theme
- `frontend/components.json` — shadcn-svelte config
- `frontend/tsconfig.json` — TypeScript config

### Frontend Components
- `frontend/src/lib/components/Header.svelte` — Top navigation
- `frontend/src/lib/components/ActivityTable.svelte` — Content grid
- `frontend/src/lib/components/UserSelector.svelte` — User filter
- `frontend/src/lib/components/PageWrapper.svelte` — Page wrapper

### Frontend Queries
- `frontend/src/lib/queries/client.ts` — GraphQL client
- `frontend/src/lib/queries/content.ts` — Content queries
- `frontend/src/lib/queries/users.ts` — User queries

## Naming Conventions

### Backend (Go)

**Files:**
- Domain models: `{entity}.go` (e.g., `content.go`, `user.go`)
- Services: `{entity}_service.go` (e.g., `content_service.go`)
- Repositories: `{entity}_repository.go` (e.g., `user_repository.go`)
- Tests: `{unit}_test.go` (e.g., `content_service_test.go`)
- Migrations: `{sequence}_{description}.{up|down}.sql`

**Go Identifiers:**
- Types: PascalCase (Content, ContentService, ContentRepository)
- Methods: PascalCase (GetByID, ListContent, Create)
- Unexported functions: camelCase (contentTypeToDBValue, rowToDomain)
- Variables: camelCase (contentID, pageSize, err)
- Constants: UPPERCASE for enums (ContentTypeYouTube, PrivacyPublic)

**GraphQL Schema:**
- Types: PascalCase (Content, Perspective, User)
- Enum values: UPPERCASE (YOUTUBE, PUBLIC, PENDING)
- Field names: camelCase (createdAt, viewCount)
- Input types: PascalCase + "Input" (CreateContentFromYouTubeInput)

### Frontend (Svelte/TypeScript)

**Files:**
- Components: PascalCase (ActivityTable.svelte, Header.svelte)
- Non-components: camelCase (client.ts, content.ts, utils.ts)
- Test files: `{name}.test.ts` or `{name}.spec.ts`

**TypeScript Identifiers:**
- Types: PascalCase (ContentItem, ContentResponse, GridOptions)
- Interfaces: PascalCase (GridOptions)
- Functions: camelCase (formatDuration, formatDate)
- Constants: UPPERCASE (GRAPHQL_ENDPOINT)
- Variables: camelCase (rowData, searchText, gridApi)

**Svelte Identifiers:**
- State: camelCase via `$state()` (let items = $state([]))
- Props: camelCase via `$props()` destructuring
- Event handlers: camelCase (onClick, onGridReady)
- Stores: camelCase (userSelection)

## Where to Add New Code

### Adding a Backend Feature

1. **Domain model:** `internal/core/domain/{entity}.go`
   - Define struct with optional fields as pointers
   - Add enums if needed

2. **Repository port:** `internal/core/ports/repositories/{entity}_repository.go`
   - Define interface with CRUD methods

3. **Service:** `internal/core/services/{entity}_service.go`
   - Implement business logic
   - Depend on ports (not concrete adapters)
   - Use error wrapping: `fmt.Errorf("%w: ...", domain.Err...)`

4. **Repository implementation:** `internal/adapters/repositories/postgres/{entity}_repository.go`
   - Implement port interface
   - Use sqlx for database access

5. **GraphQL schema:** `schema.graphql`
   - Add type, inputs, queries/mutations

6. **Generate GraphQL:** `make graphql-gen`
   - Regenerates resolver skeleton

7. **GraphQL resolvers:** Update `internal/adapters/graphql/resolvers/schema.resolvers.go`
   - Implement generated methods
   - Call services, map to GraphQL types

8. **Wire in main:** Update `cmd/server/main.go`
   - Create repository with db connection
   - Create service with repo dependency
   - Add service to resolver

9. **Tests:** Add in `test/{services,resolvers}/`

### Adding a Frontend Component

1. **New component:** `src/lib/components/YourComponent.svelte`
   - Use Svelte 5 runes exclusively
   - Props via `$props()`
   - Render children via `{@render children()}`

2. **New query:** `src/lib/queries/{entity}.ts`
   - Define gql tagged template
   - Export query constant

3. **Usage:** In page/component:
   ```svelte
   const query = createQuery(() => ({
     queryKey: ['entity'],
     queryFn: () => graphqlClient.request(QUERY)
   }));
   ```

4. **Tests:** `tests/{path}/{component}.test.ts`
   - Use Vitest
   - Mock TanStack Query if needed

### Adding a Migration

```bash
cd backend
make migrate-create
# Prompts for name, creates numbered files
# Edit .up.sql and .down.sql
make migrate-up      # Test forward
make migrate-down    # Test rollback
```

## Special Directories

**`.planning/codebase/`**
- Purpose: Codebase analysis documents (ARCHITECTURE.md, STRUCTURE.md, etc.)
- Committed: Yes
- Generated: No (manually created via /gsd:map-codebase)

**`backend/migrations/`**
- Purpose: Version-controlled SQL migrations
- Committed: Yes
- Generated: No (manually created)
- Tool: golang-migrate (make migrate-*)

**`backend/internal/adapters/graphql/generated/`**
- Purpose: gqlgen auto-generated code
- Committed: Yes
- Generated: Yes (make graphql-gen)
- Edit: Never (changes only via schema.graphql)

**`backend/internal/adapters/graphql/model/`**
- Purpose: GraphQL model types auto-generated from schema
- Committed: Yes
- Generated: Yes (make graphql-gen)
- Edit: Never

**`frontend/.svelte-kit/`**
- Purpose: SvelteKit build cache and generated types
- Committed: No (.gitignore)
- Generated: Yes

**`frontend/build/`**
- Purpose: Production build output
- Committed: No (.gitignore)
- Generated: Yes (pnpm run build)

**`frontend/coverage/`**
- Purpose: Test coverage reports
- Committed: No (.gitignore)
- Generated: Yes (pnpm run test:coverage)

**`frontend/node_modules/`**
- Purpose: npm dependencies
- Committed: No (.gitignore)
- Generated: Yes (pnpm install)

---

*Structure analysis: 2026-02-07*
