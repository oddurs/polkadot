#!/bin/sh
# Bootstrap a clean Mac, then hand off to the Go installer. Keep this shell
# script deliberately small: it exists only to provide Go on a new machine.
set -eu

if [ "$(uname -s)" != "Darwin" ]; then
    echo "polkadot currently supports macOS only." >&2
    exit 1
fi

if ! xcode-select -p >/dev/null 2>&1; then
    xcode-select --install
    echo "The Xcode Command Line Tools installer has been opened." >&2
    echo "Finish it, then run ./bootstrap.sh again." >&2
    exit 1
fi

if ! command -v brew >/dev/null 2>&1; then
    if [ "${POLKADOT_BOOTSTRAP_YES:-}" != "1" ]; then
        printf '%s' 'Homebrew is missing. Run its official installer? [y/N] '
        read -r answer
        case "$answer" in
            y|Y|yes|YES) ;;
            *) echo "Homebrew was not installed. Run bootstrap again when ready." >&2; exit 1 ;;
        esac
    fi
    /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
fi

if [ -x /opt/homebrew/bin/brew ]; then
    eval "$(/opt/homebrew/bin/brew shellenv)"
elif [ -x /usr/local/bin/brew ]; then
    eval "$(/usr/local/bin/brew shellenv)"
fi

if ! command -v go >/dev/null 2>&1; then
    brew install go
fi

repo_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
cd "$repo_dir"
exec go run . "$@"
