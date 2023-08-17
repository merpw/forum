package server_test

import (
	. "backend/forum/handlers/test/server"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"
)

type TestEvent struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	TimeAndDate string `json:"time_and_date"`
}

func generateEventData() TestEvent {
	return TestEvent{
		Title:       "Title",
		Description: "Valid description",
		TimeAndDate: time.Now().Add(24 * time.Hour).Format(time.RFC3339)[0:16],
	}
}

func createEvent(t *testing.T, cli *TestClient, groupId int) (int, TestEvent) {
	t.Helper()
	testEvent := generateEventData()
	var eventId int
	_, resp := cli.TestPost(t, fmt.Sprintf("/api/groups/%d/events/create", groupId),
		testEvent, http.StatusOK)
	if err := json.Unmarshal(resp, &eventId); err != nil {
		t.Fatal(err)
	}
	return eventId, testEvent
}

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
			invalidEventBody := generateEventData()
			invalidEventBody.Title = ""
			cli1.TestPost(t, "/api/groups/1/events/create", invalidEventBody, http.StatusBadRequest)
		})
		t.Run("Title - too long", func(t *testing.T) {
			invalidEventBody := generateEventData()
			invalidEventBody.Title = strings.Repeat("a", 101)

			cli1.TestPost(t, "/api/groups/1/events/create", invalidEventBody, http.StatusBadRequest)
		})

		t.Run("Description - too short", func(t *testing.T) {
			invalidEventBody := generateEventData()
			invalidEventBody.Description = ""
			cli1.TestPost(t, "/api/groups/1/events/create", invalidEventBody, http.StatusBadRequest)
		})

		t.Run("Description - too long", func(t *testing.T) {
			// 201 characters
			invalidEventBody := generateEventData()
			invalidEventBody.Description = strings.Repeat("a", 201)
			cli1.TestPost(t, "/api/groups/1/events/create", invalidEventBody, http.StatusBadRequest)
		})

		t.Run("Time and date", func(t *testing.T) {
			invalidEventBody := generateEventData()
			invalidEventBody.TimeAndDate = ""
			cli1.TestPost(t, "/api/groups/1/events/create", invalidEventBody, http.StatusBadRequest)
		})

		t.Run("Time and date - past", func(t *testing.T) {
			invalidEventBody := generateEventData()
			invalidEventBody.TimeAndDate = "1984-01-01T10:00"
			cli1.TestPost(t, "/api/groups/1/events/create", invalidEventBody, http.StatusBadRequest)
		})

		t.Run("Invalid body", func(t *testing.T) {
			cli1.TestPost(t, "/api/groups/1/events/create", "invalid", http.StatusBadRequest)
		})

	})

	t.Run("Valid", func(t *testing.T) {
		t.Run("Create event", func(t *testing.T) {
			eventBody := generateEventData()
			cli1.TestPost(t, "/api/groups/1/events/create", eventBody, http.StatusOK)
		})
	})
}
