package server_test

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/gofrs/uuid"
)

// TestWithAuth tests all routes that require authentication
func TestWithAuth(t *testing.T) {
	db, srv, testServer, cli, err := setupServer()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	defer testServer.Close()

	// Create a user for db queries
	u := struct {
		Name   string
		Email  string
		Pass   string
		FName  sql.NullString
		LName  sql.NullString
		DoB    sql.NullString
		Gender sql.NullString
	}{
		Name:   "Steve",
		Email:  "steve@apple.com",
		Pass:   "@@@l1sa@@@",
		FName:  sql.NullString{String: "Steve", Valid: true},
		LName:  sql.NullString{String: "Jobs", Valid: true},
		DoB:    sql.NullString{String: "1955-02-24", Valid: true},
		Gender: sql.NullString{String: "male", Valid: true},
	}

	// Adds an user and a post to the database
	userId := srv.DB.AddUser(u.Name, u.Email, u.Pass, u.FName, u.LName, u.DoB, u.Gender)
	srv.DB.AddPost("test", "test", userId, "facts", "beatufiul, amazing, wonderful facts")

	// Slice of invalid users. It will cover most nonDB test cases.
	invalidUsers := getInvalidUsers()

	t.Run("signup", func(t *testing.T) {
		resp, err := signup(cli, testServer, validUser)
		if err != nil {
			t.Fatal(err)
		}

		// bug found, if the user is already present, we give back 400 instead of 409
		if resp.StatusCode != http.StatusOK {
			errBody, _ := io.ReadAll(resp.Body)
			fmt.Println(string(errBody))
			t.Fatalf("expected %d, got %d", http.StatusOK, resp.StatusCode)
		}
	})

	t.Run("invalidSignup", func(t *testing.T) {
		// signup with an invalid user (empry string)
		invalidResp, err := signup(cli, testServer, TestUser{})
		if err != nil {
			t.Fatal(err)
		}

		if invalidResp.StatusCode != http.StatusBadRequest {
			t.Fatalf("expected %d, got %d", 400, invalidResp.StatusCode)
		}

		// Loop through all the invalid users to test all the nonDB errors.
		for _, user := range invalidUsers {
			invalidResp, err := signup(cli, testServer, user)
			if err != nil {
				t.Fatal(err)
			}

			if invalidResp.StatusCode != http.StatusBadRequest {
				t.Fatalf("expected %d, got %d", http.StatusBadRequest, invalidResp.StatusCode)
			}

		}
	})

	t.Run("login", func(t *testing.T) {
		cookie = dummyLogin(t, cli, testServer, validUser)
		// TODO: Improve this test.
	})

	// TODO: Make this test not use unauthorized global exports.

	// create a test for a post not found (404) error
	t.Run("postNotFound", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, testServer.URL+"/api/posts/1234", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.AddCookie(cookie)

		resp, err := cli.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusNotFound {
			t.Fatalf("expected %d, got %d", 404, resp.StatusCode)
		}
	})

	// Invalid posts to be used in the "invalid posts" test
	invalidPosts := getInvalidPosts()

	// test create 5 posts
	t.Run("createPost", func(t *testing.T) {
		for i := 1; i <= 5; i++ {
			p := Post{
				Title:      fmt.Sprintf("Test Title %d", i),
				Content:    fmt.Sprintf("Test Content %d", i),
				Categories: []string{"facts"}, // Categories
			}
			resp, err := createPost(cli, testServer.URL, cookie, p)
			if err != nil {
				t.Fatal(err)
			}
			if resp.StatusCode != http.StatusOK {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					log.Println(err)
				}
				errorMessage := string(body)
				t.Fatalf("expected %d, got %d, ErrBody: %v", 200, resp.StatusCode, errorMessage)
			}
		}
	})

	// test revalidate func by using url route /post/id

	// TODO: DOCUMENTATION
	t.Run("createInvalidPost", func(t *testing.T) {
		resp, err := createPost(cli, testServer.URL, cookie, Post{"", "", 1, []string{""}, ""})
		if err != nil {
			t.Fatal(err)
		}

		if resp.StatusCode != http.StatusBadRequest {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Println(err)
			}
			errorMessage := string(body)
			t.Fatalf("expected %d, got %d, ErrBody: %v", 400, resp.StatusCode, errorMessage)
		}

		// Tests the invalid posts
		for _, post := range invalidPosts {
			resp, err := createPost(cli, testServer.URL, cookie, post)
			if err != nil {
				t.Fatal(err)
			}
			if resp.StatusCode != http.StatusBadRequest {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					log.Println(err)
				}
				errorMessage := string(body)
				t.Fatalf("expected %d, got %d, ErrBody: %v", 400, resp.StatusCode, errorMessage)
			}
		}
	})

	validTests := []struct {
		name             string
		url              string
		body             []byte
		expectedResponse string
	}{
		// before starting, we have already created 5 posts
		// test getting all posts [this is a GET request]
		// {"/api/posts return 5", "/api/posts", nil, "5"},

		// test liking post 1 and then unliking it by clicking again on the like button
		{"/api/posts/1/like return 1", "/api/posts/1/like", nil, "1"},
		{"/api/posts/1/like return 0", "/api/posts/1/like", nil, "0"},

		// TODO Need to discuss with team, how to test this, as post Date becomes an issue, right now returning empty array
		// test liking posts 2 and 3
		// {"/api/posts/2/like return 1", "/api/posts/2/like", nil, "1"},
		// {"/api/posts/3/like return 1", "/api/posts/3/like", nil, "1"},

		// test disliking post 1 and then undisliking it by clicking again on the dislike button
		{"/api/posts/1/dislike return -1", "/api/posts/1/dislike", nil, "-1"},
		{"/api/posts/1/dislike return 0", "/api/posts/1/dislike", nil, "0"},
		// dislike a post then like it, so that it returns like count as 1
		{"/api/posts/3/dislike return -1", "/api/posts/3/dislike", nil, "-1"},
		{"/api/posts/3/like return 1", "/api/posts/3/like", nil, "1"},
		// like a post then like it, so that it returns like count as 1
		{"/api/posts/3/dislike return -1", "/api/posts/3/dislike", nil, "-1"},

		{"/api/posts/1/comment return 1", "/api/posts/1/comment", []byte(`{"content": "test"}`), "1"},
		{"/api/posts/1/comment/1/like return 1", "/api/posts/1/comment/1/like", nil, "1"},
		{"/api/posts/1/comment/1/like return 0", "/api/posts/1/comment/1/like", nil, "0"},
		{"/api/posts/1/comment/1/dislike return -1", "/api/posts/1/comment/1/dislike", nil, "-1"},
		{"/api/posts/1/comment/1/dislike return 0", "/api/posts/1/comment/1/dislike", nil, "0"},
		// dislike a comment then like it, so that it returns like count as 1
		{"/api/posts/1/comment/1/dislike return -1", "/api/posts/1/comment/1/dislike", nil, "-1"},
		{"/api/posts/1/comment/1/like return 1", "/api/posts/1/comment/1/like", nil, "1"},
		{"/api/posts/1/comment/1/dislike return -1", "/api/posts/1/comment/1/dislike", nil, "-1"},
		{"/api/logout return empty string", "/api/logout", nil, ""},
	}

	for _, test := range validTests {
		t.Run(test.url, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, testServer.URL+test.url, bytes.NewReader(test.body))
			if err != nil {
				t.Fatal(err)
			}
			req.AddCookie(cookie)

			resp, err := cli.Do(req)
			if err != nil {
				t.Fatal(err)
			}

			if resp.StatusCode != http.StatusOK {
				t.Fatalf("expected %d, got %d", http.StatusOK, resp.StatusCode)
			}

			respBody, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}
			if string(respBody) != test.expectedResponse {
				t.Fatalf("expected %v, got %v", test.expectedResponse, respBody)
			}
		})
	}
}

