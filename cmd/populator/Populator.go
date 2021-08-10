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
	includeSuffix string
	excludeSuffix string
}
type modules []module
var mods modules

func init() {
	mods = []module{
		//{"postgresql", "driver", "https://repo1.maven.org/maven2/org/postgresql/postgresql", "", ".jre"},
		{"mssql", "driver", "https://repo1.maven.org/maven2/com/microsoft/sqlserver/mssql-jdbc",".jre11",".jre11-preview"},
		//{"mariadb", "driver", "https://repo1.maven.org/maven2/org/mariadb/jdbc/mariadb-java-client","",""},
		//{"h2", "driver", "https://repo1.maven.org/maven2/com/h2database/h2","",""},
		//{"db2", "driver", "https://repo1.maven.org/maven2/com/ibm/db2/jcc","","db2"},
		//{"snowflake", "driver", "https://repo1.maven.org/maven2/net/snowflake/snowflake-jdbc","",""},
		//{"sybase", "driver", "https://repo1.maven.org/maven2/net/sf/squirrel-sql/plugins/sybase","",""},
		//{"firebird", "driver", "https://repo1.maven.org/maven2/net/sf/squirrel-sql/plugins/firebird","",""},
		//{"sqlite", "driver", "https://repo1.maven.org/maven2/org/xerial/sqlite-jdbc","",""},
		//{"oracle", "driver", "https://repo1.maven.org/maven2/com/oracle/ojdbc/ojdbc8","",""},
		//{"mysql", "driver", "https://repo1.maven.org/maven2/mysql/mysql-connector-java","",""},
		//{"liquibase-cache", "extension", "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-cache","",""},
		//{"liquibase-cassandra", "extension", "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-cassandra","",""},
		//{"liquibase-cosmosdb", "extension", "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-cosmosdb","",""},
		//{"liquibase-db2i", "extension", "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-db2i","",""},
		//{"liquibase-filechangelog", "extension", "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-filechangelog","",""},
		//{"liquibase-hanadb", "extension", "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-hanadb","",""},
		//{"liquibase-hibernate5", "extension", "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-hibernate5","",""},
		//{"liquibase-maxdb", "extension", "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-maxdb","",""},
		//{"liquibase-modify-column", "extension", "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-modify","",""},
		//{"liquibase-mongodb", "extension", "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-mongodb","",""},
		//{"liquibase-mssql", "extension", "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-mssql","",""},
		//{"liquibase-neo4j", "extension", "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-neo4j","",""},
		//{"liquibase-oracle", "extension", "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-oracle","",""},
		//{"liquibase-percona", "extension", "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-percona","",""},
		//{"liquibase-postgresql", "extension", "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-postgresql","",""},
		//{"liquibase-redshift", "extension", "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-redshift","",""},
		//{"liquibase-snowflake", "extension", "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-snowflake","",""},
		//{"liquibase-sqlfire", "extension", "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-sqlfire","",""},
		//{"liquibase-teradata", "extension", "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-teradata","",""},
		//{"liquibase-verticaDatabase", "extension", "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-verticaDatabase","",""},
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
			if m.excludeSuffix != "" && m.includeSuffix == "" {
				if !strings.Contains(f.Text, m.excludeSuffix) {
					versionsRaw = append(versionsRaw, strings.TrimRight(f.Text, "/"))
				}
			}
			if m.excludeSuffix == "" && m.includeSuffix != "" {
				if strings.Contains(f.Text, m.includeSuffix) {
					versionsRaw = append(versionsRaw, strings.TrimRight(f.Text, "/"))
				}
			}
			if m.excludeSuffix != "" && m.includeSuffix != "" {
				if strings.Contains(f.Text, m.includeSuffix) && !strings.Contains(f.Text, m.includeSuffix) {
					versionsRaw = append(versionsRaw, strings.TrimRight(f.Text, "/"))
				}
			}
			if m.excludeSuffix == "" && m.includeSuffix == "" {
				versionsRaw = append(versionsRaw, strings.TrimRight(f.Text, "/"))
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

		// Older versions might have bad version patters ending up with a missing sha. Don't add them.
		if ver.CheckSum != "" {
			p.Versions = append(p.Versions, ver)
		}
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
		} else {
			newPacks = append(newPacks, p)
		}
	}

	//Write all packages back to manifest.
	app.WritePackages(newPacks)
}