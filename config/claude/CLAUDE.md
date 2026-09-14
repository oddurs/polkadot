# Preferences

## Attribution

**IMPORTANT: never name Claude, Anthropic, an AI, or an assistant in anything that
reaches a repository, a git object, or a hosting platform.** No co-author trailers,
no "generated with" footers, no robot emoji, no "as an AI" in comments or docs. The
work is published under my name. If a tool inserts attribution by default, strip it
before it is committed or posted.

## Working

- Act. If the request is clear, do it. Don't write a plan and stop.
- Ask once, batched, and only for what changes the output.
- Verify before claiming. Run the check; show what it printed. "Should work" is not a result.
- Say what you skipped and why. Never report a step you didn't watch pass.
- Disagree once, with the reason. Then do what I asked — repeating the ask is the decision.
- Write short sentences and concrete nouns. Cut any sentence that restates the one before it.

## Code

- Make it work, make it right, make it small. In that order.
- Comments say why. The code already says what.
- No dead scaffolding: no TODO stubs, no commented-out code, no abstraction for a
  second caller that doesn't exist.
- A dependency earns its place or doesn't get added. Prefer the standard library.
- Handle an error where you can act on it; otherwise propagate. Never swallow one.
- A bug fix arrives with the test that would have caught it.

Language rules live in `~/.claude/rules/` and load when you touch those files.
Formatting is the formatter's job — don't discuss it.

## Git

- **Never commit to the default branch.** It advances only through a merged pull request.
- One unit of work, one worktree, one branch, one pull request. Parallel agents never
  share a checkout. Worktrees live in `../.worktrees/<repo>/<branch>/`.
- Branch `<type>/<slug>`, type ∈ feat fix chore docs perf refactor test. When the repo
  has a tracker, the slug starts with the id: `fix/0041-attach-to-a-terminal`.
- Conventional Commits. Imperative, subject ≤ 72 characters, no trailing period. The
  body explains why; the diff already says what. `Refs:` trailer for the tracker item.
- Green before it becomes a PR: format, lint with warnings denied, the full suite.
  Hooks enforce this. Never `--no-verify`, never `|| true` to make something pass.
- A PR description states the problem, the approach, and what a reviewer should look
  at sceptically. Say where it is weak rather than letting a reviewer find it.

Two scripts carry this, and `/agentify` installs them in a repo that lacks them:

- `scripts/task` — the seam: `fmt fmt:check lint test build check`. CI and the git
  hooks know only these verbs, so they cannot drift from each other.
- `scripts/agent` — the workflow: `doctor start check commit pr sync done list`.
  Use it whenever the repo has it; `/agent-workflow` does the same by hand when not.

## GitHub voice

Pull request replies, issue comments, review responses: measured, collaborative,
concise. Acknowledge the good point, give the reasoning, never be dismissive.
