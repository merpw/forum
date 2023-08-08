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
			DoB       string `json:"dob"`
			Gender    string `json:"gender"`
			Avatar    string `json:"avatar"`
			Bio       string `json:"bio"`
			Privacy   bool   `json:"privacy"`
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

		if respData.DoB != cli.DoB {
			t.Fatalf("invalid DoB, expected %s, got %s", cli.DoB, respData.DoB)
		}

		if respData.Gender != cli.Gender {
			t.Fatalf("invalid gender, expected %s, got %s", cli.Gender, respData.Gender)
		}

		if respData.Avatar != cli.Avatar {
			t.Fatalf("invalid avatar, expected %s, got %s", cli.Avatar, respData.Avatar)
		}

		if respData.Bio != cli.Bio {
			t.Fatalf("invalid bio, expected %s, got %s", cli.Bio, respData.Bio)
		}

		if respData.Privacy != true {
			t.Fatalf("invalid privacy, expected %t, got %t", true, respData.Privacy)
		}
	})
}

func TestMePrivacy(t *testing.T) {
	testServer := NewTestServer(t)

	t.Run("Unauthorized", func(t *testing.T) {
		cli := testServer.TestClient()
		cli.TestPost(t, "/api/me/privacy", nil, http.StatusUnauthorized)
	})

	t.Run("Valid", func(t *testing.T) {
		cli := testServer.TestClient()
		cli.TestAuth(t)

		var meRespData struct {
			Privacy bool `json:"privacy"`
		}

		checkPrivacy := func(expected bool) {
			t.Helper()
			_, respBody := cli.TestGet(t, "/api/me", http.StatusOK)
			err := json.Unmarshal(respBody, &meRespData)
			if err != nil {
				t.Fatal(err)
			}
			if meRespData.Privacy != expected {
				t.Fatalf("invalid privacy, expected %t, got %t", expected, meRespData.Privacy)
			}
		}

		checkPrivacy(true)

		t.Run("Private to public", func(t *testing.T) {
			var privacy bool
			_, response := cli.TestPost(t, "/api/me/privacy", nil, http.StatusOK)
			err := json.Unmarshal(response, &privacy)
			if err != nil {
				t.Fatal(err)
			}
			if privacy != false {
				t.Fatalf("invalid privacy, expected %t, got %t", false, privacy)
			}

			checkPrivacy(false)
		})

		t.Run("Public to private", func(t *testing.T) {
			var privacy bool
			_, response := cli.TestPost(t, "/api/me/privacy", nil, http.StatusOK)
			err := json.Unmarshal(response, &privacy)
			if err != nil {
				t.Fatal(err)
			}
			if privacy != true {
				t.Fatalf("invalid privacy, expected %t, got %t", true, privacy)
			}

			checkPrivacy(true)
		})
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

func TestMeFollowers(t *testing.T) {
	testServer := NewTestServer(t)
	cli1 := testServer.TestClient()

	t.Run("Unauthorized", func(t *testing.T) {
		cli1.TestGet(t, "/api/me/followers", http.StatusUnauthorized)
	})

	cli2 := testServer.TestClient()
	cli1.TestAuth(t)
	cli2.TestAuth(t)

	t.Run("Valid", func(t *testing.T) {

		t.Run("No followers", func(t *testing.T) {
			_, respData := cli1.TestGet(t, "/api/me/followers", http.StatusOK)
			var followers []int
			err := json.Unmarshal(respData, &followers)
			if err != nil {
				t.Fatal(err)
			}

			if len(followers) != 0 {
				t.Errorf("invalid followers, expected %d, got %d", 0, len(followers))
			}

			_, respData = cli1.TestGet(t, "/api/me/following", http.StatusOK)
			var following []int
			err = json.Unmarshal(respData, &following)
			if err != nil {
				t.Fatal(err)
			}

			if len(following) != 0 {
				t.Errorf("invalid following, expected %d, got %d", 0, len(following))
			}
		})

		t.Run("With followers", func(t *testing.T) {
			cli1.TestPost(t, "/api/me/privacy", nil, http.StatusOK)
			cli2.TestPost(t, "/api/users/1/follow", nil, http.StatusOK)

			_, respData := cli1.TestGet(t, "/api/me/followers", http.StatusOK)
			var followers []int
			err := json.Unmarshal(respData, &followers)
			if err != nil {
				t.Fatal(err)
			}
			if len(followers) != 1 {
				t.Errorf("invalid followers, expected %d, got %d", 1, len(followers))
			}

			cli2.TestPost(t, "/api/me/privacy", nil, http.StatusOK)
			cli1.TestPost(t, "/api/users/2/follow", nil, http.StatusOK)

			_, respData = cli1.TestGet(t, "/api/me/following", http.StatusOK)
			var following []int
			err = json.Unmarshal(respData, &following)
			if err != nil {
				t.Fatal(err)
			}

			if len(following) != 1 {
				t.Errorf("invalid following, expected %d, got %d", 1, len(following))
			}
		})
	})
}
