# AdaptableTools (AG Grid Extension) — Evaluation

**Date:** 2026-02-13
**Status:** Research complete — not adopting at this time
**Revisit when:** Framework migration, 5+ complex grids, or enterprise client requirements

---

## What Is AdaptableTools?

[AdaptableTools](https://www.adaptabletools.com/) is a commercial extension for AG Grid that provides a power-user UI layer on top of AG Grid. It adds:

- **Advanced filtering** — query builder, named filters, saved filter sets
- **Calculated columns** — formula-based virtual columns without backend changes
- **Custom editors** — inline editing with validation
- **Alerts & notifications** — conditional alerts when data changes
- **Dashboard/layout management** — end users save, name, and switch between grid configurations
- **Export** — enhanced Excel export with formatting
- **Cell styling** — conditional formatting, gradients, icons

## Framework Support

| Framework | Package | Status |
|-----------|---------|--------|
| React | `@adaptabletools/adaptable-react-aggrid` | Supported |
| Angular | `@adaptabletools/adaptable-angular-aggrid` | Supported |
| Vue | Available | Supported |
| **Svelte** | **None** | **Not supported** |
| Vanilla JS | `@adaptabletools/adaptable` | Supported (framework-agnostic) |

## Pricing

| License Tier | Cost | Limitations |
|-------------|------|-------------|
| Developer/Evaluation | Free | POCs only — cannot ship to production |
| Trial | Free | Full features, time-limited, watermark |
| Startup | Heavily discounted (3-year) | Business < 5 years, revenue < GBP 500K |
| Commercial | Contact sales (est. $5K-$20K+/yr) | Per-application, unlimited users |

**Important:** AdaptableTools license does NOT include AG Grid. AG Grid Enterprise ($999/dev/year) is a separate purchase and is needed for many of AdaptableTools' best features.

## Svelte Integration Analysis

### The Core Problem

Both `ag-grid-svelte5` and AdaptableTools want to own AG Grid instance creation:

```
ag-grid-svelte5:  Component → creates AG Grid → Svelte manages reactivity
AdaptableTools:   AdaptableTools.init({ gridOptions }) → creates AG Grid → imperative API
```

### Option 1: Ditch ag-grid-svelte5, Use Vanilla JS

Mount AdaptableTools to a raw DOM element via Svelte `use:action` or `bind:this`. Loses Svelte reactivity for grid data (no `$state()` binding). Must manually call AdaptableTools/AG Grid APIs for all updates.

**Verdict:** Technically feasible but painful. Writing a vanilla JS app inside a Svelte shell.

### Option 2: Post-Attach Hack

Let `ag-grid-svelte5` create the grid, then attach AdaptableTools after.

**Verdict:** Undocumented, fragile, likely breaks on version updates.

### Conclusion

Neither option is viable for a Svelte-first codebase. The integration cost negates AdaptableTools' "saves months" value proposition.

## AG Grid Enterprise vs Community

Features relevant to Perspectize:

| Feature | License | Current Need |
|---------|---------|-------------|
| Sorting, filtering, pagination | Community | Already using |
| Custom cell renderers | Community | Already using |
| Themes, resize, reorder | Community | Already using |
| Floating filters | Community | Already using |
| **Tool Panels** (column/filter sidebar) | Enterprise | Would enable end-user column customization |
| **Server-Side Row Model** | Enterprise | Would replace manual cursor pagination |
| Row Grouping | Enterprise | Maybe — group by channel, date |
| Master/Detail | Enterprise | Maybe — expand video to see perspectives |
| Excel Export | Enterprise | Low priority |
| Range Selection | Enterprise | Nice-to-have |
| Integrated Charts | Enterprise | Not core |

**Assessment:** AG Grid Enterprise ($999/dev/year) is not justified at current stage. Manual pagination works, and column chooser can be built with Community column API.

## Alternatives That Support Svelte

| Tool | Svelte Support | Notes | Pricing |
|------|---------------|-------|---------|
| TanStack Table | First-class Svelte 5 bindings | Headless — maximum flexibility, you build UI | Free (MIT) |
| SVAR DataGrid | Native Svelte widget | Less powerful than AG Grid, Svelte-native | Free Community + paid Pro |
| Bryntum Grid | Svelte wrapper | Enterprise-grade, closest AG Grid competitor | ~$400-$2K/dev/year |
| AG Grid + DIY | Via `ag-grid-svelte5` | Build only the features you need | Free (Community) |

## Recommendation

**Stay with AG Grid Community + `ag-grid-svelte5`. Build features incrementally.**

Rationale:
1. No Svelte support for AdaptableTools is a dealbreaker
2. Codebase is early-stage — unknown which power-user features users will actually need
3. Current ActivityTable already has server-side pagination, sorting, filtering, custom renderers
4. $6K-$21K+/year cost is not justified for a tool fighting framework integration

### What to Build Instead (When Needed)

- **Saved grid views** — localStorage + backend model
- **Column chooser** — AG Grid's built-in column API (`setColumnVisible()`)
- **Conditional formatting** — `cellStyle` callbacks (AG Grid Community)
- **Calculated columns** — `valueGetter` functions (already in use)

### When to Revisit

- Migrating to React or Angular
- 10+ complex grid views in the application
- Enterprise client requirements for financial-grade grid features
- AdaptableTools releases a Svelte wrapper

## Sources

- [AdaptableTools](https://www.adaptabletools.com/)
- [AdaptableTools Licensing](https://docs.adaptabletools.com/guide/licensing)
- [AdaptableTools Trial](https://www.adaptabletools.com/post/trial-licence)
- [AG Grid Community vs Enterprise](https://www.ag-grid.com/javascript-data-grid/community-vs-enterprise/)
- [AG Grid Tools & Extensions](https://www.ag-grid.com/community/tools-extensions/)
- [AG Grid Pricing](https://www.ag-grid.com/license-pricing/)
- [AdaptableTools npm](https://www.npmjs.com/package/@adaptabletools/adaptable)
