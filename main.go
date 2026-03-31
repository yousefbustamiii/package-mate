package main

import (
	"fmt"
	"os"

	"github.com/yousefbustamiii/package-mate/cmd/bg"
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
			fmt.Println(ui.C(ui.Cyan, "    mate                  ") + ui.C(ui.Dim, "— list all available tools (dashboard)"))
			fmt.Println(ui.C(ui.Cyan, "    mate <tool>           ") + ui.C(ui.Dim, "— interactive tool menu (Install/Uninstall/Info)"))
			fmt.Println(ui.C(ui.Cyan, "    mate bg               ") + ui.C(ui.Dim, "— view and manage background downloads"))
			fmt.Println(ui.C(ui.Cyan, "    mate cleanup bg       ") + ui.C(ui.Dim, "— remove all finished background jobs"))
			fmt.Println()
			fmt.Println(ui.C(ui.Dim, "  Examples:"))
			fmt.Println()
			fmt.Println(ui.C(ui.Dim, "    mate node"))
			fmt.Println(ui.C(ui.Dim, "    mate redis"))
			fmt.Println(ui.C(ui.Dim, "    mate bg"))
			fmt.Println(ui.C(ui.Dim, "    mate cleanup bg"))
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
		Short: "Clean up completed jobs and temporary files",
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
	root.AddCommand(bgExecCmd)

	if err := root.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
