package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"package-manager/internal/app/errors"
)

// uninstallCmd represents the install command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall [PACKAGE]",
	Short: "Uninstall Package",

	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		p := packs.GetByName(name)
		v := p.GetInstalledVersion(classpathFiles)
		if p.Name == "" {
			errors.Exit("Package '" + name + "' not found.", 1)
		}
		if !v.InClassPath(classpathFiles) {
			errors.Exit(name + " is not installed.", 1)
		}
		err := os.Remove(classpath + v.GetFilename())
		if err != nil {
			errors.Exit("Unable to delete " + v.GetFilename() + " from classpath.", 1)
		}
		fmt.Println(v.GetFilename() + " successfully uninstalled from classpath.")
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}
