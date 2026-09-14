---
name: roadmap
description: Build a complete cairn roadmap for a project — v0.1 MVP through to v1.0 — with milestones, dependency-ordered items, priorities, effort, areas, and spikes for the unknowns. Use when a project has no backlog, when asked to plan a release, define milestones, scope an MVP, or lay out the path to v1.
argument-hint: [what the project is, if cairn isn't set up yet]
---

# A roadmap from MVP to v1.0

You are turning an idea, or a half-built project, into a backlog somebody could work
through without asking you a single question. The output is cairn items, not a document.

Execute. Ask at most one batched round of questions, at §1, and only about things you
genuinely cannot infer.

## 0. Know what you are planning

Before writing a single item, read:

- `README.md`, the docs, and any design notes — what it claims to be.
- The code. Entry points, module layout, what already works. An honest roadmap starts
  from what exists, not from zero.
- `git log --oneline | head -50` — the direction of travel.
- `cairn list --all` if a backlog already exists. You are extending it, not replacing it.

State in three sentences: what this is, who it is for, what it does today. If you cannot,
you do not know enough to plan it — go read more.

## 1. Fix the destination first

**v1.0 is the promise, not the wish list.** It is the version where the thing does what
it claims, reliably, for someone who is not you: installable, documented, tested, and
stable enough that breaking it would be a bug rather than a Tuesday.

Decide and write down, before any milestone exists:

- **The one sentence.** What v1.0 does. If it needs "and", you have two projects.
- **The user.** Who runs this, on what, to do what.
- **Out of scope, explicitly.** The three things people will assume are in v1 that are
  not. This is the most load-bearing list in the roadmap and the one most often skipped.
- **What "done" is measured by.** Works on the platforms you name, the suite is green,
  the docs let a stranger start, install is one command.

Ask the owner only what you truly cannot infer: the one sentence, the user, the platforms,
and anything deliberately out of scope. One batched question, then proceed.

## 2. Set up the schema

If there is no `cairn.toml`, `cairn init`, then shape the schema to the project rather
than accepting the default:

- **Types** — `feature bug chore spike docs milestone`. Give `feature`, `bug`, and
  `spike` body templates so every item arrives with the same headings. A `spike`
  template must have an **Answer** section: a spike that closes without one was wasted.
- **Statuses** — `backlog planned` (open), `doing blocked` (active), `done dropped`.
  Don't invent more. Columns are not progress.
- **`area`** — an enum of this project's actual subsystems, taken from the module layout
  you read in §0. This is what makes `cairn list --area` useful later, and it is worth
  ten minutes of thought.
- **`priority`** — `p0 p1 p2 p3`, where **p0 means "this milestone does not ship without
  it"**. Priority is relative to a milestone, not to the universe.
- **`effort`** — `s m l xl`. A rough size, never an estimate, never a date.
- **Views** — `now` (active work), `next` (ready, ranked), `triage` (open with no
  milestone). They cost one line each and get used every day.

Run `cairn config` and read back what resolved. Then `cairn check`.

## 3. Shape the milestones

Three to five, never more. Each one **ships**: it is a version somebody could install and
use, not a phase of construction. "Backend done" is not a milestone; "reads a file and
prints it" is.

Work **backwards** from v1.0 — it is the only ordering that keeps the MVP honest, because
everything in v0.1 has to be there for something later.

**v0.1 — the MVP.** The smallest version that is *real*: it does one useful thing, end to
end, for one user, badly but correctly. Not a prototype, not a skeleton with `todo!()` in
it. If you removed any item from v0.1 it would stop being usable; if you can remove one
and it still works, that item belongs in v0.2.

The test: **could a stranger install it, use it for its one purpose, and get value?** If
no, it is not an MVP. If it does three things, it is not minimal.

**v0.2 … v0.n — the middle.** Each takes v0.1 from "works for me" toward "works for
anyone": the second real use case, the errors that currently panic, the configuration
people immediately want, the platform you skipped. Order them by what the *previous*
milestone taught you it needed.

**v1.0 — the promise.** No new surface. Stability, completeness against §1, documentation,
packaging, and the API you are willing to be held to. If v1.0 contains an exciting
feature, it is in the wrong milestone.

**`later`** — a real milestone that catches everything good that is not on the path. Items
go here instead of being argued about.

Create them as items with dates that are decisions, not hopes:

```sh
cairn new "v0.1" --type milestone --set due=YYYY-MM-DD \
  --body "$(cat <<'EOF'
## Ships

One sentence on what a user can do when this lands.

## Done when

- [ ] …

## Explicitly not in this milestone

- …
EOF
)"
```

## 4. Fill the milestones

For each milestone, from v0.1 forward, decompose until every item is one unit of work: a
single branch, a single PR, a day or less. Larger than that and it is two items; smaller
and it is a line in a checklist, not an item.

Every item gets, at creation:

```sh
cairn new "Attach to a running terminal" \
  --type feature --milestone v0.1 \
  --set priority=p0 --set effort=m --set area=runtime \
  --depends-on 12 --depends-on 14
```

- **`--milestone`** — always. An item with no milestone is invisible to planning; that is
  what the `triage` view is for catching.
- **`--depends-on`** — the real edges, not a wish for ordering. These are what make
  `cairn next` able to tell someone what to start, and what stops ten items being claimed
  that all block on the same one.
- **`priority`** — inside this milestone. Every p0 is a thing that stops the release.
  If everything is p0, nothing is; expect roughly a third.
- **`effort`** — s/m/l/xl. An `xl` is a signal you have not decomposed far enough. Split
  it or replace it with a spike.
- **`area`** — always, so the backlog can be sliced by subsystem later.

The body carries what the title cannot. A feature needs **Problem / Proposal /
Acceptance criteria**, and the acceptance criteria are checkboxes someone else could
verify without asking you what you meant. "Handles errors" is not acceptance criteria;
"a missing file exits 1 with the path in the message, covered by a test" is.

**Every milestone needs the unglamorous items**, and they are the ones that get left out:
CI, the test harness, the error type, `--help`, the README section, the install path, the
release workflow. Put them in the milestone that needs them, not in a "polish" bucket at
the end, which never happens.

**Unknowns become spikes, not features.** Anything you cannot write acceptance criteria
for is a question, and it gets a `spike` item that other items depend on. A spike is
timeboxed, closes with a written answer, and usually spawns the real items. Planning
around an unknown by writing a confident feature item is how roadmaps become fiction.

## 5. Check the shape

Before you call it done, interrogate what you built:

```sh
cairn check                                   # schema-valid: must pass
cairn roadmap --items                         # read the whole thing top to bottom
cairn next -m v0.1 --blocked                  # is there ready work on day one?
cairn list -m v0.1 --filter 'priority=p0' --count
cairn list --filter 'milestone=' --count      # unfiled items: should be 0
```

Then the questions a tool cannot ask for you:

- **Is there ready work?** If `cairn next` on the first milestone is empty or every item
  is blocked, the dependency graph is wrong. Someone must be able to start today.
- **Is the MVP still minimal?** Look at v0.1 and remove the least essential item. Does it
  still work end to end? If yes, remove it for real and move it to v0.2. Repeat.
- **Could someone else work this?** Pick three items at random. Could a competent stranger
  implement each from the body alone? If not, the body is too thin.
- **Where will it slip?** Name the two items most likely to be three times their effort.
  Say so in the report. Usually: anything touching a platform you haven't tested on,
  anything with "just" in the description.
- **Does each milestone ship something?** Read the "Ships" line of each. If one of them
  is construction rather than a release, merge it into its neighbour.
- **Is anything in v1.0 a new feature?** Move it to `later`.

Finally: `cairn render` to regenerate `ROADMAP.md`, and commit it as a `docs:` change
through a branch and a PR like anything else. Never hand-edit the generated file.

## 6. Report

- The one sentence, the user, and what is explicitly out of scope.
- Each milestone: its date, what it ships, item count, p0 count.
- What is ready to start right now (`cairn next -m <first>`).
- The spikes, and what each one blocks.
- The two most likely places this slips, and why.
- Any question you had to answer with an assumption, stated plainly so it can be
  corrected.

## Anti-patterns

- A v0.1 that cannot be used — a skeleton, a "foundation", a milestone of plumbing.
- Milestones named after layers (`backend`, `frontend`, `polish`) rather than releases.
- Dates chosen to look ambitious. A date nobody believes is worse than no date.
- Items with no acceptance criteria, so "done" is whoever gets bored first.
- Every item p0, or no dependencies anywhere — both mean the ranking is decorative.
- `xl` items left whole because splitting them is hard. That is exactly why.
- An unknown written as a confident feature instead of a spike.
- v1.0 used as a dumping ground for everything not yet built.
- Writing the roadmap into a `PLAN.md` instead of cairn items. The board is the plan.
