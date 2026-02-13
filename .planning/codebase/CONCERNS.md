# Codebase Concerns

**Analysis Date:** 2026-02-07

Comprehensive review of technical debt, bugs, and risk areas. Primary source: `KNOWN_BUGS.md` (comprehensive audit 2026-02-07), verified against actual codebase. Organized by severity and impact.

---

## CRITICAL ISSUES

### C-01: No Authentication or Authorization

**Risk:** Any client can CRUD any user's data. Complete security bypass.

**Files:**
- `backend/cmd/server/main.go` (line 74-93: no auth middleware)
- `backend/internal/adapters/graphql/resolvers/schema.resolvers.go` (no auth checks in any mutation)

**Problem:** GraphQL resolvers process all queries without verifying user identity or permissions. A malicious client can:
- List all users with email addresses exposed
- Create perspectives/content for any user
- Modify/delete any user's data
- Modify/delete any content

**Impact:** Application unsuitable for multi-user deployment. Data integrity cannot be guaranteed.

**Fix approach:**
1. Add authentication middleware (JWT, OAuth2, or session-based)
2. Inject authenticated user into request context
3. Add authorization checks in all mutations (e.g., `perspective.userID == currentUser.ID`)
4. Update GraphQL schema to include `Query.me` endpoint
5. Add `user` input parameter to mutations instead of deriving from auth context

---

### C-02: Cursor Pagination Broken for Non-ID Sorts

**Risk:** Wrong pages returned when sorting by name/date.

**Files:**
- `backend/internal/adapters/repositories/postgres/content_repository.go:207-336`
- `backend/internal/adapters/repositories/postgres/perspective_repository.go:233-362`

**Problem:** Keyset pagination cursor only encodes `id`:
```go
// cursor format: base64("id:<id>")
```

When sorting by `CREATED_AT` or `NAME`, the next page query uses the last ID but wrong sort direction, producing duplicates or missing rows.

**Correct keyset pagination requires:**
- Encode both `id` AND the sort column value in cursor
- Construct WHERE clause with compound condition: `(sortCol, id) > (lastSortVal, lastId)`

**Impact:** Pagination UX broken for any content/perspective list sorted by date or name. Users see duplicates or missing items.

**Fix approach:**
1. Redesign cursor to encode `{id, sortColumnValue}` as JSON then base64
2. Update `decodeCursor` to extract both values
3. Refactor WHERE clause construction to use compound keyset logic
4. Add tests for pagination with non-ID sorts (currently zero coverage)

---

### C-03: XSS Vulnerability in AG Grid cellRenderer

**Risk:** User-controlled data interpolated into HTML without escaping.

**Files:** `fe/src/lib/components/ActivityTable.svelte:64-70`

**Problem:**
```svelte
cellRenderer: (params: { data?: ContentRow }) => {
    if (!params.data) return '';
    if (params.data.url) {
        return `<a href="${params.data.url}" target="_blank" rel="noopener noreferrer" class="text-primary hover:underline">${params.data.name}</a>`;
    }
    return params.data.name;
}
```

If `params.data.name` contains `<img src=x onerror=alert(1)>` or `params.data.url` contains `javascript:alert(1)`, it executes in the grid cell.

**Impact:** XSS attack via backend response. Attacker can steal session tokens, post as user, modify page.

**Fix approach:**
1. Use `document.createElement()` instead of template string interpolation
2. Or use AG Grid's built-in sanitization (check ag-grid-svelte5 options)
3. Or use Svelte template syntax with `{@html ...}` guards for explicit trust boundary

---

### C-04: No GraphQL Query Complexity Limiting

**Risk:** DoS vector via deeply nested queries.

**Files:** `backend/cmd/server/main.go:75` (handler initialization)

**Problem:** gqlgen server has no complexity calculator configured. Query like:
```graphql
query {
  perspectives { perspectives { perspectives { ... } } }
}
```
will cause unbounded recursion or O(n²) query execution.

**Impact:** Attacker can crash backend with single malicious query.

**Fix approach:**
1. Add `complexity.go` with `ComplexityCalculator` function
2. Register in handler config before `NewDefaultServer()`
3. Set complexity budget (e.g., 1000) and check before execution
4. Test with nested query bombs

---

### C-05: Wildcard CORS Configuration

