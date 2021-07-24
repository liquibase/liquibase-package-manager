package utils

import (
	"io/ioutil"
	"net/http"
	"package-manager/internal/app/errors"
)

//HTTPUtil struct
type HTTPUtil struct {}

//Get contencts from URL as bytes
func (h HTTPUtil) Get(url string) []byte {
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}
	r, err := client.Get(url)
	if err != nil {
		errors.Exit("Unable to download from " + url, 1)

	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errors.Exit(err.Error(), 1)
	}
	return body
}
