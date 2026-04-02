package bg

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/yousefbustamiii/package-mate/internal/background"
	"github.com/yousefbustamiii/package-mate/internal/components"
	"github.com/yousefbustamiii/package-mate/internal/installer"
	"github.com/yousefbustamiii/package-mate/internal/ui"
	"golang.org/x/term"
)

// Run launches the interactive background jobs viewer.
func Run() {
	background.CleanOld()
	jobs := background.LoadAll()

	ui.Blank()

	if len(jobs) == 0 {
		ui.Header("Background Jobs")
		fmt.Println("  " + ui.C(ui.Dim, "No background jobs found."))
		ui.Blank()
		fmt.Println("  " + ui.C(ui.Dim, "Start one by typing ") +
			ui.C(ui.BrightCyan, "INSTALL/BG") +
			ui.C(ui.Dim, ", ") +
			ui.C(ui.BrightCyan, "UPDATE/BG") +
			ui.C(ui.Dim, " or ") +
			ui.C(ui.BrightCyan, "OVERRIDE/BG") +
			ui.C(ui.Dim, " when prompted."))
		ui.Footer()
		return
	}

	showBgMenu(jobs)
}

// ── Interactive menu ────────────────────────────────────────────────────────

type bgState struct {
	jobs       []*background.Job
	selected   int
	totalLines int
}

func showBgMenu(jobs []*background.Job) {
	s := &bgState{jobs: jobs}

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		printFallback(jobs)
		return
	}

	fmt.Print("\033[?25l") // hide cursor

	cleanup := func() {
		_ = term.Restore(int(os.Stdin.Fd()), oldState)
		fmt.Print("\033[?25h")
		if s.totalLines > 0 {
			fmt.Printf("\033[%dA\033[J", s.totalLines)
		}
	}
	defer cleanup()

	for {
		// Reload jobs to reflect background process updates.
		s.jobs = background.LoadAll()
		if s.selected >= len(s.jobs) {
			s.selected = max(0, len(s.jobs)-1)
		}

		// Erase previous render then redraw.
		if s.totalLines > 0 {
			fmt.Printf("\033[%dA\033[J", s.totalLines)
		}
		s.totalLines = renderMenu(s.jobs, s.selected)

		var b [4]byte
		n, _ := os.Stdin.Read(b[:])
		if n == 0 {
			continue
		}

		if n >= 3 && b[0] == 27 && b[1] == '[' {
			switch b[2] {
			case 'A': // arrow up
				if s.selected > 0 {
					s.selected--
				}
			case 'B': // arrow down
				if s.selected < len(s.jobs)-1 {
					s.selected++
				}
			}
			continue
		}

		if n == 1 {
			switch b[0] {
			case 'k', 'K':
				if s.selected > 0 {
					s.selected--
				}
			case 'j', 'J':
				if s.selected < len(s.jobs)-1 {
					s.selected++
				}
			case 13: // Enter
				if s.selected < len(s.jobs) {
					j := s.jobs[s.selected]
					if j.Status == background.StatusRunning {
						doAbort(s, j, oldState)
					}
				}
			case 'q', 'Q', 3, 27: // q, Ctrl+C, Esc
				return
			}
		}
	}
}

