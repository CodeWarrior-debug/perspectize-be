# Architecture

**Analysis Date:** 2026-02-07

## Monorepo Overview

Perspectize is a **two-stack monorepo** with active development on both frontend and backend:

```
perspectize-be/ (repository root)
├── perspectize-go/     # Go GraphQL API backend (ACTIVE)
├── perspectize-fe/     # SvelteKit web app frontend (ACTIVE)
└── perspectize-be/     # Legacy C# (DEPRECATED - ignore)
```

**Key Principle:** Frontend and backend are loosely coupled via GraphQL HTTP API. No shared code between stacks. Stacks are independently deployable.

---

## Backend Architecture (perspectize-go/)

### Pattern: Hexagonal Architecture (Ports & Adapters)

The Go backend implements **Hexagonal Architecture** with strict dependency inversion:

```
┌─────────────────────────────────────────────────────┐
│         PRIMARY ADAPTERS (Driving/Input)            │
│    GraphQL Resolvers (gqlgen schema-first)          │
│    HTTP Handlers (net/http + CORS middleware)       │
└──────────────┬──────────────────────────────────────┘
               │ depends on
┌──────────────▼──────────────────────────────────────┐
│         CORE DOMAIN LAYER (Domain Logic)            │
│  ├── domain/       (Entities, enums, errors)        │
│  ├── services/     (Business logic, orchestration)  │
│  └── ports/        (Interface contracts)            │
│                                                      │
│  Critical Rule: Dependencies point INWARD ONLY      │
│  Domain never depends on adapters or frameworks     │
└──────────────┬──────────────────────────────────────┘
               │ implements
┌──────────────▼──────────────────────────────────────┐
│      SECONDARY ADAPTERS (Driven/Output)             │
│  ├── repositories/postgres/  (Database access)      │
│  └── youtube/                (External API client)  │
└─────────────────────────────────────────────────────┘
```

**Critical:** Services depend on port interfaces, not concrete adapters. This enables testability and clean separation.

### Layers

**Domain Layer** (`internal/core/domain/`):
- **Purpose:** Pure business models and rules, zero external dependencies
- **Location:** `internal/core/domain/`
- **Contains:**
  - `content.go` — Content entity (name, URL, type, length, JSONB response)
  - `perspective.go` — Perspective entity (ratings: quality/agreement/importance/confidence, privacy, category, labels, categorizedRatings)
  - `user.go` — User entity (username, email)
  - `pagination.go` — Pagination types (ContentListParams, PaginatedContent, cursors)
  - `errors.go` — Domain error constants (ErrNotFound, ErrAlreadyExists, ErrInvalidInput, ErrInvalidURL, ErrYouTubeAPI)
- **Depends on:** Standard library only (time, encoding/json)
- **Used by:** Services, repositories map domain models

**Ports Layer** (`internal/core/ports/`):
- **Purpose:** Define contracts (interfaces) between core and adapters
- **Location:** `internal/core/ports/repositories/`, `internal/core/ports/services/`
- **Contains:**
  - `repositories/content_repository.go` — ContentRepository interface (Create, GetByID, GetByURL, List)
  - `repositories/user_repository.go` — UserRepository interface (Create, GetByID, GetByUsername, GetByEmail)
  - `repositories/perspective_repository.go` — PerspectiveRepository interface (Create, GetByID, GetByUserAndClaim, Update, Delete, List)
  - `services/youtube_client.go` — YouTubeClient interface (GetVideoMetadata)
- **Depends on:** Domain types only
- **Key Rule:** Core depends on these, adapters implement these

**Services Layer** (`internal/core/services/`):
- **Purpose:** Business logic orchestration, validation, error handling
- **Location:** `internal/core/services/`
- **Contains:**
  - `content_service.go` — ContentService (CreateFromYouTube, GetByID, ListContent)
  - `user_service.go` — UserService (Create, GetByID, GetByUsername)
  - `perspective_service.go` — PerspectiveService (Create, GetByID, Update, Delete, ListPerspectives)
- **Depends on:** Domain models + port interfaces (NOT concrete adapters)
- **Pattern:** Constructor injection: `NewContentService(repo ContentRepository, yt YouTubeClient) *ContentService`
- **Error handling:** Explicit error returns with fmt.Errorf("%w: context", domainError)

