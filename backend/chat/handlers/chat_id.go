package handlers

import (
	"backend/chat/ws"
	"log"
	"strconv"
	"strings"
)

type chatIdResponseData struct {
	Id            int  `json:"id"`
	CompanionId   *int `json:"userId,omitempty"`
	GroupId       *int `json:"groupId,omitempty"`
	LastMessageId int  `json:"lastMessageId"`
}

func (h *Handlers) chatId(message ws.Message, client *ws.Client) {
	chatIdStr := strings.Split(message.Item.URL, "/")[2]
	// "/chat/1" -> "1"

	chatId, err := strconv.Atoi(chatIdStr)
	if err != nil {
		log.Println(err)
		return
	}

	isMember := h.DB.IsChatMember(chatId, client.UserId)
	if !isMember {
		h.Hub.Broadcast(ws.BuildResponseMessage(message, nil), client.UserId)
		return
	}

	data := h.DB.GetChatData(client.UserId, chatId)

	responseData := chatIdResponseData{
		Id:            chatId,
		LastMessageId: data.LastMessageId,
		CompanionId:   data.CompanionId,
		GroupId:       data.GroupId,
	}
	responseMessage := ws.BuildResponseMessage(message, responseData)

	h.Hub.Broadcast(responseMessage, client.UserId)
}

// chatIdMessages responds with all messages from a chat
func (h *Handlers) chatIdMessages(message ws.Message, client *ws.Client) {
	chatIdStr := strings.Split(message.Item.URL, "/")[2]
	// "/chat/1/messages" -> "1"

	chatId, err := strconv.Atoi(chatIdStr)
	if err != nil {
		log.Println(err)
		return
	}

	isMember := h.DB.IsChatMember(chatId, client.UserId)
	if !isMember {
		h.Hub.Broadcast(ws.BuildResponseMessage(message, nil), client.UserId)
		return
	}

	chats := h.DB.GetChatMessages(chatId)
	responseMessage := ws.BuildResponseMessage(message, chats)
	h.Hub.Broadcast(responseMessage, client.UserId)
	if err != nil {
		log.Println(err)
		return
	}
}
