package server_test

import (
	. "backend/forum/handlers/test/server"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func TestUser(t *testing.T) {
	testServer := NewTestServer(t)

	cli1 := testServer.TestClient()
	cli1.TestAuth(t)

	cliAnon := testServer.TestClient()

	t.Run("Not found", func(t *testing.T) {
		cli1.TestGet(t, "/api/users/10", http.StatusNotFound)
		cli1.TestGet(t, "/api/users/214748364712312214748364712312", http.StatusNotFound)
	})

	t.Run("Valid", func(t *testing.T) {
		var responseData struct {
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
		t.Run("Private", func(t *testing.T) {

			_, respBody := cli1.TestGet(t, "/api/users/1", http.StatusOK)

			checkRespBody := func() {
				err := json.Unmarshal(respBody, &responseData)
				if err != nil {
					t.Fatal(err)
				}

				if responseData.Id != 1 {
					t.Fatalf("invalid id, expected 1, got %d", responseData.Id)
				}

				if responseData.Username != cli1.Username {
					t.Fatalf("invalid username, expected %s, got %s", cli1.Username, responseData.Username)
				}

				if responseData.Bio != cli1.Bio {
					t.Fatalf("invalid bio, expected %s, got %s", cli1.Bio, responseData.Bio)
				}

				if responseData.Avatar != cli1.Avatar {
					t.Fatalf("invalid avatar, expected %s, got %s", cli1.Avatar, responseData.Avatar)
				}
			}
			checkRespBody()

			_, respBody = cliAnon.TestGet(t, "/api/users/1", http.StatusOK)

			checkRespBody()

		})
		t.Run("Public", func(t *testing.T) {
			cli1.TestPost(t, "/api/me/privacy", nil, http.StatusOK)

			_, respBody := cli1.TestGet(t, "/api/users/1", http.StatusOK)

			checkRespBody := func() {
				err := json.Unmarshal(respBody, &responseData)
				if err != nil {
					t.Fatal(err)
				}

				if responseData.Id != 1 {
					t.Fatalf("invalid id, expected 1, got %d", responseData.Id)
				}

				if responseData.Username != cli1.Username {
					t.Fatalf("invalid username, expected %s, got %s", cli1.Username, responseData.Username)
				}

				if responseData.Bio != cli1.Bio {
					t.Fatalf("invalid bio, expected %s, got %s", cli1.Bio, responseData.Bio)
				}

				if responseData.Avatar != cli1.Avatar {
					t.Fatalf("invalid avatar, expected %s, got %s", cli1.Avatar, responseData.Avatar)
				}

				if responseData.Email != cli1.Email {
					t.Fatalf("invalid email, expected %s, got %s", cli1.Email, responseData.Email)
				}

				if responseData.FirstName != cli1.FirstName {
					t.Fatalf("invalid first name, expected %s, got %s", cli1.FirstName, responseData.FirstName)
				}

				if responseData.LastName != cli1.LastName {
					t.Fatalf("invalid last name, expected %s, got %s", cli1.LastName, responseData.LastName)
				}

				if responseData.DoB != cli1.DoB {
					t.Fatalf("invalid dob, expected %s, got %s", cli1.DoB, responseData.DoB)
				}

				if responseData.Gender != cli1.Gender {
					t.Fatalf("invalid gender, expected %s, got %s", cli1.Gender, responseData.Gender)
				}
			}
			checkRespBody()

			_, respBody = cliAnon.TestGet(t, "/api/users/1", http.StatusOK)

			checkRespBody()

		})

	})
}

func TestUserPosts(t *testing.T) {
	testServer := NewTestServer(t)

	cli1 := testServer.TestClient()
	cli1.TestAuth(t)

	cli2 := testServer.TestClient()
	cli2.TestAuth(t)

	t.Run("Not found", func(t *testing.T) {
		cli1.TestGet(t, "/api/users/10/posts", http.StatusNotFound)
		cli1.TestGet(t, "/api/users/214748364712312214748364712312/posts", http.StatusNotFound)
	})

	t.Run("Valid", func(t *testing.T) {
		createPosts(t, cli1, 5)

		_, respBody := cli1.TestGet(t, "/api/users/1/posts", http.StatusOK)

		var cli1Posts []PostData
		err := json.Unmarshal(respBody, &cli1Posts)
		if err != nil {
			t.Fatal(err)
		}

		if len(cli1Posts) != 5 {
			t.Fatalf("invalid posts count, expected 5, got %d", len(cli1Posts))
		}

		_, resp2Body := cli2.TestGet(t, "/api/users/1/posts", http.StatusOK)

		if string(respBody) != string(resp2Body) {
			t.Fatalf("responses mismatch, expected %s, got %s", string(respBody), string(resp2Body))
		}
	})
}

func TestUserFollow(t *testing.T) {
	testServer := NewTestServer(t)

	cli1 := testServer.TestClient()
	cli1.TestAuth(t)

	cli2 := testServer.TestClient()
	cli2.TestAuth(t)

	t.Run("Not found", func(t *testing.T) {
		cli1.TestPost(t, "/api/users/100/follow", nil, http.StatusNotFound)
		cli1.TestPost(t, "/api/users/214748364712312214748364712312/follow",
			nil, http.StatusNotFound)
	})

	t.Run("Valid", func(t *testing.T) {
		var followStatus int

		var responseData struct {
			Gender string `json:"gender"`
		}

		t.Run("Request to follow", func(t *testing.T) {
			_, response := cli1.TestPost(t, "/api/users/2/follow", nil, http.StatusOK)

			err := json.Unmarshal(response, &followStatus)
			if err != nil {
				t.Fatal(err)
			}
			if followStatus != 2 {
				t.Fatalf("invalid followStatus, expected %d, got %d", 2, followStatus)
			}
		})

		t.Run("Abort request to follow", func(t *testing.T) {
			_, response := cli1.TestPost(t, "/api/users/2/follow", nil, http.StatusOK)
			err := json.Unmarshal(response, &followStatus)
			if err != nil {
				t.Fatal(err)
			}
			if followStatus != 0 {
				t.Fatalf("invalid followStatus, expected %d, got %d", 0, followStatus)
			}
		})

		t.Run("Follow", func(t *testing.T) {

			t.Run("Make public", func(t *testing.T) {
				cli2.TestPost(t, "/api/me/privacy", nil, http.StatusOK)
				_, response := cli1.TestPost(t, "/api/users/2/follow", nil, http.StatusOK)

				err := json.Unmarshal(response, &followStatus)
				if err != nil {
					t.Fatal(err)
				}
				if followStatus != 1 {
					t.Fatalf("invalid followStatus, expected %d, got %d", 1, followStatus)
				}
			})

			t.Run("Make private", func(t *testing.T) {
				cli2.TestPost(t, "/api/me/privacy", nil, http.StatusOK)
				_, respBody := cli1.TestGet(t, "/api/users/2", http.StatusOK)
				err := json.Unmarshal(respBody, &responseData)
				if err != nil {
					t.Fatal(err)
				}
				if responseData.Gender != cli2.Gender {
					t.Errorf("invalid gender, expected %s, got %s", cli2.Gender, responseData.Gender)
				}

			})

		})

		t.Run("Unfollow", func(t *testing.T) {

			// Unfollow here
			_, response := cli1.TestPost(t, "/api/users/2/follow", nil, http.StatusOK)

			err := json.Unmarshal(response, &followStatus)
			if err != nil {
				t.Fatal(err)
			}

			if followStatus != 0 {
				t.Fatalf("invalid followStatus, expected %d, got %d", 0, followStatus)
			}

			// Reset the gender (part of public/followed profile).
			responseData.Gender = ""
			_, respBody := cli1.TestGet(t, "/api/users/2", http.StatusOK)
			err = json.Unmarshal(respBody, &responseData)
			if err != nil {
				t.Fatal(err)
			}

			if len(responseData.Gender) != 0 {
				t.Errorf("invalid gender, expected %s, got %s", "", responseData.Gender)
			}
		})
	})

}

const UsersCount = 5

func TestUsers(t *testing.T) {
	testServer := NewTestServer(t)

	for i := 0; i < UsersCount; i++ {
		// create users with random names
		cli := testServer.TestClient()
		cli.TestClientData = NewClientData()
		if i%2 == 0 {
			cli.TestClientData.Username = strings.ToUpper(cli.TestClientData.Username)
		}
		cli.TestPost(t, "/api/signup", cli.TestClientData, http.StatusOK)
	}

	cliAnon := testServer.TestClient()

	_, respBody := cliAnon.TestGet(t, "/api/users", http.StatusOK)

	var userIds []int
	err := json.Unmarshal(respBody, &userIds)
	if err != nil {
		t.Fatal(err)
	}

	if len(userIds) != UsersCount {
		t.Fatalf("invalid users count, expected %d, got %d", UsersCount, len(userIds))
	}

	type User struct {
		Id       int    `json:"id"`
		Username string `json:"username"`
	}
	var users []User

	for _, id := range userIds {
		_, respBody := cliAnon.TestGet(t, fmt.Sprintf("/api/users/%d", id), http.StatusOK)

		var responseData User

		err := json.Unmarshal(respBody, &responseData)
		if err != nil {
			t.Fatal(err)
		}

		if responseData.Id != id {
			t.Fatalf("invalid id, expected %d, got %d", id, responseData.Id)
		}

		users = append(users, responseData)
	}

	for i := 0; i < len(users); i++ {
		for j := i + 1; j < len(users); j++ {
			// should be sorted by username
			if strings.ToLower(users[i].Username) > strings.ToLower(users[j].Username) {
				t.Fatalf("users are not sorted by username, expected %s < %s", users[i].Username, users[j].Username)
			}
		}
	}
}
