function openrouter-key --description 'Check an OpenRouter API key and save it for this machine'
    # OpenCode reads OPENROUTER_API_KEY from the environment. The key belongs
    # to this machine, so it goes in conf.d/90-local.fish, which git ignores.
    read --silent --prompt-str 'OpenRouter key (nothing will show): ' key
    set key (string trim -- $key)
    if not string match -q 'sk-or-*' -- $key
        echo 'That does not look like an OpenRouter key: it should start with sk-or-. Nothing was saved.' >&2
        return 1
    end

    set -l code (curl -s -o /dev/null -w '%{http_code}' --max-time 15 -H "Authorization: Bearer $key" https://openrouter.ai/api/v1/key)
    if test "$code" != 200
        echo "OpenRouter did not accept the key (HTTP $code). Nothing was saved." >&2
        return 1
    end

    set -l file $__fish_config_dir/conf.d/90-local.fish
    touch $file
    # Replace an older key; keep every other line.
    string match -v -r '^set -gx OPENROUTER_API_KEY ' <$file >$file.new
    echo "set -gx OPENROUTER_API_KEY $key" >>$file.new
    command mv -f $file.new $file
    chmod 600 $file
    set -gx OPENROUTER_API_KEY $key
    echo 'OpenRouter accepted the key. Saved for this machine; new terminals have it.'
end
