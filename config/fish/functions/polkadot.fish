function polkadot --description 'Run the dotfiles installer from anywhere'
    set -l repo $POLKADOT_ROOT
    if test -z "$repo"
        set -l source_dir (path dirname (status filename))
        set repo (path resolve "$source_dir/../../..")
    end

    # -C, not `go run <path>`: Go builds a package only from inside its
    # module, so a path alone fails from any other folder.
    command go -C "$repo" run . $argv
end
