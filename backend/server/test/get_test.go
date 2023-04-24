package server_test

import (
	"database/sql"
	"forum/server"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestGet tests all GET routes for valid status codes
func TestGet(t *testing.T) {
	db, err := sql.Open("sqlite3", "./test.db?_foreign_keys=true")
	if err != nil {
		t.Fatal(err)
	}
	srv := server.Connect(db)
	err = srv.DB.InitDatabase()
	if err != nil {
		t.Fatal(err)
	}

	// Opens the available routes
	router := srv.Start()
	testServer := httptest.NewServer(router)
	defer testServer.Close()

	cli := testServer.Client()

	firstName := sql.NullString{String: "Steven", Valid: true}
	lastName := sql.NullString{String: "Smith", Valid: true}
	dob := sql.NullString{String: "2023-04-08", Valid: true}
	gender := sql.NullString{String: "male", Valid: true}

	userId := srv.DB.AddUser("Steve", "steve@apple.com", "@@@l1sa@@@", firstName, lastName, dob, gender)
	srv.DB.AddPost("test", "test", "test", userId, "facts")

	tests := []struct {
		url          string
		expectedCode int
	}{
		{"/api/posts", http.StatusOK},
		{"/api/posts/", http.StatusOK},

		{"/api/posts/1", http.StatusOK},

		{"/api/user/1", http.StatusOK},
		{"/api/user/1/posts", http.StatusOK},

		{"/api/posts/categories", http.StatusOK},
		{"/api/posts/categories/rumors", http.StatusOK},

		{"/api/posts/-1", http.StatusNotFound},
		{"/api/posts/cat", http.StatusNotFound},

		{"/api/user/-1", http.StatusNotFound},
		{"/api/user/cat", http.StatusNotFound},
		{"/api/user/cat/posts", http.StatusNotFound},
		{"/api/user/", http.StatusNotFound},

		{"/api/posts/categories/cat", http.StatusNotFound},

		{"/cat/", http.StatusNotFound},
		{"/api/cat/", http.StatusNotFound},
		{"/api/", http.StatusNotFound},
		{"/", http.StatusNotFound},

		{"/api/posts/0/comments", http.StatusNotFound},
		{"/api/posts/1/comment/1/reaction", http.StatusUnauthorized},

		{"/api/me", http.StatusUnauthorized},
		{"/api/me/posts", http.StatusUnauthorized},

		{"/api/me/posts/liked", http.StatusMethodNotAllowed},
		{"/api/posts/create", http.StatusMethodNotAllowed},
		{"/api/posts/1/like", http.StatusMethodNotAllowed},
		{"/api/posts/1/dislike", http.StatusMethodNotAllowed},
		{"/api/posts/1/comment", http.StatusMethodNotAllowed},

		{"/api/login", http.StatusMethodNotAllowed},
		{"/api/signup", http.StatusMethodNotAllowed},
		{"/api/logout", http.StatusMethodNotAllowed},
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
