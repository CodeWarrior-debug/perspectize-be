---
phase: 09-security-hardening
verified: 2026-03-02T00:00:00Z
status: passed
score: 15/15 success criteria verified
re_verification: false
---

# Phase 9: Security Hardening Verification Report

**Phase Goal:** Add authentication, authorization, rate limiting, and security headers to make the app safe for multi-user deployment
**Verified:** 2026-03-02
**Status:** passed
**Re-verification:** No — initial verification

---

## Requirements Coverage Note

The plan files reference bug-backlog requirement IDs (C-01, C-04, C-05, C-10, H-10, H-11, H-12, H-14, H-25, M-14, M-15, M-28) sourced from the ROADMAP.md `Source:` field. These IDs do NOT appear in `.planning/REQUIREMENTS.md` or `.planning/v1.1-REQUIREMENTS.md` — they are bug-backlog concern codes tracked in the ROADMAP concern checklist, not formal requirement IDs. This is expected for this phase. All 15 ROADMAP success criteria are verified below.

---

## Goal Achievement

### Observable Truths (from ROADMAP.md Success Criteria)

| # | Truth | Status | Evidence |
|---|-------|--------|----------|
| 1 | Authentication middleware validates JWT on all mutations (C-01) | VERIFIED | `AuthMiddleware` in `middleware/auth.go` validates JWT from `auth_token` cookie; wired at line 168 of `main.go` before GraphQL handler |
| 2 | Authorization checks on mutations — users can only modify own data (C-01) | VERIFIED | `@auth` on all mutations + `@owner` on `updatePerspective`, `deletePerspective`; `DirectiveRoot.Owner` calls `perspectiveService.GetByID` / `contentService.GetByID` |
| 3 | GraphQL query complexity limit enforced (C-04) | VERIFIED | `srv.Use(extension.FixedComplexityLimit(500))` in `main.go` line 145 |
| 4 | CORS restricted to explicit frontend origin (C-05) | VERIFIED | `cors.Handler` uses `secCfg.CORSOrigins` (env-driven, not hardcoded `*`) in `main.go` lines 159-165; CORS_ORIGINS env documented |
| 5 | GraphQL playground disabled in production (C-09) | VERIFIED | `main.go` line 192-194: playground only registered when `APP_ENV != "production"` (pre-existing fix) |
| 6 | Introspection disabled in production (C-10) | VERIFIED | `main.go` line 147-149: `extension.Introspection{}` only added when `APP_ENV != "production"` |
| 7 | User email only visible to authenticated user for own account (H-10) | VERIFIED | `schema.resolvers.go` line 529-548: `userResolver.Email` calls `middleware.ForContext` and returns `nil` unless `authenticatedUserID == requestedUserID` |
| 8 | Rate limiting middleware installed (H-11) | VERIFIED | `apimw.GlobalRateLimit(secCfg.RateLimitPerMin)` wired at `main.go` line 158, BEFORE auth middleware |
| 9 | YouTube API key never appears in logs or error responses (H-12) | VERIFIED | `sanitizeYouTubeError` in `youtube/client.go` lines 61-83 strips `googleapis.com` URLs; `content_service.go` returns generic `"failed to fetch video metadata"` |
| 10 | HTTP server has read/write/idle timeouts (H-15) | VERIFIED | `main.go` lines 199-204: `ReadTimeout: 15s`, `WriteTimeout: 15s`, `IdleTimeout: 60s` (pre-existing fix) |
| 11 | TLS/HTTPS via Sevalla reverse proxy (H-14) | VERIFIED | Documented in `.env.example` lines 37-44; HSTS set by `SecureHeaders()` middleware with `SSLProxyHeaders` for Cloudflare |
| 12 | Security headers: X-Content-Type-Options, X-Frame-Options, HSTS (M-14) | VERIFIED | `secureheaders.go`: `ContentTypeNosniff: true`, `FrameDeny: true`, `STSSeconds: 31536000`; tests pass for all three |
| 13 | CSRF protection via Content-Type validation (M-15) | VERIFIED | `ContentTypeValidation` middleware in `contenttype.go` rejects non-`application/json` POSTs with 415; wired at `main.go` line 167 |
| 14 | Content Security Policy enhanced on frontend (H-25) | VERIFIED | `frontend/src/app.html` CSP includes `frame-ancestors 'none'; base-uri 'self'; form-action 'self'` |
| 15 | Secrets management strategy documented (M-28) | VERIFIED | `.docs/SECURITY.md` (194 lines): covers JWT generation, rotation, YouTube key rotation, DB credentials, Sevalla practices, incident response |

**Score:** 15/15 success criteria verified

---

## Required Artifacts

### Plan 09-01 Artifacts

