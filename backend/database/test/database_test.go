package database_test

import (
	"database/sql"
	"forum/database"
	"log"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

var DB *database.DB

func TestMain(m *testing.M) {
	tmpDB, err := os.CreateTemp(".", "test.db")
	if err != nil {
		log.Fatal(err)
	}
	db, err := sql.Open("sqlite3", tmpDB.Name()+"?_foreign_keys=true")
	if err != nil {
		log.Fatal(err)
	}

	DB = &database.DB{DB: db}

	// srv = server.Connect(db)
	err = DB.InitDatabase()
	if err != nil {
		log.Fatal(err)
	}

	code := m.Run()
	if code != 0 {
		os.Exit(code)
	}

	// Clean only if all tests are successful (code=0)
	DB.Close()
	tmpDB.Close()
	os.Remove(tmpDB.Name())
}
