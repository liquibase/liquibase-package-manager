package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"package-manager/internal/app"
	"package-manager/internal/app/packages"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List Installed Packages",
	Aliases: []string{"ls"},
	Run: func(cmd *cobra.Command, args []string) {

		// Set global vs local classpath
		app.SetClasspath(global, globalpath, globalpathFiles)

		// Collect installed packages
		var installed packages.Packages
		for _, e := range packs {
			v := e.GetInstalledVersion(app.ClasspathFiles)
			if v.InClassPath(app.ClasspathFiles) {
				installed = append(installed, e)
			}
		}

		// Format output
		fmt.Println(app.Classpath)
		if len(installed) == 0 {
			os.Exit(1)
		}
		installed.Display(app.ClasspathFiles)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVarP(&global, "global", "g", false, "list global packages")
}
