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
	e := `{"dependencies":[]}`
	err = json.Unmarshal([]byte(e), &d)
	if err != nil {
		errors.Exit(err.Error(), 1)
	}
	err = ioutil.WriteFile(fileLocation, []byte(e), 0664)
	if err != nil {
		errors.Exit(err.Error(), 1)
	}
}

func (d Dependencies) FileExists() bool {
	_, err := os.Stat(fileLocation)
	return err == nil
}