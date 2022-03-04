package utils

import (
	"archive/zip"
	"bufio"
	"github.com/hashicorp/go-version"
	"io"
	"log"
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
	r, err := zip.OpenReader(hp + "liquibase.jar")
	if err != nil {
		panic(err)
	}

	props := map[string]string{}
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
						props[key] = value
					}
				}
				if err == io.EOF {
					break
				}
			}
		}
	}
	v, _ := version.NewVersion(props["build.version"])
	return Liquibase{
		Homepath:        hp,
		Version:         v,
		BuildProperties: props,
	}
}
