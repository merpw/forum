package server_test

import (
	. "backend/forum/handlers/test/server"
	"encoding/json"
	"log"
	"net/http"
	"testing"
)

func TestInvitationsAll(t *testing.T) {
	testServer := NewTestServer(t)
	cli1 := testServer.TestClient()
	cli2 := testServer.TestClient()
	cli1.TestAuth(t)
	t.Run("Unauthorized", func(t *testing.T) {
		cli2.TestGet(t, "/api/invitations", http.StatusUnauthorized)
	})

	t.Run("Valid", func(t *testing.T) {
		cli2.TestAuth(t)
		var invitations []int

		t.Run("All invitations", func(t *testing.T) {
			cli2.TestPost(t, `/api/users/1/follow`, nil, http.StatusOK)

			_, response := cli1.TestGet(t, "/api/invitations", http.StatusOK)

			err := json.Unmarshal(response, &invitations)
			if err != nil {
				t.Fatal(err)
			}

			if len(invitations) != 1 {
				t.Errorf("invalid invitations amount, expected %d, got %d", 1, len(invitations))
			}
		})
	})
}

func TestInvitationsId(t *testing.T) {
	testServer := NewTestServer(t)
	cli1 := testServer.TestClient()
	cli2 := testServer.TestClient()
	cli1.TestAuth(t)
	cli2.TestAuth(t)

	t.Run("Invalid", func(t *testing.T) {
		t.Run("Not found", func(t *testing.T) {
			cli1.TestGet(t, "/api/invitations/10", http.StatusNotFound)
			cli1.TestGet(t, "/api/invitations/214748364712312214748364712312", http.StatusNotFound)
		})
	})

	t.Run("Valid", func(t *testing.T) {
		respBody := struct {
			Id         int    `json:"id"`
			Type       int    `json:"type"`
			FromUserId int    `json:"from_user_id"`
			ToUserId   int    `json:"to_user_id"`
			TimeStamp  string `json:"timestamp"`
		}{}

		var invitations []int
		t.Run("Request to follow", func(t *testing.T) {
			cli2.TestPost(t, "/api/users/1/follow", nil, http.StatusOK)
			_, resp := cli1.TestGet(t, "/api/invitations/1", http.StatusOK)
			err := json.Unmarshal(resp, &respBody)
			if err != nil {
				t.Fatal(err)
			}

			if respBody.Id != 1 {
				t.Errorf("invalid invitation id, expected %d, got %d", 1, respBody.Id)
			}

			if respBody.Type != 0 {
				t.Errorf("invalid type, expected %d, got %d", 0, respBody.Type)
			}

			if respBody.FromUserId != 2 {
				t.Errorf("invalid from_user_id, expected %d, got %d", 2, respBody.FromUserId)
			}

			if respBody.ToUserId != 1 {
				t.Errorf("invalid to_user_id, expected %d, got %d", 1, respBody.ToUserId)
			}
		})
		t.Run("Revoke follow request", func(t *testing.T) {
			var status int

			_, resp := cli2.TestPost(t, `/api/users/1/follow`, nil, http.StatusOK)
			if err := json.Unmarshal(resp, &status); err != nil {
				log.Panic(err)
			}

			if status != 0 {
				t.Errorf("expected %d, got %d", 0, status)
			}

			_, response := cli1.TestGet(t, "/api/invitations", http.StatusOK)

			err := json.Unmarshal(response, &invitations)
			if err != nil {
				t.Fatal(err)
			}

			if len(invitations) != 0 {
				t.Errorf("invalid invitations amount, expected %d, got %d", 0, len(invitations))
			}
		})
	})

}

func TestInvitationsIdRespond(t *testing.T) {

	t.Run("Invalid", func(t *testing.T) {
		testServer := NewTestServer(t)
		cli1 := testServer.TestClient()
		cli2 := testServer.TestClient()
		cli1.TestAuth(t)
		cli2.TestAuth(t)
		t.Run("Not found", func(t *testing.T) {
			cli1.TestPost(t,
				"/api/invitations/100/respond",
				nil,
				http.StatusNotFound)
			cli1.TestPost(t,
				"/api/invitations/214748364712312214748364712312/respond",
				nil,
				http.StatusNotFound)
		})
		t.Run("Bad Request", func(t *testing.T) {
			cli2.TestPost(t, "/api/users/1/follow", nil, http.StatusOK)
			cli1.TestPost(t,
				"/api/invitations/1/respond",
				"invalid",
				http.StatusBadRequest)
		})

		t.Run("Valid", func(t *testing.T) {
			testServer := NewTestServer(t)
			cli1 := testServer.TestClient()
			cli2 := testServer.TestClient()
			cli1.TestAuth(t)
			cli2.TestAuth(t)

			responseBody := struct {
				Response bool `json:"response"`
			}{
				Response: false,
			}

			t.Run("Reject invite", func(t *testing.T) {
				cli2.TestPost(t, "/api/users/1/follow", nil, http.StatusOK)
				cli1.TestPost(t,
					"/api/invitations/1/respond",
					responseBody,
					http.StatusOK)
			})

			t.Run("Accept invite", func(t *testing.T) {
				cli2.TestPost(t, "/api/users/1/follow", nil, http.StatusOK)
				responseBody.Response = true
				cli1.TestPost(t,
					"/api/invitations/2/respond",
					responseBody,
					http.StatusOK)
			})

		})
	})
}
