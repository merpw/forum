package server

import (
	"database/sql"
	"forum/server"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestAuth tests auth routes
func TestAuth(t *testing.T) {
	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		t.Fatal(err)
	}
	srv := server.Connect(db)
	router := srv.Start()
	testServer := httptest.NewServer(router)
	defer testServer.Close()

	cli := testServer.Client()

	t.Run("signup", func(t *testing.T) {
		resp, err := cli.Post(testServer.URL+"/api/auth/signup", "application/json", nil)
		if err != nil {
			t.Fatal(err)
		}
		// TODO: add request body
		if resp.StatusCode != 200 {
			t.Fatalf("expected %d, got %d", 200, resp.StatusCode)
		}
	})

	t.Run("login", func(t *testing.T) {
		resp, err := cli.Post(testServer.URL+"/api/auth/login", "application/json", nil)
		if err != nil {
			t.Fatal(err)
		}
		// TODO: add request body
		if resp.StatusCode != 200 {
			t.Fatalf("expected %d, got %d", 200, resp.StatusCode)
		}
	})
}

// TestPost tests post routes
func TestPost(t *testing.T) {
	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		t.Fatal(err)
	}
	srv := server.Connect(db)
	router := srv.Start()
	testServer := httptest.NewServer(router)
	defer testServer.Close()

	cli := testServer.Client()

	tests := []struct {
		url string
	}{
		{"/api/posts/create"},
		{"/api/posts/1/like"},
		{"/api/posts/1/dislike"},
		{"/api/posts/1/comment"},

		{"/api/posts/1/comment/1/like"},
		{"/api/posts/1/comment/1/dislike"},
	}

	for _, test := range tests {
		t.Run(test.url, func(t *testing.T) {
			resp, err := cli.Post(testServer.URL+test.url, "application/json", nil)
			if err != nil {
				t.Fatal(err)
			}
			if resp.StatusCode != http.StatusUnauthorized {
				t.Fatalf("expected %d, got %d", http.StatusUnauthorized, resp.StatusCode)
			}
		})
	}
}
