# Pending Design Decisions

Design patches that work but need a proper long-term solution.

---

## PDD-001: Header background color hardcoded to white

**Problem:** The sticky header used `bg-background` (Tailwind utility), but `--color-background` was never defined in `app.css`. This caused the header to be transparent, letting page content bleed through when scrolling.

**Patch:** Changed `bg-background` to `bg-white` in `Header.svelte`.

**SHA:** b42c457

**Long-term solution:** Define a complete color theme in `app.css` under `@theme` with `--color-background`, `--color-foreground`, `--color-border`, `--color-muted`, etc. This would let all `bg-background`, `text-foreground`, `border-border` utilities resolve correctly and support dark mode in the future. The header should then revert to `bg-background`.

---
