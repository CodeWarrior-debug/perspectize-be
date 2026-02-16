# Phase 12: Authentication (Go Backend) - Clerk Integration Research

**Researched:** 2026-02-16
**Domain:** Clerk authentication, Go backend JWT verification, user sync
**Confidence:** HIGH

## Summary

This research covers integrating Clerk as the authentication provider for the Perspectize Go backend. The Clerk Go SDK (`clerk-sdk-go/v2`) is mature, well-maintained, and provides built-in HTTP middleware for JWT verification that works with any `http.Handler`-compatible router including chi. The SDK handles JWKS fetching, key caching, and token verification automatically.

The key architectural decisions are: (1) use Clerk's `WithHeaderAuthorization` middleware on the chi router to validate Bearer tokens, (2) extract `SessionClaims` from context in a custom middleware that maps Clerk user IDs to local database users, (3) use Clerk webhooks (verified via Svix) to sync user creation/updates/deletions to the local database, and (4) add a `clerk_user_id` column to the existing `users` table to link Clerk identities with local records.

**Primary recommendation:** Use `github.com/clerk/clerk-sdk-go/v2` with its built-in HTTP middleware. Do NOT hand-roll JWT verification, JWKS fetching, or token parsing.

## Standard Stack

### Core

| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| `github.com/clerk/clerk-sdk-go/v2` | v2.5.1 | Clerk Backend API client + JWT middleware | Official SDK, handles JWKS caching, token verification, session claims extraction |
| `github.com/clerk/clerk-sdk-go/v2/http` | v2.5.1 | HTTP middleware (WithHeaderAuthorization) | Built-in middleware compatible with any http.Handler router |
| `github.com/clerk/clerk-sdk-go/v2/user` | v2.5.1 | User management API (create, get, list, update) | Backend-initiated user operations |
| `github.com/svix/svix-webhooks/go` | latest | Webhook signature verification | Clerk uses Svix for webhook delivery; this verifies HMAC-SHA256 signatures |

### Supporting

| Library | Version | Purpose | When to Use |
|---------|---------|---------|-------------|
| `github.com/clerk/clerk-sdk-go/v2/session` | v2.5.1 | Session management API | Only if needing to revoke/list sessions server-side |
| `github.com/clerk/clerk-sdk-go/v2/jwks` | v2.5.1 | Manual JWKS client | Only if custom JWT verification is needed (rare) |

### Alternatives Considered

| Instead of | Could Use | Tradeoff |
|------------|-----------|----------|
| Clerk SDK middleware | go-chi/jwtauth + manual JWKS | More control but must hand-roll JWKS fetch, caching, key rotation. Not recommended. |
| Svix Go library | Manual HMAC verification | Svix handles timestamp tolerance, replay protection. Don't hand-roll. |
| Clerk webhook sync | Polling Clerk API | Webhooks are real-time; polling wastes API calls against rate limits (1000/10s prod) |

**Installation:**
```bash
cd backend
go get -u github.com/clerk/clerk-sdk-go/v2
go get -u github.com/svix/svix-webhooks/go
```

## Architecture Patterns

### Recommended Project Structure

```
backend/internal/
├── core/
│   ├── domain/
│   │   └── user.go              # Add ClerkUserID field
│   └── ports/
│       ├── repositories/
│       │   └── user_repository.go  # Add GetByClerkID method
│       └── services/
│           └── auth_service.go     # New: auth port interface
├── adapters/
│   ├── auth/
│   │   ├── clerk_middleware.go     # Chi middleware wrapping Clerk SDK
│   │   ├── context.go             # Context key, ForContext helper
│   │   └── webhook_handler.go     # Clerk webhook endpoint
│   └── repositories/postgres/
│       └── gorm_user_repository.go # Add GetByClerkID impl
├── middleware/                      # (existing pkg/middleware)
└── cmd/server/main.go              # Wire auth middleware
```

### Pattern 1: Clerk Middleware Chain for Chi

**What:** Two-layer middleware: (1) Clerk SDK verifies JWT and writes SessionClaims to context, (2) custom middleware maps Clerk user ID to local DB user.

