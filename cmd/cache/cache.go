package cache

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/yousefbustamiii/package-mate/internal/components"
	"github.com/yousefbustamiii/package-mate/internal/ui"
	"golang.org/x/term"
)

// Run executes the interactive cache command.
func Run() error {
	ui.Header("Cache Management")

	// Show cache menu
	choice := promptCacheAction()
	if choice == 0 {
		ui.Blank()
		fmt.Println("  " + ui.C(ui.Dim, "Cache management cancelled."))
		ui.Blank()
		return nil
	}

	// Execute the action
	switch choice {
	case 1:
		return clearCache()
	case 2:
		return updateCache()
	}

	return nil
}

// promptCacheAction shows the cache management selection menu.
func promptCacheAction() int {
	fmt.Println("  " + ui.C(ui.Bold+ui.BrightCyan, "1.") + ui.C(ui.Bold+ui.White, "  Clear Cache        ") + ui.C(ui.Dim, "(Deletes local catalog and metadata)"))
	ui.Blank()
	fmt.Println("  " + ui.C(ui.Bold+ui.BrightCyan, "2.") + ui.C(ui.Bold+ui.White, "  Update Cache       ") + ui.C(ui.Dim, "(Force sync with remote catalog now)"))
	ui.Blank()

	fmt.Println("  " + ui.C(ui.Bold+ui.Dim, "Press ") + ui.C(ui.Bold+ui.White, "1 - 2") + ui.C(ui.Bold+ui.Dim, " to choose your action"))
	fmt.Println("  " + ui.C(ui.Bold+ui.Dim, "or press ") + ui.C(ui.Bold+ui.White, "Enter") + ui.C(ui.Bold+ui.Dim, " to cancel"))
	ui.Blank()

	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return 0
	}

	fmt.Print("  " + ui.C(ui.Bold+ui.White, "❯ "))
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())
		switch input {
		case "1":
			return 1
		case "2":
			return 2
		default:
			return 0
		}
	}

	return 0
}

func clearCache() error {
	ui.Blank()
	ui.Doing("Analyzing local catalog")
	deleted, err := components.ClearCache()
	ui.Stop()

	if err != nil {
		ui.Fail("Failed to manage cache: %v", err)
		return err
	}

	if !deleted {
		// Professional, non-error message for "nothing to do"
		fmt.Printf("  " + ui.C(ui.Yellow, "~") + "  " + ui.C(ui.Bold+ui.White, "No local cache was found. Your system is already clean!") + "\n\n")
		return nil
	}

	ui.Done("Cache cleared successfully!")
	ui.Blank()
	return nil
}

func updateCache() error {
	ui.Blank()
	ui.Doing("Fetching latest catalog from remote")
	err := components.ForceUpdate()
	ui.Stop()

	if err != nil {
		ui.Fail("Failed to update cache: %v", err)
		ui.Hint("Check your internet connection or try again later.")
		return nil
	}

	ui.Done("Catalog updated and synchronized!")
	ui.Blank()
	return nil
}
