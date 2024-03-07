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

// Liquibase struct
type Liquibase struct {
	Homepath        string
	Version         *version.Version
	BuildProperties map[string]string
}

// LoadLiquibase loads liquibase struct from home path
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
		r, err = zip.OpenReader(hp + "internal/lib/liquibase-core.jar")
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

			if l.BuildProperties["build.version"] != "DEV" {
				v, _ := version.NewVersion(l.BuildProperties["build.version"])
				l.Version = v
			} else {
				l.Version, _ = version.NewVersion("999.0.0")
			}

		}
	}
	goto end
end:
	return l
}
