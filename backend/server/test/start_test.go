package server

import (
	"database/sql"
	"forum/server"
	"net/http/httptest"
	"testing"
)

func TestStart(t *testing.T) {
	db, err := sql.Open("sqlite3", "./test.db?_foreign_keys=true")
	if err != nil {
		t.Fatal(err)
	}
	srv := server.Connect(db)
	router := srv.Start()
	testServer := httptest.NewServer(router)
	defer testServer.Close()
}
