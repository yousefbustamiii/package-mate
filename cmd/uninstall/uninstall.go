package uninstall

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/yousefbustamiii/package-mate/internal/components"
	"github.com/yousefbustamiii/package-mate/internal/installer"
	"github.com/yousefbustamiii/package-mate/internal/sys"
	"github.com/yousefbustamiii/package-mate/internal/ui"
)

// Run executes the uninstallation logic for a resolved item.
func Run(item *components.InstallItem) error {
	// ── 1. App/Cask Protection ─────────────────────────────────────────────────
	// If it's a Cask, we stick to the original simple logic as requested.
	if item.Cask != "" {
		if !ui.PromptConfirmation("DELETE", item.Name) {
			return nil
		}

		ui.Header(fmt.Sprintf("Uninstalling %s", item.Name))

		if sys.IsRunning(item.Name) {
			ui.ShowAppRunning(item.Name)
			ui.Footer()
			return nil
		}

		det := installer.IsInstalled(*item)
		if det.Status == installer.DetectionManualApp || det.Status == installer.DetectionTrashedApp {
			ui.ShowManualAppUninstall(item.Name)
			ui.Footer()
			return nil
		}
		if err := installer.Uninstall(*item); err != nil {
			ui.Fail("Failed to uninstall %s: %s", item.Name, err)
			ui.Footer()
			return nil
		}
		ui.Footer()
		return nil
	}

	// ── 2. CLI Tool Logic (Formulae) ───────────────────────────────────────────
	versions, _ := installer.GetAllVersions(*item)
	if len(versions) == 0 {
		ui.Fail("No installed versions of %s detected.", item.Name)
		return nil
	}

	// Prepare the interactive menu
	var opts []ui.SelectionOption
	for _, v := range versions {
		label := ""
		switch v.Type {
		case installer.VersionManaged:
			label = ui.C(ui.Bold+ui.White, "Managed") + ":                 " + ui.C(ui.Cyan, v.Path) + " " + ui.C(ui.Dim, "("+v.Version+")")
		case installer.VersionManagedOlder:
			label = ui.C(ui.Bold+ui.White, "Managed (older version)") + ": " + ui.C(ui.Dim, v.Path) + " " + ui.C(ui.Dim, "("+v.Version+")")
		case installer.VersionUnmanaged:
			label = ui.C(ui.Bold+ui.White, "Unmanaged") + ":               " + ui.C(ui.Yellow, v.Path)
		}
		opts = append(opts, ui.SelectionOption{Label: label, Value: ""})
	}

	ui.Header("Manage " + item.Name)
	choiceIdx := ui.PromptSelection("Choose a version to remove:", opts)
	if choiceIdx == -1 {
		return nil
	}

	selected := versions[choiceIdx]

	// Specific UX for unmanaged removal
	if selected.Type == installer.VersionUnmanaged {
		ui.Blank()
		ui.Hint(ui.C(ui.Bold+ui.Red, "!") + "  " + ui.C(ui.Dim, "[Sudo password will be required next]"))
	}

	if !ui.PromptConfirmation("DELETE", item.Name) {
		return nil
	}

	ui.Blank()
	ui.Doing("Uninstalling %s", item.Name)

	if selected.Type == installer.VersionUnmanaged {
		// ❯ Stop the spinner so the sudo password prompt doesn't glitch
		ui.Stop()
		ui.Blank()

		// ❯ Execute standard sudo rm with interactive terminal
		cmd := exec.Command("sudo", "rm", "-f", selected.Path)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = nil

		if err := cmd.Run(); err != nil {
			ui.Blank() // ❯ Add vertical space after Password prompt
			errStr := err.Error()
			if strings.Contains(errStr, "exit status 1") {
				// Standard rm failure (like SIP or missing file)
				ui.Fail("Failed to remove unmanaged binary at %s", selected.Path)
				ui.Blank() // ❯ Add space between failure and explanation

				if sys.IsProtectedPath(selected.Path) {
					ui.Hint("This file is protected by macOS System Integrity Protection (SIP).")
					ui.Hint("Package Mate is strictly forbidden from modifying system core files.")
				}
			} else {
				ui.Fail("Unexpected error during removal: %s", err)
			}
			ui.Footer()
			return nil
		}
		ui.Blank()
		ui.Done("Successfully removed unmanaged binary: %s", selected.Path)
	} else {
		if err := installer.UninstallFormula(selected.Formula, false); err != nil {
			if strings.Contains(err.Error(), "is required by") {
				if ui.PromptDependencyRemoval(item.Name, err.Error()) {
					if err := installer.UninstallFormula(selected.Formula, true); err != nil {
						ui.Fail("Force uninstall failed: %s", err)
					}
				}
				ui.Footer()
				return nil
			}
			ui.Fail("Failed to uninstall managed version: %s", err)
			ui.Footer()
			return nil
		}
		ui.Done("Successfully removed managed version.")
	}

	ui.Footer()
	return nil
}
