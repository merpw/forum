package server_test

import (
	. "backend/forum/handlers/test/server"
	"encoding/json"
	"net/http"
	"testing"
)

func TestUser(t *testing.T) {
	testServer := NewTestServer(t)

	cli1 := testServer.TestClient()
	cli1.TestAuth(t)

	cliAnon := testServer.TestClient()

	t.Run("Not found", func(t *testing.T) {
		cli1.TestGet(t, "/api/user/10", http.StatusNotFound)
		cli1.TestGet(t, "/api/user/214748364712312214748364712312", http.StatusNotFound)
	})

	t.Run("Valid", func(t *testing.T) {
		var responseData struct {
			Id   int    `json:"id"`
			Name string `json:"name"`
		}

		_, respBody := cli1.TestGet(t, "/api/user/1", http.StatusOK)

		checkRespBody := func() {
			err := json.Unmarshal(respBody, &responseData)
			if err != nil {
				t.Fatal(err)
			}

			if responseData.Id != 1 {
				t.Fatalf("invalid id, expected 1, got %d", responseData.Id)
			}

			if responseData.Name != cli1.Name {
				t.Fatalf("invalid name, expected %s, got %s", cli1.Name, responseData.Name)
			}
		}

		checkRespBody()

		_, respBody = cliAnon.TestGet(t, "/api/user/1", http.StatusOK)

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
		cli1.TestGet(t, "/api/user/10/posts", http.StatusNotFound)
		cli1.TestGet(t, "/api/user/214748364712312214748364712312/posts", http.StatusNotFound)
	})

	t.Run("Valid", func(t *testing.T) {
		createPosts(t, cli1, 5)

		_, respBody := cli1.TestGet(t, "/api/user/1/posts", http.StatusOK)

		var cli1Posts []PostData
		err := json.Unmarshal(respBody, &cli1Posts)
		if err != nil {
			t.Fatal(err)
		}

		if len(cli1Posts) != 5 {
			t.Fatalf("invalid posts count, expected 5, got %d", len(cli1Posts))
		}

		_, resp2Body := cli2.TestGet(t, "/api/user/1/posts", http.StatusOK)

		if string(respBody) != string(resp2Body) {
			t.Fatalf("responses mismatch, expected %s, got %s", string(respBody), string(resp2Body))
		}
	})
}
