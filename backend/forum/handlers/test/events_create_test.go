package server_test

import (
	. "backend/forum/handlers/test/server"
	"net/http"
	"strings"
	"testing"
)

func TestCreateEvent(t *testing.T) {
	testServer := NewTestServer(t)

	cli1 := testServer.TestClient()

	group := CreateGroup("testGroup", "testDesc", []int{})

	t.Run("Invalid", func(t *testing.T) {

		t.Run("Unauthorized", func(t *testing.T) {
			cli1.TestPost(t, "/api/groups/1/events/create", nil, http.StatusUnauthorized)
		})

		cli1.TestAuth(t)

		cli1.TestPost(t, "/api/groups/create", group, http.StatusOK)
		t.Run("Method", func(t *testing.T) {
			cli1.TestGet(t, "/api/groups/1/events/create", http.StatusMethodNotAllowed)
		})

		t.Run("Title - too short", func(t *testing.T) {
			invalidEventBody := CreateEvent("", "testdesc", "2024-01-01")
			cli1.TestPost(t, "/api/groups/1/events/create", invalidEventBody, http.StatusBadRequest)
		})
		t.Run("Title - too long", func(t *testing.T) {
			// 101 characters
			invalidEventBody := CreateEvent(strings.Repeat("a", 101), "test", "2024-01-01")

			cli1.TestPost(t, "/api/groups/1/events/create", invalidEventBody, http.StatusBadRequest)
		})

		t.Run("Description - too short", func(t *testing.T) {
			invalidEventBody := CreateEvent("test", "", "2024-01-01")
			cli1.TestPost(t, "/api/groups/1/events/create", invalidEventBody, http.StatusBadRequest)
		})

		t.Run("Description - too long", func(t *testing.T) {
			// 201 characters
			invalidEventBody := CreateEvent("test", strings.Repeat("a", 201), "2024-01-01")
			cli1.TestPost(t, "/api/groups/1/events/create", invalidEventBody, http.StatusBadRequest)
		})

		t.Run("Time and date", func(t *testing.T) {
			invalidEventBody := CreateEvent("test", "desc", "")
			cli1.TestPost(t, "/api/groups/1/events/create", invalidEventBody, http.StatusBadRequest)
		})

		t.Run("Time and date - past", func(t *testing.T) {
			invalidEventBody := CreateEvent("test", "desc", "1984-01-01")
			cli1.TestPost(t, "/api/groups/1/events/create", invalidEventBody, http.StatusBadRequest)
		})

		t.Run("Invalid body", func(t *testing.T) {
			cli1.TestPost(t, "/api/groups/1/events/create", "invalid", http.StatusBadRequest)
		})

	})

	t.Run("Valid", func(t *testing.T) {
		cli1.TestAuth(t)
		t.Run("Create event", func(t *testing.T) {
			eventBody := CreateEvent("test", "test", "2024-01-01")
			cli1.TestPost(t, "/api/groups/1/events/create", eventBody, http.StatusOK)
		})
	})
}
