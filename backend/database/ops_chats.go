package database

import (
	"log"
	"time"
)

// AddChat adds chat to database, returns id of new chat
func (db DB) AddChat() int {
	result, err := db.Exec("INSERT INTO chats (last_message_date) VALUES (?)", time.Now().Format(time.RFC3339))
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
		err = query.Scan(&chat.Id, &chat.LastMessageDate)
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
		"SELECT * FROM chats WHERE id IN (SELECT chat_id FROM memberships WHERE user_id = ?) AND "+
			"id IN (SELECT chat_id FROM memberships GROUP BY chat_id HAVING COUNT(*) = 2)", userId)
	if err != nil {
		log.Panic(err)
	}

	var chats []Chat
	for query.Next() {
		var chat Chat
		err = query.Scan(&chat.Id, &chat.LastMessageDate)
		if err != nil {
			log.Panic(err)
		}
		chats = append(chats, chat)
	}
	query.Close()

	return chats
}

// TODO: perhaps move it later to models.go. Not approved yet.
type PrivateChatOponent struct {
	Id                  int    // opoent id
	Name                string // oponent name
	ChatId              int    // chat id
	ChatLastMessageDate string // last message date, to sort chats by last message date
}

// TODO: remove it later. Not approved yet.
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
func (db DB) GetPrivateChatOponentsByUserId(userId int) []PrivateChatOponent {
	query, err := db.Query(
		"SELECT users.id, users.name, chats.id, chats.last_message_date FROM users "+
			"JOIN memberships ON memberships.user_id = users.id "+
			"JOIN chats ON chats.id = memberships.chat_id "+
			"WHERE users.id != ? AND chats.id IN (SELECT chat_id FROM memberships GROUP BY chat_id HAVING COUNT(*) = 2)",
		userId)
	if err != nil {
		log.Panic(err)
	}

	var oponents []PrivateChatOponent
	for query.Next() {
		var oponent PrivateChatOponent
		err = query.Scan(&oponent.Id, &oponent.Name, &oponent.ChatId, &oponent.ChatLastMessageDate)
		if err != nil {
			log.Panic(err)
		}
		oponents = append(oponents, oponent)
	}
	query.Close()

	return oponents
}

// AddMembership adds membership to database, returns id of new membership
func (db DB) AddMembership(chatId, userId int) int {
	result, err := db.Exec(
		"INSERT INTO memberships (chat_id, user_id, date) VALUES (?, ?, ?)",
		chatId, userId, time.Now().Format(time.RFC3339))
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
func (db DB) AddMessage(chatId, userId int, content string) int {
	timeNow := time.Now().Format(time.RFC3339)
	result, err := db.Exec(
		"INSERT INTO messages (chat_id, user_id, content, date) VALUES (?, ?, ?, ?)",
		chatId, userId, content, timeNow)
	if err != nil {
		log.Panic(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Panic(err)
	}
	_, err = db.Exec("UPDATE chats SET last_message_date = ? WHERE id = ?", timeNow, chatId)
	if err != nil {
		log.Panic(err)
	}
	return int(id)
}
