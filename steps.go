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
type linkMap struct{ from, to string }

var links = []linkMap{
	{"config/fish", ".config/fish"},
	{"config/starship.toml", ".config/starship.toml"},
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
	// opencode reads these from ~/.config/opencode. node_modules, bun.lock and
	// the skills are machine state and stay out of the repo.
	{"config/opencode/opencode.jsonc", ".config/opencode/opencode.jsonc"},
	{"config/opencode/AGENTS.md", ".config/opencode/AGENTS.md"},
	{"config/opencode/agents", ".config/opencode/agents"},
	{"config/opencode/commands", ".config/opencode/commands"},
	{"config/opencode/templates", ".config/opencode/templates"},
	{"config/claude/settings.json", ".claude/settings.json"},
	// Instructions, not state: CLAUDE.md loads every session, rules/ load when
	// a matching file is touched, skills/ load on demand. Everything else
	// under ~/.claude — history, projects, plugins, auto memory — is state
	// this machine writes, and stays out of the repo.
	{"config/claude/CLAUDE.md", ".claude/CLAUDE.md"},
	{"config/claude/rules", ".claude/rules"},
	{"config/claude/skills", ".claude/skills"},
	{"config/ripgrep", ".config/ripgrep"},
	{"config/tmux/tmux.conf", ".tmux.conf"},
	{"config/vscode/settings.json", "Library/Application Support/Code/User/settings.json"},
	{"home/zshrc", ".zshrc"},
	{"home/gitconfig", ".gitconfig"},
	{"home/gitconfig.local", ".gitconfig.local"},
	{"home/gitignore_global", ".gitignore_global"},
	{"home/git_commit_template", ".git_commit_template"},
}

// doUnlink removes only links that point back into this checkout. It leaves
// ordinary files and somebody else's links alone, so it is safe for retiring
// a machine or moving the repository.
func doUnlink(l *link.Linker, t *tally) {
	ui.Section("dotfiles")
	home, _ := os.UserHomeDir()
	for _, m := range append(links, themeLinks...) {
		src := filepath.Join(l.Root, m.from)
		dst := filepath.Join(home, m.to)
		managed, err := l.Unlink(src, dst)
		if err != nil {
			ui.Result("~/"+m.to, "failed", err.Error())
			t.fail++
			continue
		}
		if !managed {
			ui.Result("~/"+m.to, "skipped", "not managed")
			t.skip++
			continue
		}
		if *dryRun {
			ui.Result("~/"+m.to, "would", "remove link")
		} else {
			ui.Result("~/"+m.to, "removed", "link only")
		}
		t.done++
	}
	ui.Blank()
}

func doLink(l *link.Linker, t *tally) {
	ui.Section("dotfiles")
	linkMappings(l, t, append(links, themeLinks...))
	ui.Blank()
}

func linkMappings(l *link.Linker, t *tally, mappings []linkMap) {
	home, _ := os.UserHomeDir()
	for _, m := range mappings {
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
		installer := `NONINTERACTIVE=1 /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"`
		if !confirmInstaller("Homebrew", installer) {
			ui.Result("homebrew", "skipped", "installer not approved")
			t.skip++
			ui.Blank()
			return
		}
		ui.Note("installing homebrew…")
		if err := step.Stream("/bin/bash", "-c",
			installer); err != nil {
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
		if !confirmInstaller(a.bin, a.how) {
			ui.Result(a.bin, "skipped", "installer not approved")
			t.skip++
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

func doDoctor(home, root string) {
	ui.Section("links")
	for _, m := range append(links, themeLinks...) {
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
		switch {
		case step.Runs(b):
			ui.Result(b, "already", "")
		case step.Has(b):
			// On PATH but cannot start: a half-installed package, or a shim
			// whose real binary went missing. Worth saying loudly, because
			// everything else about it looks fine.
			ui.Result(b, "failed", "on PATH but will not run — reinstall it")
		default:
			ui.Result(b, "skipped", "absent")
		}
	}
	ui.Blank()
	ui.Section("shell")
	ui.Result("login shell", "already", os.Getenv("SHELL"))
	ui.Blank()
}
