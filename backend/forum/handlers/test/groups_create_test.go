package server_test

import (
	. "backend/forum/handlers/test/server"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func TestCreateGroup(t *testing.T) {
	testServer := NewTestServer(t)

	invalidBody := CreateGroup("test title", "test desc", []int{1})

	t.Run("Invalid", func(t *testing.T) {
		cli := testServer.TestClient()

		t.Run("Unauthorized", func(t *testing.T) {
			cli.TestPost(t, "/api/groups/create", invalidBody, http.StatusUnauthorized)
		})

		cli.TestAuth(t)

		t.Run("Method", func(t *testing.T) {
			cli.TestGet(t, "/api/groups/create", http.StatusMethodNotAllowed)
		})

		t.Run("Invite yourself", func(t *testing.T) {
			cli.TestPost(t, "/api/groups/create", invalidBody, http.StatusBadRequest)
		})

		t.Run("Invite invalid user", func(t *testing.T) {
			invalidBody.Invite = []int{1337}
			cli.TestPost(t, "/api/groups/create", invalidBody, http.StatusBadRequest)
		})

		t.Run("Invalid requestBody", func(t *testing.T) {
			cli.TestPost(t, "/api/groups/create", "test", http.StatusBadRequest)
		})

		t.Run("Too short title", func(t *testing.T) {
			invalidBody.Title = ""
			cli.TestPost(t, "/api/groups/create", invalidBody, http.StatusBadRequest)
		})

		invalidBody.Title = strings.Repeat("t", 201)
		t.Run("Too long title", func(t *testing.T) {
			cli.TestPost(t, "/api/groups/create", invalidBody, http.StatusBadRequest)
		})

		t.Run("Too short description", func(t *testing.T) {
			invalidBody.Title = "Test"
			invalidBody.Description = ""
			cli.TestPost(t, "/api/groups/create", invalidBody, http.StatusBadRequest)
		})
		t.Run("Too long description", func(t *testing.T) {
			invalidBody.Description = strings.Repeat("t", 201)
			cli.TestPost(t, "/api/groups/create", invalidBody, http.StatusBadRequest)
		})
	})

	var invites []int

	// Creates a group with valid title, description, and with 6 users in it
	for i := 0; i <= 5; i++ {
		client := testServer.TestClient()

		client.TestAuth(t)
		var meResp struct {
			Id int `json:"id"`
		}
		_, respBody := client.TestGet(t, "/api/me", http.StatusOK)
		err := json.Unmarshal(respBody, &meResp)
		if err != nil {
			t.Fatal(err)
		}

		invites = append(invites, meResp.Id)
	}

	validBody := CreateGroup("test title", "test desc", invites)

	t.Run("Valid", func(t *testing.T) {
		cli := testServer.TestClient()
		cli.TestAuth(t)

		t.Run("Full", func(t *testing.T) {
			var groupId int
			_, respBody := cli.TestPost(t, "/api/groups/create", validBody, http.StatusOK)

			err := json.Unmarshal(respBody, &groupId)
			if err != nil {
				t.Fatal(err)
			}

			if groupId != 1 {
				t.Errorf("unexpected group id, expected %d, got %d", 1, groupId)
			}

			cli.TestGet(t, fmt.Sprintf("/api/groups/%d", groupId), http.StatusOK)

			_, resp := cli.TestGet(t, "/api/groups", http.StatusOK)
			var groupIds []int
			if err := json.Unmarshal(resp, &groupIds); err != nil {
				t.Fatal(err)
			}
			// Check if the group is in the list of groups
			found := false
			for _, id := range groupIds {
				if id == groupId {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("group not found in list of groups")
			}

			endpoint := fmt.Sprintf("/api/groups/%d", groupId)

			_, respBody = cli.TestGet(t, endpoint, http.StatusOK)

			responseBody := struct {
				Id           int    `json:"id"`
				Title        string `json:"title"`
				Description  string `json:"description"`
				MemberStatus int    `json:"member_status"`
				Members      int    `json:"member_count"`
			}{}

			err = json.Unmarshal(respBody, &responseBody)
			if err != nil {
				t.Fatal(err)
			}

			if responseBody.Id != groupId {
				t.Errorf("invalid id, expected %d, got %d", groupId, responseBody.Id)
			}

			if responseBody.Title != validBody.Title {
				t.Errorf("invalid title, expected %s, got %s", validBody.Title, responseBody.Title)
			}

			if responseBody.Description != validBody.Description {
				t.Errorf("invalid description, expected %s, got %s",
					validBody.Description, responseBody.Description)
			}
			if responseBody.MemberStatus != 1 {
				t.Errorf("invalid member status, expected %d, got %d", 1, responseBody.MemberStatus)
			}

			if responseBody.Members != 1 {
				t.Errorf("invalid members, expected %d, got %d", 1, responseBody.Members)
			}
		})
	})
	t.Run("No invites", func(t *testing.T) {
		cli := testServer.TestClient()
		cli.TestAuth(t)

		validBody.Invite = []int{}
		cli.TestPost(t, "/api/groups/create", validBody, http.StatusOK)
	})
}
