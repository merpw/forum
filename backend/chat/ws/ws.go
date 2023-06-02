// Package ws implements websocket connections
//
// It calls the specified primary MessageHandler on each message from the Client
package ws

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	mu      sync.Mutex
	Clients []*Client

	MessageHandler MessageHandler
}

type MessageHandler func(p []byte, client *Client)

// NewHub creates new Hub
func NewHub(messageHandler MessageHandler) *Hub {
	return &Hub{
		MessageHandler: messageHandler,
	}
}

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	// TODO: remove this in production
	CheckOrigin: func(r *http.Request) (b bool) { return true },
}

// UpgradeHandler upgrades HTTP connection to WebSocket
func (h *Hub) UpgradeHandler(w http.ResponseWriter, r *http.Request) {
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

func (h *Hub) Register(client *Client) {
	h.mu.Lock()
	h.Clients = append(h.Clients, client)
	client.Conn.SetCloseHandler(func(code int, text string) error {
		h.Unregister(client)
		return nil
	})
	h.mu.Unlock()
}

func (h *Hub) Unregister(client *Client) {
	h.mu.Lock()
	for i, c := range h.Clients {
		if c == client {
			h.Clients = append(h.Clients[:i], h.Clients[i+1:]...)
			break
		}
	}
	h.mu.Unlock()
}

type BroadcastFunc func(data interface{}, clients ...int)

// Broadcast sends data as JSON to all clients
//
// If no clients specified, sends to all
func (h *Hub) Broadcast(data interface{}, clients ...int) {
	if clients == nil {
		// not specified, send to all
		for _, c := range h.Clients {
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
