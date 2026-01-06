package commands

import (
	"github.com/spf13/cobra"
	"io"
	"io/fs"
	"os"
	"package-manager/internal/app"
	"package-manager/internal/app/errors"
	"package-manager/internal/app/packages"
	"package-manager/internal/app/utils"
)

var (
	category        string
	liquibase       utils.Liquibase
	globalpath      string
	globalpathFiles []fs.FileInfo
	packs           packages.Packages
	global          bool
	dryRun          bool
	skipExisting    bool
)

var rootCmd = &cobra.Command{
	Use:   "lpm",
	Short: "Liquibase Package Manager",
	Long: `Easily manage external dependencies for Database Development.
Search for, install, and uninstall liquibase drivers, extensions, and utilities.`,
}

// Execute main entry point for CLI
func Execute(cp string, s string) {
	var err error
	liquibase = utils.LoadLiquibase(cp)
	globalpath = liquibase.Homepath + "lib" + s
	globalpathFiles, err = utils.ReadDir(globalpath)
	if err != nil {
		errors.Exit(err.Error(), 1)
	}
	if err := rootCmd.Execute(); err != nil {
		errors.Exit(err.Error(), 1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	//Global params
	//rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.PersistentFlags().StringVar(&category, "category", "", "extension, driver, or utility")
	rootCmd.Version = app.Version()
	rootCmd.SetVersionTemplate("{{with .Name}}{{printf \"%s \" .}}{{end}}{{with .Short}}{{printf \"(%s) \" .}}{{end}}{{printf \"version %s\" .Version}}\n")
}

func initConfig() {
	//Install Embedded Package File
	if !app.PackagesInClassPath(globalpath) {
		app.CopyPackagesToClassPath(globalpath, app.PackagesJSON)
	}

	//Get Bytes from Package File
	jsonFile, err := os.Open(globalpath + app.PackageFile)
	if err != nil {
		errors.Exit(err.Error(), 1)
	}
	b, _ := io.ReadAll(jsonFile)

	//Load Bytes to Packages
	packs = app.LoadPackages(b)
	if category != "" {
		packs = packs.FilterByCategory(category)
	}

	// Set global vs local classpath
	app.SetClasspath(global, globalpath, globalpathFiles)
}
