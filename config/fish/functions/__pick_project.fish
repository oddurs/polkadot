# Jump to any project under ~/Code, ranked by how recently it was touched.
function __pick_project --description 'fzf a project in ~/Code'
    set -l dir (
        fd --type d --max-depth 1 . ~/Code 2>/dev/null \
        | fzf --prompt='project ❯ ' --height=40% \
              --preview 'eza --icons --git-ignore --tree --level=2 --color=always {} 2>/dev/null | head -40' \
              --preview-window=right:55%:border-left
    )
    if test -n "$dir"
        cd $dir
        commandline -f repaint
    end
    commandline -f repaint
end