**Risk:** `Access-Control-Allow-Origin: *` allows any origin.

**Files:** `backend/cmd/server/main.go:80`

**Problem:**
```go
w.Header().Set("Access-Control-Allow-Origin", "*")
```

A malicious website can make requests to the GraphQL API on behalf of your users' browsers (if authentication were in place). Combined with C-01 (no auth), this is less critical but still wrong.

**Impact:** Cross-site request forgery (CSRF) possible once auth is added. Currently moot due to C-01.

**Fix approach:**
1. Replace wildcard with explicit frontend URL (e.g., `https://perspectize.com`)
2. Use environment variable for origin (dev = `http://localhost:5173`, prod = frontend domain)
3. Update `fe/CLAUDE.md` to document required CORS setup

---

### C-06: Silent JSON Unmarshal Failure

**Risk:** Corrupted data silently omitted from responses.

**Files:** `backend/internal/adapters/repositories/postgres/perspective_repository.go:419-426`

**Problem:** `categorizedRatings` JSON field is unmarshaled without error checking:
```go
json.Unmarshal([]byte(dbCategorizedRatings), &categorizedRatings)
// error is ignored — bad data becomes nil array
```

If database contains invalid JSON, response silently drops the field instead of failing or logging.

**Impact:** Users see incomplete perspective data without knowing why. Data loss appears random.

**Fix approach:**
1. Add error check: `if err := json.Unmarshal(...); err != nil { return nil, fmt.Errorf(...) }`
2. Add structured logging to all JSON unmarshal operations
3. Add repository tests for malformed JSON handling

---

### C-07: Silent Duration Parse Failure

**Risk:** Bad duration defaults to 0 seconds, indistinguishable from real short video.

**Files:** `backend/internal/adapters/youtube/client.go:90-93`

**Problem:**
```go
duration, _ := time.ParseDuration(durationStr)
// error is ignored
```

If YouTube API returns unparseable duration string, the field silently becomes 0 without logging.

**Impact:** UI displays "0 seconds" for videos with bad metadata. No visibility into data quality.

**Fix approach:**
1. Add error handling with structured log
2. Or return `*int` (nil if unparseable) instead of silent 0
3. Add tests for non-ISO8601 duration formats

---

### C-08: Five Silent Parse Failures in Domain Conversion

**Risk:** Response, viewCount, likeCount, commentCount all silently become nil on parse error.

**Files:** `backend/internal/adapters/graphql/resolvers/helpers.go:36-63`

**Problem:** `domainToModel` unmarshal errors are silently discarded:
```go
json.Unmarshal(c.Response, &response)  // error ignored
strconv.Atoi(c.ViewCount)              // error ignored
strconv.Atoi(c.LikeCount)              // error ignored
strconv.Atoi(c.CommentCount)           // error ignored
```

**Impact:** Incomplete data in GraphQL responses with no error signal. Users can't tell if counts are 0 or corrupted.

**Fix approach:**
1. Add error handling and structured logging for all conversions
2. Return error from `domainToModel` or use proper nullable fields
3. Add resolver tests for malformed input

---

### C-09: GraphQL Playground Exposed Unconditionally

**Risk:** Introspection enabled without environment check.

**Files:** `backend/cmd/server/main.go:92`

**Problem:**
```go
http.Handle("/", playground.Handler("GraphQL Playground", "/graphql"))
```

Playground is accessible in production, exposes schema to anyone.

**Impact:** Schema enumeration enables targeted attacks. Best practice is to disable in production.

**Fix approach:**
1. Check `APP_ENV` or `DEBUG` env var
2. Only register playground in dev mode
3. Also disable GraphQL introspection in production (see C-10)

---

### C-10: GraphQL Introspection Enabled Without Restriction

**Risk:** Full schema introspection available to all clients.

**Files:** `backend/cmd/server/main.go:75` (no introspection config)

**Problem:** gqlgen server has default `IntrospectionEnabled: true`. Combined with exposed playground (C-09), attackers enumerate all queries/mutations.

**Impact:** Complete API surface visible. Enables reconnaissance for targeted attacks.

**Fix approach:**
1. Add introspection config check:
   ```go
   cfg := generated.Config{
       IntrospectionEnabled: os.Getenv("ENABLE_INTROSPECTION") == "true",
   }
   ```
