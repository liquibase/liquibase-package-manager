package main

import (
	"package-manager/internal/app"
	"package-manager/internal/app/packages"
)

func main(){
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