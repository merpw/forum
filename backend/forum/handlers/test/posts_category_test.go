package server_test

import (
	. "backend/forum/handlers/test/server"
	"encoding/json"
	"net/http"
	"testing"
)

func TestCategory(t *testing.T) {
	testServer := NewTestServer(t)

	cli := testServer.TestClient()

	cli.TestAuth(t)

	_, respBody := cli.TestGet(t, "/api/posts/categories", http.StatusOK)

	var categories []string
	err := json.Unmarshal(respBody, &categories)
	if err != nil {
		t.Fatal(err)
	}
	if len(categories) == 0 {
		t.Fatal("expected at least 1 category")
	}

	t.Run("Valid", func(t *testing.T) {
		var posts []PostData
		for _, category := range categories {
			_, respBody = cli.TestGet(t, "/api/posts/categories/"+category, http.StatusOK)
			err = json.Unmarshal(respBody, &posts)
			if err != nil {
				t.Fatal(err)
			}
			if len(posts) != 0 {
				t.Fatal("expected 0 posts")
			}
		}

		for _, category := range categories {
			for i := 0; i < 5; i++ {
				post := generatePostData()
				post.Categories = []string{category}
				cli.TestPost(t, "/api/posts/create", post, http.StatusOK)
			}
		}

		for _, category := range categories {
			_, respBody = cli.TestGet(t, "/api/posts/categories/"+category, http.StatusOK)
			err = json.Unmarshal(respBody, &posts)
			if err != nil {
				t.Fatal(err)
			}
			if len(posts) != 5 {
				t.Fatal("expected 5 posts")
			}
		}
	})

	t.Run("Invalid", func(t *testing.T) {
		cli.TestGet(t, "/api/posts/categories/invalid", http.StatusNotFound)
	})
}
