package server_test

import (
	. "backend/forum/handlers/test/server"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gofrs/uuid"
)

type TestGroup struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Invite      []int  `json:"invite"`

	MemberStatus int `json:"member_status"`
	MemberCount  int `json:"member_count"`
}

func CreateGroup(title, desc string, invites []int) *TestGroup {
	return &TestGroup{
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

	// catch when admin id is not 1
	testServer.TestClient().TestAuth(t)

	cliCreator := testServer.TestClient()
	cliCreator.TestAuth(t)

	group := CreateGroup("test", "test", []int{})
	// Create a group to test with
	cliCreator.TestPost(t, "/api/groups/create", group, http.StatusOK)

	t.Run("Invalid", func(t *testing.T) {
		cli := testServer.TestClient()
		t.Run("Unauthorized", func(t *testing.T) {
			cli.TestGet(t, "/api/groups/1", http.StatusUnauthorized)
		})

		t.Run("Not Found", func(t *testing.T) {
			cliCreator.TestGet(t, "/api/groups/10", http.StatusNotFound)
			cliCreator.TestGet(t, "/api/groups/999999999999999999999999999", http.StatusNotFound)
		})
	})

	t.Run("Valid", func(t *testing.T) {
		t.Run("Non_member", func(t *testing.T) {
			cli := testServer.TestClient()
			cli.TestAuth(t)

			var group Group
			_, resp := cli.TestGet(t, "/api/groups/1", http.StatusOK)

			if err := json.Unmarshal(resp, &group); err != nil {
				t.Fatal(err)
			}

			if group.Title != "test" {
				t.Errorf("invalid title, expected %s, got %s", "test", group.Title)
			}

			if group.Description != "test" {
				t.Errorf("invalid description, expected %s, got %s", "test", group.Description)
			}

			if group.MemberStatus != 0 {
				t.Errorf("invalid member status, expected %d, got %d", 0, group.MemberStatus)
			}
		})

		t.Run("Member", func(t *testing.T) {
			cli := testServer.TestClient()
			cli.TestAuth(t)

			var group Group

			cli.TestPost(t, "/api/groups/1/join", nil, http.StatusOK)

			_, resp := cli.TestGet(t, "/api/groups/1", http.StatusOK)
			if err := json.Unmarshal(resp, &group); err != nil {
				t.Fatal(err)
			}

			if group.MemberStatus != 2 {
				t.Errorf("invalid member status, expected 2 (pending), got %d", group.MemberStatus)
			}

			if group.MemberCount != 1 {
				t.Errorf("invalid member count, expected 1, got %d", group.MemberCount)
			}

			AcceptAllInvitations(t, cliCreator)

			_, resp = cli.TestGet(t, "/api/groups/1", http.StatusOK)
			if err := json.Unmarshal(resp, &group); err != nil {
				t.Fatal(err)
			}

			if group.MemberStatus != 1 {
				t.Errorf("invalid member status, expected 1 (member), got %d", group.MemberStatus)
			}

			if group.MemberCount != 2 {
				t.Errorf("invalid member count, expected 2, got %d", group.MemberCount)
			}
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
