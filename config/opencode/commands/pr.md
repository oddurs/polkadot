---
description: Generate a pull request title and description from current changes
agent: build
---

Generate a pull request description for the current branch:

1. Run `git log main..HEAD --oneline` (or appropriate base branch) to see all commits
2. Run `git diff main..HEAD --stat` to see changed files
3. Read the key changed files to understand the full scope

Then generate:

**Title**: Short, descriptive (under 70 chars), imperative mood

**Body**:
```
## Summary
- Bullet points describing what changed and why (2-5 bullets)

## Changes
- Key files/components modified and what was done

## Testing
- How this was tested or should be tested
- [ ] Checklist of test scenarios

## Notes
- Any migration steps, breaking changes, or follow-up work needed
```

Additional context: $ARGUMENTS