// doAbort handles the abort confirmation flow.
// It temporarily exits raw mode to read the typed confirmation.
func doAbort(s *bgState, j *background.Job, oldState *term.State) {
	// Clear the menu before showing confirmation.
	if s.totalLines > 0 {
		fmt.Printf("\033[%dA\033[J", s.totalLines)
	}
	s.totalLines = 0

	// Restore terminal for line-buffered input.
	_ = term.Restore(int(os.Stdin.Fd()), oldState)
	fmt.Print("\033[?25h")

	// Confirmation prompt.
	ui.Blank()
	fmt.Println("  " + ui.C(ui.Yellow, strings.Repeat("─", 60)))
	fmt.Println("  " + ui.C(ui.Bold+ui.Yellow, "❯ ABORT BACKGROUND DOWNLOAD"))
	fmt.Println("  " + ui.C(ui.Yellow, strings.Repeat("─", 60)))
	ui.Blank()
	fmt.Printf("  " + ui.C(ui.Bold+ui.White, "This will abort the background download of ") +
		ui.C(ui.BrightCyan, j.Name) + ".\n")
	ui.Blank()
	fmt.Printf("  " + ui.C(ui.Bold+ui.White, "Type ") +
		ui.C(ui.Bold+ui.BrightCyan, "ABORT") +
		ui.C(ui.Bold+ui.White, " to confirm.") + "\n")
	ui.Blank()
	fmt.Print("  " + ui.C(ui.Bold+ui.White, "❯ "))

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := strings.TrimSpace(scanner.Text())

	if strings.EqualFold(input, "ABORT") {
		item, _, ok := components.Resolve(strings.ToLower(j.Name))
		if ok {
			scan := installer.PerformSystemScan()
			status, _, _ := installer.ResolveStatus(scan, *item)

			if status == components.StatusInstalled {
				ui.Blank()
				ui.Header("Abort Not Possible")
				ui.Warn("It's a bit too late to abort!")
				ui.Blank()
				fmt.Printf("  The installation of " + ui.C(ui.BrightCyan+ui.Bold, j.Name) + " has already completed successfully.\n")
				fmt.Printf("  Since the tool is now fully active on your system, an abort is no longer possible.\n")
				ui.Blank()
				fmt.Println("  " + ui.C(ui.Dim, "If you'd like to remove this tool, please use the main dashboard:"))
				fmt.Printf("  " + ui.C(ui.Dim, "Run ") + ui.C(ui.Cyan, "mate "+strings.ToLower(j.Name)) + ui.C(ui.Dim, " and choose option ") + ui.C(ui.White, "2 (Uninstall)") + ui.C(ui.Dim, ".") + "\n")
				ui.Blank()
				ui.Footer()
				fmt.Println("  " + ui.C(ui.Dim, "Press any key to return..."))

				// Re-enter raw mode temporarily to wait for any key
				rawState, _ := term.MakeRaw(int(os.Stdin.Fd()))
				var sink [4]byte
				_, _ = os.Stdin.Read(sink[:])
				_ = term.Restore(int(os.Stdin.Fd()), rawState)
			} else {
				if err := j.Abort(); err != nil {
					ui.Fail("Could not abort: %v", err)
				} else {
					ui.Blank()
					ui.Done("Aborted: %s", j.Name)
				}
			}
		} else {
			// Item not found in catalog, proceed with normal abort
			if err := j.Abort(); err != nil {
				ui.Fail("Could not abort: %v", err)
			} else {
				ui.Blank()
				ui.Done("Aborted: %s", j.Name)
			}
		}
	} else {
		ui.Blank()
		fmt.Println("  " + ui.C(ui.Dim, "Abort cancelled."))
	}

	ui.Blank()
	fmt.Println("  " + ui.C(ui.Dim, "Press any key to return..."))

	// Re-enter raw mode for the next menu loop
	newState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return
	}
	*oldState = *newState
	fmt.Print("\033[?25l")

	// Consume the "any key" press.
	var sink [4]byte
	_, _ = os.Stdin.Read(sink[:])
}

// ── Rendering ───────────────────────────────────────────────────────────────

