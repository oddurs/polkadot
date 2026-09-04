# ─── environment ──────────────────────────────────────────────────────────

set -gx EDITOR nvim
set -gx VISUAL nvim
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

# bat's base16 theme reads the terminal's own 16 colours, so it is Gotham
# wherever Gotham is loaded and follows the terminal if that changes.
set -gx BAT_THEME base16

# ripgrep respects a global config file for the flags worth always having.
set -gx RIPGREP_CONFIG_PATH $HOME/.config/ripgrep/config

# Homebrew: don't phone home, don't auto-update mid-command.
set -gx HOMEBREW_NO_ANALYTICS 1
set -gx HOMEBREW_NO_ENV_HINTS 1
set -gx HOMEBREW_NO_AUTO_UPDATE 1
