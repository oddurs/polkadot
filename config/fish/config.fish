# ─── fish ─────────────────────────────────────────────────────────────────
# Managed by polkadot: ~/Code/polkadot/config/fish/config.fish
#
# Almost nothing lives here. fish sources every file in conf.d/ automatically,
# in name order, for both interactive and non-interactive shells — so PATH and
# environment go there and stay correct under `fish -c`, editors and launchd.
# This file is only for things that genuinely need an interactive terminal.

if not status is-interactive
    exit
end

# Prezto's number-jump directory stack, which is the one habit worth keeping.
for i in (seq 1 9)
    alias $i="cd +$i"
end

# fish's own greeting is noise.
set -g fish_greeting
