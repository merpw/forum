package server_test

import (
	. "backend/forum/handlers/test/server"
	"fmt"
	"net/http"
	"testing"
)

func TestCommentReactions(t *testing.T) {
	testServer := NewTestServer(t)

	cli := testServer.TestClient()

	t.Run("Unauthorized", func(t *testing.T) {
		cli.TestPost(t, "/api/posts/1/comment/1/like", nil, http.StatusUnauthorized)
		cli.TestPost(t, "/api/posts/1/comment/1/dislike", nil, http.StatusUnauthorized)
		cli.TestGet(t, "/api/posts/1/comment/1/reaction", http.StatusUnauthorized)
	})

	cli.TestAuth(t)

	t.Run("Not found", func(t *testing.T) {
		post := createPosts(t, cli, 1)[0]
		cli.TestPost(t, fmt.Sprintf("/api/posts/%d/comment/1/like", post.Id), nil, http.StatusNotFound)
		cli.TestPost(t,
			fmt.Sprintf("/api/posts/%d/comment/214748364712312214748364712312/like", post.Id),
			nil, http.StatusNotFound)

		cli.TestPost(t, fmt.Sprintf("/api/posts/%d/comment/1/dislike", post.Id), nil, http.StatusNotFound)
		cli.TestPost(t,
			fmt.Sprintf("/api/posts/%d/comment/214748364712312214748364712312/dislike", post.Id),
			nil, http.StatusNotFound)

		cli.TestGet(t, fmt.Sprintf("/api/posts/%d/comment/1/reaction", post.Id), http.StatusNotFound)
		cli.TestGet(t,
			fmt.Sprintf("/api/posts/%d/comment/214748364712312214748364712312/reaction", post.Id),
			http.StatusNotFound)
	})

	t.Run("Like", func(t *testing.T) {
		post := createPosts(t, cli, 1)[0]
		comment := createComments(t, cli, post.Id, 1)[0]

		_, respBody := cli.TestPost(t, fmt.Sprintf("/api/posts/%d/comment/%d/like", post.Id, comment.Id), nil, http.StatusOK)
		checkReaction(t, respBody, 1)

		_, respBody = cli.TestPost(t, fmt.Sprintf("/api/posts/%d/comment/%d/like", post.Id, comment.Id), nil, http.StatusOK)
		checkReaction(t, respBody, 0)
	})

	t.Run("Dislike", func(t *testing.T) {
		post := createPosts(t, cli, 1)[0]
		comment := createComments(t, cli, post.Id, 1)[0]

		_, respBody := cli.TestPost(t,
			fmt.Sprintf("/api/posts/%d/comment/%d/dislike", post.Id, comment.Id),
			nil, http.StatusOK)
		checkReaction(t, respBody, -1)

		_, respBody = cli.TestPost(t,
			fmt.Sprintf("/api/posts/%d/comment/%d/dislike", post.Id, comment.Id),
			nil, http.StatusOK)
		checkReaction(t, respBody, 0)
	})

	t.Run("Like and dislike", func(t *testing.T) {
		post := createPosts(t, cli, 1)[0]
		comment := createComments(t, cli, post.Id, 1)[0]

		_, respBody := cli.TestPost(t,
			fmt.Sprintf("/api/posts/%d/comment/%d/like", post.Id, comment.Id),
			nil, http.StatusOK)
		checkReaction(t, respBody, 1)

		_, respBody = cli.TestPost(t,
			fmt.Sprintf("/api/posts/%d/comment/%d/dislike", post.Id, comment.Id),
			nil, http.StatusOK)
		checkReaction(t, respBody, -1)

		_, respBody = cli.TestPost(t,
			fmt.Sprintf("/api/posts/%d/comment/%d/like", post.Id, comment.Id),
			nil, http.StatusOK)
		checkReaction(t, respBody, 1)
	})

	t.Run("Dislike and like", func(t *testing.T) {
		post := createPosts(t, cli, 1)[0]
		comment := createComments(t, cli, post.Id, 1)[0]

		_, respBody := cli.TestPost(t,
			fmt.Sprintf("/api/posts/%d/comment/%d/dislike", post.Id, comment.Id),
			nil, http.StatusOK)
		checkReaction(t, respBody, -1)

		_, respBody = cli.TestPost(t,
			fmt.Sprintf("/api/posts/%d/comment/%d/like", post.Id, comment.Id),
			nil, http.StatusOK)
		checkReaction(t, respBody, 1)

		_, respBody = cli.TestPost(t,
			fmt.Sprintf("/api/posts/%d/comment/%d/dislike", post.Id, comment.Id),
			nil, http.StatusOK)
		checkReaction(t, respBody, -1)
	})

	t.Run("Reaction", func(t *testing.T) {
		post := createPosts(t, cli, 1)[0]
		comment := createComments(t, cli, post.Id, 1)[0]

		_, respBody := cli.TestGet(t, fmt.Sprintf("/api/posts/%d/comment/%d/reaction", post.Id, comment.Id), http.StatusOK)

		reactionData := parseReactionData(t, respBody)

		if reactionData.Reaction != 0 || reactionData.LikesCount != 0 || reactionData.DislikesCount != 0 {
			t.Fatalf("Wrong reaction %+v", reactionData)
		}

		cli.TestPost(t, fmt.Sprintf("/api/posts/%d/comment/%d/like", post.Id, comment.Id), nil, http.StatusOK)

		_, respBody = cli.TestGet(t, fmt.Sprintf("/api/posts/%d/comment/%d/reaction", post.Id, comment.Id), http.StatusOK)

		reactionData = parseReactionData(t, respBody)
		if reactionData.Reaction != 1 || reactionData.LikesCount != 1 || reactionData.DislikesCount != 0 {
			t.Fatalf("Wrong reaction %+v", reactionData)
		}

		cli.TestPost(t, fmt.Sprintf("/api/posts/%d/comment/%d/dislike", post.Id, comment.Id), nil, http.StatusOK)

		_, respBody = cli.TestGet(t, fmt.Sprintf("/api/posts/%d/comment/%d/reaction", post.Id, comment.Id), http.StatusOK)

		reactionData = parseReactionData(t, respBody)

		if reactionData.Reaction != -1 || reactionData.LikesCount != 0 || reactionData.DislikesCount != 1 {
			t.Fatalf("Wrong reaction %+v", reactionData)
		}

		// Multiple users
		cli2 := testServer.TestClient()
		cli2.TestAuth(t)

		cli2.TestPost(t, fmt.Sprintf("/api/posts/%d/comment/%d/like", post.Id, comment.Id), nil, http.StatusOK)

		_, respBody = cli.TestGet(t, fmt.Sprintf("/api/posts/%d/comment/%d/reaction", post.Id, comment.Id), http.StatusOK)

		reactionData = parseReactionData(t, respBody)

		if reactionData.Reaction != -1 || reactionData.LikesCount != 1 || reactionData.DislikesCount != 1 {
			t.Fatalf("Wrong reaction %+v", reactionData)
		}
	})
}
