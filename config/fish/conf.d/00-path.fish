# PATH, in priority order. fish_add_path is idempotent and prepends, so the
# last call wins — hence the reverse ordering of intent here.

# Homebrew must be set up before anything that lives inside it.
if test -x /opt/homebrew/bin/brew
    /opt/homebrew/bin/brew shellenv | source
end

fish_add_path -g $HOME/.cargo/bin
fish_add_path -g $HOME/.bun/bin
fish_add_path -g $HOME/go/bin
fish_add_path -g $HOME/.local/bin
