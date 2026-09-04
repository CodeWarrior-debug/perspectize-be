# Verification & Evidence Capture

Before marking any work complete, run interactive verification.

## 1. Start Services

The database is hosted on Sevalla (cloud PostgreSQL) — no Docker or local database setup needed.

**Check if already running:**
```bash
lsof -i :8080  # Backend
lsof -i :5173  # Frontend
```

**Start if not running:**
```bash
# Terminal 1: Backend (port 8080)
cd backend
make dev    # hot reload with air
# or: make run  # standard mode

# Terminal 2: Frontend (port 5173)
cd frontend
pnpm run dev
```

## 2. Verify Backend

```bash
curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -d '{"query": "{ __typename }"}'
# Expect: {"data":{"__typename":"Query"}}
```

Also test any frontend GraphQL queries (`src/lib/queries/*.ts`) against the live backend to catch schema drift.

## 3. Verify Frontend (Chrome DevTools MCP)

| Step | MCP Tool | Purpose |
|------|----------|---------|
| Navigate | `mcp__chrome-devtools__navigate_page` | Load frontend URL |
| Screenshot | `mcp__chrome-devtools__take_screenshot` | Visual verification |
| Snapshot | `mcp__chrome-devtools__take_snapshot` | DOM/component structure |
| Resize | `mcp__chrome-devtools__resize_page` | Responsive check (375px, 768px, 1024px) |
| Console | `mcp__chrome-devtools__list_console_messages` | Check for JS errors |
| Interact | `mcp__chrome-devtools__click` | Test buttons, toasts, navigation |

## 4. GSD Plan Verification

For each plan's `must_haves`:

| Check | Command |
|-------|---------|
| `truths` | Run actual command, verify output |
| `artifacts.path` | `test -f {path} && echo "exists"` |
| `artifacts.contains` | `grep -q "{pattern}" {path}` |
| `artifacts.min_lines` | `wc -l < {path}` >= N |
| `key_links.pattern` | `grep -q "{pattern}" {from}` |

## 5. Evidence Capture

Save screenshots to `/Users/jamesjordan/Downloads/screenshots/` with naming convention:
- **Prefix:** `sv-` (Self-Verify) — supersedes the old `ccsv-` prefix; `ccsv-` may still appear in older screenshots but new captures use `sv-`
- **Format:** `sv-{plan}-{description}-{width}.png`
- **Example:** `sv-01-02-mobile-375px.png`, `sv-01-04-ag-grid-desktop-1280px.png`
- Use `filePath` parameter on `take_screenshot` to save directly
- Take full-page screenshots (`fullPage: true`) at mobile (375px), tablet (768px), desktop (1280px)

Before creating PR:
- Screenshots at mobile (375px), tablet (768px), desktop (1280px)
- Console output showing no errors
- Verification commands output
- Upload the `sv-*` screenshots and link them in the PR — see [PR Screenshots](PR_SCREENSHOTS.md)
