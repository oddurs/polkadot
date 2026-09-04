# ─── path ─────────────────────────────────────────────────────────────────
# `brew shellenv` costs ~51ms and prints the same six static lines every time.
# They are inlined here instead. If Homebrew ever moves, run
#   brew shellenv fish
# and paste the result back.

set -gx HOMEBREW_PREFIX /opt/homebrew
set -gx HOMEBREW_CELLAR /opt/homebrew/Cellar
set -gx HOMEBREW_REPOSITORY /opt/homebrew
fish_add_path -g --move /opt/homebrew/bin /opt/homebrew/sbin
not contains /opt/homebrew/share/info $INFOPATH; and set -gx INFOPATH /opt/homebrew/share/info $INFOPATH

# The rest, in priority order. fish_add_path prepends, so the last call wins.
fish_add_path -g $HOME/.cargo/bin
fish_add_path -g $HOME/.bun/bin
fish_add_path -g $HOME/go/bin
fish_add_path -g $HOME/.local/bin
