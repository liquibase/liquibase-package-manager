package packages

import (
	"github.com/hashicorp/go-version"
	"io/fs"
)

//Package struct
type Package struct {
	Name  string `json:"name"`
	Category string `json:"category"`
	Versions []Version `json:"versions"`
}

//GetLatestVersion from Package
func (p Package) GetLatestVersion() Version {
	var ver Version
	old, _ := version.NewVersion("0.0.0")
	for _, v := range p.Versions {
		new, _ := version.NewVersion(v.Tag)
		if old.LessThan(new) {
			old = new
			ver = v
		}
	}
	return ver
}

//GetVersion from package by version name
func (p Package) GetVersion(v string) Version {
	var r Version
	for _, ver := range p.Versions {
		if ver.Tag == v {
			r = ver
		}
	}
	return r
}

//GetInstalledVersion from classpath files
func (p Package) GetInstalledVersion(files []fs.FileInfo) Version {
	var r Version
	for _, f := range files {
		for _, v := range p.Versions {
			if f.Name() == v.GetFilename() {
				r = v
			}
		}
	}
	return r
}
