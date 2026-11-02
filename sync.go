package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/oddurs/polkadot/internal/link"
	"github.com/oddurs/polkadot/internal/step"
	"github.com/oddurs/polkadot/internal/ui"
)

type managedRepo struct{ name, path string }

func managedRepos(root, home string) []managedRepo {
	return []managedRepo{
		{"polkadot", root},
		{"fresh-config", filepath.Join(home, ".config", "fresh")},
	}
}

// doSync fast-forwards clean default branches only. It is safe to schedule:
// it never merges, stashes, or changes an active branch.
func doSync(root string, l *link.Linker, t *tally) {
	ui.Section("sync")
	home, _ := os.UserHomeDir()
	changed := false

	for _, r := range managedRepos(root, home) {
		dir := r.path
		if out, err := step.RunDir(dir, "git", "rev-parse", "--is-inside-work-tree"); err != nil || out != "true" {
			ui.Result(r.name, "skipped", "not cloned")
			t.skip++
			continue
		}
		if out, err := step.RunDir(dir, "git", "status", "--porcelain"); err == nil && out != "" {
			ui.Result(r.name, "skipped", "uncommitted changes")
			t.skip++
			continue
		}
		branch, branchErr := step.RunDir(dir, "git", "branch", "--show-current")
		defaultRef, defaultErr := step.RunDir(dir, "git", "symbolic-ref", "--short", "refs/remotes/origin/HEAD")
		if branchErr != nil || defaultErr != nil || branch == "" {
			ui.Result(r.name, "skipped", "remote default branch unknown")
			t.skip++
			continue
		}
		if branch != strings.TrimPrefix(defaultRef, "origin/") {
			ui.Result(r.name, "skipped", "not on default branch")
			t.skip++
			continue
		}
		if _, err := step.RunDir(dir, "git", "rev-parse", "--abbrev-ref", "--symbolic-full-name", "@{u}"); err != nil {
			ui.Result(r.name, "skipped", "branch has no upstream")
			t.skip++
			continue
		}
		before, _ := step.RunDir(dir, "git", "rev-parse", "HEAD")
		if *dryRun {
			ui.Result(r.name, "would", "git pull --ff-only")
			t.done++
			continue
		}
		if _, err := step.RunDir(dir, "git", "pull", "--ff-only", "--quiet"); err != nil {
			ui.Result(r.name, "failed", "cannot fast-forward - diverged, or no network")
			t.fail++
			continue
		}
		after, _ := step.RunDir(dir, "git", "rev-parse", "HEAD")
		if before == after {
			ui.Result(r.name, "already", "current")
			t.skip++
			continue
		}
		short, _ := step.RunDir(dir, "git", "log", "--oneline", before+".."+after, "-3")
		ui.Result(r.name, "updated", strings.TrimSpace(strings.SplitN(short, "\n", 2)[0]))
		t.done++
		changed = true
	}
	ui.Blank()

	if changed {
		doLink(l, t)
	}
}

// Weekly is enough for configuration repositories and keeps the log quiet.
const syncPlist = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>Label</key><string>dev.polkadot.sync</string>
	<key>ProgramArguments</key>
	<array>
		<string>%s</string>
		<string>run</string>
		<string>.</string>
		<string>sync</string>
	</array>
	<key>WorkingDirectory</key><string>%s</string>
	<key>StartCalendarInterval</key>
	<dict>
		<key>Weekday</key><integer>1</integer>
		<key>Hour</key><integer>9</integer>
		<key>Minute</key><integer>0</integer>
	</dict>
	<key>RunAtLoad</key><false/>
	<key>StandardOutPath</key><string>%s/Library/Logs/polkadot-sync.log</string>
	<key>StandardErrorPath</key><string>%s/Library/Logs/polkadot-sync.log</string>
</dict>
</plist>
`

func buildSyncPlist(root, goBin, home string) string {
	return fmt.Sprintf(syncPlist,
		xmlEscape(goBin),
		xmlEscape(root),
		xmlEscape(home),
		xmlEscape(home),
	)
}

func doTimer(t *tally) {
	ui.Section("timer")
	home, _ := os.UserHomeDir()
	root, _ := repoRoot()
	dst := filepath.Join(home, "Library", "LaunchAgents", "dev.polkadot.sync.plist")

	goBin, err := step.Run("/bin/sh", "-c", "command -v go")
	if err != nil {
		ui.Result("go", "failed", "not on PATH; the timer needs it to run `go run .`")
		t.fail++
		ui.Blank()
		return
	}
	body := buildSyncPlist(root, strings.TrimSpace(goBin), home)

	if *dryRun {
		ui.Result("dev.polkadot.sync", "would", "install, Mondays at 09:00")
		t.done++
		ui.Blank()
		return
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		ui.Result("LaunchAgents", "failed", err.Error())
		t.fail++
		ui.Blank()
		return
	}
	if err := os.WriteFile(dst, []byte(body), 0o644); err != nil {
		ui.Result(dst, "failed", err.Error())
		t.fail++
		ui.Blank()
		return
	}
	uid := os.Getuid()
	_, _ = step.Sh(fmt.Sprintf("launchctl bootout gui/%d/dev.polkadot.sync 2>/dev/null", uid))
	if _, err := step.Sh(fmt.Sprintf("launchctl bootstrap gui/%d %q", uid, dst)); err != nil {
		ui.Result("dev.polkadot.sync", "failed", "written, but launchctl refused it")
		t.fail++
	} else {
		ui.Result("dev.polkadot.sync", "loaded", "Mondays at 09:00")
		t.done++
	}
	ui.Blank()
}
