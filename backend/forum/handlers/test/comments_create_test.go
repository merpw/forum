package server_test

import (
	. "backend/forum/handlers/test/server"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
)

func TestCommentCreate(t *testing.T) {
	testServer := NewTestServer(t)

	cli := testServer.TestClient()
	cli2 := testServer.TestClient()

	cli2.TestAuth(t)
	post := createPosts(t, cli2, 1)[0]

	t.Run("Unauthorized", func(t *testing.T) {
		cli.TestPost(t, "/api/posts/1/comment", nil, http.StatusUnauthorized)
	})

	cli.TestAuth(t)

	t.Run("Not found", func(t *testing.T) {
		cli.TestPost(t, "/api/posts/2/comment", nil, http.StatusNotFound)
		cli.TestPost(t, "/api/posts/214748364712312214748364712312/comment", nil, http.StatusNotFound)
	})

	t.Run("Valid", func(t *testing.T) {

		cli.TestPost(t, fmt.Sprintf("/api/posts/%d/comment", post.Id), generateComment(), http.StatusOK)
	})

	t.Run("Invalid", func(t *testing.T) {
		goodComment := generateComment()

		t.Run("Body", func(t *testing.T) {
			cli.TestPost(t, fmt.Sprintf("/api/posts/%d/comment", post.Id), nil, http.StatusBadRequest)
			cli.TestPost(t, fmt.Sprintf("/api/posts/%d/comment", post.Id), "bad", http.StatusBadRequest)
			cli.TestPost(t, fmt.Sprintf("/api/posts/%d/comment", post.Id), struct {
				Content interface{} `json:"content"`
			}{
				Content: []int{1, 2, 3},
			}, http.StatusBadRequest)
		})

		t.Run("Content", func(t *testing.T) {
			goodComment.Content = ""
			cli.TestPost(t, fmt.Sprintf("/api/posts/%d/comment", post.Id), goodComment, http.StatusBadRequest)

			goodComment.Content = strings.Repeat("a", 10001)
			cli.TestPost(t, fmt.Sprintf("/api/posts/%d/comment", post.Id), goodComment, http.StatusBadRequest)
		})
	})
}

type CommentData struct {
	Id      int    `json:"id"`
	Content string `json:"content"`
	Author  struct {
		Id       int    `json:"id"`
		Username string `json:"username"`
	} `json:"author"`
	Date          string `json:"date"`
	LikesCount    int    `json:"likes_count"`
	DislikesCount int    `json:"dislikes_count"`
}

func generateComment() CommentData {
	return CommentData{
		Content: uuid.Must(uuid.NewV4()).String(),
	}
}

func createComments(t testing.TB, cli *TestClient, postId, count int) (comments []CommentData) {
	t.Helper()

	comments = make([]CommentData, count)

	for i := range comments {
		comments[i] = generateComment()

		_, respBody := cli.TestPost(t, fmt.Sprintf("/api/posts/%d/comment", postId), comments[i], http.StatusOK)
		commentId, err := strconv.Atoi(string(respBody))
		if err != nil {
			t.Fatal(err)
		}

		comments[i].Id = commentId
	}

	return comments
}
