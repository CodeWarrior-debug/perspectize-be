# Perspectize v1.0 — Design Specification

> Use this document to structure your Figma file and prompt Figma Make.
> Each section maps to a Figma page or component set you should design.

---

## Logo

### Concept

Wireframe glasses in purple. Conveys the sophistication and clarity of a well-read, thoughtful person — someone with a bookshelf's worth of opinions distilled into an app. The wireframe style bridges wisdom with modern minimalism.

### Color

| Token | Hex | Usage |
|-------|-----|-------|
| `logo-purple` | TBD | Logo primary — needs to read well on both navy header (#1a365d) and white backgrounds |

**Considerations:** The purple must have enough contrast against navy. A lighter/brighter purple (like `#a78bfa` violet-400 or `#8b5cf6` violet-500) will pop on navy. A deeper purple (like `#7c3aed` violet-600) risks blending into the dark header.

### Open Decisions

| Decision | Option A | Option B | Notes |
|----------|----------|----------|-------|
| **Eyes** | Glasses only (empty lenses) | Subtle eyes/dots visible through lenses | Empty = more iconic/minimal, cleaner at small sizes. Eyes = more personality, hints at "someone looking." |
| **Wordmark** | Icon only | Icon + "Perspectize" text | Icon-only is more flexible (works at favicon size). Text version for header where space allows. Consider designing both — icon-only for small contexts, icon+text for header. |
| **Glasses style** | Round wireframe (intellectual) | Rectangular wireframe (modern) | Round = more "professor/wisdom" feel. Rectangular = cleaner/tech. Round recommended for the concept. |

### SVG Requirements

- Must work at 3 sizes: favicon (16x16), header icon (24-32px), and splash/about (64px+)
- Single-color (purple) with transparent fill for lenses
- Stroke-based (wireframe) rather than filled shapes
- Clean paths, no gradients — works in light and dark contexts

### Design Prompt (for generating or briefing)

> "Minimal wireframe glasses icon, single stroke weight, round or slightly rounded lenses, thin bridge and temples. Purple color on transparent background. The style should feel like a refined line drawing — think New Yorker magazine illustration meets modern app icon. No fill on the lenses. SVG format, clean paths, designed to work from 16px to 64px."

---

## Current State Audit

### What's Broken

| # | Issue | Severity | Root Cause |
|---|-------|----------|------------|
| 1 | **Dialog is transparent** — table data bleeds through the Add Video form | Critical | `bg-background` references `--color-background` which is never defined in `app.css`. Resolves to `transparent`. |
| 2 | **No color token system** — only 2 of ~15 required semantic colors are defined | Critical | `app.css` `@theme` only has `--color-primary` and `--color-primary-foreground`. Missing: `background`, `foreground`, `border`, `muted`, `card`, `accent`, `secondary`, `destructive`, `ring`, `input`, `popover`. |
| 3 | **AG Grid is unstyled** — bare themeQuartz with no brand colors | Major | `themeQuartz.withParams()` only sets font and font-size. No accent color, no row styling, no header color. |
| 4 | **No visual hierarchy** — flat white page with floating elements | Major | No card containers, no section backgrounds, no shadows, no depth. |
| 5 | **Header is sparse** — plain white bar, native `<select>`, cramped on mobile | Moderate | No visual weight, no brand color presence, native form control instead of styled component. |
| 6 | **Font mismatch** — Figma uses Inter, code uses Geist | Minor | Decision: **Dual-font system** — Geist (sans) for headings/labels/buttons/nav, Charter (serif) for body text/reviews/descriptions. Figma updated to Geist. AG Grid `fontFamily` needs update. |

### What's Working

- Layout gutters ARE present (`px-4` / `md:px-6` / `lg:px-8` via PageWrapper) — but the AG Grid's internal content extends to the gutter edge, making it feel edge-to-edge
- Mobile responsive basics work (header stacks OK, grid hides columns)
- The PageWrapper component exists and wraps content with `max-w-screen-xl mx-auto`
- Toast notifications work
- Data fetching, mutations, and error handling all function correctly

---

## Design Tokens (Color Palette)

Design these as a **Figma local styles / variables** set. These are the semantic tokens the code needs.

### Light Theme (from Figma Radix 3.0 file)

> **Previous values (before Figma alignment):** The spec originally used a slate palette instead of neutral/zinc. Key differences: `foreground` was `#0f172a` (slate-900), `secondary-foreground` was `#0f172a`, `muted-foreground` was `#64748b` (slate-500), `secondary`/`muted` were `#f1f5f9` (slate-100), `accent` was `#f1f5f9`, `accent-foreground` was `#0f172a`, `border`/`input` were `#e2e8f0` (slate-200), `primary-foreground` was `#f8fafc` (slate-50), `destructive-foreground` was `#fef2f2`. No hover or disabled tokens were defined.

| Token | Hex | Usage |
|-------|-----|-------|
| `background` | `#ffffff` | Page background |
| `foreground` | `#171717` | Default text (neutral-900) |
| `card` | `#ffffff` | Card/surface background |
| `card-foreground` | `#171717` | Text on cards |
| `popover` | `#ffffff` | Dropdown/dialog background |
| `popover-foreground` | `#171717` | Text in popovers |
| `primary` | `#1a365d` | Buttons, links, header accent (navy) |
| `primary-hover` | `#2d3748` | Primary hover state |
| `primary-foreground` | `#ffffff` | Text on primary |
| `secondary` | `#f5f5f5` | Secondary buttons, tags (neutral-100) |
| `secondary-hover` | `rgba(245,245,245,0.8)` | Secondary hover state |
| `secondary-foreground` | `#525252` | Text on secondary (neutral-600) |
| `muted` | `#f5f5f5` | Subtle backgrounds, disabled states (neutral-100) |
| `muted-foreground` | `#525252` | Placeholder text, captions (neutral-600) |
| `accent` | `#f7fafc` | Hover highlights |
| `accent-foreground` | `#2d3748` | Text on accent |
| `destructive` | `#dc2626` | Error states, delete actions (red-600) |
| `destructive-hover` | `#b91c1c` | Destructive hover state (red-700) |
| `destructive-foreground` | `#ffffff` | Text on destructive |
| `border` | `#d4d4d4` | Default borders (neutral-300) |
| `input` | `#d4d4d4` | Input borders (neutral-300) |
| `ring` | `#1a365d` | Focus ring (matches primary) |
| `disabled` | `0.5` | Disabled state opacity |

### Accent Palette (for data visualization, badges)

| Token | Hex | Usage |
|-------|-----|-------|
| `hover` | `#2d3748` | Primary hover state (from Figma spec) |
| `rating-positive` | `#16a34a` | High scores (700-1000), Agree (green-600) |
| `rating-neutral` | `#ca8a04` | Mid scores (400-699) (yellow-600) |
| `rating-negative` | `#dc2626` | Low scores (0-399), Disagree (red-600) |
| `rating-undecided` | `#64748b` | Undecided / no opinion (slate-500) |

**Note:** Null or unsupplied values show no color — render the cell with default text color, no badge or tint.

---

## Typography

### Dual-Font System

| Role | Font | Weights | Usage |
|------|------|---------|-------|
| **Headings & UI** | Geist (sans-serif) | Regular 400, Medium 500, SemiBold 600, Bold 700 | Page titles, section headings, button labels, input labels, navigation, table headers, badges, toasts |
| **Body & Content** | Charter (serif) | Regular 400, Italic 400, Bold 700 | Review text, Like text, video descriptions, longer paragraphs, any user-written content |

**Rationale:** Geist keeps the UI crisp and modern. Charter adds readability for the content-heavy parts of the app (perspective reviews, descriptions). The contrast between sans headings and serif body creates natural visual hierarchy.

### Where Each Font Applies

| Context | Font | Example |
|---------|------|---------|
| Page titles ("Activity") | Geist SemiBold | H1, H2, H3 |
| Button text | Geist Medium | "Add Video", "Submit Perspective" |
| Input labels | Geist Medium | "YouTube URL", "Quality" |
| Table headers | Geist SemiBold | Column names in AG Grid |
| Table cell text | Charter Regular | Video titles, descriptions in grid rows |
| Review/Like text | Charter Regular | User-written perspective content |
| Video descriptions | Charter Regular | Metadata from YouTube |
| Toast messages | Geist Medium | "Added: Video Title" |
| Navigation/header | Geist | Logo text, user selector |
| Captions/metadata | Geist Regular | Dates, durations, "Page 1 of 4" |

### CSS Implementation

```css
@theme {
  --default-font-family: 'Geist', system-ui, -apple-system, sans-serif;
  --font-family-serif: 'Charter', 'Georgia', 'Times New Roman', serif;
}
```

Apply Charter via utility class: `font-serif` (maps to `--font-family-serif`).

### Scale

| Name | Size | Line Height | Weight | Font |
|------|------|-------------|--------|------|
| Display | 36px | 40px | Bold | Geist |
| H1 | 30px | 36px | SemiBold | Geist |
| H2 | 24px | 32px | SemiBold | Geist |
| H3 | 20px | 28px | SemiBold | Geist |
| Body Large | 16px | 24px | Regular | Charter |
| Body | 14px | 20px | Regular | Charter |
| Label | 14px | 20px | Medium | Geist |
| Small/Caption | 12px | 16px | Regular | Geist |

---

## Layout System

### Figma Page: "00 — Layout System"

Design these **structural frames** that every page uses.

#### App Shell

```
┌─────────────────────────────────────────────────────┐
│ Header (sticky, h-16, border-bottom)                │
│  ┌──────┐                    ┌─────────┐ ┌────────┐ │
│  │ Logo │                    │ User ▾  │ │+ Video │ │
│  └──────┘                    └─────────┘ └────────┘ │
├─────────────────────────────────────────────────────┤
│                                                     │
│  ┌─────────────────────────────────────────────┐    │
│  │ PageWrapper (max-w-1280, centered)          │    │
│  │                                             │    │
│  │  Content goes here                          │    │
│  │                                             │    │
│  └─────────────────────────────────────────────┘    │
│                                                     │
└─────────────────────────────────────────────────────┘
```

**Header spec:**
- Height: 64px
- Background: `primary` (#1a365d) — this is the biggest visual win. Navy header gives brand presence immediately.
- Text/icons: `primary-foreground` (white)
- Position: sticky top
- Content: Logo (left), UserSelector + action buttons (right)
- Gutter padding: 16px mobile / 24px tablet / 32px desktop

**PageWrapper spec:**
- Max width: 1280px, centered
- Gutter padding: 16px mobile / 24px tablet / 32px desktop
- Vertical padding: 24px mobile / 32px tablet

#### Breakpoints (design these 3 widths)

| Name | Width | Usage |
|------|-------|-------|
| Mobile | 375px | iPhone SE — single column, stacked layout |
| Tablet | 768px | iPad — wider gutters, columns start appearing |
| Desktop | 1440px | Full layout with max-width container |

---

## Pages to Design

### Figma Page: "01 — Activity (Discover)"

**This is the main landing page. Design 3 breakpoints (375 / 768 / 1440).**

#### Desktop (1440px)

```
┌─ Navy Header ──────────────────────────────────────────────┐
│ Perspectize                    [User ▾]  [+ Add Video]     │
└────────────────────────────────────────────────────────────┘

  ┌─ Page Content (max-w-1280, centered) ──────────────────┐
  │                                                        │
  │  Activity                                              │
  │  Recently updated content            [🔍 Search...]    │
  │                                                        │
  │  ┌─ Card Container (rounded-lg, border, shadow-sm) ──┐ │
  │  │                                                    │ │
  │  │  ┌──────────────────────────────────────────────┐  │ │
  │  │  │ Title          │ Type │ Dur │ Added │ Updated│  │ │
  │  │  ├──────────────────────────────────────────────┤  │ │
  │  │  │ Video title... │ YT   │ 4:14│ Feb 7 │ Feb 7 │  │ │
  │  │  │ (alternating row tint)                       │  │ │
  │  │  │ Video title... │ YT   │ 9:27│ Feb 7 │ Feb 7 │  │ │
  │  │  ├──────────────────────────────────────────────┤  │ │
  │  │  │              Page 1 of 4  < >   10 ▾         │  │ │
  │  │  └──────────────────────────────────────────────┘  │ │
  │  │                                                    │ │
  │  └────────────────────────────────────────────────────┘ │
  │                                                        │
  └────────────────────────────────────────────────────────┘
```

**Key design decisions for this page:**

1. **Search bar**: Move from standalone input to inline with the page title (title left, search right). On mobile it stacks below the title.
2. **Card container**: Wrap the AG Grid in a card with `border`, `rounded-lg`, `shadow-sm`. This gives it visual definition instead of floating.
3. **Table header row**: Navy background (#1a365d) with white text — ties to the app header.
4. **Alternating rows**: Very subtle tint on odd rows — `rgba(26, 54, 93, 0.03)` (3% navy).
5. **Row hover**: `rgba(26, 54, 93, 0.06)` (6% navy).
6. **Pagination**: Inside the card, at the bottom.

#### Mobile (375px)

```
┌─ Navy Header ─────────────────┐
│ Perspectize    [👤] [+ Add]   │
└───────────────────────────────┘

  Activity
  Recently updated content

  [🔍 Search content...]

  ┌─ Card ──────────────────────┐
  │ Title              │ Type   │
  ├────────────────────┼────────┤
  │ Video title...     │ YT     │
  │ Video title...     │ YT     │
  ├────────────────────┴────────┤
  │   Page 1 of 4    < >       │
  └─────────────────────────────┘
```

**Mobile notes:**
- Header: Logo + icon-only buttons (user icon, plus icon) — text labels hidden
- Grid: Only Title + Type columns visible (Duration, dates hidden)
- Search: Full-width below title
- Card: 0px horizontal margin (edge-to-edge within the 16px page gutters)

---

### Figma Page: "02 — Add Video Dialog"

**Design as an overlay component. 2 breakpoints (375 / 768+).**

#### Desktop/Tablet

```
┌──────────── dim overlay (black/50) ───────────────────┐
│                                                       │
│         ┌─ Dialog Card (max-w-md, centered) ─┐        │
│         │                                    │        │
│         │  Add Video                     [✕] │        │
│         │  Paste a YouTube URL to add it     │        │
│         │  to your library.                  │        │
│         │                                    │        │
│         │  YouTube URL                       │        │
│         │  ┌──────────────────────────────┐  │        │
│         │  │ https://youtube.com/watch... │  │        │
│         │  └──────────────────────────────┘  │        │
│         │                                    │        │
│         │         [Cancel]  [Add Video]      │        │
│         │                                    │        │
│         └────────────────────────────────────┘        │
│                                                       │
└───────────────────────────────────────────────────────┘
```

**Critical design requirements:**
- Dialog background: **solid white** (`card` token, `#ffffff`) — NO transparency
- Overlay: `black/50` (50% opacity black backdrop)
- Border radius: `rounded-xl` (12px) for a softer feel
- Shadow: `shadow-lg` for depth/lift off the overlay
- Padding: 24px internal
- Primary button: `primary` navy background
- Cancel button: `outline` variant (border, white bg)

#### Mobile (375px)

```
┌─ dim overlay ─────────────────┐
│                               │
│ ┌─ Dialog (full-width - 32px)┐│
│ │                            ││
│ │  Add Video             [✕] ││
│ │  Paste a YouTube URL...    ││
│ │                            ││
│ │  YouTube URL               ││
│ │  ┌────────────────────┐    ││
│ │  │ https://youtube... │    ││
│ │  └────────────────────┘    ││
│ │                            ││
│ │  [  Add Video (full-w)  ]  ││
│ │  [  Cancel (full-w)     ]  ││
│ │                            ││
│ └────────────────────────────┘│
└───────────────────────────────┘
```

**Mobile dialog notes:**
- Width: `calc(100% - 32px)` (16px margin on each side)
- Buttons stack vertically (full-width each)
- Primary action on top

#### States to design:
1. **Default** — empty input, Add Video button disabled
2. **Filled** — URL entered, Add Video button enabled
3. **Loading** — input disabled, button shows "Adding..." with spinner
4. **Validation error** — red border on input, error text below
5. **Success** — (dialog closes, show toast separately)
6. **API error** — (dialog stays open, show toast separately)

---

### Figma Page: "03 — Add Perspective Dialog"

**Phase 4 — design this now so it's ready when coding starts.**

This is a more complex form. Could be a dialog or a slide-out panel.

#### Desktop

```
┌──────────── dim overlay ──────────────────────────────┐
│                                                       │
│    ┌─ Dialog Card (max-w-lg, centered) ──────────┐    │
│    │                                             │    │
│    │  Add Perspective                        [✕] │    │
│    │  Share your view on this video              │    │
│    │                                             │    │
│    │  Video                                      │    │
│    │  ┌──────────────────────────────────────┐   │    │
│    │  │ 🔍 Search or select a video...    ▾ │   │    │
│    │  └──────────────────────────────────────┘   │    │
│    │                                             │    │
│    │  ── Quick Take ──                           │    │
│    │  Agreement                                  │    │
│    │  ( Agree ) ( Undecided ) ( Disagree )       │    │
│    │                                             │    │
│    │  Comment (optional)                         │    │
│    │  ┌──────────────────────────────────────┐   │    │
│    │  │                                      │   │    │
│    │  └──────────────────────────────────────┘   │    │
│    │                                             │    │
│    │  ── Detailed Ratings (optional) ──          │    │
│    │                                             │    │
│    │  Quality          ┌─────────────┐           │    │
│    │  ━━━━━━━━━━○──── │ 720 / 1000  │           │    │
│    │                   └─────────────┘           │    │
│    │  Importance       ┌─────────────┐           │    │
│    │  ━━━━━○────────── │ 500 / 1000  │           │    │
│    │                   └─────────────┘           │    │
│    │  Confidence       ┌─────────────┐           │    │
│    │  ━━━━━━━━━━━━━━━○ │ 900 / 1000  │           │    │
│    │                   └─────────────┘           │    │
│    │                                             │    │
│    │  Review (optional)                          │    │
│    │  ┌──────────────────────────────────────┐   │    │
│    │  │ Write a longer review...             │   │    │
│    │  │                                      │   │    │
│    │  └──────────────────────────────────────┘   │    │
│    │                                             │    │
│    │           [Cancel]  [Submit Perspective]     │    │
│    │                                             │    │
│    └─────────────────────────────────────────────┘    │
│                                                       │
└───────────────────────────────────────────────────────┘
```

**Design notes:**
- Agreement: Toggle button group (Agree/Undecided/Disagree) with semantic colors
- Ratings: Sliders or number inputs with progress bar visualization
- "Detailed Ratings" section: collapsible or always-visible — design both options
- Video selector: Combobox/searchable dropdown with video titles
- Mobile: Full-height sheet from bottom (not centered dialog) since form is long

---

### Figma Page: "04 — Components"

**Design these as a component library page. This is what Figma Make needs to produce clean code.**

#### Buttons
- `primary` — navy bg, white text (default, hover, disabled, loading)
- `secondary` — slate-100 bg, dark text
- `outline` — white bg, border, dark text
- `ghost` — transparent bg, dark text, hover tint
- `destructive` — red bg, white text
- `link` — underline, navy text
- Sizes: `sm` (h-8), `default` (h-9), `lg` (h-10), `icon` (square)

#### Inputs
- Text input (default, focused, error, disabled)
- Search input (with search icon)
- Textarea (2-line, 4-line variants)
- Select/Dropdown (styled, not native `<select>`)

#### Cards
- Basic card (border + rounded-lg + shadow-sm + white bg)
- Card with header section
- Data card (for future dashboard stats)

#### Dialogs
- Standard dialog (title + description + content + footer)
- Mobile sheet variant (slides up from bottom)

#### Data Display
- Score badge: pill shape, color-coded (green/yellow/red) for ratings
- Video row: title as link, type badge, duration, dates
- Empty state: icon + message + action button
- Loading state: skeleton or spinner

#### Toast Notifications
- Success (green left border or icon)
- Error (red)
- Info (blue/navy)

#### Navigation
- Header (navy variant)
- User selector (styled dropdown, not native `<select>`)

---

### Figma Page: "05 — Toast & Feedback States"

Design the toast notification system and feedback states:

1. **Success toast**: "Added: Video Title" with checkmark icon
2. **Error toast**: "This video has already been added" with X icon
3. **Loading indicator**: Within buttons ("Adding...") and as page skeleton
4. **Empty state**: "No content found" with illustration + "Add your first video" CTA

---

## Figma Make Prompting Guide

**Do NOT paste this entire spec at once.** Figma Make works best with focused, building-block prompts. Follow this exact sequence — each step builds on the previous output.

### Before you start

- Create a new Figma file (or a new page in the existing Perspectize file)
- Each step below = one Figma Make prompt
- After each prompt, review the output, tweak if needed, then move to the next
- Name your layers as you go (Make often generates `Frame 47` — rename to semantic names)

---

### Step 1: Color Variables

> Create a set of color variables for a design system. Group them under "Theme / Light". Use these exact names and hex values:
>
> background: #ffffff
> foreground: #171717
> card: #ffffff
> card-foreground: #171717
> popover: #ffffff
> popover-foreground: #171717
> primary: #1a365d
> primary-hover: #2d3748
> primary-foreground: #ffffff
> secondary: #f5f5f5
> secondary-hover: rgba(245,245,245,0.8)
> secondary-foreground: #525252
> muted: #f5f5f5
> muted-foreground: #525252
> accent: #f7fafc
> accent-foreground: #2d3748
> destructive: #dc2626
> destructive-hover: #b91c1c
> destructive-foreground: #ffffff
> border: #d4d4d4
> input: #d4d4d4
> ring: #1a365d
> disabled: 0.5 (opacity value)
>
> Also create a "Rating" group:
> rating-positive: #16a34a
> rating-neutral: #ca8a04
> rating-negative: #dc2626
> rating-undecided: #64748b
>
> And a "Brand" group:
> logo-purple: #8b5cf6

---

### Step 2: Typography

> Set up a dual-font typography scale. Geist for headings and UI elements, Charter for body/content text. Create text styles:
>
> Headings (Geist — sans-serif):
> - Display: Geist 36px/40px bold
> - H1: Geist 30px/36px semibold
> - H2: Geist 24px/32px semibold
> - H3: Geist 20px/28px semibold
> - Label: Geist 14px/20px medium
> - Small/Caption: Geist 12px/16px regular, color #525252
>
> Body (Charter — serif):
> - Body Large: Charter 16px/24px regular
> - Body: Charter 14px/20px regular
> - Body Italic: Charter 14px/20px italic
> - Body Bold: Charter 14px/20px bold

---

### Step 3: Button Components

**NOTE: Your Figma file already has a comprehensive Button component with all variants (Default/Secondary/Outline/Ghost/Link/Destructive/Loading x Small/Default/Large/Icon x Default/Hover/Disabled). You can skip this step if working in the existing file. Only use this prompt if starting fresh.**

> Create a button component set with auto layout. Use Geist 14px medium text. Border radius 6px.
>
> Variants by style:
> - Primary: background #1a365d, text white. Hover: #2d3748.
> - Secondary: background #f5f5f5, text #525252. Hover: rgba(245,245,245,0.8).
> - Outline: background white, 1px border #d4d4d4, text #525252. Hover: background #f7fafc, text #2d3748.
> - Ghost: background transparent, text #171717. Hover: background #f7fafc, text #2d3748.
> - Destructive: background #dc2626, text white. Hover: #b91c1c.
> - Link: background transparent, text #1a365d. Hover: underline.
> - Loading: same as Primary but with spinner icon + "Loading" text.
>
> Variants by size:
> - Small: height 36px, padding 6px 12px
> - Default: height 40px, padding 8px 16px
> - Large: height 44px, padding 10px 32px
> - Icon: padding 12px, square, no text, centered icon
>
> States for each: default, hover, disabled (50% opacity), loading (swap text for a spinner icon + "Loading...")

---

### Step 4: Form Inputs

> Create input components with auto layout. Charter 14px regular for input text (what users type), Geist 14px medium for labels. Border radius 6px.
>
> Text Input:
> - Height 40px, 1px border #d4d4d4, background white, padding 0 12px
> - Placeholder text in #525252
> - States: default, focused (2px ring #1a365d around it), error (border #dc2626, red error text below), disabled (background #f5f5f5, 50% opacity text)
>
> Search Input:
> - Same as text input but with a search icon (magnifying glass) on the left inside the input
>
> Textarea:
> - Same border/color style as text input, but 80px tall and resizable. Use Charter 14px regular for textarea content (user-written text)
>
> Label:
> - Geist 14px medium, color #171717, placed above the input with 6px gap
>
> Select Dropdown (styled, not native):
> - Same dimensions as text input, with a chevron-down icon on the right
> - Show open state with a dropdown menu below (white background, border, shadow, rounded-lg)

---

### Step 5: Card Container

> Create a card component with auto layout.
> - Background: white (#ffffff)
> - Border: 1px solid #d4d4d4
> - Border radius: 8px
> - Shadow: 0 1px 3px rgba(0,0,0,0.1), 0 1px 2px rgba(0,0,0,0.06)
> - Internal padding: 0 (content controls its own padding)
>
> This card will wrap the data table. It should clip overflow and contain the table header, rows, and pagination footer.

---

### Step 6: Header (App Shell)

> Create a header bar component, 1440px wide, 64px tall, auto layout horizontal, items centered, space-between.
>
> - Background: #1a365d (navy)
> - Left side: "Perspectize" in Geist 18px bold, color white. Leave space for a logo icon to the left of the text later.
> - Right side: A styled select dropdown (white border, white text, transparent background) showing "Select user..." and a primary button "Add Video" (background white, text #1a365d — inverted colors since header is dark).
> - Horizontal padding: 32px
> - Bottom border: none (the navy background provides enough separation)
>
> Also create a 375px mobile variant:
> - "Perspectize" text at 16px
> - Right side: user icon button (24px, white) + plus icon button (24px, white) — no text labels
> - Horizontal padding: 16px

---

### Step 7: Activity Page — Desktop (1440px)

> Create a page frame at 1440x900. Background #ffffff.
>
> At the top, place the navy header component from step 6.
>
> Below the header, centered content area max 1280px wide with 32px horizontal padding:
>
> - Page title section: left side has "Activity" in Geist H1 (30px semibold, #171717) with "Recently updated content" in Charter body text (#525252) below it. Right side has the search input from step 4, 320px wide. They sit on the same horizontal line, the title left-aligned, search right-aligned.
>
> - Below that (16px gap), the card component from step 5 containing a data table:
>   - Table header row: background #1a365d, text white Geist 13px semibold. Columns: Title (flex), Type (100px), Duration (120px), Date Added (140px), Last Updated (140px).
>   - 10 data rows: alternating background — odd rows white, even rows #f7fafc. Row height ~44px. Text is Charter 14px regular #171717 (serif body font). Title column text is navy #1a365d (it's a link). Type shows "YOUTUBE" in Geist 12px #525252.
>   - Pagination footer: inside the card at the bottom. Light gray background #f5f5f5. Shows "Page Size: [10 v]  1 to 10 of 37  Page 1 of 4 [< >]". Text in #525252.
>
> Use 6 sample rows with plausible YouTube video titles.

---

### Step 8: Activity Page — Mobile (375px)

> Create a page frame at 375x812. Background #ffffff.
>
> Navy header from step 6 (mobile variant — icon buttons, no text labels).
>
> Content area with 16px horizontal padding:
>
> - "Activity" in H1 (24px bold) with subtitle below.
> - Search input at full width below the subtitle (8px gap).
> - Card container below (12px gap) containing the table:
>   - Only 2 columns visible: Title (flex) and Type (80px)
>   - Navy header row, same alternating row pattern
>   - Pagination footer wraps to 2 lines if needed
>
> The card should go edge-to-edge within the 16px gutters (so 343px wide).

---

### Step 9: Add Video Dialog — Desktop

> Create a frame at 1440x900. Show the Activity page from step 7 dimmed behind a semi-transparent black overlay (50% opacity).
>
> Centered on the overlay, a dialog card:
> - Width: 448px (max-w-md)
> - Background: solid white
> - Border radius: 12px
> - Shadow: 0 20px 60px rgba(0,0,0,0.15)
> - Padding: 24px
>
> Content:
> - Top right: X close button (ghost style, 16px icon)
> - Title: "Add Video" in Geist 18px semibold
> - Subtitle: "Paste a YouTube URL to add it to your library." in Charter 14px #525252
> - 16px gap
> - Label "YouTube URL" + text input with placeholder "https://www.youtube.com/watch?v=..."
> - 24px gap
> - Footer: right-aligned, "Cancel" button (outline) and "Add Video" button (primary navy), 8px gap between them
>
> Create 3 additional state variants side by side:
> 1. Filled: input has a URL typed in, Add Video button enabled
> 2. Loading: input disabled (dimmed), button shows "Adding..." with spinner
> 3. Error: input has red border, red text "Please enter a valid YouTube URL" below it

---

### Step 10: Add Video Dialog — Mobile (375px)

> Create a frame at 375x812. Same dimmed overlay pattern.
>
> Dialog centered vertically:
> - Width: 343px (375 - 32px margins)
> - Same white background, 12px radius, shadow, 24px padding
> - Same content as desktop but:
>   - Buttons stack vertically, full width
>   - "Add Video" button on top, "Cancel" below
>   - 8px gap between buttons

---

### Step 11: Add Perspective Dialog — Desktop

> Create a frame at 1440x900 with dimmed overlay.
>
> Centered dialog:
> - Width: 512px (max-w-lg)
> - White background, 12px radius, shadow, 24px padding
> - Scrollable content if it overflows
>
> Content from top to bottom:
> - Title: "Add Perspective" in Geist 18px semibold
> - Subtitle: "Share your view on this video" in Charter 14px #525252
> - 16px gap
>
> - Section: Video selector
>   - Label "Video"
>   - Styled dropdown showing "Search or select a video..." with search icon
>
> - 16px gap
> - Divider line (1px #d4d4d4)
> - 16px gap
>
> - Section: Quick Take
>   - Label "Agreement"
>   - Row of 3 toggle buttons: "Agree" (outline, when selected: green #16a34a bg), "Undecided" (outline, when selected: slate #64748b bg), "Disagree" (outline, when selected: red #dc2626 bg). All have white text when selected. Show "Agree" as selected.
>   - 12px gap
>   - Label "Comment (optional)"
>   - Textarea, 2 lines tall, placeholder "What stood out to you?"
>
> - 16px gap
> - Divider line (1px #d4d4d4)
> - 16px gap
>
> - Section: Detailed Ratings (optional)
>   - Three rating rows, each with:
>     - Label on left ("Quality", "Importance", "Confidence")
>     - A progress bar / slider (track in #d4d4d4, fill in #1a365d navy)
>     - Numeric value on right showing "720 / 1000" in Geist 14px
>   - 12px gap between each rating row
>
> - 16px gap
> - Label "Review (optional)"
> - Textarea, 4 lines tall, placeholder "Write a longer review..."
>
> - 24px gap
> - Footer: "Cancel" outline button + "Submit Perspective" primary button, right-aligned

---

### Step 12: Toast Notifications

> Create 3 toast notification components, each 360px wide, auto layout horizontal:
>
> - Background white, border 1px #d4d4d4, border radius 8px, shadow 0 4px 12px rgba(0,0,0,0.08), padding 12px 16px
>
> 1. Success: Green checkmark circle icon (#16a34a) on left, "Added: Video Title" in Geist 14px medium on right
> 2. Error: Red X circle icon (#dc2626) on left, "This video has already been added" in Geist 14px medium on right
> 3. Info: Navy info circle icon (#1a365d) on left, "Video added to your library" in Geist 14px medium on right
>
> Position them in the top-right corner of a 1440x900 frame with 16px margin from top and right edges. Stack vertically with 8px gap.

---

### Step 13: Empty & Loading States

> Create two state frames at 1440x900, same header + page layout as step 7 but with different content areas:
>
> 1. Empty state: Instead of the table card, show a centered block:
>    - A simple line illustration or icon (film reel or play button, in #525252)
>    - "No videos yet" in Geist 18px semibold #171717
>    - "Add your first YouTube video to get started." in Charter 14px #525252
>    - "Add Video" primary button below
>
> 2. Loading state: Same card container as step 7, but instead of real data:
>    - Table header row (navy, real column names)
>    - 5 rows of skeleton loading bars (rounded rectangles in #f5f5f5, pulsing animation). Each row has rectangles roughly matching the column widths.

---

### Tips for between steps

- **After each step**, rename Figma layers to semantic names: `Header`, `PageContent`, `CardContainer`, `TableRow`, `DialogOverlay`, etc.
- **Convert repeated elements** to Figma components (buttons, inputs, table rows) so later steps can instance them.
- **Check your colors** — if Make generates a slightly different hex, swap it to the variable from step 1.
- **Steps 7-13 reference earlier components** — if Make doesn't pick them up, manually swap in your button/input/card components from steps 3-5.

---

## Implementation Priority

When you design in Figma, work in this order:

| Priority | Figma Page | Why |
|----------|-----------|-----|
| 1 | **04 — Components** | Foundation. Get buttons, inputs, cards, and color tokens right first. Everything else uses these. |
| 2 | **01 — Activity** (desktop) | Main page. Validate the grid layout, card container, and header in one design. |
| 3 | **02 — Add Video Dialog** | Currently broken. Needs a solid design to fix Phase 3.1. |
| 4 | **01 — Activity** (mobile 375px) | Verify the responsive behavior. |
| 5 | **03 — Add Perspective** | Needed for Phase 4. More complex form, design early. |
| 6 | **05 — Toasts & States** | Polish layer. |

---

## Key Layout Mistakes to Avoid

1. **Always design with gutters**: 16px on mobile, 24px tablet, 32px desktop. The PageWrapper provides these — but content inside (like the AG Grid card) should respect them.

2. **Cards need breathing room**: The AG Grid should not touch the edges of the PageWrapper. Wrap it in a card with internal padding, or at minimum give it a visible border and rounded corners so it feels contained.

3. **The header needs visual weight**: A white header with a thin border looks like a wireframe. A navy header immediately signals "this is a real app."

4. **Dialogs must be opaque**: The see-through dialog is a code bug (missing token), but design-wise, always show dialogs as solid white on a dimmed backdrop. No glass-morphism, no transparency tricks for forms.

5. **Don't rely on native form controls**: The `<select>` for user selection looks out of place. Design a styled dropdown (shadcn Select or Combobox) that matches the button and input styling.

6. **Consistent border radius**: Pick one radius and use it everywhere — `rounded-lg` (8px) for cards, `rounded-md` (6px) for inputs/buttons, `rounded-xl` (12px) for dialogs. Don't mix.
