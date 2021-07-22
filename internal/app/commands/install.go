package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"package-manager/internal/app"
	"package-manager/internal/app/dependencies"
	"package-manager/internal/app/errors"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install Packages from liquibase.json",
	Run: func(cmd *cobra.Command, args []string) {
		// Set global vs local classpath
		app.SetClasspath(global, globalpath, globalpathFiles)
		d := dependencies.Dependencies{}
		d.Read()

		for _, dep := range d.Dependencies {
			p := packs.GetByName(dep.GetName())
			v := p.GetVersion(dep.GetVersion())

			if v.InClassPath(app.ClasspathFiles) {
				errors.Exit(p.Name+" is already installed.", 1)
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
	rootCmd.AddCommand(installCmd)
}