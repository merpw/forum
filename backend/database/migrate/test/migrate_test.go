package migrate_test

import (
	"database/sql"
	"fmt"
	"forum/database/migrate"
	"log"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func TestMain(m *testing.M) {
	tmpDB, err := os.CreateTemp(".", "test.db")
	if err != nil {
		log.Fatal(err)
	}
	db, err = sql.Open("sqlite3", tmpDB.Name()+"?_foreign_keys=true")
	if err != nil {
		log.Fatal(err)
	}

	dummyStdin, err := os.CreateTemp(".", "stdin.txt")
	if err != nil {
		log.Fatal(err)
	}
	os.Stdin = dummyStdin

	_, err = fmt.Fprintln(os.Stdin, "YES")
	if err != nil {
		log.Fatal(err)
	}

	code := m.Run()
	if code != 0 {
		os.Exit(code)
	}

	// Clean only if all tests are successful (code=0)
	db.Close()
	tmpDB.Close()
	os.Remove(tmpDB.Name())

	dummyStdin.Close()
	os.Remove(dummyStdin.Name())
}

func TestMigrate(t *testing.T) {

	t.Run("Empty database", func(t *testing.T) {
		err := migrate.Migrate(db, migrate.LATEST)
		if err != nil {
			t.Fatal(err)
		}
	})

	// Mock user input

	t.Run("Down, without YES", func(t *testing.T) {
		err := migrate.Migrate(db, 1)
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("Down, single call", func(t *testing.T) {
		inputYES()
		err := migrate.Migrate(db, 1)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Up, call for each revision", func(t *testing.T) {
		for revision := 1; revision <= migrate.LATEST; revision++ {
			err := migrate.Migrate(db, revision)
			if err != nil {
				t.Fatal(err)
			}
		}
	})

	t.Run("Down, call for each revision", func(t *testing.T) {
		for i := migrate.LATEST; i > 0; i-- {
			inputYES()
			err := migrate.Migrate(db, i)
			if err != nil {
				t.Fatal(err)
			}
		}
	})

	t.Run("Error, revision 0", func(t *testing.T) {
		err := migrate.Migrate(db, 0)
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("Error, revision -10", func(t *testing.T) {
		err := migrate.Migrate(db, -10)
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("Error, revision 1000", func(t *testing.T) {
		err := migrate.Migrate(db, 1000)
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("Error, broken database", func(t *testing.T) {
		_, err := db.Exec(`PRAGMA user_version = -100`)
		if err != nil {
			t.Fatal(err)
		}
		err = migrate.Migrate(db, migrate.LATEST)
		if err == nil {
			t.Fatal("expected error")
		}
		_, err = db.Exec(`PRAGMA user_version = 1`)
		if err != nil {
			t.Fatal(err)
		}
	})
}

// inputYES moves cursor of dummyStdin to 0,0 to make it read `YES`
func inputYES() {
	_, err := os.Stdin.Seek(0, 0)
	if err != nil {
		log.Fatal(err)
	}
}
