# Codebase Concerns

**Analysis Date:** 2026-02-04

## Migration Status: C# to Go (Active Transition)

**Overall Status:** Partial migration in progress. Legacy C# codebase exists alongside new Go implementation.

**Critical Concern:** The repository contains two complete implementations:
- `perspectize-be/` - Legacy C# ASP.NET Core implementation (do not modify)
- `perspectize-go/` - New Go implementation (active development)

**Impact:** Risk of confusion about which codebase to modify. CLAUDE.md explicitly states "Do not modify C# code. All development occurs in perspectize-go/."

**Mitigation:** Documentation is clear, but enforcement depends on developer discipline. Consider:
- Adding CI checks to block PRs modifying C# code
- Archiving or removing C# code once Go migration complete
- Git pre-commit hooks to warn on C# changes

---

## Tech Debt & Missing Features

### 1. Authentication & Authorization Not Implemented

**Issue:** No authentication middleware wired into the GraphQL API.

**Files:**
- `perspectize-go/cmd/server/main.go` - Server setup (no auth middleware)
- `perspectize-go/internal/middleware/` - **Directory does not exist**
- Schema includes `user` field on mutations/queries but no permission checks

**Current State:**
- Users can be created via `createUser(input)` but no user context is available in resolvers
- No bearer token validation
- No session management
- All endpoints are public

