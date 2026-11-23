#!/bin/sh
# Bootstrap a clean machine, then hand off to the Go installer. Keep this
# shell script deliberately small: it exists only to provide Go and git.
set -eu

repo_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)

bootstrap_linux() {
    if ! command -v apt-get >/dev/null 2>&1; then
        echo "polkadot supports Debian and Ubuntu on Linux. apt-get was not found." >&2
        exit 1
    fi
    missing=""
    command -v go >/dev/null 2>&1 || missing="$missing golang-go"
    command -v git >/dev/null 2>&1 || missing="$missing git"
    command -v curl >/dev/null 2>&1 || missing="$missing curl ca-certificates"
    [ -z "$missing" ] && return 0

    # A container runs as root and has no sudo.
    sudo=""
    [ "$(id -u)" -eq 0 ] || sudo="sudo"
    $sudo apt-get update -qq
    # shellcheck disable=SC2086 # the package list is split into words on purpose
    $sudo env DEBIAN_FRONTEND=noninteractive apt-get install -y $missing
}

bootstrap_darwin() {
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
}

case "$(uname -s)" in
    Darwin) bootstrap_darwin ;;
    Linux) bootstrap_linux ;;
    *) echo "polkadot supports macOS and Linux." >&2; exit 1 ;;
esac

cd "$repo_dir"
exec go run . "$@"
