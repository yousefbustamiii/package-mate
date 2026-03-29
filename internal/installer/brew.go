package installer

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/yousefbustamiii/package-mate/internal/components"
	"github.com/yousefbustamiii/package-mate/internal/sys"
)

// IsInstalled checks whether an item is already present on the system
func IsInstalled(item components.InstallItem) Detection {
	// ── 1. Exact Matches (Formula/Cask) ────────────────────────────────────────
	if item.Formula != "" {
		if ok, ver := isFormulaInstalled(item.Formula); ok {
			if outdated, current, latest := CheckOutdated(item); outdated {
				return Detection{Status: DetectionOutdated, Detail: fmt.Sprintf("%s -> %s", current, latest)}
			}
			return Detection{Status: DetectionExact, Detail: ver}
		}
	}

	if item.Cask != "" {
		if ok, ver := isCaskInstalled(item.Cask); ok {
			if outdated, current, latest := CheckOutdated(item); outdated {
				return Detection{Status: DetectionOutdated, Detail: fmt.Sprintf("%s -> %s", current, latest)}
			}
			return Detection{Status: DetectionExact, Detail: ver}
		}

		// ❯ Check for Unmanaged App (Applications)
		if path, ok := sys.AppExists(item.Name); ok {
			return Detection{Status: DetectionManualApp, Detail: "Installed at " + path, BinaryPath: path}
		}

		// ❯ Check Trash
		if path, ok := sys.AppInTrash(item.Name); ok {
			return Detection{Status: DetectionTrashedApp, Detail: "Found in Trash", BinaryPath: path}
		}
	}

	if item.Special != "" {
		if ok, ver := isSpecialInstalled(item); ok {
			return Detection{Status: DetectionExact, Detail: ver}
		}
	}

	// ── 2. Binary Match fallback ──────────────────────────────────────────────
	if item.Binary != "" {
		if path, err := exec.LookPath(item.Binary); err == nil {
			prefix := sys.BrewPrefix()
			// If it's inside the Homebrew prefix and in a managed subdirectory (Cellar/opt/Caskroom)
			if strings.HasPrefix(path, prefix) && (strings.Contains(path, "/Cellar/") || strings.Contains(path, "/opt/") || strings.Contains(path, "/Caskroom/")) {
				formula := ""
				isCask := false
				parts := strings.Split(path, "/")
				for i, p := range parts {
					// ❯ Check for managed subdirs and grab the next segment
					if (p == "opt" || p == "Cellar" || p == "Caskroom") && i+1 < len(parts) {
						formula = parts[i+1]
						isCask = p == "Caskroom"
						break
					}
				}
				return Detection{Status: DetectionDifferentBrew, Detail: "Installed at " + path, BinaryPath: path, BrewFormula: formula, IsBrewCask: isCask}
			}
			// Truly unmanaged (e.g. ~/.local/bin or /usr/bin)
			return Detection{Status: DetectionBinary, Detail: "Installed at " + path, BinaryPath: path}
		}
	}

	return Detection{Status: DetectionNotFound}
}

// Install installs a single InstallItem idempotently and returns a Result.
func Install(item components.InstallItem) Result {
	name := item.Name

	if item.Formula != "" {
		if ok, ver := isFormulaInstalled(item.Formula); ok {
			return Result{ItemName: name, Status: StatusAlreadyHave, Version: ver}
		}
	}
	if item.Cask != "" {
		if ok, ver := isCaskInstalled(item.Cask); ok {
			return Result{ItemName: name, Status: StatusAlreadyHave, Version: ver}
		}
	}
	if item.Binary != "" && item.Formula == "" && item.Cask == "" {
		if path, err := exec.LookPath(item.Binary); err == nil {
			return Result{ItemName: name, Status: StatusAlreadyHave, Version: "Found at " + path}
		}
	}

	if item.Special != "" {
		return dispatchSpecial(item)
	}
	if item.Formula != "" {
		return installFormula(name, item.Formula)
	}
	if item.Cask != "" {
		return installCask(name, item.Cask)
	}
	return Result{ItemName: name, Status: StatusFailed, Err: fmt.Errorf("no install method defined")}
}

// Uninstall removes a single InstallItem idempotently.
func Uninstall(item components.InstallItem) error {
	if item.Formula != "" {
		return uninstallFormula(item.Name, item.Formula, false)
	}
	if item.Cask != "" {
		return uninstallCask(item.Name, item.Cask, false)
	}
	return fmt.Errorf("no uninstall method defined for %s", item.Name)
}

// UninstallForce removes a single InstallItem ignoring dependencies.
func UninstallForce(item components.InstallItem) error {
	if item.Formula != "" {
		return uninstallFormula(item.Name, item.Formula, true)
	}
	if item.Cask != "" {
		return uninstallCask(item.Name, item.Cask, true)
	}
	return fmt.Errorf("no uninstall method defined for %s", item.Name)
}
