package main

import (
	"context"
	"github.com/gocolly/colly/v2"
	"github.com/google/go-github/v39/github"
	"github.com/hashicorp/go-version"
	"golang.org/x/oauth2"
	"os"
	"package-manager/internal/app/packages"
	"package-manager/internal/app/utils"
	"sort"
	"strings"
)

type Module struct {
	name string
	category string
	url string
	owner string
	repo string
	includeSuffix string
	excludeSuffix string
	filePrefix  string
	artifactory Artifactory
}

func (m Module) GetVersions() []*version.Version {
	return m.artifactory.GetVersions(m)
}

func (m Module) GetNewVersions(p packages.Package) packages.Package {
	return m.artifactory.GetNewVersions(m, p)
}


type Artifactory interface {
	GetVersions(Module) []*version.Version
	GetNewVersions(Module, packages.Package) packages.Package
}

type Github struct {}
func (g Github) GetVersions(m Module) []*version.Version {

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_PAT")},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)


	rr, _, _ := client.Repositories.ListReleases(context.Background(), m.owner, m.repo, &github.ListOptions{})
	versions := make([]*version.Version, len(rr))
	for i, r := range rr {
		v, _ := version.NewVersion(r.GetTagName())
		versions[i] = v
	}
	sort.Sort(version.Collection(versions))
	return versions
}
func (g Github) GetNewVersions(m Module, p packages.Package) packages.Package {

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_PAT")},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	for _, v := range g.GetVersions(m) {
		var ver packages.Version
		ver.Tag = v.Original()
		pv := p.GetVersion(ver.Tag)

		if pv.Tag != "" {
			// if remote version is already in package manifest skip it
			continue
		}

		release, _, _ := client.Repositories.GetReleaseByTag(context.Background(), m.owner, m.repo, ver.Tag)
		for _, a := range release.Assets {
			if strings.HasSuffix(a.GetName(), ".jar") {
				ver.Path = a.GetBrowserDownloadURL()
			}
			if strings.Contains(a.GetName(), "sha1") {
				ver.Algorithm = "SHA1"
				ver.CheckSum = string(utils.HTTPUtil{}.Get(a.GetBrowserDownloadURL()))[0:40] //Get first 40 character of SHA1 only
			}
		}

		// Older versions might have bad version patters ending up with a missing sha. Don't add them.
		if ver.CheckSum != "" {
			p.Versions = append(p.Versions, ver)
		}
	}
	return p
}

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
	for _, v := range mav.GetVersions(m) {
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