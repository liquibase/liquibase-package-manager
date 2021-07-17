package packages

import (
	"encoding/json"
	"fmt"
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

func (ps Packages) Display() {
	var prefix string
	h := fmt.Sprintf("%-4s %-30s %s", "   ", "Package", "Category")
	fmt.Println(h)
	for i, s := range ps {
		if (i+1) == len(ps) {
			prefix = "└──"
		} else {
			prefix = "├──"
		}
		l := fmt.Sprintf("%-4s %-30s %s", prefix, s.Name, s.Category)
		fmt.Println(l)
	}
}