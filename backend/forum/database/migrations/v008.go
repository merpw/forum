package migrations

import (
	"backend/migrate"
	"database/sql"
)

// v007 adds privacy to user table
var v008 = migrate.Migration{
	Up: func(db *sql.DB) error {
		_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS invitations (
		    id INTEGER PRIMARY KEY AUTOINCREMENT,
		    type INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			associated_id INTEGER NOT NULL,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
		    FOREIGN KEY (user_id) REFERENCES users(id)            
		);
`)
		if err != nil {
			return err
		}
		return nil
	},

	Down: func(db *sql.DB) error {

		_, err := db.Exec(`
		DROP TABLE IF EXISTS followers;
		VACUUM;
`)
		if err != nil {
			return err
		}
		return nil
	},
}
