# ─── environment ──────────────────────────────────────────────────────────

set -gx EDITOR fresh
set -gx VISUAL fresh
set -gx PAGER less
set -gx LESS '-R --mouse'
set -gx MANPAGER 'sh -c "col -bx | bat -l man -p"'
set -gx MANROFFOPT '-c'

# XDG. macOS apps that ask get the standard answer, which is how lazygit and
# friends end up reading from ~/.config instead of ~/Library/Application Support.
set -gx XDG_CONFIG_HOME $HOME/.config
set -gx XDG_DATA_HOME $HOME/.local/share
set -gx XDG_CACHE_HOME $HOME/.cache
set -gx XDG_STATE_HOME $HOME/.local/state

# lazygit merges a comma-separated list left to right, so ours comes first and
# Subway Seat's generated theme layers on top. The flavor follows themeFlavor
# in steps.go; `polkadot doctor` says so if it stops.
set -gx LG_CONFIG_FILE $HOME/.config/lazygit/config.yml,$HOME/.config/lazygit/paris-guimard.yml

# ripgrep respects a global config file for the flags worth always having.
set -gx RIPGREP_CONFIG_PATH $HOME/.config/ripgrep/config

# Homebrew: don't phone home, don't auto-update mid-command.
set -gx HOMEBREW_NO_ANALYTICS 1
set -gx HOMEBREW_NO_ENV_HINTS 1
set -gx HOMEBREW_NO_AUTO_UPDATE 1
