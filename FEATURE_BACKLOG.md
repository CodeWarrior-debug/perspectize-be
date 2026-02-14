# Feature Backlog

Ideas and future enhancements captured during development. Not committed to any milestone — evaluated when planning new work.

---

## Discover Page (New Page)

The v1 home page is an **Activity** page — a data table of user activity on videos already in the system. This is the correct approach for v1.

A future **Discover** page would be a separate page for finding new content outside the system:
- **Browse** — Show topics/tags from YouTube API endpoint, letting users explore categories
- **Search** — Live YouTube API search to discover new videos directly from YouTube

This is distinct from the Activity page's local search/filter. Discover reaches out to YouTube; Activity shows what's already tracked.

---

## Decide on Table Enhancement Libraries Long-Term

Evaluate whether to adopt a table enhancement library (AdaptableTools, TanStack Table, or similar) or continue building AG Grid features incrementally as the frontend matures.

**Research:** See [AdaptableTools evaluation](docs/research/2026-02-13-adaptabletools-evaluation.md) and the associated PR for full analysis.

**Current decision:** Stay with AG Grid Community + `ag-grid-svelte5`, build features as needed (saved views, column chooser, conditional formatting). AdaptableTools has no Svelte wrapper and costs $6K-$21K+/year.

**Revisit triggers:**
- 5+ complex grid views in the app
- Framework migration away from Svelte
- Enterprise client requirements for financial-grade grid features
- AdaptableTools releases Svelte support
