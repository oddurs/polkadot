---
name: bootstrap-repo
description: Bootstrap an empty directory into a fully operational public GitHub repo — scripts/task seam, scripts/agent worktree workflow, git hooks, CI, branch protection, and open-source furniture. Use when starting a new project from zero, or when asked to set up, scaffold, or initialise a repository.
argument-hint: [project name] [one-line description] [stack]
---

# Bootstrap a public open-source repository

You are setting up a brand-new project from zero to a fully operational public GitHub
repository. When you are done, a human or an agent must be able to clone it, run one
command, and be productive — with a workflow that makes the incorrect thing (committing
to `main`, pushing red code, two agents sharing a checkout) structurally impossible
rather than merely discouraged.

Work in the current directory unless told otherwise. Execute — do not produce a plan
and stop.

## 0. Resolve inputs

Infer everything you can from `$ARGUMENTS`, the directory name, and any files already
present. Then ask **once**, in a single batched question, only for what genuinely
changes the output and cannot be defaulted:

| Input | Default if unspecified |
|---|---|
| Repo name | current directory name, kebab-cased |
| Owner | the authenticated `gh` user |
| One-line description | derived from the name; ask if you'd be guessing |
| Language / stack | ask — it changes every generated file |
| License | `MIT` |
| Default branch | `main` |
| Visibility | `public` |

Do not ask about anything else. Everything below is decided for you.

## 1. Non-negotiables

These override any convention from a template, tool, or your own defaults.

1. **No AI or assistant attribution anywhere.** Not in commits, trailers, PR bodies,
   code comments, docs, changelogs, or release notes. No co-author trailers naming a
   model, no "generated with" footers, no robot emoji. The work is published under the
   owner's name. If a tool inserts attribution by default, strip it before it reaches a
   git object or the GitHub API.
2. **`main` only ever advances through a merged pull request.** Enforced by branch
   protection on the server and by a local hook, not by discipline.
3. **One unit of work → one worktree → one branch → one PR.** Parallel agents never
   share a checkout.
4. **A branch is green before it becomes a PR.** Format, lint (warnings denied), and
   the full test suite.
5. **Lean.** Every file you add must earn its place. No placeholder docs, no dead
   config, no dependency added "just in case", no workflow that never runs.

## 2. The stack seam

All automation talks to the project through exactly one interface, so CI, git hooks,
and the agent script never learn anything stack-specific:

```
scripts/task fmt        # format in place
scripts/task fmt:check  # verify formatting, non-zero on drift
scripts/task lint       # lint, warnings are errors
scripts/task test       # full test suite
scripts/task build      # build/compile; no-op is acceptable, must exit 0
scripts/task check      # fmt:check && lint && test && build
```

Write `scripts/task` as a POSIX `sh` dispatcher that shells out to the real toolchain
for the chosen stack (`cargo fmt`/`clippy`/`nextest`, `pnpm` scripts, `ruff`/`pytest`,
`go fmt`/`vet`/`test`, …). Every target must exist and exit 0 on the freshly scaffolded
project. If a target is genuinely not applicable, make it an explicit no-op with a
one-line comment — never leave it undefined.

This seam is the load-bearing idea. Get it right and the rest is wiring.

## 3. Scaffold the project

Create a minimal but *real* project for the chosen stack — a hello-world entry point
plus one passing test. Prefer the ecosystem's own non-interactive scaffolder
(`cargo init`, `pnpm create …@latest --yes`, `uv init`, `go mod init`) over
hand-rolling. Pin the toolchain version in the file the ecosystem expects
(`rust-toolchain.toml`, `.node-version`, `.python-version`, `go.mod`).

## 4. Open-source furniture

Create these, with real content — no `TODO` placeholders:

- `README.md` — what it is, why it exists, install, quickstart, a `Development`
  section pointing at `scripts/agent`, license badge, CI badge.
- `LICENSE` — the chosen license, correct year and copyright holder.
- `CONTRIBUTING.md` — the branch/worktree/PR workflow from §5, commit convention, how
  to run checks locally. Short and imperative.
- `CODE_OF_CONDUCT.md` — Contributor Covenant 2.1, with a real contact.
- `SECURITY.md` — supported versions, private reporting via GitHub Security Advisories,
  response expectations.
- `CHANGELOG.md` — Keep a Changelog format, `## [Unreleased]` only.
- `.gitignore` — stack-appropriate, plus `.worktrees/` and OS/editor cruft.
- `.gitattributes` — `* text=auto eol=lf` and any binary/linguist rules.
- `.editorconfig` — matching the formatter's settings exactly.
- `.github/PULL_REQUEST_TEMPLATE.md` — problem, approach, what a reviewer should look
  at sceptically, checklist tied to `scripts/task check`.
- `.github/ISSUE_TEMPLATE/bug_report.yml` and `feature_request.yml` — YAML forms, not
  markdown; plus `config.yml` disabling blank issues.
- `.github/CODEOWNERS` — `* @<owner>`.
- `.github/dependabot.yml` — `github-actions` weekly, plus the stack's ecosystem
  weekly, grouped minor+patch into a single PR.
- `CLAUDE.md` — the contract for agents working in this repo: the workflow, the
  `scripts/task` seam, the commit convention, and the attribution ban. Written as
  instructions to an agent, not as prose about the project. Keep it under 200 lines;
  put anything task-specific in `.claude/skills/` and anything language-specific in
  `.claude/rules/` with `paths:` frontmatter instead.

## 5. The agentic git workflow

### Layout

```
<repo>/                          primary checkout, always on the default branch
../.worktrees/<repo>/<branch>/   one directory per branch, created and removed by script
```

### `scripts/agent`

A single POSIX `sh` entry point. Every subcommand is idempotent, prints what it is
about to do, and fails loudly with a non-zero exit and an actionable message.

