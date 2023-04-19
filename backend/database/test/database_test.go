package database_test

import (
	"database/sql"
	"forum/server"
	"log"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB
var srv *server.Server

func TestMain(m *testing.M) {
	tmpDB, err := os.CreateTemp(".", "test.db")
	if err != nil {
		log.Fatal(err)
	}
	db, err = sql.Open("sqlite3", tmpDB.Name()+"?_foreign_keys=true")
	if err != nil {
		log.Fatal(err)
	}

	srv = server.Connect(db)
	err = srv.DB.InitDatabase()
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
}
