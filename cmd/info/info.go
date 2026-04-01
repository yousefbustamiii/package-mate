package info

import (
	"fmt"
	"strings"

	"github.com/yousefbustamiii/package-mate/internal/components"
	"github.com/yousefbustamiii/package-mate/internal/installer"
	verdetect "github.com/yousefbustamiii/package-mate/internal/installer/versions"
	"github.com/yousefbustamiii/package-mate/internal/ui"
)

// Run executes the info logic for a resolved item.
func Run(item *components.InstallItem, sec *components.Section) error {
	ui.Header(item.Name)

	allVersions, _ := installer.GetAllVersions(*item)
	_, date := verdetect.InstalledVersions(*item)

	// If it's a CLI tool with detected versions, use the detailed view
	if len(allVersions) > 0 {
		ui.Row("Description", item.Desc)

		for _, v := range allVersions {
			label := ""
			val := ""
			switch v.Type {
			case installer.VersionManaged:
				label = "Managed"
				val = ui.C(ui.BrightCyan, v.Version)
			case installer.VersionManagedOlder:
				label = "Managed (older)"
				val = ui.C(ui.Dim, v.Version)
			case installer.VersionUnmanaged:
				label = "Unmanaged"
				val = ui.C(ui.Yellow, v.Path)
			}
			ui.Row(label, val)
		}

		if !date.IsZero() {
			ui.Row("Installed At", date.Format("Jan 02, 2006 at 15:04"))
		}

		if item.Binary != "" {
			ui.Row("Binary Name", item.Binary)
		}

		ui.Footer()
		return nil
	}

	det := installer.IsInstalled(*item)
	versions, _ := verdetect.InstalledVersions(*item)

	switch det.Status {
	case installer.DetectionNotFound:
		fmt.Printf("  "+ui.C(ui.Bold+ui.White, "%s")+" is not installed.\n", item.Name)
		fmt.Printf("  To install it, run " + ui.C(ui.Cyan, "mate") + " and select " + ui.C(ui.BrightCyan, item.Name) + " then type " + ui.C(ui.BrightCyan, "1") + "\n")
		ui.Blank()
	case installer.DetectionExact, installer.DetectionOutdated:
		// Falls through to the description + installed versions block below.
	case installer.DetectionDifferentBrew:
		verStr := ""
		if len(versions) > 0 {
			verStr = " (" + versions[0] + ")"
		}
		fmt.Printf("  "+ui.C(ui.Bold+ui.White, "A different Homebrew version of %s%s is present on your system.\n"), item.Name, verStr)
		fmt.Printf("  " + ui.C(ui.Dim, det.Detail) + "\n")
		ui.Blank()
		fmt.Printf("  To install Package Mate version or override the already installed one, run " + ui.C(ui.Cyan, "mate") + " and select " + ui.C(ui.BrightCyan, item.Name) + " then type " + ui.C(ui.BrightCyan, "1") + "\n")
		ui.Blank()
	case installer.DetectionBinary:
		verStr := ""
		if len(versions) > 0 {
			verStr = " (" + versions[0] + ")"
		}
		fmt.Printf("  "+ui.C(ui.Bold+ui.White, "An unmanaged version of %s%s is present on your system.\n"), item.Name, verStr)
		fmt.Printf("  " + ui.C(ui.Dim, det.Detail) + "\n")
		ui.Blank()
		fmt.Printf("  To install Package Mate version or override the already installed one, run " + ui.C(ui.Cyan, "mate") + " and select " + ui.C(ui.BrightCyan, item.Name) + " then type " + ui.C(ui.BrightCyan, "1") + "\n")
		ui.Blank()
	case installer.DetectionManualApp:
		ui.Row("Description", item.Desc)
		ui.Row("Unmanaged", ui.C(ui.Yellow, det.BinaryPath))
		ui.Blank()
		if item.Binary != "" {
			ui.Row("Binary Name", item.Binary)
		}
		ui.Hint("Installed outside of Package Mate. Run " + ui.C(ui.Cyan, "mate") + " to manage it.")
		ui.Footer()
		return nil
	case installer.DetectionTrashedApp:
		ui.ShowManualAppUninstall(item.Name)
		ui.Blank()
		ui.Hint("Note: The application was found in your Trash.")
		ui.Footer()
		return nil
	}

	ui.Row("Description", item.Desc)

	if len(versions) > 0 {
		parts := make([]string, len(versions))
		for i, v := range versions {
			parts[i] = ui.C(ui.Dim, "(") + ui.C(ui.BrightCyan, v) + ui.C(ui.Dim, ")")
		}
		ui.Row("Installed", strings.Join(parts, ui.C(ui.Dim, " & ")))

		// Show installation date
		if !date.IsZero() {
			ui.Row("Installed At", date.Format("Jan 02, 2006 at 15:04"))
		}
	}

	if item.Binary != "" {
		ui.Row("Binary", item.Binary)
	}

	ui.Footer()
	return nil
}