| Command | Behaviour |
|---|---|
| `doctor` | Verify `git`, `gh`, auth, toolchain, hooks path, clean tree. Report every problem, not just the first. |
| `start <type>/<slug>` | Validate the branch name against `^(feat\|fix\|chore\|docs\|perf\|refactor\|test)/[a-z0-9][a-z0-9._-]*$`. `git fetch origin`, create the branch from `origin/<default>`, add the worktree, print the `cd` path. Refuse if the branch or worktree already exists. |
| `check` | `scripts/task check`. |
| `commit <msg>` | Validate Conventional Commits, then `git commit`. Never `--no-verify`. |
| `pr [--draft]` | Run `check`, push with `-u`, `gh pr create --fill` (body from the template; strip any attribution), print the URL. |
| `sync` | Rebase the current branch onto `origin/<default>`. |
| `done` | Confirm the PR is merged, remove the worktree, delete the local branch, `git worktree prune`, `git fetch --prune`. |
| `list` | Worktrees with their branches and PR state. |

Never `cd` for the caller — print the path and let them move. Resolve the default
branch dynamically; do not hardcode `main` in the script.

### Hooks

Ship them in `.githooks/`, tracked in git, wired via
`git config core.hooksPath .githooks`. Set that in `scripts/agent doctor` and in a
`scripts/setup` that `README.md` tells contributors to run once.

- `commit-msg` — reject anything that is not a valid Conventional Commit
  (`type(scope)!: subject`, subject ≤ 72 chars, no trailing period). Also reject any
  line matching the attribution patterns from §1.
- `pre-commit` — `scripts/task fmt:check` and `lint` on staged files where the
  toolchain supports it, otherwise the whole tree. Must stay fast; if it can't be,
  move it to `pre-push`.
- `pre-push` — refuse a push to the default branch. Run `scripts/task check`.

Hooks are `sh`, executable, and dependency-free — no `husky`, no `pre-commit`
framework, no `lefthook`.

## 6. CI

`.github/workflows/ci.yml` — one workflow, triggered on `push` to the default branch
and on `pull_request`:

- `concurrency: group: ${{ github.workflow }}-${{ github.ref }}` with
  `cancel-in-progress: true`.
- `permissions: contents: read` at the top level; widen only per-job.
- Actions pinned to a major version tag at minimum; prefer commit SHAs for any
  third-party action.
- Dependency and toolchain caching via the official setup action.
- One `check` job that runs `scripts/task check`, so CI and local hooks can never drift.
- A `required` job that `needs` every other job and exists purely as the single status
  check to require in branch protection.
- Fail fast, no `continue-on-error`, no matrix unless the project genuinely targets
  multiple runtimes — if it does, keep it to the minimum meaningful set.

`.github/workflows/release.yml` — tag-triggered (`v*`), `permissions: contents: write`,
runs `check`, then `gh release create --generate-notes`. Strip attribution from
generated notes if any appears.

Nothing else. No stale-bot, no auto-labeler, no greeting workflow.

## 7. Create and configure the remote

```sh
git init -b <default>
git add -A
git commit -m "chore: initial commit"        # no attribution trailers
gh repo create <owner>/<name> --public --source=. --remote=origin --push \
  --description "<description>"
```

Then configure the repo — every step verified, not assumed:

- `gh repo edit` → squash merge only (disable merge commits and rebase merges),
  auto-delete head branches on merge, enable issues, disable wiki and projects unless
  asked for.
- Default workflow permissions read-only:
  `gh api -X PUT repos/<owner>/<name>/actions/permissions/workflow -f default_workflow_permissions=read`.
- Branch protection on the default branch via `gh api`: require a PR before merging,
  require the `required` status check to pass and branches to be up to date, require
  conversation resolution, block force pushes and deletions, dismiss stale approvals.
  Set required approvals to 0 for a solo repo (so the owner is not deadlocked) and say
  so explicitly in `CONTRIBUTING.md`; use 1 if the project has other maintainers.
- Enable private vulnerability reporting and Dependabot alerts.
- Add repository topics.

Anything the current `gh` token cannot do (org-level settings, some rulesets), do not
silently skip: finish everything else and report the exact command the owner must run
themselves.

## 8. Prove it end to end

Do not report success from having written files. Actually exercise the loop:

1. `scripts/agent doctor` → clean.
2. `scripts/agent start chore/verify-workflow` → worktree created at the expected path.
3. In the worktree, make a trivial real change (e.g. add the CI badge to `README.md`),
   `scripts/agent commit`, `scripts/agent pr`.
4. Watch CI: `gh pr checks --watch`. It must go green.
5. Confirm a direct push to the default branch is rejected — locally by the hook and on
   the server by branch protection.
6. Merge the PR (squash), `scripts/agent done`, and confirm the worktree, local branch,
   and remote branch are all gone.
7. `git log --oneline` on the default branch shows exactly the initial commit and the
   squashed PR — and contains no attribution anywhere.

If any step fails, fix the cause and rerun it. A checklist item is only done once you
have watched it pass.

## 9. Report

Close with, briefly:

- Repo URL and the commands to get started.
- The `scripts/agent` command table.
- Anything the owner must configure by hand, with the exact command.
- Any deliberate omission and why.

## Anti-patterns

- Generating docs with `TODO` or `<!-- fill this in -->` left behind.
- Stack-specific commands hardcoded into CI or hooks instead of `scripts/task`.
- `--no-verify`, `continue-on-error`, or `|| true` to make something pass.
- Committing to the default branch "just for setup" after §7 — the initial commit is
  the only one.
- A README that describes aspirations rather than what the code does today.
- Adding a linter/formatter without also adding its config and making it pass.
- Reporting completion with unverified steps.
