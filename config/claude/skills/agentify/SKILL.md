---
name: agentify
description: Make an existing repository agent-native — a lean CLAUDE.md, path-scoped .claude/rules, repo skills, the scripts/task seam, scripts/agent worktree workflow, git hooks, and CI. Use when a repo has no CLAUDE.md or a bloated one, lacks the branch/worktree/PR workflow, or when asked to set up agentic structure, /init properly, or "make this repo work with Claude".
argument-hint: [path to repo, defaults to cwd]
---

# Make this repository agent-native

The repository already exists and has code in it. Bring it up to the standard a new one
would be born with, without breaking anything that works today. Execute — do not produce
a plan and stop.

Work in the current directory unless `$ARGUMENTS` names another.

## 0. Audit before you write

Read before you generate. Establish, from the repo itself:

- Stack, toolchain versions, and the real commands to format, lint, test, and build.
  Find them in `package.json` scripts, `Makefile`, `justfile`, CI workflows, and the
  README — not from your assumptions about the ecosystem.
- What already exists: `CLAUDE.md`, `AGENTS.md`, `.cursorrules`, `.github/workflows/`,
  hooks (`core.hooksPath`, `.husky/`, `.pre-commit-config.yaml`), `scripts/`.
- The default branch, whether there is a remote, and whether it is protected.
- Whether the tree is clean. If it is not, stop and say so — this skill rewrites
  tracked files and the diff must be readable.

Report what you found in a few lines, then proceed. Do not ask permission for the
steps below; they are the point of the skill.

## 1. `CLAUDE.md` — the part most repos get wrong

A CLAUDE.md is loaded into context at the start of **every** session, which is what makes
it valuable and what makes bloat expensive: when the file is long, the rules that matter
get lost among the ones that don't, and adherence to all of them drops.

**Under 200 lines. Aim for well under.** For each line, ask: *would removing this cause a
mistake?* If not, cut it.

| Belongs in CLAUDE.md | Belongs somewhere else |
|---|---|
| Commands that can't be guessed | Anything derivable by reading the code → cut it |
| Conventions that differ from the ecosystem default | Standard conventions → cut it |
| Repo etiquette: branches, commits, PRs | Formatting rules → the formatter's config |
| Architectural decisions and the reason for them | Language style → `.claude/rules/*.md`, path-scoped |
| Environment quirks: required vars, services, ports | A multi-step procedure → `.claude/skills/` |
| Non-obvious gotchas and traps | File-by-file tours, API docs → link instead |
| A "never do X" rule with teeth | Anything that must happen every time → a hook |

Write it as instructions to an agent, not prose about the project. Be concrete enough to
verify: "run `pnpm test -- <file>` for a single test" beats "test your changes". Use at
most one `IMPORTANT` in the file — emphasise everything and you emphasise nothing.

If a `CLAUDE.md` already exists, edit it rather than replacing it: keep what is true,
cut what the codebase already says, and move what is path-specific or procedural into
the right home below. Say what you cut and why.

If the repo has an `AGENTS.md` for other tools, do not duplicate it. Make `CLAUDE.md`
import it and add only Claude-specific content:

```markdown
@AGENTS.md

## Claude Code
<what is specific to this tool>
```

Every repo's CLAUDE.md ends with the workflow contract, because that is the part an agent
must not improvise:

```markdown
## Workflow

- Never commit to <default>. It advances only through a merged pull request.
- One unit of work, one worktree, one branch, one PR: `scripts/agent start fix/the-thing`.
- `scripts/task check` must pass before a PR. The hooks enforce it; never `--no-verify`.
- Conventional Commits, subject ≤ 72 chars, imperative, no trailing period.
- No AI or assistant attribution in commits, PRs, comments, or docs.
```

## 2. `.claude/rules/` — style without the context cost

Anything that only matters while touching a particular kind of file goes here, scoped
with `paths:` frontmatter so it loads on demand instead of every session:

```markdown
---
paths:
  - "src/api/**/*.ts"
---

# API handlers
- Validate input at the boundary with the shared schema; trust the type inward.
- Errors use the standard envelope in `src/api/errors.ts:1`.
```

Create one file per topic — `testing.md`, `api.md`, the languages actually present.
Point at files with `path:line` references rather than pasting code that will go stale.
Rules without `paths:` load every session, so give every rule a `paths:` unless it truly
applies to the whole repo.

## 3. `.claude/skills/` — the repeatable procedures

