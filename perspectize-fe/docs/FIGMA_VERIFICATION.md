# Figma Make Verification Guide

> How to verify that each step of the Figma Make workflow was applied correctly
> in the "Perspectize Youtube - Design 1" file.

---

## Prerequisites

- Open the correct file: **Perspectize Youtube - Design 1** in the **Perspectize** project under the **Edify+Ente...** team.
- Confirm you are in the right file by checking the file name in the tab bar and verifying it was recently edited (right-click the file name and select "Show version history").

---

## 1. Verifying Color Variables (Step 1)

### Where to look

1. Open the file in Figma.
2. Click the **Variables** icon in the right panel (a small diamond/rhombus shape), or use the menu: **Window > Variables**.
3. You should see **3 variable collections** listed.

### Expected collections and values

#### Collection: Theme / Light (22 variables)

| Variable Name | Expected Hex | Usage |
|---------------|-------------|-------|
| `background` | `#FFFFFF` | Page background |
| `foreground` | `#0A0A0A` | Default text |
| `card` | `#FFFFFF` | Card/surface background |
| `card-foreground` | `#0A0A0A` | Text on cards |
| `popover` | `#FFFFFF` | Dropdown/dialog background |
| `popover-foreground` | `#0A0A0A` | Text in popovers |
| `primary` | `#1A365D` | Buttons, links, header (navy) |
| `primary-foreground` | `#F0F4F8` | Text on primary |
| `secondary` | `#F5F5F5` | Secondary buttons, tags |
| `secondary-foreground` | `#171717` | Text on secondary |
| `muted` | `#F5F5F5` | Subtle backgrounds |
| `muted-foreground` | `#525252` | Placeholder text, captions |
| `accent` | `#F5F5F5` | Hover highlights |
| `accent-foreground` | `#171717` | Text on accent |
| `destructive` | `#DC2626` | Error states, delete actions |
| `destructive-foreground` | `#FFFFFF` | Text on destructive |
| `border` | `#E5E5E5` | Default borders |
| `input` | `#E5E5E5` | Input borders |
| `ring` | `#1A365D` | Focus ring (matches primary) |
| `chart-1` | `#1A365D` | Chart color 1 |
| `chart-2` | `#2563EB` | Chart color 2 |
| `chart-3` | `#525252` | Chart color 3 |

#### Collection: Rating (4 variables)

| Variable Name | Expected Hex | Usage |
|---------------|-------------|-------|
| `positive` | `#16A34A` | High scores, Agree |
| `negative` | `#DC2626` | Low scores, Disagree |
| `neutral` | `#525252` | Mid scores |
| `undecided` | `#D4D4D4` | No opinion |

#### Collection: Brand (1 variable)

| Variable Name | Expected Hex | Usage |
|---------------|-------------|-------|
| `logo-purple` | `#7C3AED` | Logo primary color |

### How to verify a variable value

1. Click on a variable name in the Variables panel.
2. The color swatch and hex value appear to the right.
3. Click the swatch to open the color picker and confirm the exact hex.

**Total expected: 27 color variables across 3 collections.**

---

## 2. Verifying Text Styles (Step 2)

### Where to look

**Option A -- Right panel:**
1. Select any text element on the canvas.
2. In the right panel under "Text", click the style name (four-dot icon).
3. Browse available text styles.

**Option B -- Assets panel:**
1. Open the Assets panel (left sidebar, book icon, or press **Alt+2** / **Opt+2**).
2. Search for style names like "Display", "H1", "Body".
3. Text styles appear under the "Text styles" section.

**Option C -- Local styles dialog:**
1. Click the canvas background (deselect everything).
2. In the right panel, look for "Local styles" or use the menu: **File > Local styles**.

### Expected text styles (10 total)

#### Geist styles (6)

| Style Name | Font | Size | Line Height | Weight |
|-----------|------|------|-------------|--------|
| Display | Geist | 36px | 40px | Bold (700) |
| H1 | Geist | 30px | 36px | SemiBold (600) |
| H2 | Geist | 24px | 32px | SemiBold (600) |
| H3 | Geist | 20px | 28px | SemiBold (600) |
| Label | Geist | 14px | 20px | Medium (500) |
| Small / Caption | Geist | 12px | 16px | Regular (400) |

#### Charter styles (4)

