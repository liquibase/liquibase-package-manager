package dependencies

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"package-manager/internal/app/errors"
)

var FileLocation string

func init() {
	pwd, err := os.Getwd()
	if err != nil {
		errors.Exit(err.Error(), 1)
	}
	FileLocation = pwd + "/liquibase.json"
}

type Dependencies struct {
	Dependencies []Dependency `json:"dependencies"`
}

func (d Dependencies) CreateFile() {
	file, err := os.Create(FileLocation)
	if err != nil {
		errors.Exit(err.Error(), 1)
	}
	defer file.Close()
	d.Write()
}

func (d Dependencies) Write() {
	file, err := json.MarshalIndent(d, "", " ")
	if err != nil {
		errors.Exit(err.Error(), 1)
	}
	err = ioutil.WriteFile(FileLocation, file, 0664)
	if err != nil {
		errors.Exit(err.Error(), 1)
	}
}

func (d *Dependencies) Read() {
	file, _ := os.Open(FileLocation)
	defer file.Close()
	decoder := json.NewDecoder(file)
	for decoder.More() {
		decoder.Decode(d)
	}
}

func (d Dependencies) FileExists() bool {
	_, err := os.Stat(FileLocation)
	return err == nil
}

func (d *Dependencies) Remove(n string) {
	for i, m := range d.Dependencies {
		if m.GetName() == n {
			d.Dependencies = append(d.Dependencies[:i], d.Dependencies[i+1:]...)
		}
	}
}