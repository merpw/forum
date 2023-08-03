package server_test

import (
	. "backend/forum/handlers/test/server"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func TestUser(t *testing.T) {
	testServer := NewTestServer(t)

	cli1 := testServer.TestClient()
	cli1.TestAuth(t)

	cliAnon := testServer.TestClient()

	t.Run("Not found", func(t *testing.T) {
		cli1.TestGet(t, "/api/users/10", http.StatusNotFound)
		cli1.TestGet(t, "/api/users/214748364712312214748364712312", http.StatusNotFound)
	})

	t.Run("Valid", func(t *testing.T) {
		var responseData struct {
			Id       int    `json:"id"`
			Username string `json:"username"`
			Avatar   string `json:"avatar"`
			Bio      string `json:"bio"`
		}

		_, respBody := cli1.TestGet(t, "/api/users/1", http.StatusOK)

		checkRespBody := func() {
			err := json.Unmarshal(respBody, &responseData)
			if err != nil {
				t.Fatal(err)
			}

			if responseData.Id != 1 {
				t.Fatalf("invalid id, expected 1, got %d", responseData.Id)
			}

			if responseData.Username != cli1.Username {
				t.Fatalf("invalid username, expected %s, got %s", cli1.Username, responseData.Username)
			}

			if responseData.Bio != cli1.Bio {
				t.Fatalf("invalid bio, expected %s, got %s", cli1.Bio, responseData.Bio)
			}

			if responseData.Avatar != cli1.Avatar {
				t.Fatalf("invalid avatar, expected %s, got %s", cli1.Avatar, responseData.Avatar)
			}
		}

		checkRespBody()

		_, respBody = cliAnon.TestGet(t, "/api/users/1", http.StatusOK)

		checkRespBody()
	})
}

func TestUserPosts(t *testing.T) {
	testServer := NewTestServer(t)

	cli1 := testServer.TestClient()
	cli1.TestAuth(t)

	cli2 := testServer.TestClient()
	cli2.TestAuth(t)

	t.Run("Not found", func(t *testing.T) {
		cli1.TestGet(t, "/api/users/10/posts", http.StatusNotFound)
		cli1.TestGet(t, "/api/users/214748364712312214748364712312/posts", http.StatusNotFound)
	})

	t.Run("Valid", func(t *testing.T) {
		createPosts(t, cli1, 5)

		_, respBody := cli1.TestGet(t, "/api/users/1/posts", http.StatusOK)

		var cli1Posts []PostData
		err := json.Unmarshal(respBody, &cli1Posts)
		if err != nil {
			t.Fatal(err)
		}

		if len(cli1Posts) != 5 {
			t.Fatalf("invalid posts count, expected 5, got %d", len(cli1Posts))
		}

		_, resp2Body := cli2.TestGet(t, "/api/users/1/posts", http.StatusOK)

		if string(respBody) != string(resp2Body) {
			t.Fatalf("responses mismatch, expected %s, got %s", string(respBody), string(resp2Body))
		}
	})
}

const UsersCount = 5

func TestUsers(t *testing.T) {
	testServer := NewTestServer(t)

	for i := 0; i < UsersCount; i++ {
		// create users with random names
		cli := testServer.TestClient()
		cli.TestClientData = NewClientData()
		if i%2 == 0 {
			cli.TestClientData.Username = strings.ToUpper(cli.TestClientData.Username)
		}
		cli.TestPost(t, "/api/signup", cli.TestClientData, http.StatusOK)
	}

	cliAnon := testServer.TestClient()

	_, respBody := cliAnon.TestGet(t, "/api/users", http.StatusOK)

	var userIds []int
	err := json.Unmarshal(respBody, &userIds)
	if err != nil {
		t.Fatal(err)
	}

	if len(userIds) != UsersCount {
		t.Fatalf("invalid users count, expected %d, got %d", UsersCount, len(userIds))
	}

	type User struct {
		Id       int    `json:"id"`
		Username string `json:"username"`
	}
	var users []User

	for _, id := range userIds {
		_, respBody := cliAnon.TestGet(t, fmt.Sprintf("/api/users/%d", id), http.StatusOK)

		var responseData User

		err := json.Unmarshal(respBody, &responseData)
		if err != nil {
			t.Fatal(err)
		}

		if responseData.Id != id {
			t.Fatalf("invalid id, expected %d, got %d", id, responseData.Id)
		}

		users = append(users, responseData)
	}

	for i := 0; i < len(users); i++ {
		for j := i + 1; j < len(users); j++ {
			// should be sorted by username
			if strings.ToLower(users[i].Username) > strings.ToLower(users[j].Username) {
				t.Fatalf("users are not sorted by username, expected %s < %s", users[i].Username, users[j].Username)
			}
		}
	}
}
