package handlers

import (
	"backend/chat/external"
	"backend/chat/ws"
	"encoding/json"
	"fmt"
	"log"
)

// handshake checks if the messageBody is a handshake message and if the token is valid.
//
// If the token is valid, the userId is set to the client.
func (h *Handlers) handshake(messageBody []byte, client *ws.Client) {
	var message struct {
		Type string `json:"type"`
		Item struct {
			Token string `json:"token"`
		}
	}

	err := json.Unmarshal(messageBody, &message)
	if err != nil {
		log.Println(err)
	}
	if message.Type != "handshake" {
		log.Println("ERROR: no handshake")
		return
	}

	if message.Item.Token == "" {
		log.Println("ERROR: no token")
		return
	}

	userId := external.CheckSession(message.Item.Token)
	if userId == -1 {
		log.Println("ERROR: invalid token")
		return
	}
	client.UserId = userId
	client.Token = message.Item.Token

	err = client.Conn.WriteMessage(1, []byte(fmt.Sprintf(`{"type":"handshake","item":{"data":{"userId":%d}}}`, userId)))
	if err != nil {
		log.Println(err)
	}

	h.Hub.BroadcastOnlineStatus()
}
