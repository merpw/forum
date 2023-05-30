package server_test

import (
	"backend/forum/handlers"
	"database/sql"
	"log"
	"net/http/httptest"
	"os"
	"testing"
)

var testServer *httptest.Server

// TestMain is the entry point for all tests
func TestMain(m *testing.M) {
	testServer = startServer()
	os.Exit(m.Run())
}

// startServer starts a test server with new test database file and returns a httptest.Server
func startServer() *httptest.Server {
	err := os.Remove("./test.db")
	if err != nil {
		if !os.IsNotExist(err) {
			log.Fatal(err)
		}
	}

	db, err := sql.Open("sqlite3", "./test.db?_foreign_keys=true")
	if err != nil {
		log.Fatal(err)
	}

	h := handlers.New(db)
	err = h.DB.InitDatabase()
	if err != nil {
		log.Fatal(err)
	}

	return httptest.NewServer(h.Handler())
}
