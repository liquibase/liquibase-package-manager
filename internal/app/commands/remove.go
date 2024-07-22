package commands

import (
	"fmt"

	"github.com/liquibase/liquibase-package-manager/internal/app"
	"github.com/liquibase/liquibase-package-manager/internal/app/dependencies"
	"github.com/liquibase/liquibase-package-manager/internal/app/errors"
	"github.com/spf13/cobra"
)

// removeCmd represents the install command
var removeCmd = &cobra.Command{
	Use:     "remove [PACKAGE]...",
	Short:   "Removes Package",
	Aliases: []string{"rm"},
	Args:    cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {

		d := dependencies.Dependencies{}
		if !global {
			d.Read()
		}

		// Remove Each Package
		for _, name := range args {
			p := packs.GetByName(name)
			v := p.GetInstalledVersion(app.ClasspathFiles)
			if p.Name == "" {
				errors.Exit("Package '"+name+"' not found.", 1)
			}
			if !v.InClassPath(app.ClasspathFiles) {
				errors.Exit(name+" is not installed.", 1)
			}
			err := p.Remove(app.Classpath, v)
			if err != nil {
				errors.Exit("Unable to remove "+v.GetFilename()+" from classpath.", 1)
			}
			fmt.Println(v.GetFilename() + " successfully uninstalled from classpath.")
			if !global {
				d.Remove(p.Name)
			}
		}
		if !global {
			d.Write()
		}
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
	removeCmd.Flags().BoolVarP(&global, "global", "g", false, "remove package globally")
}
