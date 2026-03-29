package sys

import (
	"os"
	"os/exec"
	"strings"

	"github.com/yousefbustamiii/package-mate/internal/ui"
)

// GetArch returns the system's architecture (e.g., "arm64", "x86_64") using uname -m.
func GetArch() string {
	out, _ := exec.Command("uname", "-m").Output()
	return strings.TrimSpace(string(out))
}

// HasXcodeCLT returns true if Xcode Command Line Tools are installed.
func HasXcodeCLT() bool {
	_, err := exec.LookPath("clang")
	return err == nil
}

// BrewPrefix returns the default Homebrew installation prefix based on system architecture.
func BrewPrefix() string {
	if GetArch() == "arm64" {
		return "/opt/homebrew"
	}
	return "/usr/local"
}

// GetShell returns the user's preferred shell from the SHELL environment variable,
// falling back to /bin/bash if it's not set or doesn't exist on the filesystem.
func GetShell() string {
	shell := os.Getenv("SHELL")
	if shell != "" {
		if _, err := os.Stat(shell); err == nil {
			return shell
		}
		ui.Warn("Preferred shell %s not found", shell)
	} else {
		ui.Hint("No $SHELL set, defaulting to /bin/bash")
	}

	// Try /bin/bash as primary fallback
	if _, err := os.Stat("/bin/bash"); err == nil {
		return "/bin/bash"
	}

	// Final fallback to /bin/sh
	ui.Warn("/bin/bash not found, falling back to /bin/sh")
	return "/bin/sh"
}

// ShellCommand returns an exec.Cmd that executes the given script string
// using the system's preferred shell with the -c flag.
func ShellCommand(script string) *exec.Cmd {
	shell := GetShell()
	// All modern Unix shells support -c
	return exec.Command(shell, "-c", script)
}
