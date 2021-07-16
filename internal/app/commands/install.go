package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"os"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install [PACKAGE]",
	Short: "Install Packages",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		p := extensions.GetByName(name)
		if p.Name == "" {
			fmt.Println("Package '" + name + "' not found.")
			os.Exit(1)
		}
		if p.InClassPath(classpathFiles) {
			fmt.Println(name + " is already installed.")
			os.Exit(1)
		}
		source, err := os.Open(p.Path)
		if err != nil {
			fmt.Println("Unable to open " + p.Path)
			os.Exit(1)
		}
		defer source.Close()
		destination, err := os.Create(classpath + p.GetFilename())
		if err != nil {
			fmt.Println("Unable to access classpath located at " + classpath)
			os.Exit(1)
		}
		defer destination.Close()
		_, err = io.Copy(destination, source)
		if err != nil {
			fmt.Println("Unable to install " + p.GetFilename() + " in classpath.")
			os.Exit(1)
		}
		fmt.Println(p.GetFilename() + " successfully installed in classpath.")
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
