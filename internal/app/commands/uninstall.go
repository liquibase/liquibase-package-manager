package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"package-manager/internal/app"
)

// uninstallCmd represents the install command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall [PACKAGE]",
	Short: "Uninstall Package",

	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		p := packs.GetByName(name)
		if p.Name == "" {
			app.Exit("Package '" + name + "' not found.", 1)
		}
		if !p.InClassPath(classpathFiles) {
			app.Exit(name + " is not installed.", 1)
		}
		err := os.Remove(classpath + p.GetFilename())
		if err != nil {
			app.Exit("Unable to delete " + p.GetFilename() + " from classpath.", 1)
		}
		fmt.Println(p.GetFilename() + " successfully uninstalled from classpath.")
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}
