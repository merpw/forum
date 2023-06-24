package database

const CreateChatMessage = "Chat created"

// CreateChat creates a new user-user chat, adds members and sends a system message about it.
func (db DB) CreateChat(creatorId, companionId int) (chatId int) {
	r, err := db.Exec("INSERT INTO chats DEFAULT VALUES")
	if err != nil {
		panic(err)
	}
	id, err := r.LastInsertId()
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
		INSERT INTO memberships (chat_id, user_id) VALUES (?, ?);
		INSERT INTO memberships (chat_id, user_id) VALUES (?, ?);`,
		id, creatorId,
		id, companionId)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
		INSERT INTO messages (chat_id, user_id, content) 
		VALUES (?, ?, ?)`,
		id, -1, CreateChatMessage)
	if err != nil {
		panic(err)
	}

	return int(id)
}
