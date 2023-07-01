package server_test

import (
	. "backend/forum/handlers/test/server"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

func TestMe(t *testing.T) {
	testServer := NewTestServer(t)

	t.Run("Unauthorized", func(t *testing.T) {
		cli := testServer.TestClient()

		cli.TestGet(t, "/api/me", http.StatusUnauthorized)

		cli.TestAuth(t)
		cli2 := testServer.TestClient()
		cli2.Cookies = cli.Cookies

		cli.TestPost(t, "/api/logout", nil, http.StatusOK)

		// logged out using another client
		cli2.TestGet(t, "/api/me", http.StatusUnauthorized)
	})

	t.Run("Valid", func(t *testing.T) {
		cli := testServer.TestClient()
		cli.TestAuth(t)

		_, respBody := cli.TestGet(t, "/api/me", http.StatusOK)

		var respData struct {
			Id        int    `json:"id"`
			Username  string `json:"username"`
			Email     string `json:"email"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
		}
		err := json.Unmarshal(respBody, &respData)
		if err != nil {
			t.Fatal(err)
		}

		if respData.Username != cli.Username {
			t.Fatalf("invalid username, expected %s, got %s", cli.Username, respData.Username)
		}

		if respData.Email != cli.Email {
			t.Fatalf("invalid email, expected %s, got %s", cli.Email, respData.Email)
		}

		if respData.FirstName != cli.FirstName {
			t.Fatalf("invalid first name, expected %s, got %s", cli.FirstName, respData.FirstName)
		}

		if respData.LastName != cli.LastName {
			t.Fatalf("invalid last name, expected %s, got %s", cli.LastName, respData.LastName)
		}
	})
}

func TestMePosts(t *testing.T) {
	testServer := NewTestServer(t)

	cli := testServer.TestClient()

	cli.TestGet(t, "/api/me/posts", http.StatusUnauthorized)

	cli.TestAuth(t)

	createPosts(t, cli, 10)

	_, respData := cli.TestGet(t, "/api/me/posts", http.StatusOK)

	var posts []PostData
	err := json.Unmarshal(respData, &posts)
	if err != nil {
		t.Fatal(err)
	}

	if len(posts) != 10 {
		t.Fatalf("invalid number of posts, expected 10, got %d", len(posts))
	}
}

func TestMeLikedPosts(t *testing.T) {
	testServer := NewTestServer(t)

	creator := testServer.TestClient()

	t.Run("Unauthorized", func(t *testing.T) {
		creator.TestGet(t, "/api/me/posts/liked", http.StatusUnauthorized)
	})

	creator.TestAuth(t)

	posts := createPosts(t, creator, 10)

	cli := testServer.TestClient()
	cli.TestAuth(t)

	for _, post := range posts {
		cli.TestPost(t, fmt.Sprintf("/api/posts/%d/like", post.Id), nil, http.StatusOK)
	}

	_, respData := cli.TestGet(t, "/api/me/posts/liked", http.StatusOK)
	var likedPosts []PostData
	err := json.Unmarshal(respData, &likedPosts)
	if err != nil {
		t.Fatal(err)
	}

	if len(likedPosts) != len(posts) {
		t.Fatalf("invalid number of liked posts, expected %d, got %d", len(posts), len(likedPosts))
	}

	for _, post := range posts[:5] {
		cli.TestPost(t, fmt.Sprintf("/api/posts/%d/dislike", post.Id), nil, http.StatusOK)
	}

	_, respData = cli.TestGet(t, "/api/me/posts/liked", http.StatusOK)
	err = json.Unmarshal(respData, &likedPosts)
	if err != nil {
		t.Fatal(err)
	}

	if len(likedPosts) != len(posts)-5 {
		t.Fatalf("invalid number of liked posts, expected %d, got %d", len(posts)-5, len(likedPosts))
	}

	// check creator liked posts, they didn't like any posts
	_, respData = creator.TestGet(t, "/api/me/posts/liked", http.StatusOK)
	err = json.Unmarshal(respData, &likedPosts)
	if err != nil {
		t.Fatal(err)
	}

	if len(likedPosts) != 0 {
		t.Fatalf("invalid number of liked posts, expected 0, got %d", len(likedPosts))
	}
}
