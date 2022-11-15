package dependencies

import (
	"encoding/json"
	"os"
	"package-manager/internal/app/errors"
)

// FileLocation exported for testing overwrite
var FileLocation string

func init() {
	pwd, err := os.Getwd()
	if err != nil {
		errors.Exit(err.Error(), 1)
	}
	FileLocation = pwd + "/liquibase.json"
}

// Dependencies main wrapper for liquibase.json objects
type Dependencies struct {
	Dependencies []Dependency `json:"dependencies"`
}

// CreateFile init liquibase.json file in pwd
func (d Dependencies) CreateFile() {
	file, err := os.Create(FileLocation)
	if err != nil {
		errors.Exit(err.Error(), 1)
	}
	defer file.Close()
	d.Write()
}

// Write dump contents to liquibase.json
func (d Dependencies) Write() {
	file, err := json.MarshalIndent(d, "", " ")
	if err != nil {
		errors.Exit(err.Error(), 1)
	}
	err = os.WriteFile(FileLocation, file, 0664)
	if err != nil {
		errors.Exit(err.Error(), 1)
	}
}

// Read get contents from liquibase.json
func (d *Dependencies) Read() {
	file, _ := os.Open(FileLocation)
	defer file.Close()
	decoder := json.NewDecoder(file)
	for decoder.More() {
		decoder.Decode(d)
	}
}

// FileExists does the liquibase.json file exist
func (d Dependencies) FileExists() bool {
	_, err := os.Stat(FileLocation)
	return err == nil
}

// Remove remove specific dependency from group
func (d *Dependencies) Remove(n string) {
	for i, m := range d.Dependencies {
		if m.GetName() == n {
			d.Dependencies = append(d.Dependencies[:i], d.Dependencies[i+1:]...)
		}
	}
}