// renderMenu draws the job list and returns the number of lines printed.
func renderMenu(jobs []*background.Job, selected int) int {
	lines := 0

	running := 0
	for _, j := range jobs {
		if j.Status == background.StatusRunning {
			running++
		}
	}

	headerSuffix := ""
	if running > 0 {
		headerSuffix = "  " + ui.C(ui.Yellow, fmt.Sprintf("(%d downloading)", running))
	}

	emit := func(s string) {
		fmt.Printf("\r\033[K%s\r\n", s)
		lines++
	}

	emit(ui.C(ui.Grey, "  "+strings.Repeat("─", 60)))
	emit("  " + ui.C(ui.Bold+ui.White, "❯ BACKGROUND JOBS") + headerSuffix)
	emit(ui.C(ui.Grey, "  "+strings.Repeat("─", 60)))
	emit("")

	for i, j := range jobs {
		isSel := i == selected

		// Cursor is always 2 visual chars — content never shifts.
		cursor := "  "
		nameStyle := ui.Dim
		if isSel {
			cursor = ui.C(ui.BrightCyan, "❯ ")
			nameStyle = ui.Bold + ui.White
		}

		icon, iconCode := statusIcon(j.Status)

		// Pad raw strings BEFORE applying ANSI — prevents escape-byte inflation
		// from breaking column alignment when the style changes on selection.
		nameStr := fmt.Sprintf("%-20s", truncate(j.Name, 20))
		labelStr := fmt.Sprintf("%-12s", statusLabel(j.Status))
		tail := statusTail(j)

		row := fmt.Sprintf("  %s[%s]  %s  %s  %s",
			cursor,
			ui.C(iconCode, icon),
			ui.C(nameStyle, nameStr),
			ui.C(statusLabelColor(j.Status), labelStr),
			tail,
		)
		emit(row)

		// Show error/note for the selected failed/warned job.
		if isSel && j.Error != "" && j.Status != background.StatusRunning {
			note := j.Error
			if len(note) > 55 {
				note = note[:52] + "..."
			}
			errCode := ui.Red
			if j.Status == background.StatusDone {
				errCode = ui.Yellow
			}
			emit("          " + ui.C(errCode, note))
		}
	}

	emit("")
	emit(ui.C(ui.Grey, "  "+strings.Repeat("─", 60)))
	emit(
		ui.C(ui.Dim, "  j/k ↑/↓") + ui.C(ui.Dim, "  Navigate   ") +
			ui.C(ui.White, "Enter") + ui.C(ui.Dim, "  Abort running job   ") +
			ui.C(ui.White, "q") + ui.C(ui.Dim, " / ") +
			ui.C(ui.White, "Esc") + ui.C(ui.Dim, "  Quit"),
	)
	emit("")

	return lines
}

func statusIcon(s background.JobStatus) (string, string) {
	switch s {
	case background.StatusRunning:
		return "↺", ui.Yellow
	case background.StatusDone:
		return "✓", ui.BrightGreen
	case background.StatusFailed:
		return "✗", ui.Red
	case background.StatusAborted:
		return "⊘", ui.Grey
	default:
		return "?", ui.Dim
	}
}

func statusLabel(s background.JobStatus) string {
	switch s {
	case background.StatusRunning:
		return "Installing"
	case background.StatusDone:
		return "Installed"
	case background.StatusFailed:
		return "Failed"
	case background.StatusAborted:
		return "Aborted"
	default:
		return ""
	}
}

func statusLabelColor(s background.JobStatus) string {
	switch s {
	case background.StatusRunning:
		return ui.Yellow
	case background.StatusDone:
		return ui.BrightGreen
	case background.StatusFailed:
		return ui.Red
	case background.StatusAborted:
		return ui.Grey
	default:
		return ui.Dim
	}
}

func statusTail(j *background.Job) string {
	switch j.Status {
	case background.StatusRunning:
		return ui.C(ui.Yellow, j.Elapsed())
	case background.StatusDone:
		if j.Error != "" {
			return ui.C(ui.Yellow, "see note ↓")
		}
		var total time.Duration
		if !j.FinishedAt.IsZero() {
			total = j.FinishedAt.Sub(j.StartedAt).Round(time.Second)
		} else {
			total = time.Since(j.StartedAt).Round(time.Second)
		}
		return ui.C(ui.BrightGreen, "done") + ui.C(ui.Dim, " in "+formatDuration(total))
	default:
		return ""
	}
}

func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	return fmt.Sprintf("%dm%ds", int(d.Minutes()), int(d.Seconds())%60)
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "…"
}

func printFallback(jobs []*background.Job) {
	ui.Header("Background Jobs")
	for _, j := range jobs {
		icon, _ := statusIcon(j.Status)
		fmt.Printf("  [%s] %-24s  %s\n", icon, j.Name, j.Status)
	}
	ui.Footer()
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
