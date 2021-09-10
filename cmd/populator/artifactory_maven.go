package main

import (
	"github.com/gocolly/colly/v2"
	"github.com/hashicorp/go-version"
	"package-manager/internal/app/packages"
	"package-manager/internal/app/utils"
	"sort"
	"strings"
)

type Maven struct {}
func (mav Maven) GetVersions(m Module) []*version.Version {
	var versionsRaw []string

	// Get Versions from Root package site
	c := colly.NewCollector()
	// Find and visit all links
	c.OnHTML("a[href]", func(f *colly.HTMLElement) {
		if !strings.Contains(f.Text, "../") && !strings.Contains(f.Text, "maven-metadata.") {
			if m.excludeSuffix != "" && m.includeSuffix == "" {
				if !strings.Contains(f.Text, m.excludeSuffix) {
					versionsRaw = append(versionsRaw, strings.TrimSuffix(f.Text, "/"))
				}
			}
			if m.excludeSuffix == "" && m.includeSuffix != "" {
				if strings.Contains(f.Text, m.includeSuffix) {
					versionsRaw = append(versionsRaw, strings.TrimSuffix(f.Text,  m.includeSuffix + "/"))
				}
			}
			if m.excludeSuffix != "" && m.includeSuffix != "" {
				if strings.Contains(f.Text, m.includeSuffix) && !strings.Contains(f.Text, m.excludeSuffix) {
					versionsRaw = append(versionsRaw, strings.TrimSuffix(f.Text,  m.includeSuffix + "/"))
				}
			}
			if m.excludeSuffix == "" && m.includeSuffix == "" {
				versionsRaw = append(versionsRaw, strings.TrimSuffix(f.Text, "/"))
			}
		}
	})
	c.Visit(m.url)

	// Sort Versions
	versions := make([]*version.Version, len(versionsRaw))
	for i, raw := range versionsRaw {
		v, _ := version.NewVersion(raw)
		versions[i] = v
	}
	sort.Sort(version.Collection(versions))

	return versions
}
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
