---
paths:
  - "**/*.sh"
  - "**/scripts/*"
  - "**/.githooks/*"
---

# Shell

- POSIX `sh`, not bash, unless a bash feature is load-bearing — these run in CI images
  and git hooks where bash is not guaranteed.
- `set -eu` on the first line after the shebang. Quote every expansion.
- No dependency on a framework: no husky, no lefthook, no pre-commit. A hook is an
  executable file in `.githooks/`.
- Every script prints what it is about to do, is safe to run twice, and exits non-zero
  with a message the reader can act on.
- Never `|| true` to hide a failure, never `--no-verify`.
- Resolve the default branch from the remote; don't hardcode `main`.
