package migrations

import "database/sql"

type Migration struct {
	Up   func(db *sql.DB) error
	Down func(db *sql.DB) error
}

var Migrations = []Migration{
	v001,
	v002,
}
