package ws

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	Conn   *websocket.Conn
	UserId int
	Token  string
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		Conn:   conn,
		UserId: -1,
	}
}

// Read calls specified MessageHandler on each message from the Client
func (c *Client) Read(messageHandler MessageHandler) {
	for {
		_, p, err := c.Conn.ReadMessage()
		if err != nil {
			return
		}

		messageHandler(p, c)
	}
}
