# When a command is missing, say what would install it rather than just
# failing. Homebrew's own handler shells out on every miss and is slow; this
# only searches when the name looks like a real package.
function fish_command_not_found
    set -l cmd $argv[1]
    echo (set_color red)"not found: "(set_color normal)$cmd

    if type -q brew
        set -l hit (brew formulae 2>/dev/null | string match -e $cmd | head -3)
        if test -n "$hit"
            echo (set_color 4e5165)"  brew install "(string join ' | ' $hit)(set_color normal)
            return 127
        end
    end

    # A near-miss on something already installed is more often the real cause.
    set -l near (string match -r ".*$cmd.*" -- (builtin complete -C "" 2>/dev/null | string split \t)[1] 2>/dev/null | head -3)
    test -n "$near"; and echo (set_color 4e5165)"  did you mean: "(string join ', ' $near)(set_color normal)
    return 127
end
