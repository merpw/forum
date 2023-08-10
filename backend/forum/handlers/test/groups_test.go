package server_test

import (
	. "backend/forum/handlers/test/server"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/gofrs/uuid"
)

type Group struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Invite      []int  `json:"invite"`
}

func createGroup(title, desc string, invites []int) *Group {
	return &Group{
		Title:       title,
		Description: desc,
		Invite:      invites,
	}
}

type GroupPostData struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Content     string `json:"content"`

	// []string for requests, string for responses
	Categories interface{} `json:"categories"`
	GroupId    int         `json:"group_id"`
}

func groupPostData(groupId int) *GroupPostData {
	return &GroupPostData{
		Title:       uuid.Must(uuid.NewV4()).String()[0:8],
		Content:     "content",
		Description: "description",
		Categories:  []string{"facts"},
		GroupId:     groupId,
	}
}

func TestGroupsId(t *testing.T) {
	testServer := NewTestServer(t)
	cli1 := testServer.TestClient()
	cli1.TestAuth(t)

	group := createGroup("test", "test", []int{})
	// Create a group to test with
	cli1.TestPost(t, "/api/groups/create", group, http.StatusOK)

	cli2 := testServer.TestClient()

	t.Run("Invalid", func(t *testing.T) {
		t.Run("Unauthorized", func(t *testing.T) {
			cli2.TestGet(t, "/api/groups/1", http.StatusUnauthorized)
		})

		t.Run("Not Found", func(t *testing.T) {
			cli1.TestGet(t, "/api/groups/10", http.StatusNotFound)
			cli1.TestGet(t, "/api/groups/999999999999999999999999999", http.StatusNotFound)
		})
	})
}

func TestGroupsIdPosts(t *testing.T) {
	testServer := NewTestServer(t)
	cli := testServer.TestClient()

	t.Run("Invalid", func(t *testing.T) {
		t.Run("Unauthorized", func(t *testing.T) {
			cli.TestGet(t, "/api/groups/1", http.StatusUnauthorized)
		})
		cli.TestAuth(t)
		t.Run("Not found", func(t *testing.T) {
			cli.TestGet(t, "/api/groups/1/posts", http.StatusNotFound)
			cli.TestGet(t, "/api/groups/999999999999999999999999999999/posts", http.StatusNotFound)
		})
	})

	t.Run("Valid", func(t *testing.T) {

		// create a group, and create a post
		group := createGroup("test", "test", []int{})
		cli.TestPost(t, "/api/groups/create", group, http.StatusOK)

		post := groupPostData(1)
		cli.TestPost(t, "/api/posts/create", post, http.StatusOK)

		var postIds []int
		_, resp := cli.TestGet(t, "/api/groups/1/posts", http.StatusOK)

		err := json.Unmarshal(resp, &postIds)
		if err != nil {
			t.Fatal(err)
		}

		if len(postIds) != 1 {
			t.Errorf("invalid postIds length, expected %d, got %d", 1, len(postIds))
		}

		if postIds[0] != 1 {
			t.Errorf("invalid post id, expected %d, got %d", 1, postIds[0])
		}
	})
}

func TestGroupIdJoin(t *testing.T) {
	testServer := NewTestServer(t)
	cli1 := testServer.TestClient()
	cli2 := testServer.TestClient()
	cli1.TestAuth(t)
	cli2.TestAuth(t)

	group := createGroup("test", "test", []int{})
	cli1.TestPost(t, "/api/groups/create", group, http.StatusOK)

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

	})
}

func TestGroupIdInviteLeave(t *testing.T) {
	testServer := NewTestServer(t)
	cli1 := testServer.TestClient()
	cli2 := testServer.TestClient()
	cli1.TestAuth(t)
	cli2.TestAuth(t)

	group := createGroup("test", "test", []int{})
	cli1.TestPost(t, "/api/groups/create", group, http.StatusOK)
	invalidBody := struct {
		UserId int `json:"user_id"`
	}{
		UserId: 666,
	}
	t.Run("Not found", func(t *testing.T) {
		cli1.TestPost(t, "/api/groups/1/invite", invalidBody, http.StatusNotFound)
		cli1.TestPost(t, "/api/groups/10/leave", nil, http.StatusNotFound)
		cli1.TestPost(t, "/api/groups/10/invite", nil, http.StatusNotFound)
		cli1.TestPost(t, "/api/groups/999999999999999999999999999999/invite", nil, http.StatusNotFound)
		cli1.TestPost(t, "/api/groups/999999999999999999999999999999/leave", nil, http.StatusNotFound)
	})

	t.Run("Bad request", func(t *testing.T) {
		cli2.TestPost(t, "/api/groups/1/leave", nil, http.StatusBadRequest)
		cli2.TestPost(t, "/api/groups/1/join", nil, http.StatusOK)
		cli2.TestPost(t, "/api/groups/1/leave", nil, http.StatusBadRequest)
		cli1.TestPost(t, "/api/groups/1/invite", "invalid", http.StatusBadRequest)
	})

	t.Run("Valid", func(t *testing.T) {
		var status int
		requestBody := struct {
			UserId int `json:"user_id"`
		}{
			UserId: 2,
		}

		t.Run("Invite request", func(t *testing.T) {
			_, resp := cli1.TestPost(t, "/api/groups/1/invite", requestBody, http.StatusOK)
			if err := json.Unmarshal(resp, &status); err != nil {
				t.Fatal(err)
			}
			if status != 2 {
				t.Errorf("unexpected status, expected %d, got %d", 2, status)
			}
		})
		t.Run("Abort request", func(t *testing.T) {
			_, resp := cli2.TestPost(t, "/api/groups/1/invite", requestBody, http.StatusOK)
			if err := json.Unmarshal(resp, &status); err != nil {
				t.Fatal(err)
			}
			if status != 0 {
				t.Errorf("unexpected status, expected %d, got %d", 2, status)
			}
		})

		response := struct {
			Response bool `json:"response"`
		}{
			Response: true,
		}

		t.Run("Invite, respond, then leave", func(t *testing.T) {
			_, resp := cli1.TestPost(t, "/api/groups/1/invite", requestBody, http.StatusOK)
			if err := json.Unmarshal(resp, &status); err != nil {
				t.Fatal(err)
			}
			if status != 2 {
				t.Errorf("unexpected status, expected %d, got %d", 2, status)
			}

			var invites []int

			_, resp = cli2.TestGet(t, "/api/invitations", http.StatusOK)
			if err := json.Unmarshal(resp, &invites); err != nil {
				t.Fatal(err)
			}

			if len(invites) != 1 {
				t.Errorf("Invalid length of invites, expected %d, got %d", 1, len(invites))
			}

			endpoint := fmt.Sprintf("/api/invitations/%d/respond", invites[0])
			cli2.TestPost(t, endpoint, response, http.StatusOK)

		})

	})
}
