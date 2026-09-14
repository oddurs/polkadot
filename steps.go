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

// links maps repo paths to where they belong. config/* goes to ~/.config/*,
// home/* goes to ~/*.
var links = []struct{ from, to string }{
	{"config/fish", ".config/fish"},
	{"config/starship.toml", ".config/starship.toml"},
	{"config/ghostty/themes", ".config/ghostty/themes"},
	{"config/ghostty/config", ".config/ghostty/config"},
	{"config/lazygit", ".config/lazygit"},
	{"config/bat", ".config/bat"},
	{"config/atuin", ".config/atuin"},
	{"config/mise", ".config/mise"},
	{"config/btop", ".config/btop"},
	{"config/herdr", ".config/herdr"},
	// Only config.yml: hosts.yml beside it holds the OAuth token.
	{"config/gh/config.yml", ".config/gh/config.yml"},
	{"config/codex/config.toml", ".codex/config.toml"},
	{"config/claude/settings.json", ".claude/settings.json"},
	// Instructions, not state: CLAUDE.md loads every session, rules/ load when
	// a matching file is touched, skills/ load on demand. Everything else
	// under ~/.claude — history, projects, plugins, auto memory — is state
	// this machine writes, and stays out of the repo.
	{"config/claude/CLAUDE.md", ".claude/CLAUDE.md"},
	{"config/claude/rules", ".claude/rules"},
	{"config/claude/skills", ".claude/skills"},
	{"config/ripgrep", ".config/ripgrep"},
	{"config/vscode/settings.json", "Library/Application Support/Code/User/settings.json"},
	{"home/zshrc", ".zshrc"},
	{"home/gitconfig", ".gitconfig"},
	{"home/gitignore_global", ".gitignore_global"},
	{"home/git_commit_template", ".git_commit_template"},
}

func doLink(l *link.Linker, t *tally) {
	ui.Section("dotfiles")
	home, _ := os.UserHomeDir()
	for _, m := range links {
		src := filepath.Join(l.Root, m.from)
		dst := filepath.Join(home, m.to)
		if _, err := os.Stat(src); err != nil {
			ui.Result("~/"+m.to, "skipped", "not in repo")
			t.skip++
			continue
		}
		st, err := l.Link(src, dst)
		if err != nil {
			ui.Result("~/"+m.to, "failed", err.Error())
			t.fail++
			continue
		}
		switch st {
		case link.Already:
			ui.Result("~/"+m.to, "already", "")
		case link.Backed:
			ui.Result("~/"+m.to, "linked", "previous moved to backup")
		case link.Would:
			ui.Result("~/"+m.to, "would", "link")
		default:
			ui.Result("~/"+m.to, "linked", "")
		}
		t.count(st)
	}
	ui.Blank()
}

func doBrew(t *tally) {
	ui.Section("homebrew")
	if !step.Has("brew") {
		if *dryRun {
			ui.Result("homebrew", "would", "install")
			t.done++
			ui.Blank()
			return
		}
		ui.Note("installing homebrew…")
		if err := step.Stream("/bin/bash", "-c",
			`NONINTERACTIVE=1 /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"`); err != nil {
			ui.Result("homebrew", "failed", err.Error())
			t.fail++
			ui.Blank()
			return
		}
		ui.Result("homebrew", "installed", "")
		t.done++
	} else {
		ui.Result("homebrew", "already", "")
		t.skip++
	}

	root, _ := repoRoot()
	bf := filepath.Join(root, "Brewfile")
	if _, err := os.Stat(bf); err != nil {
		ui.Blank()
		return
	}
	if *dryRun {
		out, _ := step.Sh(fmt.Sprintf("brew bundle check --file=%q 2>&1 | head -20", bf))
		ui.Result("Brewfile", "would", "install missing")
		for _, line := range strings.Split(out, "\n") {
			if line != "" {
				ui.Note(line)
			}
		}
		t.done++
		ui.Blank()
		return
	}
	if _, err := step.Sh(fmt.Sprintf("brew bundle check --file=%q", bf)); err == nil {
		ui.Result("Brewfile", "already", "all present")
		t.skip++
	} else {
		ui.Note("installing from Brewfile…")
		if err := step.Stream("brew", "bundle", "--file="+bf); err != nil {
			ui.Result("Brewfile", "failed", err.Error())
			t.fail++
		} else {
			ui.Result("Brewfile", "installed", "")
			t.done++
		}
	}
	ui.Blank()
}

