package link

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLinkBacksUpExistingTarget(t *testing.T) {
	home := t.TempDir()
	src := filepath.Join(t.TempDir(), "config.toml")
	if err := os.WriteFile(src, []byte("new"), 0o644); err != nil {
		t.Fatal(err)
	}

	target := filepath.Join(home, ".config", "tool", "config.toml")
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(target, []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}

	l := New("", home, false)
	status, err := l.Link(src, target)
	if err != nil {
		t.Fatal(err)
	}
	if status != Backed {
		t.Fatalf("status = %v, want Backed", status)
	}
	if !l.BackupUsed() {
		t.Fatal("BackupUsed() = false, want true")
	}
	if got, err := os.Readlink(target); err != nil || got != src {
		t.Fatalf("target = %q, %v; want symlink to %q", got, err, src)
	}

	backup := filepath.Join(l.BackupDir, ".config", "tool", "config.toml")
	got, err := os.ReadFile(backup)
	if err != nil || string(got) != "old" {
		t.Fatalf("backup = %q, %v; want old contents", got, err)
	}
}

func TestLinkDryRunDoesNotChangeTarget(t *testing.T) {
	home := t.TempDir()
	src := filepath.Join(t.TempDir(), "config.toml")
	if err := os.WriteFile(src, []byte("new"), 0o644); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(home, ".config", "config.toml")
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(target, []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}

	status, err := New("", home, true).Link(src, target)
	if err != nil {
		t.Fatal(err)
	}
	if status != Would {
		t.Fatalf("status = %v, want Would", status)
	}
	got, err := os.ReadFile(target)
	if err != nil || string(got) != "old" {
		t.Fatalf("target = %q, %v; want unchanged contents", got, err)
	}
}

func TestLinkRecognizesExistingCorrectLink(t *testing.T) {
	home := t.TempDir()
	src := filepath.Join(t.TempDir(), "config.toml")
	if err := os.WriteFile(src, []byte("new"), 0o644); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(home, ".config", "config.toml")
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(src, target); err != nil {
		t.Fatal(err)
	}

	status, err := New("", home, false).Link(src, target)
	if err != nil {
		t.Fatal(err)
	}
	if status != Already {
		t.Fatalf("status = %v, want Already", status)
	}
}

func TestUnlinkRemovesOnlyManagedLink(t *testing.T) {
	home := t.TempDir()
	src := filepath.Join(t.TempDir(), "config.toml")
	if err := os.WriteFile(src, []byte("new"), 0o644); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(home, ".config", "config.toml")
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(src, target); err != nil {
		t.Fatal(err)
	}

	managed, err := New("", home, false).Unlink(src, target)
	if err != nil {
		t.Fatal(err)
	}
	if !managed {
		t.Fatal("managed = false, want true")
	}
	if _, err := os.Lstat(target); !os.IsNotExist(err) {
		t.Fatalf("target still exists: %v", err)
	}
}

