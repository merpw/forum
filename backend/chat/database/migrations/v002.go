package migrations

import (
	"backend/migrate"
	"database/sql"
)

// v002 adds group_id column to 'chats' table
var v002 = migrate.Migration{
	Up: func(db *sql.DB) error {
		_, err := db.Exec(`
		ALTER TABLE chats ADD COLUMN group_id INTEGER;	
		`)
		if err != nil {
			return err
		}
		return nil
	},
	Down: func(db *sql.DB) error {
		_, err := db.Exec(`
		ALTER TABLE chats DROP COLUMN group_id;	
		`)
		if err != nil {
			return err
		}
		return nil
	},
}
