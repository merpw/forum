package migrate

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

// initDatabase creates a temporary database for testing and returns a clean function
func initDatabase() (db *sql.DB, clean func()) {
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

	clean = func() {
		// Clean only if all tests are successful (code=0)
		db.Close()
		tmpDB.Close()
		os.Remove(tmpDB.Name())

		dummyStdin.Close()
		os.Remove(dummyStdin.Name())
	}

	return db, clean
}

// Test initializes an empty temporary database
// and runs all migrations in order and then in reverse order
func (migrations Migrations) Test(t *testing.T) {
	db, clean := initDatabase()
	t.Run("Empty database", func(t *testing.T) {
		err := migrations.Migrate(db, migrations.Latest())
		if err != nil {
			t.Fatal(err)
		}
	})

	// Mock user input

	t.Run("Down, without YES", func(t *testing.T) {
		err := migrations.Migrate(db, 1)
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("Down, single call", func(t *testing.T) {
		inputYES()
		err := migrations.Migrate(db, 1)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Up, call for each revision", func(t *testing.T) {
		for revision := 1; revision <= migrations.Latest(); revision++ {
			err := migrations.Migrate(db, revision)
			if err != nil {
				t.Fatal(err)
			}
		}
	})

	t.Run("Down, call for each revision", func(t *testing.T) {
		for i := migrations.Latest(); i > 0; i-- {
			inputYES()
			err := migrations.Migrate(db, i)
			if err != nil {
				t.Fatal(err)
			}
		}
	})

	t.Run("Error, revision 0", func(t *testing.T) {
		err := migrations.Migrate(db, 0)
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("Error, revision -10", func(t *testing.T) {
		err := migrations.Migrate(db, -10)
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("Error, revision 1000", func(t *testing.T) {
		err := migrations.Migrate(db, 1000)
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("Error, broken database", func(t *testing.T) {
		_, err := db.Exec(`PRAGMA user_version = -100`)
		if err != nil {
			t.Fatal(err)
		}
		err = migrations.Migrate(db, migrations.Latest())
		if err == nil {
			t.Fatal("expected error")
		}
		_, err = db.Exec(`PRAGMA user_version = 1`)
		if err != nil {
			t.Fatal(err)
		}
	})

	if !t.Failed() {
		clean()
	}
}

// inputYES moves cursor of dummyStdin to 0,0 to make it read `YES`
func inputYES() {
	_, err := os.Stdin.Seek(0, 0)
	if err != nil {
		log.Fatal(err)
	}
}
