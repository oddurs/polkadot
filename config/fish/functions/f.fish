function f --description 'Browse with findr and cd to where you quit'
    set -l tmp (mktemp)
    findr --cwd-file $tmp $argv
    set -l dir (command cat $tmp)
    rm -f $tmp
    test -n "$dir"; and test "$dir" != "$PWD"; and cd $dir
end
