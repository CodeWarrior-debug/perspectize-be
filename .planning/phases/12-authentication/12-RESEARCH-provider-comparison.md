# Phase 12: Authentication Provider Comparison

**Researched:** 2026-02-16
**Domain:** Authentication providers for SvelteKit (static) + Go backend
**Confidence:** HIGH (verified against official docs, pricing pages, and SDKs)

## Summary

This research compares 8 authentication approaches for the Perspectize project: Clerk, Auth0, Supabase Auth, Firebase Auth, Lucia Auth, Hanko, Logto, and Custom JWT. The evaluation focuses on compatibility with the project's specific architecture: a SvelteKit frontend using `adapter-static` with `prerender = true` (NOT SPA mode), a Go backend with chi router and gqlgen GraphQL, and deployment on Sevalla.

**Critical architecture constraint:** The frontend uses `@sveltejs/adapter-static` with `prerender = true` and `fallback: '404.html'`. This is static prerendering, not SPA mode. Any auth provider that requires SSR server hooks (`hooks.server.ts`) will NOT work without changing the deployment model. Providers must support client-side-only auth flows, with the Go backend handling token verification.

**Primary recommendation:** Clerk is the best fit for this project. It has an official Go SDK with JWT middleware, an actively maintained community Svelte 5 SDK (`svelte-clerk`), works client-side without SSR, provides pre-built UI components, and has a generous free tier (50,000 MAU). The main tradeoff is vendor lock-in, which is acceptable at this project's scale.

---

## Provider Comparison Matrix

### Compatibility with Perspectize Architecture

| Provider | SvelteKit Static Support | Go SDK | Pre-built UI | Free Tier MAU | Verdict |
|----------|------------------------|--------|-------------|---------------|---------|
| **Clerk** | YES (client-side JS SDK) | Official | YES (full suite) | 50,000 | RECOMMENDED |
| **Auth0** | Partial (JS SDK, no Svelte SDK) | Official (go-jwt-middleware) | YES (Universal Login) | 25,000 | Viable but complex |
| **Supabase Auth** | NO (requires SSR hooks) | Community only | NO | 50,000 | Poor fit |
| **Firebase Auth** | YES (client-side JS SDK) | Official (Admin SDK) | YES (FirebaseUI) | 50,000 | Viable |
| **Lucia Auth** | DEPRECATED | N/A | N/A | N/A | Do not use |
| **Hanko** | YES (web components) | Community SDK | YES (web components) | 10,000 | Viable |
| **Logto** | NO (requires SSR hooks) | Official (Go validator) | YES (hosted) | 50,000 | Poor fit |
| **Custom JWT** | YES (manual) | Self-built | NO | Unlimited | High effort |

### Pricing Comparison

| Provider | Free Tier | Cost at 100 Users | Cost at 1K Users | Cost at 10K Users | Self-Host? |
|----------|-----------|-------------------|------------------|--------------------|-----------|
| **Clerk** | 50K MAU | $0 | $0 | $0 | No |
| **Auth0** | 25K MAU (limited features) | $0 (free) or $35/mo (prod features) | $0 or $35/mo | $0 or $35/mo | No |
| **Supabase Auth** | 50K MAU | $0 | $0 | $0 | Yes (GoTrue) |
| **Firebase Auth** | 50K MAU | $0 | $0 | $0 | No |
| **Hanko** | 10K MAU | $0 | $0 | $0 | Yes |
| **Logto** | 50K MAU | $0 | $0 | $0 | Yes |
| **Custom JWT** | Unlimited | $0 | $0 | $0 | N/A (self-built) |

**Note:** At <100 users, ALL providers are effectively free. Pricing only matters at scale.

---

## Detailed Provider Analysis

### 1. Clerk (RECOMMENDED)

**Confidence: HIGH** - Verified via official docs, pricing page, Go SDK repo, and svelte-clerk repo.

**SvelteKit Support:**
- Community SDK: `svelte-clerk` (wobsoriano) - 224 stars, actively maintained, last updated Feb 4, 2026
- Supports Svelte 5 runes natively
- Works client-side only (no SSR required)
- Components: `<SignIn>`, `<SignUp>`, `<UserButton>`, `<UserProfile>`, `<OrganizationSwitcher>`
- The older `clerk-sveltekit` (markjaquith) was DEPRECATED Aug 2025 in favor of `svelte-clerk`
- Listed as a community SDK on Clerk's official docs page

