---
name: ship
description: Take the current branch from working tree to open pull request — format, lint with warnings denied, full test suite, conventional commit, push, and a PR description that says where the change is weak. Use when work is finished and asked to ship, open a PR, commit and push, or "make it green".
argument-hint: [--draft]
---

# Green, then a pull request

Nothing here is optional and nothing here is skipped to save time. If a step fails, fix
the cause. Never `--no-verify`, never `|| true`, never `continue-on-error`.

## 1. Where are you

Refuse to continue if this is the default branch — the work needs a branch and a
worktree first (`/agent-workflow`). Confirm the diff is only what this unit of work is
about; a stray edit from another task goes to its own branch.

## 2. Green

```sh
scripts/task check
```

If the repo has no `scripts/task`, run its real equivalents, in this order, and say in
your report which commands you used:

| | Rust | Node | Python | Go |
|---|---|---|---|---|
| format | `cargo fmt --all -- --check` | `<pm> run format:check` | `uv run ruff format --check .` | `gofmt -l .` |
| lint | `cargo clippy --all-targets --all-features -- -D warnings` | `<pm> run lint` | `uv run ruff check .` | `go vet ./...` |
| test | `cargo test --all-targets` | `<pm> test` | `uv run pytest` | `go test ./...` |
| build | `cargo build --all-targets` | `<pm> run build` | — | `go build ./...` |

Use the package manager the lockfile names. Warnings are errors — a lint step that
passes while printing warnings has not passed.

Show the output. A claim that the suite is green without the run is worth nothing.

## 3. Commit

Conventional Commits. Imperative, subject ≤ 72 characters, no trailing period:

```
fix(proc): start the pty reader before the child

The reader thread was spawned after the child, so a command that wrote and
exited immediately could finish before anything was reading.

Refs: 0037
```

The body explains **why**; the diff already says what. Reference the tracker item in a
`Refs:` trailer when the repo has one.

**No attribution.** No co-author trailer naming a model, no "generated with" footer, no
robot emoji — not in the commit, not in the PR body, not in the release notes.

Several unrelated changes in the tree mean several commits, or better, several branches.

## 4. Push and open it

```sh
git push -u origin HEAD
gh pr create --fill   # or --draft
```

The description has three parts and no padding:

- **Problem** — what was wrong, what it cost. Link the issue.
- **Approach** — what you did, and the alternative you rejected if the choice was close.
- **Look at this sceptically** — the weakest part of the change. Name it yourself rather
  than letting a reviewer find it. "No test for the timeout path; it needs a fake clock
  and I'd rather do that separately" is a good line. If you genuinely can't think of
  one, say what you'd break first if you were attacking it.

If the repo has a PR template, fill it honestly; don't tick a box for a step you skipped.

## 5. Watch it

```sh
gh pr checks --watch
```

Red CI is your problem, now, not a follow-up. Report the URL, the check status you
actually observed, and anything you left undone.

## Refuse to ship

- The suite is red, or you haven't run it.
- You're on the default branch.
- The diff contains a secret, a key, or a `.env`.
- The change is two changes.
- There's a `TODO` you added, a commented-out block, or debug logging left in.
