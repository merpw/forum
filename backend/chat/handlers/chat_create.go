package handlers

import (
	"backend/chat/external"
	"backend/chat/ws"
	"encoding/json"
	"log"
)

type createChatRequestData struct {
	UserId int `json:"userId"`
}

func (h *Handlers) chatCreate(message ws.Message, client *ws.Client) {
	var data createChatRequestData
	err := json.Unmarshal(message.Item.Data, &data)
	if err != nil {
		log.Println(err)
		return
	}

	user := external.GetUser(data.UserId)
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

	responseMessage := ws.BuildResponseMessage(message, chatId)

	h.Hub.Broadcast(responseMessage, client.UserId, data.UserId)
}
