package main

import (
	"fmt"
	"os"
	"package-manager/internal/app"
	"package-manager/internal/app/packages"
)

func checkConfig() {
	_, b := os.LookupEnv("GITHUB_PAT")
	if !b {
		fmt.Println("Unable to locate GITHUB_PAT env. Check your configuration and try again.")
		os.Exit(1)
	}
}

func main(){
	checkConfig()
	var newPacks = packages.Packages{}

	// read packages from embedded file
	packs := app.LoadPackages(app.PackagesJSON)
	for _, p := range packs {
		m := modules.getByName(p.Name)
		if m.name != "" {
			// Get new versions for a package
			newPacks = append(newPacks, m.GetNewVersions(p))
		} else {
			newPacks = append(newPacks, p)
		}
	}

	//Write all packages back to manifest.
	app.WritePackages(newPacks)
}