**Impact:**
- Anyone can create, read, update users and perspectives
- No ownership validation (user A can modify user B's data)
- No privacy enforcement (public/private perspectives not enforced at API layer)

**Fix Approach:**
1. Create `internal/middleware/` directory
2. Implement JWT/bearer token extraction middleware
3. Add user context to GraphQL operation context
4. Add permission checks in resolvers before data access
5. Validate `userID` in mutations matches authenticated user

**Priority:** High - blocks production use

---

### 2. GraphQL Field Resolvers Missing for Nested Data

**Issue:** Schema defines nested object fields that likely have no resolver implementations.

**Files:**
- `perspectize-go/schema.graphql` (lines 41, 43) - Defines `user: User` and `content: Content` on Perspective
- `perspectize-go/internal/adapters/graphql/generated/generated.go` - Auto-generated (difficult to verify if field resolvers exist)

**Current State:**
- Perspective type includes `user: User!` and `content: Content` fields
- Helper function `perspectiveDomainToModel` (lines 70-107 in `internal/adapters/graphql/resolvers/helpers.go`) does NOT populate `user` or `content` fields
- Queries requesting nested User/Content data will return null without error

**Impact:**
- Clients cannot fetch user details from perspective query
- Clients cannot fetch content details from perspective query
- Incomplete data responses degrade usability

**Test Coverage:** No integration tests verify nested field resolution

**Fix Approach:**
1. Verify/implement field resolvers for `Perspective.user()` and `Perspective.content()`
2. Update helpers to populate these fields (will require additional repository calls)
3. Add integration tests that query nested fields
4. Consider DataLoader pattern for N+1 query prevention

**Priority:** Medium - breaks GraphQL contract

---

### 3. Limited Content Filtering

**Issue:** ContentFilter input type has minimal filtering options.

**Files:**
- `perspectize-go/schema.graphql` (lines 116-121) - ContentFilter definition
- Comment in `internal/adapters/graphql/generated/generated.go` (line 825): `# TODO: Add additional filters (e.g., dateRange, search) or make filters dynamic`

**Current Filters Available:**
- `contentType: ContentType` (enum: YOUTUBE only)
- `minLengthSeconds: Int`
- `maxLengthSeconds: Int`

**Missing Filters:**
- Date range (createdAt/updatedAt)
- Text search (name, URL)
- View count, like count ranges
- User-created perspectives count

**Impact:** Clients cannot efficiently discover content without loading all records

**Fix Approach:**
1. Extend ContentFilter type with additional fields
2. Update repository query builders to support new filters
3. Write tests for filter combinations

**Priority:** Low - feature enhancement, not blocker

---

## Code Quality Issues

### 4. String-to-Int Conversion Scattered in Resolvers (Anti-pattern)

**Issue:** Manual ID string-to-int conversion repeated in resolver functions instead of using custom scalar.

**Files:**
- `perspectize-go/internal/adapters/graphql/resolvers/schema.resolvers.go` (lines 159-162, 180-183, 265-268, 302-305) - Multiple `strconv.Atoi` calls
- `perspectize-go/pkg/graphql/intid.go` - IntID custom scalar exists but underutilized

**Current Pattern:**
```go
// In resolvers - manual conversion
intID, err := strconv.Atoi(id)
if err != nil {
    return false, fmt.Errorf("invalid perspective ID: %s", id)
}
```

**Why It's Debt:**
- Repetitive boilerplate scattered across resolver functions
- Error handling inconsistent (some return nil, some return false)
- IntID scalar exists but only used in input types, not on query/mutation IDs

**Impact:**
- Harder to maintain (6+ places convert IDs)
- Inconsistent error messages
- Future ID type changes require multiple edits

**Fix Approach:**
1. Update GraphQL schema to use IntID scalar for all ID fields in Query/Mutation
2. Run `make graphql-gen` to regenerate resolvers with IntID parameters
3. Remove strconv.Atoi calls from resolvers (gqlgen will handle conversion)

**Priority:** Low - code cleanup, no functional impact

---

### 5. Inconsistent Null Handling for Not Found

**Issue:** Mixed patterns for returning not-found resource.

**Files:**
- `perspectize-go/internal/adapters/graphql/resolvers/schema.resolvers.go` (lines 274, 290, 311) - Returns `nil, nil` for not found
- Line 189-190 - Returns error for ContentByID not found
- Line 141 - Returns error for Perspective update not found

**Current Patterns:**
- UserByID/UserByUsername (line 274, 290): `return nil, nil` (GraphQL convention)
- ContentByID (line 189): `return nil, fmt.Errorf("content not found with ID: %s", id)`
- PerspectiveByID (line 311): `return nil, nil` (GraphQL convention)

**Impact:**
- Inconsistent API behavior - some null returns are silent, others error
- Client code can't distinguish "not found" from null field
- Makes error handling unpredictable

**Fix Approach:**
1. Standardize: nullable fields should return `nil, nil` (GraphQL convention)
2. Update ContentByID to return `nil, nil` instead of error
3. Document pattern in CONVENTIONS.md

**Priority:** Low - behavioral inconsistency

---

## Testing Gaps

### 6. Integration Tests Skip When Database Unavailable

**Issue:** No verification that tests pass against actual PostgreSQL.

**Files:**
- `perspectize-go/test/database/postgres_test.go` (lines 9, 14) - `t.Skip()` guards on DB availability

**Current Pattern:**
```go
if os.Getenv("DATABASE_URL") == "" && !dbAvailable {
    t.Skip("Skipping test - PostgreSQL not available...")
}
```

**Problem:**
- Tests silently skip in CI/local environments without database
- Developers can't verify schema/query changes work without manual DB setup
- Migration issues only discovered in production

**Impact:**
- False confidence in test suite (100% pass with skipped tests)
- Integration bugs slip through review
- Onboarding friction (must run Docker for tests)

**Fix Approach:**
1. Use testcontainers-go to spin up PostgreSQL for each test suite
2. Run migrations automatically in test setup
3. Fail tests if database can't start (vs. silently skip)
4. CI must verify testcontainers work

**Priority:** Medium - affects test reliability

---

### 7. Limited Nested Object Testing

**Issue:** Tests don't verify nested field resolution (User, Content on Perspective).

**Files:**
- `perspectize-go/test/resolvers/content_resolver_test.go` - Content resolver tests exist
- No corresponding perspective resolver test file
- Service tests (`perspective_service_test.go`) don't test nested population

**Impact:**
- Nested field bugs won't be caught by tests
- Refactoring nested resolvers has no test safety net

**Fix Approach:**
1. Create `test/resolvers/perspective_resolver_test.go`
2. Add integration tests that query `perspectives { user { id username } content { id name } }`
3. Verify null handling when user/content not found

**Priority:** Low - missing test coverage

---

## Configuration & Secrets

### 8. Database Credentials Logged in Connection String

**Issue:** Password appears in cleartext logs when database connection fails.

**Files:**
- `perspectize-go/cmd/server/main.go` (lines 35-39) - Conditional logging
- `perspectize-go/internal/config/config.go` (lines 80-84) - GetDSN constructs connection string

**Current Mitigation:**
```go
// Mask credentials in log output
if os.Getenv("DATABASE_URL") != "" {
    log.Println("Connecting to database using DATABASE_URL...")
} else {
    log.Printf("Connecting to database at %s:%d/%s...",
        cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)
}
```

**Remaining Risk:**
- If DATABASE_URL env var is empty but Password is set, connection fails silently
- Error messages from sqlx.Connect may leak password if database connection fails
- Logs at debug/trace level might expose full DSN

**Impact:** Moderate - credential exposure in error logs

**Fix Approach:**
1. Always use DATABASE_URL in production (per CLAUDE.md guidance)
2. Never log full connection string
3. Use structured logging (slog, not fmt.Printf)
4. Mask password in all error messages

**Priority:** Medium - security hardening

---

## Fragile Areas

### 9. YouTube API Client Has Silent Duration Parsing Failure

**Issue:** ISO8601 duration parsing defaults to 0 without logging error.

**Files:**
- `perspectize-go/internal/adapters/youtube/client.go` (lines 90-93)

**Current Code:**
```go
duration, err := ParseISO8601Duration(item.ContentDetails.Duration)
if err != nil {
    duration = 0 // Default to 0 if parsing fails
}
```

**Problem:**
- Parsing errors are silently swallowed
- Content gets saved with length = 0 and lengthUnits = null
- No way to detect if video metadata was incomplete

**Impact:**
- Incorrect video durations in API
- Users can't filter by length (always falls outside min/max ranges)
- Debugging why content appears broken is hard

**Fix Approach:**
1. Return error from GetVideoMetadata if duration parsing fails
2. Log parsing failures with context (videoID, raw duration string)
3. Consider retry logic for YouTube API transients
4. Add test case for non-standard duration formats

**Priority:** Low - data quality issue

---

### 10. JSONB Array Column Type Has Custom Scanner

**Issue:** Complex custom type for jsonb[] columns adds maintenance burden.

**Files:**
- `perspectize-go/internal/adapters/repositories/postgres/perspective_repository.go` (lines 17-44)

**Current Implementation:**
```go
type JSONBArray []string
func (a *JSONBArray) Scan(src interface{}) error { ... }
func (a JSONBArray) Value() (driver.Value, error) { ... }
```

**Problem:**
- Custom scanner wraps pq.StringArray
- Not used elsewhere in codebase (potential dead code)
- If jsonb[] column schema changes, this breaks silently
- Alternative: Store as JSON array in JSONB (simpler)

**Impact:** Low - code smell, no immediate functional impact

**Fix Approach:**
1. Verify where JSONBArray is actually used in schema
2. If only one column uses it, consider standardizing to pure JSONB instead
3. Add type documentation explaining why custom type is needed
4. Add test verifying scanner handles NULL, empty array, and populated array

**Priority:** Low - refactoring candidate

---

## Scaling Considerations

### 11. No Query Complexity Limits

**Issue:** GraphQL queries have no depth or complexity restrictions.

**Files:**
- `perspectize-go/cmd/server/main.go` (lines 74-75) - Server setup

**Current Setup:**
```go
srv := handler.NewDefaultServer(generated.NewExecutableSchema(...))
```

**Missing:**
- Query complexity calculation
- Depth limiting (prevents `perspectives { user { perspectives { user { ... } } } }`)
- Timeout on resolver execution
- Rate limiting

**Potential Attack:**
```graphql
{
  perspectives(first: 1000000) {
    items {
      user {
        perspectives(first: 1000000) { items { ... } }
      }
    }
  }
}
```

**Impact:** DoS vulnerability - malicious queries can exhaust server resources

**Fix Approach:**
1. Add gqlgen complexity analyzer config
2. Set per-query complexity budget (e.g., max 1000 points)
3. Add timeout middleware (ctx.WithTimeout)
4. Add request rate limiting

**Priority:** High - security issue

---

## Missing Operational Features

### 12. No Health Check Endpoint

**Issue:** No /health or readiness probe for Kubernetes/Fly.io deployments.

**Files:**
- `perspectize-go/cmd/server/main.go` - Only /graphql and / endpoints

**Current Routes:**
- `/` - GraphQL playground
- `/graphql` - GraphQL API

**Missing:**
- `/health` - Simple status check
- `/ready` - Database connectivity check

**Impact:**
- Load balancers can't detect unhealthy instances
- Deployments may serve traffic to crashed processes
- Cold start detection impossible

**Fix Approach:**
1. Add HTTP handler for `/health` (returns 200 OK)
2. Add `/ready` that checks database connectivity
3. Wire into main.go before ListenAndServe
4. Document expected response format

**Priority:** Medium - operational necessity

---

### 13. Structured Logging Not Implemented

**Issue:** Standard library `log` package used instead of structured logging (slog).

**Files:**
- `perspectize-go/cmd/server/main.go` - All logs use `log.Println`, `log.Printf`, `log.Fatalf`
- CLAUDE.md (line 238) mentions `slog` should be used

**Current Approach:**
```go
log.Printf("PostgreSQL version: %s", version)
```

**Recommended (from CLAUDE.md):**
```go
slog.Info("PostgreSQL version", "version", version)
```

**Impact:**
- No structured JSON logs for aggregation/analysis in production
- Harder to correlate requests across logs
- Missing context (user ID, request ID, latency)

**Fix Approach:**
1. Initialize slog logger in main.go
2. Replace all log.Print* calls with slog.Info/Error/Debug
3. Add structured fields (version, duration, error context)
4. Configure JSON output for production

**Priority:** Low - doesn't block functionality, but needed for production observability

---

## Documentation Gaps

### 14. No Error Code / Status Reference

**Issue:** GraphQL error responses are ad-hoc error messages.

**Files:**
- `perspectize-go/internal/adapters/graphql/resolvers/schema.resolvers.go` - Errors are fmt.Errorf strings

**Current Errors:**
- "user already exists: %w"
- "invalid input: %w"
- "resource not found"

**Problem:**
- Clients can't reliably parse error type from message
- No standard error codes (like REST status codes)
- Error mapping is inconsistent

**Impact:** Clients must implement fragile string parsing to handle errors

**Fix Approach:**
1. Define error code constants in domain/errors.go
2. Create custom GraphQL error extension format
3. Update resolvers to include error codes in response
4. Document error catalog in CONVENTIONS.md

**Priority:** Low - future enhancement

---

## Summary Table

| Area | Issue | Priority | Blocks |
|------|-------|----------|--------|
| Auth | No authentication middleware | High | Production |
| Nested Fields | User/Content resolvers missing | Medium | GraphQL contract |
| DoS | No query complexity limits | High | Security |
| Testing | DB tests skip silently | Medium | Reliability |
| Operations | No health checks | Medium | Deployment |
| Logging | Using standard log, not slog | Low | Observability |
| Filters | Limited content filtering | Low | UX |
| ID Conversion | String-to-int scattered | Low | Code quality |
| Null Handling | Inconsistent not-found patterns | Low | API consistency |

---

*Concerns audit: 2026-02-04*
