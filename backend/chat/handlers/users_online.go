package handlers

import "backend/chat/ws"

func (h *Handlers) usersOnline(message ws.Message, client *ws.Client) {
	users := h.Hub.GetOnlineUsers()
	responseMessage := ws.BuildResponseMessage(message, users)

	h.Hub.Broadcast(responseMessage, client.UserId)
}
