package server_test

import (
	. "backend/forum/database"
	. "backend/forum/handlers/test/server"
	"encoding/json"
	"net/http"
	"testing"
)

func TestPrivatePost(t *testing.T) {
	testServer := NewTestServer(t)
	cli1 := testServer.TestClient()
	cli2 := testServer.TestClient()
	cli1.TestAuth(t)
	cli2.TestAuth(t)

	validPost := generatePostData()

	validPost.Privacy = int(Private)
	accept := struct {
		Response bool `json:"response"`
	}{
		Response: true,
	}

	cli1.TestPost(t, "/api/posts/create", validPost, http.StatusOK)

	cli2.TestPost(t, "/api/users/1/follow", nil, http.StatusOK)

	cli1.TestPost(t, "/api/invitations/1/respond", accept, http.StatusOK)

	t.Run("Private post", func(t *testing.T) {
		var post Post
		var posts []Post

		_, resp := cli2.TestGet(t, "/api/posts/1", http.StatusOK)
		if err := json.Unmarshal(resp, &post); err != nil {
			t.Fatal(err)
		}

		_, resp = cli2.TestGet(t, "/api/posts", http.StatusOK)
		if err := json.Unmarshal(resp, &posts); err != nil {
			t.Fatal(err)
		}

		if post.Id != posts[0].Id {
			t.Errorf("invalid post-id match, expected %t, got %t", true, false)
		}

		t.Run("Comments", func(t *testing.T) {
			comment := generateComment()
			cli2.TestPost(t, "/api/posts/1/comment", comment, http.StatusOK)

			t.Run("Like", func(t *testing.T) {
				cli2.TestPost(t, "/api/posts/1/comment/1/like", nil, http.StatusOK)
			})

			t.Run("Remove like", func(t *testing.T) {
				cli2.TestPost(t, "/api/posts/1/comment/1/like", nil, http.StatusOK)
			})

			t.Run("Dislike", func(t *testing.T) {
				cli2.TestPost(t, "/api/posts/1/comment/1/dislike", nil, http.StatusOK)
			})

			t.Run("Remove dislike", func(t *testing.T) {
				cli2.TestPost(t, "/api/posts/1/comment/1/dislike", nil, http.StatusOK)
			})

			t.Run("Reaction", func(t *testing.T) {
				cli2.TestGet(t, "/api/posts/1/comment/1/reaction", http.StatusOK)
			})

		})

		t.Run("Test me posts liked", func(t *testing.T) {
			cli2.TestPost(t, "/api/posts/1/like", nil, http.StatusOK)

			var posts []PostData
			_, respData := cli2.TestGet(t, "/api/me/posts/liked", http.StatusOK)
			err := json.Unmarshal(respData, &posts)
			if err != nil {
				t.Fatal(err)
			}

			if len(posts) != 1 {
				t.Errorf("unexpected amount of liked posts, expected %d, got %d", 1, len(posts))
			}

			// unfollow
			cli2.TestPost(t, "/api/users/1/follow", nil, http.StatusOK)

			_, respData = cli2.TestGet(t, "/api/me/posts/liked", http.StatusOK)

			err = json.Unmarshal(respData, &posts)
			if err != nil {
				t.Fatal(err)
			}

			if len(posts) != 0 {
				t.Errorf("unexpected amount of liked posts, expected %d, got %d", 0, len(posts))
			}

		})

	})

	cli2.TestPost(t, "/api/users/1/follow", nil, http.StatusOK)

	cli1.TestPost(t, "/api/invitations/2/respond", accept, http.StatusOK)

	t.Run("SuperPrivate post", func(t *testing.T) {

	})
}

func TestSuperPrivatePost(t *testing.T) {
	testServer := NewTestServer(t)
	cli1 := testServer.TestClient()
	cli2 := testServer.TestClient()
	cli1.TestAuth(t)
	cli2.TestAuth(t)

	validPost := generatePostData()
	validPost.Privacy = int(SuperPrivate)
	validPost.Audience = []int{2}

	followAndRespond(t, cli1, cli2, 1, 1, true)

	cli1.TestPost(t, "/api/posts/create", validPost, http.StatusOK)

	var post Post
	var posts []Post

	_, resp := cli2.TestGet(t, "/api/posts/1", http.StatusOK)
	if err := json.Unmarshal(resp, &post); err != nil {
		t.Fatal(err)
	}

	_, resp = cli2.TestGet(t, "/api/posts", http.StatusOK)
	if err := json.Unmarshal(resp, &posts); err != nil {
		t.Fatal(err)
	}

	if post.Id != posts[0].Id {
		t.Errorf("invalid post-id match, expected %t, got %t", true, false)
	}

	t.Run("Comments", func(t *testing.T) {
		comment := generateComment()
		cli2.TestPost(t, "/api/posts/1/comment", comment, http.StatusOK)

		t.Run("Like", func(t *testing.T) {
			cli2.TestPost(t, "/api/posts/1/comment/1/like", nil, http.StatusOK)
		})

		t.Run("Dislike", func(t *testing.T) {
			cli2.TestPost(t, "/api/posts/1/comment/1/dislike", nil, http.StatusOK)
		})

		t.Run("Reaction", func(t *testing.T) {
			cli2.TestGet(t, "/api/posts/1/comment/1/reaction", http.StatusOK)
		})
	})

}

