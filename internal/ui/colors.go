package ui

import (
	"fmt"
	"os"
)

// ── ANSI codes ─────────────────────────────────────────────────────────────────

const (
	Bold        = "\033[1m"
	Dim         = "\033[2m"
	Reset       = "\033[0m"
	Red         = "\033[31m"
	BrightGreen = "\033[92m"
	Yellow      = "\033[33m"
	Cyan        = "\033[36m"
	BrightCyan  = "\033[96m"
	BrightPink  = "\033[95m"
	White       = "\033[97m"
	Grey        = "\033[37m" // medium grey — used for the banner
)

// ── TTY detection ──────────────────────────────────────────────────────────────

var isTTY bool

func init() {
	if fi, _ := os.Stdout.Stat(); fi != nil {
		isTTY = (fi.Mode() & os.ModeCharDevice) != 0
	}
}

// C wraps text in an ANSI code when stdout is a TTY.
func C(code, s string) string {
	if !isTTY {
		return s
	}
	return code + s + Reset
}

// RGB returns a 24-bit truecolor ANSI escape for the given RGB values.
func RGB(r, g, b int) string {
	if !isTTY {
		return ""
	}
	return fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b)
}
