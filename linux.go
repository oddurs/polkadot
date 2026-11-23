package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/oddurs/polkadot/internal/step"
	"github.com/oddurs/polkadot/internal/ui"
)

// doPackages installs what the system package manager owns: a short
// Brewfile on macOS, a short apt list and the mise binary on Linux. mise
// then installs the command-line tools in the runtimes step.
func doPackages(t *tally) {
	if goos == "darwin" {
		doBrew(t)
		return
	}
	doApt(t)
	doMise(t)
}

func doApt(t *tally) {
	ui.Section("apt")
	defer ui.Blank()
	if !step.Has("apt-get") {
		ui.Result("apt-get", "failed", "not found: the Linux target supports Debian and Ubuntu")
		t.fail++
		return
	}
	root, _ := repoRoot()
	pkgs, err := readList(filepath.Join(root, "packages", "linux.apt"))
	if err != nil {
		ui.Result("packages/linux.apt", "failed", err.Error())
		t.fail++
		return
	}
	var missing []string
	for _, p := range pkgs {
		if _, err := step.Run("dpkg", "-s", p); err != nil {
			missing = append(missing, p)
		}
	}
	if len(missing) == 0 {
		ui.Result("packages/linux.apt", "already", "all present")
		t.skip++
		return
	}
	if *dryRun {
		ui.Result("packages/linux.apt", "would", "install "+strings.Join(missing, " "))
		t.done++
		return
	}
	ui.Note("installing " + strings.Join(missing, " ") + "…")
	line := fmt.Sprintf("%[1]sapt-get update -qq && %[1]senv DEBIAN_FRONTEND=noninteractive apt-get install -y %[2]s",
		sudo(), strings.Join(missing, " "))
	if err := step.Stream("/bin/sh", "-c", line); err != nil {
		ui.Result("packages/linux.apt", "failed", err.Error())
		t.fail++
		return
	}
	ui.Result("packages/linux.apt", "installed", "")
	t.done++
}

func doMise(t *tally) {
	ui.Section("mise")
	defer ui.Blank()
	if step.Has("mise") {
		ui.Result("mise", "already", "")
		t.skip++
		return
	}
	installer := "curl -fsSL https://mise.run | sh"
	if *dryRun {
		ui.Result("mise", "would", installer)
		t.done++
		return
	}
	if !confirmInstaller("mise", installer) {
		ui.Result("mise", "skipped", "installer not approved")
		t.skip++
		return
	}
	if err := step.Stream("/bin/sh", "-c", installer); err != nil {
		ui.Result("mise", "failed", err.Error())
		t.fail++
		return
	}
	ui.Result("mise", "installed", "")
	t.done++
}

// doRuntimes asks mise to install every tool in the linked global
// configuration. It runs after the link step, because the link step puts
// that configuration in place.
func doRuntimes(home string, t *tally) {
	ui.Section("runtimes")
	defer ui.Blank()
	if !step.Has("mise") {
		ui.Result("mise", "failed", "missing: run the packages step first")
		t.fail++
		return
	}
	// Run from $HOME so that mise reads the global configuration and not a
	// mise.toml in the directory that started polkadot.
	if missing, err := missingTools(home); err == nil && len(missing) == 0 {
		ui.Result("mise install", "already", "all present")
		t.skip++
		return
	}
	if *dryRun {
		ui.Result("mise install", "would", "install the tools in ~/.config/mise/config.toml")
		t.done++
		return
	}
	if err := step.StreamDir(home, "mise", "install", "--yes"); err != nil {
		ui.Result("mise install", "failed", err.Error())
		t.fail++
		return
	}
	_, _ = step.RunDir(home, "mise", "reshim")
	ui.Result("mise install", "installed", "")
	t.done++
}

func hcRunning() bool {
	_, err := step.Run("hc", "server", "status")
	return err == nil
}

func systemdRunning() bool {
	_, err := os.Stat("/run/systemd/system")
	return err == nil
}

const inotifyConf = "/etc/sysctl.d/60-polkadot-inotify.conf"

// The default limit is too low for a file watcher on a large repository:
// a dev server stops with ENOSPC although the disk has space.
const inotifyBody = "fs.inotify.max_user_watches = 1048576\n"

