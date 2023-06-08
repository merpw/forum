package handlers

import (
	"backend/chat/ws"
	"encoding/json"
	"log"
)

type createChatRequestData struct {
	UserId int `json:"userId"`
}

func (h *Handlers) chatCreate(message Message, client *ws.Client) {
	var data createChatRequestData
	err := json.Unmarshal(message.Item.Data, &data)
	if err != nil {
		log.Println(err)
		return
	}

	// TODO: maybe check if user exists

	chatId := h.DB.CreateChat(client.UserId, data.UserId)

	responseMessage := BuildResponseMessage(message, chatId)

	h.Hub.Broadcast(responseMessage, client.UserId, data.UserId)
}
