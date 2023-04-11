package migrations

import "database/sql"

// v002 makes all fields in all tables NOT NULL
var v002 = Migration{
	Up: func(db *sql.DB) error {
		_, err := db.Exec("PRAGMA writable_schema = 1")
		if err != nil {
			return err
		}
		_, err = db.Exec("UPDATE sqlite_master SET sql = REPLACE(sql, 'TEXT', 'TEXT NOT NULL')")
		if err != nil {
			return err
		}
		_, err = db.Exec("UPDATE sqlite_master SET sql = REPLACE(sql, 'INTEGER', 'INTEGER NOT NULL')")
		if err != nil {
			return err
		}
		_, err = db.Exec("UPDATE sqlite_master SET sql = REPLACE(sql, 'NOT NULL PRIMARY KEY', 'PRIMARY KEY')")
		if err != nil {
			return err
		}
		_, err = db.Exec("UPDATE sqlite_master SET sql = REPLACE(sql, 'NOT NULL NOT NULL', 'NOT NULL')")
		if err != nil {
			return err
		}
		_, err = db.Exec("PRAGMA writable_schema = 0")
		if err != nil {
			return err
		}
		_, err = db.Exec("VACUUM")
		if err != nil {
			return err
		}
		return nil
	},

	Down: func(db *sql.DB) error {
		_, err := db.Exec("PRAGMA writable_schema = 1")
		if err != nil {
			return err
		}
		_, err = db.Exec("UPDATE sqlite_master SET sql = REPLACE(sql, 'TEXT NOT NULL', 'TEXT')")
		if err != nil {
			return err
		}
		_, err = db.Exec("UPDATE sqlite_master SET sql = REPLACE(sql, 'INTEGER NOT NULL', 'INTEGER')")
		if err != nil {
			return err
		}
		_, err = db.Exec("PRAGMA writable_schema = 0")
		if err != nil {
			return err
		}
		_, err = db.Exec("VACUUM")
		if err != nil {
			return err
		}
		return nil
	},
}
