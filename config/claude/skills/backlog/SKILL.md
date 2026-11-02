---
name: backlog
description: Plan work in cairn so a stranger could pick up any item and build it — shaping the schema to the project, writing items with verifiable acceptance criteria, wiring real dependency edges, turning unknowns into spikes, and auditing the graph before calling it done. Use whenever work needs breaking down into cairn items: planning a feature or a release, filing tickets, adding to an existing backlog, splitting something too large, revising a plan after research, or auditing a backlog that has gone stale. Also use when asked to "plan this out", "make tickets", "break this down", or "get this into cairn" — the board is the plan, so the output is always items, never a document.
argument-hint: [what to plan, or a milestone key to fill in]
---

# Plan work in cairn

The output is cairn items. Not a `PLAN.md`, not a message full of bullets that
evaporates when the session ends. If it matters, it is an item; if it is not an item,
it did not happen.

The test every item must pass: **a competent stranger, with no memory of this
conversation, could pick it up and build the right thing.** That is a high bar and it
is the whole point. You are writing for whoever works this next, which may be you with
none of today's context.

## 1. Know what you are planning

Read before writing. `README`, the docs, the entry points, the module layout,
`git log --oneline | head -50`, and `cairn list --all` if a backlog already exists —
you are extending it, not replacing it.

Then say in three sentences what this is, who it is for, and what it does today. If you
cannot, you do not know enough to plan it. Go read more.

Planning from an imagined codebase produces items that are quietly wrong, and nobody
discovers it until someone tries to build one.

## 2. Shape the schema to this project

`cairn init` gives you a starting schema. Accepting it unchanged is a missed
opportunity, because the fields are what make the backlog queryable later.

- **`area`** — an enum of this project's actual subsystems, taken from the module layout
  you just read. This is what makes `cairn list --area` worth running. Ten minutes here
  pays for itself.
- **`priority`** — `p0 p1 p2 p3`, where **p0 means the milestone does not ship without
  it**. Priority is relative to a milestone, not to the universe. Expect roughly a
  third to be p0; if everything is p0, nothing is.
- **`effort`** — `s m l xl`. A size, never an estimate, never a date. An `xl` is a
  signal you have not decomposed far enough.
- **A `spike` type with an **Answer** section in its template.** A spike that closes
  without a written answer was wasted, and the template is what prevents it.
- **Body templates** for `feature` and `bug`, so every item arrives with the same
  headings and nobody has to remember the shape.

Run `cairn config` and read back what resolved, then `cairn check`. Both are cheap and
catch a schema that does not say what you meant.

## 3. Write items a stranger could implement

One item is one unit of work: a single branch, a single PR, a day or less. Bigger is
two items. Smaller is a line in someone's checklist, not an item.

The body carries what the title cannot:

```markdown
## Problem
Why this exists. What is broken or missing today.

## Proposal
What to do. Concrete enough to argue with — a signature, a file, a syntax sketch,
the actual flag name.

## Acceptance criteria
- [ ] checkboxes someone else could verify without asking you what you meant
```

Acceptance criteria are where items live or die. "Handles errors" is not a criterion;
"a missing file exits 2 with the path in the message, covered by a test" is. The test:
could a reviewer tick this box without a conversation?

**Every milestone needs the unglamorous items**, and they are the first to be left out:
CI, the test harness, the error type, `--help`, the README section, the install path,
the release workflow. Put each in the milestone that needs it, never in a "polish"
bucket at the end, which never happens.

## 4. Create in bulk, with the edges, in one pass

Creating fifty items by hand is fifty chances to forget a field. Write one shell script
instead, capture the ids as you go, and pass them as dependencies. This is the
technique that makes a large backlog tractable:

```sh
export CAIRN_NO_HOOKS=1          # do not re-render ROADMAP.md fifty times
n() { cairn new -q "$@"; }       # -q prints only the new id

lex=$(n "Lexer with source spans" --type feature -m v0.1 \
  --set priority=p0 --set effort=m --set area=syntax --body "$(cat <<'EOF'
## Problem
Every diagnostic needs a position; the parser needs tokens.

## Acceptance criteria
- [ ] an illegal character produces an error with its span, not a panic
EOF
)")
parse=$(n "Parser for the core language" --type feature -m v0.1 \
  --set priority=p0 --set effort=l --set area=syntax -d "$lex" --body "...")
```

