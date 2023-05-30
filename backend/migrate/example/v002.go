package example

import (
	"backend/migrate"
	"database/sql"
)

// v002 adds breed column to cats table
var v002 = migrate.Migration{
	Up: func(db *sql.DB) error {
		_, err := db.Exec(`
		ALTER TABLE cats
			ADD COLUMN breed TEXT NOT NULL default '';
		`)
		if err != nil {
			return err
		}
		return nil
	},
	Down: func(db *sql.DB) error {
		_, err := db.Exec(`
		ALTER TABLE cats
			DROP COLUMN breed
		`)
		if err != nil {
			return err
		}
		return nil
	},
}
