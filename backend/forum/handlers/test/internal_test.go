package server_test

import (
	. "backend/forum/handlers/test/server"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"testing"
)

func TestCheckSession(t *testing.T) {
	testServer := NewTestServer(t)
	cli := testServer.TestClient()
	cli.TestAuth(t)

	t.Run("Invalid header", func(t *testing.T) {
		// No header
		cli.TestGet(t, "/api/internal/check-session", http.StatusForbidden)

		err := os.Setenv("FORUM_BACKEND_SECRET", "secret")
		if err != nil {
			t.Fatal(err)
		}
		req := generateInternalRequest(t, testServer, "/api/internal/check-session")
		err = os.Setenv("FORUM_BACKEND_SECRET", "secret has changed")
		if err != nil {
			t.Fatal(err)
		}
		cli.TestRequest(t, req, http.StatusUnauthorized)
	})

	err := os.Setenv("FORUM_BACKEND_SECRET", "super secret not guessable token")
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Invalid token", func(t *testing.T) {
		req := generateInternalRequest(t, testServer, "/api/internal/check-session")
		cli.TestRequest(t, req, http.StatusBadRequest)

		req = generateInternalRequest(t, testServer, "/api/internal/check-session?token=invalid")
		_, respBody := cli.TestRequest(t, req, http.StatusOK)

		var errResp struct {
			Error string `json:"error"`
		}
		err = json.Unmarshal(respBody, &errResp)
		if err != nil {
			t.Fatal(err)
		}
		if errResp.Error == "" {
			t.Fatal("Expected error")
		}
	})

	t.Run("Invalid method", func(t *testing.T) {
		req := generateInternalRequest(t, testServer, "/api/internal/check-session?token=invalid")
		req.Method = http.MethodPost
		cli.TestRequest(t, req, http.StatusMethodNotAllowed)
	})

	t.Run("Valid", func(t *testing.T) {
		var token string
		for _, cookie := range cli.Cookies {
			if cookie.Name == "forum-token" {
				token = cookie.Value
			}
		}

		_, respBody := cli.TestRequest(t,
			generateInternalRequest(t, testServer, "/api/internal/check-session?token="+token),
			http.StatusOK)
		userId, err := strconv.Atoi(string(respBody))
		if err != nil {
			t.Fatal(err)
		}

		if userId != 1 {
			t.Fatalf("Expected user id 1, got %d", userId)
		}
	})
}

func TestBypassAuth(t *testing.T) {
	testServer := NewTestServer(t)
	cli := testServer.TestClient()

	cli.TestGet(t, "/api/me", http.StatusUnauthorized)

	err := os.Setenv("FORUM_BACKEND_SECRET", "super secret nobody knows")
	if err != nil {
		t.Fatal(err)
	}

	req := generateInternalRequest(t, testServer, "/api/me")

	cli.TestRequest(t, req, http.StatusInternalServerError)

	err = os.Setenv("FORUM_BACKEND_SECRET", "secret has changed")
	if err != nil {
		t.Fatal(err)
	}

	cli.TestRequest(t, req, http.StatusUnauthorized)
}

func generateInternalRequest(t *testing.T, testServer TestServer, path string) *http.Request {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, testServer.URL+path, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Internal-Auth", os.Getenv("FORUM_BACKEND_SECRET"))
	return req
}
