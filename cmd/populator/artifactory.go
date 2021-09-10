package main

import (
	"github.com/hashicorp/go-version"
	"package-manager/internal/app/packages"
)

type Artifactory interface {
	GetVersions(Module) []*version.Version
	GetNewVersions(Module, packages.Package) packages.Package
}
