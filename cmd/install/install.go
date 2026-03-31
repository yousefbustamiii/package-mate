package install

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/yousefbustamiii/package-mate/cmd/info"
	"github.com/yousefbustamiii/package-mate/cmd/uninstall"
	"github.com/yousefbustamiii/package-mate/internal/background"
	"github.com/yousefbustamiii/package-mate/internal/components"
	"github.com/yousefbustamiii/package-mate/internal/installer"
	"github.com/yousefbustamiii/package-mate/internal/sys"
	"github.com/yousefbustamiii/package-mate/internal/ui"
)

// Run executes the installation or update logic for a resolved item.
func Run(item *components.InstallItem, sec *components.Section) error {
	ui.Blank()
	ui.Doing("Analyzing %s", item.Name)
	det := installer.IsInstalled(*item)

	var isUpdate bool
	var isOverride bool
	var isBrewOverride bool
	var isOldCask bool
	var isBg bool
	var oldBinaryPath string
	var oldFormulaName string

	switch det.Status {
	case installer.DetectionExact:
		ui.AlreadyInstalled(item.Name, det.Detail)
		return nil

	case installer.DetectionNotFound:
		res := ui.PromptInstall(item.Name)
		if res == "" {
			return nil
		}
		if res == "INSTALL/BG" {
			isBg = true
		}

	case installer.DetectionOutdated:
		parts := strings.Split(det.Detail, " -> ")
		if len(parts) != 2 {
			return nil
		}
		res := ui.PromptOutdated(item.Name, parts[0], parts[1])
		if res == "" {
			return nil
		}
		isUpdate = true
		if res == "UPDATE/BG" {
			isBg = true
		}

	case installer.DetectionBinary:
		target := item.Formula
		if target == "" {
			target = item.Cask
		}
		isProtected := sys.IsProtectedPath(det.BinaryPath)
		isSystem := sys.IsSystemPath(det.BinaryPath)
		res := ui.PromptOverride(item.Name, target, det.BinaryPath, isProtected, isSystem)
		if res == "" {
			return nil
		}
		if res == "OVERRIDE" || res == "OVERRIDE/BG" {
			isOverride = true
			oldBinaryPath = det.BinaryPath
		}
		if strings.HasSuffix(res, "/BG") {
			isBg = true
		}

	case installer.DetectionDifferentBrew:
		target := item.Formula
		if target == "" {
			target = item.Cask
		}
		res := ui.PromptDifferentBrew(target, det.BrewFormula)
		if res == "" {
			return nil
		}
		if res == "UPDATE" || res == "UPDATE/BG" {
			isUpdate = true
		}
		if res == "OVERRIDE" || res == "OVERRIDE/BG" {
			isOverride = true
			isBrewOverride = true
			isOldCask = det.IsBrewCask
			oldFormulaName = det.BrewFormula
		}
		if strings.HasSuffix(res, "/BG") {
			isBg = true
		}

	case installer.DetectionManualApp, installer.DetectionTrashedApp:
		ui.ShowManualAppUninstall(item.Name)
		ui.Footer()
		return nil
	}

	// ── Safety: check running app before touching casks ─────────────────────────
	if item.Cask != "" && sys.IsRunning(item.Name) {
		ui.ShowAppRunning(item.Name)
		ui.Footer()
		return nil
	}

	// ── Background mode ─────────────────────────────────────────────────────────
	if isBg {
		action := bgAction(isUpdate, isOverride, isBrewOverride)
		_, err := background.Enqueue(item.Name, action, oldBinaryPath, oldFormulaName, isOldCask)
		if err != nil {
			ui.Fail("Could not queue background job: %v", err)
			ui.Footer()
			return nil
		}
		ui.ShowBgQueued(item.Name)
		return nil
	}

	// ── Foreground mode ─────────────────────────────────────────────────────────
	verb := "Installing"
	if isUpdate {
		verb = "Updating"
	}
	ui.Header(fmt.Sprintf("%s %s", verb, item.Name))

	var res installer.Result
	if isUpdate {
		err := installer.Update(*item)
		res = installer.Result{
			ItemName: item.Name,
			Status:   installer.StatusInstalled,
			Err:      err,
		}
		if err != nil {
			res.Status = installer.StatusFailed
		}
	} else {
		res = installer.Install(*item)
	}

	switch res.Status {
	case installer.StatusInstalled:
		// fall through to override cleanup
	case installer.StatusAlreadyHave:
		// Already present — nothing more to do.
		ui.Footer()
		return nil
	case installer.StatusFailed:
		errMsg := "unknown error"
		if res.Err != nil {
			errMsg = res.Err.Error()
		}
		ui.Fail("Failed to %s %s: %s", strings.ToLower(verb), item.Name, errMsg)
		ui.Footer()
		return nil
	}

	// ── Override cleanup ────────────────────────────────────────────────────────
	if isOverride && res.Status == installer.StatusInstalled {
		ui.Blank()

		if isBrewOverride {
			ui.Doing("Removing old brew item: %s", oldFormulaName)

			oldItem := components.InstallItem{Name: oldFormulaName}
			if isOldCask {
				oldItem.Cask = oldFormulaName
			} else {
				oldItem.Formula = oldFormulaName
			}

			err := installer.Uninstall(oldItem)
			if err != nil {
				if strings.Contains(err.Error(), "is required by") {
					if ui.PromptDependencyRemoval(oldFormulaName, err.Error()) {
						ui.Doing("Force removing %s", oldFormulaName)
						if err := installer.UninstallForce(oldItem); err != nil {
							ui.Fail("Force removal failed: %v", err)
						} else {
							ui.Done("Successfully removed old formula %s", oldFormulaName)
						}
					} else {
						ui.Warn("Removal skipped. Legacy version %s remains.", oldFormulaName)
					}
				} else {
					ui.Warn("Failed to remove old brew item: %v", err)
				}
			} else {
				ui.Done("Successfully removed old brew item: %s", oldFormulaName)
			}
		} else {
			if sys.IsProtectedPath(oldBinaryPath) {
				ui.ShowSafetyAlert(oldBinaryPath)
				ui.Footer()
				return nil
			}

			ui.Doing("Removing old binary")

			cmd := exec.Command("rm", "-f", oldBinaryPath)
			var stderr bytes.Buffer
			cmd.Stderr = &stderr

			if err := cmd.Run(); err != nil {
				msg := strings.TrimSpace(stderr.String())
				if strings.Contains(msg, "Permission denied") {
					ui.Doing("Retrying with sudo")
					ui.Blank()

					sudoCmd := exec.Command("sudo", "rm", "-f", oldBinaryPath)
					sudoCmd.Stdin = os.Stdin
					sudoCmd.Stdout = os.Stdout
					sudoCmd.Stderr = os.Stderr

					if err := sudoCmd.Run(); err != nil {
						ui.Fail("Sudo removal failed: %s", err)
						ui.Hint("You may need to remove it manually: sudo rm %s", oldBinaryPath)
					} else {
						ui.Done("Successfully removed unmanaged binary")
					}
				} else {
					ui.Warn("Failed to remove old binary at %s", oldBinaryPath)
					if msg != "" {
						ui.Hint("Error: %s", msg)
					}
					ui.Hint("You may need to remove it manually: sudo rm %s", oldBinaryPath)
				}
			} else {
				ui.Done("Successfully removed unmanaged binary: %s", oldBinaryPath)
			}
		}
	}

	ui.Footer()
	return nil
}

