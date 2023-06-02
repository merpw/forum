package database

import (
	"database/sql"
	"errors"
)

// GetUserChats returns a slice of all user's chats ids.
func (db DB) GetUserChats(userId int) (chats []int) {
	q, err := db.Query("SELECT chat_id FROM memberships WHERE user_id = ?", userId)
	if err != nil {
		panic(err)
	}
	defer q.Close()

	chats = make([]int, 0)
	for q.Next() {
		var chatId int
		err = q.Scan(&chatId)
		if err != nil {
			panic(err)
		}
		chats = append(chats, chatId)
	}
	return chats
}

type ChatData struct {
	LastMessageId int
	CompanionId   int
}

func (db DB) GetChatData(userId, chatId int) ChatData {
	row := db.QueryRow(`
		SELECT messages.id, memberships.user_id
		FROM messages 
		INNER JOIN memberships 
		ON memberships.chat_id = messages.chat_id
		WHERE messages.chat_id = ? AND memberships.user_id != ?
		ORDER BY messages.id DESC
		LIMIT 1
`, chatId, userId)

	var data ChatData
	err := row.Scan(&data.LastMessageId, &data.CompanionId)
	if err != nil {
		panic(err)
	}

	return data
}

// GetUsersChat returns an id of the chat between two users.
//
// If there is no such chat, it returns -1.
func (db DB) GetUsersChat(user1Id, user2Id int) int {

	if user1Id == user2Id {
		// TODO: maybe allow self-chats
		return -1
	}

	row := db.QueryRow(`
		SELECT chat_id FROM memberships
		WHERE user_id = ? AND chat_id IN (
			SELECT chat_id FROM memberships
			WHERE user_id = ?
		)`, user1Id, user2Id)

	var chatId int
	err := row.Scan(&chatId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return -1
		}
		panic(err)
	}
	return chatId
}
