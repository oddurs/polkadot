function mcdir --description 'Create a directory and enter it'
    if test (count $argv) -ne 1; or test -z "$argv[1]"
        printf '%s\n' 'usage: mcdir <directory>' >&2
        return 2
    end

    set -l destination "$argv[1]"
    if not string match -q '/*' -- "$destination"
        set destination "./$destination"
    end

    command mkdir -p -- "$destination"; or return
    cd -- "$destination"
end
