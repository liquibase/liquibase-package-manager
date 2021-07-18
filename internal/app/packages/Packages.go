package packages

import (
	"fmt"
	"io/fs"
)

type Packages []Package

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

func (ps Packages) Display(files []fs.FileInfo) {
	var prefix string
	h := fmt.Sprintf("%-4s %-38s %s", "   ", "Package", "Category")
	fmt.Println(h)
	for i, s := range ps {
		if (i+1) == len(ps) {
			prefix = "└──"
		} else {
			prefix = "├──"
		}
		var v string
		if s.GetDefaultVersion().InClassPath(files) {
			v = "@" + s.GetDefaultVersion().Tag
		}
		l := fmt.Sprintf("%-4s %-38s %s", prefix, s.Name + v, s.Category)
		fmt.Println(l)
	}
}