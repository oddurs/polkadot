// polkadot bootstraps a Mac the way I like it.
//
//	polkadot install     everything: brew, packages, links, shell, tools
//	polkadot link        symlinks only
//	polkadot brew        Brewfile only
//	polkadot doctor      report what is and isn't in place
//
// Flags: --dry-run to see the plan, --only=<step> to run one part.
//
// Deliberately stdlib-only: on a fresh machine the sole prerequisite is the Go
// toolchain, so this can't be blocked by a package manager it is meant to
// install. It never writes macOS `defaults` — system preferences stay yours.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/oddurs/polkadot/internal/link"
	"github.com/oddurs/polkadot/internal/ui"
)

var (
	dryRun = flag.Bool("dry-run", false, "show what would change, touch nothing")
	only   = flag.String("only", "", "run one step: brew, link, shell, tools, agents")
)

type tally struct{ done, skip, fail int }

func (t *tally) count(s link.Status) {
	switch s {
	case link.Created, link.Backed:
		t.done++
	case link.Already:
		t.skip++
	case link.Would:
		t.done++
	}
}

func main() {
	flag.Parse()
	cmd := "install"
	if flag.NArg() > 0 {
		cmd = flag.Arg(0)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, "cannot resolve home:", err)
		os.Exit(1)
	}
	root, err := repoRoot()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	mode := ""
	if *dryRun {
		mode = "dry run — nothing will be written"
	}
	ui.Banner(mode)

	l := link.New(root, home, *dryRun)
	t := &tally{}

	switch cmd {
	case "install":
		runIf("brew", func() { doBrew(t) })
		runIf("link", func() { doLink(l, t) })
		runIf("shell", func() { doShell(home, t) })
		runIf("tools", func() { doTools(t) })
		runIf("agents", func() { doAgents(t) })
	case "link":
		doLink(l, t)
	case "brew":
		doBrew(t)
	case "doctor":
		doDoctor(home, root)
		return
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n\n", cmd)
		flag.Usage()
		os.Exit(2)
	}

	if l.BackupUsed() {
		ui.Blank()
		ui.Note("replaced files were moved to " + strings.Replace(l.BackupDir, home, "~", 1))
	}
	ui.Summary(t.done, t.skip, t.fail)
}

func runIf(name string, fn func()) {
	if *only == "" || *only == name {
		fn()
	}
}

// repoRoot finds the checkout containing this binary's source, so the tool works
// from any working directory.
func repoRoot() (string, error) {
	if env := os.Getenv("POLKADOT_ROOT"); env != "" {
		return env, nil
	}
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for dir := wd; dir != "/"; dir = filepath.Dir(dir) {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			if _, err := os.Stat(filepath.Join(dir, "config")); err == nil {
				return dir, nil
			}
		}
	}
	home, _ := os.UserHomeDir()
	guess := filepath.Join(home, "Code", "polkadot")
	if _, err := os.Stat(filepath.Join(guess, "go.mod")); err == nil {
		return guess, nil
	}
	return "", fmt.Errorf("cannot find the polkadot checkout; run from inside it or set POLKADOT_ROOT")
}
