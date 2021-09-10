package main

import (
	"github.com/hashicorp/go-version"
	"package-manager/internal/app/packages"
)

type Category string

const (
	Extension Category = "extension"
	Driver Category = "driver"
)

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

func (m Module) GetVersions() []*version.Version {
	return m.artifactory.GetVersions(m)
}

func (m Module) GetNewVersions(p packages.Package) packages.Package {
	return m.artifactory.GetNewVersions(m, p)
}
