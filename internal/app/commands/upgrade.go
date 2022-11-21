package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"package-manager/internal/app"
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
	},
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
	upgradeCmd.Flags().BoolVarP(&global, "global", "g", false, "upgrade global packages")
	upgradeCmd.Flags().BoolVar(&dryRun, "dry-run", false, "output changes without applying")
}
