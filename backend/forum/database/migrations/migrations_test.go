package migrations_test

import (
	. "backend/forum/database/migrations"
	"testing"
)

func TestMigrations(t *testing.T) {
	Migrations.Test(t)
}
