package dependencies

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"package-manager/internal/app/errors"
)

var fileLocation string

func init() {
	pwd, err := os.Getwd()
	if err != nil {
		errors.Exit(err.Error(), 1)
	}
	fileLocation = pwd + "/liquibase.json"
}

type Dependencies struct {
	Dependencies []Dependency `json:"dependencies"`
}

type Dependency map[string]string

func (d Dependencies) CreateFile() {
	file, err := os.Create(fileLocation)
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
	err = ioutil.WriteFile(fileLocation, file, 0664)
	if err != nil {
		errors.Exit(err.Error(), 1)
	}
}

func (d *Dependencies) Read() {
	file, _ := os.Open(fileLocation)
	defer file.Close()
	decoder := json.NewDecoder(file)
	for decoder.More() {
		decoder.Decode(d)
	}
}

func (d Dependencies) FileExists() bool {
	_, err := os.Stat(fileLocation)
	return err == nil
}

func (d Dependency) GetName() string {
	var r string
	for k := range d {
		r = k
	}
	return r
}

func (d Dependency) GetVersion() string {
	var r string
	for _, v := range d {
		r = v
	}
	return r
}

func (d *Dependencies) Remove(n string) {
	for i, m := range d.Dependencies {
		if m.GetName() == n {
			d.Dependencies = append(d.Dependencies[:i], d.Dependencies[i+1:]...)
		}
	}
}