Anything with steps — release, migration, adding an endpoint, regenerating fixtures,
debugging the flaky integration suite — is a skill, not a CLAUDE.md section. Skills load
only when invoked or when the model judges them relevant, so they cost nothing the rest
of the time.

Write one for each procedure this repo has that a newcomer would get wrong. Give each a
`description` that says *when to use it*, since that is what selection is based on. Add
`disable-model-invocation: true` to any skill with side effects you want triggered by
hand.

## 4. `scripts/task` — the seam

One POSIX `sh` dispatcher, the only thing in the repo that knows what language this is:

```
fmt  fmt:check  lint  test  build  check
```

`check` runs `fmt:check`, `lint`, `test`, `build` in that order and prints a clear
pass/fail. Every target exists and exits 0 on the current tree; a genuinely inapplicable
target is an explicit no-op with a one-line comment saying why.

Wire it to the commands you found in §0 — do not invent new ones, and do not change how
the project builds as a side effect of this skill. If `package.json` scripts already
exist, `scripts/task` calls them.

From here on, CI, the hooks, and `scripts/agent` call only these six verbs. That is what
keeps a local hook and a CI job from ever drifting apart.

## 5. `scripts/agent` — the workflow

Same contract as a new repo: `doctor start check commit pr sync done list`, POSIX `sh`,
idempotent, prints what it is about to do, never `cd`s for the caller, resolves the
default branch from the remote instead of hardcoding it. Worktrees go in
`../.worktrees/<repo>/<branch>/`. Add `.worktrees/` to `.gitignore`.

Add `scripts/setup` — one command a contributor runs after cloning. It sets
`core.hooksPath` and checks the toolchain is present.

## 6. `.githooks/` — enforcement

Tracked in git, wired with `git config core.hooksPath .githooks`, plain executable `sh`
with no framework. If the repo uses husky, lefthook, or pre-commit, migrate the checks
into these and remove the framework — one fewer dependency, and it works for agents that
never ran `npm install`.

- `commit-msg` — Conventional Commit, subject ≤ 72 chars, no trailing period; reject any
  assistant attribution. Let `Merge`, `Revert`, `fixup!`, `squash!`, `amend!` through.
- `pre-commit` — `scripts/task fmt:check` and `lint`. Keep it under a second or two; if
  it can't be, move it to pre-push.
- `pre-push` — refuse a push to the default branch (allow only the very first push that
  creates it), then `scripts/task check`.

## 7. CI

One `.github/workflows/ci.yml` on `push` to the default branch and on `pull_request`.
`concurrency` with `cancel-in-progress`, top-level `permissions: contents: read`, caching
via the official setup action, third-party actions pinned to a SHA. A single `check` job
that runs `./scripts/task check`, and a `required` job that `needs` it — so branch
protection names one status check and never has to be edited when jobs are added.

If CI already exists, fold it into this shape rather than adding a second workflow.

## 8. The remote, if there is one

Squash-merge only, auto-delete head branches, branch protection requiring a PR and the
`required` check, conversation resolution, no force pushes or deletions. 0 approvals on a
solo repo so the owner isn't deadlocked, 1 otherwise. Anything the token can't do: finish
everything else and print the exact command the owner must run.

## 9. Prove it

Written files are not evidence. Watch each of these pass:

1. `scripts/setup` → hooks path set.
2. `scripts/task check` → green on the tree as it stands. If the repo is red today, say
   so plainly and fix only what your own changes broke.
3. `scripts/agent doctor` → clean.
4. `scripts/agent start chore/agentify` → worktree at the expected path; commit the
   changes there and open the PR through `scripts/agent pr`.
5. A direct push to the default branch is refused by the hook.
6. `/context` in a fresh session shows `CLAUDE.md` under **Memory files**, and the rules
   files do *not* appear until a matching file is read.

## 10. Report

- What you found, what you changed, what you removed.
- The CLAUDE.md line count before and after.
- What moved out of CLAUDE.md and where it went.
- Anything the owner must run by hand, with the command.

## Anti-patterns

- Running `/init` and shipping what it generated. It is a starting point, not the file.
- A CLAUDE.md that recaps the directory tree, lists dependencies, or explains the
  architecture an agent can read in thirty seconds.
- Style rules in prose that a linter could enforce in config.
- Rules without `paths:`, which load every session and undo the point of the split.
- Two sources of truth: a `Makefile` and a `scripts/task` that drift apart.
- Adding husky or lefthook because it is familiar.
- Rewriting the project's build as a side effect. `scripts/task` wraps what exists.
