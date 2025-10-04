# Claude Code Session Instructions

## Initialization Protocol

At the beginning of each session, perform the following:

1. **Scan the instruction set** located in [perspectize-be/.cursor/rules/](perspectize-be/.cursor/rules/)
   - Read [1-planning-workflow.mdc](perspectize-be/.cursor/rules/1-planning-workflow.mdc)
   - Read [2-pull-request-writer.mdc](perspectize-be/.cursor/rules/2-pull-request-writer.mdc)
   - Note: These rules may reference other files and directories in the project

2. **Confirm context loaded** by responding with:
   ```
   Context scanned - ready to begin.
   ```

## Core Instruction Set

Follow the guidelines defined in the `.cursor/rules/` folder:

### Planning Workflow ([1-planning-workflow.mdc](perspectize-be/.cursor/rules/1-planning-workflow.mdc))
- Use for large code changes requiring many lines of code or long-running work
- Create change plan documents in `.cursor/ai-change-plans/` (gitignored, local only)
- Name files using format: `{TICKET-ID}-{brief-description}.md`
- Structure plans with numbered sections and checkbox tasks
- Maintain no more than 10 workflow plan files
- Perform weekly cleanup of completed plans (1+ month old with all tasks checked)

### Pull Request Descriptions ([2-pull-request-writer.mdc](perspectize-be/.cursor/rules/2-pull-request-writer.mdc))
- Use standardized template for all pull requests
- Required elements: JIRA link, demo placeholder, summary, problem, solution, technical changes, testing, impact
- For large changes (>10 files), use grouping strategies: by category, key files + summary, or directory-based
- Ask user for output preference: new file, clipboard, other location, or append to workflow file

## Key Principles

- Extract ticket IDs from branch names (e.g., `feature/EX-737-add-timeline` ’ `EX-737`)
- Keep workflow files local (gitignored) for personal development workflow
- Update plans throughout development process
- Include proper JIRA links: `https://iralogix.atlassian.net/browse/{TICKET-NUMBER}`
- Use visual hierarchy in PR descriptions (emojis, bullet points, checkboxes)
