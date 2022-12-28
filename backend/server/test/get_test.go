package server

import (
	"database/sql"
	"forum/database"
	"forum/server"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestGet tests all GET routes for valid status codes
func TestGet(t *testing.T) {
	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		t.Fatal(err)
	}
	srv := server.Connect(db)
	err = srv.DB.InitDatabase()
	if err != nil {
		t.Fatal(err)
	}

	router := srv.Start()
	testServer := httptest.NewServer(router)
	defer testServer.Close()

	cli := testServer.Client()

	srv.DB.AddPost(database.Post{Title: "test", Content: "test"})
	srv.DB.AddUser(database.User{Name: "Steve", Email: "steve@apple.com", Password: "@@@l1sa@@@"})

	tests := []struct {
		url          string
		expectedCode int
	}{
		{"/api/posts", 200},
		{"/api/posts/", 200},

		{"/api/posts/1", 200},

		{"/api/user/1", 200},
		{"/api/user/1/posts", 200},

		{"/api/posts/categories", 200},
		{"/api/posts/categories/rumors", 200},

		{"/api/posts/-1", 404},
		{"/api/posts/cat", 404},

		{"/api/user/-1", 404},
		{"/api/user/cat", 404},
		{"/api/user/cat/posts", 404},
		{"/api/user/", 404},

		{"/api/posts/categories/cat", 404},

		{"/cat/", 404},
		{"/api/cat/", 404},
		{"/api/", 404},
		{"/", 404},

		{"/api/me", http.StatusUnauthorized},
		{"/api/me/posts", http.StatusUnauthorized},

		{"/api/posts/create", http.StatusMethodNotAllowed},
		{"/api/posts/1/like", http.StatusMethodNotAllowed},
		{"/api/posts/1/dislike", http.StatusMethodNotAllowed},
		{"/api/posts/1/comment", http.StatusMethodNotAllowed},

		{"/api/auth/login", http.StatusMethodNotAllowed},
		{"/api/auth/signup", http.StatusMethodNotAllowed},
		{"/api/auth/logout", http.StatusMethodNotAllowed},
	}
	for _, test := range tests {
		t.Run(test.url, func(t *testing.T) {
			resp, err := cli.Get(testServer.URL + test.url)
			if err != nil {
				t.Fatal(err)
			}
			if resp.StatusCode != test.expectedCode {
				t.Fatalf("expected %d, got %d", test.expectedCode, resp.StatusCode)
			}
		})
	}
}