// BenchmarkWithAuth benchmarks all routes that require authentication
func BenchmarkWithAuth(b *testing.B) {

	for i := 0; i < b.N; i++ {
		cli := testServer.Client()
		rand, err := uuid.NewV4()
		if err != nil {
			b.Error(err)
		}

		name := rand.String()[:15]
		email := fmt.Sprintf("%s@test.com", name)
		body := fmt.Sprintf(`{ "name": "%s", "email": "%s", "password": "notapassword" }`, name, email)

		resp, err := cli.Post(testServer.URL+"/api/signup", "application/json",
			strings.NewReader(body))
		if err != nil {
			b.Error(err)
		}

		if resp.StatusCode != http.StatusOK {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				b.Error(err)
			}
			b.Errorf("expected %d, got %d, %s", http.StatusOK, resp.StatusCode, body)
		}

		resp, err = cli.Post(testServer.URL+"/api/login", "application/json",
			strings.NewReader(fmt.Sprintf(`{ "login": "%s", "password": "notapassword" }`, email)))
		if err != nil {
			b.Error(err)
		}
		if len(resp.Cookies()) != 1 {
			b.Errorf("invalid cookies, expected 1, got %d", len(resp.Cookies()))
		}
		cookie := resp.Cookies()[0]
		if resp.StatusCode != http.StatusOK {
			b.Errorf("expected %d, got %d", http.StatusOK, resp.StatusCode)
		}

		request, err := http.NewRequest(http.MethodPost, testServer.URL+"/api/logout", nil)
		if err != nil {
			b.Error(err)
		}
		request.AddCookie(cookie)

		resp, err = cli.Do(request)
		if err != nil {
			b.Errorf("error while logging out: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			b.Errorf("expected %d, got %d", http.StatusOK, resp.StatusCode)
		}
		if resp.Cookies()[0].Expires.After(time.Now()) {
			b.Errorf("cookie should be expired")
		}
	}
}
