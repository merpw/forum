package handlers

import (
	"backend/chat/ws"
	"log"
)

// chatAll returns a slice of all user's chats.
func (h *Handlers) chatAll(message Message, client *ws.Client) {
	chats := h.DB.GetUserChats(client.UserId)
	responseMessage := BuildResponseMessage(message, chats)

	err := client.Conn.WriteJSON(responseMessage)
	if err != nil {
		log.Println(err)
		return
	}
}