| Style Name | Font | Size | Line Height | Weight |
|-----------|------|------|-------------|--------|
| Body Large | Charter | 16px | 24px | Regular (400) |
| Body | Charter | 14px | 20px | Regular (400) |
| Body Italic | Charter | 14px | 20px | Italic (400i) |
| Body Bold | Charter | 14px | 20px | Bold (700) |

### How to verify a text style

1. Find the style in the Assets panel or local styles.
2. Hover over it to see a preview.
3. Click the style to inspect its properties: font family, weight, size, and line height.

**If Geist or Charter fonts are missing:** These fonts must be installed locally or available via a Figma plugin. Geist is available from [vercel.com/font](https://vercel.com/font). Charter is available from [Matthew Carter / Bitstream](https://practicaltypography.com/charter.html). Figma will show a missing font indicator (yellow warning) if the font is not available.

---

## 3. Verifying Components (Steps 3--13)

### Where to look

1. Open the **Assets** panel (left sidebar, book icon).
2. Components created by Figma Make appear under "Local components".
3. Search by name (e.g., "Button", "Card", "Dialog").
4. Full-page designs may appear as separate **Pages** in the left sidebar page list.

### How to verify a component

1. Find the component in the Assets panel.
2. Drag it onto the canvas to create an instance.
3. In the right panel, check for **Variants** (if the component has multiple states like Default/Hover/Disabled or size variants).
4. Verify the component uses the correct color variables and text styles by selecting inner elements and checking their fill/text references.

---

## 4. If Variables or Styles Are NOT Visible

If you open the Variables panel or look for text styles and see nothing, work through this troubleshooting checklist:

### Check 1: Correct file

Make sure you opened "Perspectize Youtube - Design 1" and not a different file. Look at the tab title in Figma.

### Check 2: New page created by Make

Figma Make sometimes adds output to a new page rather than applying to existing ones. Check the **Pages list** in the left sidebar. Look for pages named something like "Make Output", "Generated", or similar.

### Check 3: Version history

Right-click the file name in the tab bar and select **Show version history**. Look for a recent auto-save entry around the time you ran the Make prompt. If there is no new version, Make did not modify the file.

### Check 4: Variables vs. Color Styles

Figma has two separate systems for colors:

| System | Panel Location | Purpose |
|--------|---------------|---------|
| **Variables** | Variables panel (diamond icon) | Modern system, supports modes (light/dark), used by Make for tokens |
| **Color Styles** | Right panel > "Color styles" section | Legacy system, still widely used |

Make might have created Color Styles instead of Variables, or vice versa. Check both panels.

### Check 5: Code output instead of applied changes

If Make displayed code (CSS, JSON, or similar) in the chat instead of applying changes to the file, it did not modify Figma. This happens when the prompt is ambiguous. See the "Tips for Better Results" section below for rephrasing guidance.

### Check 6: Re-run the prompt

If nothing was applied, re-run the Make prompt with more explicit language. For example, instead of "Create color tokens", say "Create Figma color variables in this file with the following values...".

---

## 5. Manual Creation (Fallback)

If Figma Make does not apply changes correctly, you can create them by hand.

### Creating color variables manually

1. Open the **Variables** panel (diamond icon in the right panel, or **Window > Variables**).
2. Click the **+** button to create a new collection. Name it (e.g., "Theme / Light").
3. Click **+ Create variable** inside the collection.
4. Set the variable name (e.g., `background`).
5. Click the color swatch to set the hex value.
6. Repeat for each variable.

### Creating text styles manually

1. Create a text element on the canvas.
2. Set the font family, weight, size, and line height to match the spec.
3. With the text selected, click the **four-dot icon** (style selector) in the Text section of the right panel.
4. Click the **+** button to create a new text style.
5. Name it (e.g., "H1").
6. Repeat for each style.

### Creating components manually

1. Design the element on the canvas (e.g., a button with background fill, text, padding).
2. Select all layers that make up the element.
3. Right-click and select **Create component** (or press **Ctrl/Cmd + Alt + K**).
4. Name the component (e.g., "Button").
5. To add variants: select the component, click **+** next to "Variants" in the right panel.

---

## 6. Verification Checklist (All 13 Steps)

Use this table after completing each step. Mark the "Verified" column once you confirm the output in Figma.

| Step | Output | Where to Check | What to Confirm | Verified |
|------|--------|----------------|-----------------|----------|
| 1. Color Variables | 27 color variables in 3 collections | Variables panel (diamond icon) | Collections: Theme/Light (22), Rating (4), Brand (1). Hex values match spec. | [ ] |
| 2. Typography | 10 text styles (6 Geist + 4 Charter) | Assets panel > Text styles | Style names, font families, sizes, weights, line heights all match spec. | [ ] |
| 3. Buttons | Button component with 6 variants | Assets panel > Components | Variants: Default, Secondary, Destructive, Outline, Ghost, Link. States: default, hover, disabled. Uses `primary` variable for fill. | [ ] |
| 4. Form Inputs | Input, Label, Textarea components | Assets panel > Components | Input has border using `input` variable, focus ring using `ring` variable. Label uses Geist Medium 14px. | [ ] |
| 5. Card Components | VideoCard, ReviewCard components | Assets panel > Components | Card background uses `card` variable. Border uses `border` variable. Correct padding and corner radius. | [ ] |
| 6. Dialog | Dialog component with overlay | Assets panel > Components | Semi-transparent overlay backdrop. Dialog surface uses `card` variable. Close button present. | [ ] |
| 7. Header | Header component with nav | Assets panel > Components | Background uses `primary` (#1A365D). Text/icons use `primary-foreground`. Height 64px. Logo left, actions right. | [ ] |
| 8. Home Page | Full page frame at 1440px wide | Pages list (left sidebar) | Contains Header + PageWrapper + content area. Max-width 1280px centered. | [ ] |
| 9. Add Video Dialog | Dialog overlay on page | Pages list (left sidebar) | Dialog centered over page with backdrop. Form fields for YouTube URL. Uses Dialog component from Step 6. | [ ] |
| 10. Video Detail | Video detail page frame | Pages list (left sidebar) | Video player/thumbnail area, metadata section, perspectives list. Uses Card and typography styles. | [ ] |
| 11. User Profile | Profile page frame | Pages list (left sidebar) | User info section, activity/history list. Uses established components. | [ ] |
| 12. Responsive | Mobile (375px) and tablet (768px) variants | Pages list (left sidebar) | Key pages adapted for smaller viewports. Header collapses. Grid stacks to single column. | [ ] |
| 13. Empty/Loading | Skeleton and empty state components | Assets panel > Components | Skeleton placeholder shapes (pulsing rectangles). Empty state with illustration/message and CTA. | [ ] |

---

## 7. Tips for Better Figma Make Results

### Prompt phrasing

| Instead of... | Say... |
|---------------|--------|
| "Create color tokens" | "Create Figma color variables in this file" |
| "Add typography" | "Create local text styles in this Figma file" |
| "Make a button" | "Create a Button component in this Figma file with variants for Default, Secondary, and Destructive" |
| "Design the home page" | "Create a new page called 'Home' with a frame at 1440x900px containing..." |

### General guidelines

- **Be explicit about the output type.** Say "Figma variables", "local text styles", or "component with variants" -- not generic terms like "tokens" or "design system".
- **Specify the file context.** Start with "In this design file..." to anchor Make to the current file.
- **One step at a time.** Do not combine multiple steps into a single prompt. Run Step 1 and verify, then Step 2 and verify, and so on.
- **If Make outputs code:** It misunderstood the intent. Rephrase with "Apply directly to this Figma file as [variables / styles / components]. Do not output code."
- **Include exact values.** Provide hex codes, font sizes, and weights in the prompt rather than referencing external docs.
- **Name things explicitly.** Provide the exact names you want for variables, styles, and components.

### After each step

1. Check the relevant panel (Variables, Assets, Pages) immediately after running the prompt.
2. If the output is missing, check the troubleshooting steps in Section 4.
3. If the output is present but values are wrong, you can edit variables/styles directly in their panels rather than re-running the prompt.
4. Do not proceed to the next step until the current step is verified.

---

## 8. Quick Reference: Figma Panel Locations

| Panel | How to Open | What It Shows |
|-------|------------|---------------|
| Variables | Right panel > diamond icon, or **Window > Variables** | Color variables organized by collection |
| Color Styles | Right panel > circle/palette icon (when no element selected) | Legacy color style definitions |
| Text Styles | Assets panel > search, or right panel "Text" section > four-dot icon | Text style definitions |
| Assets | Left sidebar > book icon, or **Alt/Opt + 2** | Components, styles, and variables |
| Pages | Left sidebar > top section | List of pages in the file |
| Version History | Right-click file tab > "Show version history" | Timestamped file snapshots |
| Layers | Left sidebar > layers icon, or **Alt/Opt + 1** | Layer hierarchy of current page |
