# ─── abbreviations ────────────────────────────────────────────────────────
# Abbreviations, not aliases: fish expands them in place as you type, so what
# runs is what you can see, and history stays greppable for the real command.
#
# Three kinds here:
#   plain          gs → git status
#   --set-cursor   gcm → git commit -m "▮"   (cursor lands inside the quotes)
#   --function     gpo asks git for the current branch at expansion time

# Modern replacements for the tools they replace.
abbr -a cat 'bat --paging=never'
abbr -a ls 'eza --icons --group-directories-first'
abbr -a ll 'eza --icons --group-directories-first -la --git'
abbr -a lt 'eza --icons --tree --level=2 --git-ignore'
abbr -a tree 'eza --icons --tree --git-ignore'
abbr -a top 'btop'
abbr -a ports 'lsof -iTCP -sTCP:LISTEN -P -n'

# git
abbr -a g git
abbr -a gs 'git status --short --branch'
abbr -a ga 'git add'
abbr -a gaa 'git add --all'
abbr -a gc 'git commit'
abbr -a gcm --set-cursor 'git commit -m "%"'
abbr -a gcam --set-cursor 'git commit -am "%"'
abbr -a gca 'git commit --amend --no-edit'
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
abbr -a gstp 'git stash pop'
abbr -a grh 'git reset --hard'
abbr -a lg lazygit
abbr -a ld lazydocker

# Resolved when you press space, not when this file was written.
function __abbr_push_origin
    echo "git push -u origin "(git branch --show-current 2>/dev/null)
end
abbr -a gpo --function __abbr_push_origin

function __abbr_git_root
    echo "cd "(git rev-parse --show-toplevel 2>/dev/null; or echo .)
end
abbr -a groot --function __abbr_git_root

# Expand anywhere on the line, not just in command position.
abbr -a --position anywhere -- '!!' --function __abbr_last_command
function __abbr_last_command
    echo $history[1]
end

# node / package managers
abbr -a n 'npm'
abbr -a ni 'npm install'
abbr -a nr 'npm run'
abbr -a nrd 'npm run dev'
abbr -a nrb 'npm run build'
abbr -a bi 'bun install'
abbr -a br 'bun run'
abbr -a px 'pnpm exec'

# docker, via orbstack
abbr -a d docker
abbr -a dc 'docker compose'
abbr -a dcu 'docker compose up -d'
abbr -a dcd 'docker compose down'
abbr -a dcl 'docker compose logs -f'
abbr -a dps 'docker ps --format "table {{.Names}}\t{{.Image}}\t{{.Status}}"'

# navigation
abbr -a .. 'cd ..'
abbr -a ... 'cd ../..'
abbr -a .... 'cd ../../..'
abbr -a c 'cd ~/Code'

# editors and agents
abbr -a v nvim
abbr -a cl claude
abbr -a oc opencode
abbr -a hd herdr

# mise
abbr -a mi 'mise install'
abbr -a mu 'mise use'
abbr -a ml 'mise ls'

# Safety — aliases rather than abbreviations because the point is that they
# apply even when you are not looking.
alias rm 'rm -i'
alias mv 'mv -i'
alias cp 'cp -i'
