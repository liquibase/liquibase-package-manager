package app

import (
	_ "embed"
	"encoding/json"
	"io/fs"
	"io/ioutil"
	"os"
	"package-manager/internal/app/errors"
	"package-manager/internal/app/packages"
)

//go:embed "VERSION"
var version string

//go:embed "packages.json"
var PackagesJSON []byte

var PackageFile = "packages.json"
var Classpath string
var ClasspathFiles []fs.FileInfo

func Version() string {
	return version
}

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

func PackagesInClassPath(cp string) bool {
	_, err := os.Stat(cp + PackageFile)
	return err == nil
}

func CopyPackagesToClassPath(cp string, p []byte) {
	err := ioutil.WriteFile(cp + PackageFile, p, 0664)
	if err != nil {
		errors.Exit(err.Error(), 1)
	}
}

func LoadPackages(b []byte) packages.Packages {
	var e packages.Packages
	err := json.Unmarshal(b, &e)
	if err != nil {
		errors.Exit(err.Error(), 1)
	}
	return e
}