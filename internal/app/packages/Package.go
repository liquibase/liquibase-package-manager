package packages

import (
	"io/fs"
	"path/filepath"
)

type Package struct {
	Name  string `json:"name"`
	Category string `json:"category"`
	Path string `json:"path"`
	CheckSum string `json:"checksum"`
}

func (p Package) GetFilename() string {
	_, f := filepath.Split(p.Path)
	return f
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