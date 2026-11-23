package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/oddurs/polkadot/internal/link"
	"github.com/oddurs/polkadot/internal/step"
	"github.com/oddurs/polkadot/internal/ui"
)

// links maps repo paths to where they belong. config/* goes to ~/.config/*,
// home/* goes to ~/*.
type linkMap struct{ from, to string }

var links = []linkMap{
	// One instruction file for every agent. Each agent reads it from its own
	// location, so the working rules cannot drift between them.
	{"agents/AGENTS.md", ".claude/CLAUDE.md"},
	{"agents/AGENTS.md", ".codex/AGENTS.md"},
	{"agents/AGENTS.md", ".config/opencode/AGENTS.md"},
	{"agents/AGENTS.md", ".gemini/AGENTS.md"},
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
	{"config/opencode/STACK.md", ".config/opencode/STACK.md"},
	{"config/opencode/agents", ".config/opencode/agents"},
	{"config/opencode/commands", ".config/opencode/commands"},
	{"config/opencode/templates", ".config/opencode/templates"},
	{"config/claude/settings.json", ".claude/settings.json"},
	// Instructions, not state: CLAUDE.md loads every session, rules/ load when
	// a matching file is touched, skills/ load on demand. Everything else
	// under ~/.claude — history, projects, plugins, auto memory — is state
	// this machine writes, and stays out of the repo.
	{"config/claude/rules", ".claude/rules"},
	{"config/claude/skills", ".claude/skills"},
	{"config/ripgrep", ".config/ripgrep"},
	{"config/vscode/settings.json", "Library/Application Support/Code/User/settings.json"},
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
	for _, m := range active(append(links, themeLinks...)) {
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

// retired lists links that an earlier version created and this version no
// longer owns. Their sources are gone from the checkout, so without this a
// machine that updates keeps a dangling link at each target.
var retired = []linkMap{
	{"config/tmux/tmux.conf", ".tmux.conf"},
	{"home/zshrc", ".zshrc"},
}

func doLink(l *link.Linker, t *tally) {
	ui.Section("dotfiles")
	linkMappings(l, t, active(append(links, themeLinks...)))
	removeRetired(l, t)
	ui.Blank()
}

// removeRetired removes a retired link only when it still points into this
// checkout. A file the user put there instead is not ours to remove.
func removeRetired(l *link.Linker, t *tally) {
	for _, m := range retired {
		managed, err := l.Unlink(filepath.Join(l.Root, m.from), filepath.Join(l.Home, m.to))
		if err != nil {
			ui.Result("~/"+m.to, "failed", err.Error())
			t.fail++
			continue
		}
		if !managed {
			continue
		}
		if *dryRun {
			ui.Result("~/"+m.to, "would", "remove retired link")
		} else {
			ui.Result("~/"+m.to, "removed", "retired link")
		}
		t.done++
	}
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
	bf := filepath.Join(root, "packages", "darwin.Brewfile")
	if _, err := os.Stat(bf); err != nil {
		ui.Blank()
		return
	}
	if *dryRun {
		out, _ := step.Sh(fmt.Sprintf("brew bundle check --file=%q 2>&1 | head -20", bf))
		ui.Result("darwin.Brewfile", "would", "install missing")
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
		ui.Result("darwin.Brewfile", "already", "all present")
		t.skip++
	} else {
		ui.Note("installing from packages/darwin.Brewfile…")
		if err := step.Stream("brew", "bundle", "--file="+bf); err != nil {
			ui.Result("darwin.Brewfile", "failed", err.Error())
			t.fail++
		} else {
			ui.Result("darwin.Brewfile", "installed", "")
			t.done++
		}
	}
	ui.Blank()
}

func doShell(home string, t *tally) {
	ui.Section("shell")

	fishPath, err := step.Run("/bin/sh", "-c", "command -v fish")
	if err != nil || fishPath == "" {
		ui.Result("fish", "failed", "not installed — run the packages step first")
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
			"echo "+fishPath+" | "+sudo()+"tee -a /etc/shells >/dev/null"); err != nil {
			ui.Result("/etc/shells", "failed", "run by hand: echo "+fishPath+" | sudo tee -a /etc/shells")
			t.fail++
		} else {
			ui.Result("/etc/shells", "registered", "")
			t.done++
		}
	}

	cur := loginShell()
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

// baseBins are the programs that the system package manager installs. mise
// installs every other tool, and mise itself reports which are missing.
func baseBins() []string {
	bins := []string{"fish", "git", "mise"}
	if goos == "darwin" {
		bins = append([]string{"brew"}, bins...)
	}
	return bins
}

// missingTools lists the tools that the global mise configuration names and
// this machine does not have. It runs from $HOME so that mise does not read
// a mise.toml in the directory that started polkadot. It runs offline: a
// floating version such as "latest" then means the newest installed one, so
// a new upstream release does not turn a healthy machine red.
func missingTools(home string) ([]string, error) {
	cmd := exec.Command("mise", "ls", "--missing", "--json")
	cmd.Dir = home
	cmd.Env = append(os.Environ(), "MISE_OFFLINE=1")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return parseMissing(out)
}

func parseMissing(out []byte) ([]string, error) {
	tools := map[string]json.RawMessage{}
	if len(strings.TrimSpace(string(out))) > 0 {
		if err := json.Unmarshal(out, &tools); err != nil {
			return nil, err
		}
	}
	missing := make([]string, 0, len(tools))
	for name := range tools {
		missing = append(missing, name)
	}
	sort.Strings(missing)
	return missing, nil
}

// checkTools reports the base programs and the mise tools, and returns the
// number of failures. The tools step and doctor both use it.
func checkTools(home string) int {
	failed := 0
	for _, bin := range baseBins() {
		if step.Has(bin) {
			ui.Result(bin, "already", "")
		} else {
			ui.Result(bin, "failed", "missing — run the packages step")
			failed++
		}
	}
	missing, err := missingTools(home)
	switch {
	case err != nil:
		ui.Result("mise tools", "failed", "mise could not list them")
		failed++
	case len(missing) > 0:
		ui.Result("mise tools", "failed", "missing "+strings.Join(missing, " ")+" — run the runtimes step")
		failed++
	default:
		ui.Result("mise tools", "already", "all present")
	}
	return failed
}

func doTools(home string, t *tally) {
	ui.Section("tools")
	t.fail += checkTools(home)
	ui.Blank()
}

type agent struct{ bin, how string }

// agentInstallers lists each agent with the installer its vendor documents
// for this system. herdr is the terminal client; hc is herdr-core, the
// server that the desktop and mobile apps attach to.
func agentInstallers() []agent {
	opencode := "brew install sst/tap/opencode"
	if goos != "darwin" {
		opencode = "curl -fsSL https://opencode.ai/v2/install | bash"
	}
	return []agent{
		{"claude", "curl -fsSL https://claude.ai/install.sh | bash"},
		{"codex", "npm i -g @openai/codex"},
		{"opencode", opencode},
		{"agy", "curl -fsSL https://antigravity.google/cli/install.sh | bash"},
		{"herdr", "curl -fsSL https://herdr.dev/install.sh | sh"},
		{"hc", "curl -fsSL https://releases.herdr.dev/core/install.sh | sh"},
	}
}

func doAgents(t *tally) {
	ui.Section("coding agents")
	for _, a := range agentInstallers() {
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

// doDoctor reports the health of this machine and changes nothing. It
// returns the number of failures, so a script can act on the result.
func doDoctor(home, root string) int {
	failed := 0
	ui.Section("links")
	for _, m := range active(append(links, themeLinks...)) {
		dst := filepath.Join(home, m.to)
		want := filepath.Join(root, m.from)
		got, err := os.Readlink(dst)
		switch {
		case err == nil && got == want:
			ui.Result("~/"+m.to, "already", "linked")
		case err == nil:
			ui.Result("~/"+m.to, "failed", "points elsewhere: "+got)
			failed++
		default:
			if _, err := os.Lstat(dst); err == nil {
				ui.Result("~/"+m.to, "failed", "exists but is not a link")
				failed++
			} else {
				ui.Result("~/"+m.to, "skipped", "absent")
			}
		}
	}
	ui.Blank()

	ui.Section("tools")
	failed += checkTools(home)
	ui.Blank()

	ui.Section("agents")
	for _, a := range agentInstallers() {
		switch {
		case step.Runs(a.bin):
			ui.Result(a.bin, "already", "")
		case step.Has(a.bin):
			// On PATH but cannot start: a half-installed package, or a shim
			// whose real binary went missing. Worth saying loudly, because
			// everything else about it looks fine.
			ui.Result(a.bin, "failed", "on PATH but will not run — reinstall it")
			failed++
		default:
			ui.Result(a.bin, "skipped", "absent")
		}
	}
	ui.Blank()

	ui.Section("shell")
	ui.Result("login shell", "already", loginShell())
	ui.Blank()

	ui.Section("repository")
	ui.Result("version", "already", version(root))
	// A machine that writes into the checkout stops its own sync, because
	// sync skips a dirty tree. This is a warning and not a failure while
	// Codex and Herdr still write into files that polkadot links.
	if out, err := step.RunDir(root, "git", "status", "--porcelain", "--untracked-files=no"); err == nil {
		if n := len(strings.Split(out, "\n")); out != "" {
			ui.Result("working tree", "warning", fmt.Sprintf("%d changed files; sync skips a dirty tree", n))
		} else {
			ui.Result("working tree", "already", "clean")
		}
	}
	// Compared with the last fetch: doctor does not use the network.
	if out, err := step.RunDir(root, "git", "rev-list", "--count", "HEAD..@{u}"); err == nil {
		if out != "0" {
			ui.Result("upstream", "warning", out+" commits behind, as of the last fetch")
		} else {
			ui.Result("upstream", "already", "current, as of the last fetch")
		}
	}
	ui.Blank()
	return failed
}

// version names the commit that this checkout runs: the nearest tag, the
// distance from it, and a mark when the tree has local changes.
func version(root string) string {
	out, err := step.RunDir(root, "git", "describe", "--tags", "--always", "--dirty")
	if err != nil {
		return "unknown"
	}
	return out
}
