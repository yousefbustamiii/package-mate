package installer

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/yousefbustamiii/package-mate/internal/components"
	"github.com/yousefbustamiii/package-mate/internal/ui"
)

// ── Check helpers ──────────────────────────────────────────────────────────────

func isFormulaInstalled(formula string) (bool, string) {
	out, err := exec.Command(BrewExe(), "list", "--versions", formula).Output()
	if err != nil {
		return false, ""
	}
	ver := strings.TrimSpace(string(out))
	return ver != "", ver
}

func isCaskInstalled(cask string) (bool, string) {
	out, err := exec.Command(BrewExe(), "list", "--cask", "--versions", cask).Output()
	if err != nil {
		return false, ""
	}
	ver := strings.TrimSpace(string(out))
	return ver != "", ver
}

// ── Raw brew commands ──────────────────────────────────────────────────────────

func runBrewCmd(args ...string) error {
	cmd := exec.Command(BrewExe(), args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s", strings.TrimSpace(stderr.String()))
	}
	return nil
}

func brewInstallFormula(formula string) error {
	return runBrewCmd("install", formula)
}

func brewInstallCask(cask string) error {
	return runBrewCmd("install", "--cask", cask)
}

func brewUninstallFormula(formula string, ignoreDeps bool) error {
	args := []string{"uninstall", formula}
	if ignoreDeps {
		args = []string{"uninstall", "--ignore-dependencies", formula}
	}
	return runBrewCmd(args...)
}

func brewUninstallCask(cask string) error {
	return runBrewCmd("uninstall", "--cask", cask)
}

// UninstallFormula uninstalls a specific formula name directly.
func UninstallFormula(formula string, force bool) error {
	return brewUninstallFormula(formula, force)
}

// ── Orchestrated install / uninstall ──────────────────────────────────────────

func installFormula(name, formula string) Result {
	ui.Doing("Installing %s", name)
	if err := brewInstallFormula(formula); err != nil {
		return Result{ItemName: name, Status: StatusFailed, Err: err}
	}
	ui.Done("Successfully installed %s", name)
	return Result{ItemName: name, Status: StatusInstalled}
}

func installCask(name, cask string) Result {
	ui.Doing("Installing %s", name)
	if err := brewInstallCask(cask); err != nil {
		return Result{ItemName: name, Status: StatusFailed, Err: err}
	}
	ui.Done("Successfully installed %s", name)
	return Result{ItemName: name, Status: StatusInstalled}
}

func uninstallFormula(name, formula string, force bool) error {
	// ── 1. Execution ──────────────────────────────────────────────────────────
	if force {
		ui.Doing("Force removing %s", name)
	} else {
		ui.Doing("Removing %s", name)
	}

	if err := brewUninstallFormula(formula, force); err != nil {
		return err
	}
	ui.Done("Successfully removed %s", name)
	return nil
}

func uninstallCask(name, cask string, force bool) error {
	// ── 1. Execution ──────────────────────────────────────────────────────────
	ui.Doing("Removing %s", name)
	if err := brewUninstallCask(cask); err != nil {
		return err
	}
	ui.Done("Successfully removed %s", name)
	return nil
}

// Update runs brew upgrade for the given item.
func Update(item components.InstallItem) error {
	ui.Doing("Updating %s", item.Name)

	var args []string
	if item.Formula != "" {
		args = []string{"upgrade", item.Formula}
	} else if item.Cask != "" {
		args = []string{"upgrade", "--cask", item.Cask}
	} else {
		return fmt.Errorf("no update method defined")
	}

	if err := runBrewCmd(args...); err != nil {
		return err
	}

	ui.Done("Successfully updated %s", item.Name)
	return nil
}
