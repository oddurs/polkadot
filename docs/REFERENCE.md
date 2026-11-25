# Configuration Reference

This page records the choices Polkadot makes. It is deliberately descriptive:
the operational contract lives in [OPERATIONS.md](OPERATIONS.md).

## Layout

| Source | Destination | Purpose |
|---|---|---|
| `config/` | `~/.config/` | Shell, terminal, CLI, agent, and editor settings |
| `home/` | `~/` | Git, ignore rules, and commit template |
| `theme/` | `~/.config/polkadot-theme` | London Moquette runtime assets |
| `agents/AGENTS.md` | each agent's instruction path | One instruction file for Claude, Codex, OpenCode, and Antigravity |
| `packages/darwin.Brewfile` | Homebrew | macOS: the shell, git, desktop apps |
| `packages/linux.apt` | apt | Linux: the shell, git, build tools, and rsync |
| `config/mise/config.toml` | mise | Every machine: the pinned command-line tools |

## Linked and seeded

Polkadot puts a file on a machine in one of two ways.

- **Linked.** The path on the machine is a symbolic link into the checkout. A
  change in the repository reaches every machine at its next sync. Polkadot
  links a file only when no tool writes to it.
- **Seeded.** Polkadot copies the file once, when the machine has none. The
  machine owns the file from then on. `~/.claude/settings.json`,
  `~/.config/gh/config.yml`, `~/.codex/config.toml`,
  `~/.config/herdr/config.toml` and the two Antigravity files under
  `~/.gemini/` are seeded, because their tools rewrite them.

A change to a seeded file reaches a new machine only. To change it on a
machine that has the file, edit the file there.

The tracked configuration is instructions, not state. Histories, OAuth files,
plugin caches, local Git identity, and per-machine paths stay out of the repo.

## Theme

London Moquette is the single bundled theme. It supplies Ghostty, Fish, fzf,
eza, lazygit, delta, Codex, and OpenCode.
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
| `20-fzf.fish` | Fish and fzf presentation defaults |
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
desktop install through the Brewfile. Git uses delta, a commit
template, SSH signing, a conservative rebase workflow, and an ignored local
override for identity or work-specific settings.

## Agents

Claude, Codex, OpenCode, and Antigravity read one instruction file,
`agents/AGENTS.md`. Polkadot links it to each agent's own path:
`~/.claude/CLAUDE.md`, `~/.codex/AGENTS.md`, `~/.config/opencode/AGENTS.md`,
and `~/.gemini/AGENTS.md`. A rule changed there changes for every agent.

Settings that only one agent understands stay in that agent's configuration.
Claude keeps its path rules and skills. OpenCode loads `STACK.md` after the
shared file. Histories, auth state, downloaded plugins, and local approvals
remain on the machine. Agent executables are installed only through an
explicit agent command or `--with-agents` during installation.

### The same posture in four agents

Each agent may read and build without asking, asks before it changes history
or leaves the workspace, and is refused the commands that skip a hook or
force a push.

| Agent | Default model | How it is held |
|---|---|---|
| Claude Code | `claude-fable-5-2[1m]`, high effort | Auto mode, with the allow and deny lists in `settings.json` |
| Codex | `gpt-6-astra`, high reasoning | `on-request` approval inside a `workspace-write` sandbox with network |
| OpenCode | `openrouter/sonnet`; `openrouter/haiku` for titles and summaries | `ask` by default, with an allowlist of read-only and build commands |
| Antigravity | The CLI's own default | `accept-edits` in the terminal sandbox, with the same allow and deny lists |

OpenCode reaches its models through OpenRouter and needs a key in the
environment. `openrouter-key` asks for the key, checks it with OpenRouter and
saves it in `conf.d/90-local.fish` on that machine. Run it once per machine.

## Packages

Three lists, one owner for each tool:

| List | Owner | Holds |
|---|---|---|
| `config/mise/config.toml` | mise, every machine | The command-line tools, each pinned to one version |
| `packages/darwin.Brewfile` | Homebrew, macOS | fish, git, Go, mise, desktop apps, and three tools with no macOS build in mise |
| `packages/linux.apt` | apt, Linux | fish, git, build tools, and rsync |

To update a tool, change its version in `config/mise/config.toml` in a pull
request. Each machine installs the new version at its next `sync` and
`runtimes`. A tool that one project needs belongs in that project's
`mise.toml`, not here.

## Runtimes

mise owns language runtimes and respects existing project files such as
`.tool-versions`, `.nvmrc`, `.node-version`, and `mise.toml`. It replaces the
need for global nvm and pyenv state.