**Go SDK:**
- OFFICIAL: `github.com/clerk/clerk-sdk-go/v2` - maintained by Clerk
- JWT verification via `clerk/jwt` package
- HTTP middleware for chi-compatible routers
- JWKS caching built-in
- Last published: Jan 5, 2026

**Static/SPA Compatibility:**
- YES - Clerk's JavaScript SDK (`@clerk/clerk-js`) is designed for client-side use
- SPA mode explicitly supported (documented for Remix SPA mode, same principle applies)
- Session tokens sent as Bearer tokens in Authorization header
- No server-side hooks required

**Pricing (as of Feb 2026):**
- Hobby (Free): 50,000 MAU, unlimited apps, 3 dashboard seats
- Pro: $20/mo (annual) or $25/mo (monthly), includes MFA, passkeys, custom branding
- Overage: $0.02/MAU (50K-100K), decreasing at volume

**Pre-built UI:** Full suite of customizable components (sign-in, sign-up, user profile, org management)

**Social Login:** Google, GitHub, Discord, Apple, Microsoft, and 20+ providers

**Passkeys/WebAuthn:** YES (included in Pro plan, $25/mo)

**Vendor Lock-in Risk:** MEDIUM - Users stored in Clerk's infrastructure. Export possible but migration requires re-authentication. Standard JWT tokens mean backend verification code is portable.

**GraphQL Compatibility:** No special considerations. JWT Bearer token in Authorization header works with any GraphQL server. Go middleware extracts user from JWT claims into context.

