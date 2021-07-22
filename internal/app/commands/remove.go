package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"package-manager/internal/app"
	"package-manager/internal/app/dependencies"
	"package-manager/internal/app/errors"
)

// removeCmd represents the install command
var removeCmd = &cobra.Command{
	Use:   "remove [PACKAGE]...",
	Short: "Removes Package",
	Aliases: []string{"rm"},
	Args: cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {

		// Set global vs local classpath
		app.SetClasspath(global, globalpath, globalpathFiles)

		//TODO recheck against gloabl
		d := dependencies.Dependencies{}
		d.Read()
		// Remove Each Package
		for _, name := range args {
			p := packs.GetByName(name)
			v := p.GetInstalledVersion(app.ClasspathFiles)
			if p.Name == "" {
				errors.Exit("Package '" + name + "' not found.", 1)
			}
			if !v.InClassPath(app.ClasspathFiles) {
				errors.Exit(name + " is not installed.", 1)
			}
			err := os.Remove(app.Classpath + v.GetFilename())
			if err != nil {
				errors.Exit("Unable to delete " + v.GetFilename() + " from classpath.", 1)
			}
			fmt.Println(v.GetFilename() + " successfully uninstalled from classpath.")
			d.Remove(p.Name)
		}
		d.Write()
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
	removeCmd.Flags().BoolVarP(&global, "global", "g", false, "remove package globally")
}
