# Coding Conventions

**Analysis Date:** 2026-09-04

This is a monorepo with two distinct convention sets: **Go backend** (`backend/`) and **SvelteKit frontend** (`frontend/`). Conventions below are split by stack.

---

## Backend (Go)

### Naming Patterns

**Packages:**
- Plural, lowercase, no abbreviation: `package handlers`, `package services`, `package repositories`
- Never `package handler` (singular) or `package svc` (abbreviated)

**Types:**
- Services: `PerspectiveService`, `ContentService`, `UserService`, `AuthService` — `internal/core/services/*_service.go`
- Repositories (ports/interfaces): `ContentRepository`, `UserRepository` — `internal/core/ports/repositories/`
- Repository implementations (GORM): `gorm_content_repository.go`, `gorm_user_repository.go` in `internal/adapters/repositories/postgres/`
- Domain models: singular, matches table — `Content`, `User`, `Perspective` in `internal/core/domain/`

**Methods (repository layer):**
```go
GetByID(ctx, id)           // Single by ID
GetByUserID(ctx, userID)   // Single by foreign key
List(ctx, opts)            // Multiple with options/params
Create(ctx, model)         // Insert
Update(ctx, model)         // Update
Delete(ctx, id)            // Delete
```

**Methods (service layer)** — may combine multiple repo/adapter calls into one operation:
```go
GetPerspectiveWithContent(ctx, id)
CreatePerspectiveForUser(ctx, userID, input)
CreateFromYouTube(ctx, url, userID)   // e.g. ContentService
UpdateSourceData(ctx, contentID)
```

**Constructors:** `NewXService(deps...) *XService`, `NewXRepository(db) *XRepository` — plain constructor functions, no factory pattern.

### Code Style

**Formatting:** `gofmt` via `make fmt`. Pre-commit hook (`make install-hooks`) runs gofmt + prettier automatically.

**Linting:** `make lint` — must pass with zero errors before commit (see Code Quality Checklist below).

### Import Organization

Three groups, blank-line separated, in this order:
```go
import (
    // 1. Standard library
    "context"
    "encoding/json"
    "net/http"

    // 2. Third-party
    "github.com/go-chi/chi/v5"
    "gorm.io/gorm"

    // 3. Internal (this module)
    "github.com/CodeWarrior-debug/perspectize/backend/internal/core/domain"
    "github.com/CodeWarrior-debug/perspectize/backend/internal/core/services"
)
```

### Error Handling

**Sentinel errors** defined once in `internal/core/domain/errors.go`:
```go
var (
    ErrNotFound       = errors.New("resource not found")
    ErrAlreadyExists  = errors.New("resource already exists")
    ErrInvalidInput   = errors.New("invalid input")
    ErrInvalidURL     = errors.New("invalid URL")
    ErrYouTubeAPI     = errors.New("youtube API error")
    ErrInvalidRating  = errors.New("rating must be between 0 and 10000")
    ErrSentinelUser   = errors.New("cannot modify the system sentinel user")
    ErrDeleteSentinel = errors.New("cannot delete the system sentinel user")
)
```

**Wrapping with `%w`** to preserve `errors.Is` checks, and often prefixed with the sentinel for typed matching plus a human detail:
```go
// Wrap with context (fmt.Errorf + %w)
return fmt.Errorf("failed to get resource: %w", err)

// Prefix a sentinel error with detail (still satisfies errors.Is(err, domain.ErrInvalidInput))
return nil, fmt.Errorf("%w: content id must be a positive integer", domain.ErrInvalidInput)
```

**Resolver-level translation** — services return domain errors; GraphQL resolvers/handlers translate to user-facing messages and log unexpected errors instead of leaking internals:
```go
content, err := h.service.GetByID(ctx, id)
if err != nil {
    switch {
    case errors.Is(err, domain.ErrNotFound):
        return nil, fmt.Errorf("content not found")
    case errors.Is(err, domain.ErrValidation):
        return nil, err
    default:
        h.logger.Error("unexpected error", "error", err)
        return nil, fmt.Errorf("internal error")
    }
}
```

**Never leak sanitized/external API details to GraphQL clients** — log the raw error via `slog`, return a generic message. See `internal/core/services/content_service.go` `CreateFromYouTube` (YouTube API failure path).

