package commands

import (
	"fmt"
	"github.com/hashicorp/go-version"
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

		if global {
			errors.Exit("Can not install packages from liquibase.json globally", 1)
		}

		d := dependencies.Dependencies{}
		d.Read()

		for _, dep := range d.Dependencies {
			p := packs.GetByName(dep.GetName())
			v := p.GetVersion(dep.GetVersion())
			if p.Category != "driver" {
				core, _ := version.NewVersion(v.LiquibaseCore)
				if liquibase.Version.LessThan(core) {
					errors.Exit(p.Name+"@"+v.Tag+" is not compatible with liquibase v"+liquibase.Version.String()+". Please consider updating liquibase.", 1)
				}

			}
			if v.InClassPath(app.ClasspathFiles) {
				errors.Exit(p.Name+" is already installed.", 1)
			}
			if !v.PathIsHTTP() {
				v.CopyToClassPath(app.Classpath)
			} else {
				v.DownloadToClassPath(app.Classpath)
			}
			fmt.Println(v.GetFilename() + " successfully installed in classpath.")
		}

		minVer, _ := version.NewVersion("4.6.2")
		if !liquibase.Version.GreaterThanOrEqual(minVer) {
			p := "-cp liquibase_libs/*:" + globalpath + "*:" + liquibase.Homepath + "liquibase.jar"
			fmt.Println()
			fmt.Println("---------- IMPORTANT ----------")
			fmt.Println("Add the following JAVA_OPTS to your CLI:")
			fmt.Println("export JAVA_OPTS=\"" + p + "\"")
		}
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
