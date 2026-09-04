# `up 3` is three directories up. `up` is one.
function up --description 'Go up N directories'
    set -l n (test -n "$argv[1]"; and echo $argv[1]; or echo 1)
    set -l path ""
    for i in (seq $n); set path "../$path"; end
    cd $path
end
