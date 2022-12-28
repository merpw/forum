package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sql.DB
}

// InitDatabase creates all the necessary tables in sql.DB if they don't exist
func (db DB) InitDatabase() error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY,
		name TEXT,
		email TEXT,
		password TEXT)`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS posts (
						id INTEGER PRIMARY KEY,
						title TEXT,
						content TEXT,
						author INTEGER,
						date TEXT,
						likes INTEGER,
						dislikes INTEGER,
						user_reactions TEXT,
						comments_count INTEGER,
						FOREIGN KEY(author) REFERENCES users(id))`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS comments (
						id INTEGER PRIMARY KEY,
						post INTEGER,
						author INTEGER,					
						text TEXT,
						date TEXT,
						likes INTEGER,
						dislikes INTEGER,
    	   				FOREIGN KEY(post) REFERENCES posts(id),
						FOREIGN KEY(author) REFERENCES users(id))`)
	if err != nil {
		return err
	}

	return nil
}
