package main

import (
	"encoding/xml"
	"github.com/hashicorp/go-version"
	"io/ioutil"
	"net/http"
	"package-manager/internal/app/packages"
	"package-manager/internal/app/utils"
	"sort"
	"strings"
)

//Maven artifactory implementation
type Maven struct {}

type metadata struct {
	GroupID    string `xml:"groupId"`
	ArtifactID  string        `xml:"artifactId"`
	Versioning  mavenVersions `xml:"versioning"`
	LastUpdated string        `xml:"lastUpdated"`
}
type mavenVersions struct {
	Release string `xml:"release"`
	Versions []mavenVersion `xml:"versions"`
}
type mavenVersion struct {
	Version string `xml:"version"`
}

//GetVersions from maven
func (mav Maven) GetVersions(m Module) []*version.Version {
	resp, err := http.Get(m.url + "/maven-metadata.xml")
	if err != nil {
		print(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		print(err)
	}
	var meta metadata
	xml.Unmarshal(body, &meta)

	var versionsRaw []string
	for _, version := range meta.Versioning.Versions {
		if m.excludeSuffix != "" && m.includeSuffix == "" {
			if !strings.Contains(version.Version, m.excludeSuffix) {
				versionsRaw = append(versionsRaw, strings.TrimSuffix(version.Version, "/"))
			}
		}
		if m.excludeSuffix == "" && m.includeSuffix != "" {
			if strings.Contains(version.Version, m.includeSuffix) {
				versionsRaw = append(versionsRaw, strings.TrimSuffix(version.Version,  m.includeSuffix + "/"))
			}
		}
		if m.excludeSuffix != "" && m.includeSuffix != "" {
			if strings.Contains(version.Version, m.includeSuffix) && !strings.Contains(version.Version, m.excludeSuffix) {
				versionsRaw = append(versionsRaw, strings.TrimSuffix(version.Version,  m.includeSuffix + "/"))
			}
		}
		if m.excludeSuffix == "" && m.includeSuffix == "" {
			versionsRaw = append(versionsRaw, strings.TrimSuffix(version.Version, "/"))
		}
	}

	// Sort Versions
	versions := make([]*version.Version, len(versionsRaw))
	for i, raw := range versionsRaw {
		v, _ := version.NewVersion(raw)
		versions[i] = v
	}
	sort.Sort(version.Collection(versions))
	return versions
}

//GetNewVersions from maven
func (mav Maven) GetNewVersions(m Module, p packages.Package) packages.Package {
	//Look for new versions
	for _, v := range m.GetVersions() {
		var ver packages.Version
		ver.Tag = v.Original()
		pv := p.GetVersion(ver.Tag)
		if pv.Tag != "" {
			// if remote version is already in package manifest skip it
			continue
		}

		var tag string
		if m.includeSuffix != "" {
			tag = ver.Tag + m.includeSuffix
		} else {
			tag = ver.Tag
		}

		if m.filePrefix != "" {
			ver.Path = m.url + "/" + tag + "/" + m.filePrefix + tag + ".jar"
		} else {
			ver.Path = m.url + "/" + tag + "/" + p.Name + "-" + tag + ".jar"
		}
		ver.Algorithm = "SHA1"
		sha := string(utils.HTTPUtil{}.Get(ver.Path + ".sha1"))
		if strings.Contains(sha, "html") {
			sha = ""
		}
		ver.CheckSum = sha[0:40] //Get first 40 character of SHA1 only

		// Older versions might have bad version patters ending up with a missing sha. Don't add them.
		if ver.CheckSum != "" {
			p.Versions = append(p.Versions, ver)
		}
	}
	return p
}
