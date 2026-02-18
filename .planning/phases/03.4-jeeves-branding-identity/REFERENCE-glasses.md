# Perspectize Glasses — SVGs + Implementation Plan

## SVG Frames

### C1 — Sunglasses (light)
```svg
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 70 28">
  <defs><linearGradient id="c1g" x1="0" y1="0" x2="0.2" y2="1"><stop offset="0%" stop-color="#9F7AEA" stop-opacity="0.5"/><stop offset="100%" stop-color="#44337A" stop-opacity="0.9"/></linearGradient></defs>
  <path d="M10 8 L30 7 Q34 7 34 10.5 L33 17 Q32 21 28 21 L14 21 Q10 21 9 17 L8 10.5 Q8 8 10 8 Z" fill="url(#c1g)" stroke="#44337A" stroke-width="1.4" stroke-linejoin="round"/>
  <path d="M38 8 L58 7 Q62 7 62 10.5 L61 17 Q60 21 56 21 L42 21 Q38 21 37 17 L36 10.5 Q36 8 38 8 Z" fill="url(#c1g)" stroke="#44337A" stroke-width="1.4" stroke-linejoin="round"/>
  <path d="M34 10 Q36 8 36 10" fill="none" stroke="#44337A" stroke-width="1.4" stroke-linecap="round"/>
  <line x1="8" y1="9.5" x2="3" y2="7" stroke="#44337A" stroke-width="1.4" stroke-linecap="round"/>
  <line x1="62" y1="9.5" x2="67" y2="7" stroke="#44337A" stroke-width="1.4" stroke-linecap="round"/>
</svg>
```

### C1 — Sunglasses (dark)
```svg
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 70 28">
  <defs><linearGradient id="c1d" x1="0" y1="0" x2="0.2" y2="1"><stop offset="0%" stop-color="#B794F4" stop-opacity="0.5"/><stop offset="100%" stop-color="#9F7AEA" stop-opacity="0.85"/></linearGradient></defs>
  <path d="M10 8 L30 7 Q34 7 34 10.5 L33 17 Q32 21 28 21 L14 21 Q10 21 9 17 L8 10.5 Q8 8 10 8 Z" fill="url(#c1d)" stroke="#9F7AEA" stroke-width="1.4" stroke-linejoin="round"/>
  <path d="M38 8 L58 7 Q62 7 62 10.5 L61 17 Q60 21 56 21 L42 21 Q38 21 37 17 L36 10.5 Q36 8 38 8 Z" fill="url(#c1d)" stroke="#9F7AEA" stroke-width="1.4" stroke-linejoin="round"/>
  <path d="M34 10 Q36 8 36 10" fill="none" stroke="#9F7AEA" stroke-width="1.4" stroke-linecap="round"/>
  <line x1="8" y1="9.5" x2="3" y2="7" stroke="#9F7AEA" stroke-width="1.4" stroke-linecap="round"/>
  <line x1="62" y1="9.5" x2="67" y2="7" stroke="#9F7AEA" stroke-width="1.4" stroke-linecap="round"/>
</svg>
```

### Teddy A — Round Pince-nez (light)
```svg
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 56 28">
  <defs><linearGradient id="tag" x1="0" y1="0" x2="0.2" y2="1"><stop offset="0%" stop-color="#9F7AEA" stop-opacity="0.45"/><stop offset="100%" stop-color="#44337A" stop-opacity="0.85"/></linearGradient></defs>
  <circle cx="14" cy="14" r="11" fill="url(#tag)" stroke="#44337A" stroke-width="1.5"/>
  <circle cx="42" cy="14" r="11" fill="url(#tag)" stroke="#44337A" stroke-width="1.5"/>
  <path d="M25 11 Q28 6 31 11" fill="none" stroke="#44337A" stroke-width="1.5" stroke-linecap="round"/>
</svg>
```

