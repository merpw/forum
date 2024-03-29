package ws

import (
	"encoding/json"
	"fmt"
)

// Message is a basic websocket message
type Message struct {
	Type string `json:"type"`
	Item struct {
		URL  string          `json:"url"`
		Data json.RawMessage `json:"data"`
	} `json:"item"`
}

// BuildResponseMessage builds a response Message with given data
func BuildResponseMessage(requestMessage Message, data interface{}) Message {
	dataRawMessage, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	return Message{
		Type: requestMessage.Type,
		Item: struct {
			URL  string          `json:"url"`
			Data json.RawMessage `json:"data"`
		}{
			URL:  requestMessage.Item.URL,
			Data: dataRawMessage,
		},
	}
}

func getGroupChatMessage(groupId int) Message {
	return Message{
		Type: "get",
		Item: struct {
			URL  string          `json:"url"`
			Data json.RawMessage `json:"data"`
		}{
			URL: fmt.Sprintf("/groups/%d/chat", groupId),
		},
	}
}
