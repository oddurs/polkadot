# ─── transient prompt ─────────────────────────────────────────────────────
# Once a command has run, its prompt collapses to just the character and the
# time. Scrollback becomes a clean column of `❯ command` with a timestamp on
# the right, instead of the same path and git status repeated forty times.
#
# The live prompt is unaffected — it still carries everything.

if status is-interactive; and type -q starship
    function starship_transient_prompt_func
        starship module character
    end

    function starship_transient_rprompt_func
        # Keep the clock in history: scrollback doubles as a log of when you
        # ran what.
        set_color 245361
        date "+%H:%M"
        set_color normal
    end

    enable_transience
end
