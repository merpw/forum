package database

import (
	"log"
	"time"
)

type ChatType int

const (
	ChannelChat ChatType = iota // One to many, one writer, many readers
	GroupChat                   // Several members
	PrivateChat                 // Two members
	AnyChat                     // Any chat type, should be used only for reading
)

// AddChat adds a new empty chat to the database, returns id of the new chat
func (db DB) AddChat(chatType ChatType) int {
	result, err := db.Exec("INSERT INTO chats (type, date) VALUES (?, ?)", chatType,
		time.Now().Format(time.RFC3339))
	if err != nil {
		log.Panic(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Panic(err)
	}
	return int(id)
}

// GetChats reads from database all user's chats with specified ChatType
func (db DB) GetChats(userId int, chatType ChatType) []Chat {
	qs := "SELECT * FROM chats WHERE id IN (SELECT chat_id FROM memberships WHERE user_id = ?)"
	if chatType != AnyChat {
		qs += " AND type = ?"
	}

	query, err := db.Query(qs, userId, chatType)
	if err != nil {
		log.Panic(err)
	}

	var chats []Chat
	for query.Next() {
		var chat Chat
		err = query.Scan(&chat.Id, &chat.Type, &chat.Date)
		if err != nil {
			log.Panic(err)
		}
		chats = append(chats, chat)
	}
	query.Close()

	return chats
}

// AddMembership adds membership to database, returns id of the new membership
func (db DB) AddMembership(userId, chatId int) int {
	result, err := db.Exec(
		"INSERT INTO memberships (user_id, chat_id, date) VALUES (?, ?, ?)",
		userId, chatId, time.Now().Format(time.RFC3339))
	if err != nil {
		log.Panic(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Panic(err)
	}
	return int(id)
}

// AddMessage adds message to database, returns id of the new message
func (db DB) AddMessage(userId, chatId int, content string) int {
	result, err := db.Exec(
		"INSERT INTO messages (user_id, chat_id, content, date) VALUES (?, ?, ?, ?)",
		userId, chatId, content, time.Now().Format(time.RFC3339))
	if err != nil {
		log.Panic(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Panic(err)
	}
	return int(id)
}

func (db DB) GetOnlineUsers(userId int) []User {
	query, err := db.Query(`
		  SELECT users.id, users.name FROM users 
			JOIN sessions ON sessions.user_id = users.id 
			WHERE users.id != ?
		`, userId)
	if err != nil {
		log.Panic(err)
	}

	var users []User
	for query.Next() {
		var user User
		err = query.Scan(&user.Id, &user.Name)
		if err != nil {
			log.Panic(err)
		}
		users = append(users, user)
	}
	query.Close()

	return users
}

// GetContacts reads private chat opponents from database by userId
func (db DB) GetContacts(userId int) []User {
	// TODO: check this properly. It looks not clear. But at the moment there is no data and tests to check it.
	query, err := db.Query(`
			SELECT o.id, o.name FROM users AS o 
			JOIN memberships AS om ON om.user_id = o.id 
			JOIN chats AS c ON c.id = om.chat_id 
			JOIN memberships AS um ON um.chat_id = c.id 
			JOIN users AS u ON u.id = um.user_id 
			WHERE c.type = 2 AND u.id = ?
		`, userId)
	if err != nil {
		log.Panic(err)
	}

	var users []User
	for query.Next() {
		var user User
		err = query.Scan(&user.Id, &user.Name)
		if err != nil {
			log.Panic(err)
		}
		users = append(users, user)
	}
	query.Close()

	return users
}

// GetChatsIds reads chats ids from database by userId
// TODO: implement tests, later, after approving the logic
func (db DB) GetChatsIds(userId int) []int {
	query, err := db.Query("SELECT chat_id FROM memberships WHERE user_id = ?", userId)
	if err != nil {
		log.Panic(err)
	}

	var chatsIds []int
	for query.Next() {
		var chatId int
		err = query.Scan(&chatId)
		if err != nil {
			log.Panic(err)
		}
		chatsIds = append(chatsIds, chatId)
	}
	query.Close()

	return chatsIds
}

// GetChat reads messages from database by chatId
// TODO: implement tests, later, after approving the logic
func (db DB) GetChat(chatId int) []Message {
	query, err := db.Query("SELECT * FROM messages WHERE chat_id = ?", chatId)
	if err != nil {
		log.Panic(err)
	}

	var chat []Message
	for query.Next() {
		var m Message
		err = query.Scan(&m.Id, &m.UserId, &m.ChatId, &m.Content, &m.Date)
		if err != nil {
			log.Panic(err)
		}
		chat = append(chat, m)
	}
	query.Close()

	return chat
}
