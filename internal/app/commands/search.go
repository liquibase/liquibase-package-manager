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
		var out []string
		for _, p := range packs {
			if strings.Contains(p.Name, name) || name == "" {
				out = append(out, p.Name)
			}
		}
		if len(out) == 0 {
			fmt.Println("No results found.")
		}
		for _, o := range out {
			fmt.Println(o)
		}
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
