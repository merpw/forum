package external

import (
	"bufio"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
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

// RevokedSessions creates a connection to the auth SSE stream endpoint.
//
// Returns a channel with revoked tokens.
func RevokedSessions() <-chan string {
	subscriber := make(chan string)

	reconnectDelay := 1
	maxReconnectDelay := 60

	var reconnect func()

	reconnect = func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("ERROR 500, Auth service unavailable", r)
				time.Sleep(time.Duration(reconnectDelay) * time.Second)
				if reconnectDelay < maxReconnectDelay {
					reconnectDelay *= 2
				}
				reconnect()
			} else {
				log.Println("Auth API stream closed")
				close(subscriber)
			}
		}()

		req, err := http.NewRequest(http.MethodGet, os.Getenv("AUTH_BASE_URL")+"/api/internal/revoked-sessions", nil)
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

		if reconnectDelay > 1 {
			log.Println("Auth API stream connected")
			reconnectDelay = 1
		}

		reader := bufio.NewReader(resp.Body)

		for {
			token, err := reader.ReadString('\n')
			if err != nil {
				panic(err)
			}
			subscriber <- token[:len(token)-1] // remove newline
		}
	}
	go reconnect()

	return subscriber
}
