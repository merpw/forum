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

	chatId := h.DB.GetUsersChat(client.UserId, companionId)

	var data interface{} = chatId

	if data == -1 {
		data = nil
	}

	responseMessage := ws.BuildResponseMessage(message, data)

	h.Hub.Broadcast(responseMessage, client.UserId)
}
