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
	// ❯ Search in standard locations
	roots := []string{"/Applications", filepath.Join(os.Getenv("HOME"), "Applications")}

	for _, root := range roots {
		p := filepath.Join(root, appName+".app")
		if _, err := os.Stat(p); err == nil {
			return p, true
		}
	}
	return "", false
}

// FindBundle is the "Smart Discovery" engine. It tries multiple strategies to find
// a GUI app on the system based on its Catalog metadata.
func FindBundle(name, cask, binary string) (string, bool) {
	// 1. Known stubborn Mappings (Catalog Name/Cask -> Real App Name)
	hardMappings := map[string]string{
		"iterm2": "iTerm",
		"zoom":   "zoom.us",
		"visual-studio-code": "Visual Studio Code",
		"postman": "Postman",
		"postman-agent": "Postman Agent",
		"docker": "Docker",
		"tableplus": "TablePlus",
		"postgresql": "Postgres",
	}

	// 2. Check Hard Mappings (Cask based)
	if cask != "" {
		if mapped, ok := hardMappings[strings.ToLower(cask)]; ok {
			if p, ok := AppExists(mapped); ok {
				return p, true
			}
		}
	}

	// 3. Try Name (Display Name)
	if p, ok := AppExists(name); ok {
		return p, true
	}

	// 4. Try Cask name itself (normalized)
	if cask != "" {
		// e.g. "visual-studio-code" -> "Visual Studio Code"
		normalizedCask := strings.Title(strings.ReplaceAll(cask, "-", " "))
		if p, ok := AppExists(normalizedCask); ok {
			return p, true
		}
		// e.g. "iterm2"
		if p, ok := AppExists(cask); ok {
			return p, true
		}
	}

	// 5. Try Binary Hint
	if binary != "" {
		if p, ok := AppExists(binary); ok {
			return p, true
		}
		// Capitalized binary
		if p, ok := AppExists(strings.Title(binary)); ok {
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
