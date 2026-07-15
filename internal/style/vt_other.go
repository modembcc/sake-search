//go:build !windows

package style

import "os"

// isColorTerminal reports whether stdout is an interactive terminal capable
// of ANSI escape sequences. Non-Windows terminals support ANSI natively, so
// a plain char-device check is sufficient.
func isColorTerminal() bool {
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}
