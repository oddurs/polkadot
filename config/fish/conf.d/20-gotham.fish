# ─── gotham ───────────────────────────────────────────────────────────────
# fish's syntax highlighting, in Andrea Leopardi's palette.
#
# set -g rather than set -U on purpose: universal variables persist in
# fish_variables and would quietly outlive this file, which makes theme changes
# mysterious. Everything here is owned by the repo.

set -g fish_color_normal          98d1ce
set -g fish_color_command         599caa
set -g fish_color_keyword         d26939
set -g fish_color_quote           edb54b
set -g fish_color_redirection     888ba5
set -g fish_color_end             d26939
set -g fish_color_error           c33027
set -g fish_color_param           98d1ce
set -g fish_color_comment         4e5165
set -g fish_color_operator        26a98b
set -g fish_color_escape          26a98b
set -g fish_color_autosuggestion  245361
set -g fish_color_option          33859d
set -g fish_color_valid_path      --underline
set -g fish_color_selection       --background=093748
set -g fish_color_search_match    --background=093748
set -g fish_color_history_current --bold
set -g fish_color_cancel          c33027

set -g fish_pager_color_progress            888ba5
set -g fish_pager_color_prefix              33859d --bold
set -g fish_pager_color_completion          98d1ce
set -g fish_pager_color_description         4e5165
set -g fish_pager_color_selected_background --background=093748

# fzf, same palette.
set -gx FZF_DEFAULT_OPTS "\
--height=40% --layout=reverse --border=none --info=inline \
--color=bg+:#093748,bg:-1,spinner:#26a98b,hl:#d26939 \
--color=fg:#98d1ce,header:#4e5165,info:#245361,pointer:#26a98b \
--color=marker:#edb54b,fg+:#d3ebe9,prompt:#33859d,hl+:#d26939 \
--color=border:#245361"

# Use fd for fzf so it honours .gitignore and is an order of magnitude faster.
if type -q fd
    set -gx FZF_DEFAULT_COMMAND 'fd --type f --hidden --follow --exclude .git'
    set -gx FZF_CTRL_T_COMMAND $FZF_DEFAULT_COMMAND
    set -gx FZF_ALT_C_COMMAND 'fd --type d --hidden --follow --exclude .git'
end
