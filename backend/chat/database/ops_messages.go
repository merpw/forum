package database

func (db DB) GetChatMessages(chatId int) (messages []int) {
	q, err := db.Query("SELECT id FROM messages WHERE chat_id = ?", chatId)
	if err != nil {
		panic(err)
	}
	defer q.Close()

	messages = []int{}
	for q.Next() {
		var messageId int
		err = q.Scan(&messageId)
		if err != nil {
			panic(err)
		}
		messages = append(messages, messageId)
	}

	return messages
}

type Message struct {
	Id       int
	ChatId   int
	AuthorId int
	Content  string
	Time     string
}

func (db DB) GetMessage(messageId int) (message Message) {
	q, err := db.Query("SELECT id, chat_id, user_id, content, timestamp FROM messages WHERE id = ?", messageId)
	if err != nil {
		panic(err)
	}
	defer q.Close()

	if q.Next() {
		err = q.Scan(&message.Id, &message.ChatId, &message.AuthorId, &message.Content, &message.Time)
		if err != nil {
			panic(err)
		}
	}

	return message
}

func (db DB) CreateMessage(chatId, authorId int, content string) int {
	res, err := db.Exec("INSERT INTO messages (chat_id, user_id, content) VALUES (?, ?, ?)", chatId, authorId, content)
	if err != nil {
		panic(err)
	}

	messageId, err := res.LastInsertId()
	if err != nil {
		panic(err)
	}

	return int(messageId)
}
