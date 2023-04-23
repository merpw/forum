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

	userId := srv.DB.AddUser("Steve", "steve@apple.com", "@@@l1sa@@@")
	srv.DB.AddPost("test", "test", userId, "facts")
	srv.DB.AddPost("test2", "test2", userId, "facts")

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

	// Possibly combine this test with the test below named "TestDatabaseQueries"
	t.Run("databaseQueries", func(t *testing.T) {
		comments := srv.DB.GetPostComments(1)

		if len(comments) != 0 {
			t.Errorf("Expected 0 comments, but got %d", len(comments))
		}

	})
}
