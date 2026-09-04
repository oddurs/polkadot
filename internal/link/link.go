// Package link is the symlink engine.
//
// The rule: the repo is the source of truth. Anything already at a target path
// that is not our symlink gets moved into a timestamped backup directory before
// we replace it — this tool never destroys a config it did not write.
package link

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Status int

const (
	Created Status = iota // a new symlink
	Already               // correct symlink already in place
	Backed                // something was there; it was moved aside first
	Would                 // dry run
)

type Linker struct {
	Root      string // the polkadot checkout
	Home      string
	DryRun    bool
	BackupDir string
	backedUp  bool
}

func New(root, home string, dry bool) *Linker {
	return &Linker{
		Root:      root,
		Home:      home,
		DryRun:    dry,
		BackupDir: filepath.Join(home, ".polkadot-backup", time.Now().Format("2006-01-02-150405")),
	}
}

// BackupUsed reports whether anything was actually moved aside this run.
func (l *Linker) BackupUsed() bool { return l.backedUp }

// Link points target at src. Both are absolute.
func (l *Linker) Link(src, target string) (Status, error) {
	if _, err := os.Stat(src); err != nil {
		return Would, fmt.Errorf("source missing: %s", src)
	}

	if cur, err := os.Readlink(target); err == nil {
		if cur == src {
			return Already, nil
		}
	}

	_, err := os.Lstat(target)
	exists := err == nil

	if l.DryRun {
		return Would, nil
	}

	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return Would, err
	}

	status := Created
	if exists {
		// Move it aside rather than delete it, preserving the path shape so a
		// restore is a plain `mv` back.
		rel, err := filepath.Rel(l.Home, target)
		if err != nil {
			rel = filepath.Base(target)
		}
		dest := filepath.Join(l.BackupDir, rel)
		if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
			return Would, err
		}
		if err := os.Rename(target, dest); err != nil {
			return Would, fmt.Errorf("backing up %s: %w", target, err)
		}
		l.backedUp = true
		status = Backed
	}

	if err := os.Symlink(src, target); err != nil {
		return Would, err
	}
	return status, nil
}
