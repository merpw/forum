package auth

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
)

// InternalRequest makes an internal request to the auth service.
//
// The request is authenticated with the FORUM_BACKEND_SECRET environment variable.
// The response is decoded into the v variable.
func InternalRequest(method string, url string, body io.Reader, v interface{}) {
	req, err := http.NewRequest(method, os.Getenv("AUTH_BASE_URL")+url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Internal-Auth", os.Getenv("FORUM_BACKEND_SECRET"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return
	}

	if resp.StatusCode != http.StatusOK {
		panic("status code: " + resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(v)
	if err != nil {
		panic(err)
	}

	_ = resp.Body.Close()
}
