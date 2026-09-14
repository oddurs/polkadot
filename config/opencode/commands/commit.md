---
description: Stage changes and generate a conventional commit message
agent: build
---

Help me commit the current changes:

1. Run `git status` and `git diff --staged` (and `git diff` if nothing staged)
2. Analyze ALL changes — understand what was done and why
3. Generate a commit message following Conventional Commits:
   - Format: `type(scope): description`
   - Types: feat, fix, refactor, docs, test, chore, perf, style, ci, build
   - Scope: the module/component affected (optional)
   - Description: imperative mood, lowercase, no period, under 72 chars
   - Body: explain "why" not "what" if the change isn't obvious
4. Stage relevant files (avoid .env, credentials, large binaries)
5. Show me the proposed commit message and wait for approval before committing

Additional context: $ARGUMENTS
