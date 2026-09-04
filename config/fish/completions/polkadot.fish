# Completions for the dotfiles installer.
complete -c polkadot -f
complete -c polkadot -n __fish_use_subcommand -a install -d 'Everything: brew, links, shell, agents'
complete -c polkadot -n __fish_use_subcommand -a link    -d 'Symlinks only'
complete -c polkadot -n __fish_use_subcommand -a brew    -d 'Homebrew and the Brewfile only'
complete -c polkadot -n __fish_use_subcommand -a doctor  -d 'Report state, change nothing'
complete -c polkadot -l dry-run -d 'Show the plan, write nothing'
complete -c polkadot -l only -x -a 'brew link shell tools agents' -d 'Run one step'