2. Default to false in production
3. Allow override for dev/staging only

---

## HIGH PRIORITY ISSUES

### H-01: Adapter-to-Adapter Coupling

**Risk:** Violates hexagonal architecture dependency rule.

**Files:** `backend/internal/adapters/graphql/resolvers/schema.resolvers.go:16,23`

**Problem:** Resolver imports YouTube adapter directly:
```go
import "github.com/CodeWarrior-debug/perspectize/backend/internal/adapters/youtube"
```

Dependencies should flow: adapter → service → port. Adapter never talks to adapter.

**Impact:** Services layer bypassed. Tight coupling makes testing and swapping implementations difficult.

**Fix approach:**
1. Verify all YouTube operations are in `ContentService` (they should be)
2. Remove direct youtube imports from resolvers
3. Add architecture test to prevent adapter-to-adapter imports

---

### H-02: Resolver Depends on Concrete Service Types

**Risk:** Missing service port interfaces.

**Files:** `backend/internal/adapters/graphql/resolvers/resolver.go:12-16`

**Problem:**
```go
type Resolver struct {
    contentService *services.ContentService  // concrete, not interface
    userService *services.UserService        // concrete, not interface
    perspectiveService *services.PerspectiveService
}
```

Should depend on interfaces, not concrete types.

**Impact:** Can't mock services for testing. Resolver tests must use real service implementations.

**Fix approach:**
1. Create port interfaces: `ContentServicePort`, `UserServicePort`, `PerspectiveServicePort`
2. Update resolver to accept interfaces
3. Update `cmd/server/main.go` wiring
4. Update resolver tests to use mocks

---

### H-03: ListAll() Users Has No Pagination

**Risk:** Unbounded query result set.

**Files:** `backend/internal/adapters/repositories/postgres/user_repository.go:98-114`

**Problem:**
```go
func (r *UserRepository) ListAll(ctx context.Context) ([]domain.User, error) {
    // SELECT * FROM users — no LIMIT
}
```

Query returns all users. If 10,000 users exist, all rows loaded into memory.

**Impact:** Memory exhaustion DoS. Unbounded response size (C-10 + H-03 = attacker can request massive response).

**Fix approach:**
1. Add `limit int` parameter (or use cursor pagination from H-03)
2. Default to reasonable limit (50-100)
3. Return error if limit exceeds max (1000)
4. Update GraphQL schema to require pagination

---

### H-04 & H-05: GraphQL Type Schema Issues

**Risk:** Weak API contracts.

**Files:** `backend/schema.graphql`

**Problems:**
- **H-04:** Timestamps as `String!` instead of custom `DateTime` scalar (lines 9-10, 56-58, 77-78)
- **H-05:** `contentType` uses `String!` instead of defined `ContentType` enum (line 70)

