package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/hashicorp/go-version"
	"package-manager/internal/app/packages"
	"package-manager/internal/app/utils"
	"sort"
	"strings"
)

var exentions = []string{
	"liquibase-cache",
	"liquibase-cassandra",
	//"liquibase-compat",
	"liquibase-cosmosdb",
	"liquibase-db2i",
	"liquibase-filechangelog",
	"liquibase-hanadb",
	"liquibase-hibernate5",
	//"liquibase-javalogger",
	"liquibase-maxdb",
	"liquibase-modify-column",
	"liquibase-mongodb",
	"liquibase-mssql",
	"liquibase-neo4j",
	//"liquibase-nochangeloglock",
	//"liquibase-nochangelogupdate",
	"liquibase-oracle",
	"liquibase-percona",
	"liquibase-postgresql",
	"liquibase-redshift",
	//"liquibase-sequencetable",
	"liquibase-snowflake",
	"liquibase-sqlfire",
	"liquibase-teradata",
	//"liquibase-vertica",
	"liquibase-verticaDatabase",
}

func populateJSON(n string) {
	var pack packages.Package
	pack.Name = n
	pack.Category = "extension"
	url := "https://repo1.maven.org/maven2/org/liquibase/ext/" + pack.Name

	var versionsRaw []string

	// Get Versions from Root package site
	c := colly.NewCollector()
	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		if !strings.Contains(e.Text, "../") && !strings.Contains(e.Text, "maven-metadata.") {
			versionsRaw = append(versionsRaw, strings.TrimRight(e.Text, "/"))
		}
	})
	c.Visit(url)

	// Sort Versions
	versions := make([]*version.Version, len(versionsRaw))
	for i, raw := range versionsRaw {
		v, _ := version.NewVersion(raw)
		versions[i] = v
	}
	sort.Sort(version.Collection(versions))

	for _, v := range versions {
		var ver packages.Version
		ver.Tag = v.String()
		ver.Path = url + "/" + v.String() + "/" + pack.Name + "-" + v.String() + ".jar"
		ver.Algorithm = "SHA1"
		sha := string(utils.HTTPUtil{}.Get(ver.Path + ".sha1"))
		if strings.Contains(sha, "html") {
			sha = ""
		}
		ver.CheckSum = sha
		pack.Versions = append(pack.Versions, ver)
	}

	json, err := json.MarshalIndent(pack, "", "  ")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(string(json) + ",")
}

func main(){
	for _, e := range exentions {
		populateJSON(e)
	}
}