**When to use:** Every request to `/graphql` endpoint.

```go
// Source: clerk-sdk-go/v2 docs + gqlgen authentication recipe
package auth

import (
    "context"
    "net/http"

    "github.com/clerk/clerk-sdk-go/v2"
    clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
)

// Private context key type prevents collisions (gqlgen pattern)
type contextKey struct{ name string }

var userCtxKey = &contextKey{"user"}

// AuthenticatedUser holds the resolved local user info for resolvers
type AuthenticatedUser struct {
    ID         int
    ClerkID    string
    Username   string
    Email      string
    Role       string
}

// Middleware returns chi-compatible middleware that:
// 1. Validates Clerk JWT via WithHeaderAuthorization
// 2. Resolves local user from Clerk subject claim
func Middleware(userRepo UserRepository) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        // First: Clerk SDK validates JWT, writes SessionClaims to context
        clerkMiddleware := clerkhttp.WithHeaderAuthorization()

        return clerkMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            claims, ok := clerk.SessionClaimsFromContext(r.Context())
            if !ok {
                // No valid session — allow through as unauthenticated
                // (queries may be public, mutations check in resolver)
                next.ServeHTTP(w, r)
                return
            }

            // claims.Subject is the Clerk user ID (e.g., "user_2abc123...")
            clerkUserID := claims.Subject

            // Resolve local user from database
            localUser, err := userRepo.GetByClerkID(r.Context(), clerkUserID)
            if err != nil {
                // User exists in Clerk but not locally — webhook may not have fired yet
                // Return 401 or create on-demand (see strategy section)
                http.Error(w, "user not found", http.StatusUnauthorized)
                return
            }

            // Inject local user into context
            ctx := context.WithValue(r.Context(), userCtxKey, localUser)
            next.ServeHTTP(w, r.WithContext(ctx))
        }))
    }
}

// ForContext extracts authenticated user from context.
// Returns nil if unauthenticated (public query).
func ForContext(ctx context.Context) *AuthenticatedUser {
    raw := ctx.Value(userCtxKey)
    if raw == nil {
        return nil
    }
    return raw.(*AuthenticatedUser)
}
```

### Pattern 2: Resolver Authorization Check

**What:** Resolvers call `auth.ForContext(ctx)` to get the authenticated user. Mutations require auth; queries may be public.

**When to use:** Every resolver that needs user identity.

```go
// Source: gqlgen authentication recipe
func (r *mutationResolver) CreatePerspective(ctx context.Context, input model.CreatePerspectiveInput) (*model.Perspective, error) {
    user := auth.ForContext(ctx)
    if user == nil {
        return nil, fmt.Errorf("authentication required")
    }
    // Use user.ID (local DB ID) for creating the perspective
    return r.perspectiveService.Create(ctx, user.ID, input)
}
```

### Pattern 3: Webhook Handler for User Sync

**What:** HTTP endpoint that receives Clerk webhooks, verifies Svix signature, and syncs user data to local DB.

**When to use:** Dedicated endpoint at `/webhooks/clerk`.

```go
// Source: svix-webhooks/go docs + Clerk webhook docs
package auth

import (
    "encoding/json"
    "io"
    "net/http"

    svix "github.com/svix/svix-webhooks/go"
)

type WebhookHandler struct {
    webhookSecret string
    userRepo      UserRepository
}

func (h *WebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // Read raw body (MUST use raw bytes for signature verification)
    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "bad request", http.StatusBadRequest)
        return
    }

    // Verify Svix signature
    wh, err := svix.NewWebhook(h.webhookSecret)
    if err != nil {
        http.Error(w, "internal error", http.StatusInternalServerError)
        return
    }

    err = wh.Verify(body, r.Header)
    if err != nil {
        http.Error(w, "invalid signature", http.StatusUnauthorized)
        return
    }

    // Parse event
    var event struct {
        Type string          `json:"type"`
        Data json.RawMessage `json:"data"`
    }
    if err := json.Unmarshal(body, &event); err != nil {
        http.Error(w, "bad payload", http.StatusBadRequest)
        return
    }

    switch event.Type {
    case "user.created":
        h.handleUserCreated(r.Context(), event.Data)
    case "user.updated":
        h.handleUserUpdated(r.Context(), event.Data)
    case "user.deleted":
        h.handleUserDeleted(r.Context(), event.Data)
    }

    w.WriteHeader(http.StatusOK)
}
```

