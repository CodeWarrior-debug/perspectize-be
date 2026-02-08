# Figma Reference

> Persistent reference for Claude to understand the Perspectize Figma setup.
> Updated as design system evolves.

---

## File Info

| Field | Value |
|-------|-------|
| **Team** | Edify+Ente (Professional plan) |
| **Project** | Perspectize |
| **Design File** | Perspectize Youtube - Design 1 |
| **File Key** | `SyvrP9yYbrmCorofJK4Co8` |
| **URL** | `https://www.figma.com/design/SyvrP9yYbrmCorofJK4Co8` |
| **Base Kit** | Radius 3.0 / shadcn-ui v1.2 (Tailwind CSS variables) |

### Related Files

| File | Purpose |
|------|---------|
| Perspectize Youtube - App 1 | Published component library |
| Perspectize - Radix 3.0 Implementation | Radix/shadcn base kit reference |
| Radius 3.0: Design System/UI Kit | Original UI kit source |

---

## Pages (Known)

| Page | Contents | Key Node IDs |
|------|----------|-------------|
| Typography | Tailwind font size scale (text-xs → text-9xl), font weight scale (thin → black), Charter preview | `0:1` |

> **TODO**: Add remaining pages as they are discovered or reported by user.

---

## Design System State

### Fonts

| Font | Role | Weights | Source |
|------|------|---------|--------|
| **Geist** | Headings, UI labels, buttons, navigation | 400 (Regular), 500 (Medium), 600 (SemiBold), 700 (Bold) | [vercel.com/font](https://vercel.com/font) — Variable woff2 |
| **Charter** | Body text, reviews, descriptions, long-form content | 400 (Regular), 400i (Italic), 700 (Bold) | [Matthew Carter / Bitstream](https://practicaltypography.com/charter.html) |

### Color Variables

Status: **Step 1 completed via Figma Make** — awaiting verification.

Expected collections (from DESIGN_SPEC.md):

| Collection | Count | Key Variables |
|------------|-------|---------------|
| Theme / Light | 22 | `primary` (#1A365D), `background` (#FFFFFF), `foreground` (#0A0A0A), `destructive` (#DC2626), `muted` (#F5F5F5), `border` (#E5E5E5) |
| Rating | 4 | `positive` (#16A34A), `negative` (#DC2626), `neutral` (#525252), `undecided` (#D4D4D4) |
| Brand | 1 | `logo-purple` (#7C3AED) |

> Full token list: see [DESIGN_SPEC.md](DESIGN_SPEC.md) § Color Tokens

### Text Styles

Status: **Step 2 pending** — not yet applied.

Target styles (from DESIGN_SPEC.md):

| Style | Font | Size/Line Height | Weight |
|-------|------|-----------------|--------|
| Display | Geist | 36/40 | Bold |
| H1 | Geist | 30/36 | SemiBold |
| H2 | Geist | 24/32 | SemiBold |
| H3 | Geist | 20/28 | SemiBold |
| Label | Geist | 14/20 | Medium |
| Small/Caption | Geist | 12/16 | Regular |
| Body Large | Charter | 16/24 | Regular |
| Body | Charter | 14/20 | Regular |
| Body Italic | Charter | 14/20 | Italic |
| Body Bold | Charter | 14/20 | Bold |

### Components

Status: **Not yet created** — Steps 3-13 pending.

---

## Code ↔ Figma Component Mapping

Maps Svelte components to their Figma counterparts. Updated as components are built.

| Svelte Component | Figma Component | Figma Node ID | Notes |
|-----------------|----------------|---------------|-------|
| `$lib/components/shadcn/button` | Button | — | Steps 3 (pending) |
| `$lib/components/shadcn/input` | Input | — | Step 4 (pending) |
| `$lib/components/shadcn/label` | Label | — | Step 4 (pending) |
| `$lib/components/shadcn/dialog` | Dialog | — | Step 6 (pending) |
| `$lib/components/Header.svelte` | Header | — | Step 7 (pending) |
| `$lib/components/AddVideoDialog.svelte` | Add Video Dialog | — | Step 9 (pending) |
| `$lib/components/ActivityTable.svelte` | Activity Table / AG Grid | — | — |
| `$lib/components/UserSelector.svelte` | User Selector | — | — |
| `$lib/components/PageWrapper.svelte` | Page Wrapper | — | — |

---

## Code ↔ Figma Token Mapping

Maps CSS custom properties to Figma variable names.

| CSS Token (`app.css`) | Figma Variable | Collection |
|-----------------------|---------------|------------|
| `--color-primary` | `primary` | Theme / Light |
| `--color-primary-foreground` | `primary-foreground` | Theme / Light |
| `--color-background` | `background` | Theme / Light |
| `--color-foreground` | `foreground` | Theme / Light |
| `--color-muted` | `muted` | Theme / Light |
| `--color-muted-foreground` | `muted-foreground` | Theme / Light |
| `--color-destructive` | `destructive` | Theme / Light |
| `--color-border` | `border` | Theme / Light |
| `--color-input` | `input` | Theme / Light |
| `--color-ring` | `ring` | Theme / Light |
| `--color-rating-positive` | `positive` | Rating |
| `--color-rating-negative` | `negative` | Rating |
| `--color-rating-neutral` | `neutral` | Rating |
| `--color-rating-undecided` | `undecided` | Rating |

> Pattern: Figma variable `foo` → CSS `--color-foo` → Tailwind class `bg-foo` / `text-foo`

---

## MCP Tool Usage

### Available Tools

| Tool | Purpose | Auto-approved |
|------|---------|--------------|
| `get_screenshot` | Capture visual of any node | Yes |
| `get_design_context` | Extract code/structure from a node | Yes |
| `get_metadata` | Get XML structure overview of a node/page | Yes |
| `get_variable_defs` | Extract color/spacing variables (needs selected layer in design mode) | Yes |
| `get_code_connect_map` | Get existing Code Connect mappings | Yes |
| `add_code_connect_map` | Link code components to Figma components | Yes |
| `send_code_connect_mappings` | Push Code Connect mappings to Figma | Yes |
| `create_design_system_rules` | Generate design system rules file | Yes |
| `whoami` | Check authenticated Figma user | Yes |

### Known Limitations

- `get_variable_defs` requires a layer selected in **design mode** (not dev mode)
- `get_metadata` with `nodeId: 0:0` (document root) returns an error — use page IDs like `0:1`
- File key `SyvrP9yYbrmCorofJK4Co8` is hardcoded for the main design file

---

## Figma Make Workflow Progress

Tracks the 13-step design build process. See [DESIGN_SPEC.md](DESIGN_SPEC.md) for full prompts.

| Step | Name | Status |
|------|------|--------|
| 1 | Color Variables | Done (unverified) |
| 2 | Typography | Pending |
| 3 | Buttons | Pending |
| 4 | Form Inputs | Pending |
| 5 | Card Components | Pending |
| 6 | Dialog | Pending |
| 7 | Header | Pending |
| 8 | Home Page | Pending |
| 9 | Add Video Dialog | Pending |
| 10 | Video Detail | Pending |
| 11 | User Profile | Pending |
| 12 | Responsive Variants | Pending |
| 13 | Empty/Loading States | Pending |
