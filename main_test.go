package main

import (
	"encoding/xml"
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
