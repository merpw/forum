package migrations

import (
	"backend/migrate"
	"database/sql"
)

// v005 renames column "name" to "username" in table users
var v005 = migrate.Migration{
	Up: func(db *sql.DB) error {
		_, err := db.Exec(`
		ALTER TABLE users RENAME COLUMN name TO username;
		UPDATE users SET username = REPLACE(username, ' ', '_');
`)
		if err != nil {
			return err
		}
		return nil
	},

	Down: func(db *sql.DB) error {
		_, err := db.Exec(`
		ALTER TABLE users RENAME COLUMN username TO name;
`)
		if err != nil {
			return err
		}
		return nil
	},
}
