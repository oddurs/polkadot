# Open a file in $EDITOR, previewed with syntax highlighting.
function __pick_file --description 'fzf a file and open it'
    set -l file (
        fd --type f --hidden --exclude .git 2>/dev/null \
        | fzf --prompt='open ❯ ' --height=60% \
              --preview 'bat --style=numbers --color=always --line-range=:120 {}' \
              --preview-window=right:60%:border-left
    )
    test -n "$file"; and $EDITOR $file
    commandline -f repaint
end