### Pattern 4: Wiring in main.go

**What:** Register auth middleware and webhook endpoint in chi router.

```go
// In cmd/server/main.go

import (
    "github.com/clerk/clerk-sdk-go/v2"
    "github.com/CodeWarrior-debug/perspectize/backend/internal/adapters/auth"
)

func main() {
    // ... existing setup ...

    // Initialize Clerk
    clerk.SetKey(os.Getenv("CLERK_SECRET_KEY"))

    // Setup chi router
    r := chi.NewRouter()

    // Existing middleware
    r.Use(middleware.RequestID)
    r.Use(middleware.RealIP)
    r.Use(perfmw.RequestTimer)
    r.Use(perfmw.Recoverer)

    // CORS — MUST allow Authorization header for Bearer tokens
    r.Use(corsMiddleware) // Updated to include "Authorization" in allowed headers

    // Auth middleware — runs on all routes, permissive (doesn't block unauthenticated)
    r.Use(auth.Middleware(userRepo))

    // Webhook endpoint — NO auth middleware (Clerk calls this, not users)
    r.Post("/webhooks/clerk", webhookHandler)

    // GraphQL
    r.Handle("/graphql", srv)
    // ...
}
```

### Anti-Patterns to Avoid

- **Storing Clerk session tokens in DB:** Session tokens are short-lived (60 seconds). Never store them. Verify on every request.
- **Using Clerk user ID as primary key:** Keep local integer IDs. Store Clerk ID as a separate indexed column. This avoids coupling your entire DB schema to Clerk.
- **Blocking on webhook for user creation:** Webhooks are async. If a user authenticates before the `user.created` webhook fires, handle it gracefully (create on-demand or return a "pending" state).
- **Putting auth check in middleware instead of resolvers:** The middleware should be permissive (allow unauthenticated through). Individual resolvers decide what requires auth. This preserves public queries.

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| JWT verification | Custom JWT parsing + JWKS fetching | `clerkhttp.WithHeaderAuthorization()` | Handles JWKS caching, key rotation, clock skew (leeway), token format v2 |
| JWKS key caching | In-memory cache with TTL | Clerk SDK built-in JWKS client | SDK handles refresh intervals, concurrent access, error recovery |
| Webhook signature verification | Manual HMAC-SHA256 | `svix.NewWebhook(secret).Verify()` | Handles timestamp tolerance, replay protection, multiple signature formats |
| Session management | Custom session store | Clerk manages sessions | Clerk handles session lifecycle, multi-device, revocation |
| Password hashing | Argon2id/bcrypt implementation | Clerk handles all password auth | Clerk manages password policies, breach detection, rate limiting |
| Token refresh | Custom refresh token flow | Clerk frontend SDK auto-refreshes | 60-second token lifetime, frontend SDK handles refresh transparently |

**Key insight:** The CONTEXT.md describes building custom JWT auth with Argon2id, refresh tokens, and password hashing. With Clerk, ALL of this is unnecessary. Clerk handles authentication entirely. The backend only needs to verify tokens and sync user data.

## Common Pitfalls

### Pitfall 1: CORS Missing Authorization Header

**What goes wrong:** Frontend sends `Authorization: Bearer <token>` but backend CORS config only allows `Content-Type` header. Browser blocks the request with a preflight error.
**Why it happens:** Current CORS middleware in main.go has `Access-Control-Allow-Headers: Content-Type` only.
**How to avoid:** Update CORS to include `Authorization`:
```go
w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
```
**Warning signs:** Preflight OPTIONS requests succeed but actual requests fail with CORS errors in browser console.

### Pitfall 2: Webhook Body Already Read

