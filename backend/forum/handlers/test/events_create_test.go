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
	cli2 := testServer.TestClient()

	group := CreateGroup("testGroup", "testDesc", []int{})
	cli1.TestAuth(t)
	cli2.TestAuth(t)

	t.Run("Invalid", func(t *testing.T) {

		cli1.TestPost(t, "/api/groups/create", group, http.StatusOK)

		t.Run("Method", func(t *testing.T) {
			cli1.TestGet(t, "/api/groups/1/events/create", http.StatusMethodNotAllowed)
		})

		t.Run("Not Found", func(t *testing.T) {
			cli1.TestPost(t, "/api/groups/99999999999999999999999999999999999999999999/events/create",
				nil, http.StatusNotFound)
		})

		t.Run("Forbidden", func(t *testing.T) {
			cli2.TestPost(t, "/api/groups/1/events/create",
				nil, http.StatusForbidden)
		})

		t.Run("Title - too short", func(t *testing.T) {
			invalidEventBody := createEventData("", "testdesc", "2024-01-01")
			cli1.TestPost(t, "/api/groups/1/events/create", invalidEventBody, http.StatusBadRequest)
		})
		t.Run("Title - too long", func(t *testing.T) {
			// 101 characters
			invalidEventBody := createEventData(strings.Repeat("a", 101), "test", "2024-01-01")

			cli1.TestPost(t, "/api/groups/1/events/create", invalidEventBody, http.StatusBadRequest)
		})

		t.Run("Description - too short", func(t *testing.T) {
			invalidEventBody := createEventData("test", "", "2024-01-01")
			cli1.TestPost(t, "/api/groups/1/events/create", invalidEventBody, http.StatusBadRequest)
		})

		t.Run("Description - too long", func(t *testing.T) {
			// 201 characters
			invalidEventBody := createEventData("test", strings.Repeat("a", 201), "2024-01-01")
			cli1.TestPost(t, "/api/groups/1/events/create", invalidEventBody, http.StatusBadRequest)
		})

		t.Run("Time and date", func(t *testing.T) {
			invalidEventBody := createEventData("test", "desc", "")
			cli1.TestPost(t, "/api/groups/1/events/create", invalidEventBody, http.StatusBadRequest)
		})

		t.Run("Time and date - past", func(t *testing.T) {
			invalidEventBody := createEventData("test", "desc", "1984-01-01")
			cli1.TestPost(t, "/api/groups/1/events/create", invalidEventBody, http.StatusBadRequest)
		})

		t.Run("Invalid body", func(t *testing.T) {
			cli1.TestPost(t, "/api/groups/1/events/create", "invalid", http.StatusBadRequest)
		})

	})

	t.Run("Valid", func(t *testing.T) {
		t.Run("Create event", func(t *testing.T) {
			eventBody := createEventData("test", "test", "2024-01-01")
			cli1.TestPost(t, "/api/groups/1/events/create", eventBody, http.StatusOK)
		})
	})
}
