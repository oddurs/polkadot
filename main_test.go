package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
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
// terminal, the desktop editor, or the macOS git settings.
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

// gitConfig asks git itself what one config file holds, so the tests read
// each form that git reads: a subsection, a quoted value, a comment after a
// value. With --file, git does not follow an include. git exits with 1 when
// no key matches: that is an empty answer, not an error.
func gitConfig(t *testing.T, path string, args ...string) []string {
	t.Helper()
	out, err := exec.Command("git", append([]string{"config", "--file", path}, args...)...).Output()
	var exit *exec.ExitError
	if errors.As(err, &exit) && exit.ExitCode() == 1 {
		return nil
	}
	if err != nil {
		t.Fatalf("git config --file %s %v: %v", path, args, err)
	}
	return strings.Split(strings.TrimSpace(string(out)), "\n")
}

// ~/.gitconfig is one file for both systems. A macOS-only setting in it
// reaches the Linux machine too. Linux has no credential-osxkeychain, so
// each HTTPS credential lookup printed "'credential-osxkeychain' is not a
// git command". The fsmonitor daemon has a macOS and a Windows build only.
// The pattern also matches a helper for one URL: credential.<url>.helper.
func TestSharedGitconfigHoldsNoMacOnlySetting(t *testing.T) {
	for _, line := range gitConfig(t, "home/gitconfig", "--get-regexp", `^(credential\..*helper|core\.fsmonitor)$`) {
		t.Errorf("home/gitconfig sets %s: move it to home/gitconfig.darwin", line)
	}
	for _, name := range []string{"credential.helper", "core.fsmonitor"} {
		if gitConfig(t, "home/gitconfig.darwin", "--get", name) == nil {
			t.Errorf("home/gitconfig.darwin does not set %s", name)
		}
	}
}

// Three parts carry the macOS git settings: the link on macOS, no link on
// Linux, and the include in the shared file. The include must come before
// ~/.gitconfig.local, because the machine's own file has the last word.
func TestMacGitconfigIsLinkedOnDarwinOnly(t *testing.T) {
	linked := func() bool {
		for _, m := range active(links) {
			if m.from == "home/gitconfig.darwin" && m.to == ".gitconfig.darwin" {
				return true
			}
		}
		return false
	}
	withGOOS(t, "darwin")
	if !linked() {
		t.Error("macOS does not link home/gitconfig.darwin to ~/.gitconfig.darwin")
	}
	withGOOS(t, "linux")
	if linked() {
		t.Error("linux links home/gitconfig.darwin")
	}

	mac, local := -1, -1
	for i, path := range gitConfig(t, "home/gitconfig", "--get-all", "include.path") {
		switch path {
		case "~/.gitconfig.darwin":
			mac = i
		case "~/.gitconfig.local":
			local = i
		}
	}
	if mac < 0 {
		t.Fatal("home/gitconfig does not include ~/.gitconfig.darwin")
	}
	if local < mac {
		t.Error("home/gitconfig must include ~/.gitconfig.local after ~/.gitconfig.darwin")
	}
}

