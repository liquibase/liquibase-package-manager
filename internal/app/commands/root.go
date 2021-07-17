package commands

import (
	"github.com/spf13/cobra"
	"io/fs"
	"io/ioutil"
	"package-manager/internal/app"
	"package-manager/internal/app/packages"
)

var  (
	category string
	classpath string
	classpathFiles []fs.FileInfo
	packs packages.Packages
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
		app.Exit(err.Error(), 1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	//Global params
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.PersistentFlags().StringVar(&category, "category","", "extension, driver, or utility")
	rootCmd.Version = app.Version()
}

func initConfig()  {
	packs = packages.LoadPackages()
	if category != "" {
		packs = packs.FilterByCategory(category)
	}
}