package utils

import (
	"io"
	"net/http"
	"package-manager/internal/app/errors"
)

// HTTPUtil struct
type HTTPUtil struct{}

// Get contents from URL as bytes
func (h HTTPUtil) Get(url string) []byte {
	client := http.Client{}
	r, err := client.Get(url)
	if err != nil {
		errors.Exit("Unable to download from "+url, 1)

	}
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		errors.Exit(err.Error(), 1)
	}
	return body
}