**Impact:** No type safety for timestamps (clients must parse manually). Content type values not enumerated (clients don't know valid values).

**Fix approach:**
1. Define `scalar DateTime` in schema
2. Implement DateTime scalar resolver for serialization
3. Update all timestamp fields to `DateTime!`
4. Define `ContentType` enum (e.g., `YOUTUBE`, `VIMEO`)
5. Bind enum in `gqlgen.yml`
6. Update content type storage to use enum values

---

### H-06 & H-07: Race Conditions on Uniqueness Checks

**Risk:** Duplicate inserts possible under concurrent load.

**Files:**
- `backend/internal/core/services/perspective_service.go:91-97` (duplicate claim check)
- `backend/internal/core/services/user_service.go:49-65` (duplicate user check)

**Problem:** Classic TOCTOU (time-of-check-time-of-use) race:
```go
// Thread A: Check if claim exists
existing, _ := r.FindByClaimAndUser(ctx, claim, userID)
if existing != nil {
    return nil, ErrAlreadyExists
}
// [Thread B inserts here]
// Thread A: Insert new claim
r.Create(ctx, ...)  // UNIQUE constraint violated at DB level
```

**Impact:** Under load, concurrent create requests can both pass the check and fail at database, causing errors or partial inserts.

**Fix approach:**
1. Use database UNIQUE constraint as the sole source of truth
2. Catch DB unique violation error and return `ErrAlreadyExists`
3. Remove app-level duplicate check
4. Or use database-level transactions with explicit locks

---

### H-08: YouTube API Response Stored Verbatim

**Risk:** Bloat and information leakage.

**Files:** `backend/internal/adapters/youtube/client.go:100`

**Problem:** Entire YouTube API response (with metadata, thumbnails, etc.) stored in `Content.response: JSON` field. YouTube response is ~5KB per video.

**Impact:** Database bloat. Unnecessary data stored increases backup size, query time. No use case for storing full response.

**Fix approach:**
1. Parse response and extract only needed fields (title, duration, viewCount, etc.)
2. Store structured fields, not raw JSON
3. Remove `response: JSON` from schema (or make it optional for debugging)
4. Backfill existing data to remove YouTube responses

---

### H-09: Hardcoded Config Path

**Risk:** Config not flexible for deployments.

**Files:** `backend/cmd/server/main.go:27`

**Problem:**
```go
cfg, err := config.Load("config/config.example.json")
```

Path is hardcoded. Works in dev, fails in containers where file structure differs.

**Impact:** Docker builds fail. Production deployment requires workarounds.

**Fix approach:**
1. Use environment variable: `CONFIG_PATH = os.Getenv("CONFIG_PATH")`
2. Default to reasonable path if unset
3. Test in Docker container
4. Document in `.env.example`

---

### H-10: User Email Addresses Exposed

**Risk:** GDPR/privacy violation. Enables spam/phishing.

**Files:** `backend/internal/adapters/graphql/resolvers/schema.resolvers.go:302-315` (users query)

**Problem:** `Query.users` returns list with email addresses. No access control.

**Impact:** Scraping tool can dump all user emails. Violates privacy regulations.

**Fix approach:**
1. Add authentication check to `users` query
2. Only return own email, not others'
3. Or remove email from public user type, add separate `Query.me` endpoint
4. Audit other endpoints for exposed PII

---

### H-11: No Rate Limiting

**Risk:** Brute force, enumeration, DoS attacks.

**Files:** `backend/cmd/server/main.go` (no middleware)

**Problem:** No rate limiting on GraphQL endpoint. Attacker can spam queries without throttling.

**Impact:** Username enumeration (list all users via H-03 + H-11). Password brute force if auth added. Query complexity bombs (C-04).

**Fix approach:**
1. Add rate limiting middleware (e.g., `ulule/limiter`)
2. Limit by IP address (or user ID if authenticated)
3. Different limits for mutations vs queries
4. Return 429 Too Many Requests when exceeded

---

### H-12: YouTube API Key Exposure Risk

**Risk:** Key compromise enables unauthorized YouTube API calls.

**Files:**
- `backend/internal/adapters/youtube/client.go:53-57,76` (key in URL, error messages)
- Stored in config and environment

**Problem:** API key may appear in:
- HTTP error responses if request fails
- Server logs if key validation fails
- GitHub commits if `.env` checked in (see `.gitignore`)

**Impact:** Attacker can use compromised key to spam YouTube API, incurring charges.

**Fix approach:**
1. Never log API key (sanitize in error messages)
2. Add validation at startup (try dummy request, catch error, don't log key)
3. Use Cloud Key Management (e.g., AWS Secrets Manager) instead of env vars
4. Rotate key regularly
5. Audit git history: `git log -S "AIza" --`

---

### H-13: Sensitive Data Leaked in GraphQL Errors

**Risk:** Database schema/structure exposed via error messages.

**Files:** `backend/internal/adapters/graphql/resolvers/schema.resolvers.go:31,47,99,152,173`

**Problem:** Resolvers use `%w` formatting which wraps underlying errors:
```go
return nil, fmt.Errorf("failed to find content: %w", err)
// If err = "column 'foo' does not exist", attacker sees DB schema
```

**Impact:** SQL errors expose schema structure, table names, columns. Enables targeted SQL injection attempts.

**Fix approach:**
1. Return generic error to GraphQL clients: `fmt.Errorf("internal server error")`
2. Log full error server-side with structured logging
3. Add middleware to scrub errors before returning to client
4. Test: deliberately trigger DB errors, verify no schema leakage

---

### H-14 & H-15: No HTTPS/TLS and HTTP Timeouts

**Risk:** Man-in-the-middle attacks and Slowloris DoS.

**Files:** `backend/cmd/server/main.go:99` (http.ListenAndServe with no config)

**Problems:**
- **H-14:** No TLS/HTTPS. Traffic is unencrypted.
- **H-15:** Server has no timeouts. Slowloris attacker can open slow connections forever.

**Impact:** Credentials/API keys transmitted in plaintext. Server hangs waiting for slow clients.

**Fix approach:**
1. Use `http.Server` with timeouts:
   ```go
   srv := &http.Server{
       Addr:         ":8080",
       ReadTimeout:  15 * time.Second,
       WriteTimeout: 15 * time.Second,
       IdleTimeout:  60 * time.Second,
   }
   ```
2. For HTTPS: use `ListenAndServeTLS` with cert/key or reverse proxy (Caddy/Nginx)
3. Document TLS setup in deployment guide

---

### H-16: Inconsistent Not-Found Error Handling

**Risk:** Inconsistent GraphQL contracts.

**Files:** `backend/internal/adapters/graphql/resolvers/schema.resolvers.go:274,289,326`

**Problem:** Some resolvers return `nil, nil` (errors hidden from GraphQL) while others return proper errors:
```go
// Line 274: return nil, nil
// Line 289: return nil, nil
// vs ContentByID: returns error
```

**Impact:** Clients expect errors but get nulls. Silent failures hard to debug.

**Fix approach:**
1. Standardize: always return `(T, error)` with proper error handling
2. GraphQL layer converts errors to null fields as needed
3. Add resolver tests asserting error behavior

---

### H-17 & H-18: Missing Svelte Error Boundaries

**Risk:** Unhandled errors show blank page or default error.

**Files:**
- `fe/src/routes/` (missing `+error.svelte`)
- `fe/src/` (missing `hooks.client.ts`, `hooks.server.ts`)

**Problem:** No error boundary component. Errors outside TanStack Query are invisible to users.

**Impact:** If header fails to load, entire page is broken. Users see nothing or generic "Error" message.

**Fix approach:**
1. Create `src/routes/+error.svelte` with graceful error UI
2. Create `src/hooks.client.ts` to catch client-side errors
3. Log to error tracking service (see H-24 for client infra)

---

### H-19 & H-20: .env Load Failure and Empty API Key Validation

**Risk:** Silent misconfiguration.

**Files:** `backend/cmd/server/main.go:24` (.env ignored) and `youtube/client.go:23` (no key validation)

**Problems:**
- `.env` load failure is silently ignored: `_ = godotenv.Load()`
- YouTube API key not validated at startup. Fails with cryptic 403 at runtime.

**Impact:** Typo in .env goes unnoticed. Configuration errors only surface when queries run.

**Fix approach:**
1. Warn if .env file expected but missing: `if _, err := os.Stat(".env"); err != nil && os.Getenv("ENVIRONMENT") != "production" { log.Warn(...) }`
2. Add config validation: `if cfg.YouTube.APIKey == "" { log.Fatal("YOUTUBE_API_KEY required") }`
3. Test YouTube key at startup with dummy request

---

### H-21: WriteString Return Ignored

**Risk:** Response corruption.

**Files:** `backend/pkg/graphql/intid.go:17`

**Problem:**
```go
func (i IntID) MarshalJSON() ([]byte, error) {
    _, _ = io.WriteString(...) // return value ignored
}
```

If `WriteString` fails, no error is returned. Response may be incomplete.

**Impact:** IntID serialization silently fails. Clients receive null IDs.

**Fix approach:**
1. Check return value: `if n, err := io.WriteString(...); err != nil { return nil, err }`
2. Add tests for WriteString error cases

---

### H-22: prerender = true Without SSR

**Risk:** Architectural mismatch.

**Files:** `fe/src/routes/+layout.ts:1`

**Problem:**
```typescript
export const prerender = true;
```

With `adapter-static`, this tells SvelteKit to prerender all routes as static HTML. But the app fetches dynamic GraphQL data, making prerender pointless. App runs as SPA.

**Impact:** Build time wastage. No SEO benefit (content is JS-rendered). Confusing architecture.

**Fix approach:**
1. Set `prerender = false` to use SPA mode explicitly
2. Or actually leverage prerendering by fetching data at build time (requires build-time GraphQL endpoint)

---

### H-23: No TypeScript Types from GraphQL Schema

**Risk:** Manual duplication and drift.

**Files:** `fe/src/lib/components/` (all component files)

**Problem:** Type definitions for GraphQL responses are manually written in Svelte components:
```typescript
interface ContentItem {
    id: string;
    name: string;
    // ... manually duplicated from schema
}
```

No code generation from schema. Changes to GraphQL schema require manual updates.

**Impact:** Types drift from schema. Type safety lost. Updates require multiple changes.

**Fix approach:**
1. Use `graphql-codegen` to generate TypeScript types from GraphQL schema
2. Run as part of build: `gql-codegen` before `pnpm build`
3. Import types from generated file
4. All types stay in sync with schema

---

### H-24: GraphQL Client Missing Error/Timeout Infrastructure

**Risk:** No error recovery, no timeout protection.

**Files:** `fe/src/lib/queries/client.ts:1-7`

**Problem:**
```typescript
const graphqlClient = new GraphQLClient("http://localhost:8080/graphql");
```

Client has no:
- Error interceptor to catch network errors
- Timeout configuration (requests hang forever)
- Authorization header support
- Request/response logging

**Impact:** Network errors not handled. Queries that fail are invisible. No auth infrastructure.

**Fix approach:**
1. Add error interceptor to catch network failures
2. Set timeout (e.g., 30 seconds)
3. Add `headers` callback to inject auth token (prepare for C-01 fix)
4. Add request/response logging

---

### H-25: No Content Security Policy

**Risk:** XSS/injection attacks.

**Files:** `fe/app.html` (no CSP header)

**Problem:** No CSP header restricts what scripts can run. Combined with C-05 (innerHTML XSS), attacks easier.

**Impact:** XSS payloads can load external scripts, exfiltrate data.

**Fix approach:**
1. Add CSP header in `app.html` or server middleware
2. Recommended: `default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'` (AG Grid needs unsafe-inline)
3. Report violations to logging service

---

### H-26: No CI/CD or Security Scanning

**Risk:** No automated checks for regressions, secrets, vulnerable dependencies.

**Files:** `.github/` (missing workflows)

**Problem:** No GitHub Actions workflows for:
- Testing on every commit
- Dependency scanning (Dependabot)
- Secret scanning (SAST)
- Container scanning if Docker used

**Impact:** Vulnerable packages not detected. Secrets committed. Breaking changes merged.

**Fix approach:**
1. Add `.github/workflows/test.yml` for Go tests
2. Add `.github/workflows/test-fe.yml` for frontend tests
3. Enable Dependabot in repo settings
4. Add SAST scanning (e.g., `github/super-linter`)
5. Document CI requirements in CLAUDE.md

---

## MEDIUM PRIORITY ISSUES

### M-01: Dual PostgreSQL Driver Dependencies

**Issue:** `lib/pq` and `pgx/v5` both imported

**Files:** `backend/go.mod:10-11`

**Impact:** Unnecessary bloat, potential conflicts. Choose one.

**Fix:** Remove unused driver. If using sqlx, use `pgx` driver exclusively.

---

### M-02: Hardcoded Database Connection Pool Settings

**Issue:** No configuration for pool size, timeout.

**Files:** `backend/pkg/database/postgres.go:21-23`

**Problem:** Max 25 open connections, 5 idle hard-coded.

**Fix:** Load from env vars with sensible defaults.

---

### M-03: CreateFromYouTube Returns Error Instead of Idempotent Result

**Issue:** Error response for duplicate, should be idempotent.

**Files:** `backend/internal/core/services/content_service.go:30-36`

**Fix:** Check for `ErrAlreadyExists` and return existing item, not error.

---

### M-04: Schema Type Inconsistency

**Issue:** `deletePerspective` uses `ID` scalar instead of `IntID`.

**Files:** `backend/schema.graphql:185`

**Fix:** Standardize all IDs to use `IntID` scalar.

---

### M-05: Function Parameter Instead of Dependency Injection

**Issue:** `CreateFromYouTube` accepts `extractVideoID` as parameter.

**Files:** `backend/internal/core/services/content_service.go:28`

**Fix:** Inject into service constructor instead.

---

### M-06: No Request Logging Middleware

**Issue:** Uses default `net/http` mux, no middleware chain.

**Files:** `backend/cmd/server/main.go:91-94`

**Impact:** No visibility into request/response for debugging.

**Fix:** Add request logging middleware (use `chi` router with middleware).

---

### M-07: Inconsistent Not-Found Error Handling

**Issue:** Different approaches across resolvers (see H-16).

**Files:** `backend/internal/adapters/graphql/resolvers/schema.resolvers.go`

**Fix:** Standardize error return pattern.

---

### M-08: Missing Nested Field Resolvers

**Issue:** `user` and `content` fields on Perspective return null instead of fetching.

**Files:** `backend/internal/adapters/graphql/resolvers/helpers.go:70-107`

**Problem:**
```graphql
type Perspective {
    user: User  # Always null
    content: Content  # Always null
}
```

Clients must separately fetch user/content after fetching perspective.

**Fix:** Implement field resolvers to fetch nested objects.

---

### M-09: No Graceful Shutdown Handler

**Issue:** Server kills in-flight requests on shutdown.

**Files:** `backend/cmd/server/main.go:99`

**Fix:** Add signal handler for SIGTERM with graceful shutdown timeout.

---

### M-10: No Health Check Endpoint

**Issue:** Load balancers/orchestrators have no way to check health.

**Files:** `backend/cmd/server/main.go:91-93`

**Fix:** Add `/health` (liveness) and `/ready` (readiness) endpoints.

---

### M-11: Missing Input Length Validation

**Issue:** No length checks on description, labels, categorizedRatings.

**Files:** `backend/internal/core/services/perspective_service.go`, `user_service.go`

**Impact:** Unbounded inputs can cause performance issues.

**Fix:** Add validator: `description max 1000 chars`, `labels max 10 items`, etc.

---

### M-12: DB Credentials in Logs on Failure

**Issue:** Connection string may appear in logs.

**Files:** `backend/cmd/server/main.go:43-44`, `config.go:83`

**Fix:** Sanitize DSN before logging (redact password).

---

### M-13: Unbounded JSON Field

**Issue:** `response: JSON` field stores full YouTube response (~5KB per item).

**Files:** `backend/schema.graphql:77`

**Impact:** Bloats database, unnecessary data in queries.

**Fix:** See H-08.

---

### M-14: Missing Security Headers

**Issue:** No X-Content-Type-Options, X-Frame-Options, HSTS, etc.

**Files:** `backend/cmd/server/main.go`

**Fix:** Add security headers middleware:
```go
w.Header().Set("X-Content-Type-Options", "nosniff")
w.Header().Set("X-Frame-Options", "DENY")
w.Header().Set("Strict-Transport-Security", "max-age=31536000")
```

---

### M-15: No CSRF Protection

**Issue:** No anti-CSRF tokens.

**Files:** `backend/cmd/server/main.go:93`

**Impact:** Moot until C-01 (auth) is fixed, but should be added.

**Fix:** Add CSRF middleware after auth implemented.

---

### M-16: Update Does Not Check RowsAffected

**Issue:** Race condition on concurrent updates.

**Files:** `backend/internal/adapters/repositories/postgres/perspective_repository.go:187-209`

**Problem:** Unlike Delete, Update doesn't check if row was actually modified (optimistic lock missing).

**Fix:** Check `result.RowsAffected() > 0` or add `version` column for optimistic locking.

---

### M-17: No DATABASE_URL Format Validation

**Issue:** Invalid URL accepted without error.

**Files:** `backend/internal/config/config.go:79-84`

**Fix:** Validate DSN format at startup.

---

### M-18: Duplicated Type Definitions Across Components

**Issue:** ContentItem, ContentRow, User types manually defined in multiple files.

**Files:** `fe/src/lib/components/ActivityTable.svelte`, `+page.svelte`, `UserSelector.svelte`

**Fix:** See H-23 (use codegen).

---

### M-19: No Server-Side Pagination Integration

**Issue:** Content query hard-codes 100 items fetch.

**Files:** `fe/src/routes/+page.svelte:33-34`

**Problem:**
```typescript
const query = createQuery(() => ({
    queryFn: () => listContent({ first: 100 })  // Hard-coded
}));
```

AG Grid pagination not integrated with server. Fetches all 100 items then paginates client-side.

**Fix:** Integrate AG Grid pagination with cursor pagination from server.

---

### M-20: selectedUserId Store Not Consumed

**Issue:** Wired but never used.

**Files:** `fe/src/lib/stores/userSelection.svelte.ts`, `src/routes/+page.svelte`

**Fix:** Either use in content query filter or remove.

---

### M-21: Unused Type Guards

**Issue:** ContentResponse/ContentItem interfaces declared but never used as type guards.

**Files:** `fe/src/routes/+page.svelte:8-28`

**Fix:** Remove unused types or implement runtime validation.

---

### M-22: Search Input Not Debounced

**Issue:** AG Grid filter triggered on every keystroke.

**Files:** `fe/src/routes/+page.svelte:30`, `ActivityTable.svelte:130-133`

**Impact:** Excessive queries sent to server.

**Fix:** Add 300ms debounce to search input.

---

### M-23: No Error Recovery UI

**Issue:** Error states have no retry button.

**Files:** `fe/src/routes/+page.svelte:70-73`, `UserSelector.svelte:37-40`

**Fix:** Add retry button on error, call `query.refetch()`.

---

### M-24: Dead Code in Production

**Issue:** AGGridTest.svelte never imported but in component tree.

**Files:** `fe/src/lib/components/AGGridTest.svelte`

**Fix:** Delete or comment.

---

### M-25: HTTP Fallback for GraphQL Endpoint

**Issue:** Fallback uses HTTP not HTTPS.

**Files:** `fe/src/lib/queries/client.ts:3`

**Problem:**
```typescript
const endpoint = process.env.VITE_GRAPHQL_URL || "http://localhost:8080/graphql";
```

**Fix:** Use HTTPS in production endpoint.

---

### M-26: Retry Configuration Retries All Errors

**Issue:** `retry: 1` retries 4xx errors (should only retry network/5xx).

**Files:** `fe/src/routes/+layout.svelte:15`, `+page.svelte:40`

**Fix:** Use `shouldRetry: (failureCount, error) => error.status >= 500 || !error.response`

---

### M-27: formatDate Silently Produces Invalid Date

**Issue:** Bad input produces "Invalid Date" string instead of error.

**Files:** `fe/src/lib/components/ActivityTable.svelte:48-53`

**Fix:** Add validation or return fallback with warning.

---

### M-28: No Secret Rotation or Vault Integration

**Issue:** Secrets stored in .env, no rotation mechanism.

**Impact:** Compromised secret is permanent.

**Fix:** Use AWS Secrets Manager, HashiCorp Vault, or 1Password for rotation.

---

## LOW PRIORITY ISSUES

### L-01 through L-22

Code style, unused dependencies, test organization issues. See `KNOWN_BUGS.md` lines 88-110 for full details.

**Priority fixes:**
- **L-16:** Remove `@tanstack/svelte-form` (unused dependency)
- **L-23:** Inject `ref` prop instead of using `any` type

---

## UI BUGS (Phase 2.1)

See `KNOWN_BUGS.md` "UI Bugs" section (lines 111-121) for mobile responsiveness issues.

**Critical (P1):**
- BUG-001: Header overflow at 375px
- BUG-002: Pagination bar broken at 375px
- BUG-003: Sticky header clipping persists on scroll

---

## TEST COVERAGE GAPS

**Critical (P1):**
- T-01: `PerspectiveService.Update()` — 100-line mutation with zero tests
- T-02: No resolver tests for User/Perspective queries/mutations
- T-03: No tests for `helpers.go` domain-to-model conversion with silent JSON parse failures

**High (P2):**
- T-04: No repository-layer tests
- T-05: No YouTube API client tests
- T-06: No `IntID` scalar tests

See `KNOWN_BUGS.md` lines 122-138 for complete test gap inventory.

---

## Impact Summary

| Severity | Count | Primary Risk | Key Issues |
|----------|-------|--------------|------------|
| Critical | 5 | Security (no auth), Data integrity | C-01, C-02, C-03, C-04, C-05 |
| High | 22 | Operational (errors, crashes, DoS) | H-01–H-26 |
| Medium | 28 | Code quality, maintainability | M-01–M-28 |
| Low | 22 | Style, cleanup | L-01–L-22 |
| **Total** | **77** | **Multi-layer** | Requires coordinated fixes |

---

*Audit completed 2026-02-07. See `KNOWN_BUGS.md` for complete metadata and source documentation.*
