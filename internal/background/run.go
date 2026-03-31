package background

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/yousefbustamiii/package-mate/internal/components"
	"github.com/yousefbustamiii/package-mate/internal/installer"
	"github.com/yousefbustamiii/package-mate/internal/sys"
)

// Execute is called by the background child process (mate _bg-exec <jobID>).
// It runs the queued installation and updates the job state file.
func Execute(jobID string) error {
	j, err := LoadJob(filepath.Join(Dir(), jobID+".json"))
	if err != nil {
		return fmt.Errorf("job %s not found: %w", jobID, err)
	}

	// Confirm we're the right process.
	j.PID = os.Getpid()
	j.Status = StatusRunning
	_ = j.Save()

	item, _, ok := components.Resolve(strings.ToLower(j.Name))
	if !ok {
		j.Status = StatusFailed
		j.Error = fmt.Sprintf("tool %q not found in catalog", j.Name)
		return j.Save()
	}

	var runErr error
	var noteErr string

	switch j.Action {
	case "install":
		res := installer.Install(*item)
		if res.Status == installer.StatusFailed {
			runErr = res.Err
		}

	case "update":
		runErr = installer.Update(*item)

	case "install-override":
		res := installer.Install(*item)
		if res.Status == installer.StatusFailed {
			runErr = res.Err
		} else if res.Status == installer.StatusInstalled && j.OldBinaryPath != "" {
			if !sys.IsProtectedPath(j.OldBinaryPath) {
				cmd := exec.Command("rm", "-f", j.OldBinaryPath)
				var stderr bytes.Buffer
				cmd.Stderr = &stderr
				if err := cmd.Run(); err != nil {
					msg := strings.TrimSpace(stderr.String())
					// Cannot run sudo in background — note it.
					noteErr = fmt.Sprintf("installed OK; could not remove old binary at %s (%s) — remove manually with: sudo rm %s",
						j.OldBinaryPath, msg, j.OldBinaryPath)
				}
			}
		}

	case "brew-override":
		res := installer.Install(*item)
		if res.Status == installer.StatusFailed {
			runErr = res.Err
		} else if res.Status == installer.StatusInstalled && j.OldFormula != "" {
			oldItem := components.InstallItem{Name: j.OldFormula}
			if j.IsOldCask {
				oldItem.Cask = j.OldFormula
			} else {
				oldItem.Formula = j.OldFormula
			}
			if err := installer.Uninstall(oldItem); err != nil {
				if strings.Contains(err.Error(), "is required by") {
					_ = installer.UninstallForce(oldItem)
				}
				// Even if uninstall fails, new version is installed — not a fatal error.
			}
		}

	default:
		runErr = fmt.Errorf("unknown action: %s", j.Action)
	}

	if runErr != nil {
		j.Status = StatusFailed
		if j.Error == "" {
			j.Error = runErr.Error()
		}
	} else {
		j.Status = StatusDone
		if noteErr != "" {
			j.Error = noteErr // Surfaced as a warning in mate bg
		}
	}

	j.FinishedAt = time.Now()
	j.NotifyFinish()
	return j.Save()
}
