package app

import (
	_ "embed" // Embed Import for Package Files
	"encoding/json"
	"io/fs"
	"os"

	"github.com/liquibase/liquibase-package-manager/internal/app/errors"
	"github.com/liquibase/liquibase-package-manager/internal/app/packages"
	"github.com/liquibase/liquibase-package-manager/internal/app/utils"
)

var version string

// PackagesJSON is embedded for first time run
//
//go:embed "packages.json"
var PackagesJSON []byte

// PackageFile exported for overwrite
var PackageFile = "packages.json"

// Classpath exported for overwrite
var Classpath string

// ClasspathFiles exported for overwrite
var ClasspathFiles []fs.FileInfo

// Version output from build metadata
func Version() string {
	return version
}

// SetClasspath to switch between global and local modules
func SetClasspath(global bool, globalpath string, globalpathFiles []fs.FileInfo) {
	if global {
		Classpath = globalpath
		ClasspathFiles = globalpathFiles
	} else {
		pwd, err := os.Getwd()
		if err != nil {
			errors.Exit(err.Error(), 1)
		}
		Classpath = pwd + "/liquibase_libs/"
		ClasspathFiles, _ = utils.ReadDir(Classpath)
	}
}

// PackagesInClassPath is the packages.json file in global classpath
func PackagesInClassPath(cp string) bool {
	_, err := os.Stat(cp + PackageFile)
	return err == nil
}

// CopyPackagesToClassPath install packages.json to global classpath
func CopyPackagesToClassPath(cp string, p []byte) {
	err := os.WriteFile(cp+PackageFile, p, 0664)
	if err != nil {
		errors.Exit(err.Error(), 1)
	}
}

// LoadPackages get packages from bytes from file
func LoadPackages(b []byte) packages.Packages {
	var e packages.Packages
	err := json.Unmarshal(b, &e)
	if err != nil {
		errors.Exit(err.Error(), 1)
	}
	return e
}

// WritePackages write packages back to file
func WritePackages(p packages.Packages) {
	b, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		errors.Exit(err.Error(), 1)
	}
	pwd, err := os.Getwd()
	if err != nil {
		errors.Exit(err.Error(), 1)
	}
	err = os.WriteFile(pwd+"/internal/app/packages.json", b, 0664)
	if err != nil {
		errors.Exit(err.Error(), 1)
	}
}
