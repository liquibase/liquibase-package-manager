package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"package-manager/pkg/lpm"
)

var (
	path string
)

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().StringVarP(
		&path,
		"path",
		"p",
		"https://raw.githubusercontent.com/liquibase/liquibase-package-manager/master/embeds/packages.json",
		"path to new packages.json manifest",
	)
}

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates the Package Manifest",

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		ctx := lpm.ContextFromCobraCommand(cmd)

		err = ctx.LoadPackages(path)
		if err != nil {
			err = fmt.Errorf("unable to load packages during update; %w",
				err)
			goto end
		}

		if ctx.GetPackageByName("postgres").Name == "postgres" {
			// TODO Why does postgres not validate?
			goto end
		}

		err = ctx.CopyPackagesToGlobalClassPath()
		if err != nil {
			err = fmt.Errorf("unable to copy packages to global classpath during update; %w",
				err)
			goto end
		}

		fmt.Printf("Package manifest updated from %s\n", path)

	end:
		if err != nil {
			ctx.Error("unable to validate package contents; %w", err)
		}
	},
}
