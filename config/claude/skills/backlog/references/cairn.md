# cairn: the detail worth having to hand

Command and flag specifics that cost time to rediscover. `man cairn` is the full
reference; this is the subset that comes up while planning.

## Creating

```sh
cairn new "Title" --type feature -m v0.1 \
  --set priority=p0 --set effort=m --set area=syntax \
  -d 12 -d 14 \
  --body "$(cat <<'EOF'
## Problem
...
EOF
)"
```

- `-q` prints only the new id, which is what makes `id=$(cairn new -q ...)` work.
- `--set field=value` for any custom field; repeatable. `-d ID` for dependencies;
  repeatable.
- `-m KEY` files it under a milestone by key, not by id.
- Without `--body`, the type's template seeds the body. That is the point of templates.
- A milestone is an item: `cairn new "v0.1" --type milestone --set due=2026-10-31`.

## Changing

```sh
cairn set 12 status=doing priority=p0        # ids, then assignments, one list
cairn set 12 depends_on=18,21                # a list field takes a comma list
cairn set 12 depends_on=                     # an empty value clears it
cairn set --filter 'milestone=v0.1' status=planned -y
cairn tick 12 3                              # tick criterion 3, counting from 1
cairn tick 12 --all
cairn note 12 "Chose X over Y because ..."   # appends a dated section
cairn remove -f 15 16 17                     # no confirmation
```

`cairn note` is for prose. A checkbox written through a note lands in the dated note
section, below the acceptance criteria heading, where a reader scanning for criteria
will miss it. Put criteria in the section by editing the body.

## Rewriting a body

There is no `--body` on `set`, so rewriting means editing the file while preserving its
frontmatter:

```python
import re, pathlib
p = next(pathlib.Path("cairn/items").glob("0021-*.md"))
text = p.read_text()
head = re.match(r"(?s)^(---\n.*?\n---\n)", text).group(1)
p.write_text(head + "\n" + new_body)
```

Mind the blank line after the frontmatter. A helper that drops the body silently is
easy to write and hard to notice, so check one item with `cairn show` afterwards.

## Reading and scripting

```sh
cairn list --all --json
cairn list -m v0.1 --count
cairn list --filter 'priority=p0' --count
cairn list --filter 'milestone=' --count     # unfiled work; should be 0
cairn next -m v0.1                           # what is startable
cairn next -m v0.1 --blocked
cairn show 12 --json
cairn roadmap --items
cairn board
```

In the JSON, custom fields are nested under `fields` rather than at the top level, `ref`
is the zero-padded id as a string, `path` is the item's file, and a milestone carries
its key in `key`. So priority is `item["fields"]["priority"]`, not `item["priority"]` —
a silent `None` otherwise.

The `--count` line on `cairn next` ("14 ready · 78 blocked") counts the whole backlog,
not the milestone you filtered to.

## Hooks and rendering

The default schema re-renders `ROADMAP.md` after every create, change and remove. That
is right for single edits and wasteful for a hundred, so set `CAIRN_NO_HOOKS=1` (or
pass `--no-hooks`) for a bulk pass and `cairn render` once at the end.

`ROADMAP.md` is generated. Never hand-edit it; change the items and re-render.

## Validating

```sh
cairn check          # every item against the schema
cairn config         # the resolved schema, to confirm it says what you meant
cairn migrate -n     # dry run, when the on-disk format is behind the binary
```

Run `cairn check` after any bulk change. If cairn reports the project is an older
format than the binary, reads still work but writes need `cairn migrate`; dry-run it
first, since it rewrites every item file.

## Git

`cairn init --git` installs a merge driver so two branches that both touched the
backlog resolve `ROADMAP.md` and renumber colliding ids without asking. It lives in
`.git/hooks`, which is not cloned, so it is a once-per-checkout step.
