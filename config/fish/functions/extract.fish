# One command for every archive format, because nobody remembers tar's flags.
function extract --description 'Extract any archive'
    for f in $argv
        if not test -f $f
            echo "extract: $f is not a file" >&2
            continue
        end
        switch $f
            case '*.tar.bz2' '*.tbz2'; tar xjf $f
            case '*.tar.gz' '*.tgz';   tar xzf $f
            case '*.tar.xz' '*.txz';   tar xJf $f
            case '*.tar.zst';          tar --zstd -xf $f
            case '*.tar';              tar xf $f
            case '*.bz2';              bunzip2 $f
            case '*.gz';               gunzip $f
            case '*.zip';              unzip -q $f
            case '*.7z';               7z x $f
            case '*.rar';              unar $f
            case '*';                  echo "extract: don't know how to open $f" >&2
        end
    end
end