### Teddy A — Round Pince-nez (dark)
```svg
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 56 28">
  <defs><linearGradient id="tad" x1="0" y1="0" x2="0.2" y2="1"><stop offset="0%" stop-color="#B794F4" stop-opacity="0.45"/><stop offset="100%" stop-color="#9F7AEA" stop-opacity="0.8"/></linearGradient></defs>
  <circle cx="14" cy="14" r="11" fill="url(#tad)" stroke="#9F7AEA" stroke-width="1.5"/>
  <circle cx="42" cy="14" r="11" fill="url(#tad)" stroke="#9F7AEA" stroke-width="1.5"/>
  <path d="M25 11 Q28 6 31 11" fill="none" stroke="#9F7AEA" stroke-width="1.5" stroke-linecap="round"/>
</svg>
```

### Teddy B — Oval Pince-nez (light)
```svg
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 62 28">
  <defs><linearGradient id="tbg" x1="0" y1="0" x2="0.15" y2="1"><stop offset="0%" stop-color="#9F7AEA" stop-opacity="0.4"/><stop offset="100%" stop-color="#44337A" stop-opacity="0.85"/></linearGradient></defs>
  <ellipse cx="16" cy="14" rx="13" ry="10" fill="url(#tbg)" stroke="#44337A" stroke-width="1.5"/>
  <ellipse cx="48" cy="14" rx="13" ry="10" fill="url(#tbg)" stroke="#44337A" stroke-width="1.5"/>
  <path d="M29 11 Q32 5 35 11" fill="none" stroke="#44337A" stroke-width="1.5" stroke-linecap="round"/>
</svg>
```

### Teddy B — Oval Pince-nez (dark)
```svg
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 62 28">
  <defs><linearGradient id="tbd" x1="0" y1="0" x2="0.15" y2="1"><stop offset="0%" stop-color="#B794F4" stop-opacity="0.4"/><stop offset="100%" stop-color="#9F7AEA" stop-opacity="0.8"/></linearGradient></defs>
  <ellipse cx="16" cy="14" rx="13" ry="10" fill="url(#tbd)" stroke="#9F7AEA" stroke-width="1.5"/>
  <ellipse cx="48" cy="14" rx="13" ry="10" fill="url(#tbd)" stroke="#9F7AEA" stroke-width="1.5"/>
  <path d="M29 11 Q32 5 35 11" fill="none" stroke="#9F7AEA" stroke-width="1.5" stroke-linecap="round"/>
</svg>
```

---

## Implementation Plan

### 1. Svelte Component Architecture
- **1a)** Single `<GlassesIcon>` component with `{#if shape === 'sunglasses'}` blocks — simplest for 3 shapes
- **1b)** Separate SVG components per shape, parent delegates — cleaner separation, more files
- **1c)** Raw SVG strings in config map via `{@html}` — most flexible, needs sanitization
- **Recommended: 1a** (refactor to 1b later if shape count grows)

### 2. Color Selection
- **2a)** Preset palette picker — 6-8 curated colors, each defines stroke + gradient stops. Guarantees quality.
- **2b)** HSL hue slider — derive gradient/stroke from single hue. More freedom, still constrained.
- **2c)** Full color picker — max freedom, users will pick bad combos.
- **Recommended: 2a** now, add 2b as "custom" option later. Curated choices = "thoughtful friction."

### 3. Where Glasses Appear
- **3a)** User avatar/profile icon only — shown next to perspectives + profile
- **3b)** Avatar + watermark on perspective cards
- **3c)** Interactive onboarding — "pick your glasses" as first-run setup
- **Recommended: 3a + 3c** — onboarding creates investment, avatar gives ongoing visibility

### 4. Storage
- **4a)** Separate columns: `glasses_shape VARCHAR`, `glasses_color VARCHAR`
- **4b)** JSONB column: `preferences JSONB` → `{"glasses": {"shape": "oval", "color": "purple"}}`
- **Recommended: 4b** — extensible for future prefs (theme, typography, etc.)

### 5. Implementation Order
- Phase 1: `<GlassesIcon>` Svelte component with shape/color props
- Phase 2: Picker UI (shape toggle + color swatches), wire to local state
- Phase 3: JSONB preferences column, persist via Go API
- Phase 4: Show glasses in user avatars across app
