---
name: block-unsafe-rm-rf
enabled: true
event: bash
action: block
conditions:
  - field: command
    operator: regex_match
    pattern: rm\s+-rf
  - field: command
    operator: not_contains
    pattern: perspectize-fe
  - field: command
    operator: not_contains
    pattern: perspectize-go
  - field: command
    operator: not_contains
    pattern: node_modules
  - field: command
    operator: not_contains
    pattern: .svelte-kit
  - field: command
    operator: not_contains
    pattern: build/
  - field: command
    operator: not_contains
    pattern: dist/
---

🚫 **BLOCKED: rm -rf outside safe directories**

`rm -rf` is only allowed in these safe directories:
- `perspectize-fe/` (frontend project)
- `perspectize-go/` (backend project)
- `node_modules/` (dependencies)
- `.svelte-kit/` (SvelteKit cache)
- `build/` or `dist/` (build outputs)

**Your command was blocked because it doesn't target a safe path.**

If you need to remove something else, ask the user for explicit permission first.
