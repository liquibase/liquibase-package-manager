package main

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/hashicorp/go-version"
	"package-manager/pkg/lpm"
	"sort"
	"strings"
)

type Module struct {
	name          string
	category      string
	url           string
	includeSuffix string
	excludeSuffix string
	filePrefix    string
}

func (m Module) GetJarPath(tag string) string {
	return fmt.Sprintf("%s/%s/%s%s.jar",
		m.url,
		tag,
		m.filePrefix,
		tag)
}

func (m Module) GetSha1Path(tag string) string {
	return fmt.Sprintf("%s.sha1", m.GetJarPath(tag))
}

func (m Module) getCheckSum(tag string, algo lpm.ChecksumAlgorithm) (cs string, err error) {
	var url string
	var _sha []byte
	var sha string

	switch algo {
	case lpm.Sha1Algorithm:
		url = m.GetSha1Path(tag)
	default:
		err = fmt.Errorf("invalid checksum algorithm '%s'", string(algo))
		goto end
	}
	_sha, err = lpm.HttpGet(url)
	if err != nil {
		err = fmt.Errorf("unable to retrieve %s: %w", url, err)
		goto end
	}
	sha = string(_sha)
	if strings.Contains(sha, "html") {
		sha = ""
	}
	if 40 > len(sha) {
		goto end
	}
	cs = sha[0:40] //Get first 40 character of SHA1 only

end:
	return cs, err
}

func (m Module) onAHref(f *colly.HTMLElement, vs []string) []string {
	var text string

	if strings.Contains(f.Text, "../") {
		goto end
	}

	if strings.Contains(f.Text, "maven-metadata.") {
		goto end
	}

	switch true {

	case m.HasNoSuffixes():
		text = strings.TrimSuffix(f.Text, "/")

	case m.HasExcludeSuffix():
		if m.ContainsExcludeSuffix(f.Text) {
			goto end
		}
		text = strings.TrimSuffix(f.Text, "/")

	case m.HasIncludeSuffix():
		if !m.ContainsIncludeSuffix(f.Text) {
			goto end
		}
		text = strings.TrimSuffix(f.Text, m.includeSuffix+"/")

	case m.HasBothSuffixes():
		if m.ContainsExcludeSuffix(f.Text) {
			goto end
		}
		if !m.ContainsIncludeSuffix(f.Text) {
			goto end
		}
		text = strings.TrimSuffix(f.Text, m.includeSuffix+"/")

	}

	vs = append(vs, text)
end:
	return vs
}

func (m Module) ContainsIncludeSuffix(s string) bool {
	// @TODO Using Contains below seems like they could matching
	//       false positives. Regex would probably better, but I
	//       do not know the exact format to look for.
	return strings.Contains(s, m.includeSuffix)
}
func (m Module) ContainsExcludeSuffix(s string) bool {
	// @TODO Using Contains below seems like they could matching
	//       false positives. Regex would probably better, but I
	//       do not know the exact format to look for.
	return strings.Contains(s, m.excludeSuffix)
}
func (m Module) HasNoSuffixes() bool {
	return m.excludeSuffix == "" && m.includeSuffix != ""
}
func (m Module) HasExcludeSuffix() bool {
	return m.excludeSuffix != "" && m.includeSuffix == ""
}
func (m Module) HasIncludeSuffix() bool {
	return m.includeSuffix != "" && m.excludeSuffix == ""
}
func (m Module) HasBothSuffixes() bool {
	return m.excludeSuffix != "" && m.includeSuffix != ""
}

func (m Module) retrieveVersionsViaHttp() (versions []*version.Version, err error) {
	var versionsRaw []string

	// Get Versions from Root package site
	c := colly.NewCollector()
	// Find and visit all links
	c.OnHTML("a[href]", func(f *colly.HTMLElement) {
		versionsRaw = m.onAHref(f, versionsRaw)
	})

	err = c.Visit(m.url)
	if err != nil {
		err = fmt.Errorf("unable to visit '%s': %w", m.url, err)
		goto end
	}

	// Make Sorted Versions
	versions = make([]*version.Version, len(versionsRaw))
	for i, raw := range versionsRaw {
		v, _ := version.NewVersion(raw)
		versions[i] = v
	}
	sort.Sort(version.Collection(versions))

end:
	return versions, err

}

func (m Module) getNewVersion(p lpm.Package, tag string) (ver lpm.Version, err error) {

	ver.Tag = tag

	if m.includeSuffix != "" {
		tag += m.includeSuffix
	}

	if m.filePrefix == "" {
		m.filePrefix = p.Name
	}

	ver.Path = m.GetJarPath(tag)

	ver.Algorithm = lpm.Sha1Algorithm
	ver.CheckSum, err = m.getCheckSum(tag, ver.Algorithm)

	return ver, err
}

func (m Module) getNewVersions(p lpm.Package) (lpm.Package, error) {
	var vs []*version.Version
	var ver lpm.Version
	var err error

	// Retrieve the version by URL
	vs, err = m.retrieveVersionsViaHttp()
	if err != nil {
		err = fmt.Errorf("unable to retrieve vs for package '%s': %w",
			p.Name,
			err)
		goto end
	}

	//Look for new versions
	for _, v := range vs {
		tag := v.Original()

		pv := p.GetVersion(tag)

		if pv.Tag != "" {
			// if remote version is already in package manifest skip it
			continue
		}

		ver, err = m.getNewVersion(p, v.Original())

		if err != nil {
			// Older vs might have bad version patterns
			// ending up with a missing sha. Don't add them.
			continue
		}

		p.Versions = append(p.Versions, ver)

	}
end:
	return p, err
}
