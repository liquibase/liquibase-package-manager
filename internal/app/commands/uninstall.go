package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

// uninstallCmd represents the install command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall [PACKAGE]",
	Short: "Uninstall Package",

	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		p := packs.GetByName(name)
		if p.Name == "" {
			fmt.Println("Package '" + name + "' not found.")
			os.Exit(1)
		}
		if !p.InClassPath(classpathFiles) {
			fmt.Println(name + " is not installed.")
			os.Exit(1)
		}
		err := os.Remove(classpath + p.GetFilename())
		if err != nil {
			fmt.Println("Unable to delete " + p.GetFilename() + " from classpath.")
			os.Exit(1)
		}
		fmt.Println(p.GetFilename() + " successfully uninstalled from classpath.")
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}
