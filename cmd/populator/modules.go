package main

var modules Modules

//Modules main module splice
type Modules []Module

func (mm Modules) getByName(n string) Module {
	var m Module
	for _, e := range mm {
		if e.name == n {
			m = e
			goto end
		}
	}
end:
	return m
}

func init() {
	modules = []Module{
		{
			name:        "aws-secrets-manager",
			category:    Pro,
			url:         "https://maven.liquibase.com/org/liquibase/ext/secrets/aws-secrets-manager",
			artifactory: Maven{},
		},
		{
			name:        "cloudbees-feature-management",
			category:    Pro,
			url:         "https://maven.liquibase.com/org/liquibase/ext/precondition/cloudbees-feature-management",
			artifactory: Maven{},
		},
		{
			name:        "configcat",
			category:    Pro,
			url:         "https://maven.liquibase.com/org/liquibase/ext/precondition/configcat",
			artifactory: Maven{},
		},
		{
			name:        "ff4j",
			category:    Extension,
			owner:       "liquibase",
			repo:        "ff4j-extension",
			artifactory: Github{},
		},
		{
			name:        "flagr",
			category:    Extension,
			owner:       "liquibase",
			repo:        "flagr-extension",
			artifactory: Github{},
		},
		{
			name:        "flagsmith",
			category:    Pro,
			url:         "https://maven.liquibase.com/org/liquibase/ext/precondition/flagsmith",
			artifactory: Maven{},
		},
		{
			name:        "flipt",
			category:    Extension,
			owner:       "liquibase",
			repo:        "flipt-extension",
			artifactory: Github{},
		},
		{
			name:        "gitlab-feature-flags",
			category:    Pro,
			url:         "https://maven.liquibase.com/org/liquibase/ext/precondition/gitlab-feature-flags",
			artifactory: Maven{},
		},
		{
			name:        "growthbook",
			category:    Pro,
			url:         "https://maven.liquibase.com/org/liquibase/ext/precondition/growthbook",
			artifactory: Maven{},
		},
		{
			name:        "launchdarkly",
			category:    Pro,
			url:         "https://maven.liquibase.com/org/liquibase/ext/precondition/launchdarkly",
			artifactory: Maven{},
		},
		{
			name:        "custom-datatype-converter",
			category:    Extension,
			url:         "https://maven.liquibase.com/org/liquibase/ext/datatype/custom-datatype-converter",
			artifactory: Maven{},
		},
		{
			name:        "mongodb",
			category:    Driver,
			url:         "https://repo1.maven.org/maven2/org/mongodb/mongo-java-driver",
			filePrefix:  "mongo-java-driver-",
			artifactory: Maven{},
		},
		{
			name:        "custom-hosts",
			category:    Extension,
			owner:       "liquibase",
			repo:        "custom-hosts-extension",
			artifactory: Github{},
		},
		{
			name:        "cyberarkpam-vault",
			category:    Pro,
			url:         "https://maven.liquibase.com/org/liquibase/ext/vaults/cyberarkpam-vault",
			artifactory: Maven{},
		},
		{
			name:        "hashicorp-vault",
			category:    Pro,
			url:         "https://repo1.maven.org/maven2/org/liquibase/ext/vaults/hashicorp-vault",
			artifactory: Maven{},
		},
		{
			name:        "liquibase-data",
			category:    Extension,
			owner:       "liquibase",
			repo:        "liquibase-data",
			artifactory: Github{},
		},
		{
			name:          "postgresql",
			category:      Driver,
			url:           "https://repo1.maven.org/maven2/org/postgresql/postgresql",
			excludeSuffix: ".jre",
			artifactory:   Maven{},
		},
		{
			name:          "mssql",
			category:      Driver,
			url:           "https://repo1.maven.org/maven2/com/microsoft/sqlserver/mssql-jdbc",
			includeSuffix: ".jre11",
			excludeSuffix: ".jre11-preview",
			filePrefix:    "mssql-jdbc-",
			artifactory:   Maven{},
		},
		{
			name:        "mariadb",
			category:    Driver,
			url:         "https://repo1.maven.org/maven2/org/mariadb/jdbc/mariadb-java-client",
			filePrefix:  "mariadb-java-client-",
			artifactory: Maven{},
		},
		{
			name:        "h2",
			category:    Driver,
			url:         "https://repo1.maven.org/maven2/com/h2database/h2",
			artifactory: Maven{},
		},
		{
			name:          "db2",
			category:      Driver,
			url:           "https://repo1.maven.org/maven2/com/ibm/db2/jcc",
			excludeSuffix: "db2",
			filePrefix:    "jcc-",
			artifactory:   Maven{},
		},
		{
			name:        "snowflake",
			category:    Driver,
			url:         "https://repo1.maven.org/maven2/net/snowflake/snowflake-jdbc",
			filePrefix:  "snowflake-jdbc-",
			artifactory: Maven{},
		},
		{
			name:        "sybase",
			category:    Driver,
			url:         "https://repo1.maven.org/maven2/net/sf/squirrel-sql/plugins/sybase",
			artifactory: Maven{},
		},
		{
			name:        "firebird",
			category:    Driver,
			url:         "https://repo1.maven.org/maven2/net/sf/squirrel-sql/plugins/firebird",
			artifactory: Maven{},
		},
		{
			name:        "sqlite",
			category:    Driver,
			url:         "https://repo1.maven.org/maven2/org/xerial/sqlite-jdbc",
			filePrefix:  "sqlite-jdbc-",
			artifactory: Maven{},
		},
		{
			name:        "oracle",
			category:    Driver,
			url:         "https://repo1.maven.org/maven2/com/oracle/ojdbc/ojdbc8",
			filePrefix:  "ojdbc8-",
			artifactory: Maven{},
		},
		{
			name:        "mysql",
			category:    Driver,
			url:         "https://repo1.maven.org/maven2/mysql/mysql-connector-java",
			filePrefix:  "mysql-connector-java-",
			artifactory: Maven{},
		},
		{
			name:        "liquibase-cache",
			category:    Extension,
			url:         "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-cache",
			artifactory: Maven{},
		},
		{
			name:        "liquibase-cassandra",
			category:    Extension,
			url:         "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-cassandra",
			artifactory: Maven{},
		},
		{
			name:        "liquibase-cosmosdb",
			category:    Extension,
			url:         "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-cosmosdb",
			artifactory: Maven{},
		},
		{
			name:        "liquibase-db2i",
			category:    Extension,
			url:         "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-db2i",
			artifactory: Maven{},
		},
		{
			name:        "liquibase-filechangelog",
			category:    Extension,
			url:         "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-filechangelog",
			artifactory: Maven{},
		},
		{
			name:        "liquibase-hanadb",
			category:    Extension,
			url:         "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-hanadb",
			artifactory: Maven{},
		},
		{
			name:        "liquibase-hibernate5",
			category:    Extension,
			url:         "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-hibernate5",
			artifactory: Maven{},
		},
		{
			name:        "liquibase-maxdb",
			category:    Extension,
			url:         "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-maxdb",
			artifactory: Maven{},
		},
		{
			name:        "liquibase-modify-column",
			category:    Extension,
			url:         "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-modify-column",
			artifactory: Maven{},
		},
		{
			name:        "liquibase-mongodb",
			category:    Extension,
			url:         "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-mongodb",
			artifactory: Maven{},
		},
		{
			name:        "liquibase-mssql",
			category:    Extension,
			url:         "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-mssql",
			artifactory: Maven{},
		},
		{
			name:        "liquibase-neo4j",
			category:    Extension,
			url:         "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-neo4j",
			artifactory: Maven{},
		},
		{
			name:        "liquibase-oracle",
			category:    Extension,
			url:         "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-oracle",
			artifactory: Maven{},
		},
		{
			name:        "liquibase-percona",
			category:    Extension,
			url:         "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-percona",
			artifactory: Maven{},
		},
		{
			name:        "liquibase-postgresql",
			category:    Extension,
			url:         "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-postgresql",
			artifactory: Maven{},
		},
		{
			name:        "liquibase-redshift",
			category:    Extension,
			url:         "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-redshift",
			artifactory: Maven{},
		},
		{
			name:        "liquibase-snowflake",
			category:    Extension,
			url:         "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-snowflake",
			artifactory: Maven{},
		},
		{
			name:        "liquibase-sqlfire",
			category:    Extension,
			url:         "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-sqlfire",
			artifactory: Maven{},
		},
		{
			name:        "liquibase-teradata",
			category:    Extension,
			url:         "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-teradata",
			artifactory: Maven{},
		},
		{
			name:        "liquibase-verticaDatabase",
			category:    Extension,
			url:         "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-verticaDatabase",
			artifactory: Maven{},
		},
	}
}
