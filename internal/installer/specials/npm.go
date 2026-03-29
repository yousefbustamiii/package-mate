package specials

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/yousefbustamiii/package-mate/internal/ui"
)

// ── npm binary ─────────────────────────────────────────────────────────────────

func npmCheck() (bool, string) {
	if _, err := exec.LookPath("npm"); err == nil {
		out, _ := exec.Command("npm", "--version").Output()
		return true, strings.TrimSpace(string(out))
	}
	return false, ""
}

func npmInstall() error {
	return fmt.Errorf("npm not found — install Node.js first")
}

// ── npm global packages ────────────────────────────────────────────────────────

func npmGCheck(pkg string) (bool, string) {
	out, _ := exec.Command("npm", "list", "-g", "--depth=0", pkg).Output()
	if strings.Contains(string(out), pkg) {
		return true, ""
	}
	return false, ""
}

func npmGInstall(name, pkg string) error {
	ui.Doing("Installing package %s", pkg)
	cmd := exec.Command("npm", "install", "-g", pkg)
	var stderr strings.Builder
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s", strings.TrimSpace(stderr.String()))
	}
	ui.Done("Successfully installed %s", name)
	return nil
}