**Validation pattern:** simple guard clauses at the top of service methods, returning `%w`-wrapped `ErrInvalidInput`, not a validation library, e.g.:
```go
if id <= 0 {
    return nil, fmt.Errorf("%w: content id must be a positive integer", domain.ErrInvalidInput)
}
```
(`go-playground/validator` is a listed dependency but service methods largely hand-roll simple checks like this.)

### Logging

Structured logging via `log/slog`, always with key-value pairs (never string interpolation into the message):
```go
slog.Error("failed to fetch YouTube metadata",
    "videoID", videoID,
    "userID", userID,
    "error", err)

logger.Info("perspective created",
    "perspective_id", p.ID,
    "user_id", p.UserID,
    "quality", p.Quality,
)
```

### Comments

- Doc comments on every exported type/function, in standard Go style (`// ContentService implements business logic for content operations`).
- Multi-line "why" comments above non-obvious logic blocks (e.g., numbered steps: `// 1. Extract video ID first (validates URL format)`).
- `*TEMP*` marker (project-wide convention, see root `CLAUDE.md`) for temporary/learning comments meant to be grepped and removed later:
  ```go
  // *TEMP* - defer runs after function returns, ensures cleanup
  defer db.Close()
  ```

### Function Design

- Service methods take `context.Context` as the first parameter, always.
- Guard-clause validation first, then delegate to repository/adapter, then map/return.
- Prefer named/struct input types for multi-field operations, e.g. `portservices.CreateClaimInput`, over long positional parameter lists.

### Module Design

- **Hexagonal / ports-and-adapters.** Domain (`core/domain/`) has zero external imports. Services depend only on port interfaces (`core/ports/`), never on concrete adapters. Adapters (`adapters/repositories/`, `adapters/graphql/`, `adapters/youtube/`) implement those ports. Wiring happens once, in `cmd/server/main.go`.
- **GORM "hex-clean separate model" pattern:** domain models never carry `gorm:` tags. Separate GORM-tagged structs live in `adapters/repositories/postgres/gorm_models.go`, with bidirectional mapper functions in `gorm_mappers.go`. Never add GORM tags directly to `core/domain` structs.
- **Enums:** UPPERCASE string constants in domain, bound via `gqlgen.yml` model binding — never hand-written switch statements for GraphQL enum conversion (see backend `CLAUDE.md` "Enum & ID Handling").
- No barrel files — Go doesn't use them; each package is imported by its own path.

### Configuration

Environment variables with typed fallback helpers:
```go
func LoadConfig() *Config {
    return &Config{
        Port:        getEnv("PORT", "8080"),
        DatabaseURL: getEnv("DATABASE_URL", ""),
        LogLevel:    getEnv("LOG_LEVEL", "info"),
        Env:         getEnv("ENV", "development"),
    }
}
```
Precedence: env vars override `config/config.json`. See `internal/config/`.

### Code Quality Checklist (pre-commit, backend)

- [ ] `make fmt` passes
- [ ] `make lint` passes with no errors
- [ ] `make test` passes
- [ ] New code has tests
- [ ] Error messages are helpful
- [ ] No hardcoded credentials
- [ ] Context is propagated
- [ ] Errors are wrapped with context (`%w`)

---

## Frontend (SvelteKit / Svelte 5 / TypeScript)

### Naming Patterns

**Files:**
- Components: PascalCase `.svelte` — `Header.svelte`, `ActivityTable.svelte`, `PerspectivePopover.svelte`, in `src/lib/components/` (feature subfolders like `discover/` for grouped components, `shadcn/` for primitives).
- Query/hook modules: camelCase `.ts` — `src/lib/queries/content.ts`, `src/lib/queries/hooks/useCreateClaim.ts`, `src/lib/queries/hooks/useUpdatePerspective.ts` (one hook per file, `use`-prefixed).
- Utility modules: camelCase or kebab-lowercase, purpose-named — `src/lib/utils/formatting.ts`, `src/lib/utils/grid-config.ts`, `src/lib/utils/gridUrlState.ts`, `src/lib/utils/ratings.ts`, `src/lib/utils/sanitize.ts`.
- Test files mirror the source name with `.test.ts`, but live in a parallel `tests/` tree, not co-located (see TESTING.md).

