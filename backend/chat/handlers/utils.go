package handlers

type ResponseMessage struct {
	Type string `json:"type"`
	Item struct {
		URL  string      `json:"url"`
		Data interface{} `json:"data"`
	} `json:"item"`
}

// BuildResponseMessage builds a ResponseMessage with the given data
func BuildResponseMessage(requestMessage Message, data interface{}) ResponseMessage {
	return ResponseMessage{
		Type: requestMessage.Type,
		Item: struct {
			URL  string      `json:"url"`
			Data interface{} `json:"data"`
		}{
			URL:  requestMessage.Item.URL,
			Data: data,
		},
	}
}