**What goes wrong:** Svix signature verification fails because the body was already consumed by another middleware or JSON decoder.
**Why it happens:** `r.Body` is an `io.ReadCloser` -- once read, it's gone. Signature verification requires the exact raw bytes.
**How to avoid:** Read body with `io.ReadAll(r.Body)` FIRST, verify signature, THEN unmarshal JSON. Do NOT use `json.NewDecoder(r.Body)` before verification.
**Warning signs:** Webhook verification always fails even with correct secret.

### Pitfall 3: Race Between Auth and Webhook

**What goes wrong:** User signs up via Clerk frontend, immediately makes a GraphQL request, but the `user.created` webhook hasn't fired yet. Middleware can't find local user.
**Why it happens:** Webhooks are asynchronous; there's a delay between Clerk user creation and webhook delivery.
**How to avoid:** Implement "create on demand" fallback: if Clerk JWT is valid but local user not found, call Clerk Backend API to fetch user details and create local record immediately.
```go
if localUser == nil {
    clerkUser, err := clerkuser.Get(ctx, clerkUserID)
    if err == nil {
        localUser = createLocalUser(ctx, clerkUser)
    }
}
```
**Warning signs:** New users get 401 errors on their first request after signup.

### Pitfall 4: Token Expiry Confusion (60-Second Lifetime)

**What goes wrong:** Developer thinks tokens last 15 minutes (like custom JWTs) and doesn't understand why users get logged out.
**Why it happens:** Clerk session tokens expire every 60 seconds by default. The Clerk frontend SDK auto-refreshes them.
**How to avoid:** Understand the flow: frontend SDK handles refresh automatically. Backend just verifies whatever token arrives. The SDK's `Leeway()` option accommodates clock skew.
**Warning signs:** Intermittent 401 errors, especially on slow connections.

### Pitfall 5: Using String Clerk IDs as Foreign Keys

**What goes wrong:** Changing from integer `user_id` foreign keys to string Clerk IDs throughout the schema. Major migration, breaks existing queries.
**Why it happens:** Developer tries to use Clerk's `user_2abc123` IDs directly.
**How to avoid:** Keep existing integer IDs. Add `clerk_user_id TEXT UNIQUE` column to users table. Middleware resolves Clerk ID to local integer ID once per request.
**Warning signs:** Schema changes cascading across all tables with user references.

### Pitfall 6: Email Not in Default Session Token

**What goes wrong:** Developer expects email in JWT claims and tries to extract it. It's not there.
**Why it happens:** Clerk's default session token (v2) only contains: `sub`, `sid`, `azp`, `exp`, `iat`, `iss`, `nbf`, `jti`, `v`, `fva`. Email is NOT included by default.
**How to avoid:** Either: (a) add email as a custom claim in Clerk Dashboard JWT template, or (b) fetch email from local DB user record (populated via webhook sync). Option (b) is recommended to avoid token bloat.
**Warning signs:** `claims.Custom` is empty or doesn't contain email fields.

## Code Examples

### Clerk SDK Initialization

```go
// Source: https://pkg.go.dev/github.com/clerk/clerk-sdk-go/v2
import "github.com/clerk/clerk-sdk-go/v2"

// Set globally (recommended for single-key apps)
clerk.SetKey(os.Getenv("CLERK_SECRET_KEY"))
```

### Session Claims Access

```go
// Source: https://pkg.go.dev/github.com/clerk/clerk-sdk-go/v2
claims, ok := clerk.SessionClaimsFromContext(r.Context())
if !ok {
    // No valid session
    return
}

// Available fields (v2 token):
userID := claims.Subject                    // "user_2abc123..." (Clerk user ID)
sessionID := claims.SessionID              // "sess_xyz..."
authorizedParty := claims.AuthorizedParty  // Your frontend origin
orgID := claims.ActiveOrganizationID       // Empty if no org active
orgRole := claims.ActiveOrganizationRole   // Empty if no org active
```

### SessionClaims Struct (from SDK)

