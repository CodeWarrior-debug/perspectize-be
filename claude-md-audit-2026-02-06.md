# CLAUDE.md Quality Report — 2026-02-06

## Summary
- Files found: 3
- Files needing update: 0 (all recently restructured)
- Reference docs linked: 4 (`docs/` files)

---

## File-by-File Assessment

### 1. `./CLAUDE.md` (Project Root)
**Score: 82/100 (Grade: B)**

| Criterion | Score | Notes |
|-----------|-------|-------|
| Commands/workflows | 18/20 | gh CLI commands well-documented with copy-paste examples. GitHub Projects v2 correctly delegated to reference doc. |
| Architecture clarity | 15/20 | Good monorepo overview with stack pointers. Architecture details correctly delegated to package-level files. Missing: no mention of how root/package CLAUDE.md files relate to each other for new contributors. |
| Non-obvious patterns | 13/15 | Projects Classic deprecation workaround, token scope refresh, GSD plan verification checks — all good gotchas. |
| Conciseness | 13/15 | Tight. The qmd section is the densest area (two tables + workflow + GSD agents) but user chose to keep it. |
| Currency | 12/15 | References to `.planning/` and GSD workflow assume those exist. Branch naming uses `INI` prefix which may become stale as project phases change. |
| Actionability | 11/15 | Most commands are copy-pasteable. Agent delegation table is helpful but `{owner}/{repo}` placeholders in gh commands require substitution. |

**Strengths:**
- Clean separation of shared vs. stack-specific concerns
- All GitHub workarounds documented (Projects Classic, gh pr edit)
- Good progressive disclosure with `docs/` links

**Flagged issues:**
- `{owner}/{repo}` placeholders in gh commands — could be replaced with actual values for this project
- `INI` initiative prefix will become stale when project moves past initialization phase

---

### 2. `./perspectize-go/CLAUDE.md` (Backend)
**Score: 88/100 (Grade: B)**

