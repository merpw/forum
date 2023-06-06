package server_test

import (
	. "backend/forum/handlers/test/server"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"testing"
)

type ReactionData struct {
	Reaction      int `json:"reaction"`
	LikesCount    int `json:"likes_count"`
	DislikesCount int `json:"dislikes_count"`
}

func TestReactionsPost(t *testing.T) {
	testServer := NewTestServer(t)

	cli := testServer.TestClient()

	t.Run("Unauthorized", func(t *testing.T) {
		cli.TestPost(t, "/api/posts/1/like", nil, http.StatusUnauthorized)
		cli.TestPost(t, "/api/posts/1/dislike", nil, http.StatusUnauthorized)
		cli.TestGet(t, "/api/posts/1/reaction", http.StatusUnauthorized)
	})

	cli.TestAuth(t)

	t.Run("Not found", func(t *testing.T) {
		cli.TestPost(t, "/api/posts/1/like", nil, http.StatusNotFound)
		cli.TestPost(t, "/api/posts/214748364712312214748364712312/like", nil, http.StatusNotFound)

		cli.TestPost(t, "/api/posts/1/dislike", nil, http.StatusNotFound)
		cli.TestPost(t, "/api/posts/214748364712312214748364712312/dislike", nil, http.StatusNotFound)

		cli.TestGet(t, "/api/posts/1/reaction", http.StatusNotFound)
		cli.TestGet(t, "/api/posts/214748364712312214748364712312/reaction", http.StatusNotFound)
	})

	t.Run("Like", func(t *testing.T) {
		post := createPosts(t, cli, 1)[0]

		_, respBody := cli.TestPost(t, fmt.Sprintf("/api/posts/%d/like", post.Id), nil, http.StatusOK)
		checkReaction(t, respBody, 1)

		_, respBody = cli.TestPost(t, fmt.Sprintf("/api/posts/%d/like", post.Id), nil, http.StatusOK)
		checkReaction(t, respBody, 0)
	})

	t.Run("Dislike", func(t *testing.T) {
		post := createPosts(t, cli, 1)[0]

		_, respBody := cli.TestPost(t, fmt.Sprintf("/api/posts/%d/dislike", post.Id), nil, http.StatusOK)
		checkReaction(t, respBody, -1)

		_, respBody = cli.TestPost(t, fmt.Sprintf("/api/posts/%d/dislike", post.Id), nil, http.StatusOK)
		checkReaction(t, respBody, 0)
	})

	t.Run("Like + dislike", func(t *testing.T) {
		post := createPosts(t, cli, 1)[0]

		_, respBody := cli.TestPost(t, fmt.Sprintf("/api/posts/%d/like", post.Id), nil, http.StatusOK)
		checkReaction(t, respBody, 1)

		_, respBody = cli.TestPost(t, fmt.Sprintf("/api/posts/%d/dislike", post.Id), nil, http.StatusOK)
		checkReaction(t, respBody, -1)
	})

	t.Run("Dislike + like", func(t *testing.T) {
		post := createPosts(t, cli, 1)[0]

		_, respBody := cli.TestPost(t, fmt.Sprintf("/api/posts/%d/dislike", post.Id), nil, http.StatusOK)
		checkReaction(t, respBody, -1)

		_, respBody = cli.TestPost(t, fmt.Sprintf("/api/posts/%d/like", post.Id), nil, http.StatusOK)
		checkReaction(t, respBody, 1)
	})
}

func checkReaction(t *testing.T, respBody []byte, expectedReaction int) {
	t.Helper()

	reaction, err := strconv.Atoi(string(respBody))
	if err != nil {
		t.Fatal(err)
	}

	if reaction != expectedReaction {
		t.Fatalf("Wrong reaction %+v", string(respBody))
	}
}

func TestReactionsGet(t *testing.T) {
	testServer := NewTestServer(t)
	creator := testServer.TestClient()
	creator.TestAuth(t)

	cli1 := testServer.TestClient()
	cli1.TestAuth(t)

	cli2 := testServer.TestClient()
	cli2.TestAuth(t)

	cli3 := testServer.TestClient()
	cli3.TestAuth(t)

	t.Run("/api/posts/{id}", func(t *testing.T) {
		post := createPosts(t, creator, 1)[0]

		cli1.TestPost(t, fmt.Sprintf("/api/posts/%d/like", post.Id), nil, http.StatusOK)
		cli2.TestPost(t, fmt.Sprintf("/api/posts/%d/like", post.Id), nil, http.StatusOK)

		cli3.TestPost(t, fmt.Sprintf("/api/posts/%d/dislike", post.Id), nil, http.StatusOK)

		_, respBody := creator.TestGet(t, fmt.Sprintf("/api/posts/%d", post.Id), http.StatusOK)
		reactionData := parseReactionData(t, respBody)

		if reactionData.LikesCount != 2 || reactionData.DislikesCount != 1 {
			t.Fatalf("Wrong reaction %+v", reactionData)
		}
	})

	t.Run("/api/posts/{id}/reaction", func(t *testing.T) {
		post := createPosts(t, creator, 1)[0]

		cli1.TestPost(t, fmt.Sprintf("/api/posts/%d/like", post.Id), nil, http.StatusOK)
		cli2.TestPost(t, fmt.Sprintf("/api/posts/%d/like", post.Id), nil, http.StatusOK)

		cli3.TestPost(t, fmt.Sprintf("/api/posts/%d/dislike", post.Id), nil, http.StatusOK)

		_, respBody := creator.TestGet(t, fmt.Sprintf("/api/posts/%d/reaction", post.Id), http.StatusOK)
		reactionData := parseReactionData(t, respBody)

		if reactionData.LikesCount != 2 || reactionData.DislikesCount != 1 || reactionData.Reaction != 0 {
			t.Fatalf("Wrong reaction %+v", reactionData)
		}

		_, respBody = cli1.TestGet(t, fmt.Sprintf("/api/posts/%d/reaction", post.Id), http.StatusOK)
		reactionData = parseReactionData(t, respBody)

		if reactionData.LikesCount != 2 || reactionData.DislikesCount != 1 || reactionData.Reaction != 1 {
			t.Fatalf("Wrong reaction %+v", reactionData)
		}
	})
}

func parseReactionData(t *testing.T, respBody []byte) ReactionData {
	t.Helper()

	var reaction ReactionData
	err := json.Unmarshal(respBody, &reaction)
	if err != nil {
		t.Fatal(err)
	}
	return reaction
}
