#!/usr/bin/env bash
# SessionStart hook: on the first session after the 1st of each month, if at
# least 10 merges have landed on main since the last recorded monthly-routine
# run, inject context prompting Claude to run the monthly-maintenance skill.
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
ROUTINES_FILE="$REPO_ROOT/ROUTINES.md"

CURRENT_MONTH="$(date +%Y-%m)"

# Bail quietly if the tracking file is missing (nothing to check against).
if [ ! -f "$ROUTINES_FILE" ]; then
  exit 0
fi

# Already have a row for this month (completed or not) -> already handled,
# don't prompt again this month.
if grep -qE "^\| *${CURRENT_MONTH} *\|" "$ROUTINES_FILE"; then
  exit 0
fi

# Find the most recent recorded run's month-year (last data row of the table).
LAST_MONTH="$(grep -E '^\| *[0-9]{4}-[0-9]{2} *\|' "$ROUTINES_FILE" | tail -1 | awk -F'|' '{gsub(/^ +| +$/, "", $2); print $2}' || true)"

if [ -n "${LAST_MONTH:-}" ]; then
  CUTOFF="${LAST_MONTH}-01"
else
  # No prior run recorded yet — bootstrap with a 90-day lookback.
  CUTOFF="$(date -v-90d +%Y-%m-%d 2>/dev/null || date -d '90 days ago' +%Y-%m-%d)"
fi

cd "$REPO_ROOT"

MERGE_COUNT="$(git log --oneline --since="$CUTOFF" origin/main 2>/dev/null | wc -l | tr -d ' ')"

if [ -z "$MERGE_COUNT" ] || [ "$MERGE_COUNT" -lt 10 ]; then
  exit 0
fi

CONTEXT="Monthly maintenance routine is due: ${MERGE_COUNT} merges have landed on main since the last recorded run (cutoff ${CUTOFF}), and no routine has been logged for ${CURRENT_MONTH} yet in ROUTINES.md. Invoke the monthly-maintenance skill to run it (branch cleanup, dependabot/security PR merges, graphify update, gsd map-codebase refresh), then record the outcome as a new row in ROUTINES.md."

jq -n --arg ctx "$CONTEXT" '{hookSpecificOutput: {hookEventName: "SessionStart", additionalContext: $ctx}}'
