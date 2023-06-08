package external

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func init() {
	if os.Getenv("FRONTEND_REVALIDATE_URL") == "" {
		fmt.Println("WARNING: FRONTEND_REVALIDATE_URL is not set, frontend revalidation will not work")
	}
}

// RevalidateURL creates POST request to frontend to revalidate url
//
// Uses environment variables FRONTEND_REVALIDATE_URL and optional FRONTEND_REVALIDATE_TOKEN
//
// Does nothing if FRONTEND_REVALIDATE_URL is not set
func RevalidateURL(url string) {
	apiURL := os.Getenv("FRONTEND_REVALIDATE_URL")
	if apiURL == "" {
		return
	}
	req, err := http.NewRequest(http.MethodPost, apiURL, nil)
	if err != nil {
		log.Printf("revalidation of %s failed: %s", url, err)
		return
	}

	q := req.URL.Query()
	q.Add("url", url)
	q.Add("token", os.Getenv("FRONTEND_REVALIDATE_TOKEN"))
	req.URL.RawQuery = q.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("revalidation of %s failed: %s", url, err)
		return
	}
	if res.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			log.Printf("revalidation of %s failed: %s", url, err)
			return
		}
		log.Printf("revalidation of %s failed: %s", url, string(bodyBytes))
	}
}
