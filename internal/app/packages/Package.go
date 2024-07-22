package packages

import (
	"io/fs"
	"os"

	"github.com/hashicorp/go-version"
)

// Package struct
type Package struct {
	Name     string    `json:"name"`
	Category string    `json:"category"`
	Versions []Version `json:"versions"`
}

// GetLatestVersion from Package
func (p Package) GetLatestVersion(lb *version.Version) Version {
	var ver Version
	old, _ := version.NewVersion("0.0.0")
	for _, v := range p.Versions {
		if p.Category != "driver" {
			req, _ := version.NewVersion(v.LiquibaseCore)
			if lb.LessThan(req) {
				continue
			}
		}
		n, _ := version.NewVersion(v.Tag)
		if old.LessThan(n) {
			old = n
			ver = v
		}
	}
	return ver
}

// GetVersion from package by version tag
func (p Package) GetVersion(v string) Version {
	var r Version
	for _, ver := range p.Versions {
		if ver.Tag == v {
			r = ver
		}
	}
	return r
}

// GetInstalledVersion from classpath files
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

func (p Package) InClassPath(files []fs.FileInfo) bool {
	v := p.GetInstalledVersion(files)
	return v.Tag != ""
}

// DeleteVersion from Package
func (p Package) DeleteVersion(ver Version) []Version {
	var s int
	for i, v := range p.Versions {
		if v.Tag == ver.Tag {
			s = i
		}
	}
	return append(p.Versions[:s], p.Versions[s+1:]...)
}

// Remove installed Package Version from classpath
func (p Package) Remove(cp string, v Version) error {
	err := os.Remove(cp + v.GetFilename())
	if err != nil {
		return err
	}
	return nil
}
