package server_test

import (
	. "backend/forum/handlers/test/server"
	"encoding/json"
	"net/http"
	"testing"
)

func TestComment(t *testing.T) {
	testServer := NewTestServer(t)

	cli := testServer.TestClient()

	cli.TestAuth(t)

	t.Run("Not found", func(t *testing.T) {
		cli.TestGet(t, "/api/posts/1/comments", http.StatusNotFound)
		cli.TestGet(t, "/api/posts/214748364712312214748364712312/comments", http.StatusNotFound)
	})

	createPosts(t, cli, 1)

	var comments []CommentData
	_, respBody := cli.TestGet(t, "/api/posts/1/comments", http.StatusOK)
	err := json.Unmarshal(respBody, &comments)
	if err != nil {
		t.Fatal(err)
	}

	if len(comments) != 0 {
		t.Fatal("Expected 0 comments")
	}

	requestedComments := createComments(t, cli, 1, 10)

	_, respBody = cli.TestGet(t, "/api/posts/1/comments", http.StatusOK)
	err = json.Unmarshal(respBody, &comments)
	if err != nil {
		t.Fatal(err)
	}
	if len(comments) != len(requestedComments) {
		t.Fatalf("Expected %d comments, got %d", len(requestedComments), len(comments))
	}

	for i := range comments {
		if comments[i].Content != requestedComments[i].Content {
			t.Fatalf("Expected comment content %s, got %s", requestedComments[i].Content, comments[i].Content)
		}
	}

	cli.TestGet(t, "/api/posts/1/comments", http.StatusOK)
}
