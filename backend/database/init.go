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
						likes_count INTEGER,
						dislikes_count INTEGER,
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
						likes_count INTEGER,
						dislikes_count INTEGER,
    	   				FOREIGN KEY(post) REFERENCES posts(id),
						FOREIGN KEY(author) REFERENCES users(id))`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS sessions (
						id INTEGER PRIMARY KEY,
						token TEXT,
						expire INTEGER,
						user_id INTEGER,
						FOREIGN KEY(user_id) REFERENCES users(id))`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS likes (
		id INTEGER PRIMARY KEY,
		parent_id INTEGER,
		author_id INTEGER,
		post_like INTEGER,
		value INTEGER)`)
	if err != nil {
		return err
	}

	return nil
}
