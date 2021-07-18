package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"package-manager/internal/app/errors"
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
			errors.Exit("Package '" + name + "' not found.", 1)
		}
		v := p.GetDefaultVersion()
		if v.InClassPath(classpathFiles) {
			errors.Exit(name + " is already installed.", 1)
		}
		if !v.PathIsHttp() {
			v.CopyToClassPath(classpath)
		} else {
			v.DownloadToClassPath(classpath)
		}

		fmt.Println(v.GetFilename() + " successfully installed in classpath.")
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
