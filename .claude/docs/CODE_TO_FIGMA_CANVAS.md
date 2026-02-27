# Code-to-Figma Canvas Workflow

Workflow for capturing the running Perspectize frontend and pushing it into Figma to keep designs in sync with code.

## When to Use

- Figma designs are **outdated** relative to the running app
- You need **real data** in Figma (thumbnails, actual content, live pagination)
- New features were built code-first and need Figma documentation
- Responsive variants need updating after frontend changes

## Key Files

| File | Purpose |
|------|---------|
| `frontend/docs/FIGMA.md` | Figma file keys, pages, variables, code-to-Figma mapping |
| `frontend/docs/DESIGN_SPEC.md` | Design tokens, typography, layout specs |
| `.claude/docs/FIGMA_DESIGN_SYSTEM_RULES.md` | Figma-to-code translation rules (reverse direction) |

## Figma Files

| File | Key | Role in this workflow |
|------|-----|---------------------|
| Design 1 | `K1HaZLeNwCckWvhoyAfRhj` | **Target** — captured pages go here |
| Radix 3.0 | `SyvrP9yYbrmCorofJK4Co8` | Reference — design system variables |
| App 1 | `dAiiWM7FOsob5upzUjtocY` | Published Make components |

---

## Pre-Flight Checklist

Before starting a capture session:

