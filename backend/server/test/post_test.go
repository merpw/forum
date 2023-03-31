package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	// from package database import Post

	"github.com/gofrs/uuid"
)

// TestUser to be used in all Post tests.
type TestUser struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// cookie to simulate logged in user.
var cookie *http.Cookie

// TestWithAuth tests all routes that require authentication
func TestWithAuth(t *testing.T) {
	cli := testServer.Client()

	validUser := TestUser{
		Name:     "test",
		Email:    "test@test.com",
		Password: "SuperAmazingPassword()!@*#)(!@#",
	}

	// Slice of invalid users. It will cover most nonDB test cases.
	invalidUsers := []TestUser{
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
			Email:    "😂😂😂😂😂😂😂😂😂😂😂", // invalid email
			Password: "ValidPassword123",
		},
		{
			Name:     "ValidName",
			Email:    "valid@test.com",
			Password: "1234", // invalid password: too short
		},
		{}, // empty user
	}

	t.Run("signup", func(t *testing.T) {
		body, err := json.Marshal(validUser)
		if err != nil {
			t.Fatal(err)
		}

		resp, err := cli.Post(testServer.URL+"/api/signup", "application/json", bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}

		// bug found, if thee user is already present, we give back 400 instead of 409
		if resp.StatusCode != 200 {
			t.Fatalf("expected %d, got %d", 200, resp.StatusCode)
		}

	})

	t.Run("invalidSignup", func(t *testing.T) {
		invalidBody, err := json.Marshal("")
		if err != nil {
			t.Fatal(err)
		}

		invalidResp, err := cli.Post(testServer.URL+"/api/signup", "application/json", bytes.NewReader(invalidBody))
		if err != nil {
			t.Fatal(err)
		}

		if invalidResp.StatusCode != 400 {
			t.Fatalf("expected %d, got %d", 400, invalidResp.StatusCode)
		}

		for _, user := range invalidUsers {
			body, err := json.Marshal(user)
			if err != nil {
				t.Fatal(err)
			}

			resp, err := cli.Post(testServer.URL+"/api/signup", "application/json", bytes.NewReader(body))
			if err != nil {
				t.Fatal(err)
			}
			// bug found, if the user is already present, we give back 400 instead of 409
			if resp.StatusCode != 400 {
				t.Fatalf("expected %d, got %d", http.StatusBadRequest, resp.StatusCode)
			}
		}

	})

	t.Run("login", func(t *testing.T) {
		cookie = dummyLogin(t, cli, testServer, validUser)
		// for _, user := range invalidUsers {
		// 	cookie = dummyLogin(t, cli, testServer, user)
		// }
	})

	// test create 5 posts
	t.Run("createPost", func(t *testing.T) {
		for i := 1; i <= 5; i++ {
			body := struct {
				Title      string   `json:"title"`
				Content    string   `json:"content"`
				Categories []string `json:"categories"`
			}{
				Title:      fmt.Sprintf("Test Title %d", i),
				Content:    fmt.Sprintf("Test Content %d", i),
				Categories: []string{"facts"},
			}
			requestBodyBytes, err := json.Marshal(body)
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest(http.MethodPost, testServer.URL+"/api/posts/create/",
				bytes.NewReader(requestBodyBytes))
			if err != nil {
				t.Fatal(err)
			}
			req.AddCookie(cookie)

			resp, err := cli.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			if resp.StatusCode != 200 {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					log.Println(err)
				}
				errorMessage := string(body)
				t.Fatalf("expected %d, got %d, ErrBody: %v", 200, resp.StatusCode, errorMessage)
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
		// TODO DISCUSS WITH TEAM: should we return all posts on Post request also or only on GET request?
		// in the current implementation, we return all posts on GET request only
		// when we make a POST request currently to fetch allPosts, we get a method not allowed error,

		// test liking post 1 and then unliking it by clicking again on the like button
		{"/api/posts/1/like return 1", "/api/posts/1/like", nil, "1"},
		{"/api/posts/1/like return 0", "/api/posts/1/like", nil, "0"},

		// TODO Need to discuss with team, how to test this, as post Date becomes an issue, right now returning empty array
		// test liking posts 2 and 3
		// {"/api/posts/2/like return 1", "/api/posts/2/like", nil, "1"},
		// {"/api/posts/3/like return 1", "/api/posts/3/like", nil, "1"},

		// test getting the posts liked by the user
		{"/api/me/posts/liked return 0 POSTS", "/api/me/posts/liked", nil, "[]"},

		// test disliking post 1 and then undisliking it by clicking again on the dislike button
		{"/api/posts/1/dislike return -1", "/api/posts/1/dislike", nil, "-1"},
		{"/api/posts/1/dislike return 0", "/api/posts/1/dislike", nil, "0"},

		{"/api/posts/1/comment return 1", "/api/posts/1/comment", []byte(`{"content": "test"}`), "1"},
		{"/api/posts/1/comment/1/like return 1", "/api/posts/1/comment/1/like", nil, "1"},
		{"/api/posts/1/comment/1/like return 0", "/api/posts/1/comment/1/like", nil, "0"},
		{"/api/posts/1/comment/1/dislike return -1", "/api/posts/1/comment/1/dislike", nil, "-1"},
		{"/api/posts/1/comment/1/dislike return 0", "/api/posts/1/comment/1/dislike", nil, "0"},
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
			if resp.StatusCode != 200 {
				// print the full response
				body, _ := io.ReadAll(resp.Body)
				fmt.Println(string(body))
				t.Fatalf("expected %d, got %d", 200, resp.StatusCode)
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

		if resp.StatusCode != 200 {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				b.Error(err)
			}
			b.Errorf("expected %d, got %d, %s", 200, resp.StatusCode, body)
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
		if resp.StatusCode != 200 {
			b.Errorf("expected %d, got %d", 200, resp.StatusCode)
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
		if resp.StatusCode != 200 {
			b.Errorf("expected %d, got %d", 200, resp.StatusCode)
		}
		if resp.Cookies()[0].Expires.After(time.Now()) {
			b.Errorf("cookie should be expired")
		}
	}
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

	if resp.StatusCode != 200 {
		t.Fatalf("expected %d, got %d", 200, resp.StatusCode)
	}
	return cookie
}