Then `cairn render` once at the end. Keep the script in a scratch directory — it is
scaffolding, not a deliverable, and the items are the artifact.

**`--depends-on` is what makes the backlog usable.** `cairn next` can only tell someone
what to start if the edges are real. Wire the true ones, not a wish for ordering: ten
items that all block on the same unbuilt thing is a discovery, not a problem.

## 5. Unknowns become spikes, not confident features

Anything you cannot write acceptance criteria for is a question, and a question gets a
`spike`: timeboxed, closing with a written answer, usually spawning the real items.
Other items depend on it.

Planning around an unknown by writing a confident feature item is how roadmaps become
fiction. When you notice you are inventing detail to fill a section, stop and make it a
spike. Put the leading option in the body — a spike with a hypothesis closes faster
than one with an open shrug.

## 6. Audit the graph before calling it done

`cairn check` validates against the schema. It cannot tell you the plan is shaped
wrong. Run the bundled audit, which catches what review misses:

```sh
cairn list --all --json | python3 "$SKILL_DIR/scripts/audit-backlog.py"
```

where `$SKILL_DIR` is this skill's own directory. Pass `--root` if you are not running
from the project root.

It reports backward milestone edges (an early item depending on a later one, which
blocks a release forever), unfiled items, `xl` work still on the critical path,
milestones with no startable item, feature items with no acceptance criteria, and
criteria stranded outside their section. Fix what it finds, then run it again.

Then ask the questions a script cannot:

- **Is there ready work?** `cairn next -m <first>`. If it is empty or everything is
  blocked, the edges are wrong. Someone must be able to start today.
- **Is the central idea load-bearing in the first milestone?** The thing that makes
  this project worth building should be exercised in the first release, not decorated.
  If the first milestone would still be shippable with the distinctive part removed,
  it is in the wrong milestone.
- **Could someone else work this?** Pick three items at random. Could a stranger
  implement each from the body alone? If not, the bodies are too thin.
- **Where will it slip?** Name the two items most likely to take three times their
  effort — usually anything touching an untested platform, and anything whose
  description contains "just".

## 7. Argue with the plan, then revise it

A first draft is always wrong in a specific way, and the way is usually one of these:

- **A research rabbit hole nobody asked for.** A subsystem appearing across four
  milestones for a property the project never promised. Cut it to a deferred item that
  records why.
- **Milestones named after layers** — `backend`, `frontend`, `polish` — rather than
  releases. Each milestone must ship something a person could use.
- **A decoration mistaken for the idea.** See the load-bearing test above.
- **Items that are two items** wearing one title.

So argue the plan against itself before anyone builds it, and make the revision in
cairn rather than in prose. Deleting twenty items and recording why is progress.
A deferral with a reason on it is a decision; a vague `later` pile is an argument
waiting to happen again.

## 8. Record decisions on the items

A decision made in conversation is lost. `cairn note <id> "..."` puts it where the
person doing the work will read it:

```sh
cairn note 35 "Chose a journal of old values over user-written inverses: a hand-written
inverse can be wrong, and a wrong inverse silently breaks the one thing rollback promises."
```

Note the **why**. The diff will show what.

One gotcha: `cairn note` appends a dated section, so a checkbox added through a note
lands *below* the Acceptance criteria heading rather than in it. Criteria belong in the
section. Edit the body for those and keep notes for prose.

## 9. Report

Close with the destination in one sentence, each milestone with what it ships and its
counts, what is ready to start right now, the spikes and what each blocks, the two
places it will slip, and any assumption you had to make so it can be corrected.

State assumptions plainly. An unstated assumption is the one that turns out wrong.

## Reference

`references/cairn.md` — the command and flag detail worth having to hand: bulk editing,
scripting with `--json`, filters, clearing a field, and the gotchas that cost time.

## Anti-patterns

- Writing the plan into a document instead of items. The board is the plan.
- Acceptance criteria nobody can verify, so "done" is whoever gets bored first.
- Every item p0, or no dependency edges anywhere. Both mean the ranking is decorative.
- A first milestone that cannot be used: a skeleton, a "foundation", a layer of plumbing.
- An `xl` left whole because splitting it is hard. That is exactly why.
- An unknown written as a confident feature instead of a spike.
- A deferral with no recorded reason, which is just an argument postponed.
- Reporting a backlog as done without running `cairn check` and the audit.
