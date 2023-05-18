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

var cookie *http.Cookie

func getInvalidPosts() []Post {
	return []Post{
		{
			Title:      "",
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
			Content:    "",
			Categories: []string{"facts"},
		},
		{
			Title:      "Valid Title",
			Content:    "Valid title",
			Categories: []string{"Invalid category"},
		},
	}
}

func getInvalidUsers() []TestUser {
	return []TestUser{
		{
			Name:     "test1",
			Email:    "test@test.com",
			Password: "SuperAmazingPassword()!@*#)(!@#",
		},
		{
			Name:     "",
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
			Email:    "😂😂😂😂😂😂😂😂😂😂😂",
			Password: "ValidPassword123",
		},
		{
			Name:     "ValidName",
			Email:    "valid@test.com",
			Password: "1234",
		},
		{
			Name:     "Steve",
			Email:    "valid@test.com",
			Password: "12345678",
		},
		{
			Name:     "ValidName",
			Email:    "steve@apple.com",
			Password: "12345678",
		},
		{
			Name:     "  WsName  ",
			Email:    "valid@test.com",
			Password: "ValidPassword123",
		},
	}
}

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