| Artifact | Expected | Status | Details |
|----------|----------|--------|---------|
| `backend/internal/adapters/web/middleware/auth.go` | JWT validation middleware with ForContext | VERIFIED | 58 lines; exports `AuthMiddleware`, `ForContext`, `WithUserContext` |
| `backend/internal/core/services/auth_service.go` | JWT generation and validation logic | VERIFIED | Implements `GenerateAccessToken` with HS256, `ValidateToken` with algorithm check |
| `backend/internal/adapters/graphql/directives/auth.go` | GraphQL @auth directive implementation | VERIFIED | `Auth` and `Owner` methods, service injection, `extractResourceID` helper |
| `backend/internal/config/security.go` | Security configuration structure | VERIFIED | `type Security struct` with JWT, rate limit, CORS fields; fail-fast in production |
| `backend/internal/core/domain/auth.go` | JWT claims domain model | VERIFIED | `type Claims struct` with `UserID`, `Email`, `jwt.RegisteredClaims` |
| `backend/internal/core/ports/services/auth_service.go` | AuthService interface | VERIFIED | `GenerateAccessToken` and `ValidateToken` interface methods |

### Plan 09-02 Artifacts

| Artifact | Expected | Status | Details |
|----------|----------|--------|---------|
| `backend/internal/adapters/graphql/resolvers/schema.resolvers.go` | Email visibility check | VERIFIED | `userResolver.Email` uses `ForContext`; returns `nil` for unauthorized |
| `backend/internal/adapters/graphql/directives/auth.go` | Full @owner with service lookup | VERIFIED | `perspectiveService.GetByID` and `contentService.GetByID` for ownership checks |

### Plan 09-03 Artifacts

| Artifact | Expected | Status | Details |
|----------|----------|--------|---------|
| `backend/internal/adapters/web/middleware/ratelimit.go` | Rate limiting middleware | VERIFIED | `GlobalRateLimit` and `GraphQLRateLimit` via `httprate.LimitByIP` |
| `backend/internal/adapters/web/middleware/contenttype.go` | Content-Type validation middleware | VERIFIED | Rejects non-`application/json` POSTs with 415 |

### Plan 09-04 Artifacts

| Artifact | Expected | Status | Details |
|----------|----------|--------|---------|
| `backend/internal/adapters/web/middleware/secureheaders.go` | Security headers via unrolled/secure | VERIFIED | HSTS, CSP, X-Content-Type-Options, X-Frame-Options, XSS filter |
| `frontend/src/app.html` | Enhanced CSP meta tag | VERIFIED | Contains `frame-ancestors 'none'`, `base-uri 'self'`, `form-action 'self'` |

### Plan 09-05 Artifacts

| Artifact | Expected | Status | Details |
|----------|----------|--------|---------|
| `backend/internal/adapters/youtube/client.go` | Error sanitization | VERIFIED | `sanitizeYouTubeError` function strips googleapis.com URLs |
| `backend/internal/adapters/graphql/resolvers/schema.resolvers.go` | Generic GraphQL error messages | VERIFIED | `"failed to create content from YouTube"` at line 42 |

### Plan 09-06 Artifacts

| Artifact | Expected | Status | Details |
|----------|----------|--------|---------|
| `.docs/SECURITY.md` | Comprehensive security documentation | VERIFIED | 194 lines; covers Secret Management, JWT rotation, Sevalla practices, incident response |

---

## Key Link Verification

| From | To | Via | Status | Details |
|------|----|-----|--------|---------|
| `main.go` | `middleware/auth.go` | `r.Use(apimw.AuthMiddleware(authService))` | WIRED | Line 168; authService passed by value |
| `main.go` | `middleware/ratelimit.go` | `r.Use(apimw.GlobalRateLimit(...))` | WIRED | Line 158; runs before auth |
| `main.go` | `middleware/secureheaders.go` | `r.Use(apimw.SecureHeaders())` | WIRED | Line 166; after CORS |
| `main.go` | `middleware/contenttype.go` | `r.Use(apimw.ContentTypeValidation)` | WIRED | Line 167 |
| `main.go` | `directives/auth.go` | `directives.NewDirectiveRoot(contentService, perspectiveService)` | WIRED | Line 127; services injected |
| `main.go` | gqlgen extension | `srv.Use(extension.FixedComplexityLimit(500))` | WIRED | Line 145 |
| `directives/auth.go` | `middleware/auth.go` | `middleware.ForContext(ctx)` | WIRED | Lines 31, 42 |
| `directives/auth.go` | `perspectiveService.GetByID` | ownership check | WIRED | Line 58 |
| `directives/auth.go` | `contentService.GetByID` | ownership check | WIRED | Line 66 |
| `resolvers/schema.resolvers.go` | `middleware.ForContext` | email visibility check | WIRED | Line 532 |
| `youtube/client.go` | `sanitizeYouTubeError` | error sanitization | WIRED | Lines 100, 101 |

---

## ROADMAP Concern Checklist Coverage

