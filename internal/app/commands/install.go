package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install [PACKAGE]",
	Short: "Install Packages",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		p := packs.GetByName(name)

		if p.Name == "" {
			fmt.Println("Package '" + name + "' not found.")
			os.Exit(1)
		}

		if p.InClassPath(classpathFiles) {
			fmt.Println(name + " is already installed.")
			os.Exit(1)
		}

		if !p.PathIsHttp() {
			p.CopyToClassPath(classpath)
		} else {
			p.DownloadToClassPath(classpath)
		}

		fmt.Println(p.GetFilename() + " successfully installed in classpath.")
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
