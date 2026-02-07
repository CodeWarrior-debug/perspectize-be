# Pending Design Decisions

Design patches that work but need a proper long-term solution.

| ID | Problem | Patch | SHA | Long-term Solution |
|----|---------|-------|-----|-------------------|
| PDD-001 | Sticky header used `bg-background` but `--color-background` was never defined in `app.css`, causing transparent header and scroll bleed-through | Changed `bg-background` to `bg-white` in `Header.svelte` | b42c457 | Define complete color theme in `app.css` (`--color-background`, `--color-foreground`, `--color-border`, etc.) so semantic utilities resolve correctly and support dark mode. Revert header to `bg-background`. |
