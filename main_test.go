package main

import (
	"encoding/xml"
	"os"
	"strings"
	"testing"
)

func TestValidateArgs(t *testing.T) {
	original := *only
	t.Cleanup(func() { *only = original })

	*only = "not-a-step"
	if err := validateArgs("install"); err == nil {
		t.Fatal("validateArgs accepted an unknown step")
	}

	*only = "theme"
	if err := validateArgs("sync"); err == nil {
		t.Fatal("validateArgs accepted --only with sync")
	}

	*only = "theme"
	if err := validateArgs("install"); err != nil {
		t.Fatalf("validateArgs rejected a valid step: %v", err)
	}
}

func TestManagedReposUsesConfiguredRoot(t *testing.T) {
	root := "/Volumes/Work Config/polkadot"
	repos := managedRepos(root, "/Users/example")
	if repos[0].path != root {
		t.Fatalf("polkadot path = %q, want %q", repos[0].path, root)
	}
}

func TestSyncPlistAvoidsShellInterpolation(t *testing.T) {
	if strings.Contains(syncPlist, "/bin/sh") || strings.Contains(syncPlist, "cd %s") {
		t.Fatal("sync plist still uses a shell")
	}
	if !strings.Contains(syncPlist, "WorkingDirectory") {
		t.Fatal("sync plist has no working directory")
	}
}

func TestBuildSyncPlistEscapesPaths(t *testing.T) {
	body := buildSyncPlist("/Volumes/Work & Play/polkadot", "/opt/go&bin/go", "/Users/A & B")
	if !strings.Contains(body, "Work &amp; Play") || !strings.Contains(body, "go&amp;bin") {
		t.Fatal("sync plist did not escape XML-sensitive paths")
	}
	if err := xml.Unmarshal([]byte(body), new(any)); err != nil {
		t.Fatalf("sync plist is invalid XML: %v", err)
	}
}

func withGOOS(t *testing.T, value string) {
	t.Helper()
	original := goos
	goos = value
	t.Cleanup(func() { goos = original })
}

// A Linux machine is a server. It must not get a link for the desktop
// terminal or the desktop editor, whose targets do not exist there.
func TestLinuxSkipsDesktopLinks(t *testing.T) {
	withGOOS(t, "linux")
	for _, m := range active(append(links, themeLinks...)) {
		if darwinOnly[m.from] {
			t.Fatalf("linux links %s", m.from)
		}
		if strings.HasPrefix(m.to, "Library/") {
			t.Fatalf("linux links into ~/Library: %s", m.to)
		}
	}
	withGOOS(t, "darwin")
	if got, want := len(active(links)), len(links); got != want {
		t.Fatalf("darwin links %d of %d mappings", got, want)
	}
}

func TestDarwinOnlyNamesRealSources(t *testing.T) {
	known := map[string]bool{}
	for _, m := range append(links, themeLinks...) {
		known[m.from] = true
	}
	for from := range darwinOnly {
		if !known[from] {
			t.Fatalf("darwinOnly names %s, which no mapping links", from)
		}
	}
}

func TestLoginShellParsers(t *testing.T) {
	if got := shellFromPasswd("oddur:x:1000:1000:,,,:/home/oddur:/usr/bin/fish\n"); got != "/usr/bin/fish" {
		t.Fatalf("passwd shell = %q", got)
	}
	if got := shellFromPasswd("garbage"); got != "" {
		t.Fatalf("short passwd line gave %q", got)
	}
	if got := shellFromDscl("UserShell: /opt/homebrew/bin/fish\n"); got != "/opt/homebrew/bin/fish" {
		t.Fatalf("dscl shell = %q", got)
	}
}

func TestAgentInstallersPerSystem(t *testing.T) {
	find := func(bin string) string {
		for _, a := range agentInstallers() {
			if a.bin == bin {
				return a.how
			}
		}
		return ""
	}
	withGOOS(t, "linux")
	if strings.Contains(find("opencode"), "brew") {
		t.Fatal("linux installs opencode with Homebrew")
	}
	for _, bin := range []string{"hc", "agy"} {
		if find(bin) == "" {
			t.Fatalf("no installer for %s", bin)
		}
	}
}

func TestSyncServiceQuotesPaths(t *testing.T) {
	body := buildSyncService("/home/a b/polkadot", "/usr/bin/go")
	if !strings.Contains(body, `WorkingDirectory="/home/a b/polkadot"`) {
		t.Fatalf("working directory is not quoted:\n%s", body)
	}
	if !strings.Contains(body, `ExecStart="/usr/bin/go" run . sync`) {
		t.Fatalf("ExecStart is wrong:\n%s", body)
	}
}

func TestReadListSkipsCommentsAndBlanks(t *testing.T) {
	path := t.TempDir() + "/list"
	if err := os.WriteFile(path, []byte("# why\n\nfish\n  git  \n"), 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := readList(path)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(got, ",") != "fish,git" {
		t.Fatalf("readList = %v", got)
	}
}

// A link whose source is missing is skipped without an error, so a renamed
// file would silently stop being linked. gitconfig.local is the exception:
// it is ignored by git and each machine creates its own.
func TestEveryLinkSourceExists(t *testing.T) {
	for _, m := range append(links, themeLinks...) {
		if m.from == "home/gitconfig.local" {
			continue
		}
		if _, err := os.Stat(m.from); err != nil {
			t.Errorf("link source %s does not exist", m.from)
		}
	}
}

// One tool, one owner. A tool in both lists installs twice on a Mac, and the
// two copies drift apart.
func TestBrewfileAndMiseShareNoTool(t *testing.T) {
	brew, err := os.ReadFile("packages/darwin.Brewfile")
	if err != nil {
		t.Fatal(err)
	}
	mise, err := os.ReadFile("config/mise/config.toml")
	if err != nil {
		t.Fatal(err)
	}
	// mise itself is the installer, and go runs polkadot before mise exists.
	bootstrap := map[string]bool{"mise": true, "go": true}
	for _, line := range strings.Split(string(brew), "\n") {
		if !strings.HasPrefix(line, "brew \"") {
			continue
		}
		name := strings.Split(line, "\"")[1]
		if bootstrap[name] {
			continue
		}
		for _, entry := range strings.Split(string(mise), "\n") {
			if strings.HasPrefix(entry, name+" = \"") {
				t.Errorf("%s is in the Brewfile and unrestricted in the mise configuration", name)
			}
		}
	}
}

func TestEveryAgentReadsTheSharedInstructions(t *testing.T) {
	want := map[string]bool{
		".claude/CLAUDE.md": false, ".codex/AGENTS.md": false,
		".config/opencode/AGENTS.md": false, ".gemini/AGENTS.md": false,
	}
	for _, m := range links {
		if _, ok := want[m.to]; ok {
			if m.from != "agents/AGENTS.md" {
				t.Errorf("%s links %s, not the shared file", m.to, m.from)
			}
			want[m.to] = true
		}
	}
	for to, seen := range want {
		if !seen {
			t.Errorf("no link for %s", to)
		}
	}
}

func TestParseMissing(t *testing.T) {
	got, err := parseMissing([]byte(`{"ripgrep":[{"version":"15.2.0"}],"fd":[{"version":"10.5.2"}]}`))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(got, ",") != "fd,ripgrep" {
		t.Fatalf("parseMissing = %v", got)
	}
	for _, empty := range []string{"", "{}", "\n"} {
		if got, err := parseMissing([]byte(empty)); err != nil || len(got) != 0 {
			t.Fatalf("parseMissing(%q) = %v, %v", empty, got, err)
		}
	}
}
