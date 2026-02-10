---
name: block-rm-system-dirs
enabled: true
event: bash
action: block
pattern: rm\s+.*(/etc|/usr|/var|/bin|/sbin|/lib|/boot|/root|/home(?!/jamesjordan/GitHub)|/System|/Library|/Applications|\s+/$|\s+/\s|\s+~(?!/GitHub))
---

🚫 **BLOCKED: Dangerous rm command targeting system directory**

You attempted to remove files from a protected system location.

**Protected paths include:**
- `/etc`, `/usr`, `/var`, `/bin`, `/sbin`, `/lib`, `/boot`
- `/root`, `/home` (except project paths)
- `/System`, `/Library`, `/Applications` (macOS)
- `/` (root filesystem)
- `~` (home directory, except ~/GitHub)

**If you need to remove files:**
1. Use explicit paths within the project: `perspectize-fe/` or `perspectize-go/`
2. Ask the user for confirmation first
3. Never use wildcards with system paths
