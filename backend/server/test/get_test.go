package server_test

import (
	"database/sql"
	"forum/database"
	"forum/server"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestGet tests all GET routes for valid status codes
func TestGet(t *testing.T) {
	// Opens sqlite3
	db, err := sql.Open("sqlite3", "./test.db?_foreign_keys=true")
	if err != nil {
		t.Fatal(err)
	}
	// connects server to DB, and initiates the DB
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

	longDescription := `Lorem ipsum dolor sit amet, consectetur adipiscing elit.
	Quisque euismod, nibh nec aliquam ultricies, velit diam aliquet nunc, eget
	lobortis diam diam vitae velit. Donec euismod, nisl eget aliquam
	ullamcorper, nisl nisl aliquet nunc, eget lobortis diam diam vitae velit.`

	userId := srv.DB.AddUser("Steve", "steve@apple.com", "@@@l1sa@@@", firstName, lastName, dob, gender)
	srv.DB.AddPost("test", "test", userId, "fact", longDescription)

	tests := []struct {
		url          string
		expectedCode int
	}{
		{"/api/posts", 200},
		{"/api/posts/", 200},

		{"/api/posts/1", 200},

		{"/api/user/1", 200},
		{"/api/user/1/posts", 200},
		{"/api/me/posts/liked", 401}, // not logged in

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
		// Call GetPostComments to get comments for post 1
		comments := srv.DB.GetPostComments(1)

		// Check that the comments slice contains no comments
		if len(comments) != 0 {
			t.Errorf("Expected 0 comments, but got %d", len(comments))
		}

	})

}

// signup test user by sending a POST request to /api/signup

// test getting the posts liked by the user
// {"/api/me/posts/liked return 0 POSTS", "/api/me/posts/liked", nil, "[]"},

func TestQueries(t *testing.T) {
	// Create a test database with some sample data
	db, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	t.Run("GetPostComments", func(t *testing.T) {
		_, err = db.Exec(`
        CREATE TABLE posts (
            id INTEGER PRIMARY KEY,
            title TEXT,
            content TEXT
        );
        CREATE TABLE comments (
            id INTEGER PRIMARY KEY,
            post_id INTEGER,
            author_id INTEGER,
            content TEXT,
            date TEXT,
            likes_count INTEGER,
            dislikes_count INTEGER,
            FOREIGN KEY(post_id) REFERENCES posts(id)
        );
        INSERT INTO posts VALUES (1, 'Test Post', 'Lorem ipsum dolor sit amet.');
        INSERT INTO comments VALUES
            (1, 1, 1, 'Nice post!', '2023-03-30', 10, 0),
            (2, 1, 2, 'Thanks for sharing!', '2023-03-31', 5, 2),
            (3, 1, 3, 'I have a question...', '2023-04-01', 2, 3);
    `)
		if err != nil {
			t.Fatal(err)
		}

		// Call GetPostComments() with the ID of the test post
		dbInstance := database.DB{DB: db}
		comments := dbInstance.GetPostComments(1)

		// Verify that the returned comments match the expected values
		expectedComments := []database.Comment{
			{
				Id:            1,
				PostId:        1,
				AuthorId:      1,
				Content:       "Nice post!",
				Date:          "2023-03-30",
				LikesCount:    10,
				DislikesCount: 0,
			},
			{
				Id:            2,
				PostId:        1,
				AuthorId:      2,
				Content:       "Thanks for sharing!",
				Date:          "2023-03-31",
				LikesCount:    5,
				DislikesCount: 2,
			},
			{
				Id:            3,
				PostId:        1,
				AuthorId:      3,
				Content:       "I have a question...",
				Date:          "2023-04-01",
				LikesCount:    2,
				DislikesCount: 3,
			},
		}
		if len(comments) != len(expectedComments) {
			t.Fatalf("Expected %d comments, but got %d", len(expectedComments), len(comments))
		}
		for i := range expectedComments {
			if comments[i] != expectedComments[i] {
				t.Errorf("Expected comment %+v, but got %+v", expectedComments[i], comments[i])
			}
		}
	})

}

func TestRemoveExpiredSessions(t *testing.T) {
	// Open a test database
	db, err := sql.Open("sqlite3", "./test.db?_foreign_keys=true")
	if err != nil {
		t.Fatal(err)
	}

	// Create some expired sessions
	expireTime := time.Now().Add(-time.Hour).Unix()
	_, err = db.Exec("INSERT INTO sessions (token, expire, user_id) VALUES (?, ?, ?)", "expired_token_1", expireTime, 1)
	if err != nil {
		t.Fatal(err)
	}

	// Create a database instance and call the RemoveExpiredSessions function
	testDB := server.Connect(db)
	testDB.DB.RemoveExpiredSessions()

	// Check that the expired sessions have been removed from the database
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM sessions WHERE token IN (?)", "expired_token_1").Scan(&count)
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Errorf("Expected 0 sessions with expired tokens, but found %d", count)
	}
}
