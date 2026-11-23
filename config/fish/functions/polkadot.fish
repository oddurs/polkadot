function polkadot --description 'Run the dotfiles installer from anywhere'
    if set -q POLKADOT_ROOT
        command go run "$POLKADOT_ROOT" $argv
        return
    end

    set -l source_dir (path dirname (status filename))
    set -l repo (path resolve "$source_dir/../../..")
    command go run "$repo" $argv
end
