package migrations

import (
	"backend/migrate"
	"database/sql"
)

// v002 makes all fields in all tables NOT NULL
var v002 = migrate.Migration{
	Up: func(db *sql.DB) error {
		_, err := db.Exec(`
		PRAGMA writable_schema = 1;
		UPDATE sqlite_master SET sql = REPLACE(sql, 'TEXT', 'TEXT NOT NULL');
		UPDATE sqlite_master SET sql = REPLACE(sql, 'INTEGER', 'INTEGER NOT NULL');
		UPDATE sqlite_master SET sql = REPLACE(sql, 'NOT NULL PRIMARY KEY', 'PRIMARY KEY');
		UPDATE sqlite_master SET sql = REPLACE(sql, 'NOT NULL NOT NULL', 'NOT NULL');
		PRAGMA writable_schema = 0;
		VACUUM;
`)
		if err != nil {
			return err
		}
		return nil
	},

	Down: func(db *sql.DB) error {
		_, err := db.Exec(`
		PRAGMA writable_schema = 1;
		UPDATE sqlite_master SET sql = REPLACE(sql, 'TEXT NOT NULL', 'TEXT');
		UPDATE sqlite_master SET sql = REPLACE(sql, 'INTEGER NOT NULL', 'INTEGER');
		PRAGMA writable_schema = 0;
		VACUUM;
`)
		if err != nil {
			return err
		}
		return nil
	},
}
