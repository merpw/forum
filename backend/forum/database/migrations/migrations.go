package migrations

import (
	"backend/migrate"
)

// Migrations defines all available Migrations and their order
var Migrations = migrate.Migrations{
	v001,
	v002,
	v003,
	v004,
	v005,
}
