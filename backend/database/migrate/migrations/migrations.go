package migrations

import "database/sql"

// Migration defines operations that should be done to move Up or Down between revisions
type Migration struct {
	Up   func(db *sql.DB) error
	Down func(db *sql.DB) error
}

// Migrations defines all available Migrations and their order
var Migrations = []Migration{
	v001,
	v002,
	v004,
}
