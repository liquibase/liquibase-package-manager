package app

import (
	_ "embed"
	"encoding/json"
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

func Version() string {
	return version
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