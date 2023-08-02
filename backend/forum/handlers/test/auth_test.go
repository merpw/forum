package server_test

import (
	. "backend/forum/handlers/test/server"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestSignUp(t *testing.T) {
	testServer := NewTestServer(t)

	t.Run("Valid", func(t *testing.T) {
		cli := testServer.TestClient()
		client := NewClientData()
		cli.TestPost(t, "/api/signup", client, http.StatusOK)

		client.Email = "another@email.com"
		client.Username = ""
		client.Bio = ""
		cli.TestPost(t, "/api/signup", client, http.StatusOK)

	})

	t.Run("Invalid", func(t *testing.T) {
		cli := testServer.TestClient()

		t.Run("Method", func(t *testing.T) {
			cli.TestGet(t, "/api/signup", http.StatusMethodNotAllowed)
		})

		t.Run("Body", func(t *testing.T) {
			cli.TestPost(t, "/api/signup", nil, http.StatusBadRequest)
			cli.TestPost(t, "/api/signup", "cat", http.StatusBadRequest)
			cli.TestPost(t, "/api/signup", struct {
				Login string `json:"login"`
			}{
				Login: "cat",
			}, http.StatusBadRequest)
		})

		t.Run("Username", func(t *testing.T) {
			badClientData := NewClientData()

			badClientData.Username = "u123"
			cli.TestPost(t, "/api/signup", badClientData, http.StatusBadRequest)

			badClientData.Username = "E"
			cli.TestPost(t, "/api/signup", badClientData, http.StatusBadRequest)

			badClientData.Username = "THISUSERNAMEISTOOLONG"
			cli.TestPost(t, "/api/signup", badClientData, http.StatusBadRequest)

			badClientData.Username = "with spaces"
			cli.TestPost(t, "/api/signup", badClientData, http.StatusBadRequest)

		})

		t.Run("Email", func(t *testing.T) {
			badClientData := NewClientData()

			badClientData.Email = ""
			cli.TestPost(t, "/api/signup", badClientData, http.StatusBadRequest)

			badClientData.Email = "bad@@format.com"
			cli.TestPost(t, "/api/signup", badClientData, http.StatusBadRequest)

			badClientData.Email = "cat@" + strings.Repeat("a", 1000) + ".com"

		})

		t.Run("Password", func(t *testing.T) {
			badClientData := NewClientData()

			badClientData.Password = ""
			cli.TestPost(t, "/api/signup", badClientData, http.StatusBadRequest)

			badClientData.Password = "short"
			cli.TestPost(t, "/api/signup", badClientData, http.StatusBadRequest)

			badClientData.Password = strings.Repeat("a", 1000)
			cli.TestPost(t, "/api/signup", badClientData, http.StatusBadRequest)
		})

		t.Run("First name", func(t *testing.T) {
			badClientData := NewClientData()

			badClientData.FirstName = ""
			cli.TestPost(t, "/api/signup", badClientData, http.StatusBadRequest)

			badClientData.FirstName = strings.Repeat("a", 1000)
			cli.TestPost(t, "/api/signup", badClientData, http.StatusBadRequest)
		})

		t.Run("Last name", func(t *testing.T) {
			badClientData := NewClientData()

			badClientData.LastName = ""
			cli.TestPost(t, "/api/signup", badClientData, http.StatusBadRequest)

			badClientData.LastName = strings.Repeat("a", 1000)
			cli.TestPost(t, "/api/signup", badClientData, http.StatusBadRequest)
		})

		t.Run("DoB", func(t *testing.T) {
			badClientData := NewClientData()

			badClientData.DoB = ""
			cli.TestPost(t, "/api/signup", badClientData, http.StatusBadRequest)

			badClientData.DoB = "bad format"
			cli.TestPost(t, "/api/signup", badClientData, http.StatusBadRequest)

			badClientData.DoB = time.Now().AddDate(0, 0, 10).Format("2006-01-02")
			cli.TestPost(t, "/api/signup", badClientData, http.StatusBadRequest)

			badClientData.DoB = time.Now().AddDate(-200, 0, 0).Format("2006-01-02")
			cli.TestPost(t, "/api/signup", badClientData, http.StatusBadRequest)
		})

		t.Run("Gender", func(t *testing.T) {
			badClientData := NewClientData()

			badClientData.Gender = ""
			cli.TestPost(t, "/api/signup", badClientData, http.StatusBadRequest)

			badClientData.Gender = "AH-64 Apache attack helicopter"
			cli.TestPost(t, "/api/signup", badClientData, http.StatusBadRequest)
		})

		t.Run("Bio", func(t *testing.T) {
			badClientData := NewClientData()

			badClientData.Bio = strings.Repeat("a", 1000)
			cli.TestPost(t, "/api/signup", badClientData, http.StatusBadRequest)
		})

		t.Run("Avatar", func(t *testing.T) {
			badClientData := NewClientData()

			badClientData.Bio = strings.Repeat("a", 1000)
			cli.TestPost(t, "/api/signup", badClientData, http.StatusBadRequest)
		})

		t.Run("User already exists", func(t *testing.T) {
			cli := testServer.TestClient()

			clientData := NewClientData()
			cli.TestPost(t, "/api/signup", clientData, http.StatusOK)

			cli.TestPost(t, "/api/signup", clientData, http.StatusBadRequest)

			clientData.Username = "not" + clientData.Username
			cli.TestPost(t, "/api/signup", clientData, http.StatusBadRequest)
			clientData.Username = clientData.Username[3:]

			clientData.Email = "not" + clientData.Email
			cli.TestPost(t, "/api/signup", clientData, http.StatusBadRequest)
		})

		t.Run("Already logged in", func(t *testing.T) {
			cli := testServer.TestClient()
			cli.TestAuth(t)

			cli.TestPost(t, "/api/signup", NewClientData(), http.StatusBadRequest)
		})
	})
}

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func TestLogin(t *testing.T) {
	testServer := NewTestServer(t)

	t.Run("Valid", func(t *testing.T) {
		cli := testServer.TestClient()

		cli.TestAuth(t)
	})

	t.Run("Already logged in", func(t *testing.T) {
		cli := testServer.TestClient()
		cli.TestAuth(t)

		cli.TestPost(t, "/api/login", LoginRequest{
			Login:    cli.Email,
			Password: cli.Password,
		}, http.StatusBadRequest)
	})

	t.Run("Invalid", func(t *testing.T) {
		cli := testServer.TestClient()

		clientData := NewClientData()
		cli.TestPost(t, "/api/signup", clientData, http.StatusOK)

		t.Run("Body", func(t *testing.T) {
			cli.TestPost(t, "/api/login", nil, http.StatusBadRequest)
			cli.TestPost(t, "/api/login", "invalid body", http.StatusBadRequest)

			cli.TestPost(t, "/api/login", LoginRequest{
				Login: clientData.Email,
			}, http.StatusBadRequest)

			cli.TestPost(t, "/api/login", LoginRequest{
				Login: clientData.Username,
			}, http.StatusBadRequest)

			cli.TestPost(t, "/api/login", LoginRequest{
				Password: clientData.Password,
			}, http.StatusBadRequest)
		})

		t.Run("Credentials", func(t *testing.T) {
			cli := testServer.TestClient()

			cli.TestPost(t, "/api/login", LoginRequest{
				Login:    clientData.Email,
				Password: "wrong password",
			}, http.StatusBadRequest)

			cli.TestPost(t, "/api/login", LoginRequest{
				Login:    "nosuch@email.com",
				Password: clientData.Password,
			}, http.StatusBadRequest)
		})
	})
}

func TestLogout(t *testing.T) {
	testServer := NewTestServer(t)

	cli := testServer.TestClient()

	// cli.TestPost(t, "/api/logout", nil, http.StatusUnauthorized)

	cli.TestAuth(t)

	cli.TestPost(t, "/api/logout", nil, http.StatusOK)
	cli.Cookies = nil

	cli.TestPost(t, "/api/logout", nil, http.StatusUnauthorized)
}
