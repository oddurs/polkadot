# ─── tool integrations ────────────────────────────────────────────────────
# Everything here is guarded twice: on the binary existing, so a half-finished
# bootstrap still gets a working shell; and on the shell being interactive,
# because prompt hooks, history search and key bindings are worth ~110ms and
# do nothing at all inside `fish -c`.

# --- every shell, interactive or not ---

# mise resolves the right node/python for the directory, which scripts and
# editors need as much as a prompt does. The shims path is the cheap half and
# is enough on its own for non-interactive use.
if test -d $HOME/.local/share/mise/shims
    fish_add_path -g -p $HOME/.local/share/mise/shims
end

# --- interactive only ---
# A block rather than an early `exit`: inside conf.d, `exit` terminates the
# whole shell, not just this file.
if status is-interactive

    # The full activation adds the cd hook, so versions switch as you move.
    type -q mise; and mise activate fish | source

    type -q starship; and starship init fish | source
    type -q zoxide; and zoxide init fish --cmd z | source

    # --disable-up-arrow keeps plain up-arrow as fish's own prefix search, which is
    # better for the last command; ctrl-r is where the full search lives.
    type -q atuin; and atuin init fish --disable-up-arrow | source

    type -q fzf; and fzf --fish | source
    type -q direnv; and direnv hook fish | source

    # 1Password shell plugins, if the CLI is set up.
    if test -f $HOME/.config/op/plugins.sh
        for line in (grep '^alias ' $HOME/.config/op/plugins.sh 2>/dev/null)
            set -l name (string replace -r '^alias ([^=]+)=.*' '$1' -- $line)
            set -l body (string replace -r '^alias [^=]+="?([^"]*)"?$' '$1' -- $line)
            test -n "$name" -a -n "$body"; and alias $name="$body"
        end
end
end
