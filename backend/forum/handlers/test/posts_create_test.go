package server_test

import (
	. "backend/forum/database"
	. "backend/forum/handlers/test/server"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
)

func TestPostsCreate(t *testing.T) {
	testServer := NewTestServer(t)
	cli := testServer.TestClient()

	t.Run("Unauthorized", func(t *testing.T) {
		cli.TestPost(t, "/api/posts/create", nil, http.StatusUnauthorized)
	})

	cli.TestAuth(t)

	t.Run("Valid", func(t *testing.T) {

		t.Run("Good", func(t *testing.T) {
			goodPostData := generatePostData()

			cli.TestPost(t, "/api/posts/create", goodPostData, http.StatusOK)
		})

		t.Run("Bad description", func(t *testing.T) {
			checkDescription := func(postData PostData) {
				t.Helper()
				_, postIdBody := cli.TestPost(t, "/api/posts/create", postData, http.StatusOK)
				postId, err := strconv.Atoi(string(postIdBody))
				if err != nil {
					t.Fatal(err)
				}

				_, respBody := cli.TestGet(t, fmt.Sprintf("/api/posts/%d", postId), http.StatusOK)
				var post PostData
				err = json.Unmarshal(respBody, &post)
				if err != nil {
					t.Fatal(err)
				}

				if post.Description == "" {
					t.Fatal("description is empty")
				}

				if len(post.Description) > 200 {
					t.Fatalf("description is too long: %d, expected %d", len(post.Description), 200)
				}
			}

			badPostData := generatePostData()

			badPostData.Description = ""
			checkDescription(badPostData)

			badPostData.Description = strings.Repeat("CROP ME, I'm too long", 100)
			checkDescription(badPostData)

			badPostData.Content = strings.Repeat("CROP ME to gen description", 100)
			badPostData.Description = ""
			checkDescription(badPostData)
		})
	})

	t.Run("Invalid", func(t *testing.T) {
		t.Run("Method", func(t *testing.T) {
			cli.TestGet(t, "/api/posts/create", http.StatusMethodNotAllowed)
		})

		t.Run("Body", func(t *testing.T) {
			cli.TestPost(t, "/api/posts/create", nil, http.StatusBadRequest)
			cli.TestPost(t, "/api/posts/create", "Post", http.StatusBadRequest)
			cli.TestPost(t, "/api/posts/create", struct {
				Broken bool `json:"broken"`
			}{
				Broken: true,
			}, http.StatusBadRequest)
		})

		t.Run("Title", func(t *testing.T) {
			badPostData := generatePostData()

			badPostData.Title = ""
			cli.TestPost(t, "/api/posts/create", badPostData, http.StatusBadRequest)

			badPostData.Title = strings.Repeat("spam", 100)
			cli.TestPost(t, "/api/posts/create", badPostData, http.StatusBadRequest)
		})

		t.Run("Content", func(t *testing.T) {
			badPostData := generatePostData()

			badPostData.Content = ""
			cli.TestPost(t, "/api/posts/create", badPostData, http.StatusBadRequest)

			badPostData.Content = strings.Repeat("SPAM, too long", 1000)
			cli.TestPost(t, "/api/posts/create", badPostData, http.StatusBadRequest)
		})

		t.Run("Categories", func(t *testing.T) {
			badPostData := generatePostData()

			badPostData.Categories = []string{}
			cli.TestPost(t, "/api/posts/create", badPostData, http.StatusBadRequest)

			badPostData.Categories = []string{"spamCategoryThatDoesNotExist"}
			cli.TestPost(t, "/api/posts/create", badPostData, http.StatusBadRequest)
		})

		t.Run("Group", func(t *testing.T) {
			badPostData := generatePostData()

			var groupId = new(int)
			*groupId = 1
			badPostData.GroupId = groupId
			cli.TestPost(t, "/api/posts/create", badPostData, http.StatusBadRequest)

		})
	})
}

func BenchmarkCreatePosts(b *testing.B) {
	testServer := NewTestServer(b)
	for i := 0; i < b.N; i++ {
		b.Run(strconv.Itoa(i), func(b *testing.B) {
			cli := testServer.TestClient()
			cli.TestAuth(b)

			posts := createPosts(b, cli, 10)

			for _, postData := range posts {
				cli.TestGet(b, fmt.Sprintf(`/api/posts/%d`, postData.Id), http.StatusOK)
			}
		})
	}
}

type PostData struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Content     string `json:"content"`
	Audience    []int  `json:"audience"`
	GroupId     *int   `json:"group_id"`
	Privacy     int    `json:"privacy"`

	// []string for requests, string for responses
	Categories interface{} `json:"categories"`
}

func generatePostData() PostData {
	return PostData{
		Title:       uuid.Must(uuid.NewV4()).String()[0:8],
		Content:     "content",
		Description: "description",
		Categories:  []string{"facts"},
		Audience:    []int{},
		GroupId:     nil,
		Privacy:     int(Public),
	}
}

func createPosts(t testing.TB, cli *TestClient, count int) (posts []PostData) {
	t.Helper()
	posts = make([]PostData, count)

	for i := range posts {
		posts[i] = generatePostData()
		_, respBody := cli.TestPost(t, "/api/posts/create", posts[i], http.StatusOK)

		postId, err := strconv.Atoi(string(respBody))
		if err != nil {
			t.Fatal(err)
		}
		posts[i].Id = postId
	}

	return posts
}
