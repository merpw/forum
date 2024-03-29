package migrations

import (
	"backend/migrate"
	"database/sql"
)

// v010 adds 'privacy' column to 'posts' table
var v010 = migrate.Migration{
	Up: func(db *sql.DB) error {
		_, err := db.Exec(`
		ALTER TABLE posts ADD COLUMN privacy INT;
		UPDATE posts
			SET privacy = 0 WHERE privacy IS NULL;
		CREATE TABLE IF NOT EXISTS post_audience (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			post_id INTEGER NOT NULL,
			follow_id INTEGER NOT NULL,
			FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
			FOREIGN KEY (follow_id) REFERENCES followers(id) ON DELETE CASCADE,
		    UNIQUE (post_id, follow_id)
		);
`)
		if err != nil {
			return err
		}
		return nil
	},
	Down: func(db *sql.DB) error {
		_, err := db.Exec(`
		ALTER TABLE posts DROP COLUMN privacy;
		DROP TABLE IF EXISTS post_audience;
		VACUUM;
`)
		if err != nil {
			return err
		}
		return nil
	},
}
