package commands

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/liquibase/liquibase-package-manager/internal/app"
	"github.com/liquibase/liquibase-package-manager/internal/app/dependencies"
	"github.com/liquibase/liquibase-package-manager/internal/app/errors"
	"github.com/liquibase/liquibase-package-manager/internal/app/packages"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add [PACKAGE]...",
	Short: "Add Packages",
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {

		d := dependencies.Dependencies{}
		if !global {
			d.Read()
		}

		for _, name := range args {
			var p packages.Package
			var v packages.Version
			if strings.Contains(name, "@") {
				p = packs.GetByName(strings.Split(name, "@")[0])
				if p.Name == "" {
					errors.Exit("Package '"+name+"' not found.", 1)
				}
				v = p.GetVersion(strings.Split(name, "@")[1])
				if v.Tag == "" {
					errors.Exit("Version '"+strings.Split(name, "@")[1]+"' not available.", 1)
				}
				if p.Category != "driver" {
					core, _ := version.NewVersion(v.LiquibaseCore)
					if liquibase.Version.LessThan(core) {
						errors.Exit(name+" is not compatible with liquibase v"+liquibase.Version.String()+". Please consider updating liquibase.", 1)
					}
				}
			} else {
				p = packs.GetByName(name)
				if p.Name == "" {
					errors.Exit("Package '"+name+"' not found.", 1)
				}
				v = p.GetLatestVersion(liquibase.Version)
				if v.Tag == "" {
					errors.Exit("Unable to find compatible version of "+name+" for liquibase v"+liquibase.Version.String()+". Please consider updating liquibase.", 1)
				}
			}
			if p.InClassPath(app.ClasspathFiles) {
				v := p.GetInstalledVersion(app.ClasspathFiles)
				fmt.Println(p.Name + "@" + v.Tag + " is already installed.")
				fmt.Println(name + " can not be installed.")
				errors.Exit("Consider running `lpm upgrade`.", 1)
			}
			if !v.PathIsHTTP() {
				v.CopyToClassPath(app.Classpath)
			} else {
				v.DownloadToClassPath(app.Classpath)
			}
			fmt.Println(v.GetFilename() + " successfully installed in classpath.")
			d.Dependencies = append(d.Dependencies, dependencies.Dependency{p.Name: v.Tag})
		}

		if !global {
			//Add package to local manifest
			if !d.FileExists() {
				d.CreateFile()
			}
			d.Write()

			minVer, _ := version.NewVersion("4.6.2")
			if !liquibase.Version.GreaterThanOrEqual(minVer) {
				p := "-cp liquibase_libs/*:" + globalpath + "*:" + liquibase.Homepath + "liquibase.jar"
				fmt.Println()
				fmt.Println("---------- IMPORTANT ----------")
				fmt.Println("Add the following JAVA_OPTS to your CLI:")
				fmt.Println("export JAVA_OPTS=\"" + p + "\"")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().BoolVarP(&global, "global", "g", false, "add package globally")
}
