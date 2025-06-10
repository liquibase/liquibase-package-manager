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
		fmt.Printf("[DEBUG] Package Update: Starting manifest update from %s\n", path)
		
		var bytes []byte
		if strings.HasPrefix(path, "http") {
			// Update Package from Remote URL
			fmt.Printf("[DEBUG] Package Update: Fetching manifest from remote URL\n")
			bytes = utils.HTTPUtil{}.Get(path)
			fmt.Printf("[DEBUG] Package Update: Remote manifest fetched, size: %d bytes\n", len(bytes))
		} else {
			// Update Packages from Local File
			fmt.Printf("[DEBUG] Package Update: Reading manifest from local file\n")
			file, err := os.Open(path)
			if err != nil {
				fmt.Printf("[ERROR] Package Update: Failed to open local file: %v\n", err)
				errors.Exit(err.Error(), 1)
			}
			bytes, err = io.ReadAll(file)
			if err != nil {
				fmt.Printf("[ERROR] Package Update: Failed to read local file: %v\n", err)
				errors.Exit(err.Error(), 1)
			}
			fmt.Printf("[DEBUG] Package Update: Local manifest read, size: %d bytes\n", len(bytes))
		}
		
		fmt.Printf("[DEBUG] Package Update: Validating package manifest contents\n")
		//Verify bytes are valid
		p := app.LoadPackages(bytes)
		if reflect.TypeOf(p) != reflect.TypeOf(packages.Packages{}) {
			fmt.Printf("[ERROR] Package Update: Manifest validation failed - invalid package structure\n")
			errors.Exit("Unable to validate package contents.", 1)
		}
		if p.GetByName("postgres").Name == "postgres" {
			fmt.Printf("[ERROR] Package Update: Manifest validation failed - test package check failed\n")
			errors.Exit("Unable to validate package contents.", 1)
		}
		fmt.Printf("[DEBUG] Package Update: Manifest validation passed\n")
		
		fmt.Printf("[DEBUG] Package Update: Copying manifest to classpath\n")
		app.CopyPackagesToClassPath(globalpath, bytes)
		fmt.Printf("[DEBUG] Package Update: Manifest copy completed successfully\n")
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
