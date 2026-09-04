function fkill --description 'Pick a process with fzf and kill it'
    set -l pid (ps -ef | sed 1d | fzf -m --header='select processes to kill' | awk '{print $2}')
    test -n "$pid"; and echo $pid | xargs kill -9
end
