package packages

import (
	"fmt"
	"io/fs"
)

//Packages type
type Packages []Package

//GetByName individual package from packages
func (ps Packages) GetByName(n string) Package {
	var r Package
	for _, p := range ps {
		if p.Name == n {
			r = p
		}
	}
	return r
}

//FilterByCategory get packages by catetory
func (ps Packages) FilterByCategory(c string) Packages {
	var r Packages
	for _, p := range ps {
		if p.Category == c {
			r = append(r, p)
		}
	}
	return r
}

//Display generate display table for packages
func (ps Packages) Display(files []fs.FileInfo) []string {
	var r []string
	var prefix string
	r = append(r, fmt.Sprintf("%-4s %-38s %s", "   ", "Package", "Category"))
	for i, s := range ps {
		if (i+1) == len(ps) {
			prefix = "└──"
		} else {
			prefix = "├──"
		}
		//TODO installed vs uninstalled
		v := "@" + s.GetInstalledVersion(files).Tag
		r = append(r, fmt.Sprintf("%-4s %-38s %s", prefix, s.Name + v, s.Category))
	}
	return r
}