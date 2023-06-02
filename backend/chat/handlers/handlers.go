// Package handlers handles WebSocket Message's
package handlers

import (
	"backend/chat/database"
	"backend/chat/external"
	"backend/chat/ws"
	"database/sql"
	"encoding/json"
	"log"
	"regexp"
	"runtime/debug"
)

// PrimaryHandler returns ws.MessageHandler with all routes registered
func (h *Handlers) PrimaryHandler() ws.MessageHandler {

	var events = []Event{
		// method GET endpoints
		newEvent("get", `/chat/all`, h.chatAll),
		newEvent("get", `/chat/\d+`, h.chatId),
		newEvent("get", `/chat/\d+/messages`, h.chatIdMessages),

		newEvent("get", `/message/\d+`, h.messageId),

		newEvent("get", `/users/\d+/chat`, h.usersIdChat),

		// method POST endpoints
		newEvent("post", `/chat/create`, h.chatCreate),
		newEvent("post", `/chat/\d+/message`, h.chatIdMessage),
	}

	return func(messageBody []byte, client *ws.Client) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("ERROR 500, %s\n%s", r, debug.Stack())
			}
		}()

		log.Println(string(messageBody))

		if client.UserId == -1 {
			h.handshake(messageBody, client)
			return
		}

		// TODO: maybe remove, as we have a revoked-sessions endpoint
		userId := external.CheckSession(client.Token)
		if userId != client.UserId {
			log.Println("ERROR: invalid token")
			_ = client.Conn.Close()
			return
		}

		var message Message
		err := json.Unmarshal(messageBody, &message)
		if err != nil {
			log.Println(err)
		}

		for _, event := range events {
			if event.Type == message.Type && event.Path.MatchString(message.Item.URL) {
				event.Handler(message, client)
				return
			}
		}

		log.Printf("WARN: no handler for %s %s\n", message.Type, message.Item.URL)
	}
}

type Handlers struct {
	DB *database.DB

	Broadcast ws.BroadcastFunc
}

// New connects database to Handlers
func New(db *sql.DB) *Handlers {
	return &Handlers{DB: database.New(db)}
}

// Message is a basic websocket message
type Message struct {
	Type string `json:"type"`
	Item struct {
		URL  string          `json:"url"`
		Data json.RawMessage `json:"data"`
	} `json:"item"`
}

// Event is a websocket event (server.Route analog)
type Event struct {
	Type    string
	Path    *regexp.Regexp
	Handler func(message Message, client *ws.Client)
}

func newEvent(method, path string, handler func(message Message, client *ws.Client)) Event {
	return Event{
		Type:    method,
		Path:    regexp.MustCompile(`^` + path + `$`),
		Handler: handler,
	}
}
