# Homebrew installs what must be part of the system on a Mac: the shell, git,
# the desktop apps, and the few tools that have no macOS build in mise. Every
# other command-line tool is pinned in config/mise/config.toml, so each
# machine runs the same version.
#
# Do not regenerate this file with `brew bundle dump`. A dump lists every
# formula on the machine. This list is the foundation only: a tool that one
# project needs belongs in that project.

# ─── system ───────────────────────────────────────────────────────────────
brew "fish"
brew "git"
# Runs polkadot itself
brew "go"
brew "mise"

# ─── no macOS build in mise ───────────────────────────────────────────────
brew "btop"
brew "difftastic"
brew "tokei"

# ─── workstation ──────────────────────────────────────────────────────────
# Monitor a directory for changes and run a shell command
brew "fswatch"
# Continuous file synchronization
brew "syncthing", restart_service: :changed
# Unarchiving tools for many formats
brew "unar"

# ─── casks ────────────────────────────────────────────────────────────────
cask "1password-cli"
# Agent orchestration: the desktop app. The agents step installs its CLI.
cask "antigravity"
# The typeface everything is set in
cask "font-jetbrains-mono-nerd-font"
# The terminal
cask "ghostty"
# Replacement for Docker Desktop
cask "orbstack"
cask "visual-studio-code"
