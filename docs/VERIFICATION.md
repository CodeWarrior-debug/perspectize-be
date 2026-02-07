# Verification & Evidence Capture

## GSD Plan Verification

For each plan's `must_haves`:

| Check | Command |
|-------|---------|
| `truths` | Run actual command, verify output |
| `artifacts.path` | `test -f {path} && echo "exists"` |
| `artifacts.contains` | `grep -q "{pattern}" {path}` |
| `artifacts.min_lines` | `wc -l < {path}` >= N |
| `key_links.pattern` | `grep -q "{pattern}" {from}` |

## Evidence Capture

Before creating PR:
- Screenshot at mobile (375px), tablet (768px), desktop (1024px+)
- Console output showing no errors
- Verification commands output
