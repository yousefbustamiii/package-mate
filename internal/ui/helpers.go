package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

// ── Print helpers ──────────────────────────────────────────────────────────────

func Fail(format string, args ...any) {
	stopDoing()
	_, _ = fmt.Fprintf(os.Stderr, "  "+C(Red, "✗")+"  "+format+"\n\n", args...)
}

func Warn(format string, args ...any) {
	stopDoing()
	fmt.Printf("  "+C(Yellow, "!")+"  "+format+"\n", args...)
}

func Row(label, value string) {
	fmt.Printf("  "+C(Dim, "%-16s")+"%s\n", label, value)
	Blank()
}

// Stop explicitly halts any running spinner animation.
func Stop() {
	stopDoing()
}

func Doing(format string, args ...any) {
	startDoing(fmt.Sprintf(format, args...))
}

func Done(format string, args ...any) {
	stopDoing()
	fmt.Printf("  " + C(BrightGreen, "✓") + "  " + C(Bold+White, fmt.Sprintf(format, args...)) + "\n")
}

func Hint(format string, args ...any) {
	fmt.Println(C(Dim, "  "+fmt.Sprintf(format, args...)))
}

func Blank() {
	fmt.Println()
}

func Rule() {
	fmt.Println("  " + C(Grey, strings.Repeat("─", 60)))
}

func Header(title string) {
	stopDoing()
	Blank()
	Rule()
	fmt.Println("  " + C(Bold+White, "❯ "+strings.ToUpper(title)))
	Rule()
	Blank()
}

func Footer() {
	Rule()
	Blank()
}

// ── Safety ─────────────────────────────────────────────────────────────────────

// ShowSafetyAlert displays a red warning when a user tries to modify a system path.
func ShowSafetyAlert(path string) {
	stopDoing()
	Blank()
	fmt.Println("  " + C(Bold+Red, "✗ SAFETY ALERT:"))
	Blank()
	fmt.Println("  " + C(Bold+White, "❯ Protected System Path"))
	Blank()
	fmt.Printf("  Package Mate is strictly forbidden from modifying "+C(Bold+White, "%s")+".\n", path)
	fmt.Println("  This file is part of the macOS core and must not be deleted.")
	Blank()
}

func AlreadyInstalled(name, version string) {
	stopDoing()
	v := ""
	if version != "" {
		v = " (" + version + ")"
	}
	fmt.Printf("  "+C(Cyan, "~")+"  "+C(Bold+White, "Already installed: ")+"%s%s\n", name, C(Dim, v))
	Blank()
}

// PromptConfirmation asks the user to type the action (e.g. INSTALL) to proceed.
func PromptConfirmation(action, name string) bool {
	Blank()
	fmt.Printf("  "+C(Bold+White, "Type ")+C(Bold+BrightCyan, "%s")+" to %s %s\n", action, strings.ToLower(action), name)
	Blank()
	fmt.Print("  " + C(Bold+White, "❯ "))

	input := readInput()

	if strings.EqualFold(input, action) {
		return true
	}

	Blank()
	Fail("Aborted. Input did not match %s", action)
	return false
}

// PromptInstall asks the user to confirm a fresh install, with an option to run in background.
// Returns "INSTALL", "INSTALL/BG", or "" if declined.
func PromptInstall(name string) string {
	stopDoing()
	Blank()
	fmt.Println("  " + C(Grey, strings.Repeat("─", 60)))
	fmt.Println("  " + C(Bold+White, "❯ INSTALL "+strings.ToUpper(name)))
	fmt.Println("  " + C(Grey, strings.Repeat("─", 60)))
	Blank()
	fmt.Printf("  "+C(Bold+White, "Type ")+C(Bold+BrightCyan, "INSTALL")+C(Bold+White, " to install %s now.\n"), name)
	Blank()
	fmt.Printf("  " + C(Dim, "Or type ") + C(BrightCyan, "INSTALL/BG") + C(Dim, " to download it in the background.") + "\n")
	Blank()
	fmt.Print("  " + C(Bold+White, "❯ "))

	input := readInput()

	switch {
	case strings.EqualFold(input, "INSTALL"):
		return "INSTALL"
	case strings.EqualFold(input, "INSTALL/BG"):
		return "INSTALL/BG"
	}

	Blank()
	Fail("Operation declined. No changes made.")
	return ""
}

