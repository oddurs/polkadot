# polkadot

A practical, personal macOS configuration that is safe to fork and pleasant to
live with. It bootstraps a clean Mac, links declarative settings, and keeps
credentials and machine state where they belong: off the Git history.

```sh
git clone https://github.com/oddurs/polkadot.git ~/Code/polkadot
cd ~/Code/polkadot
./bootstrap.sh
```

The bootstrap installs the Xcode Command Line Tools, Homebrew, and Go when
needed, then runs `go run . install`. It is deliberately macOS-only.

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
| `polkadot brew` | Install Homebrew packages from `Brewfile` |
| `polkadot shell` | Configure Fish as the login shell |
| `polkadot tools` | Check command-line tools |
| `polkadot agents` | Explicitly install coding agents |
| `polkadot theme` | Link the bundled London Moquette theme |
| `polkadot sync` | Fast-forward managed repositories and relink changes |
| `polkadot timer` | Install the weekly `sync` launchd agent |
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

## Documentation

[Operations guide](docs/OPERATIONS.md) covers first-machine setup,
customization, recovery, and scheduled sync. [Configuration reference](docs/REFERENCE.md)
describes the tracked tools, theme, shell, editor, runtimes, and agents.

London Moquette is bundled as a small runtime subset of
[Subway Seat](https://github.com/oddurs/subway-seat); its upstream MIT license
is preserved at [`theme/SUBWAY-SEAT-LICENSE`](theme/SUBWAY-SEAT-LICENSE).

Polkadot is MIT licensed. See [CONTRIBUTING.md](CONTRIBUTING.md),
[SECURITY.md](SECURITY.md), and [SUPPORT.md](SUPPORT.md) for project guidance.
