# ─── fzf ──────────────────────────────────────────────────────────────────
# Layout only. Colours come from London Moquette in 70-london-moquette.fish,
# which appends its own --color to these options.

set -gx FZF_DEFAULT_OPTS "--height=40% --layout=reverse --border=none --info=inline"

# Use fd for fzf so it honours .gitignore and is an order of magnitude faster.
if type -q fd
    set -gx FZF_DEFAULT_COMMAND 'fd --type f --hidden --follow --exclude .git'
    set -gx FZF_CTRL_T_COMMAND $FZF_DEFAULT_COMMAND
    set -gx FZF_ALT_C_COMMAND 'fd --type d --hidden --follow --exclude .git'
end