// PromptOverride asks the user to confirm...
// Returns "INSTALL", "INSTALL/BG", "OVERRIDE", "OVERRIDE/BG", or "" if declined.
func PromptOverride(name, target, path string, isProtected, isSystem bool) string {
	stopDoing()
	Blank()
	fmt.Println("  " + C(Yellow, strings.Repeat("─", 60)))
	fmt.Println("  " + C(Bold+Yellow, "❯ UNMANAGED BINARY DETECTED"))
	fmt.Println("  " + C(Yellow, strings.Repeat("─", 60)))
	Blank()

	fmt.Printf("  "+C(White, "A version of %s is already present at %s, but it's not managed by Package Mate.\n"), name, path)
	Blank()

	fmt.Printf("  "+C(Cyan, "Current")+"  "+C(Dim, ":")+"  %s\n", path)
	fmt.Printf("  "+C(Cyan, "Target ")+"  "+C(Dim, ":")+"  %s\n", target)

	Blank()

	if isProtected {
		fmt.Printf("  " + C(White, "Would you like to install the Package Mate version alongside it? ") + C(Dim, "(Protected path, cannot override)") + "\n")
	} else {
		fmt.Printf("  " + C(White, "Would you like to install the Package Mate version or fully OVERRIDE the existing one?") + "\n")
	}

	Blank()

	// ── Single-line Command Hints ──────────────────────────────────────────────
	fmt.Printf("  " + C(Bold+White, "Type "))

	// Install options
	fmt.Printf(C(Bold+BrightCyan, "INSTALL") + " " + C(Dim, "or") + " " + C(BrightCyan, "INSTALL/BG") + " " + C(Dim, "to install"))

	// Override options (if not protected)
	if !isProtected {
		fmt.Printf(" " + C(Dim, "or") + " " + C(Bold+BrightCyan, "OVERRIDE"))

		// Background override (if not a system path)
		if !isSystem {
			fmt.Printf(" " + C(Dim, "or") + " " + C(BrightCyan, "OVERRIDE/BG"))
		}

		fmt.Printf(" " + C(Dim, "to replace the existing binary"))
	}
	fmt.Printf("\n")

	Blank()
	fmt.Print("  " + C(Bold+White, "❯ "))

	input := readInput()

	switch {
	case strings.EqualFold(input, "INSTALL"):
		return "INSTALL"
	case strings.EqualFold(input, "INSTALL/BG"):
		return "INSTALL/BG"
	case !isProtected && strings.EqualFold(input, "OVERRIDE"):
		return "OVERRIDE"
	case !isProtected && !isSystem && strings.EqualFold(input, "OVERRIDE/BG"):
		return "OVERRIDE/BG"
	}

	Blank()
	Fail("Operation declined. No changes made.")
	return ""
}

// PromptDifferentBrew asks for UPDATE/OVERRIDE choice when another Brew version exists.
// Returns "UPDATE", "UPDATE/BG", "OVERRIDE", "OVERRIDE/BG", or "" if declined.
func PromptDifferentBrew(target, oldFormula string) string {
	stopDoing()
	Blank()
	fmt.Println("  " + C(Yellow, strings.Repeat("─", 60)))
	fmt.Println("  " + C(Bold+Yellow, "❯ DIFFERENT BREW VERSION DETECTED"))
	fmt.Println("  " + C(Yellow, strings.Repeat("─", 60)))
	Blank()

	fmt.Printf("  " + C(Bold+White, "A different Homebrew version (") + C(BrightCyan, oldFormula) + C(Bold+White, ") is managed on your system.\n"))
	Blank()

	fmt.Printf("  "+C(Cyan, "Current")+"  "+C(Dim, ":")+"  %s\n", oldFormula)
	fmt.Printf("  "+C(Cyan, "Target ")+"  "+C(Dim, ":")+"  %s\n", target)

	Blank()
	fmt.Printf("  " + C(Bold+White, "Would you like to UPDATE to the Package Mate version or fully OVERRIDE the existing one?") + "\n")
	Blank()
	fmt.Printf("  "+C(Bold+White, "Type ")+
		C(Bold+BrightCyan, "UPDATE")+" "+C(Dim, "or")+
		" "+C(BrightCyan, "UPDATE/BG")+" "+C(Dim, "to upgrade")+"\n")
	fmt.Printf("  "+C(Bold+White, "     ")+
		C(Bold+BrightCyan, "OVERRIDE")+" "+C(Dim, "or")+
		" "+C(BrightCyan, "OVERRIDE/BG")+" "+C(Dim, "to replace the existing formula")+"\n")
	Blank()
	fmt.Print("  " + C(Bold+White, "❯ "))

	input := readInput()

	switch {
	case strings.EqualFold(input, "UPDATE"):
		return "UPDATE"
	case strings.EqualFold(input, "UPDATE/BG"):
		return "UPDATE/BG"
	case strings.EqualFold(input, "OVERRIDE"):
		return "OVERRIDE"
	case strings.EqualFold(input, "OVERRIDE/BG"):
		return "OVERRIDE/BG"
	}

	Blank()
	Fail("Operation declined. No changes made.")
	return ""
}

