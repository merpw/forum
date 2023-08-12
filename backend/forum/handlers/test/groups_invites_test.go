package server_test

import (
	. "backend/forum/handlers/test/server"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

func TestGroupIdJoin(t *testing.T) {
	testServer := NewTestServer(t)
	cli1 := testServer.TestClient()
	cli2 := testServer.TestClient()
	cli1.TestAuth(t)
	group := CreateGroup("test", "test", []int{})
	cli1.TestPost(t, "/api/groups/create", group, http.StatusOK)

	cli2.TestAuth(t)

	t.Run("Not found", func(t *testing.T) {
		cli1.TestPost(t, "/api/groups/2/join", nil, http.StatusNotFound)
		cli1.TestPost(t, "/api/groups/999999999999999999999999999999/join", nil, http.StatusNotFound)
	})

	t.Run("Bad request", func(t *testing.T) {
		cli1.TestPost(t, "/api/groups/1/join", nil, http.StatusBadRequest)
	})

	t.Run("Valid", func(t *testing.T) {
		var status int
		t.Run("Request to join", func(t *testing.T) {
			_, resp := cli2.TestPost(t, "/api/groups/1/join", nil, http.StatusOK)
			if err := json.Unmarshal(resp, &status); err != nil {
				t.Fatal(err)
			}
			if status != 2 {
				t.Errorf("unexpected status, expected %d, got %d", 2, status)
			}
		})
		t.Run("Abort request", func(t *testing.T) {
			_, resp := cli2.TestPost(t, "/api/groups/1/join", nil, http.StatusOK)
			if err := json.Unmarshal(resp, &status); err != nil {
				t.Fatal(err)
			}
			if status != 0 {
				t.Errorf("unexpected status, expected %d, got %d", 2, status)
			}
		})

		t.Run("Request, approve, and leave", func(t *testing.T) {
			_, resp := cli2.TestPost(t, "/api/groups/1/join", nil, http.StatusOK)
			if err := json.Unmarshal(resp, &status); err != nil {
				t.Fatal(err)
			}
			if status != 2 {
				t.Errorf("unexpected status, expected %d, got %d", 2, status)
			}
			response := struct {
				Response bool `json:"response"`
			}{
				Response: true,
			}
			var invites []int
			_, resp2 := cli1.TestGet(t, "/api/invitations", http.StatusOK)
			if err := json.Unmarshal(resp2, &invites); err != nil {
				t.Fatal(err)
			}

			if len(invites) != 1 {
				t.Errorf("Invalid length of invites, expected %d, got %d", 1, len(invites))
			}

			endpoint := fmt.Sprintf("/api/invitations/%d/respond", invites[0])
			cli1.TestPost(t, endpoint, response, http.StatusOK)

		})
	})
}

func TestGroupIdInviteLeave(t *testing.T) {
	testServer := NewTestServer(t)
	cli1 := testServer.TestClient()
	cli2 := testServer.TestClient()
	cli1.TestAuth(t)

	group := CreateGroup("test", "test", []int{})
	cli1.TestPost(t, "/api/groups/create", group, http.StatusOK)

	invalidBody := struct {
		UserId int `json:"user_id"`
	}{
		UserId: 666,
	}

	t.Run("Unauthorized", func(t *testing.T) {
		cli2.TestPost(t, "/api/groups/1/leave", nil, http.StatusUnauthorized)
	})

	cli2.TestAuth(t)

	t.Run("Not found", func(t *testing.T) {
		cli1.TestPost(t, "/api/groups/1/invite", invalidBody, http.StatusNotFound)
		cli1.TestPost(t, "/api/groups/10/leave", nil, http.StatusNotFound)
		cli1.TestPost(t, "/api/groups/10/invite", nil, http.StatusNotFound)
		cli1.TestPost(t, "/api/groups/999999999999999999999999999999/invite", nil, http.StatusNotFound)
		cli1.TestPost(t, "/api/groups/999999999999999999999999999999/leave", nil, http.StatusNotFound)
	})

	t.Run("Bad request", func(t *testing.T) {
		cli1.TestPost(t, "/api/groups/1/leave", nil, http.StatusBadRequest)
		cli2.TestPost(t, "/api/groups/1/leave", nil, http.StatusBadRequest)
		cli1.TestPost(t, "/api/groups/1/invite", "invalid", http.StatusBadRequest)
	})

	t.Run("Valid", func(t *testing.T) {
		requestBody := struct {
			UserId int `json:"user_id"`
		}{
			UserId: 2,
		}

		t.Run("Invite request", func(t *testing.T) {
			var status int
			_, resp := cli1.TestPost(t, "/api/groups/1/invite", requestBody, http.StatusOK)
			if err := json.Unmarshal(resp, &status); err != nil {
				t.Fatal(err)
			}
			if status != 2 {
				t.Errorf("unexpected status, expected %d, got %d", 2, status)
			}
		})

		t.Run("Leave while pending invite", func(t *testing.T) {
			cli2.TestPost(t, "/api/groups/1/leave", nil, http.StatusBadRequest)
		})

		t.Run("Invite again", func(t *testing.T) {
			// This is an error test
			cli1.TestPost(t, "/api/groups/1/invite", requestBody, http.StatusBadRequest)
		})

		t.Run("Respond to invite", func(t *testing.T) {
			response := struct {
				Response bool `json:"response"`
			}{
				Response: true,
			}
			var invites []int
			_, resp := cli2.TestGet(t, "/api/invitations", http.StatusOK)
			if err := json.Unmarshal(resp, &invites); err != nil {
				t.Fatal(err)
			}

			if len(invites) != 1 {
				t.Errorf("Invalid length of invites, expected %d, got %d", 1, len(invites))
			}

			endpoint := fmt.Sprintf("/api/invitations/%d/respond", invites[0])
			cli2.TestPost(t, endpoint, response, http.StatusOK)
		})

		t.Run("Invite after accept", func(t *testing.T) {
			cli1.TestPost(t, "/api/groups/1/invite", requestBody, http.StatusBadRequest)
		})

		t.Run("Leave after accept", func(t *testing.T) {
			var status int
			_, resp := cli2.TestPost(t, "/api/groups/1/leave", nil, http.StatusOK)
			if err := json.Unmarshal(resp, &status); err != nil {
				t.Fatal(err)
			}

			if status != 0 {
				t.Errorf("unexpected status, expected %d, got %d", 0, status)
			}
		})

	})
}
