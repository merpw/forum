package handlers

import (
	"backend/chat/ws"
	"backend/common/integrations/auth"
	"encoding/json"
	"log"
)

type createChatRequestData struct {
	UserId  *int `json:"userId"`
	GroupId *int `json:"groupId"`
}

type createChatResponseData struct {
	ChatId  int  `json:"chatId"`
	UserId  *int `json:"userId,omitempty"`
	GroupId *int `json:"groupId,omitempty"`
}

func (h *Handlers) chatCreate(message ws.Message, client *ws.Client) {
	var data createChatRequestData
	err := json.Unmarshal(message.Item.Data, &data)
	if err != nil {
		log.Println(err)
		return
	}

	if data.UserId == nil && data.GroupId == nil {
		log.Println("ERROR: no userId or groupId specified")
		return
	}

	if data.UserId != nil && data.GroupId != nil {
		log.Println("ERROR: both userId and groupId specified")
		return
	}

	if data.UserId != nil {
		user := auth.GetUser(*data.UserId)
		if user == nil {
			log.Println("ERROR: user not found")
			return
		}

		if h.DB.GetPrivateChat(client.UserId, *data.UserId) != nil {
			log.Println("ERROR: chat already exists")
			return
		}

		if *data.UserId == client.UserId {
			// TODO: maybe allow chats with yourself
			log.Println("ERROR: cannot create chat with yourself")
			return
		}

		chatId := h.DB.CreatePrivateChat(client.UserId, *data.UserId)

		h.Hub.Broadcast(ws.BuildResponseMessage(message, createChatResponseData{
			ChatId: chatId,
			UserId: data.UserId,
		}), client.UserId)

		h.Hub.Broadcast(ws.BuildResponseMessage(message, createChatResponseData{
			ChatId: chatId,
			UserId: &client.UserId,
		}), *data.UserId)
	}

	members := auth.GetGroupMembers(*data.GroupId)
	if members == nil {
		log.Println("ERROR: group not found")
		return
	}

	found := false
	for _, member := range *members {
		if member == client.UserId {
			found = true
			break
		}
	}

	if !found {
		log.Println("ERROR: user is not a member of the group")
		return
	}

	if h.DB.GetGroupChat(*data.GroupId) != nil {
		log.Println("ERROR: chat already exists")
		return
	}

	chatId := h.DB.CreateGroupChat(*data.GroupId, *members)

	h.Hub.Broadcast(ws.BuildResponseMessage(message, createChatResponseData{
		ChatId:  chatId,
		GroupId: data.GroupId,
	}), *members...)
}
