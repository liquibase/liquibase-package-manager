package packages

import (
	"encoding/json"
	"io/fs"
	"io/ioutil"
	"path/filepath"
)

type Package struct {
	Name  string `json:"name"`
	Category string `json:"category"`
	Path string `json:"path"`
	CheckSum string `json:"checksum"`
	IsInstalled bool `json:"isInstalled"`
}

type Packages []Package

func LoadPackages() Packages {
	data, _ := ioutil.ReadFile("./packages.json")
	var e []Package
	json.Unmarshal(data, &e)
	return e
}

func (p Package) GetFilename() string {
	_, f := filepath.Split(p.Path)
	return f
}

func (ps Packages) GetByName(n string) Package {
	var r Package
	for _, p := range ps {
		if p.Name == n {
			r = p
		}
	}
	return r
}

func (ps Packages) FilterByCategory(c string) Packages {
	var r Packages
	for _, p := range ps {
		if p.Category == c {
			r = append(r, p)
		}
	}
	return r
}

func (p Package) InClassPath(cp []fs.FileInfo) bool {
	r := false
	for _, f := range cp {
		if f.Name() == p.GetFilename() {
			r = true
		}
	}
	return r
}