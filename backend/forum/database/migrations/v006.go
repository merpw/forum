package migrations

import (
	"backend/migrate"
	"database/sql"
)

// v006 adds 'avatar' and 'bio' columns to 'users' table
var v006 = migrate.Migration{
	Up: func(db *sql.DB) error {
		_, err := db.Exec(`
		ALTER TABLE users ADD COLUMN avatar TEXT;
		ALTER TABLE users ADD COLUMN bio TEXT;
`)
		if err != nil {
			return err
		}
		return nil
	},

	Down: func(db *sql.DB) error {

		_, err := db.Exec(`
		ALTER TABLE users DROP COLUMN avatar;
		ALTER TABLE users DROP COLUMN bio;
		VACUUM;
`)
		if err != nil {
			return err
		}
		return nil
	},
}