func doShell(home string, t *tally) {
	ui.Section("shell")

	fishPath, err := step.Run("/bin/sh", "-c", "command -v fish")
	if err != nil || fishPath == "" {
		ui.Result("fish", "failed", "not installed — run the brew step first")
		t.fail++
		ui.Blank()
		return
	}
	ui.Result("fish", "already", fishPath)
	t.skip++

	// fish must be a known shell before chsh will accept it. Both of these
	// need a password, so they run attached to the terminal — on a fresh
	// machine you want to type it once, not be handed homework.
	shells, _ := os.ReadFile("/etc/shells")
	switch {
	case strings.Contains(string(shells), fishPath):
		ui.Result("/etc/shells", "already", "fish registered")
		t.skip++
	case *dryRun:
		ui.Result("/etc/shells", "would", "register fish (sudo)")
		t.done++
	default:
		ui.Note("registering fish in /etc/shells — sudo will ask for your password")
		if err := step.Stream("/bin/sh", "-c",
			"echo "+fishPath+" | sudo tee -a /etc/shells >/dev/null"); err != nil {
			ui.Result("/etc/shells", "failed", "run by hand: echo "+fishPath+" | sudo tee -a /etc/shells")
			t.fail++
		} else {
			ui.Result("/etc/shells", "registered", "")
			t.done++
		}
	}

	cur := os.Getenv("SHELL")
	switch {
	case cur == fishPath:
		ui.Result("login shell", "already", "fish")
		t.skip++
	case *dryRun:
		ui.Result("login shell", "would", "chsh -s "+fishPath)
		t.done++
	default:
		ui.Note("changing the login shell — chsh will ask for your password")
		if err := step.Stream("chsh", "-s", fishPath); err != nil {
			ui.Result("login shell", "failed", "run by hand: chsh -s "+fishPath)
			t.fail++
		} else {
			ui.Result("login shell", "changed", "fish — open a new terminal")
			t.done++
		}
	}
	ui.Blank()
}

func doTools(t *tally) {
	ui.Section("tools")
	for _, tool := range []struct{ bin, note string }{
		{"mise", "runtimes"}, {"starship", "prompt"}, {"zoxide", "jump"},
		{"atuin", "history"}, {"fzf", "fuzzy"}, {"rg", "search"},
		{"fd", "find"}, {"bat", "cat"}, {"eza", "ls"}, {"delta", "diff"},
		{"lazygit", "git"}, {"lazydocker", "docker"}, {"btop", "monitor"},
		{"fresh", "editor"}, {"gh", "github"},
	} {
		if step.Has(tool.bin) {
			ui.Result(tool.bin, "already", tool.note)
			t.skip++
		} else {
			ui.Result(tool.bin, "failed", "missing — brew step should have installed it")
			t.fail++
		}
	}
	ui.Blank()
}

func doAgents(t *tally) {
	ui.Section("coding agents")
	for _, a := range []struct{ bin, how string }{
		{"claude", "curl -fsSL https://claude.ai/install.sh | bash"},
		{"codex", "npm i -g @openai/codex"},
		{"opencode", "brew install sst/tap/opencode"},
		{"herdr", "curl -fsSL https://herdr.dev/install.sh | sh"},
	} {
		if step.Has(a.bin) {
			ui.Result(a.bin, "already", "")
			t.skip++
			continue
		}
		if *dryRun {
			ui.Result(a.bin, "would", a.how)
			t.done++
			continue
		}
		ui.Note("installing " + a.bin + "…")
		if err := step.Stream("/bin/sh", "-c", a.how); err != nil {
			ui.Result(a.bin, "failed", err.Error())
			t.fail++
		} else {
			ui.Result(a.bin, "installed", "")
			t.done++
		}
	}
	ui.Blank()
}

// doTheme keeps the colour scheme out of this repo. Subway Seat generates a
// theme file per app per flavor; vendoring 117 of those here would mean
// regenerating them whenever the palette moves. Instead its own installer
// places absolute symlinks into ~/Code/subway-seat, which is why .gitignore
// excludes every directory it writes into.
func doTheme(t *tally) {
	ui.Section("theme")

	home, _ := os.UserHomeDir()
	repo := filepath.Join(home, "Code", "subway-seat")

	if _, err := os.Stat(repo); err != nil {
		if *dryRun {
			ui.Result("subway-seat", "would", "clone into ~/Code/subway-seat")
			t.done++
			ui.Blank()
			return
		}
		ui.Note("cloning subway-seat…")
		if err := step.Stream("git", "clone", "--quiet",
			"git@github.com:oddurs/subway-seat.git", repo); err != nil {
			ui.Result("subway-seat", "failed", err.Error())
			t.fail++
			ui.Blank()
			return
		}
		ui.Result("subway-seat", "cloned", "~/Code/subway-seat")
		t.done++
	} else {
		ui.Result("subway-seat", "already", "~/Code/subway-seat")
		t.skip++
	}

	args := "install --flavor " + themeFlavor + " --yes"
	if *dryRun {
		args += " --dry-run"
	}
	if err := step.Stream("/bin/sh", "-c",
		"cd "+repo+" && sh install.sh "+args+" >/dev/null"); err != nil {
		ui.Result(themeFlavor, "failed", err.Error())
		t.fail++
	} else {
		ui.Result(themeFlavor, "linked", "themes placed into the apps present")
		t.done++
	}
	ui.Blank()
}

