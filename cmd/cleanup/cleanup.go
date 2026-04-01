package cleanup

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/yousefbustamiii/package-mate/internal/ui"
	"golang.org/x/term"
)

// Run executes the interactive cleanup command.
func Run() error {
	ui.Blank()

	// Show cleanup level menu
	choice := promptCleanupLevel()
	if choice == 0 {
		ui.Blank()
		fmt.Println("  " + ui.C(ui.Dim, "Cleanup cancelled."))
		ui.Blank()
		return nil
	}

	// Execute the cleanup
	return executeCleanup(choice)
}

// promptCleanupLevel shows the cleanup level selection menu.
func promptCleanupLevel() int {
	// Print the menu header
	fmt.Println(ui.C(ui.Grey, "  "+strings.Repeat("─", 60)))
	fmt.Println("  " + ui.C(ui.Bold+ui.White, "❯ WHICH CLEANUP LEVEL WOULD YOU LIKE?"))
	fmt.Println(ui.C(ui.Grey, "  "+strings.Repeat("─", 60)))
	ui.Blank()

	fmt.Println("  " + ui.C(ui.BrightCyan, "1.") + "  Standard Cleanup      " + ui.C(ui.Dim, "(Removes old versions, keeps current)"))
	ui.Blank()
	fmt.Println("  " + ui.C(ui.BrightCyan, "2.") + "  Deep Cleanup          " + ui.C(ui.Dim, "(Removes ALL cached files, including unused ones)"))
	ui.Blank()

	fmt.Println(ui.C(ui.Grey, "  "+strings.Repeat("─", 60)))
	fmt.Println("  " + ui.C(ui.Dim, "Press ") + ui.C(ui.White, "1 - 2") + ui.C(ui.Dim, " to choose your action"))
	fmt.Println("  " + ui.C(ui.Dim, "or press ") + ui.C(ui.White, "Enter") + ui.C(ui.Dim, " to cancel"))
	ui.Blank()

	// Check if we're in a terminal
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return 1 // Default to standard cleanup in non-interactive mode
	}

	// Read input
	fmt.Print("  " + ui.C(ui.Bold+ui.White, "❯ "))
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())
		switch input {
		case "1":
			return 1
		case "2":
			return 2
		default:
			return 0 // Cancel
		}
	}

	return 0
}

// executeCleanup runs the brew cleanup command with the specified option.
func executeCleanup(choice int) error {
	var cmdArgs []string
	var cleanupType string

	if choice == 1 {
		cmdArgs = []string{"cleanup"}
		cleanupType = "Standard Cleanup"
	} else {
		cmdArgs = []string{"cleanup", "--prune=all"}
		cleanupType = "Deep Cleanup"
	}

	// Show cleaning animation
	fmt.Println()
	fmt.Print("  " + ui.C(ui.BrightCyan, "❯ ") + ui.C(ui.Bold+ui.White, cleanupType) + " in progress")

	// Create the command
	cmd := exec.Command("brew", cmdArgs...)

	// Capture stdout and stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("\r\033[K")
		ui.Blank()
		ui.Fail("Failed to start cleanup: %v", err)
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Printf("\r\033[K")
		ui.Blank()
		ui.Fail("Failed to start cleanup: %v", err)
		return err
	}

	if err := cmd.Start(); err != nil {
		fmt.Printf("\r\033[K")
		ui.Blank()
		ui.Fail("Failed to start cleanup: %v", err)
		return err
	}

	// Animated dots
	done := make(chan struct{})
	go func() {
		dots := 0
		for {
			select {
			case <-done:
				return
			default:
				dots = (dots % 4) + 1
				fmt.Printf("\r\033[K  %s%s", ui.C(ui.BrightCyan, "❯ ")+ui.C(ui.Bold+ui.White, cleanupType)+" in progress", strings.Repeat(".", dots))
				time.Sleep(300 * time.Millisecond)
			}
		}
	}()

	// Read output (we'll parse it for the freed space)
	var output strings.Builder
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		output.WriteString(scanner.Text() + "\n")
	}

	// Also read stderr
	errScanner := bufio.NewScanner(stderr)
	for errScanner.Scan() {
		output.WriteString(errScanner.Text() + "\n")
	}

	// Wait for command to finish
	_ = cmd.Wait()
	close(done)

	// Clear the animation line
	fmt.Printf("\r\033[K")
	ui.Blank()

	// Parse the freed space from output
	freedSpace := parseFreedSpace(output.String())

	// Show result
	if freedSpace != "" {
		fmt.Println("  " + ui.C(ui.Grey, strings.Repeat("─", 60)))
		fmt.Println("  " + ui.C(ui.Bold+ui.BrightGreen, "✓ CLEANUP COMPLETE"))
		fmt.Println(ui.C(ui.Grey, "  "+strings.Repeat("─", 60)))
		ui.Blank()
		fmt.Printf("  %s %s freed\n", ui.C(ui.BrightGreen, "❯"), ui.C(ui.BrightGreen, freedSpace))
		ui.Blank()
	} else {
		// Nothing to clean
		fmt.Println("  " + ui.C(ui.Grey, strings.Repeat("─", 60)))
		fmt.Println("  " + ui.C(ui.Bold+ui.Yellow, "✓ NOTHING TO CLEAN"))
		fmt.Println(ui.C(ui.Grey, "  "+strings.Repeat("─", 60)))
		ui.Blank()
		fmt.Printf("  %s Your Homebrew cache is already clean!\n", ui.C(ui.Yellow, "❯"))
		ui.Blank()
	}

	return nil
}

// parseFreedSpace extracts the freed space from brew cleanup output.
// Matches patterns like: "2GB", "300MB", "1.5GB", "500MB", etc.
func parseFreedSpace(output string) string {
	// Brew outputs: "This operation has freed approximately 2GB of disk space."
	// or similar patterns with MB, GB, KB
	re := regexp.MustCompile(`freed approximately\s+([\d.]+\s*(?:KB|MB|GB|TB))`)
	matches := re.FindStringSubmatch(output)

	if len(matches) >= 2 {
		// Normalize the format (remove space between number and unit)
		space := strings.ReplaceAll(matches[1], " ", "")
		return space
	}

	// Alternative pattern: "Freed 2.5GB" or "cleaned 300MB"
	re2 := regexp.MustCompile(`(?:freed|cleaned)\s+([\d.]+\s*(?:KB|MB|GB|TB))`)
	matches2 := re2.FindStringSubmatch(output)

	if len(matches2) >= 2 {
		space := strings.ReplaceAll(matches2[1], " ", "")
		return space
	}

	return ""
}
