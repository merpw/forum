package migrations

import "database/sql"

// v003 - create tables: chats, memberships, messages
var v003 = Migration{
	Up: func(db *sql.DB) error {
		_, err := db.Exec(`
		CREATE TABLE chats (
			id INTEGER PRIMARY KEY,
			type INTEGER,
			date TEXT);
		CREATE TABLE memberships (
			id INTEGER PRIMARY KEY,
			user_id INTEGER,
			chat_id INTEGER,
			date TEXT,
			FOREIGN KEY(user_id) REFERENCES users(id),
			FOREIGN KEY(chat_id) REFERENCES chats(id));
		CREATE TABLE messages (
			id INTEGER PRIMARY KEY,
			user_id INTEGER,
			chat_id INTEGER,
			content TEXT,
			date TEXT,
			FOREIGN KEY(user_id) REFERENCES users(id),
			FOREIGN KEY(chat_id) REFERENCES chats(id));
		`)
		if err != nil {
			return err
		}
		return nil
	},
	Down: func(db *sql.DB) error {
		_, err := db.Exec(`
		DROP TABLE chats;
		DROP TABLE memberships;
		DROP TABLE messages;
		`)
		if err != nil {
			return err
		}
		return nil
	},
}
