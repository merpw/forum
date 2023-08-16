package server_test

import (
	. "backend/forum/handlers/test/server"
	"net/http"
	"os"
	"testing"
)

func TestPostsSSR(t *testing.T) {
	testServer := NewTestServer(t)
	ssr := testServer.TestClient()
	cli := testServer.TestClient()
	cli.TestAuth(t)

	createPosts(t, cli, 1)

	err := os.Setenv("FORUM_IS_PRIVATE", "true")
	if err != nil {
		t.Fatal(err)
	}

	ssr.TestGet(t, "/api/posts/1", http.StatusOK)
	ssr.TestGet(t, "/api/posts/10", http.StatusNotFound)

	err = os.Setenv("FORUM_IS_PRIVATE", "false")
	if err != nil {
		t.Fatal(err)
	}
}
