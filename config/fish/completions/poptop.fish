# poptop — a system monitor you can rewind
#
# Subcommand-style flags come first and are mutually exclusive with each other;
# the rest are settings that also exist as `key = value` in
# ~/.config/poptop/poptop.conf.

function __poptop_no_mode
    # True until one of the exclusive modes has been typed.
    not __fish_seen_argument -l once -l bench -l days -l schema -l report -l read -l export -l check-theme
end

function __poptop_days
    command poptop --days 2>/dev/null | string match -r '^\d{4}-\d{2}-\d{2}' 
end

complete -c poptop -f

# ── modes ────────────────────────────────────────────────────────────────────
complete -c poptop -n __poptop_no_mode -l once   -d 'Print one plain-text sample and exit'
complete -c poptop -n __poptop_no_mode -l report -d 'Summarise a recorded day (default: today)'
complete -c poptop -n __poptop_no_mode -l read   -d 'Open a recorded day and scrub it'
complete -c poptop -n __poptop_no_mode -l export -d 'Machine-readable output: json or line'
complete -c poptop -n __poptop_no_mode -l schema -d 'Every record, field, type and unit'
complete -c poptop -n __poptop_no_mode -l days   -d 'List the recorded days and their sizes'
complete -c poptop -n __poptop_no_mode -l bench  -d 'Time 20 collection passes (development)'
complete -c poptop -n __poptop_no_mode -l check-theme -r -d 'Measure a theme for legibility'

# A date argument, offered from what has actually been recorded.
complete -c poptop -n '__fish_seen_argument -l read -l report' -f -a '(__poptop_days)' -d 'recorded day'
complete -c poptop -n '__fish_seen_argument -l export' -f -a 'json line' -d 'format'

# ── settings ─────────────────────────────────────────────────────────────────
complete -c poptop -l glyphs   -x -a 'braille block ascii' -d 'Timeline drawing'
complete -c poptop -l color    -x -a 'auto mono 16 256 true' -d 'Colour tier (honours NO_COLOR)'
complete -c poptop -l theme    -x -a 'safe classic auto' -d 'Built-in name or a file in ~/.config/poptop/themes'
complete -c poptop -l interval -x -a '500ms 1s 2s 5s 10m' -d 'Time between samples'
complete -c poptop -l window   -x -a '5m 10m 30m 1h' -d 'History retained, as time'
complete -c poptop -l warn     -x -d "Where 'getting busy' begins (default 50)"
complete -c poptop -l critical -x -d "Where 'in trouble' begins (default 80)"

complete -c poptop -l store    -x -a 'on off' -d 'Keep history across restarts (default off)'
complete -c poptop -l signals  -x -a 'on off' -d 'Allow x and X to signal a process (default off)'
complete -c poptop -l log      -x -a 'on off' -d 'Write a daily log that outlives the process (default off)'
complete -c poptop -l log-interval -x -a '1m 10m 1h' -d 'How often a sample reaches the log'
complete -c poptop -l log-days     -x -d 'Days of log kept (default 7)'
complete -c poptop -l log-bytes    -x -a '128M 512M 2G' -d 'Bytes of log kept (default 512M)'

complete -c poptop -s h -l help    -d 'Show help'
complete -c poptop -s V -l version -d 'Show version'
