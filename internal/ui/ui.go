// Package ui is the whole presentation layer: a few colours, aligned status
// lines, and nothing that needs a dependency.
package ui

import (
	"fmt"
	"os"
	"strings"
)

var color = os.Getenv("NO_COLOR") == "" && os.Getenv("TERM") != "dumb"

func paint(code, s string) string {
	if !color {
		return s
	}
	return "\x1b[" + code + "m" + s + "\x1b[0m"
}

func Dim(s string) string    { return paint("2", s) }
func Bold(s string) string   { return paint("1", s) }
func Green(s string) string  { return paint("32", s) }
func Yellow(s string) string { return paint("33", s) }
func Red(s string) string    { return paint("31", s) }
func Blue(s string) string   { return paint("34", s) }

// Banner prints the run header.
func Banner(mode string) {
	fmt.Println()
	fmt.Println("  " + Bold("polkadot") + Dim("  ·  a mac, the way I like it"))
	if mode != "" {
		fmt.Println("  " + Dim(mode))
	}
	fmt.Println()
}

// Section starts a named group of steps.
func Section(name string) {
	fmt.Println("  " + Blue("──") + " " + Bold(name))
}

const pad = 46

// Result reports one step. Verb is past tense: linked, installed, skipped.
func Result(name, verb, detail string) {
	mark, v := Green("✓"), Green(verb)
	switch verb {
	case "skipped", "already":
		mark, v = Dim("·"), Dim(verb)
	case "would":
		mark, v = Yellow("~"), Yellow(verb)
	case "failed":
		mark, v = Red("✗"), Red(verb)
	}
	line := "  " + mark + " " + name
	if n := pad - len([]rune(name)); n > 0 {
		line += strings.Repeat(" ", n)
	} else {
		line += "  "
	}
	line += v
	if detail != "" {
		line += " " + Dim(detail)
	}
	fmt.Println(line)
}

func Note(s string) { fmt.Println("    " + Dim(s)) }
func Blank()        { fmt.Println() }

// Summary closes the run.
func Summary(done, skipped, failed int) {
	fmt.Println()
	parts := []string{fmt.Sprintf("%d changed", done), Dim(fmt.Sprintf("%d already set", skipped))}
	if failed > 0 {
		parts = append(parts, Red(fmt.Sprintf("%d failed", failed)))
	}
	fmt.Println("  " + strings.Join(parts, Dim("  ·  ")))
	fmt.Println()
}
