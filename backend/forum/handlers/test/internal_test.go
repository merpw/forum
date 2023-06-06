package server_test

import (
	. "backend/forum/handlers/test/server"
	"encoding/json"
	"net/http"
	"strconv"
	"testing"
)

func TestCheckSession(t *testing.T) {
	testServer := NewTestServer(t)
	cli := testServer.TestClient()
	cli.TestAuth(t)

	type Error struct {
		Error string `json:"error"`
	}

	cli.TestGet(t, "/api/internal/check-session", http.StatusBadRequest)

	_, respBody := cli.TestGet(t, "/api/internal/check-session?token=invalid", http.StatusOK)
	var errResp Error
	err := json.Unmarshal(respBody, &errResp)
	if err != nil {
		t.Fatal(err)
	}
	if errResp.Error == "" {
		t.Fatal("Expected error")
	}

	var token string
	for _, cookie := range cli.Cookies {
		if cookie.Name == "forum-token" {
			token = cookie.Value
		}
	}

	_, respBody = cli.TestGet(t, "/api/internal/check-session?token="+token, http.StatusOK)
	userId, err := strconv.Atoi(string(respBody))
	if err != nil {
		t.Fatal(err)
	}

	if userId != 1 {
		t.Fatalf("Expected user id 1, got %d", userId)
	}
}
