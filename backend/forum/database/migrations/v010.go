package migrations

import (
	"backend/migrate"
	"database/sql"
)

// v010 adds groups and group_members tables
var v010 = migrate.Migration{
	Up: func(db *sql.DB) error {
		_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS groups (
		    id INTEGER PRIMARY KEY AUTOINCREMENT,
		    title TEXT NOT NULL,
		    description TEXT NOT NULL,
		    creator_id INTEGER NOT NULL,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
		    FOREIGN KEY (creator_id) REFERENCES users(id)
		);
		CREATE TABLE IF NOT EXISTS group_members (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			group_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (group_id) REFERENCES groups(id),
			FOREIGN KEY (user_id) REFERENCES users(id)
		);
		ALTER TABLE posts 
		ADD COLUMN group_id INTEGER REFERENCES groups(id) DEFAULT NULL;

		ALTER TABLE invitations
		ADD COLUMN associated_id INTEGER DEFAULT NULL
	`)
		if err != nil {
			return err
		}
		return nil
	},

	Down: func(db *sql.DB) error {

		_, err := db.Exec(`
		DROP TABLE IF EXISTS groups;
		DROP TABLE IF EXISTS group_members;
		ALTER TABLE posts DROP COLUMN group_id;
		ALTER TABLE invitations DROP COLUMN associated_id;
		VACUUM;
`)
		if err != nil {
			return err
		}
		return nil
	},
}
