package database

func (db DB) GetChatMembers(chatId int) []int {
	rows, err := db.Query(`
		SELECT user_id FROM memberships
		WHERE chat_id = ?`, chatId)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var userIds []int
	for rows.Next() {
		var userId int
		err := rows.Scan(&userId)
		if err != nil {
			panic(err)
		}
		userIds = append(userIds, userId)
	}
	return userIds
}

// IsChatMember returns true if the user is a member of a chat.
func (db DB) IsChatMember(chatId, userId int) bool {
	row := db.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM memberships
			WHERE user_id = ? AND chat_id = ?
		)`, userId, chatId)
	var exists int
	err := row.Scan(&exists)
	if err != nil {
		panic(err)
	}
	return exists == 1
}
