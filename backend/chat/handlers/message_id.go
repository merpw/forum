package handlers

import (
	"backend/chat/ws"
	"log"
	"strconv"
	"strings"
)

func (h *Handlers) messageId(message ws.Message, client *ws.Client) {
	messageIdStr := strings.Split(message.Item.URL, "/")[2]
	// "/message/1" -> "1"

	messageId, err := strconv.Atoi(messageIdStr)
	if err != nil {
		log.Println(err)
		return
	}

	msg := h.DB.GetMessage(messageId)

	if !h.DB.IsChatMember(msg.ChatId, client.UserId) {
		h.Hub.Broadcast(ws.BuildResponseMessage(message, nil), client.UserId)
		return
	}

	type ResponseData struct {
		Id        int    `json:"id"`
		ChatId    int    `json:"chatId"`
		Content   string `json:"content"`
		AuthorId  int    `json:"authorId"`
		Timestamp string `json:"timestamp"`
	}

	responseMessage := ws.BuildResponseMessage(message, ResponseData{
		Id:        msg.Id,
		ChatId:    msg.ChatId,
		Content:   msg.Content,
		AuthorId:  msg.AuthorId,
		Timestamp: msg.Time,
	})

	h.Hub.Broadcast(responseMessage, client.UserId)
}
