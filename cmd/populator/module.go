package main

import (
	"github.com/hashicorp/go-version"
	"package-manager/internal/app/packages"
)

//Category module category type
type Category string

const (
	//Extension category
	Extension Category = "extension"
	//Driver category
	Driver Category = "driver"
)

//Module struct
type Module struct {
	name string
	category Category
	url string
	includeSuffix string
	excludeSuffix string
	filePrefix  string
	owner string
	repo string
	artifactory Artifactory
}

//GetVersions for module
func (m Module) GetVersions() []*version.Version {
	return m.artifactory.GetVersions(m)
}

//GetNewVersions for module
func (m Module) GetNewVersions(p packages.Package) packages.Package {
	return m.artifactory.GetNewVersions(m, p)
}
