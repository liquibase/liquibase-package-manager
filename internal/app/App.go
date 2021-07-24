package app

import (
	_ "embed" // Embed Import for Package Files
	"encoding/json"
	"io/fs"
	"io/ioutil"
	"os"
	"package-manager/internal/app/errors"
	"package-manager/internal/app/packages"
)

//go:embed "VERSION"
var version string

//PackagesJSON is embedded for first time run
//go:embed "packages.json"
var PackagesJSON []byte

//PackageFile exported for overwrite
var PackageFile = "packages.json"
//Classpath exported for overwrite
var Classpath string
//ClasspathFiles exported for overwrite
var ClasspathFiles []fs.FileInfo

//Version output from embedded file
func Version() string {
	return version
}

//SetClasspath to switch between global and local modules
func SetClasspath(global bool, globalpath string, globalpathFiles []fs.FileInfo) {
	if global {
		Classpath = globalpath
		ClasspathFiles = globalpathFiles
	} else {
		pwd, err := os.Getwd()
		if err != nil {
			errors.Exit(err.Error(), 1)
		}
		Classpath = pwd + "/liquibase_modules/"
		os.Mkdir(Classpath, 0775)
		if err != nil {
			errors.Exit(err.Error(), 1)
		}
		ClasspathFiles, err = ioutil.ReadDir(Classpath)
		if err != nil {
			errors.Exit(err.Error(), 1)
		}
	}
}

//PackagesInClassPath is the packages.json file in global classpath
func PackagesInClassPath(cp string) bool {
	_, err := os.Stat(cp + PackageFile)
	return err == nil
}

//CopyPackagesToClassPath install packages.json to global classpath
func CopyPackagesToClassPath(cp string, p []byte) {
	err := ioutil.WriteFile(cp + PackageFile, p, 0664)
	if err != nil {
		errors.Exit(err.Error(), 1)
	}
}

//LoadPackages get packages from bytes from file
func LoadPackages(b []byte) packages.Packages {
	var e packages.Packages
	err := json.Unmarshal(b, &e)
	if err != nil {
		errors.Exit(err.Error(), 1)
	}
	return e
}