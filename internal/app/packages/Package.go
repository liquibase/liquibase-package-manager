package packages

import (
	"crypto/sha1"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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

func (p Package) PathIsHttp() bool {
	return strings.HasPrefix(p.Path, "http")
}

func writeToDestination(d string, r io.Reader, f string) {
	destination, err := os.Create(d)
	if err != nil {
		fmt.Println("Unable to access classpath located at " + d)
		os.Exit(1)
	}
	defer destination.Close()
	_, err = io.Copy(destination, r)
	if err != nil {
		fmt.Println("Unable to install " + f + " in classpath.")
		os.Exit(1)
	}
}

func (p Package) CopyToClassPath(cp string) {
	source, err := os.Open(p.Path)
	if err != nil {
		fmt.Println("Unable to open " + p.Path)
		os.Exit(1)
	}
	defer source.Close()
	writeToDestination(cp + p.GetFilename(), source, p.GetFilename())
}

func (p Package) DownloadToClassPath(cp string) {
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}
	r, err := client.Get(p.Path)
	if err != nil {
		fmt.Println("Unable to download from " + p.Path)
		os.Exit(1)

	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	sha := fmt.Sprintf("%x", sha1.Sum(body))
	if sha == p.CheckSum {
		fmt.Println("Checksum verified. Installing " + p.GetFilename() + " to " + cp)
	} else {
		fmt.Println("Checksum validation failed. Aborting download.")
		os.Exit(1)
	}
	writeToDestination(cp + p.GetFilename(), r.Body, p.GetFilename())
}