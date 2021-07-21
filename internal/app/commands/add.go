package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"package-manager/internal/app"
	"package-manager/internal/app/errors"
	"package-manager/internal/app/packages"
	"strings"
)

// addCmd represents the install command
var addCmd = &cobra.Command{
	Use:   "add [PACKAGE]...",
	Short: "Add Packages",
	Args: cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {

		// Set global vs local classpath
		app.SetClasspath(global, globalpath, globalpathFiles)

		for _, name := range args {
			var p packages.Package
			var v packages.Version
			if strings.Contains(name, "@") {
				p = packs.GetByName(strings.Split(name, "@")[0])
				v = p.GetVersion(strings.Split(name, "@")[1])
				if v.Tag == "" {
					errors.Exit("Version '"+strings.Split(name, "@")[1]+"' not available.", 1)
				}
			} else {
				p = packs.GetByName(name)
				v = p.GetLatestVersion()
			}
			if p.Name == "" {
				errors.Exit("Package '"+name+"' not found.", 1)
			}
			if v.InClassPath(app.ClasspathFiles) {
				errors.Exit(name+" is already installed.", 1)
			}
			if !v.PathIsHttp() {
				v.CopyToClassPath(app.Classpath)
			} else {
				v.DownloadToClassPath(app.Classpath)
			}

			fmt.Println(v.GetFilename() + " successfully installed in classpath.")
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().BoolVarP(&global, "global", "g", false, "add package globally")
}
