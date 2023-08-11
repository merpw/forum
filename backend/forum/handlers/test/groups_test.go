package server_test

import (
	. "backend/forum/handlers/test/server"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gofrs/uuid"
)

type Group struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Invite      []int  `json:"invite"`
}

func CreateGroup(title, desc string, invites []int) *Group {
	return &Group{
		Title:       title,
		Description: desc,
		Invite:      invites,
	}
}

type GroupPost struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Content     string `json:"content"`

	// []string for requests, string for responses
	Categories interface{} `json:"categories"`
	GroupId    int         `json:"group_id"`
}

func GroupPostData(groupId int) *GroupPost {
	return &GroupPost{
		Title:       uuid.Must(uuid.NewV4()).String()[0:8],
		Content:     "content",
		Description: "description",
		Categories:  []string{"facts"},
		GroupId:     groupId,
	}
}

func TestGroups(t *testing.T) {
	testServer := NewTestServer(t)
	cli := testServer.TestClient()
	cli.TestAuth(t)

	var clients []*TestClient

	// Create four more users
	for i := 0; i < 3; i++ {
		client := testServer.TestClient()
		client.TestAuth(t)
		clients = append(clients, client)
	}

	t.Run("Create groups", func(t *testing.T) {
		response := struct {
			Response bool `json:"response"`
		}{
			Response: true,
		}

		// Create group with 1 user
		group1 := CreateGroup("test", "test", []int{})
		cli.TestPost(t, "/api/groups/create", group1, http.StatusOK)

		// Create group with 2 users
		group2 := CreateGroup("test", "test", []int{2})
		cli.TestPost(t, "/api/groups/create", group2, http.StatusOK)
		clients[0].TestPost(t, "/api/invitations/1/respond", response, http.StatusOK)

		// Create group with 3 users
		group3 := CreateGroup("test", "test", []int{2, 3})
		cli.TestPost(t, "/api/groups/create", group3, http.StatusOK)
		clients[1].TestPost(t, "/api/invitations/2/respond", response, http.StatusOK)
		clients[2].TestPost(t, "/api/invitations/3/respond", response, http.StatusOK)
	})

	t.Run("Check order of groups", func(t *testing.T) {
		var groupIds []int
		_, resp := cli.TestGet(t, "/api/groups", http.StatusOK)

		if err := json.Unmarshal(resp, &groupIds); err != nil {
			t.Fatal(err)
		}

		var groupId = 3

		for _, id := range groupIds {
			if groupId != id {
				t.Errorf("unexpected id, expected %d, got %d", groupId, id)
			}
			groupId--
		}

	})

}

func TestGroupsId(t *testing.T) {
	testServer := NewTestServer(t)
	cli1 := testServer.TestClient()
	cli1.TestAuth(t)

	group := CreateGroup("test", "test", []int{})
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
		group := CreateGroup("test", "test", []int{})
		cli.TestPost(t, "/api/groups/create", group, http.StatusOK)

		post := GroupPostData(1)
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
