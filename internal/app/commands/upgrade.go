package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"package-manager/internal/app"
	"package-manager/internal/app/dependencies"
	"package-manager/internal/app/errors"
	"strconv"
)

var (
	dryRun bool
)

// upgradeCmd represents the update command
var upgradeCmd = &cobra.Command{
	Use:     "upgrade [PACKAGE]...",
	Short:   "Upgrades Installed Packages to the Latest Versions",
	Aliases: []string{"up"},

	Run: func(cmd *cobra.Command, args []string) {
		outdated := packs.GetOutdated(liquibase.Version, app.ClasspathFiles)
		if len(outdated) == 0 {
			fmt.Println("You have no outdated packages installed.")
			fmt.Println(app.Classpath)
			os.Exit(0)
		}
		var r []string
		var prefix string
		r = append(r, fmt.Sprintf("%-4s %-38s %-38s %s", "   ", "Package", "Installed", "Available"))
		for i, p := range outdated {
			ins := p.GetInstalledVersion(app.ClasspathFiles)
			latest := p.GetLatestVersion(liquibase.Version)
			if (i + 1) == len(outdated) {
				prefix = "└──"
			} else {
				prefix = "├──"
			}
			r = append(r, fmt.Sprintf("%-4s %-38s %-38s %s", prefix, p.Name, ins.Tag, latest.Tag))
		}
		fmt.Println("You have " + strconv.Itoa(len(outdated)) + " outdated package(s) installed.")
		fmt.Println(app.Classpath)
		for _, out := range r {
			fmt.Println(out)
		}
		if !dryRun {
			d := dependencies.Dependencies{}
			if !global {
				d.Read()
			}
			for _, p := range outdated {
				ins := p.GetInstalledVersion(app.ClasspathFiles)
				latest := p.GetLatestVersion(liquibase.Version)
				fmt.Println()
				fmt.Println("removing " + p.Name + "@" + ins.Tag + " from classpath")
				err := p.Remove(app.Classpath, ins)
				if err != nil {
					errors.Exit("Unable to remove "+ins.GetFilename()+" from classpath.", 1)
				}
				fmt.Println(ins.GetFilename() + " successfully uninstalled from classpath.")
				if !global {
					d.Remove(p.Name)
				}
				fmt.Println()
				fmt.Println("adding " + p.Name + "@" + latest.Tag + " to classpath")
				if !latest.PathIsHTTP() {
					latest.CopyToClassPath(app.Classpath)
				} else {
					latest.DownloadToClassPath(app.Classpath)
				}
				fmt.Println(latest.GetFilename() + " successfully installed in classpath.")
				d.Dependencies = append(d.Dependencies, dependencies.Dependency{p.Name: latest.Tag})
			}
			if !global {
				d.Write()
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
	upgradeCmd.Flags().BoolVarP(&global, "global", "g", false, "upgrade global packages")
	upgradeCmd.Flags().BoolVar(&dryRun, "dry-run", false, "output changes without applying")
}