// bgAction maps install flags to a background job action string.
func bgAction(isUpdate, isOverride, isBrewOverride bool) string {
	if isUpdate {
		return "update"
	}
	if isOverride && isBrewOverride {
		return "brew-override"
	}
	if isOverride {
		return "install-override"
	}
	return "install"
}

// LaunchSearch runs the interactive search UI, scanning the system once upfront.
func LaunchSearch() {
	scan := installer.PerformSystemScan()

	var wg sync.WaitGroup

	groups := make([]ui.SectionGroup, len(components.AllSections))

	for i, sec := range components.AllSections {
		groups[i] = ui.SectionGroup{Label: sec.Name}
		groups[i].Entries = make([]ui.SectionEntry, len(sec.Items))

		for j, item := range sec.Items {
			groups[i].Entries[j] = ui.SectionEntry{
				Name:   item.Name,
				Desc:   item.Desc,
				Color:  item.Color,
				Status: components.StatusNotInstalled,
			}

			wg.Add(1)
			go func(i, j int, it components.InstallItem) {
				defer wg.Done()
				status, multiple := installer.ResolveStatus(scan, it)
				groups[i].Entries[j].Status = status
				groups[i].Entries[j].HasMultiple = multiple
			}(i, j, item)
		}
	}

	wg.Wait()

	ui.ShowSearch(groups, func(itemName string, choice int) {
		item, sec, ok := components.Resolve(strings.ToLower(itemName))
		if !ok {
			return
		}
		switch choice {
		case 1:
			_ = Run(item, sec)
		case 2:
			_ = uninstall.Run(item)
		case 3:
			_ = info.Run(item, sec)
		}
	})
}
