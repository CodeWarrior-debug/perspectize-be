#!/bin/bash
# PreToolUse hook (Bash) — non-blocking reminder to check test coverage
# before committing. Converted from hookify rule "require-tests"
# (action: warn).

input=$(cat)
command=$(echo "$input" | jq -r '.tool_input.command // empty')

# Anchor on command position so this doesn't fire on `git commit` merely
# mentioned inside a commit message or other quoted text on the line.
if ! echo "$command" | grep -qE '(^|[;&|]|`|\$\()\s*git\s+commit\b'; then
  exit 0
fi

reason='Check for missing unit tests before committing. Review the staged changes (git diff --cached). If any files under frontend/src/ or backend/internal/ contain new or modified functions, components, hooks, or business logic, verify corresponding test files exist:

Frontend:
- New .svelte components -> tests/components/{Component}.test.ts
- New/modified hooks in lib/queries/hooks/ -> tests/unit/hooks-{name}.test.ts
- New/modified utilities in lib/utils/ -> tests/unit/{name}.test.ts
- New/modified query definitions in lib/queries/ -> tests/unit/queries-{name}.test.ts

Backend:
- New/modified domain types in internal/core/domain/ -> {file}_test.go
- New/modified services in internal/core/services/ -> {service}_test.go
- New/modified handlers in internal/adapters/handler/ -> {handler}_test.go
- New/modified repositories in internal/adapters/repository/ -> {repo}_test.go
- New/modified middleware in internal/adapters/middleware/ -> {middleware}_test.go

Does NOT need tests: config files, type-only files, style-only changes, docs, test files themselves, pure reformatting, generated code (graph/generated/, graph/model/), migration SQL, GraphQL schema files.

If tests are missing, write them before committing. If the change genuinely does not need tests, proceed.'

jq -n --arg reason "$reason" '{
  hookSpecificOutput: {
    hookEventName: "PreToolUse",
    permissionDecision: "allow",
    permissionDecisionReason: $reason
  }
}'
exit 0