// PromptDependencyRemoval asks the user to confirm a forced removal of a dependency.
func PromptDependencyRemoval(name, errMsg string) bool {
	stopDoing()
	Blank()
	fmt.Println("  " + C(Yellow, strings.Repeat("─", 60)))
	fmt.Println("  " + C(Bold+Yellow, "❯ DEPENDENCY CONFLICT DETECTED"))
	fmt.Println("  " + C(Yellow, strings.Repeat("─", 60)))
	Blank()

	fmt.Printf("  "+C(Bold+White, "Homebrew refuses to uninstall %s because other tools depend on it.\n"), name)
	if errMsg != "" {
		fmt.Printf("  " + C(Dim, errMsg) + "\n")
	}
	Blank()
	fmt.Printf("  " + C(Bold+White, "Type ") + C(Bold+BrightCyan, "DEPENDENCY") + C(Bold+White, " to force removal.") + "\n")
	fmt.Printf("  " + C(Dim, "  (This may break tools that rely on this library)") + "\n")
	Blank()
	fmt.Print("  " + C(Bold+White, "❯ "))

	input := readInput()

	if strings.EqualFold(input, "DEPENDENCY") {
		return true
	}

	Blank()
	Fail("Force removal declined. Keeping the current version.")
	return false
}

func ShowManualAppUninstall(appName string) {
	Header("Manual Uninstall Required")
	Warn("%s was installed manually (outside of Package Mate).", appName)
	Blank()

	fmt.Println("  Package Mate cannot safely remove unmanaged GUI apps.")
	fmt.Println("  Please follow these steps to remove it cleanly:")
	Blank()

	fmt.Println(C(Cyan, "  1.") + " Open " + C(Bold+White, "System Settings"))
	fmt.Println(C(Cyan, "  2.") + " Go to " + C(Bold+White, "General") + " → " + C(Bold+White, "Storage"))
	fmt.Println(C(Cyan, "  3.") + " Click the " + C(Bold+White, "ⓘ") + " icon next to " + C(Bold+White, "Applications"))
	fmt.Println(C(Cyan, "  4.") + " Find " + C(Bold+BrightPink, appName) + " and click " + C(Red, "Delete..."))
	Blank()

	Hint("Once deleted, run " + C(Cyan, "mate "+strings.ToLower(appName)) + " to manage it here.")
}

func ShowAppRunning(appName string) {
	Header("Close " + appName)
	Warn("%s is currently running.", appName)
	Blank()
	fmt.Printf("  Please quit "+C(Bold+White, "%s")+" to allow Package Mate to perform this update.\n", appName)
	Blank()
	fmt.Println("  " + C(Dim, "Check your menu bar or use Command+Q to close it."))
	Blank()
}

// PromptOutdated asks the user to confirm an update.
// Returns "UPDATE", "UPDATE/BG", or "" if declined.
func PromptOutdated(name, current, latest string) string {
	Blank()
	Warn("Update Available")
	fmt.Printf("  "+C(Bold+White, "A newer version of %s is available.\n"), name)
	Blank()
	fmt.Printf("  "+C(Dim, "Current:  ")+"%s\n", current)
	fmt.Printf("  "+C(Dim, "Latest:   ")+"%s\n", latest)
	Blank()
	fmt.Printf("  "+C(Bold+White, "Type ")+C(Bold+BrightCyan, "UPDATE")+C(Bold+White, " to update now.\n"))
	fmt.Printf("  "+C(Dim, "Or type ")+C(BrightCyan, "UPDATE/BG")+C(Dim, " to update in the background.")+"\n")
	Blank()
	fmt.Print("  " + C(Bold+White, "❯ "))

	input := readInput()

	switch {
	case strings.EqualFold(input, "UPDATE"):
		return "UPDATE"
	case strings.EqualFold(input, "UPDATE/BG"):
		return "UPDATE/BG"
	}

	Blank()
	Fail("Update declined. Skipping for now.")
	return ""
}

