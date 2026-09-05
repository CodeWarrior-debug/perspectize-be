#!/bin/bash
# PreToolUse hook (Bash) — non-blocking reminder to run Prettier before
# committing frontend files. Converted from hookify rule
# "prettier-precommit" (action: warn).

input=$(cat)
command=$(echo "$input" | jq -r '.tool_input.command // empty')

# Anchor on command position so this doesn't fire on `git commit` merely
# mentioned inside a commit message or other quoted text on the line.
if ! echo "$command" | grep -qE '(^|[;&|]|`|\$\()\s*git\s+commit\b'; then
  exit 0
fi

reason='Before this commit, check if any staged files are under frontend/. If so, run:
cd frontend && pnpm exec prettier --write $(git diff --cached --name-only --diff-filter=ACMR -- "frontend/" | sed "s|^frontend/||")
then re-stage the formatted files with `git add` before committing.'

jq -n --arg reason "$reason" '{
  hookSpecificOutput: {
    hookEventName: "PreToolUse",
    permissionDecision: "allow",
    permissionDecisionReason: $reason
  }
}'
exit 0
