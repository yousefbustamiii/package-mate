package ui

import (
	"fmt"
	"strings"

	"github.com/yousefbustamiii/package-mate/internal/components"
)

// ── Section listing ────────────────────────────────────────────────────────────

// SectionEntry is a single tool entry for display purposes.
type SectionEntry struct {
	Name        string // display name e.g. "MySQL"
	HasMultiple bool   // True if 2+ installations found
	Desc        string
	Color       string                     // raw ANSI prefix from RGB()
	Status      components.DashboardStatus // Current status
}

// SectionGroup is a named group of entries for the help listing.
type SectionGroup struct {
	Label   string
	Entries []SectionEntry
}

// ShowHomebrewMissing prints a friendly prompt when Homebrew is not found.
func ShowHomebrewMissing() {
	Banner()

	mateColor := RGB(255, 170, 50)

	fmt.Println("  " + C(Bold+White, "Install ") + C(Bold+mateColor, "Homebrew") + C(Bold+White, " to be able to use ") + C(Bold+Cyan, "Package Mate"))
	Blank()
	Blank()

	fmt.Println("  " + C(Bold+Grey, "❯ run ") + C(Bold+Cyan, "mate homebrew") + C(Bold+Grey, " to install ") + C(Bold+mateColor, "Homebrew"))
	Blank()
	Blank()
}

// ShowAllTools prints tool sections as horizontal columns with vertical lists underneath.
func ShowAllTools(groups []SectionGroup) {
	Banner()

	fmt.Println("  " + C(Bold+White, "usage: ") + C(Cyan, "mate <tool> ") + C(Dim, "(interactive menu)"))
	Blank()

	totalCount := 0
	installedCount := 0
	for _, g := range groups {
		totalCount += len(g.Entries)
		for _, e := range g.Entries {
			if e.Status != components.StatusNotInstalled {
				installedCount++
			}
		}
	}
	fmt.Printf("  "+C(Bold+White, "Installed: ")+C(BrightGreen+Bold, "%d")+" / "+C(Bold+White, "%d")+"\n", installedCount, totalCount)
	Blank()
	fmt.Printf("  " + C(BrightGreen+Bold, "✓") + "  " + C(Dim, "---> Installed & Managed (up to date)") + "\n")
	Blank()
	fmt.Printf("  " + C(Yellow+Bold, "↻") + "  " + C(Dim, "---> Installed & Managed, update available") + "\n")
	Blank()
	fmt.Printf("  " + C(BrightCyan+Bold, "⚙") + "  " + C(Dim, "---> Installed & Unmanaged") + "\n")
	Blank()
	fmt.Printf("  " + C(Bold+White, "?") + "  " + C(Dim, "---> Multiple Installations") + "\n")
	Blank()

	if len(groups) == 0 {
		return
	}

	// Constants for status prefix
	statusLen := 4 // "[✓] " or "[✗] "

	// 1. Determine column widths and max height
	colWidths := make([]int, len(groups))
	maxHeight := 0
	for i, g := range groups {
		// Header width calculation
		headerW := len(g.Label) + 4 // "[ ]"
		if len(g.Entries) > maxHeight {
			maxHeight = len(g.Entries)
		}
		w := headerW
		for _, e := range g.Entries {
			// Item width is statusPrefix + Name
			itemW := statusLen + len(e.Name)
			if itemW > w {
				w = itemW
			}
		}
		colWidths[i] = w + 4 // 4 chars spacing between columns
	}

	// 2. Print Headers horizontally
	fmt.Print("  ")
	for i, g := range groups {
		header := C(Bold+White, "[ "+strings.ToUpper(g.Label)+" ]")
		padding := colWidths[i] - (len(g.Label) + 4)
		fmt.Printf("%s%s", header, strings.Repeat(" ", padding))
	}
	Blank()
	fmt.Println("  " + C(Grey, strings.Repeat("─", termWidth()-4)))
	Blank()

	// 3. Print Items row-by-row
	for row := 0; row < maxHeight; row++ {
		fmt.Print("  ")
		for col, g := range groups {
			if row < len(g.Entries) {
				e := g.Entries[row]

				// Status Indicator
				statusPrefix := C(Dim, "[ ] ")
				switch e.Status {
				case components.StatusNotInstalled:
					// Default "[ ]" prefix already set above.
				case components.StatusInstalled:
					statusIcon := C(BrightGreen+Bold, "✓")
					statusPrefix = C(Dim, "[") + statusIcon + C(Dim, "] ")
				case components.StatusOutdated:
					statusIcon := C(Yellow+Bold, "↻")
					statusPrefix = C(Dim, "[") + statusIcon + C(Dim, "] ")
				case components.StatusUnmanaged:
					statusIcon := C(BrightCyan+Bold, "⚙")
					statusPrefix = C(Dim, "[") + statusIcon + C(Dim, "] ")
				}

				nameStr := C(e.Color+Bold, e.Name)
				if e.HasMultiple {
					nameStr += " " + C(Bold+White, "?")
				}

				paddingLen := colWidths[col] - (statusLen + len(e.Name))
				if e.HasMultiple {
					paddingLen -= 2
				}
				if paddingLen < 0 {
					paddingLen = 0
				}

				fmt.Printf("%s%s%s", statusPrefix, nameStr, strings.Repeat(" ", paddingLen))
			} else {
				fmt.Print(strings.Repeat(" ", colWidths[col]))
			}
		}
		Blank()
		Blank() // Extra space between rows
		Blank()
	}
	Blank()
}
