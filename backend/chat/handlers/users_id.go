package handlers

import (
	"backend/chat/ws"
	"log"
	"strconv"
	"strings"
)

func (h *Handlers) usersIdChat(message ws.Message, client *ws.Client) {
	companionIdStr := strings.Split(message.Item.URL, "/")[2]
	// "/users/1/chat" -> "1"

	companionId, err := strconv.Atoi(companionIdStr)
	if err != nil {
		log.Println(err)
		return
	}

	chatId := h.DB.GetPrivateChat(client.UserId, companionId)
	responseMessage := ws.BuildResponseMessage(message, chatId)

	h.Hub.Broadcast(responseMessage, client.UserId)
}
