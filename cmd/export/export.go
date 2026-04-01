package export

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/yousefbustamiii/package-mate/internal/components"
	"github.com/yousefbustamiii/package-mate/internal/installer"
	"github.com/yousefbustamiii/package-mate/internal/ui"
)

// Run executes the export command.
func Run() error {
	ui.Blank()

	// ── 1. Perform one-time system scan (LaunchSearch-style) ───────────────────
	ui.Doing("Scanning system for managed packages")
	scan := installer.PerformSystemScan()
	ui.Stop()

	// Collect all managed (installed) Homebrew items from our catalog
	var installedItems []string
	seen := make(map[string]bool) // Avoid duplicates

	// Iterate through all sections and items in our catalog
	for _, section := range components.AllSections {
		for _, item := range section.Items {
			// Skip items that aren't managed by Homebrew
			if item.Formula == "" && item.Cask == "" {
				continue
			}

			// Check if the item is installed using the scan data
			status, _, isRequested := installer.ResolveStatus(scan, item)

			// Include only if it's managed (Installed or Outdated) AND was explicitly requested
			// (This filters out dependencies like 'lua' when you only asked for 'vim')
			if (status == components.StatusInstalled || status == components.StatusOutdated) && isRequested {
				// Use Formula if available, otherwise use Cask
				id := item.Formula
				if id == "" {
					id = item.Cask
				}
				if id != "" && !seen[id] {
					installedItems = append(installedItems, id)
					seen[id] = true
				}
			}
		}
	}

	// Sort for consistent output
	sort.Strings(installedItems)

	// Check if we found anything
	if len(installedItems) == 0 {
		fmt.Println("  " + ui.C(ui.Grey, strings.Repeat("─", 60)))
		fmt.Println("  " + ui.C(ui.Bold+ui.Yellow, "✓ NO MANAGED PACKAGES FOUND"))
		fmt.Println(ui.C(ui.Grey, "  "+strings.Repeat("─", 60)))
		ui.Blank()
		fmt.Printf("  %s Your Homebrew installation is completely clean.\n", ui.C(ui.Yellow, "❯"))
		fmt.Printf("  %s Install some tools via "+ui.C(ui.Cyan, "mate")+" to get started.\n", ui.C(ui.Dim, "→"))
		ui.Blank()
		return nil
	}

	// Create JSON array
	jsonData, err := json.Marshal(installedItems)
	if err != nil {
		ui.Fail("Failed to encode formulae: %v", err)
		return err
	}

	// Base64 encode
	encoded := base64.StdEncoding.EncodeToString(jsonData)

	// Display the export string
	fmt.Println("  " + ui.C(ui.Grey, strings.Repeat("─", 60)))
	fmt.Println("  " + ui.C(ui.Bold+ui.BrightCyan, "❯ EXPORT STRING"))
	fmt.Println(ui.C(ui.Grey, "  "+strings.Repeat("─", 60)))
	ui.Blank()

	// Print the encoded string
	fmt.Println("  " + ui.C(ui.White, encoded))
	ui.Blank()

	// Show usage hint
	fmt.Println("  " + ui.C(ui.Dim, "Note: To consume this, run ") + ui.C(ui.Cyan, "mate consume") + ui.C(ui.Dim, " on your other Mac."))
	ui.Blank()

	// Show packaged export summary
	fmt.Println("  " + ui.C(ui.Grey, strings.Repeat("─", 60)))
	fmt.Printf("  %s %s\n", ui.C(ui.Bold+ui.BrightGreen, "✓ PACKAGED EXPORT"), ui.C(ui.Dim, fmt.Sprintf("(%d packages)", len(installedItems))))
	fmt.Println(ui.C(ui.Grey, "  "+strings.Repeat("─", 60)))
	ui.Blank()

	// Explain dependency exclusion
	fmt.Println("  " + ui.C(ui.Bold+ui.White, " Note:"))
	ui.Blank()
	fmt.Print("   " + ui.C(ui.Dim, "Only top-level tools you've explicitly installed are listed here."))
	fmt.Print("\n   " + ui.C(ui.Dim, "Even if a dependency has a checkmark in the dashboard, it is"))
	fmt.Print("\n   " + ui.C(ui.Dim, "omitted here to ensure a clean, portable setup for your next Mac."))
	ui.Blank()
	ui.Blank()

	// List all packages
	for _, formula := range installedItems {
		fmt.Printf("  %s %s\n", ui.C(ui.BrightGreen, "❯"), ui.C(ui.White, formula))
	}

	ui.Blank()
	return nil
}
