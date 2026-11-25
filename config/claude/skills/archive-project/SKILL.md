---
name: archive-project
description: Take a project off the machine without losing work — audit everything the project owns, save uncommitted and unpushed work to the remote, write a verified cold copy, then remove the checkout, worktrees, processes and scratch files and mark the GitHub repository read-only. Use when the user stops work on a project, asks to archive, shelve, retire or clean up a project, or wants disk and memory back from a project they no longer run.
argument-hint: [project name]
---

# Archive a project

A project is more than its checkout. It also owns worktrees, uncommitted work,
ignored files with settings or secrets, running processes, services, containers
and agent scratch files. The script finds all of them and deletes nothing before
the work exists in two places.

```sh
ARCHIVE="$SKILL_DIR/scripts/archive-project"
"$ARCHIVE" audit <name>
```

`$SKILL_DIR` is the directory of this skill. The project is `~/Code/<name>`. Run
the script from a directory outside it.

## 1. Audit, and read the output

`audit` changes nothing. It prints the remote, each worktree with its count of
uncommitted files, the ignored files that the cold copy keeps, live processes,
services and containers, and the agent transcripts and scratch files.

Two findings are blockers. The script does not solve them, because each needs
a decision from the user:

| Blocker | What to do |
|---|---|
| No remote named `origin` | Ask for the repository name. Create a **private** repository (`gh repo create <owner>/<name> --private --source . --remote origin`). Never a public one unasked. |
| A systemd unit or a container belongs to the project | It holds state: a database volume, a served site. Ask whether to export the data, then stop and remove it by hand. Look also for a Caddy route and a remote deployment (`railway.toml`, `vercel.json`, `fly.toml`, a `deploy/` folder); the script cannot see those. Report a live deployment, do not take it down. |

Tell the user what the audit found before you go on. Uncommitted work in a
worktree is the finding that most often surprises them.

## 2. Run

```sh
"$ARCHIVE" run <name> [--snapshot] [--kill]
```

- `--snapshot` is necessary when a worktree has uncommitted work. Each such
  worktree becomes one commit on `archive/wip/<branch>`. The branch itself and
  the default branch do not change.
- `--kill` is necessary when processes run inside the project. Ask first: a
  process may be another agent session or a dev server the user still wants.
  The script never stops the session that runs it.

The order is fixed: snapshot, push every branch and tag, write the cold copy,
test-clone the bundle, stop processes, remove the directories, mark the GitHub
repository as archived. A failure before the removal leaves the project intact.

The cold copy is `~/Code/.archive/<name>/`:

| File | Content |
|---|---|
| `repo.bundle` | Every ref. Restores the project when the remote is gone. |
| `ignored-N.tar` | Ignored files that are not build output, for example `.env`. |
| `agent/` | Transcripts of the agent sessions. |
| `MANIFEST` | Date, host, path, remote, and each branch with its commit. |

When the session that runs the script works inside the project, the script
leaves that session's transcripts and scratch files and prints their paths.
Tell the user; remove them from a session that works somewhere else.

## 3. Report

State what was saved and where: the remote, the snapshot branches by name, the
cold copy path, and the disk space returned. State what the script left: the
caller's own scratch files, a cache directory under `~/.cache` that only looks
related, a live deployment.

## Restore

```sh
"$ARCHIVE" restore <name>
```

`restore` clones from the remote, or from the bundle when the remote does not
answer, and puts back the ignored files of the main checkout. The GitHub
repository stays read-only until `gh repo unarchive <owner>/<name>`.
