package handlers

import "backend/chat/ws"

func (h *Handlers) usersOnline(message Message, client *ws.Client) {
	users := h.Hub.GetOnlineUsers()
	responseMessage := BuildResponseMessage(message, users)

	err := client.Conn.WriteJSON(responseMessage)
	if err != nil {
		return
	}
}