**shadcn-svelte primitives:** live specifically in `src/lib/components/shadcn/` (not the CLI's default `ui/`) — always verify install location and add new components to the `shadcn/index.ts` barrel export.

### Code Style

**Formatting:** Prettier, config in `frontend/.prettierrc`:
```json
{
  "useTabs": true,
  "singleQuote": true,
  "trailingComma": "all",
  "printWidth": 120,
  "plugins": ["prettier-plugin-svelte"],
  "overrides": [{ "files": "*.svelte", "options": { "parser": "svelte" } }]
}
```
Run via `pnpm run format` / `pnpm run format:check`. Pre-commit hookify rule warns to run `pnpm exec prettier --write` on staged files.

**Type checking:** `pnpm run check` (svelte-check + TypeScript), no dedicated ESLint config present in `frontend/` — type-checking + Prettier are the enforced style gates, not ESLint.

### Svelte 5 Runes (REQUIRED — no Svelte 4 syntax)

| Use | Do NOT use |
|---|---|
| `let count = $state(0)` | `let count = 0` with `$:` |
| `let doubled = $derived(count * 2)` | `$: doubled = count * 2` |
| `let { data, children } = $props()` | `export let data` |
| `$effect(() => { ... })` | `onMount` / `$:` side effects |
| `{@render children()}` | `<slot />` |
| `onclick={handler}` | `on:click={handler}` |

Never use `$effect` for pure derivation — use `$derived` instead.

### TanStack Query Pattern (Svelte 5 function-wrapper API)

Query/mutation options are reactive — always pass a **function** returning the options object, not the object directly. Results are accessed as plain object properties, never with `$` store prefix:
```ts
const query = createQuery(() => ({
    queryKey: ['content'],
    queryFn: () => graphqlClient.request(LIST_CONTENT),
}));
```
```svelte
{#if query.isLoading}Loading...{/if}
```
Query definitions (`gql` tagged templates) live in `src/lib/queries/*.ts` (e.g. `content.ts`, `claims.ts`, `perspectives.ts`, `users.ts`); mutation hooks live in `src/lib/queries/hooks/use*.ts` and wrap `createMutation`, handling `onSuccess` (toast + `queryClient.invalidateQueries`) and `onError` (toast, with string-matching on error message to pick a user-facing copy — see `useCreateClaim.ts` pattern below).

### Icons

Per-icon imports for tree-shaking, kebab-case icon name:
```svelte
import XIcon from '@lucide/svelte/icons/x';
```

### Import Organization

No enforced import-order tool observed; conventionally: external packages first, then `$lib/*` aliased internal imports, then relative/local imports. `$lib` alias maps to `src/lib`.

### Error Handling (frontend)

- Mutation hooks catch/branch on the caught `Error`'s message text (case-insensitive substring match) to choose a specific toast message, falling back to a generic "Failed to X. Please try again." — see `src/lib/queries/hooks/useCreateClaim.ts`.
- User feedback exclusively via `svelte-sonner` toast (`toast.success(...)`, `toast.error(...)`), not inline alert banners, for mutation results.

### Comments

Sparse in-source comments; used mainly to flag non-obvious gotchas inline (timezone handling, AG Grid quirks) — many of these are promoted to `frontend/CLAUDE.md` "Gotchas"/"Testing Gotchas" sections rather than kept purely as code comments, so check that file before re-deriving an explanation.

### Module Design

- No barrel files project-wide except `src/lib/components/shadcn/index.ts` (shadcn primitives) — mandatory to keep updated when adding a new shadcn component.
- One hook per file under `src/lib/queries/hooks/`.
- Pure/testable logic is deliberately extracted out of AG-Grid-coupled components into plain TS modules (`$lib/utils/grid-config.ts`, `$lib/utils/formatting.ts`) specifically so it's unit-testable without a DOM grid (see TESTING.md "AG Grid testing strategy").

---

## Cross-Cutting

**No chained bash commands** in any tooling/agent-authored scripts (project-wide dev workflow rule, not code style, but affects how Makefile/CI-adjacent scripts should be written/invoked).

**Commit messages:** Conventional Commits (`feat`, `fix`, `refactor`, `chore`, `docs`, `test`). One logical change per commit.

---

*Convention analysis: 2026-09-04*
