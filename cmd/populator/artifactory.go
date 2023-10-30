package main

import (
		"encoding/xml"
		"github.com/hashicorp/go-version"
		"github.com/vifraa/gopom"
		"io"
		"log"
		"net/http"
		"package-manager/internal/app/packages"
		"strings"
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
	body, err := io.ReadAll(resp.Body)
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
	if pom.Dependencies != nil {
		for _, dep := range *pom.Dependencies {
			if *dep.ArtifactID == "liquibase-core" {
				if dep.Version != nil {
					if strings.Contains(*dep.Version, "${") {
						v := strings.TrimPrefix(*dep.Version, "${")
						v = strings.TrimSuffix(v, "}")
						for k, prop := range pom.Properties.Entries {
							if k == v {
								version = prop
							}
						}
					} else {
						version = *dep.Version
					}
				} 
			}
		}
	}
	if version == "" {
		for k, prop := range pom.Properties.Entries {
			if k == "liquibase.version" {
				version = prop
			}
		}
	}
	return version
}

