package database

const CreateChatMessage = "Chat created"

// CreatePrivateChat creates a new user-user chat, adds members and sends a system message about it.
func (db DB) CreatePrivateChat(creatorId, userId int) (chatId int) {
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
		id, userId)
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

// CreateGroupChat creates a new group chat, adds members and sends a system message about it.
func (db DB) CreateGroupChat(groupId int, members []int) (chatId int) {
	r, err := db.Exec("INSERT INTO chats (group_id, type) VALUES (?, 1)", groupId)
	if err != nil {
		panic(err)
	}
	id, err := r.LastInsertId()
	if err != nil {
		panic(err)
	}

	for _, memberId := range members {
		_, err = db.Exec(`
			INSERT INTO memberships (chat_id, user_id) VALUES (?, ?);`,
			id, memberId)
		if err != nil {
			panic(err)
		}
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
