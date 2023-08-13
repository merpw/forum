package handlers

import (
	"backend/chat/ws"
	"backend/common/integrations/auth"
	"encoding/json"
	"log"
)

type createChatRequestData struct {
	UserId int `json:"userId"`
}

type createChatResponseData struct {
	ChatId int `json:"chatId"`
	UserId int `json:"userId"`
}

func (h *Handlers) chatCreate(message ws.Message, client *ws.Client) {
	var data createChatRequestData
	err := json.Unmarshal(message.Item.Data, &data)
	if err != nil {
		log.Println(err)
		return
	}

	user := auth.GetUser(data.UserId)
	if user == nil {
		log.Println("ERROR: user not found")
		return
	}

	if h.DB.GetUsersChat(client.UserId, data.UserId) != -1 {
		log.Println("ERROR: chat already exists")
		return
	}

	if data.UserId == client.UserId {
		// TODO: maybe allow chats with yourself
		log.Println("ERROR: cannot create chat with yourself")
		return
	}

	chatId := h.DB.CreateChat(client.UserId, data.UserId)

	h.Hub.Broadcast(ws.BuildResponseMessage(message, createChatResponseData{
		ChatId: chatId,
		UserId: data.UserId,
	}), client.UserId)

	h.Hub.Broadcast(ws.BuildResponseMessage(message, createChatResponseData{
		ChatId: chatId,
		UserId: client.UserId,
	}), data.UserId)
}
