package auth

import (
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

// CheckSession creates an API request to the auth service to check if the session is valid.
func CheckSession(token string) (userId int) {
	req, err := http.NewRequest(http.MethodGet, os.Getenv("AUTH_BASE_URL")+"/api/internal/check-session?token="+token, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Internal-Auth", os.Getenv("FORUM_BACKEND_SECRET"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Panicf("status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	_ = resp.Body.Close()

	userId, err = strconv.Atoi(string(body))
	if err != nil {
		return -1
	}

	return userId
}