**Sources:**
- [Clerk Pricing](https://clerk.com/pricing)
- [Clerk Go SDK](https://github.com/clerk/clerk-sdk-go)
- [svelte-clerk](https://github.com/wobsoriano/svelte-clerk)
- [Clerk Community SDKs](https://clerk.com/docs/references/community-sdk/overview)
- [Clerk JavaScript SDK](https://clerk.com/docs/reference/javascript/overview)
- [Clerk Go JWT Verification](https://clerk.com/docs/guides/sessions/verifying)
- [Clerk New Pricing Feb 2026](https://clerk.com/changelog/2026-02-05-new-plans-more-value)

---

### 2. Auth0

**Confidence: HIGH** - Verified via official docs and Go middleware repo.

**SvelteKit Support:**
- NO official Svelte/SvelteKit SDK
- Can use vanilla `@auth0/auth0-spa-js` (client-side JavaScript SDK)
- Manual integration required - no pre-built Svelte components
- Community examples exist but are fragmented and often outdated

**Go SDK:**
- OFFICIAL: `github.com/auth0/go-jwt-middleware/v3` - maintained by Auth0
- v3 requires Go 1.24+, uses generics for type-safe claims
- RS256 verification against Auth0's public keys
- DPoP support for enhanced security
- Works with chi router via standard `net/http` middleware

**Static/SPA Compatibility:**
- YES via `@auth0/auth0-spa-js` - designed for SPAs
- Universal Login (redirect-based) works with static sites
- No SSR required

**Pricing:**
- Free: 25,000 MAU (B2C), but lacks MFA and RBAC
- Essentials B2C: $35/mo for 500 MAU (production-grade features)
- Professional: $240/mo for 500 MAU
- Free tier described as "free plan illusion" - high MAU limit but missing essential production features

**Pre-built UI:** Universal Login page (hosted by Auth0, customizable)

**Social Login:** 30+ providers

**Passkeys/WebAuthn:** YES (paid plans)

**Vendor Lock-in Risk:** MEDIUM-HIGH - Strong ecosystem lock-in. More complex setup means more Auth0-specific code.

**Verdict:** Overkill for a <100 user project. The Go SDK is excellent but the lack of Svelte SDK means significant manual frontend work. Pricing becomes expensive quickly once you need production features.

**Sources:**
- [Auth0 Pricing](https://auth0.com/pricing)
- [Auth0 Go JWT Middleware](https://github.com/auth0/go-jwt-middleware)
- [Auth0 Go Quickstart](https://auth0.com/docs/quickstart/backend/golang/interactive)

---

### 3. Supabase Auth (GoTrue)

**Confidence: HIGH** - Verified via official docs and GitHub issues.

**SvelteKit Support:**
- Official `@supabase/ssr` package - but REQUIRES server-side hooks (`hooks.server.ts`)
- Cookie-based auth that relies on SvelteKit's server runtime
- GitHub issue #882 confirms: static adapter does NOT support the recommended auth flow
- Client-side-only `@supabase/supabase-js` auth works but loses SSR cookie management

**Go SDK:**
- COMMUNITY only: `github.com/supabase-community/auth-go` - pre-release
- Original `gotrue-go` also community-maintained
- Not production-ready for a standalone Go backend

**Static/SPA Compatibility:**
- POOR - The recommended `@supabase/ssr` flow requires server hooks
- Can use client-side JS SDK directly but this is not the recommended path
- Would need to use Supabase as the auth source and verify JWTs in Go backend manually

**Pricing:**
- Free: 50,000 MAU, 2 projects, 500MB DB
- Projects PAUSED after 7 days of inactivity on free tier
- Pro: $25/mo, 100K MAU, no pausing

**Self-Hosting:** YES - GoTrue is open source, can self-host the auth server

**Vendor Lock-in Risk:** LOW (open source) but HIGH coupling if using Supabase's full platform

**Verdict:** Poor fit. The project already has its own PostgreSQL database and Go backend. Supabase Auth is designed for the Supabase platform, not standalone use. The static adapter incompatibility is a dealbreaker. The Go SDK is community-maintained and pre-release.

**Sources:**
- [Supabase Pricing](https://supabase.com/pricing)
- [Supabase SvelteKit Auth](https://supabase.com/docs/guides/auth/server-side/sveltekit)
- [Static Adapter Issue #882](https://github.com/supabase/supabase-js/issues/882)
- [auth-go](https://github.com/supabase-community/auth-go)

---

### 4. Firebase Auth

**Confidence: HIGH** - Verified via official docs and Go Admin SDK.

**SvelteKit Support:**
- NO official Svelte SDK
- Client-side JS SDK (`firebase/auth`) works well in static/prerendered apps
- Multiple community guides for Firebase + SvelteKit client-side auth
- Svelte 5 runes can manage auth state via `onAuthStateChanged` listener
- SDK is large (~100KB) but supports lazy loading from CDN

**Go SDK:**
- OFFICIAL: `firebase.google.com/go/v4` - Firebase Admin Go SDK
- `VerifyIDToken()` validates client tokens with JWKS caching (24h)
- `VerifyIDTokenAndCheckRevoked()` for additional revocation checks
- Well-documented, production-proven

**Static/SPA Compatibility:**
- YES - Firebase Auth JS SDK is entirely client-side
- No server hooks needed
- `onAuthStateChanged` listener provides reactive auth state
- ID tokens sent to backend as Bearer tokens

**Pricing:**
- Free: 50,000 MAU for email/password and social login
- Phone auth NOT free (per-verification charges)
- SAML/OIDC: 50 MAU free, then paid
- No monthly base fee

**Pre-built UI:** FirebaseUI (drop-in auth widget) - functional but dated appearance

**Social Login:** Google, Apple, Microsoft, GitHub, Facebook, Twitter, Yahoo

**Passkeys/WebAuthn:** Limited support (via Google Identity Platform upgrade)

**Vendor Lock-in Risk:** MEDIUM - Google ecosystem. Users exportable. JWT verification is standard.

**Verdict:** Viable option. Excellent Go SDK, works with static sites, generous free tier. Downsides: no Svelte-specific SDK (manual integration), FirebaseUI looks dated, Google dependency. Less polished developer experience than Clerk.

**Sources:**
- [Firebase Auth Pricing](https://firebase.google.com/pricing)
- [Firebase Auth Limits](https://firebase.google.com/docs/auth/limits)
- [Firebase Admin Go SDK](https://pkg.go.dev/firebase.google.com/go/auth)
- [Client-side Firebase + SvelteKit](https://www.okupter.com/blog/client-side-authentication-firebase-sveltekit)

---

### 5. Lucia Auth

**Confidence: HIGH** - Verified via official announcement and GitHub discussion.

**Status: DEPRECATED as of March 2025.**

The creator deprecated Lucia, stating the library was "not working" in its current form. Database adapters were too complex for the value provided. Lucia is now a learning resource for implementing sessions from scratch, not a usable library.

**Verdict:** Do not use. The npm package is no longer maintained. The suggested replacement is Better Auth, but Better Auth is JavaScript/TypeScript only and does not integrate with a Go backend.

**Sources:**
- [Lucia Deprecation Discussion #1714](https://github.com/lucia-auth/lucia/discussions/1714)
- [Lucia Migration Guide](https://lucia-auth.com/lucia-v3/migrate)

---

### 6. Hanko

**Confidence: MEDIUM** - Verified via official docs and GitHub. Go SDK less thoroughly documented.

**SvelteKit Support:**
- Official web components: `@teamhanko/hanko-elements`
- SvelteKit starter template exists: `hanko-sveltekit-starter`
- Web components work in any framework including static sites
- Framework-agnostic approach (custom elements, not framework-specific components)

**Go SDK:**
- Community: `github.com/teamhanko/hanko-sdk-golang`
- FIDO2/WebAuthn focused
- Less documented than Clerk or Auth0 Go SDKs
- JWT verification possible but requires manual setup

**Static/SPA Compatibility:**
- YES - Web components are client-side by design
- Passkey-first approach works entirely client-side
- Backend verification via JWKS

**Pricing:**
- Starter (Free): 10,000 MAU, 2 projects (paused after 7 days inactive)
- Pro: $29/mo + $0.01/MAU over 10K
- Startup program: 1M MAU free (application required)

**Pre-built UI:** YES - Web components for login, registration, profile

**Passkeys/WebAuthn:** EXCELLENT - This is Hanko's primary focus. FIDO2 certified.

**Self-Hosting:** YES - Fully open source

**Vendor Lock-in Risk:** LOW - Open source, self-hostable, standards-based

**Verdict:** Interesting if passkeys are a priority. Lower free tier (10K vs 50K MAU). Smaller community and less mature Go SDK. The passkey-first approach is forward-looking but may confuse users who expect traditional password login.

**Sources:**
- [Hanko Pricing](https://www.hanko.io/pricing)
- [Hanko GitHub](https://github.com/teamhanko/hanko)
- [Hanko SvelteKit Starter](https://github.com/teamhanko/hanko-sveltekit-starter)
- [Hanko Go SDK](https://pkg.go.dev/github.com/teamhanko/hanko-sdk-golang)

---

### 7. Logto

**Confidence: MEDIUM** - Verified via official docs. SvelteKit constraint verified.

**SvelteKit Support:**
- Official SDK: `@logto/sveltekit`
- REQUIRES SSR - uses `hooks.server.ts` and `handleLogto` server hook
- Session management via encrypted cookies (server-side only)
- DOES NOT work with `adapter-static` or prerendered pages

**Go SDK:**
- Official Go validator: `github.com/logto-io/go`
- JWT verification with JWKS, audience, issuer validation
- RBAC and scope checking documented
- Well-documented API protection guide

**Static/SPA Compatibility:**
- NO - The SvelteKit SDK explicitly requires SSR server hooks
- Would need to use vanilla OIDC flow with a different client library
- Redirect-based auth flow requires server callback handling

**Pricing:**
- Free: 50,000 MAU, 100K tokens
- Pro: $16/mo, unlimited MAU
- MFA add-on: $48/mo
- Organizations add-on: $48/mo

**Self-Hosting:** YES - Open source, Docker deployment

**Vendor Lock-in Risk:** LOW - Open source, OIDC/OAuth2 standards

**Verdict:** Good product but incompatible with the project's static deployment. The SvelteKit SDK requiring SSR is a dealbreaker unless the frontend deployment model changes. The Go SDK is solid but the frontend integration gap makes this impractical.

**Sources:**
- [Logto SvelteKit Quickstart](https://docs.logto.io/quick-starts/sveltekit)
- [Logto Go API Protection](https://docs.logto.io/api-protection/go)
- [Logto Go SDK](https://github.com/logto-io/go)
- [Logto Pricing](https://blog.logto.io/auth0-pricing-explain)

---

### 8. Custom JWT (Current Phase 9/12 Plan)

**Confidence: HIGH** - Based on existing CONTEXT.md research and Go ecosystem knowledge.

**What is Planned:**
- Argon2id password hashing (OWASP 2026 standard)
- JWT access tokens (15 min) + httpOnly refresh cookies (7 days)
- `go-chi/jwtauth` middleware
- Custom registration/login GraphQL mutations
- Custom frontend auth state management with Svelte 5 runes

**SvelteKit Support:** Manual - build all auth UI components from scratch

**Go SDK:** Self-built - full control but full maintenance burden

**Static/SPA Compatibility:** YES - fully custom, works however you build it

**Pricing:** $0 (only infrastructure costs)

**Pre-built UI:** NONE - must build login, signup, password reset, user profile from scratch

**Social Login:** Must implement OAuth2 flows manually per provider

**Passkeys/WebAuthn:** Must implement from scratch (complex - FIDO2 spec is non-trivial)

**Vendor Lock-in Risk:** ZERO

**Estimated Development Effort:**
- Basic auth (login/signup/JWT): 3-5 days
- Password reset flow: 1-2 days
- Email verification: 1-2 days
- Social login (per provider): 1-2 days each
- Pre-built UI components: 3-5 days
- Security hardening (rate limiting, CSRF, etc.): 2-3 days
- **Total: 10-17 days** for a production-ready system

**Risks:**
- Security vulnerabilities from hand-rolling auth
- Missing edge cases (token rotation, concurrent sessions, account lockout)
- Ongoing maintenance burden (security patches, new attack vectors)
- No pre-built UI means more frontend work

**Verdict:** Maximum control but maximum effort. The 10-17 day estimate assumes everything goes right. Auth is a solved problem with well-known security pitfalls. For a small project, the development time exceeds the cost savings vs. a hosted provider.

---

## Architecture Compatibility Summary

### The Static Adapter Constraint

The Perspectize frontend uses:
```javascript
// svelte.config.js
adapter: adapter({ fallback: '404.html', strict: true })

// +layout.ts
export const prerender = true;
```

This means:
1. No `hooks.server.ts` available at runtime
2. No server-side load functions at runtime
3. All auth MUST be client-side JavaScript
4. Token verification happens in the Go backend, not in SvelteKit

**Providers that WORK with static adapter:**
- Clerk (client-side JS SDK + svelte-clerk)
- Auth0 (auth0-spa-js)
- Firebase Auth (firebase/auth)
- Hanko (web components)
- Custom JWT (manual)

**Providers that DO NOT WORK without changing deployment:**
- Supabase Auth (requires @supabase/ssr + hooks.server.ts)
- Logto (requires @logto/sveltekit server hooks)

### Recommended Auth Flow for Static SvelteKit + Go Backend

```
1. User clicks "Sign In" in SvelteKit app
2. Auth provider's client-side SDK handles login UI
3. Provider returns JWT/session token to client
4. Client stores token in memory (access) or httpOnly cookie (refresh)
5. Client sends Bearer token with every GraphQL request
6. Go backend middleware verifies JWT using provider's JWKS
7. Verified user ID injected into GraphQL resolver context
```

This flow works with Clerk, Auth0, Firebase, and Hanko.

---

## Recommendation Ranking

### Tier 1: Recommended

**1. Clerk** - Best overall fit
- Official Go SDK with JWT middleware
- Active Svelte 5 community SDK (svelte-clerk, 224 stars, updated Feb 2026)
- Client-side components work with static adapter
- 50K MAU free tier covers this project indefinitely
- Pre-built UI saves weeks of frontend work
- Passkeys support on Pro plan ($25/mo)
- Tradeoff: Vendor lock-in, community (not official) Svelte SDK

### Tier 2: Viable Alternatives

**2. Firebase Auth** - If you want Google ecosystem
- Official Go Admin SDK (excellent)
- Client-side JS SDK works with static sites
- 50K MAU free
- Downside: No Svelte SDK, manual integration, dated UI widgets, Google dependency

**3. Auth0** - If enterprise features needed later
- Official Go JWT middleware (excellent)
- Client-side SPA SDK works
- Free tier limited (no MFA/RBAC without paid plan)
- Downside: No Svelte SDK, complex setup, expensive at scale

**4. Hanko** - If passkeys are the priority
- Web components work anywhere
- Open source, self-hostable
- FIDO2 certified
- Downside: Smaller community, less mature Go SDK, 10K free tier

### Tier 3: Not Recommended

**5. Custom JWT** - Too much effort for this project size
- 10-17 days of work vs. hours with Clerk
- Security risk from hand-rolling auth
- No pre-built UI

**6. Logto** - Good product, wrong deployment model
- Requires SSR, incompatible with static adapter

**7. Supabase Auth** - Wrong architecture
- Designed for Supabase platform, not standalone Go backends
- Requires SSR, community-only Go SDK

**8. Lucia Auth** - Deprecated
- Do not use

---

## Migration Path: Custom JWT to Clerk

If Phase 12 adopts Clerk instead of custom JWT, the existing CONTEXT.md plan changes as follows:

| CONTEXT.md Plan | With Clerk |
|-----------------|-----------|
| Argon2id password hashing | Not needed (Clerk handles) |
| JWT token generation | Not needed (Clerk issues JWTs) |
| Refresh token storage table | Not needed (Clerk manages sessions) |
| `go-chi/jwtauth` middleware | Replace with `clerk-sdk-go` middleware |
| Custom login/register mutations | Remove (Clerk UI handles) |
| Frontend auth state (runes) | Use svelte-clerk's auth state |
| User dropdown replacement | Replace with `<UserButton>` component |
| Password set flow for existing users | Use Clerk's user import/invitation |

**What stays the same:**
- GraphQL resolver authorization checks (`auth.ForContext(ctx)`)
- User-scoped TanStack Query cache keys
- Bearer token in Authorization header

**Migration of existing users:**
- Clerk supports user import via Backend API
- Can bulk-create users with metadata
- Users would need to set passwords on first Clerk login

---

## Open Questions

1. **svelte-clerk with prerender=true**: The `svelte-clerk` SDK has been verified to support Svelte 5, but explicit testing with `adapter-static` + `prerender=true` (not SPA mode) should be done. The Clerk JS SDK itself is client-side only, so it should work, but the SvelteKit integration layer may have assumptions. **Recommendation:** Test early in implementation.

2. **Clerk Free Tier for Production**: The free tier now includes 50K MAU with no pausing behavior (unlike Supabase/Hanko). Verify there are no feature restrictions that would block production use (e.g., custom domain appears to be included in free tier).

3. **Sevalla Deployment**: Clerk's client-side SDK should work on any static hosting. The Go backend needs outbound access to Clerk's JWKS endpoint for token verification. Verify Sevalla allows this.

4. **Existing User Migration**: The project currently has users in a PostgreSQL `users` table. These would need to be imported into Clerk. The Clerk Backend API supports user creation, but the existing `users` table would need a `clerk_user_id` foreign key instead of local auth.

---

## Sources

### Primary (HIGH confidence)
- [Clerk Pricing Page](https://clerk.com/pricing) - Verified Feb 2026 pricing
- [Clerk Go SDK](https://github.com/clerk/clerk-sdk-go) - Official, v2, last published Jan 2026
- [svelte-clerk](https://github.com/wobsoriano/svelte-clerk) - Community SDK, 224 stars, Feb 2026
- [Clerk Community SDKs](https://clerk.com/docs/references/community-sdk/overview) - Official listing
- [Clerk JS SDK](https://clerk.com/docs/reference/javascript/overview) - Client-side architecture
- [Auth0 Go JWT Middleware](https://github.com/auth0/go-jwt-middleware) - Official, v3
- [Firebase Admin Go SDK](https://pkg.go.dev/firebase.google.com/go/auth) - Official
- [Logto SvelteKit Docs](https://docs.logto.io/quick-starts/sveltekit) - SSR requirement verified

### Secondary (MEDIUM confidence)
- [Hanko Pricing](https://www.hanko.io/pricing) - Verified via WebFetch
- [Hanko Go SDK](https://pkg.go.dev/github.com/teamhanko/hanko-sdk-golang) - Community
- [Supabase Static Adapter Issue](https://github.com/supabase/supabase-js/issues/882) - Confirms incompatibility

### Tertiary (LOW confidence)
- [Auth Provider Comparison 2026](https://designrevision.com/blog/auth-providers-compared) - Third-party blog
- [Clerk vs Supabase](https://www.devtoolsacademy.com/blog/supabase-vs-clerk/) - Third-party comparison

## Metadata

**Confidence breakdown:**
- Provider pricing: HIGH - Verified directly from pricing pages
- SvelteKit compatibility: HIGH - Verified via SDK docs and architecture analysis
- Go SDK availability: HIGH - Verified via pkg.go.dev and GitHub repos
- Static adapter constraint: HIGH - Verified from project's own svelte.config.js
- svelte-clerk + static adapter: MEDIUM - SDK supports Svelte 5 but not explicitly tested with prerender=true

**Research date:** 2026-02-16
**Valid until:** 2026-04-16 (pricing may change; SDK versions will update)