// A Mac that gets this checkout with a plain `git pull` has the new
// ~/.gitconfig at once, through the existing link, and no ~/.gitconfig.darwin
// until `polkadot link` runs. Doctor must report that state as a failure.
func TestDoctorFailsWithoutTheMacGitconfigLink(t *testing.T) {
	root, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	home := t.TempDir()
	withGOOS(t, "darwin")
	if got := doctorLinks(home, root); got != 1 {
		t.Fatalf("doctor counts %d failures in a home without ~/.gitconfig.darwin, want 1", got)
	}
	if err := os.Symlink(filepath.Join(root, "home/gitconfig.darwin"), filepath.Join(home, ".gitconfig.darwin")); err != nil {
		t.Fatal(err)
	}
	if got := doctorLinks(home, root); got != 0 {
		t.Fatalf("doctor counts %d failures with the link in place, want 0", got)
	}
	withGOOS(t, "linux")
	if got := doctorLinks(t.TempDir(), root); got != 0 {
		t.Fatalf("doctor counts %d failures on linux, want 0", got)
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
	// One installer for both systems: two installers gave two major versions.
	withGOOS(t, "darwin")
	mac := find("opencode")
	withGOOS(t, "linux")
	if find("opencode") != mac || strings.Contains(mac, "brew") {
		t.Fatalf("opencode installers differ or use Homebrew: %q", mac)
	}
	// The installer must not write to a shell file: polkadot owns PATH, and
	// the installer would bring back the ~/.zshrc that polkadot retired.
	if !strings.Contains(mac, "--no-modify-path") {
		t.Fatalf("the opencode installer may edit shell files: %q", mac)
	}
	for _, bin := range []string{"hc", "agy"} {
		if find(bin) == "" {
			t.Fatalf("no installer for %s", bin)
		}
	}
}

func TestSyncServiceQuotesPaths(t *testing.T) {
	body := buildSyncService("/home/a b/polkadot", "/usr/bin/go")
	// systemd rejects a quoted WorkingDirectory with "path is not absolute".
	if !strings.Contains(body, "WorkingDirectory=/home/a b/polkadot\n") {
		t.Fatalf("working directory must be the bare path:\n%s", body)
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
	for _, m := range append(append(links, themeLinks...), seeds...) {
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
	got, err := parseMissing([]byte(`{
		"ripgrep":[{"version":"15.2.0","source":{"type":"mise.toml","path":"/home/a/.config/mise/config.toml"}}],
		"fd":[{"version":"10.5.2","source":{"type":"mise.toml","path":"/home/a/.config/mise/config.toml"}}],
		"node":[{"version":"26.13.1","source":null}]
	}`))
	if err != nil {
		t.Fatal(err)
	}
	// node has no source: no configuration asks for that version.
	if strings.Join(got, ",") != "fd,ripgrep" {
		t.Fatalf("parseMissing = %v", got)
	}
	for _, empty := range []string{"", "{}", "\n"} {
		if got, err := parseMissing([]byte(empty)); err != nil || len(got) != 0 {
			t.Fatalf("parseMissing(%q) = %v, %v", empty, got, err)
		}
	}
}

// sysctl reads one setting per line. A body without a real line break, or
// with the two characters backslash and n, is rejected as an invalid value.
func TestInotifyBodyIsOneValidLine(t *testing.T) {
	if !strings.HasSuffix(inotifyBody, "\n") || strings.Count(inotifyBody, "\n") != 1 {
		t.Fatalf("inotify body must end with one line break: %q", inotifyBody)
	}
	if strings.Contains(inotifyBody, `\n`) {
		t.Fatalf("inotify body holds a literal backslash-n: %q", inotifyBody)
	}
}

func TestWriteRootFileKeepsTheBody(t *testing.T) {
	// The debian CI job runs the tests as root. A developer machine would
	// need a sudo password here.
	if os.Getuid() != 0 {
		t.Skip("needs root")
	}
	path := t.TempDir() + "/60-test.conf"
	if err := writeRootFile(path, inotifyBody); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != inotifyBody {
		t.Fatalf("file holds %q, want %q", got, inotifyBody)
	}
}

// A file is linked or seeded, never both: a link to a seeded file would
// send the tool's writes back into the checkout.
func TestNoFileIsBothLinkedAndSeeded(t *testing.T) {
	linked := map[string]bool{}
	for _, m := range append(links, themeLinks...) {
		linked[m.to] = true
	}
	for _, m := range seeds {
		if linked[m.to] {
			t.Errorf("%s is linked and seeded", m.to)
		}
		for to := range linked {
			if strings.HasPrefix(m.to, to+"/") {
				t.Errorf("%s is seeded inside the linked folder %s", m.to, to)
			}
		}
	}
}

// "lts" and "latest" name a different version on different days, so two
// machines set up a month apart would run different tools. Node moved from
// 24 to 26 that way. go stays on "latest" on purpose: polkadot itself needs
// a current Go, and Go keeps old programs building.
func TestMiseConfigPinsNode(t *testing.T) {
	data, err := os.ReadFile("config/mise/config.toml")
	if err != nil {
		t.Fatal(err)
	}
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "node = ") && (strings.Contains(line, "lts") || strings.Contains(line, "latest")) {
			t.Fatalf("node follows a moving name: %s", line)
		}
	}
}

// herdr-core passes its environment to every terminal. An installer run with
// NO_COLOR=1 once started it, and every agent on that machine then ran
// without colour, under bash, with no node on PATH.
func TestHerdrCoreStartsFromACleanLogin(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	args := cleanStart([]string{"HOME=/home/a", "LANG=en_US.UTF-8"})
	if args[0] != "-i" {
		t.Fatalf("env must start empty, got %v", args)
	}
	joined := strings.Join(args, " ")
	if strings.Contains(joined, "NO_COLOR") {
		t.Fatalf("NO_COLOR reached herdr-core: %s", joined)
	}
	if !strings.Contains(joined, "SHELL=") || !strings.HasSuffix(joined, "fish -l -c hc server start") {
		t.Fatalf("herdr-core must start from a fish login shell: %s", joined)
	}
}

// The server prompt repeats the modules of the main one, because starship
// has no include. This keeps the copy honest: outside the parts that give
// the server its look, the two files are the same text.
func TestServerPromptDiffersOnlyInItsLook(t *testing.T) {
	look := map[string]bool{"[character]": true, "[username]": true, "[hostname]": true}
	body := func(path string) string {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		var out []string
		skip := false
		for _, line := range strings.Split(string(data), "\n") {
			trimmed := strings.TrimSpace(line)
			switch {
			case strings.HasPrefix(trimmed, "[palettes."):
				return strings.Join(out, "\n")
			case strings.HasPrefix(trimmed, "["):
				skip = look[trimmed]
			}
			if skip || trimmed == "" || strings.HasPrefix(trimmed, "#") ||
				strings.HasPrefix(trimmed, "palette = ") || strings.HasPrefix(trimmed, "$username") {
				continue
			}
			out = append(out, trimmed)
		}
		return strings.Join(out, "\n")
	}
	if a, b := body("config/starship.toml"), body("config/starship-server.toml"); a != b {
		t.Fatal("config/starship-server.toml has drifted from config/starship.toml outside its look")
	}
}

// A seed is copied to a new machine as it is. A path of one machine in it is
// useless there, and this repository is public: the path names a project.
func TestSeedsHoldNoMachinePath(t *testing.T) {
	for _, m := range seeds {
		data, err := os.ReadFile(m.from)
		if err != nil {
			t.Fatal(err)
		}
		for _, bad := range []string{"/Users/", "/home/"} {
			if strings.Contains(string(data), bad) {
				t.Errorf("%s holds a path under %s", m.from, bad)
			}
		}
	}
}

func TestJSONSeedsParse(t *testing.T) {
	for _, m := range seeds {
		if !strings.HasSuffix(m.from, ".json") {
			continue
		}
		data, err := os.ReadFile(m.from)
		if err != nil {
			t.Fatal(err)
		}
		if err := json.Unmarshal(data, new(any)); err != nil {
			t.Errorf("%s is not valid JSON: %v", m.from, err)
		}
	}
}
