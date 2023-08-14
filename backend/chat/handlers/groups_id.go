package handlers

import (
	"backend/chat/ws"
	"log"
	"strconv"
	"strings"
)

func (h *Handlers) groupsIdChat(message ws.Message, client *ws.Client) {
	groupIdStr := strings.Split(message.Item.URL, "/")[2]
	// "/groups/1/chat" -> "1"

	groupId, err := strconv.Atoi(groupIdStr)
	if err != nil {
		log.Println(err)
		return
	}

	chatId := h.DB.GetGroupChat(groupId)
	responseMessage := ws.BuildResponseMessage(message, chatId)

	h.Hub.Broadcast(responseMessage, client.UserId)
}
