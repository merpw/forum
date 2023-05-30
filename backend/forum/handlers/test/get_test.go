package server_test

import (
	"backend/forum/handlers"
	"database/sql"
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
	h := handlers.New(db)
	err = h.DB.InitDatabase()
	if err != nil {
		t.Fatal(err)
	}

	router := h.Handler()
	testServer := httptest.NewServer(router)
	defer testServer.Close()

	cli := testServer.Client()

	firstName := sql.NullString{String: "Steven", Valid: true}
	lastName := sql.NullString{String: "Smith", Valid: true}
	dob := sql.NullString{String: "2023-04-08", Valid: true}
	gender := sql.NullString{String: "male", Valid: true}

	userId := h.DB.AddUser("Steve", "steve@apple.com", "@@@l1sa@@@", firstName, lastName, dob, gender)
	h.DB.AddPost("test", "test", "test", userId, "facts")

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
