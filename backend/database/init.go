package database

import (
	"database/sql"
	"forum/database/migrations"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sql.DB
}

const REVISION = 2

func (db DB) InitDatabase() error {
	return migrations.Migrate(db.DB, REVISION)
}
