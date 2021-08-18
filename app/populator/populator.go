package main

import (
	"fmt"
	"os"
	"package-manager/pkg/lpm"
)

func main() {

	var newPacks = make(lpm.Packages, 0)
	var pack lpm.Package
	var err error

	// @TODO Verify '/' is the right thing here...
	cfg := lpm.NewContext("/")
	err = cfg.Initialize()

	// read packages from embedded file
	packs, err := cfg.UnmarshalPackages(lpm.PackagesJSON)
	if err != nil {
		goto end
	}
	for _, p := range packs {
		m := modules.getByName(p.Name)
		if m.name == "" {
			newPacks = append(newPacks, p)
			continue
		}
		pack, err = m.getNewVersions(p)
		if err != nil {
			err = fmt.Errorf("failed to get versions for %s: %w",
				m.name,
				err)
			goto end
		}
		// Get new versions for a package
		newPacks = append(newPacks, pack)
	}

	//WriteManifest all packages back to manifest.
	err = cfg.WritePackages(newPacks)

end:
	if err != nil {
		fmt.Printf("ERROR %s\n", err.Error())
		os.Exit(1)
	}

}
