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
