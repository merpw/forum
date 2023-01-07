package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gofrs/uuid"
	"io"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"
)

// TestWithAuth tests all routes that require authentication
func TestWithAuth(t *testing.T) {
	cli := testServer.Client()

	testUser := struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		Name:     "test",
		Email:    "test@test.com",
		Password: "SuperAmazingPassword()!@*#)(!@#",
	}

	t.Run("signup", func(t *testing.T) {
		body, err := json.Marshal(testUser)
		if err != nil {
			t.Fatal(err)
		}

		resp, err := cli.Post(testServer.URL+"/api/signup", "application/json", bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}

		// bug found, if the user is already present, we give back 400 instead of 409
		if resp.StatusCode != 200 {
			t.Fatalf("expected %d, got %d", 200, resp.StatusCode)
		}
	})

	var cookie *http.Cookie
	t.Run("login", func(t *testing.T) {
		body := fmt.Sprintf(`{ "login": "%v", "password": "%v" }`, testUser.Email, testUser.Password)
		resp, err := cli.Post(testServer.URL+"/api/login", "application/json",
			strings.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}
		if len(resp.Cookies()) == 0 {
			t.Fatal("no cookies after login")
		}
		cookie = resp.Cookies()[0]

		if resp.StatusCode != 200 {
			t.Fatalf("expected %d, got %d", 200, resp.StatusCode)
		}
	})

	// test create post
	t.Run("createPost", func(t *testing.T) {
		body := struct {
			Title      string   `json:"title"`
			Content    string   `json:"content"`
			Categories []string `json:"categories"`
		}{
			Title:      "Test Title",
			Content:    "Test Content",
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
	})

	validTests := []struct {
		url              string
		body             []byte
		expectedResponse string
	}{
		{"/api/posts/1/like", nil, "1"},
		{"/api/posts/1/like", nil, "0"},

		{"/api/posts/1/dislike", nil, "-1"},
		{"/api/posts/1/dislike", nil, "0"},

		{"/api/posts/1/comment", []byte(`{"content": "test"}`), "1"},
		{"/api/posts/1/comment/1/like", nil, "1"},
		{"/api/posts/1/comment/1/like", nil, "0"},
		{"/api/posts/1/comment/1/dislike", nil, "-1"},
		{"/api/posts/1/comment/1/dislike", nil, "0"},
		{"/api/logout", nil, ""},
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
