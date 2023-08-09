package server_test

import (
	. "backend/forum/handlers/test/server"
	"encoding/json"
	"net/http"
	"testing"
)

func TestGroups(t *testing.T) {
	testServer := NewTestServer(t)
	cli := testServer.TestClient()
	cli2 := testServer.TestClient()

	requestBody := struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Invite      []int  `json:"invite"`
	}{
		Title:       "Test title",
		Description: "Test description",
		Invite:      []int{1},
	}

	t.Run("Invalid", func(t *testing.T) {

		t.Run("Unauthorized", func(t *testing.T) {
			cli.TestPost(t, "/api/groups/create", requestBody, http.StatusUnauthorized)
		})

		cli.TestAuth(t)

		t.Run("Method", func(t *testing.T) {
			cli.TestGet(t, "/api/groups/create", http.StatusMethodNotAllowed)
		})

		t.Run("Invite yourself", func(t *testing.T) {
			cli.TestPost(t, "/api/groups/create", requestBody, http.StatusBadRequest)
		})

		t.Run("Invalid requestBody", func(t *testing.T) {
			cli.TestPost(t, "/api/groups/create", "test", http.StatusBadRequest)
		})
	})

	// Log in user 2 so they can be invited
	cli2.TestAuth(t)
	t.Run("Valid", func(t *testing.T) {
		requestBody.Invite = []int{2}
		_, respBody := cli.TestPost(t, "/api/groups/create", requestBody, http.StatusOK)

		t.Run("Add group", func(t *testing.T) {
			var groupId int
			err := json.Unmarshal(respBody, &groupId)
			if err != nil {
				t.Fatal(err)
			}

			if groupId != 1 {
				t.Errorf("unexpected group id, expected %d, got %d", 1, groupId)
			}
		})

		t.Run("Get all groups by members", func(t *testing.T) {
			_, respBody := cli.TestGet(t, "/api/groups", http.StatusOK)
			var groups []int
			err := json.Unmarshal(respBody, &groups)
			if err != nil {
				t.Fatal(err)
			}

			if len(groups) != 1 {
				t.Errorf("invalid groups length, expected %d, got %d", 1, len(groups))
			}

		})

	})

}
