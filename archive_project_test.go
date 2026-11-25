package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

const archiveScript = "config/claude/skills/archive-project/scripts/archive-project"

// archiveBox is a throwaway home with one project, its remote, and a
// worktree that holds uncommitted work.
type archiveBox struct {
	t      *testing.T
	home   string
	code   string
	proj   string
	remote string
	script string
}

func newArchiveBox(t *testing.T) *archiveBox {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not installed")
	}
	script, err := filepath.Abs(archiveScript)
	if err != nil {
		t.Fatal(err)
	}
	home, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	b := &archiveBox{t: t, home: home, script: script}
	b.code = filepath.Join(home, "Code")
	b.proj = filepath.Join(b.code, "demo")
	b.remote = filepath.Join(home, "remote.git")

	b.git(home, "init", "--quiet", "--bare", "-b", "main", b.remote)
	b.git(home, "init", "--quiet", "-b", "main", b.proj)
	b.write("demo/.gitignore", ".env\ntarget/\n")
	b.write("demo/README.md", "demo\n")
	b.git(b.proj, "add", "-A")
	b.git(b.proj, "commit", "--quiet", "-m", "feat: start")
	return b
}

func (b *archiveBox) addRemote() {
	b.git(b.proj, "remote", "add", "origin", b.remote)
	b.git(b.proj, "push", "--quiet", "-u", "origin", "main")
}

func (b *archiveBox) env() []string {
	return append(os.Environ(),
		"HOME="+b.home,
		"ARCHIVE_CODE_DIR="+b.code,
		"GIT_CONFIG_GLOBAL="+os.DevNull,
		"GIT_CONFIG_NOSYSTEM=1",
		"GIT_AUTHOR_NAME=Test", "GIT_AUTHOR_EMAIL=test@example.com",
		"GIT_COMMITTER_NAME=Test", "GIT_COMMITTER_EMAIL=test@example.com",
	)
}

func (b *archiveBox) git(dir string, args ...string) string {
	b.t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir, cmd.Env = dir, b.env()
	out, err := cmd.CombinedOutput()
	if err != nil {
		b.t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
	return string(out)
}

func (b *archiveBox) write(rel, body string) {
	b.t.Helper()
	path := filepath.Join(b.code, rel)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		b.t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		b.t.Fatal(err)
	}
}

