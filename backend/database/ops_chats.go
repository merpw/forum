package database

import (
	"log"
	"time"
)

/*
AddChat adds chat to database, returns id of new chat

	chatType:
	  2 (private 1vs1) or
	  1 (group chat) or
	  0 (the channel owner is posting to subscribers)
*/
func (db DB) AddChat(chatType int) int {
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

// GetChatsByUserId reads chats from database by user_id, does not require user to be logged in
func (db DB) GetChatsByUserId(userId int) []Chat {
	query, err := db.Query("SELECT * FROM chats WHERE id IN (SELECT chat_id FROM memberships WHERE user_id = ?)", userId)
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

/** GetPrivateChatsByUserId reads private chats from database by user_id, does not require user to be logged in
 *  Returns nil if no such chats exist.
 */
func (db DB) GetPrivateChatsByUserId(userId int) []Chat {
	query, err := db.Query(
		"SELECT * FROM chats WHERE id IN (SELECT chat_id FROM memberships WHERE user_id = ?) AND type = 2",
		userId,
	)
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

// AddMembership adds membership to database, returns id of new membership
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

// AddMessage adds message to database, returns id of new message, plus updates last_message_date in chats table
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

// TODO: perhaps move it later to models.go. Not approved yet.
type OnlineUser struct {
	Id   int
	Name string
}

/*
GetOnlineUsersIdsAndNames reads online users from database, does not require user to be logged in.

	Use this to get the list of online users for a user side panel with available users to chat with.
*/
func (db DB) GetOnlineUsersIdsAndNames(userId int) []OnlineUser {
	query, err := db.Query(
		"SELECT users.id, users.name FROM users "+
			"JOIN sessions ON sessions.user_id = users.id "+
			"WHERE users.id != ?",
		userId)
	if err != nil {
		log.Panic(err)
	}

	var onlineUsers []OnlineUser
	for query.Next() {
		var onlineUser OnlineUser
		err = query.Scan(&onlineUser.Id, &onlineUser.Name)
		if err != nil {
			log.Panic(err)
		}
		onlineUsers = append(onlineUsers, onlineUser)
	}
	query.Close()

	return onlineUsers
}

// GetPrivateChatOponentsByUserId reads private chat oponents from database by user_id
func (db DB) GetPrivateChatOponentsByUserId(userId int) []OnlineUser {
	// requirements:
	// chat.type = 2 to get private chats,
	// users.id != userId to get oponents, not the user itself
	// o = oponent, u = user
	// TODO: check this properly. It looks not clear. But at the moment there is no data and tests to check it.
	query, err := db.Query(
		"SELECT o.id, o.name FROM users AS o "+
			"JOIN memberships AS om ON om.user_id = o.id "+
			"JOIN chats AS c ON c.id = om.chat_id "+
			"JOIN memberships AS um ON um.chat_id = c.id "+
			"JOIN users AS u ON u.id = um.user_id "+
			"WHERE c.type = 2 AND o.id != ? AND u.id = ?",
		userId, userId)
	if err != nil {
		log.Panic(err)
	}

	var oponents []OnlineUser
	for query.Next() {
		var oponent OnlineUser
		err = query.Scan(&oponent.Id, &oponent.Name)
		if err != nil {
			log.Panic(err)
		}
		oponents = append(oponents, oponent)
	}
	query.Close()

	return oponents
}
