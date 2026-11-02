# Configuration Reference

This page records the choices Polkadot makes. It is deliberately descriptive:
the operational contract lives in [OPERATIONS.md](OPERATIONS.md).

## Layout

| Source | Destination | Purpose |
|---|---|---|
| `config/` | `~/.config/` | Shell, terminal, CLI, agent, and editor settings |
| `home/` | `~/` | Git, Zsh, ignore rules, and commit template |
| `theme/` | `~/.config/polkadot-theme` | London Moquette runtime assets |
| `Brewfile` | Homebrew | Formulae, casks, and VS Code extensions |

The tracked configuration is instructions, not state. Histories, OAuth files,
plugin caches, local Git identity, and per-machine paths stay out of the repo.

## Theme

London Moquette is the single bundled theme. It supplies Ghostty, Fish, fzf,
eza, lazygit, delta, tmux, Codex, OpenCode, and zsh syntax highlighting.
Starship and Herdr carry their palette directly in tracked configuration; bat
uses the terminal's Base16 colors. `polkadot theme` restores the discovery
links required by Ghostty, Codex, and OpenCode.

The full Subway Seat generator remains separate. Polkadot vendors only the
small, stable London Moquette runtime surface it actually uses.

## Shell

Fish is the login shell and Starship is the prompt. `conf.d/` loads in this
order:

| File | Responsibility |
|---|---|
| `00-path.fish` | Homebrew and executable search path |
| `10-env.fish` | Editor, XDG, pager, and tool defaults |
| `20-gotham.fish` | Fish and fzf presentation defaults |
| `30-tools.fish` | mise, zoxide, atuin, fzf, direnv, 1Password, lazygit |
| `40-abbr.fish` | Everyday command abbreviations and safety aliases |
| `50-prompt.fish` | Starship initialization |
| `60-keys.fish` | Key bindings |
| `70-london-moquette.fish` | London Moquette selection and theme paths |

Abbreviations expand into visible commands; destructive `rm`, `mv`, and `cp`
remain aliases because their safety behavior must apply even when typed from
muscle memory.

## Editor and Git

Fresh is `$EDITOR` and `$VISUAL`. VS Code has tracked user settings and a
small managed extension set through the Brewfile. Git uses delta, a commit
template, SSH signing, a conservative rebase workflow, and an ignored local
override for identity or work-specific settings.

## Agents

Claude, Codex, OpenCode, and Herdr each keep declarative instructions in this
repository. Their histories, auth state, downloaded plugins, and local
approvals remain on the machine. Agent executables are installed only through
an explicit agent command or `--with-agents` during installation.

## Runtimes

mise owns language runtimes and respects existing project files such as
`.tool-versions`, `.nvmrc`, `.node-version`, and `mise.toml`. It replaces the
need for global nvm and pyenv state.
