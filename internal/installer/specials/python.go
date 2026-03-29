package specials

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/yousefbustamiii/package-mate/internal/ui"
)

// ── pytest ─────────────────────────────────────────────────────────────────────

func pytestCheck() (bool, string) {
	out, _ := exec.Command("pip3", "show", "pytest").Output()
	if strings.Contains(string(out), "Name: pytest") {
		return true, ""
	}
	return false, ""
}

func pytestInstall() error {
	ui.Doing("Installing %s", "Pytest")
	cmd := exec.Command("pip3", "install", "--user", "pytest")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	ui.Done("Successfully installed %s", "Pytest")
	return nil
}

// ── pipx global packages ───────────────────────────────────────────────────────

func pipxCheck(pkg string) (bool, string) {
	out, _ := exec.Command("pipx", "list").Output()
	if strings.Contains(string(out), pkg) {
		return true, ""
	}
	return false, ""
}

func pipxInstall(name, pkg string) error {
	ui.Doing("Installing package %s", pkg)
	cmd := exec.Command("pipx", "install", pkg)
	var stderr strings.Builder
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s", strings.TrimSpace(stderr.String()))
	}
	ui.Done("Successfully installed %s", name)
	return nil
}
