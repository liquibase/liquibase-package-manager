package main

import (
	"github.com/gocolly/colly/v2"
	"github.com/hashicorp/go-version"
	"package-manager/internal/app"
	"package-manager/internal/app/packages"
	"package-manager/internal/app/utils"
	"sort"
	"strings"
)

type module struct {
	name string
	category string
	url string
}
type modules []module
var mods modules

func init() {
	mods = []module{
		{"liquibase-cache", "module", "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-cache"},
		{"liquibase-cassandra", "module", "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-cassandra"},
		{"liquibase-cosmosdb", "module", "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-cosmosdb"},
		{"liquibase-db2i", "module", "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-db2i"},
		{"liquibase-filechangelog", "module", "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-filechangelog"},
		{"liquibase-hanadb", "module", "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-hanadb"},
		{"liquibase-hibernate5", "module", "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-hibernate5"},
		{"liquibase-maxdb", "module", "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-maxdb"},
		{"liquibase-modify-column", "module", "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-modify"},
		{"liquibase-mongodb", "module", "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-mongodb"},
		{"liquibase-mssql", "module", "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-mssql"},
		{"liquibase-neo4j", "module", "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-neo4j"},
		{"liquibase-oracle", "module", "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-oracle"},
		{"liquibase-percona", "module", "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-percona"},
		{"liquibase-postgresql", "module", "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-postgresql"},
		{"liquibase-redshift", "module", "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-redshift"},
		{"liquibase-snowflake", "module", "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-snowflake"},
		{"liquibase-sqlfire", "module", "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-sqlfire"},
		{"liquibase-teradata", "module", "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-teradata"},
		{"liquibase-verticaDatabase", "module", "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-verticaDatabase"},
		//"liquibase-compat",
		//"liquibase-javalogger",
		//"liquibase-nochangeloglock",
		//"liquibase-nochangelogupdate",
		//"liquibase-sequencetable",
		//"liquibase-vertica",
	}
}

func (es modules) getByName(n string) module {
	var r module
	for _, e := range es {
		if e.name == n {
			r = e
		}
	}
	return r
}

func getNewVersions(m module, p packages.Package) packages.Package {
	var versionsRaw []string

	// Get Versions from Root package site
	c := colly.NewCollector()
	// Find and visit all links
	c.OnHTML("a[href]", func(f *colly.HTMLElement) {
		if !strings.Contains(f.Text, "../") && !strings.Contains(f.Text, "maven-metadata.") {
			versionsRaw = append(versionsRaw, strings.TrimRight(f.Text, "/"))
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

	//Look for new versions
	for _, v := range versions {
		pv := p.GetVersion(v.String())
		if pv.Tag != "" {
			// if remove version is already in package manifest skip it
			continue
		}
		var ver packages.Version
		ver.Tag = v.String()
		ver.Path = m.url + "/" + v.String() + "/" + p.Name + "-" + v.String() + ".jar"
		ver.Algorithm = "SHA1"
		sha := string(utils.HTTPUtil{}.Get(ver.Path + ".sha1"))
		if strings.Contains(sha, "html") {
			sha = ""
		}
		ver.CheckSum = sha
		p.Versions = append(p.Versions, ver)
	}
	return p
}

func main(){

	var newPacks = packages.Packages{}

	// read packages from embedded file
	packs := app.LoadPackages(app.PackagesJSON)
	for _, p := range packs {
		m := mods.getByName(p.Name)
		if m.name != "" {
			// Get new versions for a package
			newPacks = append(newPacks, getNewVersions(m, p))
		}
	}

	//Write all packages back to manifest.
	app.WritePackages(newPacks)
}