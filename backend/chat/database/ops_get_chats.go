package database

import (
	"database/sql"
	"errors"
)

// GetUserChats returns a slice of all user's chats ids.
//
// Sorted by the last message's id in each chat.
func (db DB) GetUserChats(userId int) (chats []int) {
	q, err := db.Query(`
		SELECT chat_id FROM memberships WHERE user_id = ? 
		ORDER BY (SELECT id FROM messages WHERE chat_id = memberships.chat_id ORDER BY id DESC LIMIT 1) DESC
`, userId)
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
	CompanionId   *int
	GroupId       *int
}

// GetChatData returns data about a chat.
func (db DB) GetChatData(userId, chatId int) ChatData {
	row := db.QueryRow(`
SELECT last_message_id, companion_id, group_id
FROM (SELECT id as last_message_id
      FROM messages
      WHERE chat_id = :chatId
      ORDER BY id DESC
      LIMIT 1),
     (SELECT CASE
         -- Private chat -- 
                 WHEN type = 0
                     THEN (SELECT user_id
                           FROM memberships
                           WHERE chat_id = :chatId
                             and user_id != :userId
                           LIMIT 1)
                 ELSE (SELECT null)
                 END as companion_id
      FROM chats
      WHERE id = :chatId),
     (SELECT CASE
         -- Group chat --
                 WHEN type = 1
                     THEN (SELECT group_id
                           FROM chats
                           WHERE id = :chatId
                           LIMIT 1)
                 ELSE (SELECT null)
                 END as group_id
      FROM chats
      WHERE id = :chatId)
`, sql.Named("chatId", chatId), sql.Named("userId", userId))

	var data ChatData
	err := row.Scan(&data.LastMessageId, &data.CompanionId, &data.GroupId)
	if err != nil {
		panic(err)
	}

	return data
}

// GetPrivateChat returns an id of the chat between two users.
//
// If there is no such chat, it returns -1.
func (db DB) GetPrivateChat(user1Id, user2Id int) (chatId *int) {

	if user1Id == user2Id {
		// TODO: maybe allow self-chats
		return nil
	}

	row := db.QueryRow(`
		SELECT chat_id FROM memberships
		WHERE user_id = ? AND chat_id IN (
			SELECT chat_id FROM memberships
			WHERE user_id = ?
		)`, user1Id, user2Id)

	err := row.Scan(&chatId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		panic(err)
	}
	return chatId
}

func (db DB) GetGroupChat(groupId int) (chatId *int) {
	row := db.QueryRow(`
		SELECT id FROM chats WHERE group_id = ?
`, groupId)

	err := row.Scan(&chatId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		panic(err)
	}
	return chatId
}

type GroupChat struct {
	Id      int
	GroupId int
}

func (db DB) GetAllGroupChats() (chats []GroupChat) {
	q, err := db.Query(`
		SELECT id,group_id FROM chats WHERE type = 1
`)
	if err != nil {
		panic(err)
	}
	defer q.Close()

	chats = make([]GroupChat, 0)
	for q.Next() {
		var chat GroupChat
		err = q.Scan(&chat.Id, &chat.GroupId)
		if err != nil {
			panic(err)
		}
		chats = append(chats, chat)
	}
	return chats
}
