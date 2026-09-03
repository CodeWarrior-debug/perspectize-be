# PR Screenshots (Release-Hosted)

Screenshots captured during self-verification (see [VERIFICATION.md](VERIFICATION.md), `sv-` prefix) live locally in `~/Downloads/screenshots/` — they aren't committed to the repo. To reference them in a PR, upload them as assets on a dedicated GitHub Release that acts as a permanent asset bucket, then link the resulting URLs in the PR body.

> **Chrome DevTools MCP note:** the MCP's `take_screenshot` sandbox only permits paths **inside the repo** (`~/Downloads/...` is rejected). When capturing via MCP, save to the gitignored `.screenshots/` dir at the repo root and `gh release upload` from there instead.

## One-time setup

Check whether the bucket release already exists before creating it:

```bash
gh release view screenshots
```

If it doesn't exist yet, create it once:

```bash
gh release create screenshots \
  --title "Verification Screenshots" \
  --notes "Asset bucket for PR self-verification screenshots. Not a versioned release — do not use for changelog/version tracking." \
  --prerelease
```

## Per-PR workflow

1. **Upload this PR's screenshots** (`--clobber` overwrites same-named assets safely on re-runs):
   ```bash
   gh release upload screenshots ~/Downloads/screenshots/sv-<plan>-*.png --clobber
   ```

2. **Get shareable URLs** for just the files you uploaded:
   ```bash
   gh release view screenshots --json assets --jq '.assets[] | select(.name | startswith("sv-<plan>-")) | .browser_download_url'
   ```

3. **Paste them into the PR body's Testing section** as markdown images:
   ```markdown
   ### Screenshots
   | Viewport | Screenshot |
   |---|---|
   | Mobile (375px) | ![mobile](https://github.com/CodeWarrior-debug/perspectize/releases/download/screenshots/sv-01-02-mobile-375px.png) |
   | Tablet (768px) | ![tablet](https://github.com/CodeWarrior-debug/perspectize/releases/download/screenshots/sv-01-03-tablet-768px.png) |
   | Desktop (1280px) | ![desktop](https://github.com/CodeWarrior-debug/perspectize/releases/download/screenshots/sv-01-04-desktop-1280px.png) |
   ```

## Why a release instead of committing the images

- Keeps screenshot binaries out of git history entirely (no repo bloat, no LFS needed)
- GitHub renders release-asset image URLs inline in PR markdown same as any other image
- One persistent `screenshots` release accumulates assets across every PR — no per-PR release/tag needed
- `--clobber` makes re-uploading after a fix idempotent — same filename just replaces the old asset and the PR's existing markdown link keeps working

## Required permission

`gh release create`/`gh release upload`/`gh release view` need the `Bash(gh release:*)` allow rule in `.claude/settings.json`.
