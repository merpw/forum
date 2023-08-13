package server_test

import (
	"backend/common/integrations/auth"
	. "backend/forum/handlers/test/server"
	"bufio"
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
		req := GenerateInternalRequest(t, testServer, "/api/internal/check-session")
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
		req := GenerateInternalRequest(t, testServer, "/api/internal/check-session")
		cli.TestRequest(t, req, http.StatusBadRequest)

		req = GenerateInternalRequest(t, testServer, "/api/internal/check-session?token=invalid")
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
		req := GenerateInternalRequest(t, testServer, "/api/internal/check-session?token=invalid")
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
			GenerateInternalRequest(t, testServer, "/api/internal/check-session?token="+token),
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

	req := GenerateInternalRequest(t, testServer, "/api/me")

	cli.TestRequest(t, req, http.StatusInternalServerError)

	err = os.Setenv("FORUM_BACKEND_SECRET", "secret has changed")
	if err != nil {
		t.Fatal(err)
	}

	cli.TestRequest(t, req, http.StatusUnauthorized)
}

func TestEvents(t *testing.T) {
	testServer := NewTestServer(t)
	cli := testServer.TestClient()

	err := os.Setenv("FORUM_BACKEND_SECRET", "super secret nobody knows")
	if err != nil {
		t.Fatal(err)
	}

	tokens := make([]string, 5)
	go func() {
		for i := 0; i < 5; i++ {
			if t.Failed() {
				return
			}
			cli := testServer.TestClient()
			cli.TestAuth(t)
			for _, cookie := range cli.Cookies {
				if cookie.Name == "forum-token" {
					tokens[i] = cookie.Value
				}
			}
			cli.TestPost(t, "/api/logout", nil, http.StatusOK)
		}
	}()

	req := GenerateInternalRequest(t, testServer, "/api/internal/events")
	resp, err := cli.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	reader := bufio.NewReader(resp.Body)

	for i := 0; i < 5; i++ {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			t.Fatal(err)
		}

		var event auth.Event

		err = json.Unmarshal(line, &event)
		if err != nil {
			t.Fatal(err)
		}

		if event.Type != auth.EventTypeTokenRevoked {
			t.Fatalf("Expected event type %s, got %s", auth.EventTypeTokenRevoked, event.Type)
		}

		if event.Item != tokens[i] {
			t.Fatal("Expected token", tokens[i], "got", event.Item)
		}
	}
}

func GenerateInternalRequest(t *testing.T, testServer TestServer, path string) *http.Request {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, testServer.URL+path, nil)
	if err != nil {
		t.Fatal(err)
	}
	var secret = os.Getenv("FORUM_BACKEND_SECRET")

	if secret == "" {
		secret = "super secret nobody knows"
	}

	req.Header.Add("Internal-Auth", secret)
	return req
}
