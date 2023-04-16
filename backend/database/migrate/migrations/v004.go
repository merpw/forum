package migrations

import "database/sql"

/*
create table "chats":

	type:
		2 (private 1vs1) or
		1 (group chat) or
		0 (the channel owner is posting to subscribers)
	date - date of creation

create table "memberships":

	user_id - fk to users.id
	chat_id - fk to chats.id
	date - date of joining

create table "messages":

	user_id - fk to users.id
	chat_id - fk to chats.id
	content - text of message, includes links to images (not approved yet)
	date - date of creation
*/
var v004 = Migration{
	Up: func(db *sql.DB) error {
		_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS chats (
			id INTEGER PRIMARY KEY,
			type INTEGER,
			date TEXT);
		CREATE TABLE IF NOT EXISTS memberships (
			id INTEGER PRIMARY KEY,
			user_id INTEGER,
			chat_id INTEGER,
			date TEXT,
			FOREIGN KEY(user_id) REFERENCES users(id),
			FOREIGN KEY(chat_id) REFERENCES chats(id));
		CREATE TABLE IF NOT EXISTS messages (
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
		DROP TABLE IF EXISTS chats;
		DROP TABLE IF EXISTS memberships;
		DROP TABLE IF EXISTS messages;
		`)
		if err != nil {
			return err
		}
		return nil
	},
}