```go
// Source: https://pkg.go.dev/github.com/clerk/clerk-sdk-go/v2
type SessionClaims struct {
    RegisteredClaims                          // sub, exp, iat, iss, nbf, jti
    Version                       int         `json:"v"`
    SessionID                     string      `json:"sid"`
    AuthorizedParty               string      `json:"azp"`
    ActiveOrganizationID          string      `json:"org_id"`
    ActiveOrganizationSlug        string      `json:"org_slug"`
    ActiveOrganizationRole        string      `json:"org_role"`
    ActiveOrganizationPermissions []string    `json:"org_permissions"`
    Actor                         json.RawMessage `json:"act,omitempty"`
    FactorVerificationAge         [2]int64    `json:"fva"`
}
```

### Backend User Management

```go
// Source: https://pkg.go.dev/github.com/clerk/clerk-sdk-go/v2/user
import clerkuser "github.com/clerk/clerk-sdk-go/v2/user"

// Fetch user details from Clerk
clerkUser, err := clerkuser.Get(ctx, "user_2abc123")

// Create user in Clerk (for migration)
newUser, err := clerkuser.Create(ctx, &clerkuser.CreateParams{
    EmailAddresses: &[]string{"user@example.com"},
    Username:       clerk.String("jsmith"),
    Password:       clerk.String("securepassword"),
})

// List all users
users, err := clerkuser.List(ctx, &clerkuser.ListParams{})
```

### Webhook Verification

```go
// Source: https://www.svix.com/guides/receiving/receive-webhooks-with-go/
import svix "github.com/svix/svix-webhooks/go"

wh, err := svix.NewWebhook(os.Getenv("CLERK_WEBHOOK_SIGNING_SECRET"))
if err != nil {
    return fmt.Errorf("invalid webhook secret: %w", err)
}

// body must be raw bytes, headers from http.Request
err = wh.Verify(body, r.Header)
if err != nil {
    return fmt.Errorf("webhook verification failed: %w", err)
}
```

### Database Migration for Clerk Integration

```sql
-- Add Clerk user ID to existing users table
ALTER TABLE users ADD COLUMN clerk_user_id TEXT UNIQUE;

-- Index for fast lookup by Clerk ID (middleware does this on every request)
CREATE INDEX idx_users_clerk_user_id ON users (clerk_user_id);
```

### Updated Domain Model

```go
// backend/internal/core/domain/user.go
type User struct {
    ID         int
    ClerkUserID string    // NEW: Clerk's "user_2abc..." identifier
    Username   string
    Email      string
    Role       UserRole
    Active     bool
    CreatedAt  time.Time
    UpdatedAt  time.Time
}
```

### Updated CORS Middleware

```go
// Must allow Authorization header for Bearer tokens
r.Use(func(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        origin := r.Header.Get("Origin")
        if origin == "" {
            origin = "*"
        }
        w.Header().Set("Access-Control-Allow-Origin", origin) // Or specific frontend URL
        w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        w.Header().Set("Access-Control-Allow-Credentials", "true")
        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusOK)
            return
        }
        next.ServeHTTP(w, r)
    })
})
```

## JWT Claims Reference

### Default Clerk Session Token (v2)

| Claim | Type | Description | Always Present |
|-------|------|-------------|----------------|
| `sub` | string | Clerk user ID (`user_2abc...`) | Yes |
| `sid` | string | Session ID (`sess_xyz...`) | Yes |
| `azp` | string | Authorized party (frontend origin) | Yes |
| `exp` | int64 | Expiration (Unix timestamp, ~60s from issue) | Yes |
| `iat` | int64 | Issued at (Unix timestamp) | Yes |
| `iss` | string | Issuer (Clerk instance URL) | Yes |
| `nbf` | int64 | Not before (Unix timestamp) | Yes |
| `jti` | string | JWT ID (unique per token) | Yes |
| `v` | int | Token version (2) | Yes |
| `fva` | [2]int64 | Factor verification age | Yes |
| `org_id` | string | Active organization ID | Only if org active |
| `org_slug` | string | Active organization slug | Only if org active |
| `org_role` | string | User's role in active org | Only if org active |

**NOT in default token:** email, name, username, metadata. These must be synced via webhooks or added as custom claims.

## User Sync Strategy

### Recommended: Webhook Sync + On-Demand Fallback

