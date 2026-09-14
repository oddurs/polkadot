---
description: Refactor code for clarity, maintainability, or performance
agent: build
---

Refactor the specified code. Follow these principles:

1. **Read first** — Understand the full context before changing anything
2. **Preserve behavior** — Refactoring must not change functionality. Run tests before and after
3. **Small steps** — Make incremental changes, not a big bang rewrite
4. **Apply these patterns**:
   - Extract repeated logic into well-named functions
   - Replace complex conditionals with early returns or guard clauses
   - Simplify deeply nested code
   - Improve naming to reveal intent
   - Remove dead code and unused imports
   - Split large functions (>40 lines) into focused units
5. **Skip** — Don't add comments, types, or docs unless specifically asked

Focus on: $ARGUMENTS
