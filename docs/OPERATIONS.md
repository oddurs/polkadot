# Operations Guide

Polkadot is a macOS setup you can run repeatedly. The repository is the source
of truth for configuration; each machine keeps its own credentials, histories,
caches, and local overrides.

For the tool and layout choices, see the [configuration reference](REFERENCE.md).

## First machine

Clone your fork with HTTPS, then run the bootstrap script. It opens the Xcode
Command Line Tools installer when necessary, installs Homebrew and Go, and
starts the normal install. It does not install coding agents unless you ask.

```sh
git clone https://github.com/YOU/polkadot.git ~/Code/polkadot
cd ~/Code/polkadot
./bootstrap.sh
go run . doctor
```

To install the configured coding agents too, opt in explicitly:

```sh
./bootstrap.sh --with-agents
```

Agent installation may execute third-party installers. Read the commands in
`steps.go` first, or install an individual agent with `go run . agents`.
Homebrew and agent installers display their exact command and require a `y`
confirmation. For intentionally non-interactive automation, add `--yes` to
the Go command, or set `POLKADOT_BOOTSTRAP_YES=1` for `bootstrap.sh`.

## Normal use

Run `polkadot doctor` after changing the repository or moving to a new
computer. It reports link and executable health without changing anything.

```sh
polkadot doctor
polkadot --dry-run install
polkadot --only=link install
polkadot sync
```

Flags always precede the command. Invalid `--only` values and `--only` with a
non-install command fail instead of silently doing nothing.

`sync` only fast-forwards clean repositories that are checked out on their
remote default branch. It skips feature branches, diverged repositories, dirty
trees, and repositories without a known default branch. `polkadot timer`
installs the weekly version of that same operation; its log is at
`~/Library/Logs/polkadot-sync.log`.

## Customization

Tracked files are shared defaults. Keep per-machine settings in files ignored
by Git:

```sh
cp config/fish/conf.d/90-local.example.fish config/fish/conf.d/90-local.fish
cp home/gitconfig.local.example home/gitconfig.local
polkadot link
```

`home/gitconfig.local` is linked to `~/.gitconfig.local` and included after
the shared Git config, so it is the correct place for your identity and work
overrides. The repository does not set a Git name or email for fork users.

The Fish wrapper discovers the checkout containing its own source. If you run
the binary from another location, set `POLKADOT_ROOT` to that checkout.

Never put API keys, OAuth tokens, SSH keys, session files, plugin caches, or
absolute machine paths in tracked configuration.

## Theme

London Moquette is bundled under `theme/london-moquette`; it is not cloned or
generated during installation. `polkadot theme` links the assets that need
app-specific discovery paths and creates `~/.config/polkadot-theme` for the
shared files used by fish, fzf, eza, lazygit, delta, tmux, and zsh.

The bundle is a focused runtime subset of Subway Seat, with the upstream MIT
license preserved at `theme/SUBWAY-SEAT-LICENSE`. Update it deliberately from
the source project when the palette changes; do not add the full generator,
website, or unrelated app ports here.

## Recovery and removal

Before replacing a file, `polkadot link` moves the existing target to a dated
directory under `~/.polkadot-backup`. Restore a file by moving it back after
removing the symlink at its original path.

```sh
polkadot --dry-run unlink
polkadot unlink
```

`unlink` removes only symlinks that point to this checkout. It does not delete
regular files, links owned by another checkout, Homebrew packages, application
state, or the repository itself. This is the right command before moving a
checkout or retiring a machine.
