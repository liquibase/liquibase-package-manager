package packages

import (
	"encoding/json"
	"io/ioutil"
)

type Packages []Package

func LoadPackages() Packages {
	data, _ := ioutil.ReadFile("./packages.json")
	var e Packages
	err := json.Unmarshal(data, &e)
	if err != nil {
		return nil
	}
	return e
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