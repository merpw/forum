package migrations

import (
	"backend/migrate"
	"database/sql"
	"fmt"
)

// v001 is the initial revision. It creates the tables users, posts, comments, sessions and reactions.
var v001 = migrate.Migration{
	Up: func(db *sql.DB) error {
		_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY,
			name TEXT,
			email TEXT,
			password TEXT
		);
		CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY,
			title TEXT,
			content TEXT,
			author INTEGER,
			date TEXT,
			likes_count INTEGER,
			dislikes_count INTEGER,
			comments_count INTEGER,
			categories TEXT,
			FOREIGN KEY(author) REFERENCES users(id)
		);
		CREATE TABLE IF NOT EXISTS comments (
			id INTEGER PRIMARY KEY,
			post_id INTEGER,
			author_id INTEGER,					
			content TEXT,
			date TEXT,
			likes_count INTEGER,
			dislikes_count INTEGER,
			FOREIGN KEY(post_id) REFERENCES posts(id),
			FOREIGN KEY(author_id) REFERENCES users(id)
		);
		CREATE TABLE IF NOT EXISTS sessions (
			id INTEGER PRIMARY KEY,
			token TEXT,
			expire INTEGER,
			user_id INTEGER,
			FOREIGN KEY(user_id) REFERENCES users(id)
        );
		CREATE TABLE IF NOT EXISTS post_reactions (
			id INTEGER PRIMARY KEY,
			post_id INTEGER,
			author_id INTEGER,
			reaction INTEGER,
			FOREIGN KEY(author_id) REFERENCES users(id),
    		FOREIGN KEY(post_id) REFERENCES posts(id)
		);
		CREATE TABLE IF NOT EXISTS comment_reactions (
			id INTEGER PRIMARY KEY,
			comment_id INTEGER,
			author_id INTEGER,
			reaction INTEGER,
			FOREIGN KEY(author_id) REFERENCES users(id),
    		FOREIGN KEY(comment_id) REFERENCES comments(id)
		)`)
		if err != nil {
			return err
		}
		return nil
	},
	Down: func(db *sql.DB) error {
		return fmt.Errorf("cannot rollback initial revision")
	},
}
