package commands

import (
	"fmt"
	"github.com/spf13/cobra"
)

// upgradeCmd represents the update command
var upgradeCmd = &cobra.Command{
	Use:   "upgrade [PACKAGE]...",
	Short: "Upgrades Installed Packages to the Latest Versions",

	Run: func(cmd *cobra.Command, args []string) {
		//todo implement this
		fmt.Println("upgrade")
	},
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
	upgradeCmd.Flags().BoolVarP(&global, "global", "g", false, "upgrade global packages")
}
