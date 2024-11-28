package commands

import (
	"fmt"
	"os"

	"github.com/liquibase/liquibase-package-manager/internal/app"
	"github.com/liquibase/liquibase-package-manager/internal/app/packages"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List Installed Packages",
	Aliases: []string{"ls"},
	Run: func(cmd *cobra.Command, args []string) {

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
			os.Exit(0)
		}
		for _, out := range installed.Display(app.ClasspathFiles) {
			fmt.Println(out)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVarP(&global, "global", "g", false, "list global packages")
}
