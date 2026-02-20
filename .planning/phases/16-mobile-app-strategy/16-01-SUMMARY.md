# Phase 16-01 Summary: Capacitor + PWA Setup

**Completed:** 2026-02-20
**Branch:** `claude/plan-mobile-app-capacitor-hO2s7`

## What Was Done

### Capacitor 7.5.0 Integration
- Installed `@capacitor/core@7.5.0`, `@capacitor/haptics@7.0.3` (runtime deps)
- Installed `@capacitor/cli@7.5.0`, `@capacitor/ios@7.5.0`, `@capacitor/android@7.5.0` (dev deps)
- Initialized Capacitor with `appId: com.perspectize.app`, `webDir: build`
- Added `androidScheme: 'https'` for secure cookies/storage on Android
- Generated iOS platform (`ios/`) — Xcode project created; `pod install` requires macOS
- Generated Android platform (`android/`) — fully synced with haptics plugin

### PWA via @vite-pwa/sveltekit 1.1.0
- Installed `@vite-pwa/sveltekit@1.1.0` with Workbox service worker
- Configured manifest: name "Perspectize", navy theme (#1a365d), standalone display
- Created placeholder icons (192px, 512px, 512px-maskable) in `static/icons/`
- Disabled SvelteKit built-in service worker (`serviceWorker.register: false`)
- Added PWA manifest link injection in `+layout.svelte` using Svelte 5 `$derived`
- Added type declarations for `virtual:pwa-info` module

### Platform Detection Utility
- Created `src/lib/utils/native.ts` with three exports:
  - `isNativePlatform()` — branch native vs web behavior
  - `getPlatform()` — returns 'ios', 'android', or 'web'
  - `nativeFeedback()` — haptic feedback for App Store compliance

### Build Infrastructure
- Added npm scripts: `cap:sync`, `cap:open:ios`, `cap:open:android`, `mobile:build`
- Updated `.gitignore` to exclude Capacitor build cache (public/ dirs) while keeping ios/ and android/ in git

## Verification Results

| Check | Result |
|-------|--------|
| `pnpm run build` | PASS — 25 precache entries, 1506 KiB PWA bundle |
| `pnpm run test:run` | PASS — 284 tests, 20 files, 0 failures |
| `capacitor.config.ts` exists | PASS — webDir: 'build' |
| iOS platform created | PASS — ios/App/ exists (pod install needs macOS) |
| Android platform created | PASS — fully synced with haptics plugin |
| PWA service worker generated | PASS — sw.js + workbox generated |
| Platform detection utility | PASS — native.ts with 3 exports |

## Must-Haves Verification

- [x] Capacitor is installed and configured with webDir pointing to frontend/build
- [x] iOS and Android platform projects are generated and checked into git
- [x] PWA manifest is configured with Perspectize branding (navy theme, icons)
- [x] SvelteKit built-in service worker is disabled to avoid conflict with vite-pwa
- [x] Platform detection utility exists for branching native vs web behavior
- [x] Frontend builds successfully with both Capacitor and PWA configurations
- [x] Haptic feedback utility exists for App Store compliance

## Files Modified

| File | Change |
|------|--------|
| `frontend/package.json` | Added Capacitor + PWA deps, mobile scripts |
| `frontend/capacitor.config.ts` | NEW — Capacitor configuration |
| `frontend/vite.config.ts` | Added SvelteKitPWA plugin with manifest |
| `frontend/svelte.config.js` | Disabled built-in service worker |
| `frontend/src/routes/+layout.svelte` | Added PWA manifest link |
| `frontend/src/lib/utils/native.ts` | NEW — platform detection + haptics |
| `frontend/src/pwa.d.ts` | NEW — type declarations for virtual:pwa-info |
| `frontend/static/icons/` | NEW — placeholder PWA icons (192, 512, 512-maskable) |
| `frontend/.gitignore` | Added Capacitor build cache exclusions |
| `frontend/ios/` | NEW — generated Xcode project |
| `frontend/android/` | NEW — generated Android Studio project |

## Notes

- Used Capacitor **7.x** (not 8.x) — Cap 8 requires Xcode 26+ which may not be available yet. Cap 7 requires Xcode 16+ (widely available).
- iOS `pod install` fails in non-macOS environments — expected. The Xcode project structure is created; run `pod install` and `npx cap sync` on a macOS dev machine.
- PWA precache covers 25 entries at 1506 KiB — well within acceptable range.
- No changes to `prerender = true` or `fallback: '404.html'` — current config is compatible with Capacitor.

## Next Steps

- **Plan 16-02:** POC verification in iOS Simulator (requires macOS + Xcode)
- Test CORS with native WebView origins (`capacitor://localhost`)
- Replace placeholder icons with designed assets before store submission
- Integrate with Phase 12 (Authentication) when complete
