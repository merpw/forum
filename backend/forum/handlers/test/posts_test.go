package server_test

import (
	. "backend/forum/handlers/test/server"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"testing"
)

func TestPosts(t *testing.T) {
	testServer := NewTestServer(t)

	t.Run("Initial", func(t *testing.T) {
		cli := testServer.TestClient()

		_, respBody := cli.TestGet(t, "/api/posts", http.StatusOK)
		var respData []PostData
		err := json.Unmarshal(respBody, &respData)
		if err != nil {
			t.Fatal(err)
		}

		if len(respData) != 0 {
			t.Fatalf("expected 0, got %d", len(respData))
		}

		t.Run("[GET]/api/posts/{id}_404", func(t *testing.T) {
			cli.TestGet(t, "/api/posts/not_found", http.StatusNotFound)
			cli.TestGet(t, "/api/posts/1", http.StatusNotFound)
			cli.TestGet(t, "/api/posts/214748364712312214748364712312", http.StatusNotFound)
			cli.TestGet(t, "/api/posts/-1.5", http.StatusNotFound)
			cli.TestGet(t, "/api/posts/1000", http.StatusNotFound)
		})
	})

	var posts []PostData
	var lock sync.Mutex

	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		t.Run(fmt.Sprintf("create_10_posts/%d", i), func(t *testing.T) {
			t.Parallel()

			cli := testServer.TestClient()
			cli.TestAuth(t)

			lock.Lock()
			posts = append(posts, createPosts(t, cli, 10)...)
			lock.Unlock()

			wg.Done()
		})
	}

	t.Run("[GET]/api/posts", func(t *testing.T) {
		t.Parallel()
		wg.Wait()

		cli := testServer.TestClient()
		_, respBody := cli.TestGet(t, "/api/posts", http.StatusOK)
		var respData []PostData

		err := json.Unmarshal(respBody, &respData)
		if err != nil {
			t.Fatal(err)
		}

		if len(respData) != len(posts) {
			t.Fatalf("expected %d, got %d", len(posts), len(respData))
		}

	checkPosts:
		for _, postData := range posts {
			for _, respPostData := range respData {
				if postData.Id == respPostData.Id && postData.Title == respPostData.Title {
					continue checkPosts
				}
			}

			t.Fatalf("post with id %d not found", postData.Id)
		}
	})

	t.Run("[GET]/api/posts/{id}", func(t *testing.T) {
		t.Parallel()
		wg.Wait()

		cli := testServer.TestClient()

		for _, postData := range posts {
			_, respBody := cli.TestGet(t, fmt.Sprintf(`/api/posts/%d`, postData.Id), http.StatusOK)
			var respData PostData

			err := json.Unmarshal(respBody, &respData)
			if err != nil {
				t.Fatal(err)
			}

			if postData.Id != respData.Id {
				t.Fatalf("expected %d, got %d", postData.Id, respData.Id)
			}
		}
	})
}
