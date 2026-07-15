// Package style provides a botan (peony) colored terminal theme for CLI output.
package style

import (
	"os"
	"strings"
)

// Botan (牡丹) is the traditional Japanese peony color: a deep pink-magenta.
// The palette layers a bright bloom tone for headings, the core botan tone
// for labels/accents, and a deep wine tone for muted/secondary text.
const (
	reset = "\x1b[0m"
	bold  = "\x1b[1m"
	dim   = "\x1b[2m"

	bloom = "\x1b[38;2;227;107;174m" // bright peony bloom
	botan = "\x1b[38;2;170;76;143m"  // core botan pink-magenta
	wine  = "\x1b[38;2;107;37;69m"   // deep muted wine
)

var enabled = shouldColor()

func shouldColor() bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	if os.Getenv("TERM") == "dumb" {
		return false
	}
	return isColorTerminal()
}

func wrap(code, s string) string {
	if !enabled {
		return s
	}
	return code + s + reset
}

// Title styles a top-level heading, e.g. an episode title.
func Title(s string) string { return wrap(bold+bloom, s) }

// Heading styles a section heading, e.g. an alcohol entry name.
func Heading(s string) string { return wrap(bold+botan, s) }

// Accent styles inline emphasis, e.g. a Japanese name alongside its English one.
func Accent(s string) string { return wrap(botan, s) }

// Label styles a field label, e.g. "Type:", "Brand:".
func Label(s string) string { return wrap(wine, s) }

// Muted styles secondary/supporting text, e.g. air dates.
func Muted(s string) string { return wrap(dim, s) }

// Bullet returns a botan-colored peony bullet glyph.
func Bullet() string { return wrap(botan, "❀") }

// Divider returns a botan-colored horizontal rule of the given width.
func Divider(width int) string {
	rule := ""
	for i := 0; i < width; i++ {
		rule += "─"
	}
	return wrap(wine, rule)
}

// Banner returns a decorative botan-colored box carrying the given title,
// for use as CLI branding on the usage/help screen.
func Banner(title string) string {
	const pad = 3
	inner := pad*2 + len([]rune(title))
	rule := strings.Repeat("═", inner)
	mid := "║" + strings.Repeat(" ", pad) + title + strings.Repeat(" ", pad) + "║"

	return strings.Join([]string{
		wrap(bloom, "❀"+rule+"❀"),
		wrap(bold+botan, mid),
		wrap(bloom, "❀"+rule+"❀"),
	}, "\n")
}
