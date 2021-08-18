package lpm

import (
	"bytes"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type ChecksumAlgorithm string

const Sha1Algorithm ChecksumAlgorithm = "SHA1"

//Version struct
type Version struct {
	Tag       string            `json:"tag"`
	Path      string            `json:"path"`
	Algorithm ChecksumAlgorithm `json:"algorithm"`
	CheckSum  string            `json:"checksum"`
}

func (v Version) CopyFilesToClassPath(cp string) (err error) {
	if !v.PathIsHTTP() {
		err = v.CopyToClassPath(cp)
	} else {
		err = v.DownloadToClassPath(cp)
	}
	return err
}

//GetFilename from VersionNumber
func (v Version) GetFilename() string {
	_, f := filepath.Split(v.Path)
	return f
}

//InClassPath VersionNumber is installed in classpath
func (v Version) InClassPath(files ClasspathFiles) bool {
	r := false
	for _, f := range files {
		if f.Name() == v.GetFilename() {
			r = true
			break
		}
	}
	return r
}

//PathIsHTTP remote or local file path
func (v Version) PathIsHTTP() bool {
	return strings.HasPrefix(v.Path, "http")
}

//CopyToClassPath install local VersionNumber to classpath
func (v Version) CopyToClassPath(cp string) error {
	var b []byte
	source, err := os.Open(v.Path)
	if err != nil {
		err = fmt.Errorf("unable to open %s when copying to classpath; %w",
			v.Path,
			err)
		goto end
	}
	//goland:noinspection GoUnhandledErrorResult
	defer source.Close()
	b, err = ioutil.ReadAll(source)
	if err != nil {
		err = fmt.Errorf("unable to read from %s when copying to classpath; %w",
			v.Path,
			err)
		goto end
	}
	err = writeToDestination(cp+v.GetFilename(), b)
end:
	return err
}

func writeToDestination(dest string, buf []byte) error {
	destination, err := os.Create(dest)
	if err != nil {
		err = fmt.Errorf("unable to create file %s while writing to destination; %w",
			dest,
			err)
		goto end
	}
	//goland:noinspection GoUnhandledErrorResult
	defer destination.Close()
	_, err = io.Copy(destination, bytes.NewReader(buf))
	if err != nil {
		err = fmt.Errorf("unable to write contents to file %s; %w",
			dest,
			err)
		goto end
	}
end:
	return err
}

func (v Version) CalcChecksum(b []byte) (r string, err error) {
	switch v.Algorithm {
	case "SHA1":
		r = fmt.Sprintf("%x", sha1.Sum(b))
	case "SHA256":
		r = fmt.Sprintf("%x", sha256.Sum256(b))
	default:
		err = fmt.Errorf("unknown checksum algorithm: '%s'", string(b))
	}
	return r, err
}

//DownloadToClassPath install remote VersionNumber to classpath
func (v Version) DownloadToClassPath(cp string) (err error) {
	var body []byte
	var sha, msg, fn, fp string

	body, err = HttpGet(v.Path)
	if sha != v.CheckSum {
		err = fmt.Errorf("failed to download class path %s; %w", cp, err)
		goto end
	}

	sha, err = v.CalcChecksum(body)
	if err != nil {
		msg = "unable to calculate checksum"
		goto end
	}
	if sha != v.CheckSum {
		msg = "checksum not valid"
		goto end
	}
	fn = v.GetFilename()
	fmt.Printf("Checksum verified. Installing %s to %s ", fn, cp)

	fp = fmt.Sprintf("%s%s", cp, fn)
	err = writeToDestination(fp, body)
	if err != nil {
		msg = fmt.Sprintf("unable to write to destination %s", fp)
		goto end
	}

end:
	if msg != "" {
		err = fmt.Errorf("%s when downloading class path %s; %w",
			msg,
			cp,
			err)
	}
	return err
}
