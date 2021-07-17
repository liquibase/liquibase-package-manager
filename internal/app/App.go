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
var packagesJSON []byte

var packageFile = "packages.json"

func Version() string {
	return version
}

func PackagesInClassPath(cp string) bool {
	_, err := os.Stat(cp + packageFile)
	return err == nil
}

func CopyPackagesToClassPath(cp string) {
	err := ioutil.WriteFile(cp + packageFile, packagesJSON, 0664)
	if err != nil {
		errors.Exit(err.Error(), 1)
	}
}

func LoadPackages(cp string) packages.Packages {
	var e packages.Packages
	jsonFile, err := os.Open(cp + packageFile)
	if err != nil {
		errors.Exit(err.Error(), 1)
	}
	b, err := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(b, &e)
	if err != nil {
		return nil
	}
	return e
}