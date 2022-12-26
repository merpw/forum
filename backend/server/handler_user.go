package server

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func (srv *Server) usersHandler(w http.ResponseWriter, r *http.Request) {

	url := strings.TrimRight(r.URL.String(), "/")
	if url == "/api/user" {
		r.Cookies()

		// todo database stuff for "user info" + Error handling during managing data

		sendObject(w, "User's info")
		return
	}

	commands := strings.Split(strings.TrimPrefix(url, "/api/user/"), "/")
	userId, err := strconv.Atoi(commands[0])
	lenOfCommands := len(commands)

	if err != nil {
		errorResponse(w, http.StatusBadRequest)
		return
	}

	if lenOfCommands == 2 && commands[1] != "posts" || lenOfCommands > 2 {
		errorResponse(w, http.StatusBadRequest)
		return
	}

	switch lenOfCommands {
	case 1:
		srv.userHandler(w, r, userId, url)
	case 2:
		srv.userPostsByIdHandler(w, r, userId, url)
	}

}

// Handling own user's info
func (srv *Server) userHandler(w http.ResponseWriter, r *http.Request, userId int, pathToCheck string) {

	// todo database stuff for "own user's info" + Error handling during managing data

	sendObject(w, fmt.Sprintf("show own user %v info", userId))
}

// Handling user's posts by Id
func (srv *Server) userPostsByIdHandler(w http.ResponseWriter, r *http.Request, userId int, pathToCheck string) {

	// todo database stuff for "posts by user's Id" + Error handling during managing data

	sendObject(w, fmt.Sprintf("show posts by user's Id %v", userId))
}
