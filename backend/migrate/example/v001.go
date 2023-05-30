package example

import (
	"backend/migrate"
	"database/sql"
	"fmt"
)

// v001 creates the cats table with id, name and age columns
var v001 = migrate.Migration{
	Up: func(db *sql.DB) error {
		_, err := db.Exec(`
		CREATE TABLE cats (
			id INTEGER PRIMARY KEY,
		    name TEXT,
		    age INTEGER
		);`)
		if err != nil {
			return err
		}
		return nil
	},
	Down: func(db *sql.DB) error {
		return fmt.Errorf("cannot rollback initial revision")
	},
}
