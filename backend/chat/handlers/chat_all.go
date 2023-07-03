package handlers

import (
	"backend/chat/ws"
)

// chatAll returns a slice of all user's chats.
func (h *Handlers) chatAll(message ws.Message, client *ws.Client) {
	chats := h.DB.GetUserChats(client.UserId)
	responseMessage := ws.BuildResponseMessage(message, chats)

	h.Hub.Broadcast(responseMessage, client.UserId)
}
