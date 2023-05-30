package migrations

import (
	"backend/migrate"
	"database/sql"
)

// v003 adds a description field to the post table
//
// For already existing posts, the description is set to the first 200 characters of the content
var v003 = migrate.Migration{
	Up: func(db *sql.DB) error {
		_, err := db.Exec(`
		ALTER TABLE posts
			ADD COLUMN description TEXT NOT NULL default '';
		UPDATE posts
			SET description = SUBSTR(content, 0, 200)
			WHERE description = '';
		UPDATE posts
			SET description = description || '...'
			WHERE length(description) = 199;
		`)
		if err != nil {
			return err
		}
		return nil
	},

	Down: func(db *sql.DB) error {
		_, err := db.Exec(`
		ALTER TABLE posts
			DROP COLUMN description
		`)
		if err != nil {
			return err
		}
		return nil
	},
}
