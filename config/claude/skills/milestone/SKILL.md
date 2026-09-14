---
name: milestone
description: Drive a whole cairn milestone to done — work each ready item in its own worktree and branch, review the diff before it becomes a PR, keep the backlog honest as you go, and close the milestone only when cairn check passes and nothing is left open. Use when asked to finish, complete, clear, or ship a milestone, or to work through the backlog end to end.
argument-hint: [milestone key, e.g. v0.1]
---

# Finish a milestone

A milestone is a queue of items, each of which is a separate unit of work with its own
worktree, branch, review, and pull request. You are running that loop until the queue is
empty — not doing everything on one branch, and not batching the reviews to the end.

Execute. This skill runs for a long time; report as you go rather than in one dump at the
end.

## 0. Establish the milestone

```sh
cairn roadmap "$MILESTONE" --items          # what is in it, and progress
cairn list -m "$MILESTONE" --json           # the full set, machine-readable
cairn next -m "$MILESTONE" --blocked        # what is ready, and what is blocked by what
cairn check                                 # the backlog must be valid before you start
```

If `$ARGUMENTS` is empty, print the roadmap and ask which milestone — that is the one
question worth asking here.

Then look at the queue as a whole before touching anything:

- **Order.** `cairn next` already excludes work blocked by unfinished dependencies and
  ranks what is left. Take its order unless you can see a reason not to, and say the
  reason.
- **Scope.** An item that is really three items gets split now (`cairn new`, `part_of`),
  not discovered halfway through a branch.
- **Rot.** An item that no longer makes sense is `cairn close`d or dropped with a note
  saying why. Do not implement something because it is written down.
- **Blocked.** Anything blocked by work outside this milestone is called out now. It is
  the most likely reason the milestone doesn't finish.

Report the plan in a few lines: the order, anything you split, anything you dropped and
why, anything blocked. Then start.

## 1. The per-item loop

For each ready item, in order. One item, one worktree, one branch, one PR.

**Claim it.**

```sh
cairn claim <ID>        # or: cairn claim --next
```

Claiming before you start is what stops two agents doing the same item.

**Branch.** `<type>/<id>-<slug>`, where the type matches the item's type
(`feature`→`feat`, `bug`→`fix`, `chore`→`chore`, `docs`→`docs`, `spike`→`refactor` or
`chore`). The id leads so the branch, the commit, and the item are one thing:

```sh
scripts/agent start fix/0041-attach-to-a-terminal
```

No `scripts/agent`? Use `/agent-workflow`. Work in the worktree; the primary checkout
stays on the default branch.

**Implement.** Smallest change that finishes the item. Record what you learn in the item
as you go, not afterwards from memory:

```sh
cairn note <ID> "Used a channel rather than a mutex: the reader thread needs to outlive the parse."
cairn set <ID> effort=l          # when reality disagrees with the estimate
```

A bug item is finished when the test that would have caught it exists and fails without
the fix.

**Green.**

```sh
scripts/task check
```

**Review before the PR, not after.** Run `/code-review` on the diff. Fix what is real;
for anything you consciously decline, say why in the PR description rather than silently
dropping it. This is the step people skip when there are eleven items left, and it is the
one that keeps item eleven from inheriting a mistake made in item two.

**Ship.** `/ship` — conventional commit, `Refs: <ID>` trailer, push, PR. The description
states the problem, the approach, and the weakest part of the change.

**Close the loop.**

```sh
gh pr checks --watch            # it must go green
# merge (squash), then:
scripts/agent done
cairn close <ID>
cairn render                    # regenerate ROADMAP.md from the items
```

Commit the regenerated roadmap with the next item's work, or as its own `docs:` change —
never by hand-editing the generated file.

## 2. Between items

Each item starts from a fresh `origin/<default>`: `git fetch origin` before the next
`start`, so item *n+1* builds on *n* rather than diverging from it.

If an item you have not started is invalidated by one you just finished, update it now
with a note saying what changed. A backlog that lies is worse than no backlog.

Every few items, re-run `cairn next -m "$MILESTONE" --blocked`. Finishing work unblocks
work, and the ready set is not the one you computed at the start.

## 3. When an item fights back

| Situation | What to do |
|---|---|
| It's bigger than the item says | Stop. Split it in cairn, finish the piece you're on, queue the rest. Don't let one branch grow. |
| It's blocked by something not in the milestone | `cairn set <ID> status=blocked`, note the blocker, move to the next ready item. Raise it in the report. |
| It needs a decision only the owner can make | `cairn propose <ID> <field>=<value> --why "..."` and move on. Don't guess and don't stall the queue. |
| It turns out to be wrong | `cairn close` or drop it with a note. Deleting the work is a legitimate outcome. |
| Review finds a problem in already-merged work | New item, new branch. Never reopen a merged PR's branch. |

## 4. Closing the milestone

Only when the queue is empty:

```sh
cairn roadmap "$MILESTONE" --items --all    # every item done or deliberately dropped
cairn list -m "$MILESTONE" --count          # open items: must be 0
cairn check                                 # must pass
cairn render
```

Then, on the default branch, verify the thing you actually shipped:

- `scripts/task check` green on the merged result, not just on the last branch.
- `git log --oneline` reads as one clean commit per item, every subject a conventional
  commit, no attribution anywhere.
- The milestone item itself is closed, and `ROADMAP.md` reflects it.
- If the milestone is a release, tag it and let the release workflow run — don't hand-write
  notes that `--generate-notes` produces from the merged PRs.

## 5. Report

- Items completed, with PR numbers.
- Items split, dropped, or deferred, and why.
- Anything left blocked, with what blocks it.
- Review findings you declined, and the reason.
- The state of `cairn check` and of the default branch.

## Rules that do not bend

- **Never commit to the default branch**, however small the item.
- One item per PR. A second fix noticed along the way is a new cairn item and a new branch.
- Review every diff before it becomes a PR. Eleven unreviewed PRs is not a milestone.
- Never hand-edit `ROADMAP.md`. Change items, run `cairn render`.
- Never `--no-verify`. `cairn check` and `scripts/task check` both pass before you report
  anything as finished.
- No ad-hoc `TODO.md`, `PLAN.md`, or `NOTES.md`. The backlog is cairn; that is the point
  of it.
