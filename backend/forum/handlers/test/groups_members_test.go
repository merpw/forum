package server_test

import (
	. "backend/forum/handlers/test/server"
	"encoding/json"
	"strconv"
	"testing"
)

func getId(t *testing.T, cli *TestClient) int {
	t.Helper()
	var resp struct {
		Id int `json:"id"`
	}
	_, body := cli.TestGet(t, "/api/me", 200)
	if err := json.Unmarshal(body, &resp); err != nil {
		t.Fatal(err)
	}
	return resp.Id
}

func getMembers(t *testing.T, cli *TestClient, groupId int) []int {
	t.Helper()
	var resp []int
	_, body := cli.TestGet(t, "/api/groups/"+strconv.Itoa(groupId)+"/members", 200)
	if err := json.Unmarshal(body, &resp); err != nil {
		t.Fatal(err)
	}
	return resp
}

func includes(arr []int, val int) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}

func acceptAllInvitations(t *testing.T, cli *TestClient) {
	t.Helper()
	var invitationIds []int
	_, body := cli.TestGet(t, "/api/invitations", 200)
	if err := json.Unmarshal(body, &invitationIds); err != nil {
		t.Fatal(err)
	}
	for _, id := range invitationIds {
		cli.TestPost(t, "/api/invitations/"+strconv.Itoa(id)+"/respond", struct {
			Response bool `json:"response"`
		}{Response: true}, 200)
	}
}

func TestGroupMembers(t *testing.T) {
	testServer := NewTestServer(t)
	cliCreator := testServer.TestClient()
	cliCreator.TestAuth(t)

	creatorId := getId(t, cliCreator)

	cliCreator.TestPost(t, "/api/groups/create", CreateGroup("test", "test", []int{}), 200)

	members := getMembers(t, cliCreator, 1)
	if !includes(members, creatorId) {
		t.Fatal("expected creator to be in the group members")
	}

	var cliMemberIds []int
	var cliMembers []*TestClient
	for i := 0; i < 10; i++ {
		cliMember := testServer.TestClient()
		cliMember.TestAuth(t)
		cliMemberIds = append(cliMemberIds, getId(t, cliMember))
		cliMembers = append(cliMembers, cliMember)
	}

	t.Run("Invalid", func(t *testing.T) {
		t.Run("No such group", func(t *testing.T) {
			cliCreator.TestGet(t, "/api/groups/2/members", 404)

			cliCreator.TestGet(t, "/api/groups/214748364712312214748364712312/members", 404)

		})

		t.Run("Not authorized", func(t *testing.T) {
			cliAnon := testServer.TestClient()
			cliAnon.TestGet(t, "/api/groups/1/members", 401)
		})

		t.Run("Not a member", func(t *testing.T) {
			cli := testServer.TestClient()
			cli.TestAuth(t)
			cli.TestGet(t, "/api/groups/1/members", 403)
		})
	})

	cliCreator.TestPost(t, "/api/groups/create", CreateGroup("test", "test", cliMemberIds), 200)
	for _, cliMember := range cliMembers {
		acceptAllInvitations(t, cliMember)
	}

	members = getMembers(t, cliCreator, 2)

	if len(members) != len(cliMemberIds)+1 {
		t.Fatalf("expected %d members, got %d", len(cliMemberIds), len(members))
	}

	for _, id := range cliMemberIds {
		if !includes(members, id) {
			t.Fatalf("expected member %d to be in the group members", id)
		}
	}

	cliMembers[0].TestPost(t, "/api/groups/2/leave", nil, 200)

	members = getMembers(t, cliCreator, 2)
	if includes(members, cliMemberIds[0]) {
		t.Fatal("expected member to not be in the group members")
	}

	t.Run("withPending", func(t *testing.T) {
		cli := testServer.TestClient()
		cli.TestAuth(t)

		cliId := getId(t, cli)

		cliCreator.TestPost(t, "/api/groups/create", CreateGroup("test", "test", []int{cliId}), 200)

		members = getMembers(t, cliCreator, 3)

		if includes(members, cliId) {
			t.Fatal("expected pending member to not be in the group members")
		}

		acceptAllInvitations(t, cli)

		members = getMembers(t, cliCreator, 3)

		if !includes(members, cliId) {
			t.Fatal("expected member to be in the group members")
		}
	})
}
