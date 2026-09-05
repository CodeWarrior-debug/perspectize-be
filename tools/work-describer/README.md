# Work Describer / Prompt Builder

A standalone, static HTML + TypeScript tool for building structured work-request
prompts to hand off to a Claude Code session (or any coding assistant).

**Fully deterministic** — the output is assembled from your form inputs by plain
string templates. No network calls, no AI/tokens are used to generate the text.

## Usage

Just open `index.html` in a browser — no server or build step is required to
use it, since the compiled `dist/main.js` is checked in.

Fill in:

- **Task name** / **Task type** (feature, bugfix, chore, docs, refactor, test, other)
- **Intent** — free text description; put narrative notes about links/attachments here
- **Links & attachments** — one URL or file reference per line
- **Repo** — free-text combo box; click "Set default" to remember it as the
  suggested value next time (stored in `localStorage`, most-recent-first)
- **Starting branch** — same combo-box pattern, defaults to `main`
- **Working location** — Main worktree / Available non-main worktree
- **PR target branch** — defaults to `main`
- **Testing approach** — "follow codebase instructions, then patterns, then
  judgment; stop and ask if really unsure" (default) or "stop and ask if
  really unsure" outright, plus optional free-text notes
- **Codebase instructions to ignore (this session only)** — one per line

The right-hand panel renders a live "rich preview" (basic markdown → HTML) plus
the raw Markdown source. Use **Copy to clipboard** to paste directly into a
Claude session, or **Download .md** to save the file.

Your last-used field values, and any repos/branches you mark as default, persist
in `localStorage` so the form comes back pre-filled next time.

## Development

Source lives in `src/main.ts`; the committed `dist/main.js` is its compiled
output. To rebuild after editing:

```bash
npm install
npm run build       # one-off build
npm run watch       # rebuild on save
```

No bundler or framework is used — `main.ts` compiles to a single ES module
loaded via `<script type="module" src="./dist/main.js">`.
