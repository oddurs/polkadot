function gclone --description 'Clone into ~/Code and cd there'
    set -l url $argv[1]
    set -l name (string replace -r '\.git$' '' -- (basename $url))
    git clone $url $HOME/Code/$name; and cd $HOME/Code/$name
end
