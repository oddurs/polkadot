# Instructions

Managed by polkadot. opencode loads this in every session.

The workflow below is the same contract `~/.claude/CLAUDE.md` gives Claude Code.
Two agents working in the same repositories under the same git hooks need the
same rules, or the one with the looser rules is the one that breaks the build.

## Attribution

**Never name an AI or an assistant in anything that reaches a repository, a git
object, or a hosting platform.** No co-author trailers, no "generated with"
footers, no robot emoji, no "as an AI" in comments or docs. The work is
published under my name. If a tool inserts attribution by default, strip it.

## Working

- Act. If the request is clear, do it. Don't write a plan and stop.
- Ask once, batched, and only for what changes the output.
- Verify before claiming. Run the check; show what it printed.
- Say what you skipped and why. Never report a step you didn't watch pass.
- Disagree once, with the reason. Then do what I asked.

## Git

- **Never commit to the default branch.** It advances only through a merged
  pull request.
- One unit of work, one worktree, one branch, one pull request. Worktrees live
  in `../.worktrees/<repo>/<branch>/`.
- Branch `<type>/<slug>`, type ∈ feat fix chore docs perf refactor test.
- Conventional Commits. Imperative, subject ≤ 72 characters, no trailing
  period. The body explains why; the diff already says what.
- Green before it becomes a PR. The hooks enforce it. Never `--no-verify`,
  never `|| true` to make something pass.

`scripts/task` is the seam every repo here answers: `fmt fmt:check lint test
build check`. Use it rather than guessing at the toolchain. `scripts/agent`
drives the branch and worktree workflow where a repo has it.

## Code

- Make it work, make it right, make it small. In that order.
- Comments say why. The code already says what.
- No dead scaffolding: no TODO stubs, no commented-out code, no abstraction
  for a second caller that doesn't exist.
- A dependency earns its place or doesn't get added.
- Handle an error where you can act on it; otherwise propagate.
- A bug fix arrives with the test that would have caught it.

## TypeScript

- `strict: true`. No `any` — use `unknown` and narrow. No `@ts-ignore` without
  a comment saying what is ignored and why.
- Parse input at the boundary into a type, then trust the type inward.
- `type` for unions and shapes; `interface` only when merging is the point.
- Explicit return types on anything exported.
- ESM, named exports. No barrel files re-exporting a directory.

## React and Next.js

- Server Components by default. `"use client"` only on the leaves that need it.
- Colocate state with what uses it; extract a hook when it is reused, not
  before.
- Server Components fetch directly or through server actions; client
  components use TanStack Query.
- Every fetch has a loading state and an error state. Both are visible in the
  UI, not just handled.
- Tailwind utilities, variants through `cva`/`cn`, mobile-first, theming via
  CSS variables.

## Data and security

- Migrations for every schema change. Parameterised queries, never string
  interpolation. Index what you query often. RLS on in Supabase.
- Validate input on the server, whatever the client already did.
- Secrets come from the environment. Never commit a key — and if one is
  already in the diff, say so rather than committing around it.

## Tests

- Test behaviour, not implementation. Mock at the boundary you own.
- A test with no assertion is not a test.
- No sleeps and no retries. A flaky test is a failing test.
