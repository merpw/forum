package server_test

import (
	. "backend/forum/handlers/test/server"
	"net/http"
	"os"
	"testing"
)

func TestPrivateMode(t *testing.T) {
	testServer := NewTestServer(t)

	cli := testServer.TestClient()

	t.Run("Disabled", func(t *testing.T) {
		cli.TestGet(t, "/api/posts/categories", http.StatusOK)
	})

	t.Run("Enabled", func(t *testing.T) {
		err := os.Setenv("FORUM_IS_PRIVATE", "true")
		if err != nil {
			t.Fatal(err)
		}
		cli.TestGet(t, "/api/posts/categories", http.StatusUnauthorized)

		err = os.Setenv("FORUM_IS_PRIVATE", "false")
		if err != nil {
			t.Fatal(err)
		}
		cli.TestGet(t, "/api/posts/categories", http.StatusOK)
	})

}
