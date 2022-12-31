package server

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"forum/server"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

// TODO: move server init to separate func (to prevent `email is already taken` false test fail)
// initServer initializes the server and database connection
func initServer(db *sql.DB) *server.Server {
	srv := server.Connect(db)
	err := srv.DB.InitDatabase()
	if err != nil {
		log.Fatal(err)
	}
	return srv
}

// TestAuth tests auth routes
func TestAuth(t *testing.T) {
	// prepare database and server for testing
	os.Remove("./test.db")
	db, err := sql.Open("sqlite3", "./test.db?_foreign_keys=true")
	if err != nil {
		t.Fatal(err)
	}
	srv := initServer(db)
	router := srv.Start()
	testServer := httptest.NewServer(router)
	defer testServer.Close()

	cli := testServer.Client()
	aUsr := fmt.Sprintf("%dtest%d", rand.Intn(100000), rand.Intn(100000))
	anEmail := fmt.Sprintf("%s@test.com", aUsr)
	t.Run("signup", func(t *testing.T) {
		rand.Seed(time.Now().UnixNano())
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
		body := fmt.Sprintf(`{ "login": "%s", "password": "notapassword" }`, anEmail)
		resp, err := cli.Post(testServer.URL+"/api/login", "application/json",
			strings.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}
		cookies = resp.Cookies()

		// TODO: add request body
		if resp.StatusCode != 200 {
			t.Fatalf("expected %d, got %d", 200, resp.StatusCode)
		}
	})

	// test create post
	t.Run("createPost", func(t *testing.T) {
		body := struct {
			Title    string `json:"title"`
			Content  string `json:"content"`
			Category string `json:"category"`
		}{
			Title:    "Test Title",
			Content:  "Test Content",
			Category: "facts",
		}
		requestBodyBytes, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPost, testServer.URL+"/api/posts/create/",
			bytes.NewBuffer(requestBodyBytes))
		if err != nil {
			t.Fatal(err)
		}
		req.AddCookie(cookies[0])

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

	t.Run("logout", func(t *testing.T) {
		testAuthLogout(t, cli, testServer, cookies)
	})
}

// BenchmarkAuth benchmarks auth routes, with i number of users in the for loop
func BenchmarkAuth(b *testing.B) {
	// delete the database file
	os.Remove("./test.db")
	db, err := sql.Open("sqlite3", "./test.db?_foreign_keys=true")
	if err != nil {
		b.Fatal(err)
	}
	srv := initServer(db)

	router := srv.Start()
	testServer := httptest.NewServer(router)
	defer testServer.Close()
	cli := testServer.Client()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		numUsers := 100
		wg := sync.WaitGroup{}
		wg.Add(numUsers)

		counter := 0
		mutex := &sync.Mutex{}
		for i := 0; i < numUsers; i++ {
			go func() {
				mutex.Lock()
				defer mutex.Unlock()

				// signup
				counter++
				aUsr := fmt.Sprintf("user%d", counter)
				anEmail := fmt.Sprintf("%s@test.com", aUsr)
				body := fmt.Sprintf(`{ "name": "%s", "email": "%s", "password": "notapassword" }`, aUsr, anEmail)
				resp, err := cli.Post(testServer.URL+"/api/signup", "application/json",
					strings.NewReader(body))
				if err != nil {
					b.Error(err)
				}

				// bug found, if the user is already present, we give back 400 instead of 409
				if resp.StatusCode != 200 {
					b.Errorf("expected %d, got %d", 200, resp.StatusCode)
				}

				// login
				resp, err = cli.Post(testServer.URL+"/api/login", "application/json",
					strings.NewReader(fmt.Sprintf(`{ "login": "%s", "password": "notapassword" }`, anEmail)))
				if err != nil {
					b.Error(err)
				}
				cookies := resp.Cookies()
				if resp.StatusCode != 200 {
					b.Errorf("expected %d, got %d", 200, resp.StatusCode)
				}

				// logout
				logoutB(cli, testServer.URL, cookies[0], b)

				wg.Done()
			}()
		}

		wg.Wait()
	}
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

func testAuthLogout(t *testing.T, cli *http.Client, testServer *httptest.Server, cookies []*http.Cookie) {
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
}

func logoutB(cli *http.Client, url string, cookie *http.Cookie, b *testing.B) {
	// logout
	req, err := http.NewRequest(http.MethodPost, url+"/api/logout", nil)
	if err != nil {
		b.Error(err)
	}
	req.AddCookie(cookie)

	resp, err := cli.Do(req)
	if err != nil {
		b.Error(err)
	}
	if resp.StatusCode != 200 {
		b.Errorf("expected %d, got %d", 200, resp.StatusCode)
	}
}