// archive runs the script from the home directory, outside the project.
func (b *archiveBox) archive(args ...string) (string, error) {
	cmd := exec.Command(b.script, args...)
	cmd.Dir, cmd.Env = b.home, b.env()
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func (b *archiveBox) exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func TestArchiveProjectSavesEverythingBeforeItRemoves(t *testing.T) {
	b := newArchiveBox(t)
	b.addRemote()

	// Work in four places that a careless delete would lose: a commit that
	// is not pushed, an edit and a new file in a worktree, and an ignored
	// file. The build output is the one thing that may go.
	b.write("demo/README.md", "demo, second line\n")
	b.git(b.proj, "commit", "--quiet", "-am", "docs: second line")
	worktree := filepath.Join(b.code, ".worktrees", "demo", "feat-x")
	b.git(b.proj, "worktree", "add", "--quiet", "-b", "feat/x", worktree)
	b.write(".worktrees/demo/feat-x/README.md", "edited in the worktree\n")
	b.write(".worktrees/demo/feat-x/new.txt", "never committed\n")
	b.write("demo/.env", "SECRET=1\n")
	b.write("demo/target/out.bin", "build output\n")

	out, err := b.archive("audit", "demo")
	if err != nil {
		t.Fatalf("audit failed: %v\n%s", err, out)
	}
	for _, want := range []string{"uncommitted: 2", "demo/.env", "commits not on the remote: 1"} {
		if !strings.Contains(out, want) {
			t.Errorf("audit output lacks %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "target") {
		t.Errorf("audit lists build output as a file to keep:\n%s", out)
	}

	out, err = b.archive("run", "demo")
	if err == nil || !strings.Contains(out, "--snapshot") {
		t.Fatalf("run accepted uncommitted work without --snapshot: %v\n%s", err, out)
	}
	if !b.exists(worktree) || !b.exists(b.proj) || b.exists(filepath.Join(b.code, ".archive")) {
		t.Fatal("a refused run changed the disk")
	}

	if out, err = b.archive("run", "demo", "--snapshot"); err != nil {
		t.Fatalf("run failed: %v\n%s", err, out)
	}
	if b.exists(b.proj) || b.exists(filepath.Join(b.code, ".worktrees", "demo")) {
		t.Fatal("run left the project or its worktrees on disk")
	}

	if got := b.git(b.home, "--git-dir", b.remote, "show", "main:README.md"); got != "demo, second line\n" {
		t.Errorf("the remote lacks the unpushed commit: %q", got)
	}
	snapshot := "archive/wip/feat/x"
	if got := b.git(b.home, "--git-dir", b.remote, "show", snapshot+":new.txt"); got != "never committed\n" {
		t.Errorf("the snapshot lacks the untracked file: %q", got)
	}
	if got := b.git(b.home, "--git-dir", b.remote, "show", snapshot+":README.md"); got != "edited in the worktree\n" {
		t.Errorf("the snapshot lacks the edit: %q", got)
	}
	if got := b.git(b.home, "--git-dir", b.remote, "show", "feat/x:README.md"); got != "demo, second line\n" {
		t.Errorf("the snapshot changed the branch itself: %q", got)
	}

	cold := filepath.Join(b.code, ".archive", "demo")
	b.git(b.home, "--git-dir", b.remote, "bundle", "verify", "--quiet", filepath.Join(cold, "repo.bundle"))
	list, err := exec.Command("tar", "-tf", filepath.Join(cold, "ignored-1.tar")).CombinedOutput()
	if err != nil {
		t.Fatalf("the cold copy has no readable ignored-1.tar: %v\n%s", err, list)
	}
	if !strings.Contains(string(list), ".env") || strings.Contains(string(list), "target") {
		t.Errorf("ignored-1.tar holds the wrong files:\n%s", list)
	}

	out, err = b.archive("run", "demo")
	if err != nil || !strings.Contains(out, "already archived") {
		t.Fatalf("a second run is not a no-op: %v\n%s", err, out)
	}

	if out, err = b.archive("restore", "demo"); err != nil {
		t.Fatalf("restore failed: %v\n%s", err, out)
	}
	env, err := os.ReadFile(filepath.Join(b.proj, ".env"))
	if err != nil || string(env) != "SECRET=1\n" {
		t.Errorf("restore did not put back .env: %v %q", err, env)
	}
	if got := b.git(b.proj, "show", "origin/"+snapshot+":new.txt"); got != "never committed\n" {
		t.Errorf("the restored clone lacks the snapshot: %q", got)
	}
}

func TestArchiveProjectRefusesWithoutARemote(t *testing.T) {
	b := newArchiveBox(t)

	out, err := b.archive("audit", "demo")
	if err == nil || !strings.Contains(out, "BLOCKER: no remote") {
		t.Fatalf("audit did not report the missing remote: %v\n%s", err, out)
	}
	out, err = b.archive("run", "demo", "--snapshot", "--kill")
	if err == nil {
		t.Fatalf("run archived a project that exists nowhere else:\n%s", out)
	}
	if !b.exists(b.proj) {
		t.Fatal("a refused run removed the project")
	}
}

func TestArchiveProjectStopsProcessesOnlyWhenTold(t *testing.T) {
	if _, err := os.Stat("/proc/self/cwd"); err != nil {
		t.Skip("the process check needs /proc")
	}
	b := newArchiveBox(t)
	b.addRemote()

	inside := exec.Command("sleep", "300")
	inside.Dir = b.proj
	if err := inside.Start(); err != nil {
		t.Fatal(err)
	}
	exited := make(chan struct{})
	go func() { _ = inside.Wait(); close(exited) }()
	t.Cleanup(func() { _ = inside.Process.Kill(); <-exited })

	// A process that works somewhere else must survive.
	outside := exec.Command("sleep", "300")
	outside.Dir = b.home
	if err := outside.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = outside.Process.Kill(); _ = outside.Wait() })

	out, err := b.archive("run", "demo")
	if err == nil || !strings.Contains(out, "--kill") {
		t.Fatalf("run removed a project with a live process: %v\n%s", err, out)
	}
	if !b.exists(b.proj) {
		t.Fatal("a refused run removed the project")
	}

	if out, err = b.archive("run", "demo", "--kill"); err != nil {
		t.Fatalf("run --kill failed: %v\n%s", err, out)
	}
	<-exited
	if _, err := os.Stat("/proc/" + strconv.Itoa(outside.Process.Pid) + "/cwd"); err != nil {
		t.Fatalf("run --kill stopped a process outside the project: %v", err)
	}
}