1. **Primary:** Clerk webhooks sync user data to local DB asynchronously
2. **Fallback:** If authenticated request arrives but local user not found, fetch from Clerk API and create locally
3. **Local ID:** Keep existing integer `users.id` as primary key for all foreign key references
4. **Link column:** Add `users.clerk_user_id TEXT UNIQUE` for Clerk-to-local mapping

### Webhook Events to Subscribe

| Event | Action |
|-------|--------|
| `user.created` | Create local user with Clerk ID, email, username |
| `user.updated` | Update local user's email, username, active status |
| `user.deleted` | Soft-delete or reassign to `[deleted]` sentinel user (existing pattern) |

### Migration of Existing Users

Two options for linking existing local users to Clerk:

**Option A: Invite existing users to Clerk (recommended)**
1. Create Clerk users via Backend API with matching emails
2. Set `clerk_user_id` on local records
3. Users receive invite email to set password / configure auth

**Option B: Use Clerk's ExternalID**
1. Set `external_id` on Clerk user to match local integer ID
2. Look up by external ID if clerk_user_id not set yet

## Rate Limiting

**Clerk Backend API limits (as of July 2025):**
- Production: 1000 requests per 10 seconds
- Development: 100 requests per 10 seconds
- Returns `429 Too Many Requests` with `Retry-After` header when exceeded

**Frontend API limits (per IP):**
- SignIn/SignUp creation: 5 per 10 seconds
- Verification attempts: 3 per 10 seconds

**Implication:** You still need your own rate limiting for your GraphQL endpoint. Clerk only rate-limits its own API, not your backend.

## Error Handling in GraphQL Resolvers

```go
// Pattern for GraphQL resolver auth errors
func (r *mutationResolver) UpdatePerspective(ctx context.Context, id int, input model.UpdatePerspectiveInput) (*model.Perspective, error) {
    user := auth.ForContext(ctx)
    if user == nil {
        // Return GraphQL error, not HTTP error
        return nil, fmt.Errorf("authentication required")
    }

    // Check ownership
    perspective, err := r.perspectiveService.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }
    if perspective.UserID != user.ID {
        return nil, fmt.Errorf("not authorized to modify this perspective")
    }

    return r.perspectiveService.Update(ctx, id, input)
}
```

**Error scenarios:**
| Scenario | Middleware Behavior | Resolver Behavior |
|----------|--------------------|--------------------|
| No token | Passes through (unauthenticated) | Returns "authentication required" for mutations |
| Expired token | `SessionClaimsFromContext` returns `ok=false` | Same as no token |
| Valid token, no local user | Creates user on-demand or returns 401 | N/A (handled in middleware) |
| Valid token, valid user | Injects user into context | Proceeds normally |
| Revoked session | Clerk SDK rejects token | Same as expired |

## Environment Variables

```bash
# Required
CLERK_SECRET_KEY=sk_live_xxx          # Clerk Backend API secret key
CLERK_WEBHOOK_SIGNING_SECRET=whsec_xxx  # From Clerk Dashboard > Webhooks

# Optional
CLERK_PUBLISHABLE_KEY=pk_live_xxx     # Only needed if backend serves frontend
```

## State of the Art

| Old Approach (CONTEXT.md) | Current Approach (Clerk) | Impact |
|---------------------------|--------------------------|--------|
| Custom JWT with HS256 | Clerk RS256 JWTs with JWKS | No secret key management for tokens |
| Argon2id password hashing | Clerk handles passwords | Remove password_hash column, no crypto code |
| Refresh token table | Clerk manages sessions | No refresh_tokens table needed |
| go-chi/jwtauth middleware | clerk-sdk-go/v2/http middleware | Built-in JWKS caching, v2 token support |
| Custom login/register mutations | Clerk frontend components | Remove auth GraphQL mutations |
| Manual CORS for cookies | CORS for Bearer tokens only | Simpler CORS (no credentials cookie concern for auth) |

**The CONTEXT.md plan for custom auth is superseded by Clerk integration.** Database changes from CONTEXT.md (password_hash column, refresh_tokens table) are NOT needed.

