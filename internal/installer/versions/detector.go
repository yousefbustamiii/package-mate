package versions

import (
	"os/exec"
	"strings"
	"time"

	"github.com/yousefbustamiii/package-mate/internal/components"
	"github.com/yousefbustamiii/package-mate/internal/installer"
)

// InstalledVersions returns all version strings for the given item and its
// installation date if available. Returns nil if nothing is detected.
func InstalledVersions(item components.InstallItem) ([]string, time.Time) {
	switch item.Special {
	case "brew-install":
		if versions, date, ok := installer.GetBrewInfo(item); ok {
			return versions, date
		}
		return nil, time.Time{}
	case "brew-update":
		return nil, time.Time{}
	case "nvm":
		return nvmVersion(), time.Time{}
	case "rustup":
		return rustVersion(), time.Time{}
	case "npm-check":
		return npmVersion(), time.Time{}
	case "claude":
		return npmGlobalVersion("@anthropic-ai/claude-code"), time.Time{}
	case "gemini":
		return npmGlobalVersion("@google/gemini-cli"), time.Time{}
	case "jest":
		return npmGlobalVersion("jest"), time.Time{}
	case "npm-g":
		return npmGlobalVersion(item.Formula), time.Time{}
	case "pytest":
		return pytestVersion(), time.Time{}
	case "pipx-g":
		return pipxVersion(item.Formula), time.Time{}
	}

	if item.Formula != "" || item.Cask != "" {
		if versions, date, ok := installer.GetBrewInfo(item); ok {
			return versions, date
		}
	}
	// Binary fallback
	if item.Binary != "" {
		if _, err := exec.LookPath(item.Binary); err == nil {
			out, err := exec.Command(item.Binary, "--version").Output()
			if err == nil {
				v := strings.SplitN(strings.TrimSpace(string(out)), "\n", 2)[0]
				if v != "" {
					return []string{v}, time.Time{}
				}
			}
		}
	}
	return nil, time.Time{}
}
