package migrate_test

import (
	"database/sql"
	"forum/database/migrate"
	"log"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func init() {
	err := os.Remove("./test.db")
	if err != nil {
		if !os.IsNotExist(err) {
			log.Fatal(err)
		}
	}
}

// TODO: improve tests

func TestMigrate(t *testing.T) {
	db, err := sql.Open("sqlite3", "./test.db?_foreign_keys=true")
	if err != nil {
		t.Fatal(err)
	}
	for revision := 1; revision <= migrate.LATEST; revision++ {
		err = migrate.Migrate(db, revision)
		if err != nil {
			t.Fatal(err)
		}
	}
}
