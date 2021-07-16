package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "lpm",
	Short: "Liquibase Package Manager",
	Long: `Easily manage external dependencies for Database Development.
Search for, install and uninstall drivers, extensions, and utilities.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}