---
description: Systematically debug an issue with root cause analysis
agent: debugger
---

Debug the following issue:

$ARGUMENTS

Follow the debugging process:

1. **Gather context** — Read relevant files, check recent changes (`git log -5`), look at error messages
2. **Form hypotheses** — List 2-3 possible causes ranked by likelihood
3. **Investigate each** — Trace the code path, check data flow, look for common pitfalls
4. **Identify root cause** — Pin down the exact location and condition
5. **Apply fix** — Make the minimal correct change
6. **Verify** — Run relevant tests or explain how to verify the fix

Be systematic. Show your reasoning at each step.
