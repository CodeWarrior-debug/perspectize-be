# Onboarding: Guest landing + checklist coach (design freeze)

Peer-reviewable freeze of the dual-path first-run experience. Engineering implements against this document; video assets are optional at ship.

**Status:** frozen. Step 1 (backend) + Step 2 (guest landing + shell eligibility) landed; coach UI is Step 3.

**Intro version constant:** `CURRENT_INTRO_VERSION = 1` (bump only on material flow changes that warrant a forced re-intro).

---

## Goals

1. **Guest (“show me”):** short product **mp4 first**, then soft Sign in — no account required to understand the product.
2. **Signed-in (“ready to go”):** skippable **checklist coach** over the **real app** that guides **add first video → leave a perspective**, with per-step skip and skip-all.
3. **Persistence:** thin user fields so finish/skip-all does not re-nag; **Help → Getting started** can replay.
4. **Videos optional at ship:** coach works with copy + CTAs if mp4s are not ready; three asset slots documented for later capture.

---

## Out of scope (v1)

- Glasses picker as required intro step
- Linear `/welcome` wizard or fake sandbox identity
- Per-step skip flags / analytics event schema in DB
- Public read-only full Activity as guest marketing
- Removing UserSelector (separate); coach keys off `me` + Clerk only
- Full help center; Discover-specific checklist steps

---

## User stories

| Actor | Outcome |
|-------|---------|
| Visitor | Watch a short film and understand Perspectize without signing in |
| New signed-in user | Gently coached to add a video and leave a perspective, or skip any/all |
| User who finished/skipped | Not auto-prompted again until Help → Getting started |
| Returning power user (pre-feature) | Backfilled; never sees coach unless version policy re-introduces it |
| User on Discover | Same coach shell, not a second flow |

---

## Functional rules

### Guest (signed out)

- `/` emphasizes value line + **Watch how it works** + **Sign in** (existing Clerk modal).
- Deep-link `/discover` signed out → same guest treatment (no raw empty product as hero).
- Videos: click-to-play, no autoplay sound; hide Watch if URL unset.

### Coach eligibility (signed in)

After `ClerkLoaded` + `me` settled:

```text
showCoach = signedIn && meLoaded &&
  (onboarding.displayNextSession || onboarding.version < CURRENT_INTRO_VERSION)
```

### Step machine (UI-local)

- Steps: (1) Add video → (2) Leave a perspective.
- Per-step skip is **UI-local only** (not persisted).
- Skip-all / complete / dismiss (X) whole intro → `markOnboardingSeen`.
- Esc marks intro seen (control > ceremony).
- Empty library on step 2: optional how-to mp4, back to step 1, or skip.
- Quiet graduate (recommended): on `me` load, if eligible but user already has ≥1 owned content and ≥1 perspective → call `markOnboardingSeen` once without UI.

### Help replay

- Sets `displayNextSession = true` via `setOnboardingDisplayNextSession` (keeps `version` and `completedAt`).
- Opens coach at step 1.

### Videos (config slots)

```ts
CURRENT_INTRO_VERSION = 1
ONBOARDING_VIDEOS = {
  guestProduct?,   // guest landing
  howAddVideo?,    // coach step 1 optional
  howPerspective?, // coach step 2 optional
}
```

Asset paths (optional): `frontend/static/onboarding/*.mp4` or CDN env URLs.

---

## Data model (backend)

Prefer JSONB column `users.onboarding` default `{}` semantics with structured keys:

```json
{
  "version": 0,
  "displayNextSession": true,
  "completedAt": null
}
```

| Field | New user default | Meaning |
|-------|------------------|---------|
| `version` | `0` | Last intro version the user completed/dismissed |
| `displayNextSession` | `true` | Soft flag: show coach next eligible session |
| `completedAt` | `null` | ISO timestamp when last marked seen (nullable) |

### Backfill (migration)

Existing rows at ship:

- `displayNextSession = false`
- `version = CURRENT_INTRO_VERSION` (1)
- `completedAt = now()` (UTC)

Sentinel / all existing users are treated as already seen so they are never ambushed.

### GraphQL

```graphql
type UserOnboarding {
  version: Int!
  displayNextSession: Boolean!
  completedAt: String  # ISO-8601; null if never completed
}

extend type User {
  onboarding: UserOnboarding!
}

markOnboardingSeen(version: Int!): UserOnboarding! @auth
setOnboardingDisplayNextSession(displayNextSession: Boolean!): UserOnboarding! @auth
```

**Auth:** both mutations operate only on the authenticated current user (no client-supplied `userId`).

**`markOnboardingSeen(version)`:**

- `displayNextSession = false`
- `completedAt = now()`
- `version =` argument (typically `CURRENT_INTRO_VERSION`)

**`setOnboardingDisplayNextSession(displayNextSession)`:**

- Updates only that boolean; leaves `version` and `completedAt` unchanged.

**`me.onboarding`:** always non-null; missing/empty JSON maps to new-user defaults (`version: 0`, `displayNextSession: true`, `completedAt: null`).

---

## UX / NFR notes

- Calm: non-blocking sheet/drawer; app usable underneath.
- Mobile: bottom sheet; videos `playsinline`.
- Gate coach on ClerkLoaded + me settled (avoid flash).
- Mutations auth-safe; owner-only.
- Version bump policy: only bump `CURRENT_INTRO_VERSION` on material flow change; document in release notes.

---

## Architecture (summary)

```text
Signed out → GuestLanding (mp4 + Sign in)
Signed in  → ClerkLoaded → me + onboarding → OnboardingCoach (layout shell)
             Help replay → displayNextSession true → coach
             Coach done  → markOnboardingSeen
```

Coach mounts once at app shell so Activity and Discover share one instance.

---

## Testing (this freeze + later)

**Backend (Step 1):** defaults, backfill semantics, owner-only mutations, `me.onboarding` shape.

**Frontend (later steps):** eligibility helper, mark-seen / replay flags, step navigation.

**Manual smoke:** guest mp4 path, new-user coach, skip-all, Help replay, Discover same coach, backfilled user never auto-sees coach.

---

## Delivery map

| Step | Deliverable |
|------|-------------|
| 1 | This doc + migration + GraphQL onboarding on `me` + mutations + Go tests |
| 2 | Guest landing + shell eligibility + config + `OnboardingVideo` |
| 3 | `OnboardingCoach` + real Add Video / Perspective wiring |
| 4 | Help replay + shot lists / VO / ffmpeg notes |

---

## Smoke checklist (full feature)

- [ ] Guest: landing, video control (if URL set), Sign in works
- [ ] New user: coach shows after sign-in
- [ ] Skip all → reload → no coach
- [ ] Complete flows → mark seen → no coach
- [ ] Help → Getting started → coach returns; finish → stays off
- [ ] Discover signed-in: same coach
- [ ] Backfilled user at CURRENT version: never auto-sees coach
- [ ] Quiet graduate: activated library marks seen without UI
