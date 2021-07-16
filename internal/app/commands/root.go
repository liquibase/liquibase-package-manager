package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/fs"
	"io/ioutil"
	"os"
	"package-manager/internal/app/packages"
)

var  (
	category string
	classpath string
	classpathFiles []fs.FileInfo
	extensions packages.Packages
)

var rootCmd = &cobra.Command{
	Use:   "lpm",
	Short: "Liquibase Package Manager",
	Long: `Easily manage external dependencies for Database Development.
Search for, install, and uninstall liquibase drivers, extensions, and utilities.`,
}

func Execute(cp string) {
	var err error
	classpath = cp
	classpathFiles, err = ioutil.ReadDir(cp)
	if err != nil {
		panic(err)
	}
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	//Global params
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.PersistentFlags().StringVar(&category, "category","", "extension, driver, or utility")
	rootCmd.Version = "0.0.1"
}

func initConfig()  {
	extensions = packages.LoadPackages()
	if category != "" {
		extensions = extensions.FilterByCategory(category)
	}
}