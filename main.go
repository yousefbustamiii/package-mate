package main

import (
	"fmt"
	"os"

	"github.com/yousefbustamiii/package-mate/cmd/bg"
	"github.com/yousefbustamiii/package-mate/cmd/cleanup"
	"github.com/yousefbustamiii/package-mate/cmd/consume"
	"github.com/yousefbustamiii/package-mate/cmd/export"
	"github.com/yousefbustamiii/package-mate/cmd/info"
	"github.com/yousefbustamiii/package-mate/cmd/install"
	"github.com/yousefbustamiii/package-mate/cmd/uninstall"
	"github.com/yousefbustamiii/package-mate/internal/background"
	"github.com/yousefbustamiii/package-mate/internal/components"
	"github.com/yousefbustamiii/package-mate/internal/installer"
	"github.com/yousefbustamiii/package-mate/internal/ui"

	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{
		Use:           "mate",
		Short:         "Package Mate — macOS dev environment installer",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				if !installer.IsBrewInstalled() {
					ui.ShowHomebrewMissing()
					return nil
				}
				install.LaunchSearch()
				return nil
			}

			alias := args[0]
			item, sec, ok := components.Resolve(alias)
			if ok {
				choice := ui.PromptActionMenu(item.Name)
				switch choice {
				case 1:
					return install.Run(item, sec)
				case 2:
					return uninstall.Run(item)
				case 3:
					return info.Run(item, sec)
				default:
					return nil
				}
			}

			// Unknown argument — show help banner.
			ui.Banner()
			fmt.Println(ui.C(ui.Dim, "  Usage:"))
			fmt.Println()
			fmt.Println(ui.C(ui.Cyan, "    mate                  ") + ui.C(ui.Dim, "— open interactive search dashboard"))
			fmt.Println(ui.C(ui.Cyan, "    mate cleanup          ") + ui.C(ui.Dim, "— clean up Homebrew cache (interactive)"))
			fmt.Println(ui.C(ui.Cyan, "    mate cleanup bg       ") + ui.C(ui.Dim, "— remove all finished background jobs"))
			fmt.Println(ui.C(ui.Cyan, "    mate export           ") + ui.C(ui.Dim, "— export installed packages as base64 string"))
			fmt.Println(ui.C(ui.Cyan, "    mate consume          ") + ui.C(ui.Dim, "— sync packages from an export string"))
			fmt.Println()
			fmt.Println(ui.C(ui.Dim, "  Examples:"))
			fmt.Println()
			fmt.Println(ui.C(ui.Dim, "    mate                  ") + ui.C(ui.Dim, "— browse tools, type to search"))
			fmt.Println(ui.C(ui.Dim, "    mate cleanup          ") + ui.C(ui.Dim, "— free up disk space from old formulae"))
			fmt.Println(ui.C(ui.Dim, "    mate cleanup bg       ") + ui.C(ui.Dim, "— clear finished jobs"))
			fmt.Println(ui.C(ui.Dim, "    mate export           ") + ui.C(ui.Dim, "— generate portable package list"))
			fmt.Println(ui.C(ui.Dim, "    mate consume          ") + ui.C(ui.Dim, "— install packages from another machine"))
			fmt.Println()
			return nil
		},
	}

	// mate bg — interactive background job viewer.
	bgCmd := &cobra.Command{
		Use:   "bg",
		Short: "View and manage background downloads",
		RunE: func(cmd *cobra.Command, args []string) error {
			bg.Run()
			return nil
		},
	}

	// mate cleanup bg — removes all finished background jobs.
	cleanupCmd := &cobra.Command{
		Use:   "cleanup",
		Short: "Clean up Homebrew cache and old formula versions",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cleanup.Run()
		},
	}

	cleanupBgCmd := &cobra.Command{
		Use:   "bg",
		Short: "Remove all finished background jobs",
		RunE: func(cmd *cobra.Command, args []string) error {
			background.CleanAllFinished()
			ui.Blank()
			ui.Done("All finished background jobs have been cleared.")
			ui.Blank()
			return nil
		},
	}

	cleanupCmd.AddCommand(cleanupBgCmd)

	// mate export — export installed packages as base64 encoded JSON array.
	exportCmd := &cobra.Command{
		Use:   "export",
		Short: "Export installed Homebrew packages as a portable base64 string",
		RunE: func(cmd *cobra.Command, args []string) error {
			return export.Run()
		},
	}

	// mate consume — consume an export string and install missing tools in background.
	consumeCmd := &cobra.Command{
		Use:   "consume",
		Short: "Sync your machine using an export string from another Mac",
		RunE: func(cmd *cobra.Command, args []string) error {
			return consume.Run()
		},
	}

	// mate _bg-exec <jobID> — hidden command used internally by background jobs.
	bgExecCmd := &cobra.Command{
		Use:    "_bg-exec",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("usage: mate _bg-exec <jobID>")
			}
			return background.Execute(args[0])
		},
	}

	root.AddCommand(bgCmd)
	root.AddCommand(cleanupCmd)
	root.AddCommand(exportCmd)
	root.AddCommand(consumeCmd)
	root.AddCommand(bgExecCmd)

	if err := root.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
