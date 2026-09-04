# ─── abbreviations ────────────────────────────────────────────────────────
# Abbreviations, not aliases: fish expands them in place as you type, so what
# runs is what you can see, and history stays greppable for the real command.

# Modern replacements for the tools they replace.
abbr -a cat 'bat --paging=never'
abbr -a ls 'eza --icons --group-directories-first'
abbr -a ll 'eza --icons --group-directories-first -la --git'
abbr -a lt 'eza --icons --tree --level=2 --git-ignore'
abbr -a tree 'eza --icons --tree --git-ignore'
abbr -a du 'dust'
abbr -a top 'btop'

# git — the ones worth not typing.
abbr -a g git
abbr -a gs 'git status --short --branch'
abbr -a ga 'git add'
abbr -a gaa 'git add --all'
abbr -a gc 'git commit'
abbr -a gca 'git commit --amend'
abbr -a gco 'git checkout'
abbr -a gcb 'git checkout -b'
abbr -a gd 'git diff'
abbr -a gds 'git diff --staged'
abbr -a gl 'git log --oneline --graph --decorate -20'
abbr -a gp 'git push'
abbr -a gpf 'git push --force-with-lease'
abbr -a gpl 'git pull --rebase'
abbr -a gsw 'git switch'
abbr -a gst 'git stash'
abbr -a lg lazygit
abbr -a ld lazydocker

# navigation
abbr -a .. 'cd ..'
abbr -a ... 'cd ../..'
abbr -a .... 'cd ../../..'
abbr -a c 'cd ~/Code'

# editors and agents
abbr -a v nvim
abbr -a n nvim
abbr -a cl claude
abbr -a oc opencode

# safety — these are aliases rather than abbreviations because the point is
# that they apply even when you are not looking.
alias rm 'rm -i'
alias mv 'mv -i'
alias cp 'cp -i'
