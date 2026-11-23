# London Moquette runtime bundle for Polkadot. Upstream attribution is in SUBWAY-SEAT-LICENSE.
# London Moquette for zsh-syntax-highlighting. Needs a truecolor terminal.

typeset -gA ZSH_HIGHLIGHT_STYLES

# General
ZSH_HIGHLIGHT_STYLES[default]='fg=#D8DEEA'
ZSH_HIGHLIGHT_STYLES[unknown-token]='fg=#FE8474'
ZSH_HIGHLIGHT_STYLES[comment]='fg=#73819C,italic'

# Commands
ZSH_HIGHLIGHT_STYLES[reserved-word]='fg=#DE8946'
ZSH_HIGHLIGHT_STYLES[alias]='fg=#F2C03F'
ZSH_HIGHLIGHT_STYLES[suffix-alias]='fg=#F2C03F'
ZSH_HIGHLIGHT_STYLES[global-alias]='fg=#F2C03F'
ZSH_HIGHLIGHT_STYLES[builtin]='fg=#F2C03F'
ZSH_HIGHLIGHT_STYLES[function]='fg=#F2C03F'
ZSH_HIGHLIGHT_STYLES[command]='fg=#F2C03F'
ZSH_HIGHLIGHT_STYLES[hashed-command]='fg=#F2C03F'
ZSH_HIGHLIGHT_STYLES[precommand]='fg=#F2C03F,italic'
ZSH_HIGHLIGHT_STYLES[autodirectory]='fg=#F2C03F,italic'
ZSH_HIGHLIGHT_STYLES[arg0]='fg=#F2C03F'

# Separators and redirection
ZSH_HIGHLIGHT_STYLES[commandseparator]='fg=#DE8946'
ZSH_HIGHLIGHT_STYLES[redirection]='fg=#AE9EDC'
ZSH_HIGHLIGHT_STYLES[named-fd]='fg=#AE9EDC'
ZSH_HIGHLIGHT_STYLES[numeric-fd]='fg=#AE9EDC'

# Arguments and options
ZSH_HIGHLIGHT_STYLES[single-hyphen-option]='fg=#54B4B5'
ZSH_HIGHLIGHT_STYLES[double-hyphen-option]='fg=#54B4B5'
ZSH_HIGHLIGHT_STYLES[assign]='fg=#D8DEEA'
ZSH_HIGHLIGHT_STYLES[path]='fg=#C1C9D8,underline'
ZSH_HIGHLIGHT_STYLES[path_pathseparator]='fg=#8F9AB0,underline'
ZSH_HIGHLIGHT_STYLES[path_prefix]='fg=#C1C9D8,underline'
ZSH_HIGHLIGHT_STYLES[path_prefix_pathseparator]='fg=#8F9AB0,underline'
ZSH_HIGHLIGHT_STYLES[globbing]='fg=#AE9EDC'
ZSH_HIGHLIGHT_STYLES[history-expansion]='fg=#AE9EDC'
ZSH_HIGHLIGHT_STYLES[arithmetic-expansion]='fg=#FE8474'

# Strings
ZSH_HIGHLIGHT_STYLES[single-quoted-argument]='fg=#77C581'
ZSH_HIGHLIGHT_STYLES[single-quoted-argument-unclosed]='fg=#FE8474'
ZSH_HIGHLIGHT_STYLES[double-quoted-argument]='fg=#77C581'
ZSH_HIGHLIGHT_STYLES[double-quoted-argument-unclosed]='fg=#FE8474'
ZSH_HIGHLIGHT_STYLES[dollar-quoted-argument]='fg=#77C581'
ZSH_HIGHLIGHT_STYLES[dollar-quoted-argument-unclosed]='fg=#FE8474'
ZSH_HIGHLIGHT_STYLES[rc-quote]='fg=#AE9EDC'
ZSH_HIGHLIGHT_STYLES[dollar-double-quoted-argument]='fg=#AE9EDC'
ZSH_HIGHLIGHT_STYLES[back-double-quoted-argument]='fg=#AE9EDC'
ZSH_HIGHLIGHT_STYLES[back-dollar-quoted-argument]='fg=#AE9EDC'

# Substitutions
ZSH_HIGHLIGHT_STYLES[command-substitution]='none'
ZSH_HIGHLIGHT_STYLES[command-substitution-unquoted]='none'
ZSH_HIGHLIGHT_STYLES[command-substitution-quoted]='fg=#77C581'
ZSH_HIGHLIGHT_STYLES[command-substitution-delimiter]='fg=#AE9EDC'
ZSH_HIGHLIGHT_STYLES[command-substitution-delimiter-unquoted]='fg=#AE9EDC'
ZSH_HIGHLIGHT_STYLES[command-substitution-delimiter-quoted]='fg=#AE9EDC'
ZSH_HIGHLIGHT_STYLES[process-substitution]='none'
ZSH_HIGHLIGHT_STYLES[process-substitution-delimiter]='fg=#AE9EDC'
ZSH_HIGHLIGHT_STYLES[back-quoted-argument]='none'
ZSH_HIGHLIGHT_STYLES[back-quoted-argument-unclosed]='fg=#FE8474'
ZSH_HIGHLIGHT_STYLES[back-quoted-argument-delimiter]='fg=#AE9EDC'

# brackets and cursor highlighters (only used if enabled in ZSH_HIGHLIGHT_HIGHLIGHTERS)
ZSH_HIGHLIGHT_STYLES[bracket-error]='fg=#FE8474,bold'
ZSH_HIGHLIGHT_STYLES[bracket-level-1]='fg=#F2C03F'
ZSH_HIGHLIGHT_STYLES[bracket-level-2]='fg=#DE8946'
ZSH_HIGHLIGHT_STYLES[bracket-level-3]='fg=#77C581'
ZSH_HIGHLIGHT_STYLES[bracket-level-4]='fg=#54B4B5'
ZSH_HIGHLIGHT_STYLES[cursor-matchingbracket]='fg=#FFD36C,bold'
ZSH_HIGHLIGHT_STYLES[cursor]='fg=#121826,bg=#F2C03F'

# zsh-autosuggestions (fish's autosuggestion color)
ZSH_AUTOSUGGEST_HIGHLIGHT_STYLE='fg=#73819C'
