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