| Criterion | Score | Notes |
|-----------|-------|-------|
| Commands/workflows | 20/20 | Complete: setup, daily dev, migrations, docker, graphql-gen all present and copy-pasteable. |
| Architecture clarity | 18/20 | Compact directory tree with link to full version. Hex arch guidelines give clear implementation order. Domain guide correctly delegated. |
| Non-obvious patterns | 14/15 | Excellent gotchas section: gqlgen defaults, JSON scalar, cursor pagination, enum binding. All are real foot-guns. |
| Conciseness | 13/15 | Well-trimmed. Error handling and DB queries delegated to docs. Enum section is the longest remaining block but user chose to keep it (it's marked REQUIRED). |
| Currency | 12/15 | `go-playground/validator` listed in stack but no usage pattern shown. `github.com/yourorg/perspectize-go` placeholder in gqlgen.yml example should use actual module path. |
| Actionability | 11/15 | "Adding a New Feature" checklist is excellent — 8 concrete steps with file paths. Self-verification curl is immediately runnable. |

**Strengths:**
- Best-in-class "Adding a New Feature" checklist
- Gotchas section captures real gqlgen/PostgreSQL foot-guns
- Development commands are complete and minimal

**Flagged issues:**
- `github.com/yourorg/perspectize-go` in gqlgen.yml example should use the actual Go module path
- Sevalla production note may be outdated if hosting has changed
- `clearConfigEnvVars` reference in Testing section — verify this helper still exists at that path

---

### 3. `./perspectize-fe/CLAUDE.md` (Frontend)
**Score: 72/100 (Grade: B)**

| Criterion | Score | Notes |
|-----------|-------|-------|
| Commands/workflows | 12/20 | Basic pnpm commands present. Missing: lint, format, type-check, storybook (if applicable). |
| Architecture clarity | 8/20 | No directory structure, no mention of routing patterns, component organization, or state management approach (TanStack Query). |
| Non-obvious patterns | 14/15 | AG Grid gotcha is critical and well-documented — this alone justifies the file. |
| Conciseness | 15/15 | Extremely lean. Every line earns its place. |
| Currency | 12/15 | AG Grid v32.2.x pinned versions — verify these are still current. |
| Actionability | 11/15 | Chrome DevTools MCP table is immediately actionable for verification workflows. |

**Strengths:**
- AG Grid section prevents a costly version conflict mistake
- Chrome DevTools verification table is practical and actionable

**Flagged issues:**
- Missing architecture section (SvelteKit routes, component structure, lib/ organization)
- No mention of TanStack Query patterns despite it being in the tech stack
- No lint/format/type-check commands
- No testing patterns (Vitest? Playwright?)

---

## Extended Analysis

### Instruction Count (per session)

| Session | Root | Package | Combined | Limit | Status |
|---------|------|---------|----------|-------|--------|
| root + backend | 87 | 110 | **197** | 200 | WARNING (3 remaining) |
| root + frontend | 87 | 30 | **117** | 200 | PASS (83 remaining) |

### Context Budget (per session)

| Session | Tokens | % of ~30k Budget | Status |
|---------|--------|-------------------|--------|
| root + backend | 3,048 | 10.2% | Healthy |
| root + frontend | 1,953 | 6.5% | Healthy |
| All files combined | 3,543 | 11.8% | Healthy |

**Per-file breakdown:**

| File | Lines | Chars | Tokens | Instructions | Category |
|------|-------|-------|--------|-------------|----------|
| `CLAUDE.md` | 162 | 5,835 | 1,458 | 87 | root |
| `perspectize-go/CLAUDE.md` | 189 | 6,362 | 1,590 | 110 | backend |
| `perspectize-fe/CLAUDE.md` | 52 | 1,982 | 495 | 30 | frontend |

### Performance Optimizations

| # | Severity | Finding | Detail |
|---|----------|---------|--------|
| 1 | **medium** | Backend session near limit | 197/200 instructions. Any new backend instructions should replace existing ones or go to `docs/`. |
| 2 | **medium** | Backend file density | 110 instructions in one file. Functional but tight. |
| 3 | **info** | Frontend has room | 83 instruction slots available for architecture docs, testing patterns, state management. |
| 4 | **info** | Context budget healthy | Both sessions under 11% — plenty of room for conversation context. |

### Manual Review Findings

| Check | Result |
|-------|--------|
| Instruction priority | Critical rules (dependency rule, hex arch, enum handling) are near top of backend file. Good. |
| Stale references | All 12 linked file paths verified as existing. `github.com/yourorg` placeholder needs updating. |
| Duplicate instructions | No duplicates found across root/backend/frontend files. Clean split. |
| Generic advice | None found — all instructions are project-specific. |
| Version-pinned info | AG Grid v32.2.x pins and Go 1.25+ should be verified periodically. |

---

## Recommendations

| Priority | Change | Impact | File |
|----------|--------|--------|------|
| P1 | Replace `github.com/yourorg/perspectize-go` with actual Go module path in enum example | Prevents copy-paste errors | `perspectize-go/CLAUDE.md:173` |
| P1 | Replace `{owner}/{repo}` in gh commands with `jamesjordan/perspectize-be` (or actual values) | Commands become directly copy-pasteable | `CLAUDE.md:27-38` |
| P2 | Add architecture section to frontend CLAUDE.md (routes, components, lib/) | Helps frontend agents navigate codebase | `perspectize-fe/CLAUDE.md` |
| P2 | Add lint/format/type-check commands to frontend | Completeness for frontend dev workflow | `perspectize-fe/CLAUDE.md` |
| P3 | Add TanStack Query patterns to frontend when established | Prevents re-discovery each session | `perspectize-fe/CLAUDE.md` |
| P3 | Verify AG Grid version pins are still current | Prevents stale version conflicts | `perspectize-fe/CLAUDE.md:20` |