func TestInvalidPostPrivacy(t *testing.T) {
	testServer := NewTestServer(t)

	cli1 := testServer.TestClient()
	cli2 := testServer.TestClient()

	cli1.TestAuth(t)

	t.Run("Wack privacy", func(t *testing.T) {
		invalidPost := generatePostData()
		invalidPost.Privacy = 666
		cli1.TestPost(t, "/api/posts/create", invalidPost, http.StatusBadRequest)
	})

	cli2.TestAuth(t)

	t.Run("Super private", func(t *testing.T) {

		invalidPost := generatePostData()
		invalidPost.Privacy = int(SuperPrivate)

		t.Run("No followers", func(t *testing.T) {
			cli1.TestPost(t, "/api/posts/create", invalidPost, http.StatusBadRequest)
		})

		invalidPost.Audience = []int{1}
		t.Run("Add self to followers", func(t *testing.T) {
			cli1.TestPost(t, "/api/posts/create", invalidPost, http.StatusBadRequest)
		})

		invalidPost.Audience = []int{2}
		t.Run("Add non-follower to followers", func(t *testing.T) {
			cli1.TestPost(t, "/api/posts/create", invalidPost, http.StatusBadRequest)
		})

	})

	t.Run("Private post", func(t *testing.T) {

		t.Run("Add follower to post", func(t *testing.T) {
			invalidPost := generatePostData()
			invalidPost.Privacy = int(Private)
			invalidPost.Audience = []int{2}
			cli1.TestPost(t, "/api/posts/create", invalidPost, http.StatusBadRequest)
		})

		invalidPost := generatePostData()
		invalidPost.Privacy = int(Private)

		cli1.TestPost(t, "/api/posts/create", invalidPost, http.StatusOK)

		t.Run("Post reactions and comment", func(t *testing.T) {

			t.Run("Like", func(t *testing.T) {
				cli2.TestPost(t, "/api/posts/1/like", nil, http.StatusNotFound)
			})

			t.Run("Dislike", func(t *testing.T) {
				cli2.TestPost(t, "/api/posts/1/dislike", nil, http.StatusNotFound)
			})

			t.Run("Comment", func(t *testing.T) {
				cli2.TestPost(t, "/api/posts/1/comment", nil, http.StatusNotFound)
			})

			t.Run("Get comments", func(t *testing.T) {
				cli2.TestGet(t, "/api/posts/1/comments", http.StatusNotFound)
			})
		})
	})

	validPost := generatePostData()
	validPost.Privacy = int(Private)

	comment := generateComment()

	cli1.TestPost(t, "/api/posts/create", validPost, http.StatusOK)
	cli1.TestPost(t, "/api/posts/1/comment", comment, http.StatusOK)

	t.Run("Private post/id/comments", func(t *testing.T) {

		t.Run("Like", func(t *testing.T) {
			cli2.TestPost(t, "/api/posts/1/comment/1/like", nil, http.StatusBadRequest)
		})

		t.Run("Dislike", func(t *testing.T) {
			cli2.TestPost(t, "/api/posts/1/comment/1/dislike", nil, http.StatusBadRequest)
		})

		t.Run("Reaction", func(t *testing.T) {
			cli2.TestGet(t, "/api/posts/1/comment/1/reaction", http.StatusBadRequest)
		})
	})

	cli3 := testServer.TestClient()
	cli3.TestAuth(t)

	validPost.Privacy = int(SuperPrivate)
	validPost.Audience = []int{2}

	followAndRespond(t, cli1, cli2, 1, 1, true)

	t.Run("SuperPrivate comments", func(t *testing.T) {

		cli1.TestPost(t, "/api/posts/create", validPost, http.StatusOK)
		comment := createComments(t, cli1, 1, 1)[0]

		cli1.TestPost(t, "/api/posts/2/comment", comment, http.StatusOK)

		t.Run("Like", func(t *testing.T) {
			cli3.TestPost(t, "/api/posts/2/comment/2/like", nil, http.StatusBadRequest)
		})

		t.Run("Dislike", func(t *testing.T) {
			cli3.TestPost(t, "/api/posts/2/comment/2/dislike", nil, http.StatusBadRequest)
		})

		t.Run("Reaction", func(t *testing.T) {
			cli3.TestGet(t, "/api/posts/2/comment/2/reaction", http.StatusBadRequest)
		})

	})

}
