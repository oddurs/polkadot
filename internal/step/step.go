// Package step runs shell work and reports it uniformly.
package step

import (
	"os"
	"os/exec"
	"strings"
)

// Has reports whether a binary is on PATH.
func Has(bin string) bool {
	_, err := exec.LookPath(bin)
	return err == nil
}

// Runs reports whether a binary is on PATH *and* can actually execute.
//
// The two are not the same, and the gap is not theoretical: a partially
// installed npm package leaves its launcher shim on PATH while the platform
// binary it execs is missing, so `codex` resolved, `Has` said yes, and every
// invocation died with ENOENT. A doctor that reports such a tool as present is
// worse than one that does not check at all, because it answers the question
// you were asking with the wrong fact.
//
// `--version` is the cheapest thing that proves a process started. A tool that
// does not support it is reported on PATH rather than guessed at.
func Runs(bin string) bool {
	if !Has(bin) {
		return false
	}
	cmd := exec.Command(bin, "--version")
	cmd.Env = os.Environ()
	cmd.Stdin = nil
	return cmd.Run() == nil
}

// Run executes a command, returning combined output.
func Run(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Env = os.Environ()
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

// Sh runs a line through the shell, for pipelines and substitutions.
func Sh(line string) (string, error) {
	return Run("/bin/sh", "-c", line)
}

// Stream runs a command with its output attached to the terminal, for things
// like brew bundle where the progress is the point.
func Stream(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Env = os.Environ()
	cmd.Stdout, cmd.Stderr, cmd.Stdin = os.Stdout, os.Stderr, os.Stdin
	return cmd.Run()
}
