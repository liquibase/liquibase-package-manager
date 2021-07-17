package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"package-manager/internal/app/packages"
	"strings"
)

// searchCmd represents the install command
var searchCmd = &cobra.Command{
	Use:   "search [PACKAGE]",
	Short: "Search for Packages",

	Run: func(cmd *cobra.Command, args []string) {
		var name string
		if len(args) > 0 {
			name = args[0]
			if len(name) < 3 {
				fmt.Println("Minimum of 3 characters required for search.")
				os.Exit(1)
			}
		} else {
			name = ""
		}
		var found packages.Packages
		for _, p := range packs {
			if strings.Contains(p.Name, name) || name == "" {
				found = append(found, p)
			}
		}
		if len(found) == 0 {
			fmt.Println("No results found.")
			os.Exit(1)
		}
		found.Display()
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
