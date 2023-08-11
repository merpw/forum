package server_test

import (
	"backend/forum/handlers/test/server"
	"net/http"
	"testing"
)

type TestEvent struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	TimeAndDate string `json:"time_and_date"`
}

func CreateEvent(title, desc string, timeAndDate string) *TestEvent {
	return &TestEvent{
		Title:       title,
		Description: desc,
		TimeAndDate: timeAndDate,
	}
}

func TestEvents(t *testing.T) {
	testServer := server.NewTestServer(t)
	cli := testServer.TestClient()
	cli.TestAuth(t)

	t.Run("Create event", func(t *testing.T) {
		group := CreateGroup("test", "test", []int{})
		cli.TestPost(t, "/api/groups/create", group, http.StatusOK)
		event := CreateEvent("test title", "test desc", "2024-01-01")
		cli.TestPost(t, "/api/groups/1/events/create", event, http.StatusOK)
	})

}
