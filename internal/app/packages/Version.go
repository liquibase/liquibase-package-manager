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

type Version struct {
	Tag string `json:"tag"`
	Path string `json:"path"`
	Algorithm string `json:"algorithm"`
	CheckSum string `json:"checksum"`
}

func (v Version) GetFilename() string {
	_, f := filepath.Split(v.Path)
	return f
}

func (v Version) InClassPath(files []fs.FileInfo) bool {
	r := false
	for _, f := range files {
		if f.Name() == v.GetFilename() {
			r = true
		}
	}
	return r
}

func (v Version) PathIsHttp() bool {
	return strings.HasPrefix(v.Path, "http")
}

func (v Version) CopyToClassPath(cp string) {
	source, err := os.Open(v.Path)
	if err != nil {
		errors.Exit("Unable to open " + v.Path, 1)
	}
	defer source.Close()
	b, err := ioutil.ReadAll(source)
	if err != nil {
		errors.Exit(err.Error(), 1)
	}
	writeToDestination(cp + v.GetFilename(), b, v.GetFilename())
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
	//err = os.Chmod(d, 0664)
	//if err != nil {
	//	errors.Exit(err.Error(), 1)
	//}
}

func (v Version) calcChecksum(b []byte) string {
	var r string
	switch v.Algorithm {
	case "SHA1":
		r = fmt.Sprintf("%x", sha1.Sum(b))
	case "SHA256":
		r = fmt.Sprintf("%x", sha256.Sum256(b))
	default:
		errors.Exit("Unknown Algorithm.", 1)
	}
	return r
}

func (v Version) DownloadToClassPath(cp string) {
	body := utils.HttpUtil{}.Get(v.Path)
	sha := v.calcChecksum(body)
	if sha == v.CheckSum {
		fmt.Println("Checksum verified. Installing " + v.GetFilename() + " to " + cp)
	} else {
		errors.Exit("Checksum validation failed. Aborting download.", 1)
	}
	writeToDestination(cp + v.GetFilename(), body, v.GetFilename())
}