---
name: block-sudo
enabled: true
event: bash
action: block
pattern: ^sudo\s+|;\s*sudo\s+|\|\s*sudo\s+|&&\s*sudo\s+
---

🚫 **BLOCKED: sudo command detected**

Running commands with `sudo` (superuser privileges) is not allowed.

**Why this is blocked:**
- Claude should never need root privileges for development tasks
- Accidental sudo commands could damage the system
- All project work should happen in user space

**If you genuinely need elevated privileges:**
1. Explain the situation to the user
2. Let the user run the command manually
3. Never attempt to bypass this restriction
