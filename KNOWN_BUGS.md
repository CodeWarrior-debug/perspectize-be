# Known Bugs

Issues discovered during development and testing.

| ID | Priority | Description | Where | Found | Phase |
|----|----------|-------------|-------|-------|-------|
| BUG-001 | P1 | Header overflows at 375px — "Perspectize" truncated, "Add Video" button clipped off right edge | `Header.svelte` | 2026-02-07 | 2.1 |
| BUG-002 | P1 | Pagination bar broken at 375px — "Page Size:" truncated to "e Size:", count text wraps across lines | `ActivityTable.svelte` (AG Grid pagination) | 2026-02-07 | 2.1 |
| BUG-003 | P1 | Sticky header clipping persists on scroll at 375px — same overflow visible throughout page | `Header.svelte` | 2026-02-07 | 2.1 |
| BUG-004 | P2 | No responsive header collapse — header needs to stack, shrink, or use hamburger menu at mobile widths | `Header.svelte` | 2026-02-07 | 2.1 |
| BUG-005 | P2 | Table content left-shifted beyond viewport at 375px — rows start outside visible area when scrolled down | `ActivityTable.svelte` | 2026-02-07 | 2.1 |
| BUG-006 | P3 | No visual affordance for hidden columns — only Title visible at 375px with no hint that horizontal scroll reveals more | `ActivityTable.svelte` | 2026-02-07 | 2.1 |
