package versions

import (
	"os/exec"
	"strings"
)

// npmVersion returns the installed npm version string.
func npmVersion() []string {
	out, err := exec.Command("npm", "--version").Output()
	if err == nil {
		if v := strings.TrimSpace(string(out)); v != "" {
			return []string{"npm " + v}
		}
	}
	return nil
}

// npmGlobalVersion returns ["pkg 1.2.3"] for a globally installed npm package.
func npmGlobalVersion(pkg string) []string {
	out, err := exec.Command("npm", "list", "-g", "--depth=0", pkg).Output()
	if err != nil {
		return nil
	}
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if strings.Contains(line, pkg) {
			// Format: "── @scope/name@version" or "── name@version"
			if idx := strings.LastIndex(line, "@"); idx > 0 {
				if v := line[idx+1:]; v != "" {
					return []string{pkg + " " + v}
				}
			}
		}
	}
	return nil
}
