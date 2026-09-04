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
| `polkadot doctor` | report what is and isn't in place, change nothing |

Flags: `--dry-run` prints the plan without writing anything, `--only=<step>`
runs one part (`brew`, `link`, `shell`, `tools`, `agents`).

## What it will not do

**It never writes macOS `defaults`.** No Dock tweaks, no key repeat rates, no
Finder settings. Those are personal, they drift, and reverting them is worse
than setting them. System Preferences stays yours.

**It never deletes a config it did not write.** Anything already sitting at a
target path is moved to `~/.polkadot-backup/<timestamp>/`, keeping its path
shape, so restoring is a plain `mv` back.

## Layout

```
config/     → ~/.config/*        fish, starship, nvim, lazygit, ghostty,
                                 herdr, atuin, mise, btop, bat, ripgrep
home/       → ~/*                gitconfig, gitignore_global, commit template
config/vscode/settings.json      → ~/Library/Application Support/Code/User/
Brewfile                         formulae, casks, VS Code extensions
```

Everything is symlinked, not copied. Edit in the repo, it is live immediately;
edit in place, it is already staged.

## Gotham

[Andrea Leopardi's palette](https://github.com/whatyouhide/vim-gotham), carried
across the whole stack — Ghostty, Neovim, fish, starship, lazygit, delta, bat,
fzf, btop, herdr and VS Code.

Three of those get it for free. `bat`, `delta` and `herdr` are set to render
through the terminal's own sixteen colours (`base16`, `theme = "terminal"`),
so they are Gotham because Ghostty is, and they follow automatically if the
terminal theme ever changes. The rest carry the hex values, which live in one
place: [`config/gotham/palette.md`](config/gotham/palette.md).

The previous theme, Halide, is still here. Its twelve computed Ghostty themes
are in `config/ghostty/themes/`, its Neovim colorschemes in
`config/nvim/colors/`, and its VS Code extension is still installed. Switching
back is one line in `config/ghostty/config` and one in
`config/nvim/lua/plugins/gotham.lua`.

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

## Runtimes

`mise` replaces nvm and pyenv. It reads `.tool-versions`, `.nvmrc`,
`.node-version` and `mise.toml`, so per-project pins keep working, and it costs
a few milliseconds at shell start instead of a few hundred.

## Agents

`claude`, `codex`, `opencode` and `herdr`. Herdr hosts the other three in
persistent sessions that survive sleep, network drops and restarts, and can be
reattached from another machine. Its prefix is `ctrl+a` so it does not fight
tmux, and its panes open `fish -l` so an agent's shell is the same shell you get
everywhere else.

## Notes

- `brew bundle check` reports installed-but-outdated as unmet. That is expected
  on a machine that has been running a while; it is only a true failure on a
  fresh one.
- Ghostty is in the Brewfile as a cask. If it was originally installed by hand,
  `brew install --cask ghostty --adopt` makes Homebrew take ownership.
