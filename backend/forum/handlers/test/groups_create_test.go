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

	cli1 := testServer.TestClient()

	invalidBody := CreateGroup("test title", "test desc", []int{1})

	t.Run("Invalid", func(t *testing.T) {

		t.Run("Unauthorized", func(t *testing.T) {
			cli1.TestPost(t, "/api/groups/create", invalidBody, http.StatusUnauthorized)
		})

		cli1.TestAuth(t)

		t.Run("Method", func(t *testing.T) {
			cli1.TestGet(t, "/api/groups/create", http.StatusMethodNotAllowed)
		})

		t.Run("Invite yourself", func(t *testing.T) {
			cli1.TestPost(t, "/api/groups/create", invalidBody, http.StatusBadRequest)
		})

		t.Run("Invite invalid user", func(t *testing.T) {
			invalidBody.Invite = []int{1337}
			cli1.TestPost(t, "/api/groups/create", invalidBody, http.StatusNotFound)
		})

		t.Run("Invalid requestBody", func(t *testing.T) {
			cli1.TestPost(t, "/api/groups/create", "test", http.StatusBadRequest)
		})

		t.Run("Too short title", func(t *testing.T) {
			invalidBody.Title = ""
			cli1.TestPost(t, "/api/groups/create", invalidBody, http.StatusBadRequest)
		})

		invalidBody.Title = strings.Repeat("t", 201)
		t.Run("Too long title", func(t *testing.T) {
			cli1.TestPost(t, "/api/groups/create", invalidBody, http.StatusBadRequest)
		})

		t.Run("Too short description", func(t *testing.T) {
			invalidBody.Title = "Test"
			invalidBody.Description = ""
			cli1.TestPost(t, "/api/groups/create", invalidBody, http.StatusBadRequest)
		})
		t.Run("Too long description", func(t *testing.T) {
			invalidBody.Description = strings.Repeat("t", 201)
			cli1.TestPost(t, "/api/groups/create", invalidBody, http.StatusBadRequest)
		})
	})

	var invites []int

	for i := 2; i <= 6; i++ {
		client := testServer.TestClient()
		client.TestAuth(t)
		invites = append(invites, i)
	}

	// Creates a group with valid title, description, and with 6 users in it
	validBody := CreateGroup("test title", "test desc", invites)

	t.Run("Valid", func(t *testing.T) {
		t.Run("Create", func(t *testing.T) {
			var groupId int
			_, respBody := cli1.TestPost(t, "/api/groups/create", validBody, http.StatusOK)

			err := json.Unmarshal(respBody, &groupId)
			if err != nil {
				t.Fatal(err)
			}

			if groupId != 1 {
				t.Errorf("unexpected group id, expected %d, got %d", 1, groupId)
			}
		})

		var groups []int
		t.Run("Get all groups", func(t *testing.T) {
			_, resp := cli1.TestGet(t, "/api/groups", http.StatusOK)
			if err := json.Unmarshal(resp, &groups); err != nil {
				t.Fatal(err)
			}

			if len(groups) != 1 {
				t.Errorf("unexpected amount of groups, expected %d, got %d", 1, len(groups))
			}

		})

		responseBody := struct {
			Id           int    `json:"id"`
			Title        string `json:"title"`
			Description  string `json:"description"`
			MemberStatus int    `json:"member_status"`
			Members      int    `json:"members"`
		}{}

		t.Run("Get group info", func(t *testing.T) {
			endpoint := fmt.Sprintf("/api/groups/%d", groups[0])

			_, respBody := cli1.TestGet(t, endpoint, http.StatusOK)

			err := json.Unmarshal(respBody, &responseBody)
			if err != nil {
				t.Fatal(err)
			}
			if responseBody.Id != groups[0] {
				t.Errorf("invalid id, expected %d, got %d", groups[0], responseBody.Id)
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
}
