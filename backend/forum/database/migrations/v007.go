package migrations

import (
	"backend/migrate"
	"database/sql"
)

// v007 adds privacy to user table
var v007 = migrate.Migration{
	Up: func(db *sql.DB) error {
		_, err := db.Exec(`
		CREATE TABLE followers (
		    id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			follower_id INTEGER NOT NULL,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
		    FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (follower_id) REFERENCES users(id)              
		);
		ALTER TABLE users 
		    ADD COLUMN privacy INTEGER DEFAULT 1;
		UPDATE users
			SET privacy = 1 WHERE privacy IS NULL;
`)
		if err != nil {
			return err
		}
		return nil
	},

	Down: func(db *sql.DB) error {

		_, err := db.Exec(`
		ALTER TABLE users DROP COLUMN privacy;
		DROP TABLE followers;
		VACUUM;
`)
		if err != nil {
			return err
		}
		return nil
	},
}
