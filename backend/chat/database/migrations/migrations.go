package migrations

import (
	"backend/migrate"
)

// Migrations define all available Migrations and their order
var Migrations = migrate.Migrations{
	v001,
	v002,
}