func doDoctor(home, root string) {
	ui.Section("links")
	for _, m := range links {
		dst := filepath.Join(home, m.to)
		want := filepath.Join(root, m.from)
		got, err := os.Readlink(dst)
		switch {
		case err == nil && got == want:
			ui.Result("~/"+m.to, "already", "linked")
		case err == nil:
			ui.Result("~/"+m.to, "failed", "points elsewhere: "+got)
		default:
			if _, err := os.Lstat(dst); err == nil {
				ui.Result("~/"+m.to, "failed", "exists but is not a link")
			} else {
				ui.Result("~/"+m.to, "skipped", "absent")
			}
		}
	}
	ui.Blank()
	ui.Section("binaries")
	for _, b := range []string{"brew", "fish", "starship", "mise", "zoxide", "atuin", "lazygit", "lazydocker", "btop", "gh", "fresh", "claude", "codex", "opencode", "herdr"} {
		if step.Has(b) {
			ui.Result(b, "already", "")
		} else {
			ui.Result(b, "skipped", "absent")
		}
	}
	ui.Blank()
	ui.Section("shell")
	ui.Result("login shell", "already", os.Getenv("SHELL"))
	ui.Blank()
}

// The flavor everything is themed with. `sh ~/Code/subway-seat/install.sh
// switch <flavor>` changes it for the apps already installed; changing it here
// is what a fresh machine gets.
const themeFlavor = "moquette"

// The repositories this machine's configuration is actually made of. polkadot
// holds the dotfiles, subway-seat generates the theme, fresh-config is the
// editor. Keeping them in one list is what lets `sync` be one command rather
// than three things you have to remember.
var repos = []struct{ name, path, url string }{
	{"polkadot", "Code/polkadot", "git@github.com:oddurs/polkadot.git"},
	{"subway-seat", "Code/subway-seat", "git@github.com:oddurs/subway-seat.git"},
	{"fresh-config", ".config/fresh", "git@github.com:oddurs/fresh-config.git"},
}

// doSync fast-forwards each of them and puts back anything new.
//
// Fast-forward only, and never over a dirty tree: this runs unattended from a
// launchd timer, and the one thing it must never do is touch work in progress
// or leave a merge conflict in a config file that a shell is about to read.
// A repo it cannot advance is reported and skipped, not forced.
func doSync(l *link.Linker, t *tally) {
	ui.Section("sync")
	home, _ := os.UserHomeDir()
	changed := false

	for _, r := range repos {
		dir := filepath.Join(home, r.path)
		if _, err := os.Stat(filepath.Join(dir, ".git")); err != nil {
			ui.Result(r.name, "skipped", "not cloned")
			t.skip++
			continue
		}
		if out, err := step.Sh(fmt.Sprintf("git -C %q status --porcelain", dir)); err == nil && strings.TrimSpace(out) != "" {
			ui.Result(r.name, "skipped", "uncommitted changes")
			t.skip++
			continue
		}
		before, _ := step.Sh(fmt.Sprintf("git -C %q rev-parse HEAD", dir))
		if *dryRun {
			ui.Result(r.name, "would", "git pull --ff-only")
			t.done++
			continue
		}
		if _, err := step.Sh(fmt.Sprintf("git -C %q pull --ff-only --quiet 2>&1", dir)); err != nil {
			ui.Result(r.name, "failed", "cannot fast-forward — diverged, or no network")
			t.fail++
			continue
		}
		after, _ := step.Sh(fmt.Sprintf("git -C %q rev-parse HEAD", dir))
		if before == after {
			ui.Result(r.name, "already", "current")
			t.skip++
			continue
		}
		short, _ := step.Sh(fmt.Sprintf("git -C %q log --oneline %s..%s | head -3", dir, strings.TrimSpace(before), strings.TrimSpace(after)))
		ui.Result(r.name, "updated", strings.TrimSpace(strings.SplitN(short, "\n", 2)[0]))
		t.done++
		changed = true
	}
	ui.Blank()

	// Only relink when something moved. A sync that changed nothing should
	// print three lines and touch no files.
	if changed {
		doLink(l, t)
		doTheme(t)
	}
}

// The launchd agent that runs the above. Weekly rather than daily: these
// repositories change when I change them, and a timer that fires more often
// than the thing it watches is just noise in the log.
const syncPlist = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>Label</key><string>dev.polkadot.sync</string>
	<key>ProgramArguments</key>
	<array>
		<string>/bin/sh</string>
		<string>-lc</string>
		<string>cd %s && %s run . sync</string>
	</array>
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
	body := fmt.Sprintf(syncPlist, root, strings.TrimSpace(goBin), home, home)

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
	// bootout first so a changed plist is actually reloaded; it fails
	// harmlessly when nothing is loaded yet.
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