1. **Figma MCP connected** — plugin running in Figma desktop, token not expired
2. **Dev server running** — `pnpm run dev` in `frontend/` (http://localhost:5173)
3. **Backend running** — `go run ./cmd/server` in `backend/` (real data needed for thumbnails)
4. **Test data seeded** — enough rows to show pagination, varied content types
5. **Design 1 file open** in Figma desktop app

If Figma MCP token is expired, you'll get: `"MCP server requires re-authorization"`. Reconnect in Figma desktop before proceeding.

---

## Workflow Steps

### Step 1: Plan Capture Targets

Decide which pages/breakpoints to capture. Standard set:

| Page | Breakpoints | Width | Notes |
|------|------------|-------|-------|
| Activity (home) | Mobile | 375px | 2-column grid, icon buttons |
| Activity (home) | Tablet | 768px | 5-column grid, full header |
| Activity (home) | Desktop | 1440px | All visible columns |
| Add Video Dialog | Mobile | 375px | Stacked buttons |
| Add Video Dialog | Desktop | 1440px | Side-by-side buttons |

### Step 2: Identify Existing Figma Frames to Replace

Before capturing, note which frames have **good layer names** worth preserving.

```bash
# Get the structure of the target node
# Use get_metadata to see layer names without downloading full design context
```

Document the naming map:

| Layer Path | Semantic Name | Keep? |
|-----------|---------------|-------|
| `Activity Page - Mobile > Header` | Header | Yes |
| `Activity Page - Mobile > Header > Container > Link > Logo` | Logo | Yes |
| `Activity Page - Mobile > Header > Container > Container > Button > Icon` | Icon | Yes |
| `Activity Page - Mobile > Header > Container > Container > User dropdown` | User dropdown | Yes |

### Step 3: Capture Each Breakpoint

Use `generate_figma_design` to capture from localhost. The tool:
1. Opens a browser at the target URL
2. Renders the page
3. Pushes the result into Figma

**First call** — get capture instructions (no `outputMode`):

```
generate_figma_design()
```

This returns capture instructions and a `captureId`.

**Subsequent calls** — push to existing file:

```
generate_figma_design(
  captureId: "<from-first-call>",
  outputMode: "existingFile",
  fileKey: "K1HaZLeNwCckWvhoyAfRhj"
)
```

**For each breakpoint**, resize the browser before capturing:
- Use Chrome DevTools MCP `resize_page` to set viewport width
- Or let the capture tool handle it if it supports width parameters

**Poll for completion** — call with just the `captureId` every 5 seconds (up to 10 times) until status is `completed`.

### Step 4: Rename Layers (Preserve Good Names)

After capture, the generated frames will have auto-generated names. Rename them to match the existing naming convention:

**Manual process in Figma:**
1. Keep old frame visible as reference (side by side)
2. Walk the new frame's layer tree
3. Rename each layer to match the semantic name from Step 2
4. Pay special attention to: Header, Logo, Buttons, User dropdown, Search input, Table container

**Naming conventions to follow:**

| Element | Name Pattern | Example |
|---------|-------------|---------|
| Page frames | `{Page} - {Breakpoint}` | `Activity Page - Mobile` |
| Layout sections | Descriptive noun | `Header`, `Container`, `Footer` |
| Interactive elements | Component type | `Button`, `Link`, `Input` |
| Text elements | Content role | `Logo`, `User dropdown`, `Page title` |
| Data regions | Data type | `Table`, `Row`, `Cell` |

### Step 5: Organize in Figma

1. **Move captured frames** to the correct Figma page (e.g., "youtube initial" or create new pages per DESIGN_SPEC)
2. **Delete old outdated frames** once names are transferred
3. **Group breakpoints** — place Mobile/Tablet/Desktop side by side with 100px gap
4. **Add frame annotations** if needed (Figma comments or sticky notes)

### Step 6: Convert Repeated Elements to Components

Promote repeated elements so future captures can swap instances:

| Element | Component Name | Variants |
|---------|---------------|----------|
| Header bar | `Header` | Mobile, Desktop |
| Action buttons | `Button` | Primary, Secondary, Outline, Ghost, Icon-only |
| Search input | `SearchInput` | Default, Focused, Filled |
| Table row | `TableRow` | Default, Hover, Alternating |
| Pagination bar | `Pagination` | — |

### Step 7: Verify & Update References

After capture is complete:

1. **Update `frontend/docs/FIGMA.md`** — add new node IDs for captured frames to the mapping table
2. **Update Code-to-Figma mapping** — ensure `Code ↔ Figma Component Mapping` table reflects current state
3. **Screenshot verification** — use `get_screenshot` on the new frames to confirm they look correct

---

## Handling Specific Scenarios

### Scenario: Thumbnails Missing in Figma

YouTube thumbnails render in the live app but Figma captures them as rasterized images. After capture:
- Thumbnails appear as embedded images in the Figma frame
- They're real (not placeholder) because the capture renders the actual page
- If thumbnails fail to load, ensure the backend is running and content has `thumbnailUrl` populated

### Scenario: AG Grid Doesn't Render Fully

AG Grid uses canvas rendering which may not capture perfectly. If the grid appears blank or partial:
1. Add a short delay before capture (let AG Grid hydrate)
2. Ensure data is loaded (check network tab)
3. Consider capturing with Chrome DevTools screenshot as a fallback, then manually importing

### Scenario: Preserving Names Across Multiple Captures

If you recapture frequently, maintain a **name mapping file** to automate renaming:

```
# .figma/layer-names.md (local reference, not committed)
Activity Page - Mobile:
  Frame 1 → Header
  Frame 2 → Content Area
  Frame 3 → Table Card
  ...
```

### Scenario: Capture into New File vs Existing File

| Mode | When to use |
|------|-------------|
| `existingFile` | Updating Design 1 with current code state (default) |
| `newFile` | Creating a separate "Code Snapshot" file for comparison |
| `clipboard` | Quick grab to paste into any file |

---

## Quick Reference: Tool Calls

```
# 1. Start capture (get instructions + captureId)
generate_figma_design()

# 2. Push to existing Design 1 file
generate_figma_design(
  captureId: "xxx",
  outputMode: "existingFile",
  fileKey: "K1HaZLeNwCckWvhoyAfRhj"
)

# 3. Poll for completion
generate_figma_design(captureId: "xxx")

# 4. Verify with screenshot
get_screenshot(
  fileKey: "K1HaZLeNwCckWvhoyAfRhj",
  nodeId: "<new-frame-node-id>"
)

# 5. Get metadata to see layer names
get_metadata(
  fileKey: "K1HaZLeNwCckWvhoyAfRhj",
  nodeId: "<new-frame-node-id>"
)
```

---

## Relationship to Other Workflows

| Direction | Doc | When |
|-----------|-----|------|
| **Code → Figma** (this doc) | `CODE_TO_FIGMA_CANVAS.md` | Figma is outdated, sync from code |
| **Figma → Code** | `FIGMA_DESIGN_SYSTEM_RULES.md` | Implementing new designs from Figma |
| **Design Spec** | `frontend/docs/DESIGN_SPEC.md` | Token/component reference (source of truth for tokens) |
| **Figma Reference** | `frontend/docs/FIGMA.md` | File keys, node IDs, variable mapping |
