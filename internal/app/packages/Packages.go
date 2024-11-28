package packages

import (
	"fmt"
	"io/fs"

	"github.com/hashicorp/go-version"
)

// Packages type
type Packages []Package

// GetByName individual package from packages
func (ps Packages) GetByName(n string) Package {
	var r Package
	for _, p := range ps {
		if p.Name == n {
			r = p
		}
	}
	return r
}

// GetInstalled get slice of installed packages in classpath
func (ps Packages) GetInstalled(cpFiles []fs.FileInfo) Packages {
	var r Packages
	for _, e := range ps {
		v := e.GetInstalledVersion(cpFiles)
		if v.InClassPath(cpFiles) {
			r = append(r, e)
		}
	}
	return r
}

// GetOutdated get slice of outdated packaged in classpath
func (ps Packages) GetOutdated(lb *version.Version, cpFiles []fs.FileInfo) Packages {
	var r Packages
	for _, p := range ps.GetInstalled(cpFiles) {
		i := p.GetInstalledVersion(cpFiles)
		iv, _ := version.NewVersion(i.Tag)
		l := p.GetLatestVersion(lb)
		lv, _ := version.NewVersion(l.Tag)
		if iv.LessThan(lv) {
			r = append(r, p)
		}
	}
	return r
}

// FilterByCategory get packages by catetory
func (ps Packages) FilterByCategory(c string) Packages {
	var r Packages
	for _, p := range ps {
		if p.Category == c {
			r = append(r, p)
		}
	}
	return r
}

// Display generate display table for packages
func (ps Packages) Display(files []fs.FileInfo) []string {
	var r []string
	var prefix string
	r = append(r, fmt.Sprintf("%-4s %-38s %s", "   ", "Package", "Category"))
	for i, p := range ps {
		if (i + 1) == len(ps) {
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
