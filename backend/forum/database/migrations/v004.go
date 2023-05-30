package migrations

import (
	"backend/migrate"
	"database/sql"
)

// v004 adds first_name, last_name, dob and gender to the users table
var v004 = migrate.Migration{
	Up: func(db *sql.DB) error {
		_, err := db.Exec(`
		ALTER TABLE users ADD COLUMN first_name TEXT;
		ALTER TABLE users ADD COLUMN last_name TEXT;
		ALTER TABLE users ADD COLUMN dob TEXT;
		ALTER TABLE users ADD COLUMN gender TEXT;
`)
		if err != nil {
			return err
		}
		return nil
	},

	Down: func(db *sql.DB) error {

		_, err := db.Exec(`
		ALTER TABLE users DROP COLUMN first_name;
		ALTER TABLE users DROP COLUMN last_name;
		ALTER TABLE users DROP COLUMN dob;
		ALTER TABLE users DROP COLUMN gender;
		VACUUM;
`)
		if err != nil {
			return err
		}
		return nil
	},
}
