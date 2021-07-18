package packages

import (
	"bytes"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"package-manager/internal/app/errors"
	"package-manager/internal/app/utils"
	"path/filepath"
	"strings"
)

type Package struct {
	Name  string `json:"name"`
	Category string `json:"category"`
	Path string `json:"path"`
	Algorithm string `json:"algorithm"`
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

func (p Package) PathIsHttp() bool {
	return strings.HasPrefix(p.Path, "http")
}

func writeToDestination(d string, b []byte, f string) {
	destination, err := os.Create(d)
	if err != nil {
		errors.Exit("Unable to access classpath located at " + d, 1)
	}
	defer destination.Close()
	_, err = io.Copy(destination, bytes.NewReader(b))
	if err != nil {
		errors.Exit("Unable to install " + f + " in classpath.", 1)
	}
}

func (p Package) CopyToClassPath(cp string) {
	source, err := os.Open(p.Path)
	if err != nil {
		errors.Exit("Unable to open " + p.Path, 1)
	}
	defer source.Close()
	b, err := ioutil.ReadAll(source)
	if err != nil {
		errors.Exit(err.Error(), 1)
	}
	writeToDestination(cp + p.GetFilename(), b, p.GetFilename())
}

func (p Package) calcChecksum(b []byte) string {
	var r string
	switch p.Algorithm {
	case "SHA1":
		r = fmt.Sprintf("%x", sha1.Sum(b))
	case "SHA256":
		r = fmt.Sprintf("%x", sha256.Sum256(b))
	default:
		errors.Exit("Unknown Algorithm.", 1)
	}
	return r
}

func (p Package) DownloadToClassPath(cp string) {
	body := utils.HttpUtil{}.Get(p.Path)
	sha := p.calcChecksum(body)
	if sha == p.CheckSum {
		fmt.Println("Checksum verified. Installing " + p.GetFilename() + " to " + cp)
	} else {
		errors.Exit("Checksum validation failed. Aborting download.", 1)
	}
	writeToDestination(cp + p.GetFilename(), body, p.GetFilename())
}