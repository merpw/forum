package server_test

import (
	. "backend/forum/handlers/test/server"
	"encoding/json"
	"net/http"
	"testing"
)

func TestEventsId(t *testing.T) {
	testServer := NewTestServer(t)
	cli1 := testServer.TestClient()

	cli1.TestAuth(t)

	group := CreateGroup("test", "test", []int{})

	cli1.TestPost(t, "/api/groups/create", group, http.StatusOK)

	createEvent(t, cli1, 1)

	invitationIdRespond(t, cli1, 1, true)

	cli1.TestGet(t, "/api/groups/1/events", http.StatusOK)
	cli1.TestGet(t, "/api/events/1", http.StatusOK)

	t.Run("events/id/users", func(t *testing.T) {
		var users []int
		_, resp := cli1.TestGet(t, "/api/events/1/members", http.StatusOK)

		if err := json.Unmarshal(resp, &users); err != nil {
			t.Fatal(err)
		}

		if len(users) != 1 {
			t.Errorf("unexpected users going, expected %d, got %d", 1, len(users))
		}
	})

	cli1.TestPost(t, "/api/events/1/leave", nil, http.StatusOK)

	t.Run("events/id/leave", func(t *testing.T) {
		var users []int

		_, resp := cli1.TestGet(t, "/api/events/1/members", http.StatusOK)

		if err := json.Unmarshal(resp, &users); err != nil {
			t.Fatal(err)
		}

		if len(users) != 0 {
			t.Errorf("unexpected users going, expected %d, got %d", 0, len(users))
		}
	})

	t.Run("events/id/going", func(t *testing.T) {
		cli1.TestPost(t, "/api/events/1/going", nil, http.StatusOK)
		var users []int
		_, resp := cli1.TestGet(t, "/api/events/1/members", http.StatusOK)

		if err := json.Unmarshal(resp, &users); err != nil {
			t.Fatal(err)
		}

		if len(users) != 1 {
			t.Errorf("unexpected users going, expected %d, got %d", 2, len(users))
		}
	})
}

func TestEventsIdErrors(t *testing.T) {
	testServer := NewTestServer(t)
	cli1 := testServer.TestClient()
	cli2 := testServer.TestClient()

	cli1.TestAuth(t)
	cli2.TestAuth(t)

	group := CreateGroup("test", "test", []int{})

	cli1.TestPost(t, "/api/groups/create", group, http.StatusOK)

	createEvent(t, cli1, 1)

	invitationIdRespond(t, cli1, 1, true)

	t.Run("/groups/id/events", func(t *testing.T) {
		cli1.TestGet(t, "/api/groups/9999999999999999999999999999999999/events", http.StatusNotFound)
		cli1.TestGet(t, "/api/groups/2/events", http.StatusNotFound)
		cli2.TestGet(t, "/api/groups/1/events", http.StatusForbidden)
	})

	t.Run(".../events/id", func(t *testing.T) {
		cli1.TestGet(t, "/api/events/999999999999999999999999999999999", http.StatusNotFound)
		cli1.TestGet(t, "/api/events/2", http.StatusNotFound)
		cli2.TestGet(t, "/api/events/1", http.StatusForbidden)
	})

	t.Run(".../events/id/members", func(t *testing.T) {
		cli1.TestGet(t, "/api/events/999999999999999999999999999999999/members", http.StatusNotFound)
		cli1.TestGet(t, "/api/events/2/members", http.StatusNotFound)
		cli2.TestGet(t, "/api/events/1/members", http.StatusForbidden)
	})

	t.Run(".../events/id/leave", func(t *testing.T) {
		cli1.TestPost(t, "/api/events/999999999999999999999999999999999/leave",
			nil, http.StatusNotFound)
		cli1.TestPost(t, "/api/events/2/leave", nil, http.StatusNotFound)
		cli2.TestPost(t, "/api/events/1/leave", nil, http.StatusForbidden)
	})

	t.Run(".../events/id/going", func(t *testing.T) {
		cli1.TestPost(t, "/api/events/999999999999999999999999999999999/going",
			nil, http.StatusNotFound)
		cli1.TestPost(t, "/api/events/2/going", nil, http.StatusNotFound)
		cli2.TestPost(t, "/api/events/1/going", nil, http.StatusForbidden)
	})
}
