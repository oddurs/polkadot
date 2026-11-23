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
