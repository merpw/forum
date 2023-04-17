package migrations

import (
	"database/sql"
)

var v003 = Migration{
	Up: func(db *sql.DB) error {
		_, err := db.Exec(`

		ALTER TABLE users ADD COLUMN first_name TEXT  NOT NULL DEFAULT '';
		ALTER TABLE users ADD COLUMN last_name TEXT  NOT NULL DEFAULT '';
		ALTER TABLE users ADD COLUMN dob TEXT  NOT NULL DEFAULT '';
		ALTER TABLE users ADD COLUMN gender TEXT NOT NULL DEFAULT '';
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
