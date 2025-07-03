package utils

import (
	"archive/zip"
	"bufio"
	"github.com/hashicorp/go-version"
	"io"
	"log"
	"os"
	"strings"
)

//Liquibase struct
type Liquibase struct {
	Homepath        string
	Version         *version.Version
	BuildProperties map[string]string
}

//LoadLiquibase loads liquibase struct from home path
func LoadLiquibase(hp string) Liquibase {
	l := Liquibase{
		Homepath:        hp,
		BuildProperties: map[string]string{},
	}

	var r *zip.ReadCloser
	var err error
	if _, err = os.Stat(hp + "liquibase.jar"); err == nil {
		r, err = zip.OpenReader(hp + "liquibase.jar")
	} else {
		r, err = getLiquibaseJarReader(hp)
	}
	if err != nil {
		z, _ := version.NewVersion("0.0.0")
		l.Version = z
		goto end
	}

	for _, f := range r.File {
		if f.Name == "liquibase.build.properties" {
			file, err := f.Open()
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()

			reader := bufio.NewReader(file)

			for {
				line, err := reader.ReadString('\n')

				// check if the line has = sign
				// and process the line. Ignore the rest.
				if equal := strings.Index(line, "="); equal >= 0 {
					if key := strings.TrimSpace(line[:equal]); len(key) > 0 {
						value := ""
						if len(line) > equal {
							value = strings.TrimSpace(line[equal+1:])
						}
						// assign the config map
						l.BuildProperties[key] = value
					}
				}
				if err == io.EOF {
					break
				}
			}
			v, err := version.NewVersion(l.BuildProperties["build.version"])
			if err != nil {
				// Log the parsing error before falling back to version "0.0.0"
				log.Printf("Error parsing version '%s': %v. Falling back to version '0.0.0'.", l.BuildProperties["build.version"], err)
				v, _ = version.NewVersion("0.0.0")
			}
			l.Version = v
		}
	}
	goto end
end:
	return l
}

// getLiquibaseJarReader checks for liquibase-core.jar first, then liquibase-commercial.jar
func getLiquibaseJarReader(hp string) (*zip.ReadCloser, error) {
	corePath := hp + "internal/lib/liquibase-core.jar"
	commercialPath := hp + "internal/lib/liquibase-commercial.jar"

	if _, err := os.Stat(corePath); err == nil {
		return zip.OpenReader(corePath)
	}

	if _, err := os.Stat(commercialPath); err == nil {
		return zip.OpenReader(commercialPath)
	}

	return nil, os.ErrNotExist
}
