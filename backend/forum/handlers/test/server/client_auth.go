package server

import (
	"math/rand"
	"net/http"
	"strconv"
	"testing"

	"github.com/gofrs/uuid"
)

type TestClientData struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	DoB       string `json:"dob"`
	Gender    string `json:"gender"`
	Avatar    string `json:"avatar"`
	Bio       string `json:"bio"`
	Privacy   int    `json:"privacy"`
}

// NewClientData returns random TestClientData
func NewClientData() TestClientData {
	genders := []string{"male", "female", "other"}

	return TestClientData{
		Username:  "t" + uuid.Must(uuid.NewV4()).String()[0:8],
		Email:     "t" + uuid.Must(uuid.NewV4()).String()[0:8] + "@test.com",
		Password:  uuid.Must(uuid.NewV4()).String()[0:8],
		FirstName: "John",
		LastName:  "Doe",
		DoB:       "2000-01-01",
		Gender:    genders[rand.Intn(3)], //nolint:gosec // it's ok for tests
		Avatar:    strconv.Itoa(rand.Intn(10)) + ".jpg",
		Bio:       "t" + uuid.Must(uuid.NewV4()).String()[0:8],
		Privacy:   1,
	}
}

// TestAuth registers and logg-in TestClient with random TestClientData
func (cli *TestClient) TestAuth(t testing.TB) {
	t.Helper()
	cli.TestClientData = NewClientData()

	cli.TestPost(t, "/api/signup", cli.TestClientData, http.StatusOK)

	var loginRequest = struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}{
		Login:    cli.Email,
		Password: cli.Password,
	}

	cli.TestPost(t, "/api/login", loginRequest, http.StatusOK)

	if len(cli.Cookies) == 0 {
		t.Fatal("no Cookies after login")
	}
}
