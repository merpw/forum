package server

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// upper level handler for url starts from "/api/user"
func (srv *Server) apiUserHandler(w http.ResponseWriter, r *http.Request) {
	switch {
	case reApiUser.MatchString(r.URL.Path):
		srv.userHandler(w, r)
	case reApiUserLikedPosts.MatchString(r.URL.Path):
		srv.userLikedPostsHandler(w, r)
	case reApiUserId.MatchString(r.URL.Path):
		srv.userIdHandler(w, r)
	case reApiUserIdPosts.MatchString(r.URL.Path):
		srv.userIdPostsHandler(w, r)
	}
}

// func (srv *Server) usersHandler(w http.ResponseWriter, r *http.Request) {
// 	url := strings.TrimRight(r.URL.String(), "/")
// 	if url == "/api/user" {
// 		r.Cookies()
// 		// todo database stuff for "user info" + Error handling during managing data
// 		sendObject(w, "User's info")
// 		return
// 	}
// 	commands := strings.Split(strings.TrimPrefix(url, "/api/user/"), "/")
// 	userId, err := strconv.Atoi(commands[0])
// 	lenOfCommands := len(commands)
// 	if err != nil {
// 		errorResponse(w, http.StatusBadRequest)
// 		return
// 	}
// 	if lenOfCommands == 2 && commands[1] != "posts" || lenOfCommands > 2 {
// 		errorResponse(w, http.StatusBadRequest)
// 		return
// 	}
// 	switch lenOfCommands {
// 	case 1:
// 		srv.userHandler(w, r, userId, url)
// 	case 2:
// 		srv.userPostsByIdHandler(w, r, userId, url)
// 	}
// }

// Handling own user's info
func (srv *Server) userHandler(w http.ResponseWriter, r *http.Request) {
	r.Cookies()
	// todo database stuff for "own user's info" + Error handling during managing data
	sendObject(w, "show own user info")
}

func (srv *Server) userLikedPostsHandler(w http.ResponseWriter, r *http.Request) {
	r.Cookies()
	// todo database stuff
	sendObject(w, "show own liked posts")
}

func (srv *Server) userIdHandler(w http.ResponseWriter, r *http.Request) {
	userId, _ := strconv.Atoi(strings.Split(r.URL.Path, "/")[3])
	// todo database stuff for "posts by user's Id" + Error handling during managing data
	sendObject(w, fmt.Sprintf("show info about user with Id %v", userId))
}

func (srv *Server) userIdPostsHandler(w http.ResponseWriter, r *http.Request) {
	userId, _ := strconv.Atoi(strings.Split(r.URL.Path, "/")[3])
	// todo database stuff for "posts by user's Id" + Error handling during managing data
	sendObject(w, fmt.Sprintf("show posts by user with Id %v", userId))
}
