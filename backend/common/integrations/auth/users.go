package auth

import (
	"net/http"
	"strconv"
)

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
}

// GetUser gets user data from the auth service. Returns nil if user not found.
func GetUser(userId int) (user *User) {
	InternalRequest(http.MethodGet, "/api/users/"+strconv.Itoa(userId), nil, &user)
	return user
}
