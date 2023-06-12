package server

import (
	"backend/forum/handlers"
	"database/sql"
	"net/http/httptest"
	"os"
	"testing"
)

type TestServer struct {
	*httptest.Server
}

// NewTestServer creates a new test server with a temporary database
func NewTestServer(t testing.TB) TestServer {
	dbFile, err := os.CreateTemp(".", "test.db")
	if err != nil {
		t.Fatal(err)
	}

	db, err := sql.Open("sqlite3", dbFile.Name()+"?_foreign_keys=true")
	if err != nil {
		t.Fatal(err)
	}

	h := handlers.New(db)
	err = h.DB.InitDatabase()
	if err != nil {
		t.Fatal(err)
	}

	testServer := TestServer{httptest.NewServer(h.Handler())}

	t.Cleanup(func() {
		_ = db.Close()
		_ = dbFile.Close()
		testServer.Close()
		if !t.Failed() {
			_ = os.Remove(dbFile.Name())
		} else {
			t.Logf("Database file: %s", dbFile.Name())
		}
	})

	return testServer
}