**GraphQL Adapter** (`internal/adapters/graphql/`):
- **Purpose:** PRIMARY adapter — HTTP GraphQL API entry point
- **Location:** `internal/adapters/graphql/`
- **Contains:**
  - `resolvers/resolver.go` — Dependency injection container (holds service references)
  - `resolvers/schema.resolvers.go` — Query/Mutation resolver implementations
  - `generated/generated.go` — Auto-generated gqlgen code (DO NOT EDIT)
  - `model/` — Auto-generated GraphQL type definitions
- **Workflow:** Edit `schema.graphql` → `make graphql-gen` → Implement resolvers in schema.resolvers.go
- **Key pattern:** Resolver calls service, maps domain response to GraphQL model

**Repository Adapters** (`internal/adapters/repositories/postgres/`):
- **Purpose:** SECONDARY adapter — PostgreSQL database access
- **Location:** `internal/adapters/repositories/postgres/`
- **Contains:**
  - `content_repository.go` — ContentRepository implementation (cursor pagination, keyset queries)
  - `user_repository.go` — UserRepository implementation (ID/username/email lookups)
  - `perspective_repository.go` — PerspectiveRepository implementation (complex JSONB handling)
- **Pattern:** Each implements corresponding port interface
- **Key techniques:**
  - Cursor-based pagination: `WHERE id > cursor_id ORDER BY ... LIMIT limit+1`
  - Type conversion: `contentTypeToDBValue()`, `contentTypeFromDBValue()`
  - Null handling: sql.NullString, sql.NullInt64
  - Array handling: pq.StringArray, custom JSONBArray for JSONB arrays

**YouTube Adapter** (`internal/adapters/youtube/`):
- **Purpose:** SECONDARY adapter — External API client
- **Location:** `internal/adapters/youtube/`
- **Contains:**
  - `client.go` — YouTube API HTTP client (GetVideoMetadata)
  - `parser.go` — Utility functions (ExtractVideoID regex, ParseISO8601Duration)
- **Implements:** YouTubeClient port interface
- **Pattern:** HTTP client abstraction, no direct calls from services

**Configuration** (`internal/config/`):
- **Purpose:** Load and manage configuration
- **Location:** `internal/config/config.go`
- **Behavior:** Load `config/config.example.json` + environment variable overrides
- **Critical variables:** DATABASE_URL, YOUTUBE_API_KEY

## Frontend Architecture (perspectize-fe/)

### Pattern: SvelteKit File-Based Routing + TanStack Query

**Entry Point:** `src/routes/+layout.svelte` (root) → `src/routes/+page.svelte` (home)

```
Routes (SvelteKit filesystem routing):
├── +layout.svelte          Root layout
├── +layout.ts              Prerender config
└── +page.svelte            Home page (activity feed)

Library Structure:
├── components/             Reusable Svelte 5 components
│   ├── Header.svelte
│   ├── ActivityTable.svelte (AG Grid content table)
│   ├── UserSelector.svelte
│   ├── PageWrapper.svelte
│   └── shadcn/             shadcn-svelte UI primitives
│
├── queries/                TanStack Query + GraphQL
│   ├── client.ts           GraphQLClient singleton
│   ├── content.ts          Content queries (gql)
│   └── users.ts            User queries (gql)
│
├── stores/                 Svelte runes reactive state
│   └── userSelection.svelte.ts
│
└── utils/                  Helper functions
    └── utils.ts
```

### Data Flow: TanStack Query Pattern

**1. GraphQL Client Setup** (`lib/queries/client.ts`):
```typescript
const graphqlClient = new GraphQLClient(
  import.meta.env.VITE_GRAPHQL_URL || 'http://localhost:8080/graphql'
);
```

**2. Query Definition** (`lib/queries/content.ts`):
```typescript
export const LIST_CONTENT = gql`
  query Content($first: Int, $sortBy: ContentSortBy, $sortOrder: SortOrder) {
    content(first: $first, sortBy: $sortBy, sortOrder: $sortOrder) {
      items { id name url contentType length lengthUnits createdAt updatedAt }
      pageInfo { hasNextPage endCursor }
      totalCount
    }
  }
`;
```

