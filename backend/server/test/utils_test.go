package server_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"forum/server"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestUser is struct to create testing users for handlerfuncs
type TestUser struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	DoB       string `json:"dob"`
	Gender    string `json:"gender"`
}

type Post struct {
	Title       string
	Content     string
	AuthorId    int
	Categories  []string
	Description string
}

var validUser TestUser = TestUser{
	Name:      "test",
	Email:     "test@test.com",
	Password:  "SuperAmazingPassword()!@*#)(!@#",
	FirstName: "John",
	LastName:  "Doe",
	DoB:       "2000-01-01",
	Gender:    "male",
}

// dummyLogin is a helper function to login a user and return a cookie
func dummyLogin(t *testing.T, cli *http.Client, testServer *httptest.Server, testUser TestUser) *http.Cookie {
	body := fmt.Sprintf(`{ "login": "%v", "password": "%v" }`, testUser.Email, testUser.Password)
	resp, err := cli.Post(testServer.URL+"/api/login", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Cookies()) == 0 {
		t.Fatal("no cookies after login")
	}
	cookie := resp.Cookies()[0]

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected %d, got %d", 200, resp.StatusCode)
	}
	return cookie
}

// cookie to simulate logged in user.
var cookie *http.Cookie

// getInValidPosts returns a slice of invalid posts, that should not be accepted by the server.
func getInvalidPosts() []Post {
	return []Post{
		{
			Title:      "", // invalid title (empty)
			Content:    "Valid Content",
			Categories: []string{"facts"},
		},
		{
			Title:      "invalidTitleTooLongItWillExceed25length",
			Content:    "Valid Content",
			Categories: []string{"facts"},
		},
		{
			Title:      "Valid Title",
			Content:    "", // invalid content (empty)
			Categories: []string{"facts"},
		},
		{
			Title:      "Valid Title",
			Content:    "Valid title",
			Categories: []string{"Invalid category"}, // invalid category (empty list)
		},
	}
}

// getInvalidUsers returns a slice of invalid users, that should not be accepted by the server.
func getInvalidUsers() []TestUser {
	return []TestUser{
		{
			Name:     "test1",
			Email:    "test@test.com", // email already in use
			Password: "SuperAmazingPassword()!@*#)(!@#",
		},
		{
			Name:     "", // invalid name (empty)
			Email:    "",
			Password: "",
		},
		{
			Name:     "ThisUserNameIsWayTooLong",
			Email:    "valid@test.com",
			Password: "ValidPassword123",
		},
		{
			Name:     "ValidName",
			Email:    "😂😂😂😂😂😂😂😂😂😂😂", // invalid email
			Password: "ValidPassword123",
		},
		{
			Name:     "ValidName",
			Email:    "valid@test.com",
			Password: "1234", // invalid password: too short
		},
		{
			Name:     "Steve", // Invalid name: already in use
			Email:    "valid@test.com",
			Password: "12345678",
		},
		{
			Name:     "ValidName",
			Email:    "steve@apple.com", // Invalid e-mail: Already in use
			Password: "12345678",
		},
		{
			Name:     "  WsName  ", // Invalid name: Contains leading or trailing whitespace
			Email:    "valid@test.com",
			Password: "ValidPassword123",
		},
	}
}

// setupServer sets up a test server with a SQLite database and returns the
// necessary components for testing.
func setupServer() (*sql.DB, *server.Server, *httptest.Server, *http.Client, error) {
	db, err := sql.Open("sqlite3", "./test.db?_foreign_keys=true")
	if err != nil {
		return nil, nil, nil, nil, err
	}

	srv := server.Connect(db)
	err = srv.DB.InitDatabase()
	if err != nil {
		return nil, nil, nil, nil, err
	}

	router := srv.Start()
	testServer := httptest.NewServer(router)
	cli := testServer.Client()

	return db, srv, testServer, cli, nil
}

// signup function sends a POST request to a test server's signup API endpoint with a
// JSON payload representing a user and returns the response or an error.
func signup(cli *http.Client, testServer *httptest.Server, user TestUser) (*http.Response, error) {
	body, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	resp, err := cli.Post(testServer.URL+"/api/signup", "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// createPost creates a new post by sending a POST request to a specified URL
// with a JSON request body and a cookie.
// err := createPost(cli, testServer, cookie, title, content, []string{"facts"})
func createPost(cli *http.Client, sURL string, ck *http.Cookie, p Post) (*http.Response, error) {
	requestBodyBytes, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, sURL+"/api/posts/create/", bytes.NewReader(requestBodyBytes))
	if err != nil {
		return nil, err
	}
	req.AddCookie(ck)

	return cli.Do(req)
}
