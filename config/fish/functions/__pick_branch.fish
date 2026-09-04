# Switch branch, with the branch's recent commits in the preview so you can
# tell two similarly-named ones apart.
function __pick_branch --description 'fzf a git branch and switch to it'
    git rev-parse --git-dir >/dev/null 2>&1; or begin
        commandline -f repaint
        return
    end
    set -l branch (
        git branch --all --sort=-committerdate \
            --format='%(refname:short)' 2>/dev/null \
        | string replace -r '^origin/' '' | sort -u | string match -v HEAD \
        | fzf --prompt='branch ❯ ' --height=45% \
              --preview 'git log --oneline --graph --color=always -20 {}' \
              --preview-window=right:60%:border-left
    )
    test -n "$branch"; and git switch $branch 2>/dev/null; or test -n "$branch"; and git switch -c $branch
    commandline -f repaint
end
