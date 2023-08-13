package auth

import (
	"log"
	"os"
)

func init() {
	if os.Getenv("AUTH_BASE_URL") == "" {
		log.Println("AUTH_BASE_URL is not set, using default value http://localhost:8080")
		err := os.Setenv("AUTH_BASE_URL", "http://localhost:8080")
		if err != nil {
			log.Fatal(err)
		}
	}

	if os.Getenv("FORUM_BACKEND_SECRET") == "" {
		log.Println("WARNING: FORUM_BACKEND_SECRET is not set, using default value `secret`")
		err := os.Setenv("FORUM_BACKEND_SECRET", "secret")
		if err != nil {
			log.Fatal(err)
		}
	}
}
