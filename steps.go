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
	{"config/nvim", ".config/nvim"},
	{"config/lazygit", ".config/lazygit"},
	{"config/bat", ".config/bat"},
	{"config/atuin", ".config/atuin"},
	{"config/mise", ".config/mise"},
	{"config/btop", ".config/btop"},
	{"config/herdr", ".config/herdr"},
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

	// fish must be a known shell before chsh will accept it.
	shells, _ := os.ReadFile("/etc/shells")
	if !strings.Contains(string(shells), fishPath) {
		if *dryRun {
			ui.Result("/etc/shells", "would", "register fish")
			t.done++
		} else {
			ui.Result("/etc/shells", "skipped", "needs sudo: echo "+fishPath+" | sudo tee -a /etc/shells")
			t.skip++
		}
	} else {
		ui.Result("/etc/shells", "already", "fish registered")
		t.skip++
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
		ui.Result("login shell", "skipped", "run: chsh -s "+fishPath)
		t.skip++
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
		{"nvim", "editor"}, {"gh", "github"},
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
	for _, b := range []string{"brew", "fish", "starship", "mise", "zoxide", "atuin", "nvim", "lazygit", "lazydocker", "btop", "gh", "claude", "codex", "opencode", "herdr"} {
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
