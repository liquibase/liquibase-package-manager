package lpm

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

//HttpGet connects from URL as bytes
func HttpGet(url string) (body []byte, err error) {

	callback := func(r *http.Request, via []*http.Request) error {
		r.URL.Opaque = r.URL.Path
		return nil
	}

	client := http.Client{
		CheckRedirect: callback,
	}

	r, err := client.Get(url)
	if err != nil {
		err = fmt.Errorf("unable to download from "+url, 1)
		goto end

	}

	//goland:noinspection GoUnhandledErrorResult
	defer r.Body.Close()
	body, err = ioutil.ReadAll(r.Body)
	if err != nil {
		err = fmt.Errorf("unable to read HTTP response body: %w", err)
		goto end
	}

end:
	return body, err
}
