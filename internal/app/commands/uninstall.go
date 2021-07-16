package commands

import (
	"github.com/spf13/cobra"
)

// uninstallCmd represents the install command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall Package",

	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}
