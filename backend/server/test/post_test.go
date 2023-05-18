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

	userId := srv.DB.AddUser(u.Name, u.Email, u.Pass, u.FName, u.LName, u.DoB, u.Gender)
	srv.DB.AddPost("test", "test", userId, "facts", "beatufiul, amazing, wonderful facts")

	invalidUsers := getInvalidUsers()

	t.Run("signup", func(t *testing.T) {
		resp, err := signup(cli, testServer, validUser)
		if err != nil {
			t.Fatal(err)
		}

		if resp.StatusCode != http.StatusOK {
			errBody, _ := io.ReadAll(resp.Body)
			fmt.Println(string(errBody))
			t.Fatalf("expected %d, got %d", http.StatusOK, resp.StatusCode)
		}
	})

	t.Run("invalidSignup", func(t *testing.T) {

		invalidResp, err := signup(cli, testServer, TestUser{})
		if err != nil {
			t.Fatal(err)
		}

		if invalidResp.StatusCode != http.StatusBadRequest {
			t.Fatalf("expected %d, got %d", 400, invalidResp.StatusCode)
		}

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
	})

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

	invalidPosts := getInvalidPosts()

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
		{"/api/posts/1/like return 1", "/api/posts/1/like", nil, "1"},
		{"/api/posts/1/like return 0", "/api/posts/1/like", nil, "0"},
		{"/api/posts/1/dislike return -1", "/api/posts/1/dislike", nil, "-1"},
		{"/api/posts/1/dislike return 0", "/api/posts/1/dislike", nil, "0"},
		{"/api/posts/3/dislike return -1", "/api/posts/3/dislike", nil, "-1"},
		{"/api/posts/3/like return 1", "/api/posts/3/like", nil, "1"},
		{"/api/posts/3/dislike return -1", "/api/posts/3/dislike", nil, "-1"},

		{"/api/posts/1/comment return 1", "/api/posts/1/comment", []byte(`{"content": "test"}`), "1"},
		{"/api/posts/1/comment/1/like return 1", "/api/posts/1/comment/1/like", nil, "1"},
		{"/api/posts/1/comment/1/like return 0", "/api/posts/1/comment/1/like", nil, "0"},
		{"/api/posts/1/comment/1/dislike return -1", "/api/posts/1/comment/1/dislike", nil, "-1"},
		{"/api/posts/1/comment/1/dislike return 0", "/api/posts/1/comment/1/dislike", nil, "0"},
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
