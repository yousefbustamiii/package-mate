package main

import (
	"fmt"
	"os"

	"github.com/yousefbustamiii/package-mate/cmd/info"
	"github.com/yousefbustamiii/package-mate/cmd/install"
	"github.com/yousefbustamiii/package-mate/cmd/uninstall"
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

			// If it's not a tool and not a matched command, show the help banner
			ui.Banner()
			fmt.Println(ui.C(ui.Dim, "  Usage:"))
			fmt.Println()
			fmt.Println(ui.C(ui.Cyan, "    mate                  ") + ui.C(ui.Dim, "— list all available tools (dashboard)"))
			fmt.Println(ui.C(ui.Cyan, "    mate <tool>           ") + ui.C(ui.Dim, "— interactive tool menu (Install/Uninstall/Info)"))
			fmt.Println()
			fmt.Println(ui.C(ui.Dim, "  Examples:"))
			fmt.Println()
			fmt.Println(ui.C(ui.Dim, "    mate node"))
			fmt.Println(ui.C(ui.Dim, "    mate redis"))
			fmt.Println()
			return nil
		},
	}

	if err := root.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
