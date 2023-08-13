// Package ws implements websocket connections
//
// It calls the specified primary MessageHandler on each message from the Client
package ws

import (
	"backend/chat/database"
	"backend/common/integrations/auth"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Hub is a WebSocket hub that handles all WebSocket connections.
//
// It calls the specified primary MessageHandler on each message from the Client
type Hub struct {
	mu      sync.Mutex
	Clients []*Client

	// MessageHandler is a primary MessageHandler, it is called on each message from the Client
	MessageHandler MessageHandler
}

// MessageHandler is a function that handles raw messages from the Client
type MessageHandler func(p []byte, client *Client)

// NewHub creates new Hub with specified primary MessageHandler
func NewHub(messageHandler MessageHandler, db *database.DB) *Hub {
	h := &Hub{
		MessageHandler: messageHandler,
	}

	go func() {
		for event := range auth.Events() {
			switch event.Type {
			case auth.EventTypeTokenRevoked:
				token := string(event.Item.(json.RawMessage))
				for _, c := range h.Clients {
					if c.Token == token {
						_ = c.Conn.Close()
					}
				}
			case auth.EventTypeGroupJoin:
				var item auth.EventGroupItem
				err := json.Unmarshal(event.Item.(json.RawMessage), &item)
				if err != nil {
					log.Println(err)
					continue
				}
				chatId := db.GetGroupChat(item.GroupId)
				if chatId != nil {
					db.AddChatMember(*chatId, item.UserId)
				}

				message := BuildResponseMessage(getGroupChatMessage(item.GroupId), chatId)
				h.Broadcast(message, item.UserId)
			case auth.EventTypeGroupLeave:
				var item auth.EventGroupItem
				err := json.Unmarshal(event.Item.(json.RawMessage), &item)
				if err != nil {
					log.Println(err)
					continue
				}
				chatId := db.GetGroupChat(item.GroupId)
				if chatId != nil {
					db.RemoveChatMember(*chatId, item.UserId)
				}
				message := BuildResponseMessage(getGroupChatMessage(item.GroupId), nil)
				h.Broadcast(message, item.UserId)
			}
		}
	}()

	return h
}

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	// TODO: remove this in production
	CheckOrigin: func(r *http.Request) (b bool) { return true },
}

// UpgradeHandler upgrades HTTP connection to WebSocket
func (h *Hub) UpgradeHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("ERROR 500: ", err)
		}
	}()

	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("New connection from ", r.RemoteAddr)

	client := NewClient(conn)
	h.Register(client)

	go func() {
		client.Read(h.MessageHandler)
		h.Unregister(client)
	}()
}

// Register registers new Client in the Hub
func (h *Hub) Register(client *Client) {
	h.mu.Lock()
	h.Clients = append(h.Clients, client)
	client.Conn.SetCloseHandler(func(code int, text string) error {
		h.Unregister(client)
		return nil
	})
	h.mu.Unlock()
}

// Unregister unregisters Client from the Hub
func (h *Hub) Unregister(client *Client) {
	h.mu.Lock()
	if client.Token != "" {
		defer h.BroadcastOnlineStatus()
	}
	for i, c := range h.Clients {
		if c == client {
			h.Clients = append(h.Clients[:i], h.Clients[i+1:]...)
			break
		}
	}
	h.mu.Unlock()
}

// Broadcast sends data as JSON to all clients
//
// If no clients specified, sends to all
func (h *Hub) Broadcast(data interface{}, clients ...int) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if clients == nil {
		// not specified, send to all
		for _, c := range h.Clients {
			if c.Token == "" {
				continue
			}
			err := c.Conn.WriteJSON(data)
			if err != nil {
				log.Println(err)
			}
		}
		return
	}

	for _, id := range clients {
		for _, c := range h.Clients {
			if c.UserId == id {
				err := c.Conn.WriteJSON(data)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}
}

// BroadcastOnlineStatus Broadcast shortcut for sending online status to all clients
func (h *Hub) BroadcastOnlineStatus(clients ...int) {
	message := BuildResponseMessage(Message{
		Type: "get",
		Item: struct {
			URL  string          `json:"url"`
			Data json.RawMessage `json:"data"`
		}{
			URL: "/users/online",
		},
	}, h.GetOnlineUsers())

	h.Broadcast(message, clients...)
}

func (h *Hub) GetOnlineUsers() []int {
	var onlineUsers []int

nextClient:
	for _, c := range h.Clients {
		for _, id := range onlineUsers {
			if id == c.UserId {
				continue nextClient
			}
		}
		// append only unique users
		onlineUsers = append(onlineUsers, c.UserId)
	}
	return onlineUsers
}
