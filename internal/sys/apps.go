package sys

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// IsRunning returns true if an application with the given name is currently active.
// It uses osascript to check for the process name.
func IsRunning(appName string) bool {
	// ❯ For Docker, we check the specific process name
	if strings.EqualFold(appName, "Docker") {
		cmd := exec.Command("pgrep", "-f", "Docker.app")
		return cmd.Run() == nil
	}

	// ❯ General check using osascript (this is better for GUI apps)
	// We escape double quotes to prevent AppleScript injection
	safeName := strings.ReplaceAll(appName, "\"", "\\\"")
	script := fmt.Sprintf(`tell application "System Events" to (name of processes) contains "%s"`, safeName)
	out, err := exec.Command("osascript", "-e", script).Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(out)) == "true"
}

// AppExists checks for the existence of an .app bundle in /Applications or ~/Applications.
// It returns the full path and a boolean indicating if it was found.
func AppExists(appName string) (string, bool) {
	paths := []string{
		filepath.Join("/Applications", appName+".app"),
		filepath.Join(os.Getenv("HOME"), "Applications", appName+".app"),
	}

	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p, true
		}
	}
	return "", false
}

// AppInTrash checks if the app is currently in the system Trash.
func AppInTrash(appName string) (string, bool) {
	trashPath := filepath.Join(os.Getenv("HOME"), ".Trash", appName+".app")
	if _, err := os.Stat(trashPath); err == nil {
		return trashPath, true
	}
	return "", false
}
