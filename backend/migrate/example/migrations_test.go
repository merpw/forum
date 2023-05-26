package example_test

import (
	"backend/migrate/example"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestMigrations(t *testing.T) {
	example.Migrations.Test(t)
}
