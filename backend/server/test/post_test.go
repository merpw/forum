package server

import (
	"database/sql"
	"fmt"
	"forum/server"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// TODO: move server init to separate func (to prevent `email is already taken` false test fail)

// TestAuth tests auth routes
func TestAuth(t *testing.T) {
	db, err := sql.Open("sqlite3", "./test.db?_foreign_keys=true")
	if err != nil {
		t.Fatal(err)
	}
	srv := server.Connect(db)
	router := srv.Start()
	testServer := httptest.NewServer(router)
	defer testServer.Close()

	cli := testServer.Client()

	t.Run("signup", func(t *testing.T) {
		rand.Seed(time.Now().UnixNano())
		aUsr := fmt.Sprintf("test%d", rand.Intn(100000))
		anEmail := fmt.Sprintf("%s@test.com", aUsr)
		body := fmt.Sprintf(`{ "name": "%s", "email": "%s", "password": "notapassword" }`, aUsr, anEmail)
		resp, err := cli.Post(testServer.URL+"/api/signup", "application/json",
			strings.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}

		// bug found, if the user is already present, we give back 400 instead of 409
		if resp.StatusCode != 200 {
			t.Fatalf("expected %d, got %d", 200, resp.StatusCode)
		}
	})

	var cookies []*http.Cookie
	t.Run("login", func(t *testing.T) {
		resp, err := cli.Post(testServer.URL+"/api/login", "application/json",
			strings.NewReader(`{ "login": "test@test.com", "password": "notapassword" }`))
		if err != nil {
			t.Fatal(err)
		}
		cookies = resp.Cookies()

		// TODO: add request body
		if resp.StatusCode != 200 {
			t.Fatalf("expected %d, got %d", 200, resp.StatusCode)
		}
	})

	t.Run("logout", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, testServer.URL+"/api/logout", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.AddCookie(cookies[0])

		resp, err := cli.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != 200 {
			t.Fatalf("expected %d, got %d", 200, resp.StatusCode)
		}
	})
}

// TestPost tests post routes
func TestPost(t *testing.T) {
	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		t.Fatal(err)
	}
	srv := server.Connect(db)
	router := srv.Start()
	testServer := httptest.NewServer(router)
	defer testServer.Close()

	cli := testServer.Client()

	tests := []struct {
		url string
	}{
		{"/api/posts/create"},
		{"/api/posts/1/like"},
		{"/api/posts/1/dislike"},
		{"/api/posts/1/comment"},

		{"/api/posts/1/comment/1/like"},
		{"/api/posts/1/comment/1/dislike"},
	}

	for _, test := range tests {
		t.Run(test.url, func(t *testing.T) {
			resp, err := cli.Post(testServer.URL+test.url, "application/json", nil)
			if err != nil {
				t.Fatal(err)
			}
			if resp.StatusCode != http.StatusUnauthorized {
				t.Fatalf("expected %d, got %d", http.StatusUnauthorized, resp.StatusCode)
			}
		})
	}
	//	TODO: add POST tests with auth
}
