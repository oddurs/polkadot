#!/usr/bin/env python3
"""Audit a cairn backlog for shape problems `cairn check` cannot see.

    cairn list --all --json | python3 audit-backlog.py [--root DIR]

`cairn check` validates items against the schema. It cannot tell you the plan is
shaped wrong: that an early item waits on a late one, that nothing is startable,
or that a body has no verifiable criteria. This does.

Exits 1 if anything was found, so it can gate a planning session.
"""

import argparse
import json
import re
import sys
from collections import defaultdict
from pathlib import Path

DONE = {"done", "dropped"}
NEEDS_CRITERIA = {"feature", "bug"}


def milestone_order(items):
    """Milestones in shipping order. Undated ones sort last: they are the
    deferral bucket, and work parked there is not on the critical path."""
    stones = [i for i in items if i["type"] == "milestone"]
    dated = sorted(
        (m for m in stones if m["fields"].get("due")),
        key=lambda m: m["fields"]["due"],
    )
    undated = [m for m in stones if not m["fields"].get("due")]
    order = {}
    for n, m in enumerate(dated + undated):
        key = m.get("key") or m["title"]
        order[key] = n
    dateless = {(m.get("key") or m["title"]) for m in undated}
    return order, dateless


def body_of(item, root):
    path = item.get("path")
    if not path:
        return None
    full = Path(root) / path
    if not full.exists():
        return None
    text = full.read_text(encoding="utf-8")
    return re.sub(r"(?s)^---\n.*?\n---\n", "", text, count=1)


def sections(body):
    """Body split into (heading, content). Heading is None for any preamble."""
    parts = re.split(r"(?m)^##\s+(.*)$", body)
    out = [(None, parts[0])]
    for i in range(1, len(parts), 2):
        out.append((parts[i].strip(), parts[i + 1]))
    return out


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--root", default=".", help="project root, for reading item bodies")
    args = ap.parse_args()

    items = json.load(sys.stdin)
    by_id = {i["id"]: i for i in items}
    work = [i for i in items if i["type"] != "milestone"]
    order, dateless = milestone_order(items)
    findings = defaultdict(list)

    def note(kind, message):
        findings[kind].append(message)

    for i in work:
        ref, title, ms = i["ref"], i["title"], i.get("milestone")

        # An item filed nowhere is invisible to planning.
        if not ms:
            note("unfiled", f"{ref} {title}")

        for d in i.get("depends_on", []):
            dep = by_id.get(d)
            if dep is None:
                note("missing-dep", f"{ref} depends on {d:04d}, which does not exist")
                continue
            if dep["type"] == "milestone":
                note("dep-on-milestone", f"{ref} depends on the milestone {dep['title']}")
                continue
            dms = dep.get("milestone")
            # The edge that silently blocks a release forever.
            if ms in order and dms in order and order[dms] > order[ms]:
                note(
                    "backward",
                    f"{ref} ({ms}) depends on {dep['ref']} ({dms}): {dep['title']}",
                )

        if i["fields"].get("effort") == "xl" and ms in order and ms not in dateless:
            note("xl", f"{ref} ({ms}) {title} — decompose it or make it a spike")

        body = body_of(i, args.root)
        if body is None:
            note("unreadable", f"{ref} {title}")
            continue
        if not body.strip():
            note("empty", f"{ref} {title}")
            continue

        boxes = re.findall(r"(?m)^\s*-\s*\[[ xX]\]", body)
        # A deferred item is a sketch on purpose; criteria there would be
        # invented detail for work nobody has committed to.
        scheduled = ms in order and ms not in dateless
        if i["type"] in NEEDS_CRITERIA and scheduled and not boxes:
            note("no-criteria", f"{ref} {title}")

        # A criterion added through `cairn note` lands in a dated section, where
        # a reader scanning for acceptance criteria will not see it.
        for heading, content in sections(body):
            if heading and "acceptance criteria" in heading.lower():
                continue
            if re.search(r"(?m)^\s*-\s*\[[ xX]\]", content):
                where = heading or "the preamble"
                note("stray-criteria", f"{ref} has a checkbox under '{where}'")

        if i["type"] == "spike" and "## Answer" not in body:
            note("spike-no-answer", f"{ref} {title} — no Answer section to fill in")

    # Per-milestone shape.
    keys = sorted(order, key=lambda k: order[k])
    current = next(
        (
            k
            for k in keys
            if k not in dateless
            and any(
                i.get("milestone") == k and i["status"] not in DONE for i in work
            )
        ),
        None,
    )
    for key in keys:
        mine = [i for i in work if i.get("milestone") == key]
        if not mine:
            continue
        # Only the milestone in flight needs startable work. A later one being
        # fully blocked is the dependency graph doing its job.
        if key == current:
            open_items = [i for i in mine if i["status"] not in DONE]
            startable = [
                i
                for i in open_items
                if all(
                    by_id[d]["status"] in DONE
                    for d in i.get("depends_on", [])
                    if d in by_id
                )
            ]
            if not startable:
                note(
                    "nothing-startable",
                    f"{key} is the milestone in flight and every open item is blocked",
                )
        if key not in dateless:
            p0 = sum(1 for i in mine if i["fields"].get("priority") == "p0")
            share = p0 / len(mine)
            first_dated = next((k for k in keys if k not in dateless), None)
            if key == first_dated:
                pass  # see above
            elif share > 0.85:
                note("p0-inflation", f"{key}: {p0} of {len(mine)} are p0 — ranking is decorative")
            elif share < 0.1 and key != first_dated:
                note("p0-absent", f"{key}: {p0} of {len(mine)} are p0 — what blocks the release?")

    seen = defaultdict(list)
    for i in work:
        seen[i["title"].strip().lower()].append(i["ref"])
    for title, refs in seen.items():
        if len(refs) > 1:
            note("duplicate", f"{', '.join(refs)} share a title: {title}")

    labels = {
        "backward": "Backward dependency: an item waits on a later milestone",
        "unfiled": "Filed under no milestone, so invisible to planning",
        "missing-dep": "Depends on an item that does not exist",
        "dep-on-milestone": "Depends on a milestone rather than on work",
        "nothing-startable": "No startable item: nobody could begin here today",
        "xl": "xl effort on the critical path",
        "no-criteria": "No acceptance criteria, so done is whoever gets bored first",
        "stray-criteria": "Criterion outside the acceptance criteria section",
        "spike-no-answer": "Spike with no Answer section",
        "empty": "Empty body",
        "unreadable": "Body could not be read (wrong --root?)",
        "p0-inflation": "Almost everything is p0",
        "p0-absent": "Nothing is p0",
        "duplicate": "Duplicate title",
    }
    order_of_report = [
        "backward", "missing-dep", "nothing-startable", "unfiled", "empty",
        "unreadable", "no-criteria", "stray-criteria", "spike-no-answer",
        "dep-on-milestone", "xl", "p0-inflation", "p0-absent", "duplicate",
    ]

    total = sum(len(v) for v in findings.values())
    if not total:
        print(f"backlog audit: clean ({len(work)} items, {len(order)} milestones)")
        return 0

    for kind in order_of_report:
        if kind not in findings:
            continue
        print(f"\n{labels[kind]} ({len(findings[kind])})")
        for line in findings[kind]:
            print(f"  {line}")
    print(f"\n{total} finding(s) across {len(work)} items")
    return 1


if __name__ == "__main__":
    sys.exit(main())