// PromptActionMenu shows a 1-3 selection menu for tool management.
func PromptActionMenu(name string) int {
	Blank()
	Header("Manage " + name)

	fmt.Printf("  " + C(Cyan, "1.") + "  " + C(Bold+White, "Install or Update") + "\n")
	Blank()
	fmt.Printf("  " + C(Red, "2.") + "  " + C(Bold+White, "Uninstall") + "\n")
	Blank()
	fmt.Printf("  " + C(BrightPink, "3.") + "  " + C(Bold+White, "Information & Versions") + "\n")
	Blank()

	fmt.Printf("  " + C(Bold+White, "Press ") + C(Bold+Cyan, "1 - 3") + C(Bold+White, " to choose your action") + "\n")
	fmt.Printf("  " + C(Dim, "or press ") + C(White, "Enter") + C(Dim, " to cancel") + "\n")
	Blank()
	fmt.Print("  " + C(Bold+White, "❯ "))

	input := readInput()

	switch input {
	case "1":
		return 1
	case "2":
		return 2
	case "3":
		return 3
	default:
		return 0
	}
}

// readInput reads a full line of text from Stdin and trims whitespace.
func readInput() string {
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		return strings.TrimSpace(scanner.Text())
	}
	return ""
}

// SelectionOption represents an item in an interactive menu.
type SelectionOption struct {
	Label string
	Value string
}

// PromptSelection displays an interactive menu navigable with arrows or j/k.
func PromptSelection(title string, opts []SelectionOption) int {
	if len(opts) == 0 {
		return -1
	}

	selected := 0

	// Set terminal to raw mode
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return -1
	}
	defer func() {
		_ = term.Restore(int(os.Stdin.Fd()), oldState)
	}()

	// Hide cursor
	fmt.Print("\033[?25l")
	defer fmt.Print("\033[?25h")

	for {
		// Clear lines (we'll redraw) - actually let's use a simpler approach for now
		// by just printing and using \r or moving cursor

		// Render title
		fmt.Printf("\r\033[K  %s\r\n", C(Bold+White, title))

		for i, opt := range opts {
			prefix := "    "
			label := C(Dim, opt.Label)
			if i == selected {
				prefix = C(Cyan, "  ❯ ")
				label = C(Bold+BrightCyan, opt.Label)
			}
			fmt.Printf("\033[K%s%s\r\n", prefix, label)
		}

		// Read input
		var b [3]byte
		n, _ := os.Stdin.Read(b[:])
		if n == 1 {
			switch b[0] {
			case 'j', 'J', 14: // 14 is Ctrl+N (not needed but common)
				selected = (selected + 1) % len(opts)
			case 'k', 'K', 16: // 16 is Ctrl+P
				selected = (selected - 1 + len(opts)) % len(opts)
			case 13: // Enter
				// Move cursor down to end of list before returning
				return selected
			case 3: // Ctrl+C
				return -1
			}
		} else if n == 3 && b[0] == 27 && b[1] == '[' {
			switch b[2] {
			case 'A': // Arrow Up
				selected = (selected - 1 + len(opts)) % len(opts)
			case 'B': // Arrow Down
				selected = (selected + 1) % len(opts)
			}
		}

		// Move cursor back up to redraw
		fmt.Printf("\033[%dA", len(opts)+1)
	}
}

// ── Terminal helpers ───────────────────────────────────────────────────────────

func termWidth() int {
	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || w <= 0 {
		return 120
	}
	return w
}

// ── Background job helpers ─────────────────────────────────────────────────────

// ShowBgQueued prints a confirmation that a tool has been queued for background download.
func ShowBgQueued(name string) {
	stopDoing()
	Blank()
	fmt.Println("  " + C(Grey, strings.Repeat("─", 60)))
	fmt.Println("  " + C(Bold+BrightCyan, "❯ QUEUED FOR BACKGROUND DOWNLOAD"))
	fmt.Println("  " + C(Grey, strings.Repeat("─", 60)))
	Blank()
	fmt.Printf("  "+C(Bold+White, "%s")+" is now downloading in the background.\n", name)
	Blank()
	fmt.Printf("  "+C(Dim, "Run ")+C(Cyan, "mate bg")+C(Dim, " to check progress or abort the download.")+"\n")
	Blank()
}
