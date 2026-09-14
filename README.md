# polkadot

A Mac, the way I like it. One Go binary that takes a machine with nothing on it
but the Xcode command line tools and turns it into a working one.

```bash
git clone git@github.com:oddurs/polkadot.git ~/Code/polkadot
cd ~/Code/polkadot
go run . install
```

That is the whole bootstrap. `go run .` rather than a prebuilt binary because
Go is the only prerequisite, and a machine that can't run Go can't build the
binary either.

## Commands

| | |
|---|---|
| `polkadot install` | everything: Homebrew, the Brewfile, symlinks, shell, agents |
| `polkadot link` | symlinks only |
| `polkadot brew` | Homebrew and the Brewfile only |
| `polkadot theme` | Subway Seat only |
| `polkadot sync` | fast-forward the repos, relink what moved |
| `polkadot timer` | install the weekly launchd agent that runs `sync` |
| `polkadot doctor` | report what is and isn't in place, change nothing |

Flags go **before** the command — `polkadot --dry-run install`, not
`polkadot install --dry-run`. Go stops parsing flags at the first positional
argument, so a flag after the command is silently ignored. `--dry-run` prints
the plan without writing anything; `--only=<step>` runs one part (`brew`,
`link`, `shell`, `tools`, `agents`, `theme`).

## What it will not do

**It never writes macOS `defaults`.** No Dock tweaks, no key repeat rates, no
Finder settings. Those are personal, they drift, and reverting them is worse
than setting them. System Preferences stays yours.

**It never deletes a config it did not write.** Anything already sitting at a
target path is moved to `~/.polkadot-backup/<timestamp>/`, keeping its path
shape, so restoring is a plain `mv` back.

## Layout

```
config/      → ~/.config/*     fish, starship, ghostty, lazygit, herdr,
                               atuin, mise, btop, bat, ripgrep, tmux
home/        → ~/*             gitconfig, gitignore_global, commit template
config/claude/    → ~/.claude/*        CLAUDE.md, rules/, skills/, settings.json
config/opencode/  → ~/.config/opencode/  config, AGENTS.md, agents/, commands/
config/codex/     → ~/.codex/config.toml
config/gh/        → ~/.config/gh/config.yml     (not hosts.yml — that is a token)
config/vscode/    → ~/Library/Application Support/Code/User/
Brewfile                       formulae, casks, VS Code extensions
```

The line every one of those follows: **instructions are tracked, state is not.**
`CLAUDE.md` and an agent's command definitions belong in a repo. Session
history, plugin caches, OAuth tokens and 2 MB of somebody else's skill packages
do not.

Everything is symlinked, not copied. Edit in the repo, it is live immediately;
edit in place, it is already staged.

## Subway Seat

[Subway Seat](https://github.com/oddurs/subway-seat), riding **London
Moquette** — the Tube's Corporate Blue turned right down, with brick, hazard
yellow and the standard red. Dark.

The theme files are not in this repo. Subway Seat generates one per app per
flavor across 118 ports; vendoring those here would mean regenerating them
every time the palette moves. `polkadot theme` clones it to `~/Code/subway-seat`
and runs its own installer, which places absolute symlinks into the apps
present — which is why `.gitignore` excludes every directory it writes into.

```sh
sh ~/Code/subway-seat/install.sh switch deep      # another flavor
sh ~/Code/subway-seat/install.sh status           # what drifted
```

`bat`, `delta` and `herdr` render through the terminal's own sixteen colours
(`base16`, `theme = "terminal"`), so they are London Moquette because Ghostty
is, and they follow automatically if the terminal theme ever changes.

The flavor a fresh machine gets is `themeFlavor` in `steps.go`, in one place.

## Editor

[Fresh](https://getfresh.dev), a terminal IDE, is `$EDITOR` and `$VISUAL` —
for git, for `gh`, for lazygit's `e`, and for the `v` abbreviation. Its config
is a repo of its own at `~/.config/fresh`
([fresh-config](https://github.com/oddurs/fresh-config)), because it carries a
Prolog grammar and LSP wiring that has nothing to do with dotfiles.

Neovim is gone. One editor, configured properly, beats two configured halfway.

## Shell

fish, with starship as the prompt. `conf.d/` is split by concern and sourced in
name order:

| | |
|---|---|
| `00-path.fish` | Homebrew, then PATH in priority order |
| `10-env.fish` | editor, pager, XDG, tool defaults |
| `20-gotham.fish` | fish's own syntax colours, and fzf's |
| `30-tools.fish` | mise, starship, zoxide, atuin, fzf, direnv, 1Password |
| `40-abbr.fish` | abbreviations and the three safety aliases |

Abbreviations rather than aliases for anything non-destructive: fish expands
them in place as you type, so what runs is what you can see and history stays
greppable for the real command. `rm`, `mv` and `cp` are aliases instead,
because the point of those is that they apply when you are *not* looking.

`XDG_CONFIG_HOME` is exported, which is what makes lazygit read `~/.config`
instead of `~/Library/Application Support`.

## Staying current

The configuration is three repositories — this one, `subway-seat` for the
theme, `fresh-config` for the editor. `polkadot sync` fast-forwards all three
and relinks only if something actually moved.

```sh
polkadot sync      # once
polkadot timer     # and then never again: Mondays at 09:00
```

Two rules make it safe to run unattended. It is **fast-forward only** — a repo
that has diverged is reported and skipped, never forced. And it **will not
touch a dirty tree**, so work in progress is never merged over or stashed
behind your back. A sync with nothing to do prints three lines and writes
nothing.

It deliberately does not run `brew bundle`. Installing software in the
background while you are working is not a thing a timer should do; `polkadot
brew` is there when you want it.

Log: `~/Library/Logs/polkadot-sync.log`.

## Runtimes

`mise` replaces nvm and pyenv. It reads `.tool-versions`, `.nvmrc`,
`.node-version` and `mise.toml`, so per-project pins keep working, and it costs
a few milliseconds at shell start instead of a few hundred.

## Agents

`claude`, `codex`, `opencode` and `herdr`, each with its config under
`config/` — for Claude Code that means `CLAUDE.md`, the path-scoped
`rules/` and the `skills/`, which are instructions rather than state and so
belong in a repo. History, projects, plugins and auto memory stay out.

`gh` is here too, with aliases for the loop these repos actually run: `gh mine`,
`gh cs` to watch checks, `gh done` to squash-merge and delete the branch.
`gh dash` (a gh extension, not a formula) is the same view without a browser.

opencode carries the most configuration, because it has the most to say: seven
models behind short names that show what each costs per million tokens, a
permission table that allows the read-only half of git and the whole of
`scripts/task` while denying anything that rewrites history, and five agents and
eleven commands of its own. Its skills are third-party and deliberately absent —
[`config/opencode/skills.md`](config/opencode/skills.md) says which, and how to
get them back. Herdr hosts the other three in
persistent sessions that survive sleep, network drops and restarts, and can be
reattached from another machine. Its prefix is `ctrl+a` so it does not fight
tmux, and its panes open `fish -l` so an agent's shell is the same shell you get
everywhere else.

## Notes

- Flags come before the command. `polkadot install --dry-run` silently ignores
  the flag and installs for real.
- `brew bundle check` reports installed-but-outdated as unmet. That is expected
  on a machine that has been running a while; it is only a true failure on a
  fresh one.
- Ghostty is in the Brewfile as a cask. If it was originally installed by hand,
  `brew install --cask ghostty --adopt` makes Homebrew take ownership.
