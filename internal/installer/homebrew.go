package installer

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/yousefbustamiii/package-mate/internal/sys"
	"github.com/yousefbustamiii/package-mate/internal/ui"
)

// IsBrewInstalled returns true when `brew` is on PATH or in standard prefix.
func IsBrewInstalled() bool {
	if _, err := exec.LookPath("brew"); err == nil {
		return true
	}
	// Check standard location
	_, err := os.Stat(sys.BrewPrefix() + "/bin/brew")
	return err == nil
}

// EnsureBrew installs Homebrew via the official bash script if it is missing.
func EnsureBrew() Result {
	if IsBrewInstalled() {
		return Result{ItemName: "Homebrew", Status: StatusAlreadyHave, Version: brewVersion()}
	}

	// Check for Xcode CLT
	if !sys.HasXcodeCLT() {
		ui.Warn("Xcode Command Line Tools are required for Homebrew but appear to be missing.")
		if !ui.PromptConfirmation("INSTALL", "Xcode Command Line Tools") {
			return Result{ItemName: "Homebrew", Status: StatusFailed, Err: fmt.Errorf("Xcode CLT required")}
		}
		ui.Doing("Installing %s", "Xcode CLT")
		cmd := exec.Command("xcode-select", "--install")
		if err := cmd.Run(); err != nil {
			ui.Fail("Failed to trigger Xcode CLT installation: %v", err)
			return Result{ItemName: "Homebrew", Status: StatusFailed, Err: err}
		}
		ui.Done("Triggered Xcode CLT install. Please complete the installer dialog and run this again.")
		return Result{ItemName: "Xcode CLT", Status: StatusInstalled}
	}

	cmd := sys.ShellCommand(`/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"`)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		return Result{ItemName: "Homebrew", Status: StatusFailed, Err: err}
	}
	return Result{ItemName: "Homebrew", Status: StatusInstalled}
}

// BrewExe returns the path to the brew binary, checking PATH first then standard prefixes.
func BrewExe() string {
	if _, err := exec.LookPath("brew"); err == nil {
		return "brew"
	}
	path := sys.BrewPrefix() + "/bin/brew"
	if _, err := os.Stat(path); err == nil {
		return path
	}
	return "brew" // Fallback to "brew" and let exec.Command fail naturally if not found
}

// UpdateBrew runs `brew update`.
func UpdateBrew() Result {
	cmd := exec.Command(BrewExe(), "update")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return Result{ItemName: "Homebrew Update", Status: StatusFailed, Err: err}
	}
	return Result{ItemName: "Homebrew Update", Status: StatusInstalled}
}

func brewVersion() string {
	out, _ := exec.Command(BrewExe(), "--version").Output()
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) > 0 {
		return strings.TrimPrefix(lines[0], "Homebrew ")
	}
	return ""
}
