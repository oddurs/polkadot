---
description: Find and summarize all TODOs, FIXMEs, and HACKs in the codebase
agent: build
---

Scan the codebase for technical debt markers:

1. Search for: `TODO`, `FIXME`, `HACK`, `XXX`, `WORKAROUND`, `TEMP`, `DEPRECATED`
2. Group by severity:
   - **FIXME/HACK** — Issues that need attention soon
   - **TODO** — Planned improvements
   - **DEPRECATED/TEMP** — Code that should be removed
3. For each item, show: file:line, the comment, and brief context

Output as a prioritized table:
| Priority | File | Line | Type | Description |
|----------|------|------|------|-------------|

Focus on: $ARGUMENTS (if specified, limit scan to that directory/file pattern)
