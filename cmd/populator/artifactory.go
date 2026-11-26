package main

import (
	"encoding/xml"
	"github.com/hashicorp/go-version"
	"github.com/vifraa/gopom"
	"io"
	"log"
	"net/http"
	"package-manager/internal/app/packages"
	"regexp"
	"strings"
)

// Artifactory main interface for module artifactory logic
type Artifactory interface {
	GetVersions(Module) []*version.Version
	GetNewVersions(Module, packages.Package) packages.Package
}

// GetPomFromURL get POM object from remote URL
func GetPomFromURL(url string) *gopom.Project {
	resp, err := http.Get(url)
	if err != nil {
		print(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		print(err)
	}
	var pom *gopom.Project
	err = xml.Unmarshal(body, &pom)
	if err != nil {
		log.Fatal(err)
	}
	return pom
}

var sha1Pattern = regexp.MustCompile(`^[0-9a-f]{40}$`)

// GetSHA1FromURL get SHA1 checksum from remote .sha1 file URL.
// Returns an empty string if the request fails or the response is not a valid SHA1,
// e.g. an error page or a rate limit message from the repository.
func GetSHA1FromURL(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Unable to download %s: %v", url, err)
		return ""
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Printf("Unable to download %s: HTTP %d", url, resp.StatusCode)
		return ""
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Unable to read %s: %v", url, err)
		return ""
	}
	// .sha1 files may contain "<checksum>  <filename>", so take the first field only
	fields := strings.Fields(string(body))
	if len(fields) == 0 || !sha1Pattern.MatchString(strings.ToLower(fields[0])) {
		log.Printf("Invalid SHA1 checksum in %s", url)
		return ""
	}
	return strings.ToLower(fields[0])
}

// GetCoreVersionFromPom get liquibase core version string from POM object
func GetCoreVersionFromPom(pom *gopom.Project) string {
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
