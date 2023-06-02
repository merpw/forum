package database

import (
	. "backend/chat/database/migrations"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sql.DB
}

func New(db *sql.DB) *DB {
	return &DB{
		DB: db,
	}
}

func (db DB) InitDatabase() error {
	return Migrations.Migrate(db.DB, Migrations.Latest())
}
