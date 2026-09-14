---
name: agent-workflow
description: Start and run a unit of work under the one-worktree-one-branch-one-PR rule — creates the worktree, keeps the branch in sync, and cleans up after the merge. Use when beginning a task in a repo, when asked to make a branch or worktree, when two tasks need isolating, or when a repo has no scripts/agent and the workflow must be driven by hand.
argument-hint: [type/slug] [what the work is]
---

# One unit of work, one worktree, one branch, one pull request

## First: is it scripted?

```sh
test -x scripts/agent && scripts/agent --help
```

If `scripts/agent` exists, use it and stop reading here — it is the repo's own copy of
this workflow and it wins:

| | |
|---|---|
| `scripts/agent doctor` | is this checkout ready to work in |
| `scripts/agent start <type>/<slug>` | a branch and a worktree of its own |
| `scripts/agent check` | `scripts/task check` |
| `scripts/agent commit <msg>` | validate the message, then commit |
| `scripts/agent pr [--draft]` | check, push, open the PR |
| `scripts/agent sync` | rebase onto the default branch |
| `scripts/agent done [branch]` | confirm merged, remove worktree and branch |
| `scripts/agent list` | worktrees, branches, PR state |

If it doesn't exist, do the same thing by hand as below. If this repo is going to see
more than one task, run `/agentify` first and have the script instead.

## By hand

Resolve the default branch; never assume `main`:

```sh
base=$(git symbolic-ref --quiet --short refs/remotes/origin/HEAD 2>/dev/null || true)
base=${base#origin/}; : "${base:=$(git config --get init.defaultBranch || echo main)}"
```

**Start.** The branch name is `<type>/<slug>`, type one of
`feat fix chore docs perf refactor test`, slug lowercase and hyphenated. When the repo
has a tracker, the id leads: `fix/0041-attach-to-a-terminal`.

```sh
git fetch origin
repo=$(basename "$(git rev-parse --show-toplevel)")
dir="$(dirname "$(git rev-parse --show-toplevel)")/.worktrees/$repo/$branch"
git worktree add -b "$branch" "$dir" "origin/$base"
```

Print the path. Do not `cd` on the user's behalf — tell them where it is.

**Work there.** The primary checkout stays on the default branch, always. Two tasks
never share a checkout: a second task means a second worktree, not a stash.

**Stay current.** `git fetch origin && git rebase "origin/$base"`. Rebase, don't merge —
the PR should be a straight line.

**Ship.** `/ship`.

**Finish.** Only after the PR is merged:

```sh
git worktree remove "$dir"
git branch -d "$branch"
git worktree prune && git fetch --prune
```

## Rules that do not bend

- **Never commit to the default branch.** Not for a typo. If you already did, move the
  commit: `git branch <type>/<slug> && git reset --hard origin/<base>`.
- Never `--no-verify`. The hooks are the point. A hook that is wrong gets fixed, not
  bypassed.
- Never force-push a branch someone else may have pulled. `--force-with-lease` on your
  own PR branch after a rebase is fine.
- One PR does one thing. A drive-by fix you noticed on the way is its own branch.
- If the work turns out to be two units, stop and split it before it grows a shared
  history.

## When something is already wrong

| Symptom | Fix |
|---|---|
| Committed to the default branch, not pushed | `git branch <t>/<slug>; git reset --hard origin/<base>` |
| Worktree exists, branch doesn't (or vice versa) | `git worktree prune`, then start again |
| `start` refuses: branch exists | Someone's already on it — `scripts/agent list`, or pick another slug |
| Rebase conflicts every time | The branch is too old or too large. Split it. |
| Tests pass locally, fail in CI | The hook and CI aren't running the same thing. Both must call `scripts/task check`. |
