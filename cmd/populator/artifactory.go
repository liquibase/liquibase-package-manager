package main

import (
	"encoding/xml"
	"github.com/hashicorp/go-version"
	"github.com/vifraa/gopom"
	"io/ioutil"
	"log"
	"net/http"
	"package-manager/internal/app/packages"
)

// Artifactory main interface for module artifactory logic
type Artifactory interface {
	GetVersions(Module) []*version.Version
	GetNewVersions(Module, packages.Package) packages.Package
}

// GetPomFromURL get POM object from remote URL
func GetPomFromURL(url string) gopom.Project {
	resp, err := http.Get(url)
	if err != nil {
		print(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		print(err)
	}
	var pom gopom.Project
	err = xml.Unmarshal(body, &pom)
	if err != nil {
		log.Fatal(err)
	}
	return pom
}

// GetCoreVersionFromPom get liquibase core version string from POM object
func GetCoreVersionFromPom(pom gopom.Project) string {
	var version string
	for k, prop := range pom.Properties.Entries {
		if k == "liquibase.version" {
			version = prop
		}
	}
	return version
}

