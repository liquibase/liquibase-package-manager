package commands

import (
	"github.com/spf13/cobra"
)

// searchCmd represents the install command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for Packages",

	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
