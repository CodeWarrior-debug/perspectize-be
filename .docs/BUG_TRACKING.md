# Bug Tracking

## Overview

Known bugs and technical concerns are tracked **privately** — they are never committed to the public repository. This prevents exposing security vulnerabilities, internal quality metrics, and attack surface information before issues are resolved.

## How It Works

### Private Files (gitignored)

| File | Purpose |
|------|---------|
| `KNOWN_BUGS.md` (root) | Legacy bug inventory — migrated to planning phase |
| `.planning/codebase/CONCERNS.md` | Technical debt catalog — referenced by phase plans |
| `.planning/phases/bugs/` | **Persistent GSD phase** for bug fixes — backlog, plans, closed log |

These files exist only on developer machines. They are listed in `.gitignore` and are never pushed to the remote.

### Persistent Bugs Phase

Unlike numbered phases (1, 2, 3...) that complete and close, `.planning/phases/bugs/` is a **persistent phase** that remains open indefinitely:

```
.planning/phases/bugs/
├── README.md        # Phase overview and triage process
├── BACKLOG.md       # Full bug inventory with severity levels
├── BUG-{id}-PLAN.md # Individual fix plans (standard GSD format)
└── CLOSED.md        # Completed fixes with evidence
```

### Workflow

1. **Discover** a bug during development, review, or testing
2. **Log** it in `.planning/phases/bugs/BACKLOG.md` with severity and location
3. **Plan** the fix using a standard GSD plan file (`BUG-{id}-PLAN.md`)
4. **Fix** on a `bugfix/BUGS-{id}-description` branch
5. **Close** by moving the entry to `CLOSED.md` with PR reference

### Integration with Numbered Phases

Some numbered phases (e.g., Phase 6: Error Handling, Phase 9: Security Hardening) fix bugs as part of a larger initiative. When this happens:
- The numbered phase plan references the bug ID
- On completion, the bug is closed in both the phase and the bugs backlog
- The bugs backlog remains the single source of truth

## For New Developers

After cloning the repo, you won't have the private bug files. Ask a team member for the current `KNOWN_BUGS.md` and `.planning/phases/bugs/` contents. These are shared out-of-band (not via git).

## Rationale

Public bug trackers are standard for open-source projects, but Perspectize's bug inventory includes detailed security vulnerability descriptions with exact file paths and line numbers. Publishing this information before fixes are in place would create unnecessary risk. Once the project matures and critical issues are resolved, this policy may be revisited.
