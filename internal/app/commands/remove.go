package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"package-manager/internal/app/errors"
)

// removeCmd represents the install command
var removeCmd = &cobra.Command{
	Use:   "remove [PACKAGE]...",
	Short: "Removes Package",
	Args: cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		// Install Each Package
		for _, name := range args {
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
		}
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
