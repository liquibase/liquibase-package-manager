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
	classpath      string
	classpathFiles []fs.FileInfo
	extensions packages.Packages
)

var rootCmd = &cobra.Command{
	Use:   "lpm",
	Short: "Liquibase Package Manager",
	Long: `Easily manage external dependencies for Database Development.
Search for, install and uninstall drivers, extensions, and utilities.`,
}

func Execute(cp string) {
	var err error
	classpath = cp
	classpathFiles, err = ioutil.ReadDir(cp)
	if err != nil {
		panic(err)
	}
	extensions = packages.LoadExtensions()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}