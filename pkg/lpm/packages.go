package lpm

import (
	"fmt"
)

//Packages type
type Packages []Package

//GetByName individual package from packages
func (pp Packages) GetByName(n string) Package {
	var r Package
	for _, p := range pp {
		if p.Name != n {
			continue
		}
		r = p
		break
	}
	return r
}

//FilterByCategory get packages by category
func (pp Packages) FilterByCategory(c string) Packages {
	var fpp Packages
	for _, p := range pp {
		if p.Category != c {
			continue
		}
		fpp = append(fpp, p)
	}
	return fpp
}

//Println generates display table for all packages
func (pp Packages) Println(files ClasspathFiles) {
	for _, out := range pp.Display(files) {
		fmt.Println(out)
	}
}

//Display generate display table for packages
func (pp Packages) Display(files ClasspathFiles) []string {
	var r []string
	var prefix string
	r = append(r, fmt.Sprintf("%-4s %-38s %s", "   ", "Package", "Category"))
	for i, p := range pp {
		if (i + 1) == len(pp) {
			prefix = "└──"
		} else {
			prefix = "├──"
		}
		var v string
		tag := p.GetInstalledVersion(files).Tag
		if tag != "" {
			v = "@" + tag
		} else {
			v = tag
		}
		r = append(r, fmt.Sprintf("%-4s %-38s %s", prefix, p.Name+v, p.Category))
	}
	return r
}
