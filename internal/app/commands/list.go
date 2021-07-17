package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"package-manager/internal/app/packages"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "ls",
	Short: "List Installed Packages",

	Run: func(cmd *cobra.Command, args []string) {
		// Collect installed packages
		var installed packages.Packages
		for _, e := range packs {
			if e.InClassPath(classpathFiles) {
				installed = append(installed, e)
			}
		}

		// Format output
		fmt.Println(classpath)
		if len(installed) == 0 {
			os.Exit(1)
		}
		installed.Display()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