**3. Component Query Usage** (Svelte 5 with function wrapper pattern):
```svelte
const query = createQuery(() => ({
  queryKey: ['content', { first: 100, sortBy: 'UPDATED_AT', sortOrder: 'DESC' }],
  queryFn: () => graphqlClient.request(LIST_CONTENT, { ... }),
  staleTime: 60 * 1000,
  retry: 1,
  enabled: browser  // Browser-only (prevents SSR)
}));
```

**4. Reactive Access** (NO `$` prefix — reactive objects, not stores):
```svelte
{#if query.isLoading}Loading...{/if}
{#if query.error}Error: {query.error.message}{/if}
{#if query.data}Display: {query.data.content.items}{/if}
```

### Svelte 5 Runes (REQUIRED)

This codebase uses **Svelte 5 runes exclusively**. Do NOT use Svelte 4 syntax:

| Svelte 5 (USE THIS) | Svelte 4 (DON'T) | Purpose |
|---|---|---|
| `let count = $state(0)` | `let count = 0` + `$:` | Reactive state |
| `let doubled = $derived(count * 2)` | `$: doubled = count * 2` | Derived values |
| `let { prop } = $props()` | `export let prop` | Props |
| `$effect(() => { ... })` | `onMount()` / `$:` side effects | Lifecycle/effects |
| `{@render children()}` | `<slot />` | Render children |
| `onclick={handler}` | `on:click={handler}` | Event handlers |

**Component Structure** (from ActivityTable.svelte):
```svelte
<script lang="ts">
  let { rowData = [], loading = false, searchText = '' } = $props<{
    rowData: ContentRow[];
    loading?: boolean;
    searchText?: string;
  }>();

  let gridApi = $state<GridApi | null>(null);

  let formattedDuration = $derived(rowData.map(row => formatDuration(row.length, row.lengthUnits)));

  $effect(() => {
    if (gridApi) {
      gridApi.setGridOption('quickFilterText', searchText);
    }
  });
</script>
```

### AG Grid v32.x Integration (CRITICAL)

The frontend uses **ag-grid-svelte5** wrapper which bundles AG Grid v32.x internally:

```svelte
<script lang="ts">
  import AgGridSvelte5Component from 'ag-grid-svelte5';
  import { ClientSideRowModelModule } from '@ag-grid-community/client-side-row-model';
  import { themeQuartz } from '@ag-grid-community/theming';
  import type { GridOptions } from '@ag-grid-community/core';

  const theme = themeQuartz.withParams({ fontFamily: 'Inter, sans-serif' });
  const gridOptions: GridOptions<ContentRow> = {
    columnDefs: [
      { field: 'name', headerName: 'Title', flex: 2, sortable: true },
      { field: 'contentType', headerName: 'Type', width: 100 },
      // ...
    ],
    pagination: true,
    paginationPageSize: 10,
    paginationPageSizeSelector: [10, 25, 50]
  };
</script>

<AgGridSvelte5Component {gridOptions} {rowData} {theme} modules={[ClientSideRowModelModule]} />
```

**Critical:**
- Do NOT import from `ag-grid-community` directly (use `@ag-grid-community/*`)
- Do NOT import AG Grid CSS (use `themeQuartz.withParams()`)
- Do use `AgGridSvelte5Component` (not `AgGridSvelte`)

## End-to-End Data Flow

### Fetching Content List

```
Svelte Component (+page.svelte)
    |
    | createQuery(() => LIST_CONTENT query)
    v
graphqlClient.request() [graphql-request library]
    |
    | HTTP POST /graphql with GraphQL query
    v
Go Backend (gqlgen handler)
    |
    | Parse GraphQL query
    | Route to Query.content() resolver
    v
GraphQL Resolver (schema.resolvers.go)
    |
    | contentService.ListContent(ctx, params)
    v
Service (ContentService)
    |
    | Validate pagination (1-100)
    | contentRepo.List(ctx, params)
    v
Repository (postgres/ContentRepository)
    |
    | Build cursor-keyset SQL: SELECT * FROM content WHERE id > cursor ORDER BY ...
    | sqlx.DB.Select() executes
    v
PostgreSQL 17
    |
    | SELECT * FROM content ORDER BY updated_at DESC LIMIT 101
    |
    +--- Map rows to domain.Content objects ---+
                                               |
                        PaginatedContent { items, pageInfo }
                                               |
                        Resolver maps to GraphQL types
                                               |
                        HTTP JSON response
                                               |
                        graphql-request deserializes
                                               |
                        TanStack Query caches (staleTime: 60s)
                                               |
                        Svelte reactive object updates
                                               |
                        Component re-renders
                                               |
                        ActivityTable receives rowData
                                               |
                        AG Grid renders with columns
```

## Key Abstractions

### Content Entity
- **Purpose:** Represents media items (YouTube videos initially)
- **Domain path:** `internal/core/domain/content.go`
- **Fields:** id, name, url, contentType (enum: YOUTUBE), length (seconds), lengthUnits, response (JSONB YouTube metadata), createdAt, updatedAt

### Perspective Entity
- **Purpose:** Multi-dimensional user rating of content
- **Domain path:** `internal/core/domain/perspective.go`
- **Rating fields:** quality (0-10000), agreement (0-10000), importance (0-10000), confidence (0-10000)
- **Metadata:** claim (string), privacy (PUBLIC/PRIVATE), reviewStatus (PENDING/APPROVED/REJECTED), category, labels (array), parts (array), categorizedRatings (JSONB array)

### Pagination Abstraction
- **Type:** Cursor-based (keyset) not offset-based
- **Format:** Opaque base64-encoded cursor: `cursor:<id>`
- **SQL Pattern:** `WHERE id > :cursor AND ... ORDER BY created_at DESC LIMIT limit+1`
- **Frontend:** TanStack Query handles caching and invalidation

## Entry Points

### Backend

**Server Startup** (`cmd/server/main.go`):
1. Load .env file (optional via godotenv)
2. Load config from `config/config.example.json`, override with env vars
3. Connect to PostgreSQL via DATABASE_URL or local config
4. Test connection: Ping and version query
5. Initialize adapters: YouTube client, repositories with db connection
6. Initialize services: Inject adapter dependencies
7. Create GraphQL resolver: Inject services
8. Register HTTP routes:
   - `GET /` → GraphQL Playground (gqlgen built-in)
   - `POST /graphql` → GraphQL execution
   - CORS middleware allows all origins (`*`)
9. Listen on port 8080 (configurable)

**GraphQL Endpoint:** `POST http://localhost:8080/graphql`

### Frontend

**Root Layout** (`src/routes/+layout.svelte`):
- Initializes QueryClient with browser-only queries
- Wraps app with QueryClientProvider
- Renders Toaster (top-right, 2s auto-dismiss)
- Renders Header component
- Renders {@render children()} for page routes

**Home Page** (`src/routes/+page.svelte`):
- Fetches LIST_CONTENT query with sorting/pagination
- Displays search input
- Renders ActivityTable with AG Grid

## Error Handling

### Backend Strategy
- Explicit error returns, no exceptions
- Domain error types in `internal/core/domain/errors.go`
- Error wrapping: `fmt.Errorf("%w: context", domain.ErrXxx)`
- Service layer validates, maps to domain errors
- Repository layer maps sql.ErrNoRows → domain.ErrNotFound
- Resolver layer returns GraphQL errors

### Frontend Strategy
- TanStack Query error state: `query.error`
- Conditional rendering: `{#if query.error} Show error {/if}`
- Toast notifications for user feedback via svelte-sonner
- Network errors passed through graphql-request

## Cross-Cutting Concerns

**Logging (Backend):**
- Framework: `log/slog` (standard library)
- Current: Used in main.go for startup diagnostics
- Future: Add structured logging to services

**Validation (Backend):**
- Input validation in service layer before DB operations
- Example: Pagination bounds (1-100), claim length, rating ranges (0-10000)
- GraphQL input types validated by gqlgen
- Database constraints enforce NOT NULL, UNIQUE, CHECK

**Authentication:**
- Currently: None (public API for MVP development)
- Planned: OAuth 2.0 with YouTube, JWT tokens, refresh token rotation

**CORS (Backend):**
- Middleware in main.go: Allows all origins (`*`)
- Future Phase 5: Restrict to frontend production origin

---

*Architecture analysis: 2026-02-07*