## Open Questions

1. **Custom claims vs webhook sync for email:**
   - What we know: Email is NOT in default Clerk session tokens. Can be added via custom claims (Clerk Dashboard) or fetched from local DB (synced via webhook).
   - What's unclear: Performance trade-off of larger tokens vs DB lookup on every request.
   - Recommendation: Use webhook sync. Email is already stored locally. Avoids token bloat.

2. **Existing users migration path:**
   - What we know: 3 existing users (including sentinel users `[deleted]`, `[system]`). Clerk Backend API supports creating users programmatically.
   - What's unclear: Whether to migrate sentinel users to Clerk or keep them local-only.
   - Recommendation: Sentinel users remain local-only (they never authenticate). Only real users get Clerk accounts.

3. **Authorized parties (azp) validation:**
   - What we know: Clerk SDK supports `AuthorizedPartyMatches()` option to validate which origins can generate tokens.
   - What's unclear: Exact frontend URLs for production vs development.
   - Recommendation: Configure `azp` validation in production. Skip in development for flexibility.

4. **GraphQL Playground access in development:**
   - What we know: Currently playground is open in non-production. Auth middleware will affect it.
   - What's unclear: How to allow unauthenticated playground access while testing auth.
   - Recommendation: Auth middleware is permissive (allows unauthenticated through). Playground works for public queries. For mutations, use Clerk dev tokens.

## Sources

### Primary (HIGH confidence)
- [clerk/clerk-sdk-go GitHub](https://github.com/clerk/clerk-sdk-go) - SDK structure, version, Go requirements
- [clerk-sdk-go v2 pkg.go.dev](https://pkg.go.dev/github.com/clerk/clerk-sdk-go/v2) - SessionClaims struct, API types, v2.5.1
- [clerk-sdk-go/v2/http pkg.go.dev](https://pkg.go.dev/github.com/clerk/clerk-sdk-go/v2/http) - WithHeaderAuthorization, RequireHeaderAuthorization, options
- [clerk-sdk-go/v2/user pkg.go.dev](https://pkg.go.dev/github.com/clerk/clerk-sdk-go/v2/user) - CreateParams, ListParams, UpdateParams
- [Clerk Session Tokens docs](https://clerk.com/docs/guides/sessions/session-tokens) - Token v2 claims, 60s lifetime, custom claims
- [Clerk Manual JWT Verification](https://clerk.com/docs/backend-requests/manual-jwt) - JWKS endpoints, RS256, azp claim
- [Clerk Cross-Origin Requests](https://clerk.com/docs/backend-requests/making/cross-origin) - Bearer token in Authorization header
- [Clerk Go Session Verification](https://clerk.com/docs/references/go/verifying-sessions) - Middleware usage patterns
- [Svix Go Webhook Guide](https://www.svix.com/guides/receiving/receive-webhooks-with-go/) - Go verification code
- [gqlgen Authentication Recipe](https://gqlgen.com/recipes/authentication/) - Context key pattern, ForContext helper

### Secondary (MEDIUM confidence)
- [Clerk Webhooks Overview](https://clerk.com/docs/guides/development/webhooks/overview) - Event types, payload structure, Svix integration
- [Clerk Rate Limits](https://clerk.com/docs/backend-requests/resources/rate-limits) - 1000/10s production, 100/10s development
- [Clerk Sync Data with Webhooks](https://clerk.com/docs/webhooks/sync-data) - User sync pattern, event types

### Tertiary (LOW confidence)
- None. All findings verified with primary sources.

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH - Official SDK, verified on pkg.go.dev with version numbers
- Architecture: HIGH - Based on official SDK middleware patterns + established gqlgen auth recipe
- Pitfalls: HIGH - Documented in official docs (token lifetime, CORS, claims structure)
- User sync: MEDIUM - Webhook + on-demand pattern is well-documented but specific Go examples are sparse
- Migration: MEDIUM - Clerk API supports user creation, but exact migration workflow depends on project specifics

**Research date:** 2026-02-16
**Valid until:** 2026-03-16 (stable SDK, v2 is current)
