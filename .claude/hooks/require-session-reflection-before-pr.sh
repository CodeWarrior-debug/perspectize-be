#!/bin/bash
# PreToolUse hook (Bash) — blocks `gh pr create` until the
# /revise-claude-md command (claude-md-management plugin) has been run
# this session.
# Converted from hookify rule "require-session-reflection-before-pr"
# (action: block).
#
# Detection: hookify itself could not detect completion either — this rule
# always fires on the first `gh pr create` attempt. Once you've run the
# command, re-run `gh pr create` (or use `gh api ... /pulls` per CLAUDE.md,
# which this hook does not match) and it will proceed.

input=$(cat)
command=$(echo "$input" | jq -r '.tool_input.command // empty')

# Anchor on command position (start of string, or after a shell separator)
# so this doesn't fire on `gh pr create` merely mentioned inside a commit
# message, PR body, or other quoted text elsewhere in the command line.
if ! echo "$command" | grep -qE '(^|[;&|]|`|\$\()\s*gh\s+pr\s+create\b'; then
  exit 0
fi

reason='Session reflection required before PR creation.

Before creating a PR, you MUST run the /revise-claude-md command to capture session learnings:
1. Run the /revise-claude-md slash command (from the claude-md-management plugin)
2. Review what context was missing, permissions needed, gotchas encountered
3. Apply approved changes to CLAUDE.md
4. Include a "Session Learnings" section in the PR description summarizing what was updated
5. Then create the PR via `gh api repos/CodeWarrior-debug/perspectize/pulls -f ...` (per CLAUDE.md — gh pr create is blocked by this hook and also fails due to Projects Classic deprecation)'

jq -n --arg reason "$reason" '{
  hookSpecificOutput: {
    hookEventName: "PreToolUse",
    permissionDecision: "deny",
    permissionDecisionReason: $reason
  }
}'
exit 0
