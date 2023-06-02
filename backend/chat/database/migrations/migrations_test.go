package migrations_test

import (
	. "backend/chat/database/migrations"
	"testing"
)

func TestMigrations(t *testing.T) {
	Migrations.Test(t)
}
