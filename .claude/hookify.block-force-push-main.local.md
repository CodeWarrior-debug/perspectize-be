---
name: block-force-push-main
enabled: true
event: bash
action: block
pattern: git\s+push\s+.*(-f|--force|--force-with-lease).*\s+(origin\s+)?(main|master)|\bgit\s+push\s+(-f|--force|--force-with-lease)\b
---

🚫 **BLOCKED: Force push to main/master**

Force pushing to `main` or `master` is extremely dangerous:
- Rewrites shared history
- Can cause other developers to lose work
- May break CI/CD pipelines
- Often irreversible without backups

**Safe alternatives:**
1. Create a PR and merge normally
2. Use `git revert` to undo commits instead of rewriting history
3. If you absolutely must force push, ask the user to do it manually

**This protection exists because mistakes happen, especially at 2am.**
