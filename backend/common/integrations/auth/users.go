package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
)

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
}

// GetUser gets user data from the auth service. Returns nil if user not found.
func GetUser(userId int) (user *User) {
	req, err := http.NewRequest(http.MethodGet, os.Getenv("AUTH_BASE_URL")+"/api/users/"+strconv.Itoa(userId), nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Internal-Auth", os.Getenv("FORUM_BACKEND_SECRET"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil
	}

	if resp.StatusCode != http.StatusOK {
		log.Panicf("status code: %d", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		panic(err)
	}
	_ = resp.Body.Close()

	return user
}