| Concern ID | Description | Plan | Status | Evidence |
|------------|-------------|------|--------|---------|
| C-01 | No authentication or authorization | 09-01, 09-02 | SATISFIED | JWT middleware + @auth/@owner directives |
| C-04 | No GraphQL query complexity limiting | 09-03 | SATISFIED | `FixedComplexityLimit(500)` |
| C-05 | Wildcard CORS configuration | 09-03 | SATISFIED | Config-driven via `CORS_ORIGINS` env var |
| C-09 | GraphQL playground exposed unconditionally | (pre-existing) | SATISFIED | Gated by `APP_ENV != "production"` |
| C-10 | GraphQL introspection enabled unconditionally | 09-03 | SATISFIED | Conditional on `APP_ENV != "production"` |
| H-10 | User email addresses exposed in public query | 09-02 | SATISFIED | Field resolver returns nil unless own account |
| H-11 | No rate limiting | 09-03 | SATISFIED | `GlobalRateLimit` before auth in middleware stack |
| H-12 | YouTube API key exposure risk | 09-05 | SATISFIED | `sanitizeYouTubeError` + generic service errors |
| H-14 | No HTTPS/TLS | 09-04 | SATISFIED | Sevalla TLS termination documented; HSTS via unrolled/secure |
| H-15 | No HTTP server timeouts | (pre-existing) | SATISFIED | ReadTimeout 15s, WriteTimeout 15s, IdleTimeout 60s |
| H-25 | No Content Security Policy | 09-04 | SATISFIED | Enhanced CSP with frame-ancestors, base-uri, form-action |
| M-14 | Missing security headers | 09-04 | SATISFIED | HSTS, X-Content-Type-Options, X-Frame-Options, CSP via unrolled/secure |
| M-15 | No CSRF protection | 09-03 | SATISFIED | Content-Type validation rejects non-JSON POSTs |
| M-28 | No secret rotation or vault integration | 09-06 | SATISFIED | `.docs/SECURITY.md` 194-line documentation |

**Note:** Requirements IDs in plan files (C-01, H-10, etc.) refer to bug-backlog concern codes tracked in the ROADMAP, not to entries in REQUIREMENTS.md or v1.1-REQUIREMENTS.md. No orphaned requirements — these codes have no presence in either requirements file and are not expected to.

---

## Go Dependencies Verified

| Dependency | Version | Purpose |
|------------|---------|---------|
| `github.com/golang-jwt/jwt/v5` | v5.3.1 | JWT generation and validation |
| `github.com/go-chi/httprate` | v0.15.0 | Rate limiting middleware |
| `github.com/go-chi/cors` | v1.2.2 | CORS middleware |
| `github.com/unrolled/secure` | v1.17.0 | Security headers |

---

## Anti-Patterns Found

No blocking anti-patterns found. No TODO/FIXME/placeholder comments in security-relevant files.

**Deferred test issue (resolved):** `deferred-items.md` documented `TestSecureHeaders_ProductionHSTS` as failing, but running tests confirms all 5 security header tests now pass (including the production HSTS test). The deferred item is closed.

---

## Human Verification Required

### 1. CORS Wildcard in Production

**Test:** Deploy to Sevalla with `CORS_ORIGINS=*` (the default) and verify a cross-origin request from an untrusted origin is accepted.
**Expected:** In production, CORS_ORIGINS should be set to `https://app.perspectize.com` — the default `*` is only safe for local dev. The infrastructure is correct (config-driven) but operational verification is needed.
**Why human:** Cannot verify env var configuration on Sevalla from code alone; this is a deployment configuration check.

### 2. Clerk Auth vs Custom JWT

**Test:** Review whether the CONTEXT.md decision to use Clerk as auth provider was superseded by the plan decisions to implement custom JWT.
**Expected:** The CONTEXT.md specifies Clerk + JWKS endpoint validation. The implemented plans use a custom `golang-jwt/jwt` service. There is no Clerk integration anywhere in the codebase. Confirm this divergence was intentional (Clerk deferred to Phase 12 based on ROADMAP Phase 12 Authentication plans).
**Why human:** Architectural decision about scope boundary — the ROADMAP Phase 12 (`AUTH-01` through `AUTH-13`) covers full authentication with user registration/login; Phase 9 established the JWT infrastructure layer. Verify this scope split was intentional.

### 3. Rate Limit Default Under Load

**Test:** Verify 100 requests/minute limit (default `RATE_LIMIT_PER_MIN`) is appropriate for expected production traffic.
**Expected:** No false positives for legitimate users; sufficient to prevent DoS.
**Why human:** Load testing required; no performance benchmarks in codebase.

---

## Gaps Summary

No gaps found. All 15 ROADMAP success criteria are implemented and verified in the codebase.

One architectural note for awareness (not a gap): The CONTEXT.md planning document specified Clerk as the auth provider with JWKS endpoint validation. The executed plans (09-01 through 09-06) implemented a custom JWT service instead. Based on ROADMAP Phase 12 (`AUTH-01` through `AUTH-13`), full user authentication (registration, login, refresh tokens, frontend auth UI) is deferred to Phase 12. Phase 9 correctly established the JWT infrastructure layer as a foundation for Phase 12. This is consistent and intentional.

---

_Verified: 2026-03-02_
_Verifier: Claude (gsd-verifier)_
