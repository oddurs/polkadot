package main

import (
	"github.com/oddurs/polkadot/internal/link"
	"github.com/oddurs/polkadot/internal/ui"
)

// themeLinks is the compact, runtime London Moquette bundle. The full Subway
// Seat generator remains its own project; Polkadot ships only the ports it
// configures so a fresh machine never needs a second theme checkout.
var themeLinks = []linkMap{
	{"theme/london-moquette", ".config/polkadot-theme"},
	{"theme/london-moquette/ghostty", ".config/ghostty/themes"},
	{"theme/london-moquette/codex/london-moquette.tmTheme", ".codex/themes/london-moquette.tmTheme"},
	{"theme/london-moquette/opencode/london-moquette.json", ".config/opencode/themes/london-moquette.json"},
}

func doTheme(l *link.Linker, t *tally) {
	ui.Section("theme: London Moquette")
	linkMappings(l, t, themeLinks)
	ui.Blank()
}
