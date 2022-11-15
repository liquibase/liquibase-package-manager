package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"os"
	"package-manager/internal/app"
	"package-manager/internal/app/errors"
	"package-manager/internal/app/packages"
	"package-manager/internal/app/utils"
	"reflect"
	"strings"
)

var (
	path string
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates the Package Manifest",

	Run: func(cmd *cobra.Command, args []string) {
		var bytes []byte
		if strings.HasPrefix(path, "http") {
			// Update Package from Remote URL
			bytes = utils.HTTPUtil{}.Get(path)
		} else {
			// Update Packages from Local File
			file, err := os.Open(path)
			if err != nil {
				errors.Exit(err.Error(), 1)
			}
			bytes, err = io.ReadAll(file)
			if err != nil {
				errors.Exit(err.Error(), 1)
			}
		}
		//Verify bytes are valid
		p := app.LoadPackages(bytes)
		if reflect.TypeOf(p) != reflect.TypeOf(packages.Packages{}) {
			errors.Exit("Unable to validate package contents.", 1)
		}
		if p.GetByName("postgres").Name == "postgres" {
			errors.Exit("Unable to validate package contents.", 1)
		}
		app.CopyPackagesToClassPath(globalpath, bytes)
		fmt.Println("Package manifest updated from " + path)
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().StringVarP(
		&path,
		"path",
		"p",
		"https://raw.githubusercontent.com/liquibase/liquibase-package-manager/master/internal/app/packages.json",
		"path to new packages.json manifest",
	)
}
