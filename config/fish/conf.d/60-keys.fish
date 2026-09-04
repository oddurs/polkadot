# ─── keys ─────────────────────────────────────────────────────────────────
# Only bindings fish leaves free, so nothing standard is shadowed.
#
#   ctrl-z   background job ⇄ foreground, both ways
#   ctrl-g   jump to a project in ~/Code
#   ctrl-o   open a file in the editor
#   alt-g    switch git branch
#   alt-s    prefix the line with sudo
#   alt-e    edit the line in $EDITOR   (fish's own, kept)
#   ctrl-t   insert a file path          (fzf's)
#   alt-c    cd into a directory         (fzf's)
#   ctrl-r   search history              (atuin's)

function fish_user_key_bindings
    # ctrl-z suspends into the background; pressing it again on an empty line
    # brings the job back. The single most useful missing binding in any shell.
    bind ctrl-z __toggle_fg

    bind ctrl-g __pick_project
    bind ctrl-o __pick_file
    bind alt-g __pick_branch
    bind alt-s __prefix_sudo
end

function __toggle_fg
    if jobs -q
        commandline -r ''
        commandline -f repaint
        fg >/dev/null 2>&1
        commandline -f repaint
    else
        commandline -f repaint
    end
end

function __prefix_sudo
    set -l cmd (commandline)
    if string match -qr '^sudo ' -- $cmd
        commandline -r (string replace -r '^sudo ' '' -- $cmd)
    else
        commandline -r "sudo $cmd"
    end
    commandline -f repaint
end
