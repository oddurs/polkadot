---
description: Run tests, analyze failures, and suggest fixes
agent: debugger
---

Run the test suite for this project. Follow these steps:

1. **Detect the test runner** — Look at package.json scripts, or check for vitest/jest/pytest/cargo test config
2. **Run the tests** — Execute the appropriate test command
3. **Analyze results**:
   - If all pass: Report summary (pass count, coverage if available)
   - If failures: For each failing test, explain the root cause and suggest a fix
4. **Focus on**: $ARGUMENTS

If $ARGUMENTS specifies a file or pattern, only run tests matching that scope.
