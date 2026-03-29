package versions

import (
	"os/exec"
	"strings"

	"github.com/yousefbustamiii/package-mate/internal/sys"
)

// nvmVersion returns the installed NVM version string.
func nvmVersion() []string {
	out, err := sys.ShellCommand(
		`source "$HOME/.nvm/nvm.sh" 2>/dev/null && nvm --version`).Output()
	if err == nil {
		if v := strings.TrimSpace(string(out)); v != "" {
			return []string{v}
		}
	}
	return nil
}

// rustVersion returns the installed Rust version string.
func rustVersion() []string {
	out, err := exec.Command("rustc", "--version").Output()
	if err == nil {
		v := strings.TrimSpace(string(out))
		if idx := strings.Index(v, " ("); idx > 0 {
			v = v[:idx] // "rustc 1.75.0 (abc 2024)" → "rustc 1.75.0"
		}
		if v != "" {
			return []string{v}
		}
	}
	return nil
}

// pytestVersion returns the installed pytest version string.
func pytestVersion() []string {
	out, err := exec.Command("pip3", "show", "pytest").Output()
	if err == nil {
		for _, line := range strings.Split(string(out), "\n") {
			if strings.HasPrefix(line, "Version:") {
				if v := strings.TrimSpace(strings.TrimPrefix(line, "Version:")); v != "" {
					return []string{"pytest " + v}
				}
			}
		}
	}
	return nil
}

// pipxVersion returns the installed version of a pipx-managed package.
func pipxVersion(pkg string) []string {
	out, err := exec.Command("pipx", "list", "--short").Output()
	if err == nil {
		for _, line := range strings.Split(string(out), "\n") {
			if strings.Contains(line, pkg) {
				if parts := strings.Fields(line); len(parts) >= 2 {
					return []string{pkg + " " + parts[1]}
				}
			}
		}
	}
	return nil
}
