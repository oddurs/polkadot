# A static server on the current directory. Port defaults to 8000.
function serve --description 'Serve the current directory over HTTP'
    set -l port (test -n "$argv[1]"; and echo $argv[1]; or echo 8000)
    echo "→ http://localhost:$port  ($PWD)"
    python3 -m http.server $port
end
