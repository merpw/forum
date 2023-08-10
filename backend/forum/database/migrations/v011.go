package migrations

import (
	"backend/migrate"
	"database/sql"
)

var v011 = migrate.Migration{
	Up: func(db *sql.DB) error {
		_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS events (
		    id INTEGER PRIMARY KEY AUTOINCREMENT,
		    group_id INT NOT NULL,
			creator_id INT NOT NULL,
			title TEXT NOT NULL,
			description TEXT NOT NULL,
			time_and_date  TEXT NOT NULL,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (creator_id) REFERENCES users(id),
			FOREIGN KEY (group_id) REFERENCES groups(id)
		);
		CREATE TABLE IF NOT EXISTS event_members (
		    id INTEGER PRIMARY KEY AUTOINCREMENT,
		    event_id INTEGER NOT NULL,
		    user_id INTEGER NOT NULL,
		    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
		    FOREIGN KEY (event_id) REFERENCES events(id),
		    FOREIGN KEY (user_id) REFERENCES users(id),
		    CONSTRAINT unique_member UNIQUE(event_id, user_id)
		    );
	
	`)
		if err != nil {
			return err
		}
		return nil
	},

	Down: func(db *sql.DB) error {
		_, err := db.Exec(`
		DROP TABLE IF EXISTS events;
		DROP TABLE IF EXISTS event_members;
		VACUUM;
	`)
		if err != nil {
			return err
		}
		return nil
	},
}
