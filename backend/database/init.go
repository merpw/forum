package database

import (
	"database/sql"
	"forum/database/migrate"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sql.DB
}

const REVISION = 2

func (db DB) InitDatabase() error {
	return migrate.Migrate(db.DB, REVISION)
}
