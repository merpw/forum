package handlers

import (
	"backend/chat/ws"
	"encoding/json"
	"log"
	"strconv"
	"strings"
)

type messageRequestData struct {
	Content string `json:"content"`
}

func (h *Handlers) chatIdMessage(message ws.Message, client *ws.Client) {
	chatIdStr := strings.Split(message.Item.URL, "/")[2]
	// "/chat/1/message" -> "1"

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

	var data messageRequestData
	err = json.Unmarshal(message.Item.Data, &data)
	if err != nil {
		log.Println(err)
		return
	}

	messageId := h.DB.CreateMessage(chatId, client.UserId, data.Content)
	responseMessage := ws.BuildResponseMessage(message, messageId)

	chatMembers := h.DB.GetChatMembers(chatId)
	h.Hub.Broadcast(responseMessage, chatMembers...)
}
