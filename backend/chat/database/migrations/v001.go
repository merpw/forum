package migrations

import (
	"backend/migrate"
	"database/sql"
	"fmt"
)

var v001 = migrate.Migration{
	Up: func(db *sql.DB) error {
		_, err := db.Exec(`
		CREATE TABLE chats (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			type INTEGER NOT NULL DEFAULT 0
			-- 0 - private, 1 - group      
		);
		CREATE TABLE memberships (
		    id INTEGER PRIMARY KEY AUTOINCREMENT,
			chat_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
		    FOREIGN KEY (chat_id) REFERENCES chats(id)                    
		);
		CREATE TABLE messages (
		    id INTEGER PRIMARY KEY AUTOINCREMENT,
			chat_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			content TEXT NOT NULL,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (chat_id) REFERENCES chats(id)		                      
		);
		`)
		if err != nil {
			return err
		}
		return nil
	},
	Down: func(db *sql.DB) error {
		return fmt.Errorf("cannot rollback initial revision")
	},
}
