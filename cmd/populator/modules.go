package main

var modules Modules

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
			name: "liquibase-data",
			category: "extension",
			owner: "liquibase",
			repo: "liquibase-data",
			artifactory: Github{},
		},
		{
			name:          "postgresql",
			category:      "driver",
			url:           "https://repo1.maven.org/maven2/org/postgresql/postgresql",
			excludeSuffix: ".jre",
			artifactory: Maven{},
		},
		//{
		//	name:          "mssql",
		//	category:      "driver",
		//	url:           "https://repo1.maven.org/maven2/com/microsoft/sqlserver/mssql-jdbc",
		//	includeSuffix: ".jre11",
		//	excludeSuffix: ".jre11-preview",
		//	filePrefix:    "mssql-jdbc-",
		//},
		//{"mariadb", "driver", "https://repo1.maven.org/maven2/org/mariadb/jdbc/mariadb-java-client", "", "", "mariadb-java-client-"},
		//{"h2", "driver", "https://repo1.maven.org/maven2/com/h2database/h2", "", "", ""},
		//{"db2", "driver", "https://repo1.maven.org/maven2/com/ibm/db2/jcc", "", "db2", "jcc-"},
		//{"snowflake", "driver", "https://repo1.maven.org/maven2/net/snowflake/snowflake-jdbc", "", "", "snowflake-jdbc-"},
		//{"sybase", "driver", "https://repo1.maven.org/maven2/net/sf/squirrel-sql/plugins/sybase", "", "", ""},
		//{"firebird", "driver", "https://repo1.maven.org/maven2/net/sf/squirrel-sql/plugins/firebird", "", "", ""},
		//{"sqlite", "driver", "https://repo1.maven.org/maven2/org/xerial/sqlite-jdbc", "", "", "sqlite-jdbc-"},
		//{"oracle", "driver", "https://repo1.maven.org/maven2/com/oracle/ojdbc/ojdbc8", "", "", "ojdbc8-"},
		//{"mysql", "driver", "https://repo1.maven.org/maven2/mysql/mysql-connector-java", "", "", "mysql-connector-java-"},
		//{"liquibase-cache", "extension", "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-cache", "", "", ""},
		//{"liquibase-cassandra", "extension", "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-cassandra", "", "", ""},
		//{"liquibase-cosmosdb", "extension", "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-cosmosdb", "", "", ""},
		//{"liquibase-db2i", "extension", "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-db2i", "", "", ""},
		//{"liquibase-filechangelog", "extension", "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-filechangelog", "", "", ""},
		//{"liquibase-hanadb", "extension", "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-hanadb", "", "", ""},
		//{"liquibase-hibernate5", "extension", "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-hibernate5", "", "", ""},
		//{"liquibase-maxdb", "extension", "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-maxdb", "", "", ""},
		//{"liquibase-modify-column", "extension", "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-modify-column", "", "", ""},
		//{"liquibase-mongodb", "extension", "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-mongodb", "", "", ""},
		//{"liquibase-mssql", "extension", "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-mssql", "", "", ""},
		//{"liquibase-neo4j", "extension", "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-neo4j", "", "", ""},
		//{"liquibase-oracle", "extension", "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-oracle", "", "", ""},
		//{"liquibase-percona", "extension", "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-percona", "", "", ""},
		//{"liquibase-postgresql", "extension", "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-postgresql", "", "", ""},
		//{"liquibase-redshift", "extension", "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-redshift", "", "", ""},
		//{"liquibase-snowflake", "extension", "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-snowflake", "", "", ""},
		//{"liquibase-sqlfire", "extension", "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-sqlfire", "", "", ""},
		//{"liquibase-teradata", "extension", "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-teradata", "", "", ""},
		//{"liquibase-verticaDatabase", "extension", "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-verticaDatabase", "", "", ""},
		//"liquibase-compat",
		//"liquibase-javalogger",
		//"liquibase-nochangeloglock",
		//"liquibase-nochangelogupdate",
		//"liquibase-sequencetable",
		//"liquibase-vertica",
	}
}
