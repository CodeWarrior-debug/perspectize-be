# Phase 12: Authentication & Security — Context

## Phase Goal

Replace user dropdown selector with proper JWT authentication. Secure all GraphQL mutations. Enable user-specific features.

## Problem Statement

From FEATURE_BACKLOG.md:

"The GraphQL client has empty `headers: {}` — no auth tokens, no CSRF protection, no per-user cache scoping."

Current system allows anyone to:
- Create/update/delete perspectives
- Create/update/delete content
- Act as any user (via dropdown selector)

## Research Summary

See `.planning/v1.1-research/AUTH-ARCHITECTURE.md` for full research.

**Recommended approach:** Hybrid JWT with httpOnly refresh cookies
- **Access tokens:** 15 min expiry, stored in-memory (JavaScript variable)
- **Refresh tokens:** 7 days expiry, stored in httpOnly cookie
- **Password hashing:** Argon2id (OWASP 2026 standard)
- **Middleware:** go-chi/jwtauth for JWT validation

**Why not localStorage for tokens:**
- XSS vulnerability — any script can read localStorage
- httpOnly cookies cannot be accessed by JavaScript

**Why hybrid approach:**
- Access tokens in memory: fast access, no cookie overhead
- Refresh tokens in cookies: survive page refresh, secure storage

## Current Architecture

```
Frontend                    Backend
┌─────────────┐            ┌─────────────┐
│ User Select │───────────▶│ No Auth     │
│ Dropdown    │            │ Middleware  │
└─────────────┘            └─────────────┘
     │                           │
     ▼                           ▼
┌─────────────┐            ┌─────────────┐
│ GraphQL     │───────────▶│ Resolvers   │
│ Client      │ userId     │ Trust input │
│ (no auth)   │ in request │             │
└─────────────┘            └─────────────┘
```

## Target Architecture

```
Frontend                    Backend
┌─────────────┐            ┌─────────────┐
│ Auth State  │◀──────────▶│ JWT Auth    │
│ (runes)     │ access     │ Middleware  │
└─────────────┘ token      └─────────────┘
     │                           │
     ▼                           ▼
┌─────────────┐            ┌─────────────┐
│ GraphQL     │───────────▶│ Resolvers   │
│ Client      │ Bearer     │ Check ctx   │
│ + auth hook │ token      │ user        │
└─────────────┘            └─────────────┘
     │
     ▼
┌─────────────┐
│ TanStack    │
│ Query       │
│ (user-keyed)│
└─────────────┘
```

## Database Changes

```sql
-- Add password field to users table
ALTER TABLE users ADD COLUMN password_hash TEXT;

-- Add refresh token storage
CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash TEXT NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    revoked_at TIMESTAMPTZ
);

CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens (user_id);
CREATE INDEX idx_refresh_tokens_hash ON refresh_tokens (token_hash);
```

## GraphQL Schema Changes

```graphql
input RegisterInput {
  username: String!
  password: String!
  email: String
}

input LoginInput {
  username: String!
  password: String!
}

type AuthPayload {
  user: User!
  accessToken: String!
  expiresIn: Int!  # seconds
}

type Mutation {
  register(input: RegisterInput!): AuthPayload!
  login(input: LoginInput!): AuthPayload!
  logout: Boolean!
  refreshToken: AuthPayload!
}
```

## Go Middleware Pattern

```go
// Chi middleware
func JWTAuthMiddleware(secretKey string) func(next http.Handler) http.Handler {
    return jwtauth.Verifier(jwtauth.New("HS256", []byte(secretKey), nil))
}

// Context injection
func AuthContext(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        _, claims, _ := jwtauth.FromContext(r.Context())
        if claims != nil {
            userID := claims["user_id"].(float64)
            ctx := context.WithValue(r.Context(), auth.UserIDKey, int64(userID))
            next.ServeHTTP(w, r.WithContext(ctx))
        } else {
            next.ServeHTTP(w, r)
        }
    })
}

// Resolver access
func (r *mutationResolver) CreatePerspective(ctx context.Context, input model.CreatePerspectiveInput) (*model.Perspective, error) {
    userID := auth.ForContext(ctx)
    if userID == 0 {
        return nil, errors.New("authentication required")
    }
    // ... create perspective with userID
}
```

## SvelteKit Patterns

```typescript
// Auth state with Svelte 5 runes
let accessToken = $state<string | null>(null);
let user = $state<User | null>(null);

// Token refresh
async function refreshToken() {
    const response = await fetch('/api/refresh', { credentials: 'include' });
    const data = await response.json();
    accessToken = data.accessToken;
    user = data.user;
}

// GraphQL client with auth
const graphqlClient = new GraphQLClient(GRAPHQL_ENDPOINT, {
    requestMiddleware: (request) => ({
        ...request,
        headers: {
            ...request.headers,
            ...(accessToken ? { Authorization: `Bearer ${accessToken}` } : {}),
        },
    }),
});

// TanStack Query cache scoping
const queryClient = new QueryClient({
    defaultOptions: {
        queries: {
            queryKeyHashFn: (key) => {
                // Prefix with user ID
                return JSON.stringify([user?.id, ...key]);
            },
        },
    },
});
```

## Requirements Covered

- AUTH-01 through AUTH-13 (all authentication requirements)

## Success Metrics

| Metric | Current | Target |
|--------|---------|--------|
| Auth mechanism | User dropdown | JWT tokens |
| Password storage | N/A | Argon2id |
| Token storage | N/A | httpOnly cookies |
| Mutation protection | None | All mutations |
| Cache scoping | None | User ID prefix |

## Dependencies

- Phase 11 (clean schema foundation)

## Risks

- **Security vulnerabilities:** JWT implementation errors, timing attacks
- **Token leakage:** Improper cookie flags, XSS exposure
- **Migration:** Existing users need password set flow

## Open Questions

1. Should existing users be forced to set password, or grandfather with dropdown?
2. Should we implement "remember me" checkbox (extend refresh to 30 days)?
3. Should we add login attempt rate limiting immediately, or defer?

---

*Context gathered: 2026-02-16*
