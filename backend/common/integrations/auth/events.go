package auth

import (
	"bufio"
	"log"
	"net/http"
	"os"
	"time"
)

// Events function returns a channel with revoked tokens.
//
// Returns a channel with revoked tokens.
func Events() <-chan string {
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

		req, err := http.NewRequest(http.MethodGet, os.Getenv("AUTH_BASE_URL")+"/api/internal/events", nil)
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
