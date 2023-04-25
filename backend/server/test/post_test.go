package server_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"forum/server"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofrs/uuid"
)

// cookie to simulate logged-in user.
var cookie *http.Cookie

type TestPost struct {
	Title      string
	Content    string
	Categories []string
}

var invalidPosts = []TestPost{
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

// TestUser to be used in all Post tests.
type TestUser struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	DoB       string `json:"dob"`
	Gender    string `json:"gender"`
}

var invalidUsers = []TestUser{
	{
		Name:     "test1",
		Email:    "test@test.com", // email already in use
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

// err := createPost(cli, testServer, cookie, title, content, []string{"facts"})
func createPost(cli *http.Client, sURL string, ck *http.Cookie, p TestPost) (*http.Response, error) {
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

// TestWithAuth tests all routes that require authentication
func TestWithAuth(t *testing.T) {
	db, _, testServer, cli, err := setupServer()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	defer testServer.Close()

	testUser := struct {
		Name      string `json:"name"`
		Email     string `json:"email"`
		Password  string `json:"password"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		DoB       string `json:"dob"`
		Gender    string `json:"gender"`
	}{
		Name:      "test",
		Email:     "test@test.com",
		Password:  "SuperAmazingPassword()!@*#)(!@#",
		FirstName: "John",
		LastName:  "Doe",
		DoB:       "2000-01-01",
		Gender:    "male",
	}

	t.Run("signup", func(t *testing.T) {
		resp, err := signup(cli, testServer, testUser)
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
		cookie = dummyLogin(t, cli, testServer, testUser)
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

	// test create 5 posts
	t.Run("createPost", func(t *testing.T) {
		for i := 1; i <= 5; i++ {
			p := TestPost{
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

	// TODO: DOCUMENTATION
	t.Run("createInvalidPost", func(t *testing.T) {
		resp, err := createPost(cli, testServer.URL, cookie, TestPost{"", "", []string{""}})
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
