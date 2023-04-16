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

// TODO: perhaps move it later to models.go. Not approved yet. Structure was changed. Full refactoring is needed.
// type PrivateChatOponent struct {
// 	Id                  int    // opoent id
// 	Name                string // oponent name
// 	ChatId              int    // chat id
// 	ChatLastMessageDate string // last message date, to sort chats by last message date
// }

// TODO: remove it later. Not approved yet. Structure was changed. Full refactoring is needed, include comments.
/*
* GetPrivateChatOponentsByUserId reads private chat oponents from database by user_id,
does not require user to be logged in.

 Use this to get the list of private chats for a user.
 Can be sorted by last_message_date on frontend. To show in the list of private chats.

The oponent is the user in the private chat that is not the user with the given user_id.
The oponent has Name and Id.
The ChatId is the id of the private chat. Only two users can be in a private chat.
Returns nil if no such chats exist.
*/
// func (db DB) GetPrivateChatOponentsByUserId(userId int) []PrivateChatOponent {
// 	query, err := db.Query(
// 		"SELECT users.id, users.name, chats.id, chats.last_message_date FROM users "+
// 			"JOIN memberships ON memberships.user_id = users.id "+
// 			"JOIN chats ON chats.id = memberships.chat_id "+
// 			"WHERE users.id != ? AND chats.id IN (SELECT chat_id FROM memberships GROUP BY chat_id HAVING COUNT(*) = 2)",
// 		userId)
// 	if err != nil {
// 		log.Panic(err)
// 	}

// 	var oponents []PrivateChatOponent
// 	for query.Next() {
// 		var oponent PrivateChatOponent
// 		err = query.Scan(&oponent.Id, &oponent.Name, &oponent.ChatId, &oponent.ChatLastMessageDate)
// 		if err != nil {
// 			log.Panic(err)
// 		}
// 		oponents = append(oponents, oponent)
// 	}
// 	query.Close()

// 	return oponents
// }

/*
*
GetChatsIdsByUserId reads chats ids from database by user_id, does not require user to be logged in
*/
func (db DB) GetChatsIdsByUserId(userId int) []int {
	query, err := db.Query("SELECT chat_id FROM memberships WHERE user_id = ?", userId)
	if err != nil {
		log.Panic(err)
	}

	var chatIds []int
	for query.Next() {
		var chatId int
		err = query.Scan(&chatId)
		if err != nil {
			log.Panic(err)
		}
		chatIds = append(chatIds, chatId)
	}
	query.Close()

	return chatIds
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
		userId)
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
