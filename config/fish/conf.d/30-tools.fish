# ─── tool integrations ────────────────────────────────────────────────────
# Each is guarded, so a machine part-way through a bootstrap still gets a
# working shell instead of a screenful of errors.

# mise — one version manager for node, python, go, everything. Replaces the
# nvm and pyenv evals, which cost ~200ms of shell startup between them.
if type -q mise
    mise activate fish | source
end

# starship — the prompt.
if type -q starship
    starship init fish | source
end

# zoxide — `z` learns where you go. `zi` picks interactively.
if type -q zoxide
    zoxide init fish --cmd z | source
end

# atuin — searchable, durable shell history. --disable-up-arrow keeps plain
# up-arrow as fish's own prefix search, which is better for the last command;
# ctrl-r is where the full search lives.
if type -q atuin
    atuin init fish --disable-up-arrow | source
end

# fzf — ctrl-t files, ctrl-r is atuin's, alt-c directories.
if type -q fzf
    fzf --fish | source
end

# direnv — per-project environments.
if type -q direnv
    direnv hook fish | source
end

# 1Password shell plugins, if the CLI is set up.
if test -f $HOME/.config/op/plugins.sh
    # The plugin file is bash; fish reads the aliases it defines via a shim.
    for line in (grep '^alias ' $HOME/.config/op/plugins.sh 2>/dev/null)
        set -l name (string replace -r '^alias ([^=]+)=.*' '$1' -- $line)
        set -l body (string replace -r '^alias [^=]+="?([^"]*)"?$' '$1' -- $line)
        test -n "$name" -a -n "$body"; and alias $name="$body"
    end
end
