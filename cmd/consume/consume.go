package consume

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/yousefbustamiii/package-mate/internal/background"
	"github.com/yousefbustamiii/package-mate/internal/components"
	"github.com/yousefbustamiii/package-mate/internal/installer"
	"github.com/yousefbustamiii/package-mate/internal/ui"
)

// Run executes the consume command.
func Run() error {
	ui.Header("Consume Export")
	fmt.Println("  " + ui.C(ui.Dim, "Paste an export string from another Mac to sync your tools."))
	ui.Blank()

	encoded := ui.PromptInput("Export String:")
	if encoded == "" {
		ui.Blank()
		ui.Fail("No string provided. Aborting sync.")
		return nil
	}

	// 1. Decode & Parse
	decoded, err := base64.StdEncoding.DecodeString(strings.TrimSpace(encoded))
	if err != nil {
		ui.Blank()
		ui.Fail("Invalid export string (the format doesn't look quite right).")
		return nil
	}

	var names []string
	if err := json.Unmarshal(decoded, &names); err != nil {
		ui.Blank()
		ui.Fail("Invalid export string (the package list appears to be corrupted).")
		return nil
	}

	if len(names) == 0 {
		ui.Blank()
		ui.Fail("This export string doesn't contain any packages to install.")
		return nil
	}

	ui.Blank()
	ui.Doing("Analyzing %d packages", len(names))
	scan := installer.PerformSystemScan()
	ui.Stop()

	// 2. Resolve items and determine actions
	type task struct {
		item    *components.InstallItem
		action  string // install | update | skip
		reason  string
		display string
	}

	var tasks []task
	seen := make(map[string]bool)

	for _, name := range names {
		if seen[name] {
			continue
		}
		seen[name] = true

		item, _, ok := components.Resolve(name)
		if !ok {
			// If not in our catalog, we can't manage it reliably here
			continue
		}

		status, _, _ := installer.ResolveStatus(scan, *item)
		t := task{item: item}

		// Priority check: Is it already being handled in the background?
		if running := background.GetRunningJob(item.Name); running != nil {
			t.action = "skip"
			t.display = ui.C(ui.Dim+ui.Strikethrough, item.Name) + " " + ui.C(ui.Dim, "---- (installing in background)")
		} else {
			switch status {
			case components.StatusInstalled:
				// Already installed and up to date
				t.action = "skip"
				t.display = ui.C(ui.Dim+ui.Strikethrough, item.Name) + " " + ui.C(ui.Dim, "---- (already installed)")
			
			case components.StatusOutdated:
				t.action = "update"
				t.display = ui.C(ui.Bold+ui.White, item.Name) + " " + ui.C(ui.Yellow, "[exists, will update only]")
			
			case components.StatusUnmanaged:
				// App exists but not managed by brew
				t.action = "skip"
				t.display = ui.C(ui.Dim+ui.Strikethrough, item.Name) + " " + ui.C(ui.Dim, "---- (manual app exists, skipping)")
			
			case components.StatusNotInstalled:
				t.action = "install"
				t.display = ui.C(ui.Bold+ui.White, item.Name)
			}
		}

		tasks = append(tasks, t)
	}

	// Sort tasks for display (installs first, then updates, then skips)
	sort.Slice(tasks, func(i, j int) bool {
		priority := func(t task) int {
			switch t.action {
			case "install": return 1
			case "update": return 2
			case "skip": return 3
			default: return 4
			}
		}
		if priority(tasks[i]) != priority(tasks[j]) {
			return priority(tasks[i]) < priority(tasks[j])
		}
		return tasks[i].item.Name < tasks[j].item.Name
	})

	ui.Header("Consumption Plan")
	for _, t := range tasks {
		bullet := ui.C(ui.BrightGreen, "❯")
		if t.action == "skip" {
			bullet = ui.C(ui.Dim, "•")
		}
		fmt.Printf("  %s %s\n", bullet, t.display)
	}
	ui.Blank()

	// 3. Confirm and execute
	toProcess := 0
	for _, t := range tasks {
		if t.action == "install" || t.action == "update" {
			toProcess++
		}
	}

	if toProcess == 0 {
		ui.Done("All tools are already present and up to date. Nothing to do!")
		ui.Footer()
		return nil
	}

	fmt.Printf("  " + ui.C(ui.Bold+ui.White, "Ready to process ") + ui.C(ui.BrightCyan+ui.Bold, fmt.Sprintf("%d packages", toProcess)) + " in the background.\n")
	ui.Blank()

	if !ui.PromptConfirmation("CONSUME", fmt.Sprintf("%d packages", toProcess)) {
		return nil
	}

	// 4. Enqueue background jobs
	ui.Blank()
	ui.Doing("Enqueuing background jobs")
	
	for _, t := range tasks {
		if t.action == "install" || t.action == "update" {
			// Note: consume always uses background mode
			_, err := background.Enqueue(t.item.Name, t.action, "", "", false)
			if err != nil {
				ui.Warn("Failed to queue %s: %v", t.item.Name, err)
			}
		}
	}
	ui.Stop()

	ui.Blank()
	ui.Done("Successfully queued %d background jobs!", toProcess)
	ui.Blank()
	fmt.Printf("  %s %s\n", ui.C(ui.Dim, "Run"), ui.C(ui.Cyan, "mate bg"))
	fmt.Printf("  %s %s\n", ui.C(ui.Dim, "to monitor progress."), "")
	ui.Blank()
	ui.Footer()

	return nil
}
