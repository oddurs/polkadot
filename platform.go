package main

import (
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/oddurs/polkadot/internal/step"
)

// goos is a variable so a test can check both targets on one machine.
var goos = runtime.GOOS

// darwinOnly names the link sources that only a Mac uses: the desktop
// terminal and the desktop editor. A Linux machine is a server.
var darwinOnly = map[string]bool{
	"config/ghostty/config":         true,
	"config/vscode/settings.json":   true,
	"theme/london-moquette/ghostty": true,
}

func applies(m linkMap) bool { return goos == "darwin" || !darwinOnly[m.from] }

func active(mappings []linkMap) []linkMap {
	var out []linkMap
	for _, m := range mappings {
		if applies(m) {
			out = append(out, m)
		}
	}
	return out
}

// sudo is the prefix for a command that needs root. A container runs as root
// and has no sudo, so the prefix is empty there.
func sudo() string {
	if os.Getuid() == 0 {
		return ""
	}
	return "sudo "
}

// loginShell reads the account's shell from the system. $SHELL is the wrong
// source: a process started from another shell inherits that shell's value.
func loginShell() string {
	u, err := user.Current()
	if err != nil {
		return os.Getenv("SHELL")
	}
	shell := ""
	if goos == "darwin" {
		out, _ := step.Run("dscl", ".", "-read", "/Users/"+u.Username, "UserShell")
		shell = shellFromDscl(out)
	} else {
		out, _ := step.Run("getent", "passwd", u.Username)
		shell = shellFromPasswd(out)
	}
	if shell == "" {
		return os.Getenv("SHELL")
	}
	return shell
}

func shellFromDscl(out string) string {
	return strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(out), "UserShell:"))
}

func shellFromPasswd(line string) string {
	fields := strings.Split(strings.TrimSpace(line), ":")
	if len(fields) < 7 {
		return ""
	}
	return fields[6]
}

// addLocalPaths puts the per-user install directories on PATH for this
// process. An installer that this run starts puts its binary there, and a
// later step in the same run must find it.
func addLocalPaths(home string) {
	extra := []string{
		filepath.Join(home, ".local", "share", "mise", "shims"),
		filepath.Join(home, ".local", "bin"),
		// The OpenCode installer uses its own directory.
		filepath.Join(home, ".opencode", "bin"),
	}
	os.Setenv("PATH", strings.Join(extra, string(os.PathListSeparator))+string(os.PathListSeparator)+os.Getenv("PATH"))
}

// readList reads one entry per line and skips blank lines and # comments.
func readList(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var out []string
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		out = append(out, line)
	}
	return out, nil
}
