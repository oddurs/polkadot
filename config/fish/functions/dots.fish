# The dotfiles, from anywhere.
function dots --description 'Jump to the polkadot checkout'
    cd ~/Code/polkadot
    test (count $argv) -gt 0; and $EDITOR $argv
end