func TestUnlinkLeavesForeignTargetAlone(t *testing.T) {
	home := t.TempDir()
	src := filepath.Join(t.TempDir(), "config.toml")
	foreign := filepath.Join(t.TempDir(), "foreign.toml")
	for _, path := range []string{src, foreign} {
		if err := os.WriteFile(path, []byte("config"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	target := filepath.Join(home, ".config", "config.toml")
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(foreign, target); err != nil {
		t.Fatal(err)
	}

	managed, err := New("", home, false).Unlink(src, target)
	if err != nil {
		t.Fatal(err)
	}
	if managed {
		t.Fatal("managed = true, want false")
	}
	if got, err := os.Readlink(target); err != nil || got != foreign {
		t.Fatalf("target = %q, %v; want unchanged foreign link", got, err)
	}
}

// A retired link points at a source that the checkout no longer has. Unlink
// must still recognise it as managed, or the dangling link stays forever.
func TestUnlinkRemovesDanglingManagedLink(t *testing.T) {
	root, home := t.TempDir(), t.TempDir()
	src := filepath.Join(root, "home", "zshrc")
	target := filepath.Join(home, ".zshrc")
	if err := os.Symlink(src, target); err != nil {
		t.Fatal(err)
	}

	managed, err := New(root, home, false).Unlink(src, target)
	if err != nil {
		t.Fatal(err)
	}
	if !managed {
		t.Fatal("dangling link into the checkout was not recognised as managed")
	}
	if _, err := os.Lstat(target); !os.IsNotExist(err) {
		t.Fatalf("dangling link still present: %v", err)
	}
}

func TestSeedCopiesOnceAndNeverOverwrites(t *testing.T) {
	root, home := t.TempDir(), t.TempDir()
	src := filepath.Join(root, "settings.json")
	target := filepath.Join(home, ".claude", "settings.json")
	if err := os.WriteFile(src, []byte("from the repository"), 0o644); err != nil {
		t.Fatal(err)
	}
	l := New(root, home, false)

	if st, err := l.Seed(src, target); err != nil || st != Created {
		t.Fatalf("first seed = %v, %v", st, err)
	}
	if fi, _ := os.Lstat(target); fi.Mode()&os.ModeSymlink != 0 {
		t.Fatal("seed made a link; the tool would write into the checkout")
	}

	// The tool rewrites its file. A later seed must leave that alone.
	if err := os.WriteFile(target, []byte("written by the tool"), 0o644); err != nil {
		t.Fatal(err)
	}
	if st, err := l.Seed(src, target); err != nil || st != Already {
		t.Fatalf("second seed = %v, %v", st, err)
	}
	if got, _ := os.ReadFile(target); string(got) != "written by the tool" {
		t.Fatalf("seed overwrote the tool's file: %q", got)
	}
}

// An earlier version linked the file. The tool wrote through the link, so
// the checkout holds the machine's content. The migration must keep it.
func TestSeedReplacesItsOldLinkWithACopy(t *testing.T) {
	root, home := t.TempDir(), t.TempDir()
	src := filepath.Join(root, "config.toml")
	target := filepath.Join(home, "config.toml")
	if err := os.WriteFile(src, []byte("what the machine runs"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(src, target); err != nil {
		t.Fatal(err)
	}

	if st, err := New(root, home, false).Seed(src, target); err != nil || st != Created {
		t.Fatalf("seed over its old link = %v, %v", st, err)
	}
	fi, err := os.Lstat(target)
	if err != nil || fi.Mode()&os.ModeSymlink != 0 {
		t.Fatalf("target is still a link: %v", err)
	}
	if fi.Mode().Perm() != 0o600 {
		t.Fatalf("mode = %v, want the source's 0600", fi.Mode().Perm())
	}
	if got, _ := os.ReadFile(target); string(got) != "what the machine runs" {
		t.Fatalf("content lost: %q", got)
	}
}

func TestUnfoldMovesToolStateOutOfTheCheckout(t *testing.T) {
	root, home := t.TempDir(), t.TempDir()
	srcDir := filepath.Join(root, "config", "herdr")
	target := filepath.Join(home, ".config", "herdr")
	for name, body := range map[string]string{"config.toml": "tracked", "plugins.json": "state", "session.json": "state"} {
		if err := os.MkdirAll(srcDir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(srcDir, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.MkdirAll(filepath.Join(srcDir, "plugins", "one"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(srcDir, target); err != nil {
		t.Fatal(err)
	}

	managed, err := New(root, home, false).Unfold(srcDir, target, []string{"config.toml"})
	if err != nil || !managed {
		t.Fatalf("unfold = %v, %v", managed, err)
	}
	if fi, _ := os.Lstat(target); fi == nil || fi.Mode()&os.ModeSymlink != 0 || !fi.IsDir() {
		t.Fatal("target is not a real directory")
	}
	for _, moved := range []string{"plugins.json", "session.json", "plugins/one"} {
		if _, err := os.Stat(filepath.Join(target, moved)); err != nil {
			t.Errorf("%s did not move to the real directory", moved)
		}
		if _, err := os.Stat(filepath.Join(srcDir, moved)); err == nil {
			t.Errorf("%s is still in the checkout", moved)
		}
	}
	if _, err := os.Stat(filepath.Join(srcDir, "config.toml")); err != nil {
		t.Error("the tracked file left the checkout")
	}
}

func TestUnfoldLeavesAForeignTargetAlone(t *testing.T) {
	root, home := t.TempDir(), t.TempDir()
	srcDir := filepath.Join(root, "herdr")
	target := filepath.Join(home, "herdr")
	if err := os.MkdirAll(srcDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(target, 0o755); err != nil {
		t.Fatal(err)
	}
	if managed, err := New(root, home, false).Unfold(srcDir, target, nil); err != nil || managed {
		t.Fatalf("unfold touched a real directory: %v, %v", managed, err)
	}
}
