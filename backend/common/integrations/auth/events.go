package auth

import (
	"bufio"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

type EventType string

const (
	// EventTypeTokenRevoked is sent when a token is revoked. The item is the revoked token (string).
	EventTypeTokenRevoked EventType = "token_revoked"
	// EventTypeGroupJoin is sent when a user joins a group. The item is EventGroupItem.
	EventTypeGroupJoin EventType = "group_join"
	// EventTypeGroupLeave is sent when a user leaves a group. The item is EventGroupItem.
	EventTypeGroupLeave EventType = "group_leave"
)

type Event struct {
	Type EventType   `json:"type"`
	Item interface{} `json:"item"`
}

type EventGroupItem struct {
	GroupId int `json:"group_id"`
	UserId  int `json:"user_id"`
}

// Events function connects to the auth service and returns a channel of Event objects
//
// It reconnects automatically if the connection is lost.
func Events() <-chan Event {
	subscriber := make(chan Event)

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
			line, err := reader.ReadBytes('\n')
			if err != nil {
				panic(err)
			}

			type RawEvent struct {
				Type EventType       `json:"type"`
				Item json.RawMessage `json:"item"`
			}
			var rawEvent RawEvent
			err = json.Unmarshal(line, &rawEvent)
			if err != nil {
				panic(err)
			}

			subscriber <- Event{
				Type: rawEvent.Type,
				Item: rawEvent.Item,
			}
		}
	}
	go reconnect()

	return subscriber
}
