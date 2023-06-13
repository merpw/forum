package handlers

import (
	"backend/chat/ws"
	"log"
	"strconv"
	"strings"
)

type chatIdResponseData struct {
	Id            int `json:"id"`
	CompanionId   int `json:"companionId"`
	LastMessageId int `json:"lastMessageId"`
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
		_ = client.Conn.WriteJSON(ws.BuildResponseMessage(message, nil))
		return
	}

	data := h.DB.GetChatData(client.UserId, chatId)

	responseData := chatIdResponseData{
		Id:            chatId,
		LastMessageId: data.LastMessageId,
		CompanionId:   data.CompanionId,
	}
	responseMessage := ws.BuildResponseMessage(message, responseData)

	err = client.Conn.WriteJSON(responseMessage)
	if err != nil {
		log.Println(err)
		return
	}
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
		_ = client.Conn.WriteJSON(ws.BuildResponseMessage(message, nil))
		return
	}

	chats := h.DB.GetChatMessages(chatId)
	responseMessage := ws.BuildResponseMessage(message, chats)
	err = client.Conn.WriteJSON(responseMessage)
	if err != nil {
		log.Println(err)
		return
	}
}
