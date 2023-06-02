package external

import (
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

func init() {
	if os.Getenv("AUTH_BASE_URL") == "" {
		log.Println("AUTH_BASE_URL is not set, using default value http://localhost:8080")
		err := os.Setenv("AUTH_BASE_URL", "http://localhost:8080")
		if err != nil {
			log.Fatal(err)
		}
	}
}

// CheckSession creates an API request to the auth service to check if the session is valid.
func CheckSession(token string) (userId int) {
	resp, err := http.Get(os.Getenv("AUTH_BASE_URL") + "/api/internal/check-session?token=" + token)
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
