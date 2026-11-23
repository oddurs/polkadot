# polkadot

A practical, personal configuration for a Mac and for a Debian server. It is
safe to fork and pleasant to live with. It bootstraps a clean machine, links
declarative settings, and keeps credentials and machine state where they
belong: off the Git history.

```sh
git clone https://github.com/oddurs/polkadot.git ~/Code/polkadot
cd ~/Code/polkadot
./bootstrap.sh
```

On macOS the bootstrap installs the Xcode Command Line Tools, Homebrew, and Go
when needed. On Debian or Ubuntu it installs Go and git with apt. It then runs
`go run . install`.

The two systems share one configuration. A Mac is a workstation: it also gets
the terminal, the fonts and the desktop apps. A Linux machine is a server: it
also gets lingering, a raised inotify limit and a running herdr-core.

Coding-agent installers are opt-in:

```sh
./bootstrap.sh --with-agents
```

## Make It Yours

Fork Polkadot, then edit its tracked configuration. Keep settings that differ
between machines in ignored local layers:

```sh
cp config/fish/conf.d/90-local.example.fish config/fish/conf.d/90-local.fish
cp home/gitconfig.local.example home/gitconfig.local
go run . link
```

No Git identity, credentials, OAuth state, histories, caches, or absolute
machine paths are tracked. A fork starts as yours, rather than as a copy of
someone else's account.

## Commands

| Command | Purpose |
|---|---|
| `polkadot install` | Install packages, link configuration, set up shell, tools, theme, and timer |
| `polkadot link` | Link configuration only |
| `polkadot unlink` | Remove only symlinks owned by this checkout |
| `polkadot packages` | Install `packages/darwin.Brewfile` on macOS; `packages/linux.apt` and mise on Linux |
| `polkadot repos` | Clone the managed repositories, such as mcdir |
| `polkadot runtimes` | Install the pinned tools that `config/mise/config.toml` lists |
| `polkadot version` | Print the tag and commit of this checkout |
| `polkadot system` | Linux: enable lingering, raise the inotify limit, start herdr-core |
| `polkadot shell` | Configure Fish as the login shell |
| `polkadot tools` | Check command-line tools |
| `polkadot agents` | Explicitly install coding agents |
| `polkadot theme` | Link the bundled London Moquette theme |
| `polkadot sync` | Fast-forward managed repositories and relink changes |
| `polkadot timer` | Install the weekly `sync` timer: launchd on macOS, systemd on Linux |
| `polkadot doctor` | Report configuration health without changing anything |

Put flags before the command: `polkadot --dry-run install`. `--only=<step>`
selects one install step; `--with-agents` enables the agent step; `--yes`
approves external installers for intentional unattended runs.

## Guarantees

- Existing files are backed up to `~/.polkadot-backup/<timestamp>/` before a
  link replaces them.
- `unlink` removes only links pointing into this checkout.
- `sync` is fast-forward only and skips dirty trees and non-default branches.
- Nothing writes macOS `defaults`; personal system preferences remain yours.
- On Linux, `polkadot system` makes the only changes outside `$HOME`: one
  file in `/etc/sysctl.d` and lingering for your user.

## Documentation

[Operations guide](docs/OPERATIONS.md) covers first-machine setup,
customization, recovery, and scheduled sync. [Configuration reference](docs/REFERENCE.md)
describes the tracked tools, theme, shell, editor, runtimes, and agents.

London Moquette is bundled as a small runtime subset of
[Subway Seat](https://github.com/oddurs/subway-seat); its upstream MIT license
is preserved at [`theme/SUBWAY-SEAT-LICENSE`](theme/SUBWAY-SEAT-LICENSE).

Polkadot is MIT licensed. See [SECURITY.md](SECURITY.md) to report a
vulnerability.