// doSystem prepares a Linux server for sessions that outlive a login. It is
// the one place where polkadot changes the system outside $HOME, and each
// change is listed here.
func doSystem(t *tally) {
	if goos == "darwin" {
		return
	}
	ui.Section("system")
	defer ui.Blank()

	if got, _ := os.ReadFile(inotifyConf); string(got) == inotifyBody {
		ui.Result("inotify limit", "already", "1048576")
		t.skip++
	} else if *dryRun {
		ui.Result("inotify limit", "would", "write "+inotifyConf)
		t.done++
	} else if err := writeRootFile(inotifyConf, inotifyBody); err != nil {
		ui.Result("inotify limit", "failed", err.Error())
		t.fail++
	} else if _, err := step.Sh(sudo() + "sysctl -q -p " + inotifyConf); err != nil {
		// The file is in place for the next boot. A container shares the
		// kernel of its host and may not change the live value.
		ui.Result("inotify limit", "skipped", "written, but this system does not allow the live change")
		t.skip++
	} else {
		ui.Result("inotify limit", "set", "1048576")
		t.done++
	}

	if !systemdRunning() {
		ui.Result("lingering", "skipped", "systemd is not running here")
		ui.Result("herdr-core", "skipped", "systemd is not running here")
		t.skip += 2
		return
	}

	// Lingering starts the user's services at boot, without a login. The
	// sync timer and herdr-core both need it on a server.
	name := "user"
	if u, err := user.Current(); err == nil {
		name = u.Username
	}
	if out, _ := step.Run("loginctl", "show-user", name, "-p", "Linger"); strings.TrimSpace(out) == "Linger=yes" {
		ui.Result("lingering", "already", name)
		t.skip++
	} else if *dryRun {
		ui.Result("lingering", "would", "loginctl enable-linger "+name)
		t.done++
	} else if _, err := step.Sh(sudo() + "loginctl enable-linger " + name); err != nil {
		ui.Result("lingering", "failed", "run by hand: sudo loginctl enable-linger "+name)
		t.fail++
	} else {
		ui.Result("lingering", "enabled", name)
		t.done++
	}

	// herdr-core registers itself with the systemd user manager when it
	// starts, so polkadot writes no unit for it.
	switch {
	case !step.Has("hc"):
		ui.Result("herdr-core", "skipped", "not installed: run the agents step")
		t.skip++
	case *dryRun:
		ui.Result("herdr-core", "would", "hc server start")
		t.done++
	case hcRunning():
		ui.Result("herdr-core", "already", "running")
		t.skip++
	default:
		if _, err := step.Run("hc", "server", "start"); err != nil {
			ui.Result("herdr-core", "failed", "hc server start: "+err.Error())
			t.fail++
		} else {
			ui.Result("herdr-core", "started", "")
			t.done++
		}
	}
}

const syncService = `[Unit]
Description=polkadot sync

[Service]
Type=oneshot
WorkingDirectory=%s
ExecStart=%s run . sync
`

// Weekly, the same as the launchd agent on macOS.
const syncTimer = `[Unit]
Description=polkadot sync, weekly

[Timer]
OnCalendar=Mon 09:00
Persistent=true

[Install]
WantedBy=timers.target
`

// systemdQuote quotes one word of a command line in a unit file, where a
// space would otherwise split the path into two words.
func systemdQuote(s string) string {
	return `"` + strings.NewReplacer(`\`, `\\`, `"`, `\"`, `%`, `%%`).Replace(s) + `"`
}

// buildSyncService fills the unit. WorkingDirectory takes the rest of its
// line as one path and rejects a quoted value, so only % is escaped there.
// ExecStart is a command line, so its program path is quoted.
func buildSyncService(root, goBin string) string {
	return fmt.Sprintf(syncService, strings.ReplaceAll(root, "%", "%%"), systemdQuote(goBin))
}

// writeRootFile puts body at a path that only root may write. The body goes
// through a file and not through a shell string: a shell does not turn the
// two characters \n into a line break, and sysctl then rejects the value.
func writeRootFile(path, body string) error {
	tmp, err := os.CreateTemp("", "polkadot-*")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.WriteString(body); err != nil {
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if out, err := step.Sh(fmt.Sprintf("%sinstall -m 0644 %q %q", sudo(), tmp.Name(), path)); err != nil {
		return fmt.Errorf("%s: %s", err, out)
	}
	return nil
}

func doTimerLinux(root, goBin, home string, t *tally) {
	if *dryRun {
		ui.Result("polkadot-sync.timer", "would", "install, Mondays at 09:00")
		t.done++
		return
	}
	dir := filepath.Join(home, ".config", "systemd", "user")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		ui.Result(dir, "failed", err.Error())
		t.fail++
		return
	}
	files := map[string]string{
		"polkadot-sync.service": buildSyncService(root, goBin),
		"polkadot-sync.timer":   syncTimer,
	}
	for name, body := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
			ui.Result(name, "failed", err.Error())
			t.fail++
			return
		}
	}
	if !systemdRunning() {
		ui.Result("polkadot-sync.timer", "skipped", "written, but systemd is not running here")
		t.skip++
		return
	}
	if _, err := step.Sh("systemctl --user daemon-reload && systemctl --user enable --now polkadot-sync.timer"); err != nil {
		ui.Result("polkadot-sync.timer", "failed", "written, but systemctl refused it")
		t.fail++
		return
	}
	ui.Result("polkadot-sync.timer", "enabled", "Mondays at 09:00")
	t.done++
}
