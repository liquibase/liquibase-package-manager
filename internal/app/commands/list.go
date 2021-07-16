package commands

import (
	"fmt"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "ls",
	Short: "List Installed Packages",

	Run: func(cmd *cobra.Command, args []string) {
		for _, e := range extensions {
			if e.InClassPath(classpathFiles) {
				fmt.Println(e.Name)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
