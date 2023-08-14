package database

import (
	. "backend/chat/database/migrations"
	"backend/common/integrations/auth"
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
	err := Migrations.Migrate(db.DB, Migrations.Latest())
	if err != nil {
		return err
	}

	groupChats := db.GetAllGroupChats()

	for _, chat := range groupChats {
		updatedMembers := auth.GetGroupMembers(chat.GroupId)
		db.UpdateChatMembers(chat.Id, *updatedMembers)
	}

	return nil
}
