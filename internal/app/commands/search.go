package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
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
		for _, p := range extensions {
			if strings.Contains(p.Name, name) || name == "" {
				fmt.Println(p.Name)
			} else  {
				fmt.Println("Package '" + name + "' not found." )
				os.Exit(1)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
