package handlers

import (
	"backend/chat/ws"
	"encoding/json"
	"log"
	"strconv"
	"strings"
)

// broadcast typing event to all chat members

type chatIdTypingRequest struct {
	IsTyping bool `json:"isTyping"`
}

type chatIdTypingBroadcastData struct {
	UserId   int  `json:"userId"`
	IsTyping bool `json:"isTyping"`
}

func (h *Handlers) chatIdTyping(message ws.Message, client *ws.Client) {
	chatIdStr := strings.Split(message.Item.URL, "/")[2]
	// "/chat/1/typing" -> "1"

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

	members := h.DB.GetChatMembers(chatId)

	// remove sender from members
	for i, member := range members {
		if member == client.UserId {
			members = append(members[:i], members[i+1:]...)
			break
		}
	}

	var reqData chatIdTypingRequest

	err = json.Unmarshal(message.Item.Data, &reqData)
	if err != nil {
		log.Println(err)
		return
	}

	data := chatIdTypingBroadcastData{
		UserId:   client.UserId,
		IsTyping: reqData.IsTyping,
	}

	responseMessage := ws.BuildResponseMessage(message, data)

	h.Hub.Broadcast(responseMessage, members...)
}
