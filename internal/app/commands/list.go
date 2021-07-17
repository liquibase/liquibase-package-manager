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
		// Collect installed packages
		var installed []string
		for _, e := range packs {
			if e.InClassPath(classpathFiles) {
				installed = append(installed, e.Name)
			}
		}

		// Format output
		fmt.Println(classpath)
		var prefix string
		for i, s := range installed {
			if (i+1) == len(installed) {
				prefix = "└──"
			} else {
				prefix = "├──"
			}
			l := fmt.Sprintf("%s %s",  prefix, s)
			fmt.Println(l)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
